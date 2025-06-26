package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationService interface {
	CreateVerificationCode(ctx context.Context, code *models.VerificationCode) error
	GetCodesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int64) ([]models.VerificationCodeResponse, int64, error)
	VerifyCode(ctx context.Context, code, viewType string) (*models.VerificationCode, *models.CertificateResponse, error)
}

type verificationService struct {
	repo               repository.VerificationRepository
	certificateService CertificateService
}

func NewVerificationService(repo repository.VerificationRepository, certSvc CertificateService) VerificationService {
	return &verificationService{
		repo:               repo,
		certificateService: certSvc,
	}
}

func (s *verificationService) CreateVerificationCode(ctx context.Context, code *models.VerificationCode) error {
	code.ID = primitive.NewObjectID()
	code.Code = generateRandomCode(8)
	code.CreatedAt = time.Now()

	return s.repo.Save(ctx, code)
}

func generateRandomCode(length int) string {
	return uuid.New().String()[:length]
}

func (s *verificationService) GetCodesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int64) ([]models.VerificationCodeResponse, int64, error) {
	codes, total, err := s.repo.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	var res []models.VerificationCodeResponse
	for _, code := range codes {
		minutesRemaining := int64(code.ExpiredAt.Sub(now).Minutes())
		if minutesRemaining < 0 {
			minutesRemaining = 0
		}

		res = append(res, models.VerificationCodeResponse{
			ID:               code.ID,
			Code:             code.Code,
			CanViewScore:     code.CanViewScore,
			CanViewData:      code.CanViewData,
			CanViewFile:      code.CanViewFile,
			ExpiredInMinutes: minutesRemaining,
			CreatedAt:        code.CreatedAt,
		})
	}

	return res, total, nil
}

func (s *verificationService) VerifyCode(ctx context.Context, code, viewType string) (*models.VerificationCode, *models.CertificateResponse, error) {
	vc, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, nil, errors.New("mã không tồn tại")
	}
	if time.Now().After(vc.ExpiredAt) {
		return nil, nil, errors.New("mã đã hết hạn")
	}

	switch viewType {
	case "score":
		if !vc.CanViewScore {
			return nil, nil, errors.New("không có quyền xem điểm")
		}
	case "data":
		if !vc.CanViewData {
			return nil, nil, errors.New("không có quyền xem thông tin")
		}
	case "file":
		if !vc.CanViewFile {
			return nil, nil, errors.New("không có quyền xem file")
		}
	default:
		return nil, nil, errors.New("loại dữ liệu không hợp lệ")
	}

	if viewType == "data" || viewType == "file" {
		certResp, err := s.certificateService.GetCertificateByUserID(ctx, vc.UserID)
		if err != nil {
			return nil, nil, err
		}
		return vc, certResp, nil
	}

	return vc, nil, nil
}
