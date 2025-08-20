package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateHandler struct {
	certificateService service.CertificateService
	universityService  service.UniversityService
	facultyService     service.FacultyService
	userService        service.UserService
	accountService     service.AuthService
	minioClient        *database.MinioClient
}

func NewCertificateHandler(
	certSvc service.CertificateService,
	uniSvc service.UniversityService,
	facultySvc service.FacultyService,
	userSvc service.UserService,
	authSvc service.AuthService,
	minioClient *database.MinioClient,
) *CertificateHandler {
	return &CertificateHandler{
		certificateService: certSvc,
		universityService:  uniSvc,
		facultyService:     facultySvc,
		userService:        userSvc,
		accountService:     authSvc,
		minioClient:        minioClient,
	}
}

func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	var req models.CreateCertificateRequest

	// Đọc raw body để debug input
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	fmt.Println(">>> Raw request body:", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reset lại body cho ShouldBindJSON

	// Validate JSON đầu vào
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(">>> Binding error:", err) // In lỗi gốc

		if validationErrs, ok := common.ParseValidationError(err); ok {
			fmt.Println(">>> Validation errors:", validationErrs) // In lỗi từng field
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Dữ liệu không hợp lệ",
				"details": validationErrs,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Dữ liệu không hợp lệ",
		})
		return
	}

	fmt.Printf(">>> Parsed struct: %+v\n", req) // In request đã parse xong

	// Lấy claims từ context
	claims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Không xác thực được người dùng",
		})
		return
	}

	// Gọi service
	if err := h.certificateService.CreateCertificate(c.Request.Context(), claims, &req); err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token không hợp lệ"})

		case errors.Is(err, common.ErrUserNotExisted):
			c.JSON(http.StatusNotFound, gin.H{"message": "Sinh viên không tồn tại"})

		case errors.Is(err, common.ErrCertificateAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"message": "Văn bằng đã tồn tại"})

		case errors.Is(err, common.ErrFacultyNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy khoa"})

		case errors.Is(err, common.ErrUniversityNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Không tìm thấy trường"})

		case errors.Is(err, common.ErrSerialNumberExists):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Số hiệu văn bằng đã tồn tại"})

		case errors.Is(err, common.ErrRegNoExists):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Số vào sổ gốc đã tồn tại"})

		case errors.Is(err, common.ErrMissingRequiredFieldsForDegree):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Thiếu thông tin bắt buộc cho văn bằng"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
		}
		return
	}

	// Trả về thành công
	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo chứng nhận thành công",
	})
}

func (h *CertificateHandler) GetAllCertificates(c *gin.Context) {
	certs, err := h.certificateService.GetAllCertificates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Lỗi hệ thống",
			"chi_tiet": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": certs,
	})

}

func (h *CertificateHandler) GetCertificateByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID certificate không hợp lệ"})
		return
	}

	cert, err := h.certificateService.GetCertificateByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrCertificateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy certificate"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống", "chi_tiet": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cert,
	})

}
func (h *CertificateHandler) UploadCertificateFile(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*utils.CustomClaims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng chọn file để tải lên"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chỉ hỗ trợ file PDF, JPG, JPEG, PNG"})
		return
	}

	isDegree := c.Query("is_degree") == "true"
	certificateName := c.Query("name")
	studentCode := strings.TrimSuffix(file.Filename, ext)

	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ (UniversityID sai định dạng)"})
		return
	}

	// Tìm chứng chỉ
	var certificate *models.Certificate
	if isDegree {
		certificate, err = h.certificateService.GetDegreeCertificateByStudentCodeAndUniversity(
			c.Request.Context(), studentCode, universityID)
	} else {
		if certificateName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tên chứng chỉ (query param 'name')"})
			return
		}
		certificate, err = h.certificateService.GetCertificateByStudentCodeAndNameAndUniversity(
			c.Request.Context(), studentCode, certificateName, universityID)
	}
	if err != nil || certificate == nil || certificate.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng/chứng chỉ phù hợp"})
		return
	}
	if certificate.Path != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Văn bằng/chứng chỉ đã có file, không thể ghi đè"})
		return
	}

	// Đọc file gốc
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc file"})
		return
	}

	// Upload file trực tiếp (bỏ mã hóa)
	objectPath, err := h.certificateService.UploadCertificateFileDirect(
		c.Request.Context(), certificate.ID, fileData, file.Filename, isDegree,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải lên thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tải file thành công",
		"path":    objectPath,
	})
}

func (h *CertificateHandler) GetCertificateFile(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	certificateID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	certificate, err := h.certificateService.GetRawCertificateByID(ctx, certificateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng"})
		return
	}

	if certificate.Path == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng chưa có file"})
		return
	}

	// Lấy object từ MinIO (không cần giải mã)
	object, err := h.minioClient.Client.GetObject(ctx, h.minioClient.Bucket, certificate.Path, minio.GetObjectOptions{})
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

	// Trả về file gốc
	contentType := http.DetectContentType(fileData)
	c.DataFromReader(http.StatusOK, int64(len(fileData)), contentType, bytes.NewReader(fileData), nil)
}

