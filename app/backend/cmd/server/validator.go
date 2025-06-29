package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
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

		// Add validator for date format dd/mm/yyyy
		_ = v.RegisterValidation("dateformat", func(fl validator.FieldLevel) bool {
			dateStr := fl.Field().String()
			if dateStr == "" {
				return false
			}
			_, err := time.Parse("02/01/2006", dateStr)
			return err == nil
		})

		// Add validator for citizen ID number
		_ = v.RegisterValidation("citizenid", func(fl validator.FieldLevel) bool {
			idStr := fl.Field().String()
			if idStr == "" {
				return false
			}
			// Check if the ID 12 digits (old and new format)
			re := regexp.MustCompile(`^\d{12}$`)
			return re.MatchString(idStr)
		})

		// Add validator for discipline level
		_ = v.RegisterValidation("disciplinelevel", func(fl validator.FieldLevel) bool {
			if fl.Field().IsNil() {
				return false
			}
			level := int(fl.Field().Int())
			return level >= 1 && level <= 4
		})

		v.RegisterStructValidation(models.ValidateCreateCertificateRequest, models.CreateCertificateRequest{})
	}
}
