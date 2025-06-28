package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/internal/service"
	"github.com/tuyenngduc/certificate-management-system/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService       service.AuthService
	universityService service.UniversityService
	userService       service.UserService
	facultyService    service.FacultyService
}

func NewAuthHandler(
	authService service.AuthService,
	universityService service.UniversityService,
	userService service.UserService,
	facultyService service.FacultyService,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		universityService: universityService,
		userService:       userService,
		facultyService:    facultyService,
	}
}

func (h *AuthHandler) GetAllAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	accounts, total, err := h.authService.GetAllAccounts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var resp []models.AccountResponse
	for _, acc := range accounts {
		var studentID *primitive.ObjectID
		if acc.StudentID != primitive.NilObjectID {
			studentID = &acc.StudentID
		}
		var universityID *primitive.ObjectID
		if acc.UniversityID != primitive.NilObjectID {
			universityID = &acc.UniversityID
		}

		resp = append(resp, models.AccountResponse{
			ID:            acc.ID,
			StudentID:     studentID,
			UniversityID:  universityID,
			StudentEmail:  acc.StudentEmail,
			PersonalEmail: acc.PersonalEmail,
			CreatedAt:     acc.CreatedAt.Format(time.RFC3339),
			Role:          acc.Role,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(200, gin.H{
		"data":       resp,
		"page":       page,
		"page_size":  pageSize,
		"total":      total,
		"total_page": totalPages,
	})
}

func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var input models.RequestOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	err := h.authService.RequestOTP(c.Request.Context(), input)
	if err != nil {
		switch err {
		case common.ErrUserNotExisted:
			c.JSON(http.StatusNotFound, gin.H{"error": "Email không tồn tại trong hệ thống"})
		case common.ErrPersonalAccountAlreadyExist:
			c.JSON(http.StatusConflict, gin.H{"error": "Email này đã được liên kết với tài khoản cá nhân"})
		case common.ErrCheckingPersonalAccount:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi kiểm tra tài khoản cá nhân"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã gửi mã OTP tới email sinh viên"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	res, err := h.authService.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": res,
	})

}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo tài khoản thành công"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	account, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(account.ID, account.StudentID, account.UniversityID, account.Role, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tạo được token"})
		return
	}

	resp := models.LoginResponse{
		Token: token,
		Role:  account.Role,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email là bắt buộc"})
		return
	}

	err := h.authService.DeleteAccountByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, common.ErrAccountUniversityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trường không tồn tại"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xóa tài khoản thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa tài khoản thành công"})
}
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := common.ParseValidationError(err); ok {
			c.JSON(400, gin.H{"errors": errs})
			return
		}
		c.JSON(400, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.JSON(401, gin.H{"error": "Không xác thực"})
		return
	}

	claims, ok := claimsRaw.(*utils.CustomClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Dữ liệu xác thực không hợp lệ"})
		return
	}

	accountID, err := primitive.ObjectIDFromHex(claims.AccountID)
	if err != nil {
		c.JSON(401, gin.H{"error": "ID tài khoản không hợp lệ"})
		return
	}

	err = h.authService.ChangePassword(c.Request.Context(), accountID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case common.ErrAccountNotFound:
			c.JSON(404, gin.H{"error": "Không tìm thấy tài khoản"})
		case common.ErrInvalidOldPassword:
			c.JSON(400, gin.H{"error": "Mật khẩu cũ không đúng"})
		default:
			c.JSON(500, gin.H{"error": "Lỗi hệ thống"})
		}
		return
	}

	c.JSON(200, gin.H{"message": "Đổi mật khẩu thành công"})
}

func (h *AuthHandler) GetUniversityAdmins(c *gin.Context) {
	ctx := c.Request.Context()

	adminAccounts, err := h.authService.GetAccountsByRole(ctx, "university_admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tài khoản"})
		return
	}

	type AdminWithUniversityInfo struct {
		ID primitive.ObjectID `json:"id"`

		PersonalEmail  string `json:"personal_email"`
		UniversityCode string `json:"university_code"`
		UniversityName string `json:"university_name"`
		Address        string `json:"address"`
		Status         string `json:"status"`
	}

	var result []AdminWithUniversityInfo

	for _, acc := range adminAccounts {
		univ, err := h.universityService.GetUniversityByID(ctx, acc.UniversityID)
		if err != nil {
			continue
		}

		item := AdminWithUniversityInfo{
			ID:             acc.ID,
			PersonalEmail:  acc.PersonalEmail,
			UniversityCode: univ.UniversityCode,
			UniversityName: univ.UniversityName,
			Address:        univ.Address,
			Status:         univ.Status,
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

type AccountWithDetailsResponse struct {
	ID             string `json:"id"`
	PersonalEmail  string `json:"personal_email"`
	StudentCode    string `json:"student_code"`
	StudentName    string `json:"student_name"`
	UniversityName string `json:"university_name"`
	FacultyName    string `json:"faculty_name"`
}

func (h *AuthHandler) GetStudentAccounts(c *gin.Context) {
	ctx := c.Request.Context()

	accounts, err := h.authService.GetAccountsByRole(ctx, "student")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("accounts count:", len(accounts))

	var result []AccountWithDetailsResponse

	for _, acc := range accounts {
		user, err := h.userService.GetUserByID(ctx, acc.StudentID)
		if err != nil {
			continue
		}

		univ, err := h.universityService.GetUniversityByCode(ctx, user.UniversityCode)
		if err != nil {
			continue
		}

		faculty, err := h.facultyService.GetFacultyByCode(ctx, user.FacultyCode)
		if err != nil {
			continue
		}

		item := AccountWithDetailsResponse{
			ID:             acc.ID.Hex(),
			PersonalEmail:  acc.PersonalEmail,
			StudentCode:    user.StudentCode,
			StudentName:    user.FullName,
			UniversityName: univ.UniversityName,
			FacultyName:    faculty.FacultyName,
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
