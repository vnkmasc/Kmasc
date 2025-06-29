package mapper

import (
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCertificateToResponse(cert *models.Certificate, user *models.User, faculty *models.Faculty, university *models.University) *models.CertificateResponse {
	return &models.CertificateResponse{
		ID:              cert.ID.Hex(),
		UserID:          cert.UserID.Hex(),
		StudentCode:     cert.StudentCode,
		StudentName:     user.FullName,
		CertificateType: cert.CertificateType,
		Name:            cert.Name,
		SerialNumber:    cert.SerialNumber,
		RegNo:           cert.RegNo,
		Path:            cert.Path,
		FacultyCode:     faculty.FacultyCode,
		FacultyName:     faculty.FacultyName,
		UniversityCode:  university.UniversityCode,
		UniversityName:  university.UniversityName,
		CertHash:        cert.CertHash,
		HashFile:        cert.HashFile,
		Major:           cert.Major,
		Course:          cert.Course,
		GPA:             cert.GPA,
		GraduationRank:  cert.GraduationRank,
		EducationType:   cert.EducationType,
		Signed:          cert.Signed,
		IssueDate:       cert.IssueDate.Format("02/01/2006"),
		CreatedAt:       cert.CreatedAt,
		UpdatedAt:       cert.UpdatedAt,
		Description:     cert.Description,
	}
}
func MapCertificatesToResponses(certs []*models.Certificate, userMap map[primitive.ObjectID]*models.User, facultyMap map[primitive.ObjectID]*models.Faculty, universityMap map[primitive.ObjectID]*models.University) []*models.CertificateResponse {
	var responses []*models.CertificateResponse
	for _, cert := range certs {
		user := userMap[cert.UserID]
		faculty := facultyMap[cert.FacultyID]
		university := universityMap[cert.UniversityID]
		responses = append(responses, MapCertificateToResponse(cert, user, faculty, university))
	}
	return responses
}
