package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainHandler struct {
	BlockchainSvc service.BlockchainService
}

func NewBlockchainHandler(blockchainSvc service.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{BlockchainSvc: blockchainSvc}
}

func (h *BlockchainHandler) PushCertificateToChain(c *gin.Context) {
	certIDStr := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(certIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	txID, err := h.BlockchainSvc.PushCertificateToChain(c.Request.Context(), certID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrCertificateNotSigned),
			errors.Is(err, common.ErrCertificateNoFile),
			errors.Is(err, common.ErrCertificateAlreadyOnChain),
			errors.Is(err, common.ErrCertificateMissingHash):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Không thể đưa lên blockchain",
				"detail": err.Error(),
			})
			return

		case errors.Is(err, common.ErrCertificateNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error":  "Không tìm thấy văn bằng",
				"detail": err.Error(),
			})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Lỗi hệ thống khi đưa lên blockchain",
				"detail": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Ghi văn bằng lên blockchain thành công",
		"transaction_id": txID,
		"certificate_id": certID.Hex(),
	})
}

func (h *BlockchainHandler) GetCertificateByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.BlockchainSvc.GetCertificateFromChain(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
func (h *BlockchainHandler) VerifyCertificateIntegrity(c *gin.Context) {
	certID := c.Param("id")
	if certID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu ID văn bằng"})
		return
	}

	ok, msg, onChainCert, cert, user, faculty, university, err := h.BlockchainSvc.VerifyCertificateIntegrity(c.Request.Context(), certID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "certID không hợp lệ"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "không tìm thấy"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Lỗi khi xác minh dữ liệu",
				"message": err.Error(),
			})
		}
		return
	}

	resp := mapper.MapCertificateToResponse(cert, user, faculty, university)

	if !ok {
		c.JSON(http.StatusConflict, gin.H{
			"valid":       false,
			"message":     msg,
			"on_chain":    onChainCert,
			"certificate": resp,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"message":     msg,
		"on_chain":    onChainCert,
		"certificate": resp,
	})
}

func (h *BlockchainHandler) VerifyCertificateFile(c *gin.Context) {
	certIDHex := c.Param("id")
	certID, err := primitive.ObjectIDFromHex(certIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	stream, contentType, err := h.BlockchainSvc.VerifyFileByID(c.Request.Context(), certID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()

	// Trả file trực tiếp cho người dùng (xem được trong browser)
	c.DataFromReader(http.StatusOK, -1, contentType, stream, nil)
}

func (h *BlockchainHandler) PushEDiplomasToBlockchain(c *gin.Context) {
	var req struct {
		FacultyID       string `form:"faculty_id"`
		CertificateType string `form:"certificate_type"`
		Course          string `form:"course"`
		Issued          *bool  `form:"issued"`
	}

	// Bind form-data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Lấy claims từ context
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

	// Parse university_id
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	// Gọi service
	count, err := h.BlockchainSvc.PushToBlockchain(
		c.Request.Context(),
		universityID.Hex(), // truyền lại string vào service
		req.FacultyID,
		req.CertificateType,
		req.Course,
		req.Issued,
	)
	if err != nil {
		// Phân loại lỗi dựa theo message
		switch {
		case strings.Contains(err.Error(), "invalid faculty_id"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "no eDiplomas found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "no valid eDiplomas to push"):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Đã đẩy lên chuỗi khối",
		"updated_records": count,
	})
}

type VerifyBatchRequest struct {
	UniversityID    string `json:"university_id" binding:"required"`
	FacultyID       string `form:"faculty_id" json:"faculty_id" binding:"required"`
	CertificateType string `form:"certificate_type" json:"certificate_type"`
	Course          string `form:"course" json:"course"`
}

type VerifyBatchResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

func (h *BlockchainHandler) VerifyBatch(c *gin.Context) {
	var req VerifyBatchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse university_id từ request
	universityID, err := primitive.ObjectIDFromHex(req.UniversityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	// Gọi service để verify
	verified, msg, err := h.BlockchainSvc.VerifyBatch(
		c.Request.Context(),
		universityID.Hex(),
		req.FacultyID,
		req.CertificateType,
		req.Course,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Trả response chuẩn
	c.JSON(http.StatusOK, VerifyBatchResponse{
		Verified: verified,
		Message:  msg,
	})
}