func (h *CertificateHandler) GetCertificatesByStudentID(c *gin.Context) {
	ctx := c.Request.Context()
	studentIDParam := c.Param("id")

	studentID, err := primitive.ObjectIDFromHex(studentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sinh viên không hợp lệ"})
		return
	}

	user, err := h.userService.GetUserByID(ctx, studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sinh viên"})
		return
	}

	faculty, err := h.facultyService.GetFacultyByCode(ctx, user.FacultyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy khoa"})
		return
	}

	university, err := h.universityService.GetUniversityByCode(ctx, user.UniversityCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy trường đại học"})
		return
	}

	certificate, err := h.certificateService.GetCertificateByUserID(ctx, studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng của người dùng"})
		return
	}

	result := models.CertificateResponse{
		ID:              certificate.ID,
		UserID:          certificate.UserID,
		StudentCode:     certificate.StudentCode,
		CertificateType: certificate.CertificateType,
		Name:            certificate.Name,
		SerialNumber:    certificate.SerialNumber,
		RegNo:           certificate.RegNo,
		Path:            certificate.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		Signed:          certificate.Signed,
		CreatedAt:       certificate.CreatedAt,
		UpdatedAt:       certificate.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *CertificateHandler) SearchCertificates(c *gin.Context) {
	var params models.SearchCertificateParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	certs, total, err := h.certificateService.SearchCertificates(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       certs,
		"total":      total,
		"page":       params.Page,
		"page_size":  params.PageSize,
		"total_page": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
	})
}

func (h *CertificateHandler) GetMyCertificates(c *gin.Context) {
	val, exists := c.Get(string(utils.ClaimsContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token không hợp lệ"})
		return
	}

	certificates, err := h.certificateService.GetCertificatesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certificates})
}

func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	err = h.certificateService.DeleteCertificateByID(c.Request.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng không tồn tại"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa văn bằng"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa văn bằng thành công"})
}
func (h *CertificateHandler) GetMyCertificateNames(c *gin.Context) {
	val, exists := c.Get(string(utils.ClaimsContextKey))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bạn chưa đăng nhập hoặc token không hợp lệ"})
		return
	}
	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token không hợp lệ"})
		return
	}

	certificates, err := h.certificateService.GetSimpleCertificatesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certificates})
}

func (h *CertificateHandler) ImportCertificatesFromExcel(c *gin.Context) {
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

	// Đọc file
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
	if err != nil || len(rows) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy dữ liệu trong Sheet1"})
		return
	}

	var successResults []map[string]interface{}
	var errorResults []map[string]interface{}

	for i, row := range rows {
		if i == 0 {
			continue // Bỏ qua header
		}

		result := map[string]interface{}{"row": i + 1}
		get := func(index int) string {
			if index < len(row) {
				return strings.TrimSpace(row[index])
			}
			return ""
		}

		isDegree := strings.ToLower(get(1)) == "văn bằng"

		gpa := 0.0
		if gpaStr := get(10); gpaStr != "" {
			gpa, _ = strconv.ParseFloat(gpaStr, 64)
		}
		parsedDate, err := parseDateTime(get(11)) // Cột Ngày cấp
		if err != nil {
			result["error"] = fmt.Sprintf("Ngày cấp không hợp lệ: %v", err)
			errorResults = append(errorResults, result)
			continue
		}

		req := &models.CreateCertificateRequest{
			StudentCode:     get(0),
			IsDegree:        isDegree,
			Name:            get(2),
			CertificateType: get(3),
			Course:          get(4),
			GraduationRank:  get(5),
			EducationType:   get(6),
			SerialNumber:    get(7),
			RegNo:           get(8),
			Major:           get(9),
			GPA:             gpa,
			IssueDate:       parsedDate,
			Description:     get(12),
		}

		// 👉 Debug log trước khi gọi service
		fmt.Printf(">>> [ROW %d] Creating certificate for StudentCode: '%s', IsDegree: %v\n", i+1, req.StudentCode, req.IsDegree)
		fmt.Printf(">>>        Name: %s | Serial: %s | RegNo: %s | Date: %s\n", req.Name, req.SerialNumber, req.RegNo, req.IssueDate.Format("2006-01-02"))
		fmt.Printf(">>>        UniversityID: %s\n", claims.UniversityID)

		err = h.certificateService.CreateCertificate(c.Request.Context(), claims, req)
		if err != nil {
			result["error"] = mapErrorToMessage(err)
			errorResults = append(errorResults, result)
		} else {
			result["status"] = "Tạo thành công"
			successResults = append(successResults, result)
		}
	}

	// Trả về kết quả
	if len(errorResults) == 0 {
		c.JSON(http.StatusCreated, gin.H{
			"message":       "Tất cả chứng nhận đã được tạo thành công",
			"success_count": len(successResults),
			"error_count":   0,
			"data":          gin.H{"success": successResults},
		})
	} else {
		c.JSON(http.StatusMultiStatus, gin.H{
			"message":       "Một số chứng nhận không thể tạo",
			"success_count": len(successResults),
			"error_count":   len(errorResults),
			"data": gin.H{
				"success": successResults,
				"error":   errorResults,
			},
		})
	}
}

func mapErrorToMessage(err error) string {
	switch {
	case errors.Is(err, common.ErrInvalidToken):
		return "Token không hợp lệ"
	case errors.Is(err, common.ErrUserNotExisted):
		return "Không tìm thấy sinh viên"
	case errors.Is(err, common.ErrSerialNumberExists):
		return "Số hiệu đã tồn tại"
	case errors.Is(err, common.ErrRegNoExists):
		return "Số vào sổ đã tồn tại"
	case errors.Is(err, common.ErrMissingRequiredFieldsForDegree):
		return "Thiếu thông tin bắt buộc cho văn bằng"
	case errors.Is(err, common.ErrCertificateAlreadyExists):
		return "Văn bằng/chứng chỉ này đã tồn tại"
	default:
		fmt.Printf(">>> mapErrorToMessage - Unknown error: %+v\n", err)
		return "Lỗi hệ thống hoặc không xác định"
	}
}

func parseDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("02/01/2006", s)
}
