package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
// 	templateIDHex := c.Param("id")
// 	templateID, err := primitive.ObjectIDFromHex(templateIDHex)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
// 		return
// 	}

// 	claimsRaw, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	claims, ok := claimsRaw.(*utils.CustomClaims)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
// 		return
// 	}

// 	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID in token"})
// 		return
// 	}

// 	var req struct {
// 		Name        string `json:"name"`
// 		Description string `json:"description"`
// 		HTMLContent string `json:"html_content"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
// 		return
// 	}

// 	updatedTemplate, err := h.templateService.UpdateTemplate(
// 		c.Request.Context(),
// 		templateID,
// 		universityID,
// 		req.Name,
// 		req.Description,
// 		req.HTMLContent,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":  "Template updated successfully",
// 		"template": updatedTemplate,
// 	})
// }

func (h *TemplateHandler) VerifyTemplatesByFaculty(c *gin.Context) {
	facultyIDHex := c.Param("faculty_id")
	facultyID, err := primitive.ObjectIDFromHex(facultyIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

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

	err = h.templateService.VerifyTemplatesByFaculty(c.Request.Context(), universityID, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All templates verified for faculty successfully"})
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

	// Gọi service, lấy template đã ký
	template, err := h.templateService.SignTemplateByID(c.Request.Context(), universityID, templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully signed template: %s", template.Name),
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

func (h *TemplateHandler) SignTemplatesByMinEdu(c *gin.Context) {
	universityIDStr := c.Param("university_id")
	universityID, err := primitive.ObjectIDFromHex(universityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	count, err := h.templateService.SignAllTemplatesByMinEdu(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign templates by Ministry of Education"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          fmt.Sprintf("Successfully signed %d templates by Ministry of Education", count),
		"signed_templates": count,
	})
}

func parseMinioURL(urlStr string) (bucket, objectPath string, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	trimmedPath := strings.TrimPrefix(u.Path, "/")
	parts := strings.SplitN(trimmedPath, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid MinIO file URL")
	}

	return parts[0], parts[1], nil
}
