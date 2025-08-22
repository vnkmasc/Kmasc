package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/blockchain"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainService interface {
	VerifyBatch(ctx context.Context, universityID, facultyIDStr, certificateType, course, ediplomaID string) (bool, string, *models.EDiplomaResponse, error)
	VerifyEDiploma(ctx context.Context, universityID, facultyIDStr, course, studentCode string) (*VerifyResult, error)
	PushCertificateToChain(ctx context.Context, certificateID primitive.ObjectID) (string, error)
	GetCertificateFromChain(ctx context.Context, certificateID string) (*models.CertificateOnChain, error)
	VerifyFileByID(ctx context.Context, certID primitive.ObjectID) (io.ReadCloser, string, error)
	VerifyCertificateIntegrity(ctx context.Context, certID string) (bool, string, *models.CertificateOnChain, *models.Certificate, *models.User, *models.Faculty, *models.University, error)
	PushToBlockchain(ctx context.Context, universityID, facultyIDStr, certificateType, course string, issued *bool) (int, error)
	PushToBlockchain1(ctx context.Context, universityID, facultyIDStr, certificateType, course string) (int, error)
}

type blockchainService struct {
	templateRepo   repository.TemplateRepository
	ediplomaRepo   repository.EDiplomaRepository
	certRepo       repository.CertificateRepository
	userRepo       repository.UserRepository
	facultyRepo    repository.FacultyRepository
	universityRepo repository.UniversityRepository
	fabricClient   *blockchain.FabricClient
	minioClient    *database.MinioClient
}

