package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	RequestOTP(ctx context.Context, input models.RequestOTPInput) error
	VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) (string, error)
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, email, password string) (*models.Account, error)
	DeleteAccountByEmail(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, accountID primitive.ObjectID, oldPass, newPass string) error
	GetAllAccounts(ctx context.Context, page, pageSize int) ([]*models.Account, int64, error)
	GetAccountByID(ctx context.Context, id primitive.ObjectID) (*models.Account, error)
	GetAccountsByRole(ctx context.Context, role string) ([]models.Account, error)
}

type authService struct {
	authRepo    repository.AuthRepository
	userRepo    repository.UserRepository
	emailSender utils.EmailSender
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, emailSender utils.EmailSender) AuthService {
	return &authService{
		authRepo:    authRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}
func (s *authService) GetAccountByID(ctx context.Context, id primitive.ObjectID) (*models.Account, error) {
	return s.authRepo.FindByID(ctx, id)
}

func (s *authService) RequestOTP(ctx context.Context, input models.RequestOTPInput) error {
	user, err := s.userRepo.FindByEmail(ctx, input.StudentEmail)
	if err != nil {
		return common.ErrUserNotExisted
	}
	existingAccount, err := s.authRepo.FindPersonalAccountByUserID(ctx, user.ID)
	if err != nil {
		return common.ErrCheckingPersonalAccount
	}
	if existingAccount != nil {
		return common.ErrPersonalAccountAlreadyExist
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpData := models.OTP{
		Email:     input.StudentEmail,
		Code:      otp,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	if err := s.authRepo.SaveOTP(ctx, otpData); err != nil {
		return err
	}

	body := fmt.Sprintf("Mã OTP của bạn là: %s. Có hiệu lực trong 3 phút.", otp)
	return s.emailSender.SendEmail(input.StudentEmail, "Mã xác thực OTP", body)
}

func (s *authService) VerifyOTP(ctx context.Context, input *models.VerifyOTPRequest) (string, error) {
	otpRecord, err := s.authRepo.FindLatestOTPByEmail(ctx, input.StudentEmail)
	if err != nil {
		return "", fmt.Errorf("không tìm thấy mã OTP")
	}

	if otpRecord.Code != input.OTP {
		return "", fmt.Errorf("mã OTP không đúng")
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		return "", fmt.Errorf("mã OTP đã hết hạn")
	}

	user, err := s.userRepo.FindByEmail(ctx, input.StudentEmail)
	if err != nil {
		return "", fmt.Errorf("lỗi khi tìm người dùng: %v", err)
	}
	if user == nil {
		return "", fmt.Errorf("người dùng không tồn tại")
	}

	return user.ID.Hex(), nil
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) error {
	exists, err := s.authRepo.IsPersonalEmailExist(ctx, req.PersonalEmail)
	if err != nil {
		return fmt.Errorf("lỗi kiểm tra email: %w", err)
	}
	if exists {
		return fmt.Errorf("email cá nhân đã được sử dụng")
	}

	userObjID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return fmt.Errorf("user_id không hợp lệ")
	}

	user, err := s.userRepo.GetUserByID(ctx, userObjID)
	if err != nil {
		return fmt.Errorf("không tìm thấy user: %v", err)
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("lỗi hash mật khẩu: %w", err)
	}

	account := &models.Account{
		StudentID:     user.ID,
		StudentEmail:  user.Email,
		PersonalEmail: req.PersonalEmail,
		PasswordHash:  hash,
		CreatedAt:     time.Now(),
		Role:          "student",
	}

	if err := s.authRepo.CreateAccount(ctx, account); err != nil {
		return fmt.Errorf("không tạo được tài khoản: %w", err)
	}

	return nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.Account, error) {
	account, err := s.authRepo.FindByPersonalEmail(ctx, email)
	if err != nil {
		return nil, errors.New("tài khoản không tồn tại")
	}

	if !utils.ComparePassword(account.PasswordHash, password) {
		return nil, errors.New("mật khẩu không đúng")
	}

	return account, nil
}

func (s *authService) GetAllAccounts(ctx context.Context, page, pageSize int) ([]*models.Account, int64, error) {
	return s.authRepo.GetAllAccounts(ctx, page, pageSize)
}

func (s *authService) DeleteAccountByEmail(ctx context.Context, email string) error {
	err := s.authRepo.DeleteAccountByEmail(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, accountID primitive.ObjectID, oldPass, newPass string) error {
	account, err := s.authRepo.FindByID(ctx, accountID)
	if err != nil || account == nil {
		return common.ErrAccountNotFound
	}

	if !utils.CheckPasswordHash(oldPass, account.PasswordHash) {
		return common.ErrInvalidOldPassword
	}

	newHash, err := utils.HashPassword(newPass)
	if err != nil {
		return err
	}
	return s.authRepo.UpdatePassword(ctx, accountID, newHash)
}
func (s *authService) GetAccountsByRole(ctx context.Context, role string) ([]models.Account, error) {
	return s.authRepo.FindByRole(ctx, role)
}
