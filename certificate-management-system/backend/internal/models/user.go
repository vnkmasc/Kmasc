package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CitizenIdentificationNumber string             `bson:"citizen_identification_number" json:"citizen_identification_number"`
	StudentCode                 string             `bson:"student_code" json:"student_code"`
	FullName                    string             `bson:"full_name" json:"full_name"`
	Gender                      string             `bson:"gender" json:"gender"`                                                                   // Nam / Nữ / Khác
	DateOfBirth                 time.Time          `bson:"date_of_birth" json:"date_of_birth"`                                                     // yyyy-mm-dd
	PlaceOfBirth                string             `bson:"place_of_birth" json:"place_of_birth"`                                                   // Nơi sinh
	Ethnicity                   string             `bson:"ethnicity" json:"ethnicity"`                                                             // Dân tộc
	Hometown                    string             `bson:"hometown" json:"hometown"`                                                               // Quê quán
	CurrentAddress              string             `bson:"current_address" json:"current_address"`                                                 // Nơi ở hiện tại
	JoinYouthUnionDate          *time.Time         `bson:"join_youth_union_date,omitempty" json:"join_youth_union_date,omitempty"`                 // Ngày vào đoàn
	JoinCommunistPartyDate      *time.Time         `bson:"join_communist_party_date,omitempty" json:"join_communist_party_date,omitempty"`         // Ngày vào đảng
	OfficialCommunistPartyDate  *time.Time         `bson:"official_communist_party_date,omitempty" json:"official_communist_party_date,omitempty"` // Ngày vào đảng chính thức

	Email        string             `bson:"email" json:"email"`
	FacultyID    primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`
	UniversityID primitive.ObjectID `bson:"university_id" json:"university_id"`
	Course       string             `bson:"course" json:"course"`                     // Khóa học (VD: K66)
	Status       int                `bson:"status,omitempty" json:"status,omitempty"` // 0 = đang học, 1 = tốt nghiệp, v.v.

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
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
