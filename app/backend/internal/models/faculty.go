package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Faculty struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FacultyCode  string             `bson:"faculty_code,omitempty" json:"faculty_code"`
	FacultyName  string             `bson:"faculty_name" json:"faculty_name"`
	UniversityID primitive.ObjectID `bson:"university_id" json:"university_id"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type CreateFacultyRequest struct {
	FacultyCode string `json:"faculty_code" binding:"required"`
	FacultyName string `json:"faculty_name" binding:"required"`
}
type FacultyResponse struct {
	ID          primitive.ObjectID `json:"id"`
	FacultyCode string             `json:"faculty_code"`
	FacultyName string             `json:"faculty_name"`
	CreatedAt   string             `json:"created_at"`
}
type UpdateFacultyRequest struct {
	FacultyCode string `json:"faculty_code" binding:"required"`
	FacultyName string `json:"faculty_name" binding:"required"`
}
