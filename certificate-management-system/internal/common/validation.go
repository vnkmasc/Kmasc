package common

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func TranslateError(field, tag string) string {
	messages := map[string]map[string]string{
		"StudentID": {
			"required": "Mã sinh viên không được để trống",
		},
		"FullName": {
			"required": "Họ tên không được để trống",
		},
		"Email": {
			"required": "Email không được để trống",
			"email":    "Email không hợp lệ",
		},
		"Faculty": {
			"required": "Khoa không được để trống",
		},
		"Class": {
			"required": "Lớp không được để trống",
		},
		"Course": {
			"required":   "Khóa học không được để trống",
			"courseyear": "Khóa học phải có định dạng năm, ví dụ 2021",
		},
		"PersonalEmail": {
			"required": "Email cá nhân không được để trống",
			"email":    "Email cá nhân không hợp lệ",
		},
		"Password": {
			"required": "Mật khẩu không được để trống",
		},
		"CertificateType": {
			"required": "Loại văn bằng không được để trống",
			"certtype": "Loại văn bằng phải là Cử nhân, Kỹ Sư, Thạc sĩ hoặc Tiến sĩ",
		},
		"Name": {
			"required": "Tên văn bằng không được để trống",
		},
		"Issuer": {
			"required": "Nơi cấp không được để trống",
		},
		"IssueDate": {
			"required": "Ngày cấp không được để trống",
		},
		"SerialNumber": {
			"required": "Số hiệu văn bằng không được để trống",
		},
		"RegistrationNumber": {
			"required": "Số vào sổ không được để trống",
		},
		"UniversityName": {
			"required": "Tên trường không được để trống",
		},
		"UniversityCode": {
			"required": "Mã trường không được để trống",
		},
		"Address": {
			"required": "Địa chỉ trường không được để trống",
		},
		"EmailDomain": {
			"required": "Tên miền email không được để trống",
			"email":    "Tên miền email không hợp lệ",
		},
		"Action": {
			"required": "Hành động không được để trống",
			"oneof":    "Hành động thực hiện approve hoặc reject",
		},
		"NewPassword": {
			"required": "Mật khẩu mới không được để trống"},
		"OldPassword": {
			"required": "Yêu cầu nhập mật khẩu cũ",
		},
	}

	if fieldMsg, ok := messages[field]; ok {
		if msg, ok2 := fieldMsg[tag]; ok2 {
			return msg
		}
	}
	return field + " không hợp lệ"
}

func ParseValidationError(err error) (map[string]string, bool) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(map[string]string)
		for _, e := range ve {
			field := e.Field()
			tag := e.Tag()
			if field != "" && tag != "" {
				errs[field] = TranslateError(field, tag)
			}
		}
		return errs, true
	}
	return nil, false
}
