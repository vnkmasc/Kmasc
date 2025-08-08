package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
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

func (h *EDiplomaHandler) GenerateEDiploma(c *gin.Context) {
	var req generateEDiplomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ediploma, err := h.ediplomaService.GenerateEDiploma(c.Request.Context(), req.CertificateID, req.TemplateID)
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

	c.JSON(http.StatusOK, ediplomas)
}
