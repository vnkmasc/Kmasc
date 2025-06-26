package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
)

type UniversityHandler struct {
	universityService service.UniversityService
}

func NewUniversityHandler(s service.UniversityService) *UniversityHandler {
	return &UniversityHandler{universityService: s}
}

func (h *UniversityHandler) CreateUniversity(c *gin.Context) {
	var req models.CreateUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("BindJSON error: %v", err)
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}
	err := h.universityService.CreateUniversity(c.Request.Context(), &req)
	switch err {
	case common.ErrUniversityNameExists:
		c.JSON(400, gin.H{"error": "Tên trường đã tồn tại"})
		return
	case common.ErrUniversityEmailDomainExists:
		c.JSON(400, gin.H{"error": "Tên miền email đã tồn tại"})
		return
	case common.ErrUniversityCodeExists:
		c.JSON(400, gin.H{"error": "Mã trường đã tồn tại"})
		return
	}
	c.JSON(200, gin.H{"message": "Đã gửi yêu cầu sử dụng hệ thống, chờ admin phê duyệt"})
}

func (h *UniversityHandler) ApproveOrRejectUniversity(c *gin.Context) {
	var req models.ApproveOrRejectUniversityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("BindJSON error: %v", err)
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	err := h.universityService.ApproveOrRejectUniversity(c.Request.Context(), req.UniversityID, req.Action)
	if err != nil {
		switch err {
		case common.ErrUniversityNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy trường"})
		case common.ErrUniversityAlreadyApproved:
			c.JSON(http.StatusConflict, gin.H{"error": "Trường này đã được phê duyệt"})
		case common.ErrAccountUniversityAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Trường đã có tài khoản quản trị"})
		case common.ErrUniversityCodeExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Mã trường đã tồn tại"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống: " + err.Error()})
		}
		return
	}

	switch req.Action {
	case "approve":
		c.JSON(http.StatusOK, gin.H{"message": "Trường đã được phê duyệt và đã gửi tài khoản qua email"})
	case "reject":
		c.JSON(http.StatusOK, gin.H{"message": "Đã từ chối trường sử dụng hệ thống"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hành động không hợp lệ"})
	}
}

func (h *UniversityHandler) GetAllUniversities(c *gin.Context) {
	universities, err := h.universityService.GetAllUniversities(c.Request.Context())
	if err != nil {
		log.Printf("Error getting universities: %v", err)
		c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		return
	}
	var resp []models.UniversityResponse
	for _, u := range universities {
		resp = append(resp, models.UniversityResponse{
			ID:             u.ID.Hex(),
			UniversityName: u.UniversityName,
			UniversityCode: u.UniversityCode,
			EmailDomain:    u.EmailDomain,
			Address:        u.Address,
			Status:         u.Status,
		})
	}

	c.JSON(200, gin.H{"data": resp})
}

func (h *UniversityHandler) GetUniversities(c *gin.Context) {
	status := c.Query("status")

	universities, err := h.universityService.GetUniversitiesByStatus(c.Request.Context(), status)
	if err != nil {
		log.Printf("Error getting universities by status: %v", err)
		c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	var resp []models.UniversityResponse
	for _, u := range universities {
		resp = append(resp, models.UniversityResponse{
			ID:             u.ID.Hex(),
			UniversityName: u.UniversityName,
			UniversityCode: u.UniversityCode,
			EmailDomain:    u.EmailDomain,
			Address:        u.Address,
			Status:         u.Status,
		})
	}

	c.JSON(200, gin.H{"data": resp})
}
