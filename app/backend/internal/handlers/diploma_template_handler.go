package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateHandler struct {
	templateService service.TemplateService
	minioClient     *database.MinioClient
	facultyService  service.FacultyService
}

func NewTemplateHandler(
	templateService service.TemplateService,
	minioClient *database.MinioClient,
	facultyService service.FacultyService,
) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		minioClient:     minioClient,
		facultyService:  facultyService,
	}
}

func (h *TemplateHandler) GetTemplateByID(c *gin.Context) {
	id := c.Param("id")

	template, err := h.templateService.GetTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID or template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": template,
	})
}

// POST /templates
func (h *TemplateHandler) GetTemplatesByFaculty(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id in token"})
		return
	}

	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty_id"})
		return
	}

	templates, err := h.templateService.GetTemplatesByFaculty(c.Request.Context(), universityID, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": templates,
	})
}
func (h *TemplateHandler) GetTemplatesByFacultyAndUniversity(c *gin.Context) {
	universityIDStr := c.Param("university_id")
	facultyIDStr := c.Param("faculty_id")

	universityID, err := primitive.ObjectIDFromHex(universityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id"})
		return
	}

	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty_id"})
		return
	}

	templates, err := h.templateService.GetTemplatesByFaculty(c.Request.Context(), universityID, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": templates,
	})
}
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		FacultyID        string `json:"faculty_id" binding:"required"`
		TemplateSampleID string `json:"template_sample_id" binding:"required"`
		Name             string `json:"name" binding:"required"` // tên truyền từ request
		Description      string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body", "details": err.Error()})
		return
	}

	// Lấy UniversityID từ token
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id in token"})
		return
	}

	facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty_id"})
		return
	}

	templateSampleID, err := primitive.ObjectIDFromHex(req.TemplateSampleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template_sample_id"})
		return
	}

	template, err := h.templateService.CreateTemplate(
		c.Request.Context(),
		universityID,
		facultyID,
		templateSampleID,
		req.Name,        // truyền name từ request
		req.Description, // truyền description từ request
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Template created successfully from sample",
		"template": template,
	})
}

func (h *TemplateHandler) SignTemplatesByFaculty(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")

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

	facultyID, err := primitive.ObjectIDFromHex(facultyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	count, err := h.templateService.SignTemplatesByFaculty(c.Request.Context(), universityID, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully signed %d templates", count),
	})
}

type SignTemplateRequest struct {
	Signature string `json:"signature" binding:"required"`
}

func (h *TemplateHandler) SignTemplateByID(c *gin.Context) {
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

	templateIDHex := c.Param("template_id")
	templateID, err := primitive.ObjectIDFromHex(templateIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req SignTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature is required"})
		return
	}

	// Gọi service để lưu signature
	template, err := h.templateService.SaveClientSignature(c.Request.Context(), universityID, templateID, req.Signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã lưu chữ ký cho mẫu %s", template.Name),
	})
}

func (h *TemplateHandler) SignAllPendingTemplatesOfUniversity(c *gin.Context) {
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

	count, err := h.templateService.SignAllPendingTemplatesOfUniversity(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully signed %d templates", count),
	})
}

func (h *TemplateHandler) SignTemplateByMinEdu(c *gin.Context) {
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

	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	templateIDHex := c.Param("template_id")
	templateID, err := primitive.ObjectIDFromHex(templateIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req SignTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature is required"})
		return
	}

	template, err := h.templateService.SaveMinEduSignature(
		c.Request.Context(),
		templateID,
		req.Signature,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Đã lưu chữ ký Bộ Giáo dục cho mẫu %s", template.Name),
	})
}

func (h *TemplateHandler) UpdateDiplomaTemplate(c *gin.Context) {
	templateIDHex := c.Param("template_id")
	templateID, err := primitive.ObjectIDFromHex(templateIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var req models.UpdateDiplomaTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.templateService.UpdateDiplomaTemplate(
		c.Request.Context(),
		templateID,
		req,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chỉ trả message
	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật mẫu văn bằng thành công",
	})
}
