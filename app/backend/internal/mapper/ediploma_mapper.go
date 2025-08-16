package mapper

import "github.com/vnkmasc/Kmasc/app/backend/internal/models"

// Helpers để tránh nil pointer
func getUniversityCode(u *models.University) string {
	if u == nil {
		return ""
	}
	return u.UniversityCode
}

func getUniversityName(u *models.University) string {
	if u == nil {
		return ""
	}
	return u.UniversityName
}

func getFacultyCode(f *models.Faculty) string {
	if f == nil {
		return ""
	}
	return f.FacultyCode
}

func getFacultyName(f *models.Faculty) string {
	if f == nil {
		return ""
	}
	return f.FacultyName
}

func getTemplateName(t *models.DiplomaTemplate) string {
	if t == nil {
		return ""
	}
	return t.Name
}

func getUserFullName(u *models.User) string {
	if u == nil {
		return ""
	}
	return u.FullName
}

func MapEDiplomaToDTO(
	ed *models.EDiploma,
	university *models.University,
	faculty *models.Faculty,
	template *models.DiplomaTemplate,
	user *models.User,
) *models.EDiplomaResponse {
	return &models.EDiplomaResponse{
		ID:             ed.ID,
		Name:           ed.Name,
		FacultyID:      ed.FacultyID,
		UniversityCode: getUniversityCode(university),
		UniversityName: getUniversityName(university),
		FacultyCode:    getFacultyCode(faculty),
		FacultyName:    getFacultyName(faculty),
		StudentCode:    ed.StudentCode,
		FullName:       ed.FullName,
		StudentName:    getUserFullName(user),
		TemplateName:   getTemplateName(template),

		CertificateType:    ed.CertificateType,
		Course:             ed.Course,
		EducationType:      ed.EducationType,
		GPA:                ed.GPA,
		GraduationRank:     ed.GraduationRank,
		IssueDate:          ed.IssueDate,
		SerialNumber:       ed.SerialNumber,
		RegistrationNumber: ed.RegistrationNumber,
		Issued:             ed.Issued,
		Signed:             ed.Signed,
		DataEncrypted:      ed.DataEncrypted,
		OnBlockchain:       ed.OnBlockchain,
	}
}
