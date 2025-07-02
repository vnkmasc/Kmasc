package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. Khóa học 4 chữ số
		_ = v.RegisterValidation("courseyear", func(fl validator.FieldLevel) bool {
			return regexp.MustCompile(`^\d{4}$`).MatchString(fl.Field().String())
		})

		// 2. Định dạng ngày dd/mm/yyyy (chỉ dùng nếu ngày là string)
		_ = v.RegisterValidation("dateformat", func(fl validator.FieldLevel) bool {
			dateStr := fl.Field().String()
			_, err := time.Parse("02/01/2006", dateStr)
			return err == nil
		})

		// 3. CCCD: 12 chữ số
		_ = v.RegisterValidation("citizenid", func(fl validator.FieldLevel) bool {
			return regexp.MustCompile(`^\d{12}$`).MatchString(fl.Field().String())
		})

		// 4. Mức độ kỷ luật 1–4
		_ = v.RegisterValidation("disciplinelevel", func(fl validator.FieldLevel) bool {
			if fl.Field().IsNil() {
				return false
			}
			level := int(fl.Field().Int())
			return level >= 1 && level <= 4
		})

		// 5. Struct-level validator cho CreateCertificateRequest
		v.RegisterStructValidation(models.ValidateCreateCertificateRequest, models.CreateCertificateRequest{})
	}
}
