package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EDiploma struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	CertificateID   primitive.ObjectID `bson:"certificate_id" json:"certificate_id"` // Liên kết đến văn bằng
	TemplateID      primitive.ObjectID `bson:"template_id" json:"template_id"`       // Liên kết mẫu đã được Bộ duyệt
	Name            string             `bson:"name" json:"name"`                     // Tên văn bằng
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

	SerialNumber       string    `bson:"serial_number" json:"serial_number"`             // Số hiệu
	RegistrationNumber string    `bson:"registration_number" json:"registration_number"` // Số vào sổ
	EDiplomaFileLink   string    `bson:"ediploma_file_link" json:"ediploma_file_link"`   // Đường dẫn đến file văn bằng số
	EDiplomaFileHash   string    `bson:"ediploma_file_hash" json:"ediploma_file_hash"`   // Mã băm của văn bằng số
	Signature          string    `bson:"signature" json:"signature"`                     // Chữ ký số
	Signed             bool      `bson:"signed" json:"signed"`                           // Đã ký hay chưa
	DataEncrypted      bool      `bson:"data_encrypted" json:"data_encrypted"`           //Đã mã hóa dữ liệu hay chưa
	Issued             bool      `bson:"issued" json:"issued"`                           // Đã cấp bằng số hay chưa
	SignedAt           time.Time `bson:"signed_at,omitempty" json:"signed_at,omitempty"`

	// Blockchain & trạng thái
	OnBlockchain   bool        `bson:"on_blockchain" json:"on_blockchain"`
	MerkleProof    []ProofNode `bson:"merkle_proof,omitempty" json:"merkle_proof,omitempty"`
	BlockchainTxID string      `bson:"blockchain_tx_id,omitempty" json:"blockchain_tx_id,omitempty"`

	SignatureOfUni    string `bson:"signature_of_uni,omitempty" json:"signatureOfUni,omitempty"`
	SignatureOfMinEdu string `bson:"signature_of_minedu,omitempty" json:"signatureOfMinEdu,omitempty"`

	// Metadata
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type EDiplomaResponse struct {
	ID                 primitive.ObjectID `json:"id"`
	CertificateID      primitive.ObjectID `json:"certificate_id"`
	UniversityID       primitive.ObjectID `json:"university_id"`
	Name               string             `json:"name"`
	TemplateName       string             `json:"template_name"`
	UniversityCode     string             `json:"university_code"`
	UniversityName     string             `json:"university_name"`
	FacultyID          primitive.ObjectID `json:"faculty_id"`
	FacultyCode        string             `json:"faculty_code"`
	FacultyName        string             `json:"faculty_name"`
	StudentName        string             `json:"student_name"`
	StudentCode        string             `json:"student_code"`
	FullName           string             `json:"full_name"`
	CertificateType    string             `json:"certificate_type"`
	Course             string             `json:"course"`
	EducationType      string             `json:"education_type"`
	GPA                float64            `json:"gpa"`
	GraduationRank     string             `json:"graduation_rank"`
	IssueDate          string             `json:"issue_date"`
	SerialNumber       string             `json:"serial_number"`
	RegistrationNumber string             `json:"registration_number"`
	Issued             bool               `json:"issued"`
	Signed             bool               `json:"signed"`
	DataEncrypted      bool               `json:"data_encrypted"`
	OnBlockchain       bool               `json:"on_blockchain"`
}
type EDiplomaSearchFilter struct {
	UniversityID    string `json:"university_id"`
	FacultyID       string `json:"faculty_id"`
	CertificateType string `json:"certificate_type"`
	Course          string `json:"course"`
	Issued          *bool  `json:"issued"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
}

type EDiplomaBatchOnChain struct {
	BatchID           string `json:"batch_id"`
	UniversityID      string `json:"university_id"`
	FacultyID         string `json:"faculty_id"`
	CertificateType   string `json:"certificate_type"`
	Course            string `json:"course"`
	AggregateInfoHash string `json:"aggregate_info_hash"`
	AggregateFileHash string `json:"aggregate_file_hash"`
	Count             int    `json:"count"`
}
