package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CertificateHandler struct {
	certificateService service.CertificateService
	universityService  service.UniversityService
	facultyService     service.FacultyService
	userService        service.UserService
	minioClient        *database.MinioClient
}

func NewCertificateHandler(
	certSvc service.CertificateService,
	uniSvc service.UniversityService,
	facultySvc service.FacultyService,
	userSvc service.UserService,
	minioClient *database.MinioClient,
) *CertificateHandler {
	return &CertificateHandler{
		certificateService: certSvc,
		universityService:  uniSvc,
		facultyService:     facultySvc,
		userService:        userSvc,
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

	// Lấy file upload
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

	// Parse query param
	isDegree := c.Query("is_degree") == "true"
	certificateType := c.Query("certificate_type")
	certificateName := c.Query("name")

	fmt.Println("=== DEBUG UPLOAD START ===")
	fmt.Println("Filename:", file.Filename)
	fmt.Println("Extension:", ext)
	fmt.Println("isDegree:", isDegree)
	fmt.Println("certificate_type:", certificateType)
	fmt.Println("certificate_name:", certificateName)

	// Lấy university
	universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
	if err != nil {
		fmt.Println("Lỗi universityID:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token không hợp lệ (UniversityID không đúng định dạng)"})
		return
	}
	fmt.Println("UniversityID:", universityID.Hex())

	university, err := h.universityService.GetUniversityByID(c.Request.Context(), universityID)
	if err != nil {
		fmt.Println("Lỗi lấy university:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không lấy được thông tin trường đại học"})
		return
	}

	// Giả định tên file là mã sinh viên
	studentCode := strings.TrimSuffix(file.Filename, ext)
	fmt.Println("StudentCode từ filename:", studentCode)

	// Truy vấn certificate
	var certificate *models.Certificate
	if isDegree {
		if certificateType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu loại văn bằng (query param 'certificate_type')"})
			return
		}
		certificate, err = h.certificateService.GetCertificateByStudentCodeAndTypeAndUniversity(
			c.Request.Context(), studentCode, certificateType, university.ID)
	} else {
		if certificateName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tên chứng chỉ (query param 'name')"})
			return
		}
		certificate, err = h.certificateService.GetCertificateByStudentCodeAndNameAndUniversity(
			c.Request.Context(), studentCode, certificateName, university.ID)
	}

	if err != nil {
		fmt.Println("Lỗi truy vấn certificate:", err)
	}
	if certificate == nil || certificate.ID.IsZero() {
		fmt.Println("Certificate không tìm thấy hoặc ID rỗng")
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng/chứng chỉ phù hợp"})
		return
	}
	fmt.Println("Certificate tìm thấy:", certificate.ID.Hex())

	// Check permission
	if certificate.UniversityID != university.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không được phép cập nhật văn bằng này"})
		return
	}

	if certificate.Path != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Văn bằng/chứng chỉ đã có file, không thể ghi đè"})
		return
	}

	// Đọc nội dung file
	src, err := file.Open()
	if err != nil {
		fmt.Println("Lỗi mở file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		fmt.Println("Lỗi đọc file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc file"})
		return
	}

	// Tạo tên file lưu lên MinIO: CT060344/thac-si.pdf hoặc CT060344/toeic-800+.pdf
	var typeStr string
	if certificate.IsDegree {
		typeStr = certificate.CertificateType
	} else {
		typeStr = certificate.Name
	}

	slug := utils.Slugify(typeStr)
	finalFileName := fmt.Sprintf("%s/%s%s", certificate.StudentCode, slug, ext)
	fmt.Println("File path lưu lên MinIO:", finalFileName)

	// Upload file và lưu đường dẫn
	filePath, err := h.certificateService.UploadCertificateFile(
		c.Request.Context(), certificate.ID, fileData, finalFileName, isDegree, typeStr)
	if err != nil {
		fmt.Println("Lỗi upload file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tải lên thất bại: " + err.Error()})
		return
	}

	fmt.Println("Upload thành công. Đường dẫn:", filePath)
	fmt.Println("=== DEBUG UPLOAD END ===")

	c.JSON(http.StatusOK, gin.H{
		"message": "Tải file thành công",
		"path":    filePath,
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

	certificate, err := h.certificateService.GetCertificateByID(ctx, certificateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy văn bằng"})
		return
	}

	if certificate.Path == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Văn bằng chưa có file"})
		return
	}

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
