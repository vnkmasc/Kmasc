package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiplomaService interface {
	GetByID(ctx context.Context, id string) (*models.EDiploma, error)
	GenerateEDiploma(ctx context.Context, certificateIDStr, templateIDStr string) (*models.EDiploma, error)
	GenerateBulkEDiplomas(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error)
	GetEDiplomasByFaculty(ctx context.Context, facultyID string) ([]*models.EDiplomaDTO, error)
	SearchEDiplomaDTOs(ctx context.Context, filter models.EDiplomaSearchFilter) ([]*models.EDiplomaDTO, int64, error)
	GenerateBulkEDiplomasLocal(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error)
	UploadLocalEDiplomas(ctx context.Context) []map[string]interface{}
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

func (s *eDiplomaService) GetEDiplomasByFaculty(ctx context.Context, facultyIDStr string) ([]*models.EDiplomaDTO, error) {
	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid faculty ID")
	}

	ediplomas, err := s.repo.GetByFacultyID(ctx, facultyID)
	if err != nil {
		return nil, err
	}

	faculty, err := s.facultyRepo.FindByID(ctx, facultyID)
	if err != nil {
		return nil, err
	}

	university, err := s.universityRepo.FindByID(ctx, faculty.UniversityID)
	if err != nil {
		return nil, err
	}

	var result []*models.EDiplomaDTO
	for _, ed := range ediplomas {
		var major *models.Major
		if ed.MajorID != primitive.NilObjectID {
			major, _ = s.majorRepo.GetByID(ctx, ed.MajorID)
		}

		dto := mapper.MapEDiplomaToDTO(ed, university, faculty, major)
		result = append(result, dto)
	}

	return result, nil
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

// --- Main function ---
func (s *eDiplomaService) GenerateEDiploma(ctx context.Context, certificateIDStr, templateIDStr string) (*models.EDiploma, error) {
	log.Printf("[DEBUG] Parsing certificateID: %s", certificateIDStr)
	certificateID, err := primitive.ObjectIDFromHex(certificateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate ID")
	}

	log.Printf("[DEBUG] Parsing templateID: %s", templateIDStr)
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID")
	}

	log.Printf("[DEBUG] Getting certificate by ID: %s", certificateID.Hex())
	cert, err := s.certificateRepo.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, fmt.Errorf("certificate not found")
	}

	log.Printf("[DEBUG] Getting template by ID: %s", templateID.Hex())
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}
	log.Printf("[DEBUG] Template file link: %s", template.FileLink)

	bucket, objectPath, err := parseMinioURL(template.FileLink)
	if err != nil {
		log.Printf("[ERROR] Failed to parse MinIO URL: %v", err)
		return nil, fmt.Errorf("invalid template file URL: %w", err)
	}
	log.Printf("[DEBUG] Parsed MinIO URL - bucket: %s, object: %s", bucket, objectPath)

	htmlContent, err := s.minioClient.DownloadFile(ctx, bucket, objectPath)
	if err != nil {
		log.Printf("[ERROR] Failed to download HTML from MinIO: objectPath=%s, err=%v", objectPath, err)
		return nil, fmt.Errorf("failed to download template HTML from MinIO: %w", err)
	}
	log.Printf("[DEBUG] Downloaded HTML content length: %d", len(htmlContent))
	log.Printf("[DEBUG] Getting user by ID: %s", cert.UserID.Hex())
	user, err := s.userRepo.GetUserByID(ctx, cert.UserID)
	if err != nil {
		log.Printf("[ERROR] Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("[DEBUG] Getting university by ID: %s", cert.UniversityID.Hex())
	university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
	if err != nil {
		log.Printf("[ERROR] Failed to get university: %v", err)
		return nil, fmt.Errorf("failed to get university: %w", err)
	}
	dob, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		log.Printf("[ERROR] Failed to parse DateOfBirth: %v", err)
		// fallback hoặc trả về lỗi tùy logic
	}

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

	log.Printf("[DEBUG] Template data prepared: %+v", data)

	renderedHTML, err := s.templateEngine.Render(string(htmlContent), data)
	if err != nil {
		log.Printf("[ERROR] Failed to render HTML: %v", err)
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	pdfBytes, err := s.pdfGenerator.ConvertHTMLToPDF(renderedHTML)
	if err != nil {
		log.Printf("[ERROR] Failed to convert HTML to PDF: %v", err)
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	log.Printf("[DEBUG] Generated PDF size: %d bytes", len(pdfBytes))

	hash := utils.ComputeSHA256(pdfBytes)
	pdfPath := fmt.Sprintf("ediplomas/%s.pdf", primitive.NewObjectID().Hex())

	log.Printf("[DEBUG] Uploading PDF to MinIO at path: %s", pdfPath)
	if err := s.minioClient.UploadFile(ctx, pdfPath, pdfBytes, "application/pdf"); err != nil {
		log.Printf("[ERROR] Failed to upload PDF to MinIO: %v", err)
		return nil, fmt.Errorf("failed to upload PDF to MinIO: %w", err)
	}
	log.Printf("[DEBUG] Locking template after diploma generation")
	err = s.templateRepo.UpdateIsLocked(ctx, templateID, true)
	if err != nil {
		log.Printf("[ERROR] Failed to lock template: %v", err)
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
		FileLink:           pdfPath,
		FileHash:           hash,
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

	log.Printf("[DEBUG] Saving EDiploma record to DB")
	if err := s.repo.Save(ctx, ediploma); err != nil {
		log.Printf("[ERROR] Failed to save EDiploma: %v", err)
		return nil, fmt.Errorf("failed to save EDiploma: %w", err)
	}

	log.Printf("[DEBUG] EDiploma generated successfully: ID=%s", ediploma.ID.Hex())
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

		dto := mapper.MapEDiplomaToDTO(ed, university, faculty, major)
		dtoList = append(dtoList, dto)
	}

	return dtoList, total, nil
}

func (s *eDiplomaService) GenerateBulkEDiplomas(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error) {
	var result []*models.EDiploma

	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid faculty ID")
	}
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID")
	}

	// Load template
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}
	if template.FacultyID != facultyID {
		return nil, errors.New("template does not belong to the given faculty")
	}

	// Load all certificates for this faculty
	certificates, err := s.certificateRepo.FindByFacultyID(ctx, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates: %w", err)
	}

	// Download template HTML once
	bucket, objectPath, err := parseMinioURL(template.FileLink)
	if err != nil {
		return nil, fmt.Errorf("invalid template file URL: %w", err)
	}
	htmlContent, err := s.minioClient.DownloadFile(ctx, bucket, objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to download template HTML: %w", err)
	}

	for _, cert := range certificates {
		// Load university
		university, err := s.universityRepo.FindByID(ctx, cert.UniversityID)
		if err != nil {
			log.Printf("Load university failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		// Load user
		user, err := s.userRepo.GetUserByID(ctx, cert.UserID) // Bạn cần đảm bảo repo có hàm này
		if err != nil {
			log.Printf("Load user failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		dobTime, err := time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			log.Printf("Invalid date format for user %s: %v", user.ID.Hex(), err)
			continue // bỏ qua nếu date format sai
		}
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

		renderedHTML, err := s.templateEngine.Render(string(htmlContent), data)
		if err != nil {
			log.Printf("Render failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		pdfBytes, err := s.pdfGenerator.ConvertHTMLToPDF(renderedHTML)
		if err != nil {
			log.Printf("PDF generation failed for cert %s: %v", cert.ID.Hex(), err)
			continue
		}

		hash := utils.ComputeSHA256(pdfBytes)
		pdfPath := fmt.Sprintf("ediplomas/%s.pdf", primitive.NewObjectID().Hex())

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
			FileLink:           pdfPath,
			FileHash:           hash,
			Signed:             false,
			Signature:          "",
			SignedAt:           time.Time{},
			OnBlockchain:       false,
			BlockchainTxID:     "",
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

func (s *eDiplomaService) GenerateBulkEDiplomasLocal(ctx context.Context, facultyIDStr, templateIDStr string) ([]*models.EDiploma, error) {
	var result []*models.EDiploma

	// Chuyển đổi ID
	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid faculty ID")
	}
	templateID, err := primitive.ObjectIDFromHex(templateIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template ID")
	}

	// Lấy template
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}
	if template.FacultyID != facultyID {
		return nil, errors.New("template does not belong to the given faculty")
	}

	// Lấy certificates của faculty
	certificates, err := s.certificateRepo.FindByFacultyID(ctx, facultyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates: %w", err)
	}

	// Tải template HTML từ MinIO
	bucket, objectPath, err := parseMinioURL(template.FileLink)
	if err != nil {
		return nil, fmt.Errorf("invalid template file URL: %w", err)
	}
	htmlContent, err := s.minioClient.DownloadFile(ctx, bucket, objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to download template HTML: %w", err)
	}

	localFolder := os.Getenv("EDIPLOMA_LOCAL_FOLDER")
	if localFolder == "" {
		return nil, fmt.Errorf("EDIPLOMA_LOCAL_FOLDER not set in environment")
	}

	if err := os.MkdirAll(localFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create local folder: %w", err)
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

		// Data render template
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
		renderedHTML, err := s.templateEngine.Render(string(htmlContent), data)
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

		hash := utils.ComputeSHA256(pdfBytes)

		// Lưu PDF xuống local
		localFileName := fmt.Sprintf("%s.pdf", primitive.NewObjectID().Hex())
		localFilePath := filepath.Join(localFolder, localFileName)
		if err := os.WriteFile(localFilePath, pdfBytes, 0644); err != nil {
			log.Printf("Save PDF failed for cert %s: %v", cert.ID.Hex(), err)
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
			FileLink:           localFilePath,
			FileHash:           hash,
			Signed:             false,
			Signature:          "",
			SignedAt:           time.Time{},
			OnBlockchain:       false,
			BlockchainTxID:     "",
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
			ID:        primitive.NewObjectID(),
			FileLink:  minioPath,
			FileHash:  hash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
