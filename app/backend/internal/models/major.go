package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Major struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MajorCode    string             `bson:"major_code" json:"major_code"`       // Mã chuyên ngành (VD: "CNPM")
	MajorName    string             `bson:"major_name" json:"major_name"`       // Tên chuyên ngành (VD: "Công nghệ phần mềm")
	FacultyID    primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`       // ID khoa
	UniversityID primitive.ObjectID `bson:"university_id" json:"university_id"` // ID trường
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	Quota        int                `bson:"quota" json:"quota"`                 // Chỉ tiêu tuyển sinh
	DiplomaCount int                `bson:"diploma_count" json:"diploma_count"` // Số bằng đã cấp
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
