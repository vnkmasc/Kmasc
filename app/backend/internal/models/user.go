package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentCode     string             `bson:"student_code" json:"student_code"`
	FullName        string             `bson:"full_name" json:"full_name"`
	Email           string             `bson:"email" json:"email"`
	FacultyID       primitive.ObjectID `bson:"faculty_id" json:"faculty"`
	UniversityID    primitive.ObjectID `bson:"university_id" json:"university_id"`
	Course          string             `bson:"course" json:"course"`
	Status          int                `bson:"status,omitempty"`
	CitizenIdNumber string             `bson:"citizen_id_number" json:"citizen_id_number"`
	Gender          bool               `bson:"gender" json:"gender"`
	DateOfBirth     string             `bson:"date_of_birth" json:"date_of_birth"`
	Ethnicity       string             `bson:"ethnicity" json:"ethnicity"`
	CurrentAddress  string             `bson:"current_address" json:"current_address"`
	BirthAddress    string             `bson:"birth_address" json:"birth_address"`
	UnionJoinDate   string             `bson:"union_join_date" json:"union_join_date"`
	PartyJoinDate   string             `bson:"party_join_date" json:"party_join_date"`
	Description     string             `bson:"description" json:"description"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateUserRequest struct {
	StudentCode     string `json:"student_code" binding:"required"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	FacultyCode     string `json:"faculty_code" binding:"required"`
	Course          string `json:"course" binding:"required,courseyear"`
	CitizenIdNumber string `json:"citizen_id_number" binding:"required,citizenid"`
	Gender          bool   `json:"gender"`
	DateOfBirth     string `json:"date_of_birth" binding:"required,dateformat"`
	Ethnicity       string `json:"ethnicity"`
	CurrentAddress  string `json:"current_address"`
	BirthAddress    string `json:"birth_address"`
	UnionJoinDate   string `json:"union_join_date" binding:"omitempty,dateformat"`
	PartyJoinDate   string `json:"party_join_date" binding:"omitempty,dateformat"`
	Description     string `json:"description"`
}

type UserResponse struct {
	ID              primitive.ObjectID `json:"id"`
	StudentCode     string             `json:"student_code"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	FacultyCode     string             `json:"faculty_code"`
	FacultyName     string             `json:"faculty_name"`
	UniversityCode  string             `json:"university_code"`
	UniversityName  string             `json:"university_name"`
	Course          string             `json:"course"`
	Status          int                `json:"status"`
	CitizenIdNumber string             `json:"citizen_id_number"`
	Gender          bool               `json:"gender"`
	DateOfBirth     string             `json:"date_of_birth"`
	Ethnicity       string             `json:"ethnicity"`
	CurrentAddress  string             `json:"current_address"`
	BirthAddress    string             `json:"birth_address"`
	UnionJoinDate   string             `json:"union_join_date"`
	PartyJoinDate   string             `json:"party_join_date"`
	Description     string             `json:"description"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
}

type SearchUserParams struct {
	StudentCode     string             `form:"student_code"`
	FullName        string             `form:"full_name"`
	Email           string             `form:"email"`
	Faculty         string             `form:"faculty_code"`
	Course          string             `form:"course" `
	Status          int                `form:"status"`
	CitizenIdNumber string             `form:"citizen_id_number"`
	Page            int                `form:"page,default=1"`
	PageSize        int                `form:"page_size,default=10"`
	UniversityID    primitive.ObjectID `json:"-"`
}

type UpdateUserRequest struct {
	StudentCode     *string `json:"student_code" binding:"omitempty"`
	FullName        *string `json:"full_name" binding:"omitempty"`
	Email           *string `json:"email" binding:"omitempty,email"`
	FacultyCode     *string `json:"faculty_code" binding:"omitempty"`
	Course          *string `json:"course" binding:"omitempty,courseyear"`
	CitizenIdNumber *string `json:"citizen_id_number" binding:"omitempty,citizenid"`
	Gender          *bool   `json:"gender" binding:"omitempty"`
	DateOfBirth     *string `json:"date_of_birth" binding:"omitempty,dateformat"`
	Ethnicity       *string `json:"ethnicity" binding:"omitempty"`
	CurrentAddress  *string `json:"current_address" binding:"omitempty"`
	BirthAddress    *string `json:"birth_address" binding:"omitempty"`
	UnionJoinDate   *string `json:"union_join_date" binding:"omitempty,dateformat"`
	PartyJoinDate   *string `json:"party_join_date" binding:"omitempty,dateformat"`
	Description     *string `json:"description" binding:"omitempty"`
}
