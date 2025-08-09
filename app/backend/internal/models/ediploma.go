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

	SignatureOfUni    string `bson:"signature_of_uni,omitempty" json:"signatureOfUni,omitempty"`
	SignatureOfMinEdu string `bson:"signature_of_minedu,omitempty" json:"signatureOfMinEdu,omitempty"`
	Status            string `bson:"status" json:"status"` // PENDING, VERIFIED,...
	IsLocked          bool   `bson:"is_locked" json:"isLocked"`

	// Metadata
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type EDiplomaDTO struct {
	ID                 primitive.ObjectID `json:"id"`
	TemplateID         primitive.ObjectID `json:"template_id"`
	UniversityID       primitive.ObjectID `json:"university_id"`
	UniversityCode     string             `json:"university_code"`
	UniversityName     string             `json:"university_name"`
	FacultyID          primitive.ObjectID `json:"faculty_id"`
	FacultyCode        string             `json:"faculty_code"`
	FacultyName        string             `json:"faculty_name"`
	MajorID            primitive.ObjectID `json:"major_id"`
	MajorCode          string             `json:"major_code"`
	MajorName          string             `json:"major_name"`
	UserID             primitive.ObjectID `json:"user_id"`
	StudentCode        string             `json:"student_code"`
	FullName           string             `json:"full_name"`
	CertificateType    string             `json:"certificate_type"`
	Course             string             `json:"course"`
	EducationType      string             `json:"education_type"`
	GPA                float64            `json:"gpa"`
	GraduationRank     string             `json:"graduation_rank"`
	IssueDate          time.Time          `json:"issue_date"`
	SerialNumber       string             `json:"serial_number"`
	RegistrationNumber string             `json:"registration_number"`
	FileLink           string             `json:"file_link"`
	FileHash           string             `json:"file_hash"`
	Signature          string             `json:"signature"`
	Signed             bool               `json:"signed"`
	SignedAt           time.Time          `json:"signed_at"`
	OnBlockchain       bool               `json:"on_blockchain"`
	SignatureOfUni     string             `json:"signatureOfUni,omitempty"`
	SignatureOfMinEdu  string             `json:"signatureOfMinEdu,omitempty"`
	Status             string             `json:"status"`
	IsLocked           bool               `json:"isLocked"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type EDiplomaSearchFilter struct {
	StudentCode     string `json:"student_code"`
	FacultyCode     string `json:"faculty_code"`
	CertificateType string `json:"certificate_type"`
	Course          string `json:"course"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
}
