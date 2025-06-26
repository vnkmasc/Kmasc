package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.UserResponse, error)
	SearchUsers(ctx context.Context, params models.SearchUserParams) ([]models.UserResponse, int64, error)
	CreateUser(ctx context.Context, claims *utils.CustomClaims, req *models.CreateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error
	GetUsersByFacultyCode(ctx context.Context, code string) ([]models.UserResponse, error)
	GetMyProfile(ctx context.Context) (*models.UserResponse, error)
}

type userService struct {
	userRepo       repository.UserRepository
	universityRepo repository.UniversityRepository
	facultyRepo    repository.FacultyRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	universityRepo repository.UniversityRepository,
	facultyRepo repository.FacultyRepository,
) UserService {
	return &userService{
		userRepo:       userRepo,
		universityRepo: universityRepo,
		facultyRepo:    facultyRepo,
	}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.UserResponse
	for _, u := range users {
		faculty, err := s.facultyRepo.FindByID(ctx, u.FacultyID)
		if err != nil || faculty == nil {
			continue // hoặc xử lý lỗi
		}

		university, err := s.universityRepo.FindByID(ctx, faculty.UniversityID)
		if err != nil || university == nil {
			continue // hoặc xử lý lỗi
		}

		result = append(result, models.UserResponse{
			ID:             u.ID,
			StudentCode:    u.StudentCode,
			FullName:       u.FullName,
			Email:          u.Email,
			Course:         u.Course,
			Status:         u.Status,
			FacultyCode:    faculty.FacultyCode,
			FacultyName:    faculty.FacultyName,
			UniversityCode: university.UniversityCode,
			UniversityName: university.UniversityName,
		})
	}

	return result, nil
}

func (s *userService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, user.UniversityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	return &models.UserResponse{
		ID:             user.ID,
		StudentCode:    user.StudentCode,
		FullName:       user.FullName,
		Email:          user.Email,
		Course:         user.Course,
		Status:         user.Status,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
	}, nil
}

func (s *userService) SearchUsers(ctx context.Context, params models.SearchUserParams) ([]models.UserResponse, int64, error) {
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, 0, common.ErrUnauthorized
	}
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, 0, common.ErrInvalidToken
	}

	params.UniversityID = universityID

	users, total, err := s.userRepo.SearchUsers(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.UserResponse
	for _, u := range users {
		faculty, _ := s.facultyRepo.FindByID(ctx, u.FacultyID)
		university, _ := s.universityRepo.FindByID(ctx, u.UniversityID)

		resp := models.UserResponse{
			ID:             u.ID,
			StudentCode:    u.StudentCode,
			FullName:       u.FullName,
			Email:          u.Email,
			Course:         u.Course,
			Status:         u.Status,
			FacultyCode:    "",
			FacultyName:    "",
			UniversityCode: "",
			UniversityName: "",
		}

		if faculty != nil {
			resp.FacultyCode = faculty.FacultyCode
			resp.FacultyName = faculty.FacultyName
		}
		if university != nil {
			resp.UniversityCode = university.UniversityCode
			resp.UniversityName = university.UniversityName
		}

		responses = append(responses, resp)
	}

	return responses, total, nil
}

func (s *userService) CreateUser(ctx context.Context, claims *utils.CustomClaims, req *models.CreateUserRequest) (*models.UserResponse, error) {
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	exists, err := s.userRepo.ExistsByStudentCodeAndUniversityID(ctx, req.StudentCode, universityID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrStudentIDExists
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrEmailExists
	}

	university, err := s.universityRepo.FindByID(ctx, universityID)
	if err != nil || university == nil {
		return nil, common.ErrUniversityNotFound
	}

	faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, req.FacultyCode, universityID)
	if err != nil || faculty == nil {
		return nil, common.ErrFacultyNotFound
	}

	user := &models.User{
		ID:           primitive.NewObjectID(),
		StudentCode:  req.StudentCode,
		FullName:     req.FullName,
		Email:        req.Email,
		FacultyID:    faculty.ID,
		UniversityID: universityID,
		Course:       req.Course,
		Status:       0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := &models.UserResponse{
		ID:             user.ID,
		StudentCode:    user.StudentCode,
		FullName:       user.FullName,
		Email:          user.Email,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
		Course:         user.Course,
		Status:         user.Status,
	}

	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, req models.UpdateUserRequest) error {
	update := bson.M{}

	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return common.ErrUnauthorized
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		return common.ErrInvalidToken
	}

	if req.StudentCode != nil {
		studentCode := strings.TrimSpace(*req.StudentCode)
		if studentCode != "" {
			exist, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, studentCode, universityID)
			if err != nil {
				return err
			}
			if exist != nil && exist.ID != id {
				return common.ErrStudentIDExists
			}
			update["student_code"] = studentCode
		}
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email != "" {
			exist, err := s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return err
			}
			if exist != nil && exist.ID != id {
				return common.ErrEmailExists
			}
			update["email"] = email
		}
	}

	if req.FullName != nil {
		fullName := strings.TrimSpace(*req.FullName)
		if fullName != "" {
			update["full_name"] = fullName
		}
	}

	if req.Course != nil {
		course := strings.TrimSpace(*req.Course)
		if course != "" {
			update["course"] = course
		}
	}

	if req.FacultyCode != nil {
		facultyCode := strings.TrimSpace(*req.FacultyCode)
		if facultyCode != "" {
			faculty, err := s.facultyRepo.FindByCodeAndUniversityID(ctx, facultyCode, universityID)
			if err != nil {
				return err
			}
			if faculty == nil {
				return common.ErrFacultyNotFound
			}
			update["faculty_id"] = faculty.ID
		}
	}

	update["updated_at"] = time.Now()

	if len(update) == 1 {
		return errors.New("không có trường nào để cập nhật")
	}

	return s.userRepo.UpdateUser(ctx, id, update)
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
func (s *userService) GetUsersByFacultyCode(ctx context.Context, code string) ([]models.UserResponse, error) {
	faculty, err := s.facultyRepo.FindByFacultyCode(ctx, code)
	if err != nil || faculty == nil {
		return nil, fmt.Errorf("không tìm thấy khoa với mã %s", code)
	}

	users, err := s.userRepo.FindUsersByFacultyID(ctx, faculty.ID)
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, u := range users {
		university, _ := s.universityRepo.FindByID(ctx, u.UniversityID)

		resp := models.UserResponse{
			ID:             u.ID,
			StudentCode:    u.StudentCode,
			FullName:       u.FullName,
			Email:          u.Email,
			Course:         u.Course,
			Status:         u.Status,
			FacultyCode:    faculty.FacultyCode,
			FacultyName:    faculty.FacultyName,
			UniversityCode: university.UniversityCode,
			UniversityName: university.UniversityName,
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *userService) GetMyProfile(ctx context.Context) (*models.UserResponse, error) {
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, common.ErrUnauthorized
	}
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotExisted
	}

	faculty, err := s.facultyRepo.FindByID(ctx, user.FacultyID)
	if err != nil {
		return nil, err
	}
	university, err := s.universityRepo.FindByID(ctx, user.UniversityID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:             user.ID,
		StudentCode:    user.StudentCode,
		FullName:       user.FullName,
		Email:          user.Email,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
		Course:         user.Course,
		Status:         user.Status,
	}, nil
}
