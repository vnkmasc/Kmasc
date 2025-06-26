package common

import "errors"

var (
	//auth
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid_token")

	ErrNoFieldsToUpdate               = errors.New("no_fields_to_update")
	ErrUserNotExisted                 = errors.New("user_not_exists")
	ErrInvalidUserID                  = errors.New("invalid_user_id")
	ErrStudentIDExists                = errors.New("student_id_exists")
	ErrEmailExists                    = errors.New("email_exists")
	ErrUniversityNameExists           = errors.New("university_name_exists")
	ErrUniversityEmailDomainExists    = errors.New("university_email_domain_exists")
	ErrUniversityCodeExists           = errors.New("university_code_exists")
	ErrUniversityNotFound             = errors.New("university not found")
	ErrAccountUniversityNotFound      = errors.New("university account not found")
	ErrUniversityAlreadyApproved      = errors.New("university_already_approved")
	ErrAccountUniversityAlreadyExists = errors.New("university_admin_account_already_exists")
	ErrAccountNotFound                = errors.New("account_not_found")
	ErrInvalidOldPassword             = errors.New("invalid_old_password")
	ErrPersonalAccountAlreadyExist    = errors.New("personal_account_already_exists")
	ErrCheckingPersonalAccount        = errors.New("error_checking_personal_account")

	//Faculty
	ErrFacultyNotFound   = errors.New("faculty_not_found")
	ErrFacultyCodeExists = errors.New("faculty_code_existed")

	//Certificate
	ErrCertificateNotFound      = errors.New("certificate_not_found")
	ErrCertificateAlreadyExists = errors.New("certificate_already_exists")
	ErrSerialNumberExists       = errors.New("serial_number_exists")
	ErrRegNoExists              = errors.New("reg_no_exists")

	ErrMissingRequiredFieldsForDegree      = errors.New("missing_required_fields_for_degree")
	ErrMissingRequiredFieldsForCertificate = errors.New("missing_required_fields_for_certificate")
)
