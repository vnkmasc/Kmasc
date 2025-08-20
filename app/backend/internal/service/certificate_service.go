package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CertificateService interface {
	GetDegreeCertificateByStudentCodeAndUniversity(ctx context.Context, studentCode string, universityID primitive.ObjectID) (*models.Certificate, error)
	GetCertificateByStudentCodeAndTypeAndUniversity(ctx context.Context, studentCode string, certificateType string, universityID primitive.ObjectID) (*models.Certificate, error)
	GetAllCertificates(ctx context.Context) ([]*models.CertificateResponse, error)
	GetCertificateByStudentCodeAndNameAndUniversity(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error)
	DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificateBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error)
	GetCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error)
	CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) error
	GetSimpleCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateSimpleResponse, error)
	SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error)
	GetRawCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error)
	UploadCertificateFileDirect(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, origFileName string, isDegree bool) (string, error)
}

type certificateService struct {
	certificateRepo repository.CertificateRepository
	ediplomaRepo    repository.EDiplomaRepository
	userRepo        repository.UserRepository
	facultyRepo     repository.FacultyRepository
	universityRepo  repository.UniversityRepository
	minioClient     *database.MinioClient
}

func NewCertificateService(
	certificateRepo repository.CertificateRepository,
	ediplomaRepo repository.EDiplomaRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	minioClient *database.MinioClient,
) CertificateService {
	return &certificateService{
		certificateRepo: certificateRepo,
		ediplomaRepo:    ediplomaRepo,
		userRepo:        userRepo,
		facultyRepo:     facultyRepo,
		universityRepo:  universityRepo,
		minioClient:     minioClient,
	}
}

func (s *certificateService) CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) error {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return common.ErrInvalidToken
	}

	user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, strings.TrimSpace(req.StudentCode), universityID)
	if err != nil || user == nil {
		return common.ErrUserNotExisted
	}
	if user.FacultyID.IsZero() {
		return fmt.Errorf("người dùng chưa được gán khoa")
	}

	if err := s.validateDegreeRequest(ctx, req, universityID); err != nil {
		return err
	}
	if err := s.checkDuplicateSerialAndRegNo(ctx, universityID, req); err != nil {
		return err
	}

	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return common.ErrUniversityNotFound
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return common.ErrFacultyNotFound
	}

	cert := models.NewCertificate(req, user, universityID)
	cert.CertHash = generateCertificateHash(cert, user, faculty, university)

	if err := s.certificateRepo.CreateCertificate(ctx, cert); err != nil {
		return err
	}

	ed := mapCertificateToEDiploma(cert, user)
	if err := s.ediplomaRepo.Save(ctx, ed); err != nil {
		return err
	}

	if req.IsDegree {
		s.updateUserStatusIfNeeded(ctx, user, req.CertificateType)
	}

	return nil
}

