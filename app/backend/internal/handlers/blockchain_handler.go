package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/mapper"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlockchainHandler struct {
	BlockchainSvc service.BlockchainService
	EDiplomaSvc   service.EDiplomaService
}

func NewBlockchainHandler(blockchainSvc service.BlockchainService, ediplomaSvc service.EDiplomaService) *BlockchainHandler {
	return &BlockchainHandler{
		BlockchainSvc: blockchainSvc,
		EDiplomaSvc:   ediplomaSvc,
	}
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
		case strings.Contains(err.Error(), "không tìm thấy trên chuỗi khối"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		FacultyID       string `form:"faculty_id" json:"faculty_id"`
		CertificateType string `form:"certificate_type" json:"certificate_type"`
		Course          string `form:"course" json:"course"`
		Issued          *bool  `form:"issued" json:"issued"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Lấy claims
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
		universityID.Hex(),
		req.FacultyID,
		req.CertificateType,
		req.Course,
		req.Issued,
	)

	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidFaculty):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, common.ErrNoDiplomas):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, common.ErrNoValidDiplomas):
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

func (h *BlockchainHandler) PushEDiplomasToBlockchain1(c *gin.Context) {
	var req struct {
		FacultyID       string `json:"faculty_id"`       // bind từ JSON
		CertificateType string `json:"certificate_type"` // bind từ JSON
		Course          string `json:"course"`           // bind từ JSON
	}

	// Bind raw JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Lấy claims từ context (JWT token)
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

	// Gọi service để đẩy eDiplomas lên blockchain
	count, err := h.BlockchainSvc.PushToBlockchain1(
		c.Request.Context(),
		universityID.Hex(), // UniversityID lấy từ token
		req.FacultyID,
		req.CertificateType,
		req.Course,
	)

	if err != nil {
		switch err {
		case common.ErrInvalidFaculty:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case common.ErrNoDiplomas:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case common.ErrNoValidDiplomas, common.ErrAlreadyOnChain:
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

func (h *BlockchainHandler) VerifyBatch(c *gin.Context) {
	var req struct {
		UniversityID    string `json:"university_id" binding:"required"`
		FacultyID       string `json:"faculty_id"`
		CertificateType string `json:"certificate_type"`
		Course          string `json:"course"`
		EDiplomaID      string `json:"ediploma_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	result, err := h.BlockchainSvc.VerifyBatch(
		c.Request.Context(),
		req.UniversityID,
		req.FacultyID,
		req.CertificateType,
		req.Course,
		req.EDiplomaID,
	)

	if err != nil {
		switch err {
		case common.ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case common.ErrInvalidFaculty:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Lỗi khi verify batch: %v", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"batch_id": result.BatchID,
		"verified": result.Verified,
		"details":  result.Details,
		"data":     result.EDiplomaData,
	})
}

func (h *BlockchainHandler) VerifyCertificateBatch(c *gin.Context) {
	var req struct {
		UniversityID    string `json:"university_id" binding:"required"`
		FacultyID       string `json:"faculty_id"`
		CertificateType string `json:"certificate_type"`
		Course          string `json:"course"`
		CertificateID   string `json:"certificate_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	result, err := h.BlockchainSvc.VerifyCertificateBatch(
		c.Request.Context(),
		req.UniversityID,
		req.FacultyID,
		req.CertificateType,
		req.Course,
		req.CertificateID,
	)

	if err != nil {
		switch err {
		case common.ErrBatchNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case common.ErrInvalidFaculty:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Lỗi khi verify certificate batch: %v", err),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"batch_id": result.BatchID,
		"verified": result.Verified,
		"details":  result.Details,
		"data":     result.CertificateData,
	})
}

type VerifyEDiplomaRequest struct {
	UniversityID string `json:"university_id" binding:"required"`
	FacultyID    string `json:"faculty_id" binding:"required"`
	Course       string `json:"course" binding:"required"`
	StudentCode  string `json:"student_code" binding:"required"`
}

type VerifyResult struct {
	StudentCode    string             `json:"student_code"`
	Valid          bool               `json:"valid"`
	BlockchainRoot string             `json:"blockchain_root"`
	ComputedHash   string             `json:"computed_hash"`
	Proof          []models.ProofNode `json:"proof"`
}

// Handler
func (h *BlockchainHandler) VerifyEDiploma(c *gin.Context) {
	var req VerifyEDiplomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	result, err := h.BlockchainSvc.VerifyEDiploma(
		c.Request.Context(),
		req.UniversityID,
		req.FacultyID,
		req.Course,
		req.StudentCode,
	)
	if err != nil {
		if errors.Is(err, common.ErrNoDiplomas) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BlockchainHandler) PushCertificatesToBlockchain(c *gin.Context) {
	var req struct {
		FacultyID       string `json:"faculty_id"`
		CertificateType string `json:"certificate_type"`
		Course          string `json:"course"`
	}

	// Bind input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Lấy claims từ context (JWT token)
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

	// Gọi service để đẩy Certificates lên blockchain
	count, err := h.BlockchainSvc.PushCertificatesToBlockchain(
		c.Request.Context(),
		universityID.Hex(),
		req.FacultyID,
		req.CertificateType,
		req.Course,
	)

	if err != nil {
		switch err {
		case common.ErrInvalidFaculty:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case common.ErrCertificateNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case common.ErrNoValidCertificates:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Certificates đã được đẩy lên blockchain",
		"updated_records": count,
	})
}
