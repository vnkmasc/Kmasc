package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentCode  string             `bson:"student_code" json:"student_code"`
	FullName     string             `bson:"full_name" json:"full_name"`
	Email        string             `bson:"email" json:"email"`
	FacultyID    primitive.ObjectID `bson:"faculty_id" json:"faculty"`
	UniversityID primitive.ObjectID `bson:"university_id" json:"university_id"`
	Course       string             `bson:"course" json:"course"`
	Status       int                `bson:"status,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateUserRequest struct {
	StudentCode string `json:"student_code" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	FacultyCode string `json:"faculty_code" binding:"required"`
	Course      string `json:"course" binding:"required,courseyear"`
}

type UserResponse struct {
	ID             primitive.ObjectID `json:"id"`
	StudentCode    string             `json:"student_code"`
	FullName       string             `json:"full_name"`
	Email          string             `json:"email"`
	FacultyCode    string             `json:"faculty_code"`
	FacultyName    string             `json:"faculty_name"`
	UniversityCode string             `json:"university_code"`
	UniversityName string             `json:"university_name"`
	Course         string             `json:"course"`
	Status         int                `json:"status"`
}

type SearchUserParams struct {
	StudentCode  string             `form:"student_code"`
	FullName     string             `form:"full_name"`
	Email        string             `form:"email"`
	Faculty      string             `form:"faculty_code"`
	Course       string             `form:"course" `
	Status       int                `form:"status"`
	Page         int                `form:"page,default=1"`
	PageSize     int                `form:"page_size,default=10"`
	UniversityID primitive.ObjectID `json:"-"`
}

type UpdateUserRequest struct {
	StudentCode *string `json:"student_code" binding:"omitempty"`
	FullName    *string `json:"full_name" binding:"omitempty"`
	Email       *string `json:"email" binding:"omitempty,email"`
	FacultyCode *string `json:"faculty_code" binding:"omitempty"`
	Course      *string `json:"course" binding:"omitempty"`
}
