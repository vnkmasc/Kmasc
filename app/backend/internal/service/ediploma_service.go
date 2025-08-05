package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	urlpkg "net/url"
	"strings"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiplomaService interface {
	GetByID(ctx context.Context, id string) (*models.EDiploma, error)
	GenerateEDiploma(ctx context.Context, certificateIDStr, templateIDStr string) (*models.EDiploma, error)
}

type eDiplomaService struct {
	facultyRepo     repository.FacultyRepository
	repo            repository.EDiplomaRepository
	certificateRepo repository.CertificateRepository
	templateRepo    repository.TemplateRepository
	minioClient     *database.MinioClient
	templateEngine  *models.TemplateEngine
	pdfGenerator    *utils.PDFGenerator
}

func NewEDiplomaService(
	facultyRepo repository.FacultyRepository,
	repo repository.EDiplomaRepository,
	certificateRepo repository.CertificateRepository,
	templateRepo repository.TemplateRepository,
	minioClient *database.MinioClient,
	templateEngine *models.TemplateEngine,
	pdfGenerator *utils.PDFGenerator,
) *eDiplomaService {
	return &eDiplomaService{
		facultyRepo:     facultyRepo,
		repo:            repo,
		certificateRepo: certificateRepo,
		templateRepo:    templateRepo,
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

	data := map[string]interface{}{
		"SoHieu":         cert.SerialNumber,
		"SoVaoSo":        cert.RegNo,
		"HoTen":          cert.Name,
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

	log.Printf("[DEBUG] Saving EDiploma record to DB")
	if err := s.repo.Save(ctx, ediploma); err != nil {
		log.Printf("[ERROR] Failed to save EDiploma: %v", err)
		return nil, fmt.Errorf("failed to save EDiploma: %w", err)
	}

	log.Printf("[DEBUG] EDiploma generated successfully: ID=%s", ediploma.ID.Hex())
	return ediploma, nil
}
