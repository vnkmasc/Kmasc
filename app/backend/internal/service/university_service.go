package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniversityService interface {
	CreateUniversity(ctx context.Context, req *models.CreateUniversityRequest) error
	ApproveOrRejectUniversity(ctx context.Context, idStr string, action string) error
	GetAllUniversities(ctx context.Context) ([]models.UniversityResponse, error)
	GetUniversitiesByStatus(ctx context.Context, status string) ([]*models.University, error)
	GetUniversityByID(ctx context.Context, id primitive.ObjectID) (*models.University, error)
	GetUniversityByCode(ctx context.Context, code string) (*models.University, error)
}

type universityService struct {
	universityRepo repository.UniversityRepository
	authRepo       repository.AuthRepository
	emailSender    utils.EmailSender
}

func NewUniversityService(
	universityRepo repository.UniversityRepository,
	authRepo repository.AuthRepository,
	emailSender utils.EmailSender,
) UniversityService {
	return &universityService{
		universityRepo: universityRepo,
		authRepo:       authRepo,
		emailSender:    emailSender,
	}
}

func (s *universityService) GetUniversityByID(ctx context.Context, id primitive.ObjectID) (*models.University, error) {
	university, err := s.universityRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if university == nil {
		return nil, fmt.Errorf("không tìm thấy trường đại học")
	}
	return university, nil
}

func (s *universityService) CreateUniversity(ctx context.Context, req *models.CreateUniversityRequest) error {
	conflictField, err := s.universityRepo.CheckUniversityConflicts(ctx, req.UniversityName, req.EmailDomain, req.UniversityCode)
	if err != nil {
		return err
	}
	switch conflictField {
	case "university_name":
		return common.ErrUniversityNameExists
	case "email_domain":
		return common.ErrUniversityEmailDomainExists
	case "university_code":
		return common.ErrUniversityCodeExists
	}

	uni := &models.University{
		ID:             primitive.NewObjectID(),
		UniversityName: req.UniversityName,
		Address:        req.Address,
		EmailDomain:    req.EmailDomain,
		UniversityCode: req.UniversityCode,
		Description:    req.Description,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.universityRepo.CreateUniversity(ctx, uni)
}

func (s *universityService) ApproveOrRejectUniversity(ctx context.Context, idStr string, action string) error {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return common.ErrUniversityNotFound
	}

	university, err := s.universityRepo.FindByID(ctx, objID)
	if err != nil || university == nil {
		return common.ErrUniversityNotFound
	}

	switch action {
	case "approve":
		if university.Status == "approved" {
			return common.ErrUniversityAlreadyApproved
		}
		if err := s.universityRepo.UpdateStatus(ctx, objID, "approved"); err != nil {
			return err
		}
		existingAccount, _ := s.authRepo.FindByPersonalEmail(ctx, university.EmailDomain)
		if existingAccount != nil {
			return common.ErrAccountUniversityAlreadyExists
		}

		rawPassword := utils.GenerateRandomPassword(10)
		hashed, err := utils.HashPassword(rawPassword)
		if err != nil {
			return err
		}

		account := &models.Account{
			ID:            primitive.NewObjectID(),
			UniversityID:  university.ID,
			PersonalEmail: university.EmailDomain,
			PasswordHash:  hashed,
			CreatedAt:     time.Now(),
			Role:          "university_admin",
		}
		fmt.Println("University ID:", university.ID.Hex())

		if err := s.authRepo.CreateAccount(ctx, account); err != nil {
			return err
		}

		emailBody := fmt.Sprintf(`Xin chào,

Trường %s đã được phê duyệt truy cập hệ thống.

Thông tin tài khoản:
- Email đăng nhập: %s
- Mật khẩu: %s

Vui lòng đăng nhập và thay đổi mật khẩu ngay sau lần đầu sử dụng.

Trân trọng.`, university.UniversityName, account.PersonalEmail, rawPassword)

		_ = s.emailSender.SendEmail(account.PersonalEmail, "Tài khoản quản trị trường", emailBody)
		return nil

	case "reject":
		return s.universityRepo.DeleteByID(ctx, objID)

	default:
		return errors.New("invalid action")
	}
}

func (s *universityService) GetAllUniversities(ctx context.Context) ([]models.UniversityResponse, error) {
	universities, err := s.universityRepo.GetAllUniversities(ctx)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Printf("Error loading timezone Asia/Ho_Chi_Minh: %v. Using UTC instead.", err)
		loc = time.UTC
	}

	res := make([]models.UniversityResponse, 0, len(universities))
	for _, u := range universities {
		res = append(res, models.UniversityResponse{
			ID:             u.ID.Hex(),
			UniversityName: u.UniversityName,
			UniversityCode: u.UniversityCode,
			EmailDomain:    u.EmailDomain,
			Address:        u.Address,
			Status:         u.Status,
			Description:    u.Description,
			CreatedAt:      u.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
			UpdatedAt:      u.UpdatedAt.In(loc).Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}

func (s *universityService) GetUniversitiesByStatus(ctx context.Context, status string) ([]*models.University, error) {
	return s.universityRepo.GetUniversitiesByStatus(ctx, status)
}
func (s *universityService) GetUniversityByCode(ctx context.Context, code string) (*models.University, error) {
	return s.universityRepo.GetUniversityByCode(ctx, code)
}
