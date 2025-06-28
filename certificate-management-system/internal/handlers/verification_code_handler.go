package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationHandler struct {
	verificationService service.VerificationService
	userService         service.UserService
	certificateService  service.CertificateService
	minioClient         *database.MinioClient
}

func NewVerificationHandler(
	verificationService service.VerificationService,
	userService service.UserService,
	certificateService service.CertificateService,
	minioClient *database.MinioClient,

) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
		userService:         userService,
		certificateService:  certificateService,
		minioClient:         minioClient,
	}
}

func (h *VerificationHandler) CreateVerificationCode(c *gin.Context) {
	var req models.CreateVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	claims, ok := c.Request.Context().Value(utils.ClaimsContextKey).(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID người dùng không hợp lệ"})
		return
	}

	expiredAt := time.Now().Add(time.Duration(req.DurationMinutes) * time.Minute)

	code := &models.VerificationCode{
		UserID:       userID,
		CanViewScore: req.CanViewScore,
		CanViewData:  req.CanViewData,
		CanViewFile:  req.CanViewFile,
		ExpiredAt:    expiredAt,
	}

	err = h.verificationService.CreateVerificationCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo mã xác minh"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       code.Code,
		"expired_at": code.ExpiredAt,
		"can_view": gin.H{
			"score": code.CanViewScore,
			"data":  code.CanViewData,
			"file":  code.CanViewFile,
		},
	})
}
func (h *VerificationHandler) GetMyCodes(c *gin.Context) {
	claims, ok := c.Request.Context().Value(utils.ClaimsContextKey).(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID không hợp lệ"})
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)

	codes, total, err := h.verificationService.GetCodesByUser(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy mã xác minh"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       codes,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + pageSize - 1) / pageSize,
	})
}

func (h *VerificationHandler) VerifyCode(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	ctx := c.Request.Context()
	_, certResp, err := h.verificationService.VerifyCode(ctx, req.Code, req.ViewType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.ViewType {
	case "data", "score":
		c.JSON(http.StatusOK, gin.H{
			"data": certResp,
		})

	case "file":
		if certResp == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy dữ liệu văn bằng"})
			return
		}
		if certResp.Path == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng chưa có file"})
			return
		}

		object, err := h.minioClient.Client.GetObject(ctx, h.minioClient.Bucket, certResp.Path, minio.GetObjectOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không đọc được file từ MinIO"})
			return
		}
		defer object.Close()

		fileData, err := io.ReadAll(object)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc nội dung file"})
			return
		}

		contentType := http.DetectContentType(fileData)
		c.DataFromReader(http.StatusOK, int64(len(fileData)), contentType, bytes.NewReader(fileData), nil)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại dữ liệu không hợp lệ"})
	}
}
