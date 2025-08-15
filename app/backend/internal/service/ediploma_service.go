package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiplomaService interface {
	GetByID(ctx context.Context, id string) (*models.EDiploma, error)
	GenerateEDiploma(ctx context.Context, certificateIDStr, templateIDStr string) (*models.EDiploma, error)
	GenerateBulkEDiplomas(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error)
	UploadLocalEDiplomas(ctx context.Context) []map[string]interface{}
	ProcessZip(ctx context.Context, zipPath string) ([]*models.EDiploma, error)
	GetDiplomaPDF(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, int64, string, error)
	SearchEDiplomaDTOs(ctx context.Context, filter models.EDiplomaSearchFilter) ([]*models.EDiplomaDTO, int64, error)
	GenerateBulkEDiplomasZip(ctx context.Context, facultyIDStr, certificateType, course string, issued *bool, templateIDStr string) (string, error)
}

type eDiplomaService struct {
	universityRepo  repository.UniversityRepository
	majorRepo       repository.MajorRepository
	facultyRepo     repository.FacultyRepository
	repo            repository.EDiplomaRepository
	certificateRepo repository.CertificateRepository
	templateRepo    repository.TemplateRepository
	userRepo        repository.UserRepository
	minioClient     *database.MinioClient
	templateEngine  *models.TemplateEngine
	pdfGenerator    *utils.PDFGenerator
}

func NewEDiplomaService(
	universityRepo repository.UniversityRepository,
	majorRepo repository.MajorRepository,
	facultyRepo repository.FacultyRepository,
	repo repository.EDiplomaRepository,
	certificateRepo repository.CertificateRepository,
	templateRepo repository.TemplateRepository,
	userRepo repository.UserRepository,
	minioClient *database.MinioClient,
	templateEngine *models.TemplateEngine,
	pdfGenerator *utils.PDFGenerator,
) *eDiplomaService {
	return &eDiplomaService{
		universityRepo:  universityRepo,
		majorRepo:       majorRepo,
		facultyRepo:     facultyRepo,
		repo:            repo,
		certificateRepo: certificateRepo,
		templateRepo:    templateRepo,
		userRepo:        userRepo,
		minioClient:     minioClient,
		templateEngine:  templateEngine,
		pdfGenerator:    pdfGenerator,
	}
}

func (s *eDiplomaService) GetByID(ctx context.Context, id string) (*models.EDiploma, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid diploma id")
	}
	return s.repo.FindByID(ctx, objID)
}
func parseMinioURL(url string) (bucket, object string, err error) {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return "", "", err
	}

	// Ví dụ: /certificates/diploma_template/KMA/CNTT/template.html
	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid MinIO URL: %s", url)
	}

	return parts[0], parts[1], nil
}

