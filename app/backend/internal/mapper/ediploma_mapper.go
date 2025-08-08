package mapper

import "github.com/vnkmasc/Kmasc/app/backend/internal/models"

func MapEDiplomaToDTO(
	ed *models.EDiploma,
	university *models.University,
	faculty *models.Faculty,
	major *models.Major,
) *models.EDiplomaDTO {
	return &models.EDiplomaDTO{
		ID:             ed.ID,
		TemplateID:     ed.TemplateID,
		UniversityID:   ed.UniversityID,
		UniversityCode: university.UniversityCode,
		UniversityName: university.UniversityName,
		FacultyID:      ed.FacultyID,
		FacultyCode:    faculty.FacultyCode,
		FacultyName:    faculty.FacultyName,
		MajorID:        ed.MajorID,
		// MajorCode:      ifMajorNotNil(major, major.MajorCode),
		// MajorName:      ifMajorNotNil(major, major.MajorName),

		UserID:             ed.UserID,
		StudentCode:        ed.StudentCode,
		FullName:           ed.FullName,
		CertificateType:    ed.CertificateType,
		Course:             ed.Course,
		EducationType:      ed.EducationType,
		GPA:                ed.GPA,
		GraduationRank:     ed.GraduationRank,
		IssueDate:          ed.IssueDate,
		SerialNumber:       ed.SerialNumber,
		RegistrationNumber: ed.RegistrationNumber,
		FileLink:           ed.FileLink,
		FileHash:           ed.FileHash,
		Signature:          ed.Signature,
		Signed:             ed.Signed,
		SignedAt:           ed.SignedAt,
		OnBlockchain:       ed.OnBlockchain,
		CreatedAt:          ed.CreatedAt,
		UpdatedAt:          ed.UpdatedAt,

		// New fields from Template
		SignatureOfUni:    ed.SignatureOfUni,
		SignatureOfMinEdu: ed.SignatureOfMinEdu,
		Status:            ed.Status,
		IsLocked:          ed.IsLocked,
	}
}

func ifMajorNotNil(m *models.Major, val string) string {
	if m == nil {
		return ""
	}
	return val
}
