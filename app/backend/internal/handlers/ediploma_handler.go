package handlers

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EDiplomaHandler struct {
	ediplomaService service.EDiplomaService
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

func (h *EDiplomaHandler) SearchEDiplomas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filters := models.EDiplomaSearchFilter{
		StudentCode:     c.Query("student_code"),
		FacultyCode:     c.Query("faculty_code"),
		CertificateType: c.Query("certificate_type"),
		Course:          c.Query("course"),
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

	io.Copy(out, file)

	// Gọi service xử lý zip
	results, err := h.ediplomaService.ProcessZip(c.Request.Context(), tempZipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_files": len(results),
		"results":     results,
	})
}

func (h *EDiplomaHandler) ViewEDiploma(c *gin.Context) {
	ctx := c.Request.Context()
	diplomaID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(diplomaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diploma ID"})
		return
	}

	obj, size, contentType, err := h.ediplomaService.GetDiplomaPDF(ctx, objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer obj.Close()

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filepath.Base(diplomaID+".pdf")))
	c.Header("Content-Length", fmt.Sprintf("%d", size))

	if _, err := io.Copy(c.Writer, obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error streaming file"})
		return
	}
}

func (h *EDiplomaHandler) GetEDiplomasByFaculty(c *gin.Context) {
	facultyID := c.Param("faculty_id")
	if facultyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faculty_id is required"})
		return
	}

	ediplomas, err := h.ediplomaService.GetEDiplomasByFaculty(c.Request.Context(), facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ediplomas})
}
func (h *EDiplomaHandler) GetEDiplomaByID(c *gin.Context) {
	id := c.Param("id")

	dto, err := h.ediplomaService.GetEDiplomaDTOByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "EDiploma not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dto)
}

func (h *EDiplomaHandler) GenerateBulkEDiplomasZip(c *gin.Context) {
	var req struct {
		FacultyID  string `json:"faculty_id" binding:"required"`
		TemplateID string `json:"template_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	zipFilePath, err := h.ediplomaService.GenerateBulkEDiplomasZip(
		c.Request.Context(),
		req.FacultyID,
		req.TemplateID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Đọc file ZIP để trả về client
	c.FileAttachment(zipFilePath, "ediplomas.zip")

	// Optional: Xóa file ZIP sau khi gửi xong
	_ = os.Remove(zipFilePath)
}