func NewBlockchainService(
	templateRepo repository.TemplateRepository,
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
		templateRepo:   templateRepo,
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
func (s *blockchainService) PushToBlockchain1(
	ctx context.Context,
	universityID, facultyIDStr, certificateType, course string,
) (int, error) {

	// Build dynamic filter
	filter := bson.M{}
	parts := []string{universityID} // dùng để tạo BatchID linh hoạt

	if facultyIDStr != "" {
		facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
		if err != nil {
			return 0, common.ErrInvalidFaculty
		}
		filter["faculty_id"] = facultyID
		parts = append(parts, facultyIDStr) // thêm vào BatchID
	}

	if certificateType != "" {
		filter["certificate_type"] = bson.M{"$regex": certificateType, "$options": "i"}
		parts = append(parts, certificateType) // thêm vào BatchID
	}

	if course != "" {
		filter["course"] = bson.M{"$regex": course, "$options": "i"}
		parts = append(parts, course) // thêm vào BatchID
	}

	// Lấy danh sách EDiploma
	ediplomas, err := s.ediplomaRepo.FindByDynamicFilter(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to load eDiplomas: %w", err)
	}
	if len(ediplomas) == 0 {
		return 0, common.ErrNoDiplomas
	}

	// Gom hash
	infoHashes := []string{}
	fileHashes := []string{}
	for _, ed := range ediplomas {
		if ed.DataEncrypted && ed.Issued && !ed.OnBlockchain {
			hInfo, err := hashEDiplomaInfo(ed)
			if err != nil {
				log.Printf("hash failed for %s: %v", ed.StudentCode, err)
				continue
			}
			infoHashes = append(infoHashes, hInfo)

			if ed.EDiplomaFileHash != "" {
				fileHashes = append(fileHashes, ed.EDiplomaFileHash)
			}
		}
	}

	if len(infoHashes) == 0 {
		return 0, common.ErrNoValidDiplomas
	}

	aggregateInfoHash, _ := aggregateHashes(infoHashes)
	aggregateFileHash, _ := aggregateHashes(fileHashes)
	if aggregateFileHash == "" {
		aggregateFileHash = "NO_FILE_HASH"
	}

	// Tạo BatchID linh động
	batchID := strings.Join(parts, "-")

	batchOnChain := models.EDiplomaBatchOnChain{
		BatchID:           batchID,
		UniversityID:      universityID,
		FacultyID:         facultyIDStr,
		CertificateType:   certificateType,
		Course:            course,
		AggregateInfoHash: aggregateInfoHash,
		AggregateFileHash: aggregateFileHash,
		Count:             len(infoHashes),
	}

	// Đẩy lên blockchain
	txID, err := s.fabricClient.IssueEDiplomaBatch(batchOnChain)
	if err != nil {
		if errors.Is(err, common.ErrAlreadyOnChain) {
			return 0, common.ErrAlreadyOnChain
		}
		return 0, fmt.Errorf("push blockchain failed: %w", err)
	}

	// Update DB
	updatedCount := 0
	for _, ed := range ediplomas {
		if ed.DataEncrypted && ed.Issued && !ed.OnBlockchain {
			update := bson.M{
				"$set": bson.M{
					"on_blockchain":  true,
					"transaction_id": txID,
					"batch_id":       batchID,
					"university_id":  universityID,
					"updated_at":     time.Now(),
				},
			}
			if err := s.ediplomaRepo.UpdateByID(ctx, ed.ID, update); err != nil {
				log.Printf("Failed to update %s: %v", ed.StudentCode, err)
				continue
			}
			updatedCount++
		}
	}

	return updatedCount, nil
}

func (s *blockchainService) PushToBlockchain(
	ctx context.Context,
	universityID, facultyIDStr, certificateType, course string,
	issued *bool,
) (int, error) {

	// Build dynamic filter
	filter := bson.M{}
	if facultyIDStr != "" {
		facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
		if err != nil {
			return 0, common.ErrInvalidFaculty
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

	// Lấy danh sách EDiploma
	ediplomas, err := s.ediplomaRepo.FindByDynamicFilter(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to load eDiplomas: %w", err)
	}
	if len(ediplomas) == 0 {
		return 0, common.ErrNoDiplomas
	}
	log.Printf("[BlockchainService] Found %d eDiplomas", len(ediplomas))

	// Gom hash
	infoHashes := []string{}
	fileHashes := []string{}
	hashMap := map[string]*models.EDiploma{} // map hash -> EDiploma để sinh proof sau

	for _, ed := range ediplomas {
		if ed.DataEncrypted && ed.Issued && !ed.OnBlockchain {
			hInfo, err := hashEDiplomaInfo(ed)
			if err != nil {
				log.Printf("[BlockchainService] hash failed for %s: %v", ed.StudentCode, err)
				continue
			}
			infoHashes = append(infoHashes, hInfo)
			hashMap[hInfo] = ed

			if ed.EDiplomaFileHash != "" {
				fileHashes = append(fileHashes, ed.EDiplomaFileHash)
			}
		}
	}

	if len(infoHashes) == 0 {
		return 0, common.ErrNoValidDiplomas
	}
	log.Printf("[BlockchainService] Info hashes: %+v", infoHashes)
	log.Printf("[BlockchainService] File hashes: %+v", fileHashes)

	// Tạo Merkle root
	merkleInfoTree := models.NewMerkleTreeFromStrings(infoHashes)
	merkleInfoRoot := merkleInfoTree.RootHash()
	log.Printf("[BlockchainService] Merkle Info Root: %s", merkleInfoRoot)

	merkleFileRoot := "NO_FILE_HASH"
	if len(fileHashes) > 0 {
		merkleFileTree := models.NewMerkleTreeFromStrings(fileHashes)
		merkleFileRoot = merkleFileTree.RootHash()
		log.Printf("[BlockchainService] Merkle File Root: %s", merkleFileRoot)
	}

	// Sinh BatchID
	batchID := fmt.Sprintf("%s-%s-%s", universityID, facultyIDStr, course)
	log.Printf("[BlockchainService] BatchID = %s", batchID)

	batchOnChain := models.EDiplomaBatchOnChain{
		BatchID:           batchID,
		UniversityID:      universityID,
		FacultyID:         facultyIDStr,
		CertificateType:   certificateType,
		Course:            course,
		AggregateInfoHash: merkleInfoRoot,
		AggregateFileHash: merkleFileRoot,
		Count:             len(infoHashes),
	}

	// Push lên blockchain
	txID, err := s.fabricClient.IssueEDiplomaBatch(batchOnChain)
	if err != nil {
		return 0, fmt.Errorf("push blockchain failed: %w", err)
	}
	log.Printf("[BlockchainService] TxID = %s", txID)

	// Update DB với log chi tiết
	updatedCount := 0
	for hInfo, ed := range hashMap {
		proof := merkleInfoTree.GetProof(hInfo)
		if len(proof) == 0 {
			log.Printf("[⚠️ BlockchainService] Proof is EMPTY for %s (hash=%s)", ed.StudentCode, hInfo)
		} else {
			log.Printf("[✅ BlockchainService] Proof for %s (%s): %+v", ed.StudentCode, hInfo, proof)
		}

		update := bson.M{
			"$set": bson.M{
				"on_blockchain":  true,
				"transaction_id": txID,
				"batch_id":       batchID,
				"university_id":  universityID,
				"merkle_proof":   proof,
				"updated_at":     time.Now(),
			},
		}

		err := s.ediplomaRepo.UpdateByID(ctx, ed.ID, update)
		if err != nil {
			log.Printf("[❌ BlockchainService] Failed to update %s: %v", ed.StudentCode, err)
			continue
		}
		log.Printf("[🟢 BlockchainService] Updated %s successfully", ed.StudentCode)
		updatedCount++

	}

	log.Printf("[BlockchainService] Total updated records: %d", updatedCount)
	return updatedCount, nil
}

func hashEDiplomaInfo(ed *models.EDiploma) (string, error) {
	data := ed.StudentCode +
		ed.FullName +
		ed.CertificateType +
		ed.Course +
		ed.EducationType +
		fmt.Sprintf("%.2f", ed.GPA) +
		ed.GraduationRank +
		ed.IssueDate.Format("2006-01-02") +
		ed.SerialNumber +
		ed.RegistrationNumber

	h := sha256.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func aggregateHashes(hashes []string) (string, error) {
	if len(hashes) == 0 {
		return "", nil
	}
	combined := strings.Join(hashes, "")
	h := sha256.New()
	_, err := h.Write([]byte(combined))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *blockchainService) VerifyBatch(
	ctx context.Context,
	universityID, facultyIDStr, certificateType, course, ediplomaID string,
) (bool, string, *models.EDiplomaResponse, error) {

	filter := bson.M{}

	if ediplomaID != "" {
		edID, err := primitive.ObjectIDFromHex(ediplomaID)
		if err != nil {
			return false, "", nil, fmt.Errorf("invalid ediploma_id")
		}
		filter["_id"] = edID
	} else {
		if facultyIDStr != "" {
			facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
			if err != nil {
				return false, "", nil, fmt.Errorf("invalid faculty_id")
			}
			filter["faculty_id"] = facultyID
		}
		if certificateType != "" {
			filter["certificate_type"] = bson.M{"$regex": certificateType, "$options": "i"}
		}
		if course != "" {
			filter["course"] = bson.M{"$regex": course, "$options": "i"}
		}
	}

	ediplomas, err := s.ediplomaRepo.FindByDynamicFilter(ctx, filter)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to load eDiplomas: %w", err)
	}
	if len(ediplomas) == 0 {
		return false, "", nil, fmt.Errorf("no eDiplomas found")
	}

	// Chỉ lấy DTO nếu ediplomaID được truyền
	var targetDTO *models.EDiplomaResponse
	if ediplomaID != "" {
		id, _ := primitive.ObjectIDFromHex(ediplomaID)

		// Lấy trực tiếp eDiploma
		ediploma, err := s.ediplomaRepo.FindByID(ctx, id)
		if err != nil || ediploma == nil {
			return false, "", nil, fmt.Errorf("ediploma not found")
		}

		// Lấy các đối tượng liên quan từ repository
		university, _ := s.universityRepo.FindByID(ctx, ediploma.UniversityID)
		faculty, _ := s.facultyRepo.FindByID(ctx, ediploma.FacultyID)
		template, _ := s.templateRepo.GetByID(ctx, ediploma.TemplateID)
		user, _ := s.userRepo.GetUserByID(ctx, ediploma.UserID)

		// Map sang DTO
		targetDTO = mapper.MapEDiplomaToDTO(ediploma, university, faculty, template, user)
	}

	// Kiểm tra OnBlockchain
	var notOnChain []string
	for _, ed := range ediplomas {
		if !ed.OnBlockchain {
			notOnChain = append(notOnChain, ed.StudentCode)
		}
	}
	if len(notOnChain) > 0 {
		msg := fmt.Sprintf("Có %d eDiplomas chưa được đẩy lên blockchain: %v", len(notOnChain), notOnChain)
		return false, msg, targetDTO, nil
	}

	// Hash từ DB
	infoHashes := []string{}
	fileHashes := []string{}
	for _, ed := range ediplomas {
		hInfo, err := hashEDiplomaInfo(ed)
		if err != nil {
			log.Printf("[VerifyBatch] hash failed for %s: %v", ed.StudentCode, err)
			continue
		}
		infoHashes = append(infoHashes, hInfo)
		if ed.EDiplomaFileHash != "" {
			fileHashes = append(fileHashes, ed.EDiplomaFileHash)
		}
	}

	aggregateInfoHash, _ := aggregateHashes(infoHashes)
	aggregateFileHash, _ := aggregateHashes(fileHashes)
	if aggregateFileHash == "" {
		aggregateFileHash = "NO_FILE_HASH"
	}

	parts := []string{universityID}
	if facultyIDStr != "" {
		parts = append(parts, facultyIDStr)
	}
	if certificateType != "" {
		parts = append(parts, certificateType)
	}
	if course != "" {
		parts = append(parts, course)
	}
	batchID := strings.Join(parts, "-")
	log.Printf("[VerifyBatch] batchID=%s", batchID)

	batchOnChain, err := s.fabricClient.GetEDiplomaBatch(batchID)
	if err != nil {
		return false, "", targetDTO, fmt.Errorf("failed to get batch from blockchain: %w", err)
	}

	// So sánh
	if batchOnChain.AggregateInfoHash != aggregateInfoHash {
		return false, fmt.Sprintf("Mismatch AggregateInfoHash: onChain=%s, local=%s", batchOnChain.AggregateInfoHash, aggregateInfoHash), targetDTO, nil
	}
	if batchOnChain.AggregateFileHash != aggregateFileHash {
		return false, fmt.Sprintf("Mismatch AggregateFileHash: onChain=%s, local=%s", batchOnChain.AggregateFileHash, aggregateFileHash), targetDTO, nil
	}
	if batchOnChain.Count != len(infoHashes) {
		return false, fmt.Sprintf("Mismatch Count: onChain=%d, local=%d", batchOnChain.Count, len(infoHashes)), targetDTO, nil
	}

	return true, "Dữ liệu khớp hoàn toàn trên chuỗi khối", targetDTO, nil
}

type VerifyResult struct {
	StudentCode    string             `json:"student_code"`
	Valid          bool               `json:"valid"`
	BlockchainRoot string             `json:"blockchain_root"`
	ComputedHash   string             `json:"computed_hash"`
	Proof          []models.ProofNode `json:"proof"`
}

func (s *blockchainService) VerifyEDiploma(
	ctx context.Context,
	universityID, facultyIDStr, course, studentCode string,
) (*VerifyResult, error) {

	// 1. Lấy eDiploma từ DB
	ed, err := s.ediplomaRepo.FindByStudentCode(ctx, studentCode)
	if err != nil {
		return nil, err
	}
	if ed == nil {
		return nil, common.ErrNoDiplomas
	}

	// 2. Hash lại thông tin eDiploma
	hInfo, err := hashEDiplomaInfo(ed)
	if err != nil {
		return nil, fmt.Errorf("failed to hash diploma: %w", err)
	}

	// 3. Lấy proof đã lưu từ DB
	proof := ed.MerkleProof
	if len(proof) == 0 {
		return nil, fmt.Errorf("no merkle proof stored for this diploma")
	}

	// 4. Tạo batchID dựa trên cùng công thức push
	batchID := fmt.Sprintf("%s-%s-%s", universityID, facultyIDStr, course)

	// 5. Lấy batch từ blockchain
	batchOnChain, err := s.fabricClient.GetEDiplomaBatch(batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch from blockchain: %w", err)
	}

	// 6. Lấy Merkle root từ batch
	root := batchOnChain.AggregateInfoHash
	if root == "" {
		return nil, fmt.Errorf("no merkle root stored on blockchain for batch %s", batchID)
	}

	// 7. Verify proof
	isValid := models.VerifyProof(hInfo, proof, root)

	// 8. Trả kết quả
	return &VerifyResult{
		StudentCode:    studentCode,
		Valid:          isValid,
		BlockchainRoot: root,
		ComputedHash:   hInfo,
		Proof:          proof,
	}, nil
}
