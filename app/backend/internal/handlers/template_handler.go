package handlers

import (
	"fmt"
	"io"
	"net/http"
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
	name := c.PostForm("name")
	description := c.PostForm("description")
	facultyIDStr := c.PostForm("faculty_id")

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

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	template, err := h.templateService.CreateTemplate(c.Request.Context(), name, description, universityID, facultyID, header.Filename, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Template created successfully",
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

// // PUT /templates/:id
// func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
// 	id := c.Param("id")

// 	var req struct {
// 		Name        string `json:"name"`
// 		Description string `json:"description"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	updateData := &models.DiplomaTemplate{
// 		Name:        req.Name,
// 		Description: req.Description,
// 	}

// 	if err := h.templateService.UpdateTemplate(c.Request.Context(), id, updateData); err != nil {
// 		if err.Error() == "template is locked and cannot be modified" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Template is locked"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{"message": "Template updated"})
//	}
func (h *TemplateHandler) GetTemplateFile(c *gin.Context) {

	templateID := c.Param("id")
	template, err := h.templateService.GetTemplateByID(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Parse bucket và object từ FileLink
	// Ví dụ: http://host:9000/bucket/object-path
	fileLink := template.FileLink
	parts := strings.SplitN(strings.TrimPrefix(fileLink, "http://"), "/", 2)
	if len(parts) != 2 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file link"})
		return
	}
	bucketAndHost := parts[1]
	idx := strings.Index(bucketAndHost, "/")
	if idx == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file link"})
		return
	}
	bucket := bucketAndHost[:idx]
	object := bucketAndHost[idx+1:]

	// Đọc file từ MinIO
	obj, err := h.minioClient.GetObject(c, bucket, object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file from MinIO"})
		return
	}
	defer obj.Close()

	// Đọc toàn bộ file
	fileBytes, err := io.ReadAll(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Trả về file (ví dụ PDF)
	c.Data(http.StatusOK, "text/html; charset=utf-8", fileBytes)
}