func mapCertificateToEDiploma(cert *models.Certificate, user *models.User) *models.EDiploma {
	return &models.EDiploma{
		ID:                 primitive.NewObjectID(),
		CertificateID:      cert.ID,
		Name:               cert.Name,
		UniversityID:       cert.UniversityID,
		FacultyID:          cert.FacultyID,
		UserID:             cert.UserID,
		StudentCode:        cert.StudentCode,
		FullName:           user.FullName,
		CertificateType:    cert.CertificateType,
		Course:             cert.Course,
		EducationType:      cert.EducationType,
		GPA:                cert.GPA,
		GraduationRank:     cert.GraduationRank,
		SerialNumber:       cert.SerialNumber,
		RegistrationNumber: cert.RegNo,
		Issued:             false,
		Signed:             false,
		SignedAt:           cert.SignedAt,
		DataEncrypted:      false,
		OnBlockchain:       false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func (s *certificateService) checkDuplicateSerialAndRegNo(
	ctx context.Context,
	universityID primitive.ObjectID,
	req *models.CreateCertificateRequest,
) error {
	if req.SerialNumber != "" {
		exists, err := s.certificateRepo.ExistsBySerial(ctx, universityID, req.SerialNumber, req.IsDegree)
		if err != nil {
			return err
		}
		if exists {
			return common.ErrSerialNumberExists
		}
	}
	if req.RegNo != "" {
		exists, err := s.certificateRepo.ExistsByRegNo(ctx, universityID, req.RegNo, req.IsDegree)
		if err != nil {
			return err
		}
		if exists {
			return common.ErrRegNoExists
		}
	}
	return nil
}
func (s *certificateService) GetCertificateByStudentCodeAndTypeAndUniversity(
	ctx context.Context,
	studentCode string,
	certificateType string,
	universityID primitive.ObjectID,
) (*models.Certificate, error) {
	return s.certificateRepo.FindOneByStudentCodeAndType(ctx, studentCode, certificateType, universityID)
}

func (s *certificateService) updateUserStatusIfNeeded(ctx context.Context, user *models.User, certType string) {
	var newStatus int
	switch strings.TrimSpace(certType) {
	case "Cử nhân":
		newStatus = 1
	case "Kỹ sư":
		newStatus = 2
	case "Thạc sĩ":
		newStatus = 3
	case "Tiến sĩ":
		newStatus = 4
	}
	currentStatus, _ := strconv.Atoi(fmt.Sprintf("%v", user.Status))
	if newStatus != 0 && currentStatus != newStatus {
		update := bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		}
		if err := s.userRepo.UpdateUser(ctx, user.ID, update); err == nil {
			user.Status = newStatus
		}
	}
}

func generateCertificateHash(cert *models.Certificate, user *models.User, faculty *models.Faculty, university *models.University) string {
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%f|%s",
		user.FullName,                       // họ tên
		user.DateOfBirth,                    // ngày sinh
		cert.StudentCode,                    // mã sv
		user.CitizenIdNumber,                // căn cước công dân
		user.Email,                          // email
		university.UniversityCode,           // mã trường
		faculty.FacultyCode,                 // mã khoa
		cert.Major,                          // ngành đào tạo
		cert.GPA,                            // gpa
		cert.IssueDate.Format("2006-01-02"), // ngày cấp
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *certificateService) validateDegreeRequest(ctx context.Context, req *models.CreateCertificateRequest, universityID primitive.ObjectID) error {
	// Validate bắt buộc chung (với cả văn bằng và chứng chỉ)
	if req.SerialNumber == "" || req.RegNo == "" || req.IssueDate.IsZero() || req.Name == "" {
		return common.ErrMissingRequiredFieldsForDegree
	}

	// Nếu là văn bằng thì check thêm các trường đặc thù
	if req.IsDegree {
		if req.CertificateType == "" || req.Course == "" || req.Major == "" {
			return common.ErrMissingRequiredFieldsForDegree
		}

		// Kiểm tra trùng loại văn bằng (Cử nhân, Thạc sĩ, v.v.)
		singleDegreeTypes := map[string]bool{
			"Cử nhân": true,
			"Thạc sĩ": true,
			"Tiến sĩ": true,
			"Kỹ sư":   true,
		}

		if singleDegreeTypes[req.CertificateType] {
			alreadyIssued, err := s.certificateRepo.ExistsDegreeByStudentCodeAndType(ctx, req.StudentCode, universityID, req.CertificateType)
			if err != nil {
				return err
			}
			if alreadyIssued {
				return common.ErrCertificateAlreadyExists
			}
		}
	}

	return nil
}

func (s *certificateService) GetCertificateByStudentCodeAndNameAndUniversity(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error) {
	return s.certificateRepo.FindCertificateByStudentCodeAndName(ctx, studentCode, name, universityID)
}

func (s *certificateService) GetAllCertificates(ctx context.Context) ([]*models.CertificateResponse, error) {
	certs, err := s.certificateRepo.GetAllCertificates(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.CertificateResponse, 0, len(certs))

	for _, cert := range certs {
		faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if err != nil || faculty == nil {
			faculty = &models.Faculty{
				FacultyCode: "N/A",
				FacultyName: "Không xác định",
			}
		}

		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil || university == nil {
			university = &models.University{
				UniversityCode: "N/A",
				UniversityName: "Không xác định",
			}
		}

		responses = append(responses, &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     cert.StudentCode,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		})
	}

	return responses, nil
}

