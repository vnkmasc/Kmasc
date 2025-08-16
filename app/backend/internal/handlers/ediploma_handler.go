package handlers

import (
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiplomaHandler struct {
	universityService service.UniversityService
	ediplomaService   service.EDiplomaService
	minioClient       *database.MinioClient
}

func NewEDiplomaHandler(
	ediplomaService service.EDiplomaService,
) *EDiplomaHandler {
	return &EDiplomaHandler{
		ediplomaService: ediplomaService,
	}
}

type generateEDiplomaRequest struct {
	CertificateID string `json:"certificate_id"`
	TemplateID    string `json:"template_id"`
}

// Handler
func (h *EDiplomaHandler) SearchEDiplomas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var issued *bool
	if issuedStr := c.Query("issued"); issuedStr != "" {
		v, _ := strconv.ParseBool(issuedStr)
		issued = &v
	}

	filters := models.EDiplomaSearchFilter{
		FacultyID:       c.Query("faculty_id"), // dùng trực tiếp faculty_id từ query
		CertificateType: c.Query("certificate_type"),
		Course:          c.Query("course"),
		Issued:          issued,
		Page:            page,
		PageSize:        pageSize,
	}

	dtoList, total, err := h.ediplomaService.SearchEDiplomaDTOs(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{
		"data":       dtoList,
		"page":       page,
		"page_size":  pageSize,
		"total":      total,
		"total_page": totalPage,
	})
}

func (h *EDiplomaHandler) GenerateEDiploma(c *gin.Context) {
	var req generateEDiplomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ediploma, err := h.ediplomaService.GenerateEDiploma(
		c.Request.Context(),
		req.CertificateID,
		req.TemplateID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ediploma)
}

type generateBulkEDiplomaRequest struct {
	FacultyID  string `json:"faculty_id" binding:"required"`
	TemplateID string `json:"template_id" binding:"required"`
}

func (h *EDiplomaHandler) GenerateBulkEDiplomas(c *gin.Context) {
	var req generateBulkEDiplomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ediplomas, err := h.ediplomaService.GenerateBulkEDiplomas(c.Request.Context(), req.FacultyID, req.TemplateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ediplomas)
}

func (h *EDiplomaHandler) UploadLocalEDiplomas(c *gin.Context) {
	results := h.ediplomaService.UploadLocalEDiplomas(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"total_files": len(results),
		"results":     results,
	})
}

func (h *EDiplomaHandler) UploadEDiplomasZip(c *gin.Context) {

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := claimsRaw.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing zip file"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot open uploaded file"})
		return
	}
	defer file.Close()

	tempZipPath := filepath.Join(os.TempDir(), fileHeader.Filename)
	out, err := os.Create(tempZipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save temp zip"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file"})
		return
	}

	// Gọi service xử lý zip và lấy các bản ghi EDiploma đã update
	updatedDiplomas, err := h.ediplomaService.ProcessZip(c.Request.Context(), tempZipPath, universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_files": len(updatedDiplomas),
		"results":     updatedDiplomas,
	})
}

func (h *EDiplomaHandler) GenerateBulkEDiplomasZip(c *gin.Context) {
	var req struct {
		FacultyID       string `form:"faculty_id" json:"faculty_id"`
		CertificateType string `form:"certificate_type" json:"certificate_type"` // optional
		Course          string `form:"course" json:"course"`                     // optional
		Issued          *bool  `form:"issued" json:"issued"`                     // optional
		TemplateID      string `form:"template_id" json:"template_id" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	zipFilePath, err := h.ediplomaService.GenerateBulkEDiplomasZip(
		c.Request.Context(),
		req.FacultyID,
		req.CertificateType,
		req.Course,
		req.Issued,
		req.TemplateID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if zipFilePath == "" {
		// Không có văn bằng nào được cấp
		c.JSON(http.StatusOK, gin.H{
			"message": "Không có văn bằng nào được cấp theo bộ lọc này",
		})
		return
	}

	// Nếu có zip, gửi file
	c.FileAttachment(zipFilePath, "ediplomas.zip")
	_ = os.Remove(zipFilePath)
}

func (h *EDiplomaHandler) GetEDiplomaByID(c *gin.Context) {
	ctx := c.Request.Context()
	ediplomaIDHex := c.Param("id")
	ediplomaID, err := primitive.ObjectIDFromHex(ediplomaIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ediploma ID"})
		return
	}

	dto, err := h.ediplomaService.GetEDiplomaDTOByID(ctx, ediplomaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "EDiploma not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}

func (h *EDiplomaHandler) ViewEDiplomaFile(c *gin.Context) {
	// Lấy claims từ token
	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims, ok := claimsRaw.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims"})
		return
	}
	universityID, _ := primitive.ObjectIDFromHex(claims.UniversityID)

	// Lấy ediplomaID từ URL
	ediplomaIDHex := c.Param("id")
	ediplomaID, err := primitive.ObjectIDFromHex(ediplomaIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ediploma ID"})
		return
	}

	// Lấy stream và content type từ service
	stream, contentType, err := h.ediplomaService.GetEDiplomaFile(c.Request.Context(), ediplomaID, universityID)
	if err != nil {
		if err.Error() == "EDiploma not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	// Trả file về client
	c.DataFromReader(http.StatusOK, -1, contentType, stream, nil)
}