func (s *eDiplomaService) GenerateEDiploma(ctx context.Context, certificateIDStr, templateIDStr string) (*models.EDiploma, error) {
	certificateID, err := primitive.ObjectIDFromHex(certificateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate ID")
	}
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID")
	}

	cert, err := s.certificateRepo.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, fmt.Errorf("certificate not found")
	}

	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}

	calculatedHash := utils.ComputeSHA256([]byte(template.HTMLContent))
	if calculatedHash != template.HashTemplate {
		return nil, fmt.Errorf("template content hash mismatch - data may be corrupted")
	}
	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get university: %w", err)
	}

	dob, _ := time.Parse("2006-01-02", user.DateOfBirth)

	data := map[string]interface{}{
		"SoHieu":         cert.SerialNumber,
		"SoVaoSo":        cert.RegNo,
		"HoTen":          user.FullName,
		"NgaySinh":       dob.Format("02/01/2006"),
		"TenTruong":      university.UniversityName,
		"Nganh":          cert.Major,
		"XepLoai":        cert.GraduationRank,
		"HinhThucDaoTao": cert.EducationType,
		"Khoa":           cert.Course,
		"NgayCap":        cert.IssueDate.Format("02/01/2006"),
	}

	renderedHTML, err := s.templateEngine.Render(template.HTMLContent, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	pdfBytes, err := s.pdfGenerator.ConvertHTMLToPDF(renderedHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	hash := utils.ComputeSHA256(pdfBytes)

	pdfPath := fmt.Sprintf("ediplomas/%s.pdf", primitive.NewObjectID().Hex())
	if err := s.minioClient.UploadFile(ctx, pdfPath, pdfBytes, "application/pdf"); err != nil {
		return nil, fmt.Errorf("failed to upload PDF: %w", err)
	}

	_ = s.templateRepo.UpdateIsLocked(ctx, templateID, true)

	// Lưu EDiploma
	now := time.Now()
	ediploma := &models.EDiploma{
		ID:                 primitive.NewObjectID(),
		TemplateID:         templateID,
		UniversityID:       cert.UniversityID,
		FacultyID:          cert.FacultyID,
		UserID:             cert.UserID,
		MajorID:            primitive.NilObjectID,
		StudentCode:        cert.StudentCode,
		FullName:           cert.Name,
		CertificateType:    cert.CertificateType,
		Course:             cert.Course,
		EducationType:      cert.EducationType,
		GPA:                cert.GPA,
		GraduationRank:     cert.GraduationRank,
		IssueDate:          cert.IssueDate,
		SerialNumber:       cert.SerialNumber,
		RegistrationNumber: cert.RegNo,
		EDiplomaFileLink:   pdfPath,
		EDiplomaFileHash:   hash,
		Signed:             false,
		Signature:          "",
		SignedAt:           time.Time{},
		OnBlockchain:       false,
		BlockchainTxID:     "",
		SignatureOfUni:     template.SignatureOfUni,
		SignatureOfMinEdu:  template.SignatureOfMinEdu,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.Save(ctx, ediploma); err != nil {
		return nil, fmt.Errorf("failed to save EDiploma: %w", err)
	}

	return ediploma, nil
}

func (s *eDiplomaService) SearchEDiplomaDTOs(ctx context.Context, filter models.EDiplomaSearchFilter) ([]*models.EDiplomaDTO, int64, error) {
	ediplomas, total, err := s.repo.SearchByFilters(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var dtoList []*models.EDiplomaDTO
	for _, ed := range ediplomas {
		university, _ := s.universityRepo.FindByID(ctx, ed.UniversityID)
		faculty, _ := s.facultyRepo.FindByID(ctx, ed.FacultyID)

		var major *models.Major
		if !ed.MajorID.IsZero() {
			major, _ = s.majorRepo.GetByID(ctx, ed.MajorID)
		}

		template, _ := s.templateRepo.GetByID(ctx, ed.TemplateID)
		user, _ := s.userRepo.GetUserByID(ctx, ed.UserID)

		dto := mapper.MapEDiplomaToDTO(ed, university, faculty, major, template, user)
		dtoList = append(dtoList, dto)
	}

	return dtoList, total, nil
}

func (s *eDiplomaService) GenerateBulkEDiplomas(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error) {
	var result []*models.EDiploma

	// Parse ID
	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid faculty ID")
	}
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID")
	}

	// 1. Load template từ MongoDB
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}
	if template.FacultyID != facultyID {
		return nil, errors.New("template does not belong to the given faculty")
	}

	// 2. Xác minh hash của HTMLContent
	calculatedHash := utils.ComputeSHA256([]byte(template.HTMLContent))
	if calculatedHash != template.HashTemplate {
		return nil, fmt.Errorf("template content hash mismatch - data may be corrupted")
	}

	// 3. Load tất cả certificates của faculty
	certificates, err := s.certificateRepo.FindByFacultyID(ctx, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates: %w", err)
	}

	for _, cert := range certificates {
		// Load university
		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil {
			log.Printf("Load university failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		// Load user
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
		if err != nil {
			log.Printf("Load user failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		// Parse ngày sinh
		dobTime, err := time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			log.Printf("Invalid date format for user %s: %v", user.ID.Hex(), err)
			continue
		}

		// Map dữ liệu render
		data := map[string]interface{}{
			"SoHieu":         cert.SerialNumber,
			"SoVaoSo":        cert.RegNo,
			"HoTen":          user.FullName,
			"NgaySinh":       dobTime.Format("02/01/2006"),
			"TenTruong":      university.UniversityName,
			"Nganh":          cert.Major,
			"XepLoai":        cert.GraduationRank,
			"HinhThucDaoTao": cert.EducationType,
			"Khoa":           cert.Course,
			"NgayCap":        cert.IssueDate.Format("02/01/2006"),
		}

		// Render HTML
		renderedHTML, err := s.templateEngine.Render(template.HTMLContent, data)
		if err != nil {
			log.Printf("Render failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		// Convert sang PDF
		pdfBytes, err := s.pdfGenerator.ConvertHTMLToPDF(renderedHTML)
		if err != nil {
			log.Printf("PDF generation failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		// Hash PDF
		hash := utils.ComputeSHA256(pdfBytes)

		// 📌 Nếu vẫn muốn lưu PDF ở MinIO
		pdfPath := fmt.Sprintf("ediplomas/%s/%s.pdf", university.UniversityCode, cert.StudentCode)
		if err := s.minioClient.UploadFile(ctx, pdfPath, pdfBytes, "application/pdf"); err != nil {
			log.Printf("Upload failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		now := time.Now()
		ediploma := &models.EDiploma{
			ID:                 primitive.NewObjectID(),
			TemplateID:         templateID,
			UniversityID:       cert.UniversityID,
			FacultyID:          cert.FacultyID,
			UserID:             cert.UserID,
			MajorID:            primitive.NilObjectID,
			StudentCode:        cert.StudentCode,
			FullName:           cert.Name,
			CertificateType:    cert.CertificateType,
			Course:             cert.Course,
			EducationType:      cert.EducationType,
			GPA:                cert.GPA,
			GraduationRank:     cert.GraduationRank,
			IssueDate:          cert.IssueDate,
			SerialNumber:       cert.SerialNumber,
			RegistrationNumber: cert.RegNo,
			EDiplomaFileLink:   pdfPath,
			EDiplomaFileHash:   hash,
			Signed:             false,
			Signature:          "",
			SignedAt:           time.Time{},
			OnBlockchain:       false,
			BlockchainTxID:     "",
			SignatureOfUni:     template.SignatureOfUni,
			SignatureOfMinEdu:  template.SignatureOfMinEdu,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		if err := s.repo.Save(ctx, ediploma); err != nil {
			log.Printf("Save failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		result = append(result, ediploma)
	}

	return result, nil
}

func (s *eDiplomaService) GetDiplomaPDF(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, int64, string, error) {
	diploma, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, 0, "", fmt.Errorf("diploma not found: %w", err)
	}
	if diploma.EDiplomaFileLink == "" {
		return nil, 0, "", fmt.Errorf("no file link for diploma")
	}

	log.Printf("[DEBUG] Getting file from MinIO - Bucket=%s, Object=%s", s.minioClient.Bucket, diploma.EDiplomaFileLink)

	statInfo, err := s.minioClient.Client.StatObject(
		ctx,
		s.minioClient.Bucket,
		diploma.EDiplomaFileLink,
		minio.StatObjectOptions{},
	)
	if err != nil {
		log.Printf("[ERROR] StatObject failed: %v", err)
		return nil, 0, "", fmt.Errorf("failed to stat file on MinIO: %w", err)
	}

	obj, err := s.minioClient.Client.GetObject(
		ctx,
		s.minioClient.Bucket,
		diploma.EDiplomaFileLink,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Printf("[ERROR] GetObject failed: %v", err)
		return nil, 0, "", fmt.Errorf("failed to get file from MinIO: %w", err)
	}

	return obj, statInfo.Size, statInfo.ContentType, nil
}

func (s *eDiplomaService) UploadLocalEDiplomas(ctx context.Context) []map[string]interface{} {
	localFolder := os.Getenv("EDIPLOMA_LOCAL_FOLDER")
	if localFolder == "" {
		return []map[string]interface{}{
			{"error": "EDIPLOMA_LOCAL_FOLDER not set in .env"},
		}
	}

	files, err := os.ReadDir(localFolder)
	if err != nil {
		return []map[string]interface{}{
			{"error": fmt.Sprintf("failed to read folder: %v", err)},
		}
	}

	var results []map[string]interface{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(localFolder, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			results = append(results, map[string]interface{}{
				"file":  file.Name(),
				"error": fmt.Sprintf("read failed: %v", err),
			})
			continue
		}

		hash := utils.ComputeSHA256(data) // SHA256 -> string

		// Upload file lên MinIO
		minioPath := fmt.Sprintf("ediplomas/%s", file.Name())
		if err := s.minioClient.UploadFile(ctx, minioPath, data, "application/pdf"); err != nil {
			results = append(results, map[string]interface{}{
				"file":  file.Name(),
				"error": fmt.Sprintf("upload failed: %v", err),
			})
			continue
		}

		// Lưu vào DB
		ediploma := &models.EDiploma{
			ID:               primitive.NewObjectID(),
			EDiplomaFileLink: minioPath,
			EDiplomaFileHash: hash,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := s.repo.Save(ctx, ediploma); err != nil {
			results = append(results, map[string]interface{}{
				"file":  file.Name(),
				"error": fmt.Sprintf("DB save failed: %v", err),
			})
			continue
		}

		results = append(results, map[string]interface{}{
			"file":   file.Name(),
			"hash":   hash,
			"link":   minioPath,
			"status": "uploaded",
		})
	}

	return results
}
func (s *eDiplomaService) ProcessZip(ctx context.Context, zipPath string) ([]*models.EDiploma, error) {
	// Thư mục tạm để giải nén
	extractDir := filepath.Join(os.TempDir(), fmt.Sprintf("unzipped_%d", time.Now().Unix()))
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create extract dir: %w", err)
	}

	if err := utils.Unzip(zipPath, extractDir); err != nil {
		return nil, fmt.Errorf("failed to unzip: %w", err)
	}

	files, err := os.ReadDir(extractDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read extracted folder: %w", err)
	}

	var results []*models.EDiploma

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(extractDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		hash := utils.ComputeSHA256(data)
		minioPath := fmt.Sprintf("ediplomas/%s", file.Name())

		// Upload file lên MinIO
		if err := s.minioClient.UploadFile(ctx, minioPath, data, "application/pdf"); err != nil {
			continue
		}

		// Lấy student_code từ tên file, loại bỏ phần .pdf
		studentCode := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		// Tìm bản ghi EDiploma hiện có
		ediploma, err := s.repo.FindByStudentCode(ctx, studentCode)
		if err != nil || ediploma == nil {
			continue
		}

		// Cập nhật các trường cần thiết
		updates := bson.M{
			"ediploma_file_hash": hash,
			"ediploma_file_link": minioPath,
			"data_encrypted":     true,
			"updated_at":         time.Now(),
		}

		if err := s.repo.UpdateFields(ctx, ediploma.ID, updates); err != nil {
			continue
		}

		// Lấy lại bản ghi sau khi update
		updatedEDiploma, err := s.repo.FindByID(ctx, ediploma.ID)
		if err != nil || updatedEDiploma == nil {
			continue
		}

		results = append(results, updatedEDiploma)
	}

	return results, nil
}

func (s *eDiplomaService) GenerateBulkEDiplomasZip(
	ctx context.Context,
	facultyIDStr, certificateType, course string,
	issued *bool,
	templateIDStr string,
) (string, error) {

	// 1. Convert templateID
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid template ID")
	}

	// 2. Lấy template
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil || template.HTMLContent == "" {
		return "", errors.New("template not found or empty")
	}

	// 3. Convert facultyID nếu có
	var facultyID primitive.ObjectID
	if facultyIDStr != "" {
		facultyID, err = primitive.ObjectIDFromHex(facultyIDStr)
		if err != nil {
			return "", fmt.Errorf("invalid faculty_id")
		}
	}

	// 4. Build dynamic filter
	filter := bson.M{}
	if !facultyID.IsZero() {
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

	// 5. Lấy danh sách eDiplomas
	ediplomas, err := s.repo.FindByDynamicFilter(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("failed to load eDiplomas: %w", err)
	}
	if len(ediplomas) == 0 {
		// Không tạo file zip, trả về thông báo rỗng
		return "", nil
	}

	// 6. Tạo thư mục tạm và xuất PDF
	tmpDir, err := os.MkdirTemp("", "ediplomas_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	var generatedFilePaths []string
	for _, ed := range ediplomas {
		user, _ := s.userRepo.GetUserByID(ctx, ed.UserID)
		university, _ := s.universityRepo.FindByID(ctx, ed.UniversityID)
		faculty, _ := s.facultyRepo.FindByID(ctx, ed.FacultyID)

		// Parse ngày sinh
		dobTime, _ := time.Parse("2006-01-02", user.DateOfBirth)
		data := map[string]interface{}{
			"SoHieu":         ed.SerialNumber,
			"SoVaoSo":        ed.RegistrationNumber,
			"HoTen":          user.FullName,
			"NgaySinh":       dobTime.Format("02/01/2006"),
			"TenTruong":      university.UniversityName,
			"Nganh":          faculty.FacultyName,
			"XepLoai":        ed.GraduationRank,
			"HinhThucDaoTao": ed.EducationType,
			"Khoa":           ed.Course,
			"NgayCap":        ed.IssueDate.Format("02/01/2006"),
		}

		renderedHTML, err := s.templateEngine.Render(template.HTMLContent, data)
		if err != nil {
			continue
		}

		pdfBytes, err := s.pdfGenerator.ConvertHTMLToPDF(renderedHTML)
		if err != nil {
			continue
		}

		fileName := fmt.Sprintf("%s.pdf", ed.StudentCode)
		filePath := filepath.Join(tmpDir, fileName)
		if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
			continue
		}

		// Update eDiploma
		ed.TemplateID = templateID
		ed.SignatureOfUni = template.SignatureOfUni
		ed.SignatureOfMinEdu = template.SignatureOfMinEdu
		ed.EDiplomaFileLink = filePath
		ed.EDiplomaFileHash = utils.ComputeSHA256(pdfBytes)
		ed.Issued = true
		ed.UpdatedAt = time.Now()
		_ = s.repo.Update(ctx, ed.ID, ed)

		generatedFilePaths = append(generatedFilePaths, filePath)
	}

	// 7. Tạo file zip từ các PDF
	zipFilePath := filepath.Join(tmpDir, "ediplomas.zip")
	if err := utils.CreateZipFromFiles(zipFilePath, generatedFilePaths); err != nil {
		return "", fmt.Errorf("failed to create zip: %w", err)
	}

	return zipFilePath, nil
}
