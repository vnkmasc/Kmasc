package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

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

func (s *templateService) CreateTemplate(ctx context.Context, name, description string, universityID, facultyID primitive.ObjectID, originalFilename string, fileBytes []byte) (*models.DiplomaTemplate, error) {
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

	// 3. Prepare MinIO path
	sanitizedFilename := sanitizeFileName(originalFilename)
	objectPath := fmt.Sprintf("diploma_template/%s/%s/%s", university.UniversityCode, faculty.FacultyCode, sanitizedFilename)

	err = s.minioClient.UploadFile(ctx, objectPath, fileBytes, "application/pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to upload to MinIO: %v", err)
	}

	fileURL := s.minioClient.GetFileURL(objectPath)
	hash := utils.ComputeSHA256(fileBytes)

	// 4. Build and save template
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

func sanitizeFileName(name string) string {
	// Replace space with underscore
	name = strings.ReplaceAll(name, " ", "_")

	// Only keep letters, digits, dashes, underscores, dots
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_\.]`)
	return reg.ReplaceAllString(name, "")
}

func (s *templateService) GetTemplatesByFaculty(ctx context.Context, _, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	// 💡 Tra thông tin faculty để lấy university_id
	faculty, err := s.facultyRepo.FindByID(ctx, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch faculty: %v", err)
	}

	universityID := faculty.UniversityID

	// ✅ Check lại vẫn đúng cơ chế
	belongs, err := s.facultyService.CheckFacultyBelongsToUniversity(ctx, facultyID, universityID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify faculty ownership: %v", err)
	}
	if !belongs {
		return nil, errors.New("faculty does not belong to any university")
	}

	return s.templateRepo.FindByUniversityAndFaculty(ctx, universityID, facultyID)
}

func (s *templateService) UpdateTemplate(ctx context.Context, id string, updateData *models.DiplomaTemplate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateData.UpdatedAt = time.Now()
	return s.templateRepo.UpdateIfNotLocked(ctx, objectID, updateData)
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
