package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrTemplateLocked = errors.New("template is locked and cannot be modified")
)

type TemplateService interface {
	GetTemplateByID(ctx context.Context, id string) (*models.DiplomaTemplate, error)
	UpdateTemplate(ctx context.Context, templateID, universityID primitive.ObjectID, name, description, originalFilename string, fileBytes []byte) (*models.DiplomaTemplate, error)
	CreateTemplate(ctx context.Context, name, description string, universityID, facultyID primitive.ObjectID, originalFilename string, fileBytes []byte) (*models.DiplomaTemplate, error)
	GetTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	SignTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) (int, error)
	SignAllPendingTemplatesOfUniversity(ctx context.Context, universityID primitive.ObjectID) (int, error)
	SignAllTemplatesByMinEdu(ctx context.Context, universityID primitive.ObjectID) (int, error)
	VerifyTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) error
}

type templateService struct {
	templateRepo   repository.TemplateRepository
	facultyRepo    repository.FacultyRepository
	universityRepo repository.UniversityRepository
	facultyService FacultyService
	minioClient    *database.MinioClient
}

func NewTemplateService(
	templateRepo repository.TemplateRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	facultyService FacultyService,
	minioClient *database.MinioClient,
) TemplateService {
	return &templateService{
		templateRepo:   templateRepo,
		facultyRepo:    facultyRepo,
		universityRepo: universityRepo,
		facultyService: facultyService,
		minioClient:    minioClient,
	}
}
func (s *templateService) GetTemplateByID(ctx context.Context, id string) (*models.DiplomaTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.templateRepo.GetByID(ctx, objID)
}

func (s *templateService) VerifyTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) error {
	return s.templateRepo.VerifyTemplatesByFaculty(ctx, universityID, facultyID)
}

func (s *templateService) CreateTemplate(
	ctx context.Context,
	name, description string,
	universityID, facultyID primitive.ObjectID,
	originalFilename string,
	fileBytes []byte,
) (*models.DiplomaTemplate, error) {

	// 1. Check ownership
	belongs, err := s.facultyService.CheckFacultyBelongsToUniversity(ctx, facultyID, universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify faculty ownership: %v", err)
	}
	if !belongs {
		return nil, errors.New("faculty does not belong to your university")
	}

	// 2. Get university & faculty info
	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch university: %v", err)
	}

	faculty, err := s.facultyRepo.FindByID(ctx, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch faculty: %v", err)
	}

	// 3. Generate unique filename
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".html" // hoặc .pdf tùy định dạng mặc định
	}
	randomName := fmt.Sprintf("%s_template%s", uuid.New().String(), ext)

	objectPath := fmt.Sprintf("diploma_template/%s/%s/%s", university.UniversityCode, faculty.FacultyCode, randomName)

	// 4. Upload
	err = s.minioClient.UploadFile(ctx, objectPath, fileBytes, "application/pdf") // hoặc "text/html"
	if err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %v", err)
	}

	fileURL := s.minioClient.GetFileURL(objectPath)
	hash := utils.ComputeSHA256(fileBytes)

	// 5. Save to DB
	template := &models.DiplomaTemplate{
		ID:           primitive.NewObjectID(),
		Name:         name,
		Description:  description,
		FileLink:     fileURL,
		Hash:         hash,
		Status:       "PENDING",
		IsLocked:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UniversityID: universityID,
		FacultyID:    facultyID,
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to save template: %v", err)
	}

	return template, nil
}

// func sanitizeFileName(name string) string {
// 	// Replace space with underscore
// 	name = strings.ReplaceAll(name, " ", "_")

// 	// Only keep letters, digits, dashes, underscores, dots
// 	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_\.]`)
// 	return reg.ReplaceAllString(name, "")
// }

func (s *templateService) GetTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	// ✅ Tìm faculty theo facultyID và universityID để đảm bảo khoa thuộc đúng trường
	faculty, err := s.facultyRepo.FindByIDAndUniversityID(ctx, facultyID, universityID)
	if err != nil {
		return nil, fmt.Errorf("faculty does not belong to the university or not found: %v", err)
	}

	// ✅ Nếu tìm thấy thì truy vấn template
	return s.templateRepo.FindByUniversityAndFaculty(ctx, universityID, faculty.ID)
}

