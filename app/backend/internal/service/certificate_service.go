package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/utils"
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
	CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) (*models.CertificateResponse, error)
	GetCertificatesByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.CertificateResponse, error)
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

func (s *certificateService) CreateCertificate(ctx context.Context, claims *utils.CustomClaims, req *models.CreateCertificateRequest) (*models.CertificateResponse, error) {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, req.StudentCode, universityID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}
	if user.FacultyID.IsZero() {
		return nil, fmt.Errorf("người dùng chưa được gán khoa")
	}

	if req.IsDegree {
		if err := s.validateDegreeRequest(ctx, req, universityID); err != nil {
			return nil, err
		}

		if req.SerialNumber != "" {
			exists, err := s.certificateRepo.ExistsBySerial(ctx, universityID, req.SerialNumber, true)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, common.ErrSerialNumberExists
			}
		}

		if req.RegNo != "" {
			exists, err := s.certificateRepo.ExistsByRegNo(ctx, universityID, req.RegNo, true)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, common.ErrRegNoExists
			}
		}
	} else {
		if err := s.validateCertificateRequest(ctx, req, universityID); err != nil {
			return nil, err
		}
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	now := time.Now()

	cert := &models.Certificate{
		ID:              primitive.NewObjectID(),
		UserID:          user.ID,
		FacultyID:       user.FacultyID,
		UniversityID:    universityID,
		StudentCode:     user.StudentCode,
		IsDegree:        req.IsDegree,
		Name:            req.Name,
		CertificateType: req.CertificateType,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		IssueDate:       req.IssueDate,
		Signed:          false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.certificateRepo.CreateCertificate(ctx, cert); err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] IsDegree: %v, CertificateType: %s\n", req.IsDegree, req.CertificateType)
	if req.IsDegree {

		var newStatus int
		switch strings.TrimSpace(req.CertificateType) {
		case "Cử nhân":
			newStatus = 1
		case "Kỹ sư":
			newStatus = 2
		case "Thạc sĩ":
			newStatus = 3
		case "Tiến sĩ":
			newStatus = 4
		}

		fmt.Printf("User hiện tại: %v, newStatus cần cập nhật: %d\n", user.Status, newStatus)

		currentStatus, _ := strconv.Atoi(fmt.Sprintf("%v", user.Status))
		if newStatus != 0 && currentStatus != newStatus {
			update := bson.M{
				"status":     newStatus,
				"updated_at": time.Now(),
			}
			if err := s.userRepo.UpdateUser(ctx, user.ID, update); err != nil {
				fmt.Printf("Không thể cập nhật trạng thái sinh viên: %v\n", err)
			} else {
				user.Status = newStatus
			}
		}

	}

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		IssueDate:       cert.IssueDate.Format("02/01/2006"),
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}

func (s *certificateService) GetCertificateByStudentCodeAndNameAndUniversity(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error) {
	return s.certificateRepo.FindCertificateByStudentCodeAndName(ctx, studentCode, name, universityID)
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

func (s *certificateService) validateCertificateRequest(ctx context.Context, req *models.CreateCertificateRequest, universityID primitive.ObjectID) error {
	if req.Name == "" || req.IssueDate.IsZero() {
		return common.ErrMissingRequiredFieldsForCertificate
	}

	if !req.IsDegree {
		req.SerialNumber = ""
		req.RegNo = ""
		req.CertificateType = ""
	}

	alreadyIssued, err := s.certificateRepo.ExistsCertificateByStudentCodeAndName(ctx, req.StudentCode, universityID, req.Name)
	if err != nil {
		return err
	}
	if alreadyIssued {
		return common.ErrCertificateAlreadyExists
	}

	return nil
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

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     user.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		Path:            cert.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		IssueDate:       cert.IssueDate.Format("02/01/2006"),
		Signed:          cert.Signed,
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
	}, nil
}

func (s *certificateService) DeleteCertificate(ctx context.Context, id primitive.ObjectID) error {
	return s.certificateRepo.DeleteCertificate(ctx, id)
}

func (s *certificateService) UploadCertificateFile(
	ctx context.Context,
	certificateID primitive.ObjectID,
	fileData []byte,
	filename string,
	isDegree bool,
	certificateName string,
) (string, error) {
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

	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		StudentName:     user.FullName,
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
	}, nil
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

		resp := &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     cert.StudentCode,
			StudentName:     user.FullName,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			IssueDate:       cert.IssueDate.Format("02/01/2006"),
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		}
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
			faculty = &models.Faculty{
				FacultyCode: "N/A",
				FacultyName: "Không xác định",
			}
		}

		university, _ := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if university == nil {
			university = &models.University{
				UniversityCode: "N/A",
				UniversityName: "Không xác định",
			}
		}

		resp := &models.CertificateResponse{
			ID:              cert.ID.Hex(),
			UserID:          cert.UserID.Hex(),
			StudentCode:     user.StudentCode,
			StudentName:     user.FullName,
			CertificateType: cert.CertificateType,
			Name:            cert.Name,
			SerialNumber:    cert.SerialNumber,
			RegNo:           cert.RegNo,
			Path:            cert.Path,
			FacultyCode:     faculty.FacultyCode,
			FacultyName:     faculty.FacultyName,
			IssueDate:       cert.IssueDate.Format("02/01/2006"),
			UniversityCode:  university.UniversityCode,
			UniversityName:  university.UniversityName,
			Signed:          cert.Signed,
			CreatedAt:       cert.CreatedAt,
			UpdatedAt:       cert.UpdatedAt,
		}
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
