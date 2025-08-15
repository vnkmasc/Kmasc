package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateSampleHandler struct {
	service *service.TemplateSampleService
}

// Constructor
func NewTemplateSampleHandler(service *service.TemplateSampleService) *TemplateSampleHandler {
	return &TemplateSampleHandler{service: service}
}
func (h *TemplateSampleHandler) GetTemplateSampleByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template sample ID"})
		return
	}

	sample, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template sample not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": sample,
	})
}

type CreateTemplateSampleRequest struct {
	Name        string `json:"name" binding:"required"`
	HTMLContent string `json:"html_content" binding:"required"`
}

func (h *TemplateSampleHandler) CreateTemplateSample(c *gin.Context) {
	var req CreateTemplateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON body",
			"details": err.Error(),
		})
		return
	}

	val, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id in token"})
		return
	}

	sample := &models.TemplateSample{
		Name:         req.Name,
		HTMLContent:  req.HTMLContent,
		UniversityID: universityID, // gán từ token
	}

	id, err := h.service.Create(c.Request.Context(), sample)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Gán ID để trả về
	sample.ID = id

	c.JSON(http.StatusOK, gin.H{
		"message": "Template sample created successfully",
		"data":    sample,
	})
}

type UpdateTemplateSampleRequest struct {
	Name        string `json:"name" binding:"required"`
	HTMLContent string `json:"html_content" binding:"required"`
}

func (h *TemplateSampleHandler) UpdateTemplateSample(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template sample ID"})
		return
	}

	var req UpdateTemplateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body", "details": err.Error()})
		return
	}

	sample := &models.TemplateSample{
		ID:          id,
		Name:        req.Name,
		HTMLContent: req.HTMLContent,
	}

	if err := h.service.Update(c.Request.Context(), sample); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Template sample updated successfully",
	})
}
func (h *TemplateSampleHandler) GetAllTemplateSamples(c *gin.Context) {
	val, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id in token"})
		return
	}

	samples, err := h.service.GetAllVisible(c.Request.Context(), universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, s := range samples {
		result = append(result, gin.H{
			"id":   s.ID.Hex(),
			"name": s.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    result,
	})
}
