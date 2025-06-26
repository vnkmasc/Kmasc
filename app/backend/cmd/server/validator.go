package main

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
			re := regexp.MustCompile(`^\d{4}$`)
			return re.MatchString(fl.Field().String())
		})

		_ = v.RegisterValidation("certtype", func(fl validator.FieldLevel) bool {
			val := strings.TrimSpace(fl.Field().String())
			switch val {
			case "Cử nhân", "Kỹ sư", "Thạc sĩ", "Tiến sĩ":
				return true
			default:
				return false
			}
		})

		v.RegisterStructValidation(models.ValidateCreateCertificateRequest, models.CreateCertificateRequest{})

		println("✅ Đăng ký validator thành công")
	}
}
