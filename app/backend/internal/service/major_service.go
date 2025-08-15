package service

import (
	"context"
	"errors"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MajorService interface {
	CreateMajor(ctx context.Context, major *models.Major) error
	GetMajorsByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.Major, error)
	DeleteMajor(ctx context.Context, id primitive.ObjectID) error
}

type majorService struct {
	majorRepo   repository.MajorRepository
	facultyRepo repository.FacultyRepository
}

func NewMajorService(
	majorRepo repository.MajorRepository,
	facultyRepo repository.FacultyRepository,
) MajorService {
	return &majorService{
		majorRepo:   majorRepo,
		facultyRepo: facultyRepo,
	}
}

func (s *majorService) CreateMajor(ctx context.Context, major *models.Major) error {
	// Kiểm tra khoa có tồn tại không
	faculty, err := s.facultyRepo.FindByID(ctx, major.FacultyID)
	if err != nil {
		return err
	}
	if faculty == nil {
		return errors.New("khoa không tồn tại")
	}

	// Kiểm tra chuyên ngành có trùng mã trong khoa không
	existingMajor, err := s.majorRepo.FindByCodeAndFacultyID(ctx, major.MajorCode, major.FacultyID)
	if err != nil {
		return err
	}
	if existingMajor != nil {
		return errors.New("chuyên ngành đã tồn tại trong khoa")
	}

	// Thêm chuyên ngành
	return s.majorRepo.Insert(ctx, major)
}

func (s *majorService) GetMajorsByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.Major, error) {
	return s.majorRepo.GetByFaculty(ctx, universityID, facultyID)
}
func (s *majorService) DeleteMajor(ctx context.Context, id primitive.ObjectID) error {
	return s.majorRepo.DeleteByID(ctx, id)
}