// Sửa service
func (s *certificateService) GetRawCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	cert, err := s.certificateRepo.GetCertificateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrCertificateNotFound
	}
	return cert, nil
}

func (s *certificateService) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error) {
	cert, err := s.certificateRepo.GetCertificateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrCertificateNotFound
	}

	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil || faculty == nil {
		faculty = &models.Faculty{
			FacultyCode: "N/A",
			FacultyName: "Không xác định",
		}
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil || university == nil {
		university = &models.University{
			UniversityCode: "N/A",
			UniversityName: "Không xác định",
		}
	}

	return mapper.MapCertificateToResponse(cert, user, faculty, university), nil
}

func (s *certificateService) DeleteCertificate(ctx context.Context, id primitive.ObjectID) error {
	return s.certificateRepo.DeleteCertificate(ctx, id)
}
func (s *certificateService) GetDegreeCertificateByStudentCodeAndUniversity(
	ctx context.Context,
	studentCode string,
	universityID primitive.ObjectID,
) (*models.Certificate, error) {
	return s.certificateRepo.FindOneByFilter(ctx, bson.M{
		"student_code":  studentCode,
		"university_id": universityID,
		"is_degree":     true,
	})
}

func (s *certificateService) UploadCertificateFileDirect(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, origFileName string, isDegree bool) (string, error) {
	certificate, err := s.certificateRepo.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy certificate: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, certificate.UniversityID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy trường đại học: %w", err)
	}

	// Tên file lưu trên MinIO
	ext := filepath.Ext(origFileName)
	slug := "van-bang"
	if !isDegree {
		slug = utils.Slugify(certificate.Name)
	}
	filename := fmt.Sprintf("%s/%s%s", certificate.StudentCode, slug, ext)
	objectKey := fmt.Sprintf("certificates/%s/%s", university.UniversityCode, filename)

	// Upload file trực tiếp
	contentType := http.DetectContentType(fileData)
	err = s.minioClient.UploadFile(ctx, objectKey, fileData, contentType)
	if err != nil {
		return "", fmt.Errorf("lỗi upload file lên MinIO: %w", err)
	}

	// Tính SHA256 file gốc
	hash := sha256.Sum256(fileData)
	certHash := hex.EncodeToString(hash[:])

	// Cập nhật thông tin file trong MongoDB
	update := bson.M{
		"$set": bson.M{
			"path":                 objectKey,
			"hash_file":            certHash,
			"physical_copy_issued": true,
			"updated_at":           time.Now(),
		},
	}
	if err := s.certificateRepo.UpdateCertificateByID(ctx, certificateID, update); err != nil {
		return "", fmt.Errorf("lỗi cập nhật thông tin file: %w", err)
	}

	return objectKey, nil
}

func (s *certificateService) GetCertificateBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error) {
	return s.certificateRepo.FindBySerialAndUniversity(ctx, serial, universityID)
}

