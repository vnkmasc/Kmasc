package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiploma struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	TemplateID      primitive.ObjectID `bson:"template_id" json:"template_id"` // Liên kết mẫu đã được Bộ duyệt
	UniversityID    primitive.ObjectID `bson:"university_id" json:"university_id"`
	FacultyID       primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	MajorID         primitive.ObjectID `bson:"major_id" json:"major_id"`
	StudentCode     string             `bson:"student_code" json:"student_code"`
	FullName        string             `bson:"full_name" json:"full_name"`
	CertificateType string             `bson:"certificate_type" json:"certificate_type"` // Cử nhân, Thạc sĩ...
	Course          string             `bson:"course" json:"course"`                     // Khóa học (VD: AT18)
	EducationType   string             `bson:"education_type" json:"education_type"`     // Chính quy, Tại chức...
	GPA             float64            `bson:"gpa" json:"gpa"`
	GraduationRank  string             `bson:"graduation_rank" json:"graduation_rank"`
	IssueDate       time.Time          `bson:"issue_date" json:"issue_date"`

	SerialNumber       string `bson:"serial_number" json:"serial_number"`             // Số hiệu
	RegistrationNumber string `bson:"registration_number" json:"registration_number"` // Số vào sổ

	// File văn bằng
	FileLink  string    `bson:"file_link" json:"file_link"` // Link MinIO hoặc CDN
	FileHash  string    `bson:"file_hash" json:"file_hash"` // SHA256 mã băm file PDF
	Signature string    `bson:"signature" json:"signature"` // Chữ ký số
	Signed    bool      `bson:"signed" json:"signed"`       // Đã ký hay chưa
	SignedAt  time.Time `bson:"signed_at,omitempty" json:"signed_at,omitempty"`

	// Blockchain & trạng thái
	OnBlockchain   bool   `bson:"on_blockchain" json:"on_blockchain"`
	BlockchainTxID string `bson:"blockchain_tx_id,omitempty" json:"blockchain_tx_id,omitempty"`

	// Metadata
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
