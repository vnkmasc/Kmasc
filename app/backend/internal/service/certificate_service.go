package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
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
	GetAllCertificates(ctx context.Context) ([]*models.CertificateResponse, error)
	GetCertificateByStudentCodeAndNameAndUniversity(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error)
	DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
	UploadCertificateFile(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, filename string, isDegree bool, certificateName string) (string, error)
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificateBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error)
	GetCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.CertificateResponse, error)
	GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error)
	CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) error
	GetSimpleCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateSimpleResponse, error)
	SearchCertificates(ctx context.Context, params models.SearchCertificateParams) ([]*models.CertificateResponse, int64, error)
}

type certificateService struct {
	certificateRepo repository.CertificateRepository
	userRepo        repository.UserRepository
	facultyRepo     repository.FacultyRepository
	universityRepo  repository.UniversityRepository
	minioClient     *database.MinioClient
}

func NewCertificateService(
	certificateRepo repository.CertificateRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	minioClient *database.MinioClient,
) CertificateService {
	return &certificateService{
		certificateRepo: certificateRepo,
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

	user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, req.StudentCode, universityID)
	if err != nil || user == nil {
		return common.ErrUserNotExisted
	}
	if user.FacultyID.IsZero() {
		return fmt.Errorf("người dùng chưa được gán khoa")
	}

	// Validate đầu vào
	if err := s.validateDegreeRequest(ctx, req, universityID); err != nil {
		return err
	}
	if err := s.checkDuplicateSerialAndRegNo(ctx, universityID, req); err != nil {
		return err
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return common.ErrFacultyNotFound
	}
	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return common.ErrUniversityNotFound
	}

	// Khởi tạo object văn bằng
	cert := models.NewCertificate(req, user, universityID)
	cert.CertHash = generateCertificateHash(cert, user, faculty, university)

	// Lưu vào Mongo
	if err := s.certificateRepo.CreateCertificate(ctx, cert); err != nil {
		return err
	}

	// Cập nhật trạng thái sinh viên nếu cần
	s.updateUserStatusIfNeeded(ctx, user, req.CertificateType)

	return nil
}

func (s *certificateService) checkDuplicateSerialAndRegNo(
	ctx context.Context,
	universityID primitive.ObjectID,
	req *models.CreateCertificateRequest,
) error {
	if req.SerialNumber != "" {
		exists, err := s.certificateRepo.ExistsBySerial(ctx, universityID, req.SerialNumber, true)
		if err != nil {
			return err
		}
		if exists {
			return common.ErrSerialNumberExists
		}
	}
	if req.RegNo != "" {
		exists, err := s.certificateRepo.ExistsByRegNo(ctx, universityID, req.RegNo, true)
		if err != nil {
			return err
		}
		if exists {
			return common.ErrRegNoExists
		}
	}
	return nil
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
	if req.CertificateType == "" || req.SerialNumber == "" || req.RegNo == "" || req.IssueDate.IsZero() {
		return common.ErrMissingRequiredFieldsForDegree
	}

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

func (s *certificateService) UploadCertificateFile(ctx context.Context, certificateID primitive.ObjectID, fileData []byte, filename string, isDegree bool, certificateName string) (string, error) {
	certificate, err := s.certificateRepo.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy certificate: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, certificate.UniversityID)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy trường đại học: %w", err)
	}

	var objectKey string
	if isDegree {
		objectKey = fmt.Sprintf("certificates/%s/diploma/%s", university.UniversityCode, filename)
	} else {
		cleanName := strings.ReplaceAll(strings.TrimSpace(certificateName), " ", "_")
		objectKey = fmt.Sprintf("certificates/%s/%s/%s", university.UniversityCode, cleanName, filename)
	}

	contentType := http.DetectContentType(fileData)

	err = s.minioClient.UploadFile(ctx, objectKey, fileData, contentType)
	if err != nil {
		return "", fmt.Errorf("lỗi upload file lên MinIO: %w", err)
	}

	err = s.certificateRepo.UpdateCertificatePath(ctx, certificateID, objectKey)
	if err != nil {
		return "", fmt.Errorf("lỗi cập nhật path vào MongoDB: %w", err)
	}
	hash := sha256.Sum256(fileData)
	certHash := hex.EncodeToString(hash[:])
	update := bson.M{
		"$set": bson.M{
			"path":       objectKey,
			"hash_file":  certHash,
			"updated_at": time.Now(),
		},
	}

	if err := s.certificateRepo.UpdateCertificateByID(ctx, certificateID, update); err != nil {
		return "", fmt.Errorf("lỗi cập nhật thông tin file vào MongoDB: %w", err)
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

	filter := bson.M{
		"university_id": universityID,
	}
	if params.StudentCode != "" {
		filter["student_code"] = bson.M{"$regex": params.StudentCode, "$options": "i"}
	}
	if params.CertificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": params.CertificateType, "$options": "i"}
	}
	if params.Signed != nil {
		filter["signed"] = *params.Signed
	}
	if params.FacultyCode != "" {
		faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, params.FacultyCode, universityID)
		if err != nil || faculty == nil {
			return nil, 0, fmt.Errorf("faculty not found in your university with code: %s", params.FacultyCode)
		}
		filter["faculty_id"] = faculty.ID
	}

	certs, total, err := s.certificateRepo.FindCertificate(ctx, filter, params.Page, params.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []*models.CertificateResponse
	for _, cert := range certs {
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil || user == nil {
			continue
		}
		if params.Course != "" && !strings.Contains(strings.ToLower(user.Course), strings.ToLower(params.Course)) {
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
