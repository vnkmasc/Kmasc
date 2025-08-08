package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MajorHandler struct {
	majorService service.MajorService
}

func NewMajorHandler(majorService service.MajorService) *MajorHandler {
	return &MajorHandler{majorService: majorService}
}

type CreateMajorRequest struct {
	MajorCode   string `json:"major_code" binding:"required"`
	MajorName   string `json:"major_name" binding:"required"`
	FacultyID   string `json:"faculty_id" binding:"required"`
	Description string `json:"description"`
	Quota       int    `json:"quota" binding:"required"`
}

func (h *MajorHandler) CreateMajor(c *gin.Context) {
	var req CreateMajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Chuyển facultyID từ string sang ObjectID
	facultyID, err := primitive.ObjectIDFromHex(req.FacultyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty_id format"})
		return
	}

	// Chuyển universityID từ string sang ObjectID
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university_id format in token"})
		return
	}

	major := models.Major{
		ID:           primitive.NewObjectID(),
		MajorCode:    req.MajorCode,
		MajorName:    req.MajorName,
		FacultyID:    facultyID,
		UniversityID: universityID,
		Description:  req.Description,
		Quota:        req.Quota,
		DiplomaCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.majorService.CreateMajor(c.Request.Context(), &major); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chuyên ngành: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo chuyên ngành thành công", "data": major})
}

func (h *MajorHandler) GetMajorsByFaculty(c *gin.Context) {
	facultyIDHex := c.Param("faculty_id")
	facultyID, err := primitive.ObjectIDFromHex(facultyIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "faculty_id không hợp lệ"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "university_id trong token không hợp lệ"})
		return
	}

	majors, err := h.majorService.GetMajorsByFaculty(c.Request.Context(), universityID, facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách chuyên ngành: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": majors})
}

func (h *MajorHandler) DeleteMajor(c *gin.Context) {
	majorIDHex := c.Param("id")
	majorID, err := primitive.ObjectIDFromHex(majorIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID chuyên ngành không hợp lệ"})
		return
	}

	if err := h.majorService.DeleteMajor(c.Request.Context(), majorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xoá chuyên ngành thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xoá chuyên ngành thành công"})
}
