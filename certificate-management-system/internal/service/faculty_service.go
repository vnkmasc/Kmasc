package service

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacultyService interface {
	GetAllFaculties(ctx context.Context, universityID primitive.ObjectID) ([]*models.Faculty, error)
	UpdateFaculty(ctx context.Context, idStr string, req *models.UpdateFacultyRequest) (*models.Faculty, error)
	DeleteFaculty(ctx context.Context, idStr string) error
	GetFacultyByID(ctx context.Context, id primitive.ObjectID) (*models.FacultyResponse, error)
	GetFacultyByCode(ctx context.Context, code string) (*models.Faculty, error)
	CreateFaculty(ctx context.Context, claims *utils.CustomClaims, req *models.CreateFacultyRequest) (*models.FacultyResponse, error)
}

type facultyService struct {
	universityRepo repository.UniversityRepository
	facultyRepo    repository.FacultyRepository
}

func NewFacultyService(universityRepo repository.UniversityRepository, facultyRepo repository.FacultyRepository) FacultyService {
	return &facultyService{
		universityRepo: universityRepo,
		facultyRepo:    facultyRepo,
	}
}

func (s *facultyService) CreateFaculty(ctx context.Context, claims *utils.CustomClaims, req *models.CreateFacultyRequest) (*models.FacultyResponse, error) {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, universityID)
	if err != nil {
		return nil, err
	}
	if faculty != nil {
		return nil, common.ErrFacultyCodeExists
	}

	faculty = &models.Faculty{
		ID:           primitive.NewObjectID(),
		FacultyCode:  req.FacultyCode,
		FacultyName:  req.FacultyName,
		UniversityID: universityID,
		CreatedAt:    time.Now(),
	}

	if err := s.facultyRepo.Create(ctx, faculty); err != nil {
		return nil, err
	}

	return &models.FacultyResponse{
		ID:          faculty.ID,
		FacultyCode: faculty.FacultyCode,
		FacultyName: faculty.FacultyName,
		CreatedAt:   faculty.CreatedAt.Format(time.RFC3339),
	}, nil
}
func (s *facultyService) GetFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	return s.facultyRepo.FindByFacultyCode(ctx, code)
}

func (s *facultyService) GetAllFaculties(ctx context.Context, universityID primitive.ObjectID) ([]*models.Faculty, error) {
	return s.facultyRepo.FindAllByUniversityID(ctx, universityID)
}
func (s *facultyService) GetFacultyByID(ctx context.Context, id primitive.ObjectID) (*models.FacultyResponse, error) {
	faculty, err := s.facultyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if faculty == nil {
		return nil, mongo.ErrNoDocuments
	}

	resp := &models.FacultyResponse{
		ID:          faculty.ID,
		FacultyCode: faculty.FacultyCode,
		FacultyName: faculty.FacultyName,
		CreatedAt:   faculty.CreatedAt.Format(time.RFC3339),
	}
	return resp, nil
}

func (s *facultyService) UpdateFaculty(ctx context.Context, idStr string, req *models.UpdateFacultyRequest) (*models.Faculty, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, common.ErrInvalidUserID
	}

	faculty, err := s.facultyRepo.FindByID(ctx, id)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	update := bson.M{}

	if req.FacultyCode != "" && req.FacultyCode != faculty.FacultyCode {
		existing, _ := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, faculty.UniversityID)
		if existing != nil && existing.ID != id {
			return nil, common.ErrFacultyCodeExists
		}
		update["faculty_code"] = req.FacultyCode
	}

	if req.FacultyName != "" && req.FacultyName != faculty.FacultyName {
		update["faculty_name"] = req.FacultyName
	}

	if len(update) == 0 {
		return nil, common.ErrNoFieldsToUpdate
	}

	err = s.facultyRepo.UpdateFaculty(ctx, id, update)
	if err != nil {
		return nil, err
	}

	updatedFaculty, err := s.facultyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedFaculty, nil
}
func (s *facultyService) DeleteFaculty(ctx context.Context, idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return common.ErrInvalidUserID
	}

	faculty, err := s.facultyRepo.FindByID(ctx, id)
	if err != nil || faculty == nil {
		return common.ErrFacultyNotFound
	}

	err = s.facultyRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
