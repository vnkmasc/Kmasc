package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/blockchain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainService interface {
	PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error)
	GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error)
	VerifyCertificateIntegrity(ctx context.Context, certID string) (bool, string, *models.CertificateOnChain, error)
}

type blockchainService struct {
	certRepo       repository.CertificateRepository
	userRepo       repository.UserRepository
	facultyRepo    repository.FacultyRepository
	universityRepo repository.UniversityRepository
	fabricClient   *blockchain.FabricClient
}

func NewBlockchainService(
	certRepo repository.CertificateRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	fabricClient *blockchain.FabricClient,
) BlockchainService {
	return &blockchainService{
		certRepo:       certRepo,
		userRepo:       userRepo,
		facultyRepo:    facultyRepo,
		universityRepo: universityRepo,
		fabricClient:   fabricClient,
	}
}

func (s *blockchainService) PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error) {
	cert, err := s.certRepo.GetCertificateByID(ctx, certificateID)
	if err != nil || cert == nil {
		return "", common.ErrCertificateNotFound
	}
	if cert.CertHash == "" {
		return "", fmt.Errorf("certificate chưa có cert_hash")
	}

	chainData := models.CertificateOnChain{
		CertID:              cert.ID.Hex(),
		CertHash:            cert.CertHash,
		HashFile:            cert.HashFile,
		UniversitySignature: "",
		DateOfIssuing:       cert.IssueDate.Format("2006-01-02"),
		SerialNumber:        cert.SerialNumber,
		RegNo:               cert.RegNo,
		Version:             1,
		UpdatedDate:         cert.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	txID, err := s.fabricClient.IssueCertificate(chainData)
	if err != nil {
		return "", err
	}

	update := bson.M{
		"$set": bson.M{
			"blockchain_tx_id": txID,
			"updated_at":       time.Now(),
		},
	}
	if err := s.certRepo.UpdateCertificateByID(ctx, certificateID, update); err != nil {
		return "", fmt.Errorf("không thể cập nhật blockchain_tx_id: %v", err)
	}

	return txID, nil
}

func (s *blockchainService) GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error) {
	cert, err := s.fabricClient.GetCertificateByID(certificateID)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (s *blockchainService) VerifyCertificateIntegrity(ctx context.Context, certID string) (bool, string, *models.CertificateOnChain, error) {
	onChainCert, err := s.fabricClient.GetCertificateByID(certID)
	if err != nil {
		return false, "", nil, fmt.Errorf("lỗi lấy từ blockchain: %w", err)
	}

	certificateObjID, err := primitive.ObjectIDFromHex(certID)
	if err != nil {
		return false, "", nil, fmt.Errorf("certID không hợp lệ: %w", err)
	}

	cert, err := s.certRepo.GetCertificateByID(ctx, certificateObjID)
	if err != nil {
		return false, "", nil, fmt.Errorf("không tìm thấy văn bằng trong MongoDB: %w", err)
	}

	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil {
		return false, "", nil, fmt.Errorf("không tìm thấy sinh viên: %w", err)
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil {
		return false, "", nil, fmt.Errorf("không tìm thấy khoa: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil {
		return false, "", nil, fmt.Errorf("không tìm thấy trường đại học: %w", err)
	}

	localHash := generateCertificateHash(cert, user, faculty, university)

	if localHash != onChainCert.CertHash {
		return false, "Dữ liệu đã bị thay đổi!", onChainCert, nil
	}

	return true, "Dữ liệu khớp hoàn toàn với blockchain", onChainCert, nil
}
