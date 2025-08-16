package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/blockchain"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainService interface {
	PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error)
	GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error)
	VerifyFileByID(ctx context.Context, certID primitive.ObjectID) (io.ReadCloser, string, error)
	VerifyCertificateIntegrity(ctx context.Context, certID string) (bool, string, *models.CertificateOnChain, *models.Certificate, *models.User, *models.Faculty, *models.University, error)
	PushToBlockchain(ctx context.Context, facultyIDStr, certificateType, course string, issued *bool) (int, error)
}

type blockchainService struct {
	ediplomaRepo   repository.EDiplomaRepository
	certRepo       repository.CertificateRepository
	userRepo       repository.UserRepository
	facultyRepo    repository.FacultyRepository
	universityRepo repository.UniversityRepository
	fabricClient   *blockchain.FabricClient
	minioClient    *database.MinioClient
}

func NewBlockchainService(
	ediplomaRepo repository.EDiplomaRepository,
	certRepo repository.CertificateRepository,
	userRepo repository.UserRepository,
	facultyRepo repository.FacultyRepository,
	universityRepo repository.UniversityRepository,
	fabricClient *blockchain.FabricClient,
	minioClient *database.MinioClient,
) BlockchainService {
	return &blockchainService{
		ediplomaRepo:   ediplomaRepo,
		certRepo:       certRepo,
		userRepo:       userRepo,
		facultyRepo:    facultyRepo,
		universityRepo: universityRepo,
		fabricClient:   fabricClient,
		minioClient:    minioClient,
	}
}

func (s *blockchainService) PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error) {
	if s.fabricClient == nil {
		return "", fmt.Errorf("❌ fabricClient is nil trong blockchainService")
	}
	cert, err := s.certRepo.GetCertificateByID(ctx, certificateID)
	if err != nil || cert == nil {
		return "", common.ErrCertificateNotFound
	}

	if cert.CertHash == "" {
		return "", fmt.Errorf("%w", common.ErrCertificateMissingHash)
	}
	if !cert.PhysicalCopyIssued {
		return "", fmt.Errorf("%w", common.ErrCertificateNoFile)
	}
	// if !cert.Signed {
	// 	return "", fmt.Errorf("%w", common.ErrCertificateNotSigned)
	// }
	if cert.OnBlockchain {
		return "", fmt.Errorf("%w", common.ErrCertificateAlreadyOnChain)
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
			"transaction_id": txID,
			"on_blockchain":  true,
			"updated_at":     time.Now(),
		},
	}
	if err := s.certRepo.UpdateCertificateByID(ctx, certificateID, update); err != nil {
		return "", fmt.Errorf("không thể cập nhật blockchain_tx_id: %w", err)
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

func (s *blockchainService) VerifyCertificateIntegrity(ctx context.Context, certID string) (
	bool, string,
	*models.CertificateOnChain,
	*models.Certificate,
	*models.User,
	*models.Faculty,
	*models.University,
	error,
) {
	onChainCert, err := s.fabricClient.GetCertificateByID(certID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("lỗi lấy từ blockchain: %w", err)
	}

	certificateObjID, err := primitive.ObjectIDFromHex(certID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("certID không hợp lệ: %w", err)
	}

	cert, err := s.certRepo.GetCertificateByID(ctx, certificateObjID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("không tìm thấy văn bằng trong MongoDB: %w", err)
	}

	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("không tìm thấy sinh viên: %w", err)
	}

	faculty, err := s.facultyRepo.FindByID(ctx, cert.FacultyID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("không tìm thấy khoa: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil {
		return false, "", nil, nil, nil, nil, nil, fmt.Errorf("không tìm thấy trường đại học: %w", err)
	}

	localHash := generateCertificateHash(cert, user, faculty, university)

	if localHash != onChainCert.CertHash {
		return false, "Dữ liệu đã bị thay đổi!", onChainCert, cert, user, faculty, university, nil
	}

	return true, "Dữ liệu khớp hoàn toàn với blockchain", onChainCert, cert, user, faculty, university, nil
}

func (s *blockchainService) VerifyFileByID(ctx context.Context, certID primitive.ObjectID) (io.ReadCloser, string, error) {
	// Lấy certificate từ MongoDB
	certificate, err := s.certRepo.GetCertificateByID(ctx, certID)
	if err != nil || certificate == nil {
		return nil, "", fmt.Errorf("không tìm thấy certificate")
	}

	// Lấy hash từ blockchain
	onChainCert, err := s.fabricClient.GetCertificateByID(certID.Hex())
	if err != nil {
		return nil, "", fmt.Errorf("lỗi lấy dữ liệu từ blockchain: %w", err)
	}
	hashOnChain := onChainCert.HashFile

	// Lấy stream file từ MinIO
	stream, contentType, err := s.minioClient.DownloadFileStream(ctx, certificate.Path)
	if err != nil {
		return nil, "", fmt.Errorf("không thể tải file từ MinIO: %w", err)
	}

	// Hash lại từ stream (phải đọc stream rồi tạo lại reader)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	stream.Close() // Đóng stream cũ
	if err != nil {
		return nil, "", fmt.Errorf("lỗi đọc file từ stream: %w", err)
	}

	hash := sha256.Sum256(buf.Bytes())
	currentHash := hex.EncodeToString(hash[:])

	if !strings.EqualFold(currentHash, hashOnChain) {
		return nil, "", fmt.Errorf("file đã bị sửa đổi hoặc không khớp với blockchain")
	}

	// Tạo lại reader mới để trả cho handler
	newReader := io.NopCloser(bytes.NewReader(buf.Bytes()))
	return newReader, contentType, nil
}

func (s *blockchainService) PushToBlockchain(
	ctx context.Context,
	facultyIDStr, certificateType, course string,
	issued *bool,
) (int, error) {

	// Build dynamic filter
	filter := bson.M{}
	if facultyIDStr != "" {
		facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
		if err != nil {
			return 0, fmt.Errorf("invalid faculty_id")
		}
		filter["faculty_id"] = facultyID
	}
	if certificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": certificateType, "$options": "i"}
	}
	if course != "" {
		filter["course"] = bson.M{"$regex": course, "$options": "i"}
	}
	if issued != nil {
		filter["issued"] = *issued
	}

	// Lấy danh sách EDiploma thỏa filter
	ediplomas, err := s.ediplomaRepo.FindByDynamicFilter(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to load eDiplomas: %w", err)
	}

	updatedCount := 0
	for _, ed := range ediplomas {
		if ed.DataEncrypted && ed.Issued && !ed.OnBlockchain {
			ed.OnBlockchain = true
			ed.UpdatedAt = time.Now()
			if err := s.ediplomaRepo.Update(ctx, ed.ID, ed); err != nil {
				log.Printf("Failed to update OnBlockchain for %s: %v", ed.StudentCode, err)
				continue
			}
			updatedCount++
		}
	}

	return updatedCount, nil
}
