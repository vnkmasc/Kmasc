package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateSampleService struct {
	repo         *repository.TemplateSampleRepo
	templateRepo repository.TemplateRepository
}

func NewTemplateSampleService(repo *repository.TemplateSampleRepo, templateRepo repository.TemplateRepository) *TemplateSampleService {
	return &TemplateSampleService{repo: repo, templateRepo: templateRepo}
}

func (s *TemplateSampleService) Create(ctx context.Context, sample *models.TemplateSample) (primitive.ObjectID, error) {
	if sample.Name == "" {
		return primitive.NilObjectID, errors.New("template name is required")
	}
	if sample.HTMLContent == "" {
		return primitive.NilObjectID, errors.New("html_content is required")
	}
	if sample.UniversityID == primitive.NilObjectID {
		return primitive.NilObjectID, errors.New("university_id is required")
	}

	now := time.Now()
	sample.CreatedAt = now
	sample.UpdatedAt = now

	// Gọi repo lưu vào MongoDB
	id, err := s.repo.Create(ctx, sample)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id, nil
}

func (s *TemplateSampleService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.TemplateSample, error) {
	if id.IsZero() {
		return nil, errors.New("invalid template sample ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TemplateSampleService) Update(ctx context.Context, sample *models.TemplateSample) error {
	if sample.ID.IsZero() {
		return errors.New("invalid template sample ID")
	}
	if sample.Name == "" {
		return errors.New("template name is required")
	}
	if sample.HTMLContent == "" {
		return errors.New("html_content is required")
	}

	// Kiểm tra xem có template nào đang khóa sử dụng sample này
	templates, err := s.templateRepo.FindByTemplateSampleID(ctx, sample.ID)
	if err != nil {
		return fmt.Errorf("failed to check related templates: %w", err)
	}
	for _, tmpl := range templates {
		if tmpl.IsLocked {
			return fmt.Errorf("template %s đã bị khóa, không thể sửa giao diện", tmpl.Name)
		}
	}

	// Update TemplateSample
	return s.repo.Update(ctx, sample)
}

func (s *TemplateSampleService) GetAllVisible(ctx context.Context, universityID primitive.ObjectID) ([]*models.TemplateSample, error) {
	return s.repo.GetAllVisible(ctx, universityID)
}
