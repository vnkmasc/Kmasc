package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đưa lên blockchain", "detail": err.Error()})
		return
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

	ok, msg, onChainCert, cert, err := h.BlockchainSvc.VerifyCertificateIntegrity(c.Request.Context(), certID)
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

	if !ok {
		c.JSON(http.StatusConflict, gin.H{
			"valid":       false,
			"message":     msg,
			"on_chain":    onChainCert,
			"certificate": cert,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"message":     msg,
		"on_chain":    onChainCert,
		"certificate": cert,
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
