package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardDisciplineHandler struct {
	rdService service.RewardDisciplineService
}

func NewRewardDisciplineHandler(rdService service.RewardDisciplineService) *RewardDisciplineHandler {
	return &RewardDisciplineHandler{
		rdService: rdService,
	}
}

func (h *RewardDisciplineHandler) CreateRewardDiscipline(c *gin.Context) {
	var req models.CreateRewardDisciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.rdService.CreateRewardDiscipline(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case common.ErrUserNotExisted:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy sinh viên với mã sinh viên này"})
		case common.ErrDecisionNumberExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Số quyết định đã tồn tại"})
		default:
			if ve, ok := err.(*common.ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
			}
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *RewardDisciplineHandler) GetRewardDisciplineByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	resp, err := h.rdService.GetRewardDisciplineByID(c.Request.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == common.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bản ghi"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *RewardDisciplineHandler) GetAllRewardDisciplines(c *gin.Context) {
	resp, err := h.rdService.GetAllRewardDisciplines(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *RewardDisciplineHandler) UpdateRewardDiscipline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	var req models.UpdateRewardDisciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.rdService.UpdateRewardDiscipline(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case common.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bản ghi"})
		case common.ErrUserNotExisted:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy sinh viên với mã số này"})
		case common.ErrDecisionNumberExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Số quyết định đã tồn tại"})
		default:
			if ve, ok := err.(*common.ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
			}
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công"})
}

func (h *RewardDisciplineHandler) DeleteRewardDiscipline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	err = h.rdService.DeleteRewardDiscipline(c.Request.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bản ghi"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}

func (h *RewardDisciplineHandler) SearchRewardDisciplines(c *gin.Context) {
	var params models.SearchRewardDisciplineParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tham số không hợp lệ"})
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	resp, total, err := h.rdService.SearchRewardDisciplines(c.Request.Context(), params)
	if err != nil {
		switch err {
		case common.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Không có quyền truy cập"})
		case common.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       resp,
		"total":      total,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_page": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
	})
}

func (h *RewardDisciplineHandler) GetMyRewardDisciplines(c *gin.Context) {
	resp, err := h.rdService.GetMyRewardDisciplines(c.Request.Context())
	if err != nil {
		switch err {
		case common.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Không có quyền truy cập"})
		case common.ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		case common.ErrUserNotExisted:
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy thông tin người dùng"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *RewardDisciplineHandler) ImportRewardDisciplinesFromExcel(c *gin.Context) {
	val, exists := c.Get(string(utils.ClaimsContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	_, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	// Get is_discipline parameter
	isDisciplineStr := c.Query("is_discipline")
	isDiscipline := false
	if isDisciplineStr == "true" {
		isDiscipline = true
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng upload file Excel"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File không đúng định dạng Excel"})
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil || len(rows) == 0 {
		rows, err = f.GetRows("Sheet")
		if err != nil || len(rows) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Không đọc được sheet dữ liệu (Sheet1 hoặc Sheet)"})
			return
		}
	}

	var (
		successResults []map[string]interface{}
		errorResults   []map[string]interface{}
	)

	for i, row := range rows {
		if i == 0 {
			continue // Skip header row
		}

		result := map[string]interface{}{"row": i + 1}

		// Check minimum required columns: name, DecisionNumber, Description, StudentCode
		minCols := 4
		if isDiscipline {
			minCols = 5 // Add DisciplineLevel column for discipline
		}

		if len(row) < minCols {
			result["error"] = "Thiếu dữ liệu"
			errorResults = append(errorResults, result)
			continue
		}

		// Create request
		req := &models.CreateRewardDisciplineRequest{
			Name:           row[0],
			DecisionNumber: row[1],
			Description:    row[2],
			StudentCode:    row[3],
			IsDiscipline:   isDiscipline,
		}

		// If this is a discipline, parse DisciplineLevel
		if isDiscipline && len(row) > 4 && row[4] != "" {
			level, err := strconv.Atoi(row[4])
			if err != nil {
				result["error"] = "Mức độ kỷ luật không hợp lệ"
				errorResults = append(errorResults, result)
				continue
			}
			if level < 1 || level > 4 {
				result["error"] = "Mức độ kỷ luật phải từ 1 đến 4"
				errorResults = append(errorResults, result)
				continue
			}
			req.DisciplineLevel = &level
		}

		// Create reward/discipline
		_, err := h.rdService.CreateRewardDiscipline(c.Request.Context(), req)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrUserNotExisted):
				result["error"] = "Không tìm thấy sinh viên với mã sinh viên này"
			case errors.Is(err, common.ErrDecisionNumberExists):
				result["error"] = "Số quyết định đã tồn tại"
			case errors.Is(err, common.ErrUnauthorized):
				result["error"] = "Bạn chưa đăng nhập hoặc token không hợp lệ"
			case errors.Is(err, common.ErrInvalidToken):
				result["error"] = "Token không hợp lệ"
			default:
				if ve, ok := err.(*common.ValidationError); ok {
					result["error"] = ve.Message
				} else {
					result["error"] = err.Error()
				}
			}
			errorResults = append(errorResults, result)
		} else {
			result["status"] = "Thêm thành công"
			successResults = append(successResults, result)
		}
	}

	if len(errorResults) == 0 {
		c.JSON(http.StatusCreated, gin.H{
			"success_count": len(successResults),
			"data": gin.H{
				"success": successResults,
				"error":   []map[string]interface{}{},
			},
		})
	} else {
		c.JSON(http.StatusMultiStatus, gin.H{
			"success_count": len(successResults),
			"error_count":   len(errorResults),
			"data": gin.H{
				"success": successResults,
				"error":   errorResults,
			},
		})
	}
}