func (s *templateService) LockTemplate(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.templateRepo.LockTemplate(ctx, objectID)
}

func (s *templateService) SignTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) (int, error) {
	templates, err := s.templateRepo.FindPendingByUniversityAndFaculty(ctx, universityID, facultyID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range templates {
		signature := "SIMULATED_SIGNATURE_" + t.ID.Hex()
		status := "SIGNED_BY_UNI"

		err := s.templateRepo.UpdateStatusAndSignatureByID(ctx, t.ID, status, signature)
		if err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func (s *templateService) SignAllPendingTemplatesOfUniversity(ctx context.Context, universityID primitive.ObjectID) (int, error) {
	templates, err := s.templateRepo.FindPendingByUniversity(ctx, universityID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range templates {
		signature := "SIMULATED_SIGNATURE_" + t.ID.Hex()
		status := "SIGNED_BY_UNI"

		err := s.templateRepo.UpdateStatusAndSignatureByID(ctx, t.ID, status, signature)
		if err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func (s *templateService) SignAllTemplatesByMinEdu(ctx context.Context, universityID primitive.ObjectID) (int, error) {
	templates, err := s.templateRepo.FindSignedByUniversity(ctx, universityID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range templates {
		signature := "SIMULATED_SIGNATURE_MINEDU_" + t.ID.Hex()
		status := "SIGNED_BY_MINEDU"

		err := s.templateRepo.UpdateStatusAndMinEduSignatureByID(ctx, t.ID, status, signature)
		if err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func (s *templateService) UpdateTemplate(
	ctx context.Context,
	templateID, universityID primitive.ObjectID,
	name, description, originalFilename string,
	fileBytes []byte,
) (*models.DiplomaTemplate, error) {
	// 1. Lấy template hiện tại
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}

	// 2. Kiểm tra quyền sở hữu
	belongs, err := s.facultyService.CheckFacultyBelongsToUniversity(ctx, template.FacultyID, universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %v", err)
	}
	if !belongs {
		return nil, errors.New("you don't have permission to update this template")
	}

	// 3. Nếu đã bị khóa thì không cho sửa
	if template.IsLocked {
		return nil, errors.New("template is locked and cannot be updated")
	}

	// 4. Cập nhật tên và mô tả nếu có
	if name != "" {
		template.Name = name
	}
	if description != "" {
		template.Description = description
	}

	// 5. Nếu có file mới thì xử lý file
	if len(fileBytes) > 0 && originalFilename != "" {
		// 5.1. Lấy thông tin trường và khoa
		university, err := s.universityRepo.FindByID(ctx, universityID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch university: %v", err)
		}
		faculty, err := s.facultyRepo.FindByID(ctx, template.FacultyID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch faculty: %v", err)
		}

		// 5.2. Xoá file cũ khỏi MinIO (nếu có)
		oldPath := extractObjectPathFromURL(template.FileLink)
		if oldPath != "" {
			if err := s.minioClient.RemoveFile(ctx, oldPath); err != nil {
				// Ghi log nếu cần, nhưng không fail
				fmt.Printf("Warning: failed to remove old file from MinIO: %v\n", err)
			}
		}

		// 5.3. Tạo tên file mới ngẫu nhiên
		ext := filepath.Ext(originalFilename)
		if ext == "" {
			ext = ".html"
		}
		randomName := fmt.Sprintf("%s_template%s", uuid.New().String(), ext)

		objectPath := fmt.Sprintf("diploma_template/%s/%s/%s", university.UniversityCode, faculty.FacultyCode, randomName)

		// 5.4. Upload file mới lên MinIO
		if err := s.minioClient.UploadFile(ctx, objectPath, fileBytes, "text/html"); err != nil {
			return nil, fmt.Errorf("failed to upload to MinIO: %v", err)
		}

		// 5.5. Cập nhật thông tin file trong DB
		template.FileLink = s.minioClient.GetFileURL(objectPath)
		template.Hash = utils.ComputeSHA256(fileBytes)
	}

	// 6. Cập nhật thời gian
	template.UpdatedAt = time.Now()

	// 7. Lưu lại vào DB
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to update template: %v", err)
	}

	return template, nil
}

func extractObjectPathFromURL(url string) string {
	// Ví dụ: http://host:9000/certificates/diploma_template/...
	parts := strings.SplitN(url, "/certificates/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
