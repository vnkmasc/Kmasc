package service

import (
	"context"
	"log"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacultyService interface {
	UpdateFaculty(ctx context.Context, idStr string, req *models.UpdateFacultyRequest) (*models.Faculty, error)
	DeleteFaculty(ctx context.Context, idStr string) error
	GetFacultyByID(ctx context.Context, id primitive.ObjectID) (*models.FacultyResponse, error)
	GetFacultyByCode(ctx context.Context, code string) (*models.Faculty, error)
	CreateFaculty(ctx context.Context, claims *utils.CustomClaims, req *models.CreateFacultyRequest) error
	GetAllFaculties(ctx context.Context, universityID primitive.ObjectID) ([]models.FacultyResponse, error)
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

func (s *facultyService) CreateFaculty(ctx context.Context, claims *utils.CustomClaims, req *models.CreateFacultyRequest) error {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return common.ErrInvalidToken
	}

	faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, universityID)
	if err != nil {
		log.Printf("Error checking existing faculty by code and university ID: %v", err)
		return err
	}
	if faculty != nil {
		return common.ErrFacultyCodeExists
	}

	now := time.Now()
	newFaculty := &models.Faculty{
		ID:           primitive.NewObjectID(),
		FacultyCode:  req.FacultyCode,
		FacultyName:  req.FacultyName,
		Description:  req.Description,
		UniversityID: universityID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.facultyRepo.Create(ctx, newFaculty); err != nil {
		log.Printf("Error creating faculty in repository: %v", err)
		return err
	}

	return nil
}

func (s *facultyService) GetFacultyByCode(ctx context.Context, code string) (*models.Faculty, error) {
	return s.facultyRepo.FindByFacultyCode(ctx, code)
}

func (s *facultyService) GetAllFaculties(ctx context.Context, universityID primitive.ObjectID) ([]models.FacultyResponse, error) {
	faculties, err := s.facultyRepo.FindAllByUniversityID(ctx, universityID)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Printf("Error loading timezone Asia/Ho_Chi_Minh: %v. Using UTC instead.", err)
		loc = time.UTC
	}

	res := make([]models.FacultyResponse, 0, len(faculties))
	for _, f := range faculties {
		res = append(res, models.FacultyResponse{
			ID:          f.ID,
			FacultyCode: f.FacultyCode,
			FacultyName: f.FacultyName,
			Description: f.Description,
			CreatedAt:   f.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
			UpdatedAt:   f.UpdatedAt.In(loc).Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
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
		Description: faculty.Description,
		CreatedAt:   faculty.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   faculty.UpdatedAt.Format(time.RFC3339),
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

	update := bson.M{
		"updated_at": time.Now(),
	}

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

	if req.Description != faculty.Description {
		update["description"] = req.Description
	}

	if len(update) == 1 && update["updated_at"] != nil {
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