func (s *certificateService) GetCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.CertificateResponse, error) {
	cert, err := s.certificateRepo.FindLatestCertificateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, common.ErrCertificateNotFound
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		user = &models.User{
			FullName: "Không xác định",
		}
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil || faculty == nil {
		faculty = &models.Faculty{
			FacultyCode: "N/A",
			FacultyName: "Không xác định",
		}
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil || university == nil {
		university = &models.University{
			UniversityCode: "N/A",
			UniversityName: "Không xác định",
		}
	}

	return mapper.MapCertificateToResponse(cert, user, faculty, university), nil

}

func (s *certificateService) SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error) {
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, 0, common.ErrUnauthorized
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, 0, common.ErrInvalidToken
	}

	filter := bson.M{"university_id": universityID}
	if params.StudentCode != "" {
		filter["student_code"] = bson.M{"$regex": params.StudentCode, "$options": "i"}
	}
	if params.CertificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": params.CertificateType, "$options": "i"}
	}
	if params.Signed != nil {
		filter["signed"] = *params.Signed
	}
	if params.Course != "" {
		filter["course"] = bson.M{"$regex": params.Course, "$options": "i"}
	}

	if params.FacultyCode != "" {
		faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, params.FacultyCode, universityID)
		if err != nil || faculty == nil {
			return nil, 0, fmt.Errorf("faculty not found in your university with code: %s", params.FacultyCode)
		}
		filter["faculty_id"] = faculty.ID
	}
	if params.Year > 0 {
		from := time.Date(params.Year, 1, 1, 0, 0, 0, 0, time.UTC)
		to := from.AddDate(1, 0, 0)

		filter["issue_date"] = bson.M{
			"$gte": from,
			"$lt":  to,
		}
	}

	// Lấy toàn bộ certificates để tự nhóm và phân trang
	allCerts, _, err := s.certificateRepo.FindCertificate(ctx, filter, 0, 0)
	if err != nil {
		return nil, 0, err
	}

	// Gom theo userID
	grouped := make(map[primitive.ObjectID][]*models.Certificate)
	for _, cert := range allCerts {
		grouped[cert.UserID] = append(grouped[cert.UserID], cert)
	}

	// Sắp xếp từng nhóm theo: Văn bằng trước, chứng chỉ sau → rồi theo IssueDate
	for _, group := range grouped {
		sort.SliceStable(group, func(i, j int) bool {
			if group[i].IsDegree != group[j].IsDegree {
				return group[i].IsDegree // văn bằng trước
			}
			return group[i].IssueDate.Before(group[j].IssueDate)
		})
	}

	// Duyệt theo từng user để gom lại thứ tự mong muốn: sinh viên 1 -> sinh viên 2 -> ...
	var sortedCerts []*models.Certificate
	for userID := range grouped {
		sortedCerts = append(sortedCerts, grouped[userID]...)
	}

	// Nếu có sort tổng thể (mới -> cũ), áp dụng sau khi đã gom nhóm
	if strings.ToLower(params.SortOrder) == "desc" {
		sort.SliceStable(sortedCerts, func(i, j int) bool {
			return sortedCerts[i].CreatedAt.After(sortedCerts[j].CreatedAt)
		})
	} else if strings.ToLower(params.SortOrder) == "asc" {
		sort.SliceStable(sortedCerts, func(i, j int) bool {
			return sortedCerts[i].CreatedAt.Before(sortedCerts[j].CreatedAt)
		})
	}

	// Phân trang thủ công
	total := int64(len(sortedCerts))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if start > len(sortedCerts) {
		start = len(sortedCerts)
	}
	if end > len(sortedCerts) {
		end = len(sortedCerts)
	}
	pagedCerts := sortedCerts[start:end]

	// Map sang response
	var results []*models.CertificateResponse
	for _, cert := range pagedCerts {
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			continue
		}

		faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if err != nil || faculty == nil {
			continue
		}
		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil || university == nil {
			continue
		}
		resp := mapper.MapCertificateToResponse(cert, user, faculty, university)
		results = append(results, resp)
	}

	return results, total, nil
}

func (s *certificateService) GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error) {
	certs, err := s.certificateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.CertificateResponse
	for _, cert := range certs {
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			continue
		}

		faculty, _ := s.facultyRepo.FindByID(ctx, cert.FacultyID)
		if faculty == nil {
			faculty = &models.Faculty{FacultyCode: "N/A", FacultyName: "Không xác định"}
		}

		university, _ := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if university == nil {
			university = &models.University{UniversityCode: "N/A", UniversityName: "Không xác định"}
		}

		resp := mapper.MapCertificateToResponse(cert, user, faculty, university)
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *certificateService) DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error {
	return s.certificateRepo.DeleteCertificateByID(ctx, id)
}

func (s *certificateService) GetSimpleCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateSimpleResponse, error) {
	certs, err := s.certificateRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []*models.CertificateSimpleResponse
	for _, cert := range certs {
		responses = append(responses, &models.CertificateSimpleResponse{
			ID:   cert.ID.Hex(),
			Name: cert.Name,
		})
	}

	return responses, nil
}
