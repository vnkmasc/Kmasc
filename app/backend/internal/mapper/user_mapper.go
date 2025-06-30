package mapper

import (
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
)

func MapUserToResponse(u *models.User, faculty *models.Faculty, university *models.University, loc *time.Location) models.UserResponse {
	if u == nil {
		return models.UserResponse{}
	}
	resp := models.UserResponse{
		ID:              u.ID,
		StudentCode:     u.StudentCode,
		FullName:        u.FullName,
		Email:           u.Email,
		Course:          u.Course,
		Status:          u.Status,
		FacultyCode:     "",
		FacultyName:     "",
		UniversityCode:  "",
		UniversityName:  "",
		CitizenIdNumber: u.CitizenIdNumber,
		Gender:          u.Gender,
		DateOfBirth:     u.DateOfBirth,
		Ethnicity:       u.Ethnicity,
		CurrentAddress:  u.CurrentAddress,
		BirthAddress:    u.BirthAddress,
		UnionJoinDate:   u.UnionJoinDate,
		PartyJoinDate:   u.PartyJoinDate,
		Description:     u.Description,
		CreatedAt:       u.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
		UpdatedAt:       u.UpdatedAt.In(loc).Format("2006-01-02 15:04:05"),
	}

	if faculty != nil {
		resp.FacultyCode = faculty.FacultyCode
		resp.FacultyName = faculty.FacultyName
	}
	if university != nil {
		resp.UniversityCode = university.UniversityCode
		resp.UniversityName = university.UniversityName
	}

	return resp
}
