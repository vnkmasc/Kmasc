package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	BlockchainTxID string             `bson:"blockchain_tx_id,omitempty"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	FacultyID      primitive.ObjectID `bson:"faculty_id" json:"faculty_id"`
	UniversityID   primitive.ObjectID `bson:"university_id" json:"university_id"`
	IsDegree       bool               `bson:"is_degree" json:"is_degree"` // true: văn bằng, false: chứng chỉ

	StudentCode        string            `bson:"student_code" json:"student_code"`
	CertificateType    string            `bson:"certificate_type" json:"certificate_type"`       // Cử nhân, Thạc sĩ,.....
	Name               string            `bson:"name" json:"name"`                               // Tên văn bằng
	SerialNumber       string            `bson:"serial_number" json:"serial_number"`             // Số hiệu
	RegNo              string            `bson:"registration_number" json:"registration_number"` // Số vào sổ gốc
	Path               string            `bson:"path" json:"path"`
	CertHash           string            `bson:"cert_hash" json:"cert_hash"`
	IssueDate          time.Time         `bson:"issue_date" json:"issue_date"` // Ngày cấp
	HashFile           string            `bson:"hash_file,omitempty" json:"hash_file,omitempty"`
	CertificateFiles   []CertificateFile `bson:"certificate_files,omitempty"`
	Major              string            `bson:"major" json:"major"`                               // Ngành đào tạo
	Course             string            `bson:"course" json:"course"`                             //  Khóa học (VD: AT18)
	GPA                float64           `bson:"gpa" json:"gpa"`                                   //  GPA toàn khóa
	GraduationRank     string            `bson:"graduation_rank" json:"graduation_rank"`           //  Hạng tốt nghiệp: Xuất sắc, Giỏi, Khá...
	EducationType      string            `bson:"education_type" json:"education_type"`             //  Hệ đào tạo: Chính quy, Tại chức...
	PhysicalCopyIssued bool              `bson:"physical_copy_issued" json:"physical_copy_issued"` // Đã phát hành bản giấy
	OnBlockchain       bool              `bson:"on_blockchain" json:"on_blockchain"`               // Đã đẩy lên blockchain
	Signed             bool              `bson:"signed" json:"signed"`
	SignedAt           time.Time         `bson:"signed_at,omitempty" json:"signed_at,omitempty"`
	Description        string            `bson:"description,omitempty" json:"description,omitempty"` // Mô tả thêm

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
type CertificateFile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ID
	FileName  string             `bson:"file_name"`     // Tên file gốc
	FilePath  string             `bson:"file_path"`     // Đường dẫn lưu file mã hóa (trên MinIO)
	AESKey    []byte             `bson:"aes_key"`       // AES key (32 bytes - 256 bit)
	IV        []byte             `bson:"iv"`            // IV (16 bytes - 128 bit)
	CreatedAt time.Time          `bson:"created_at"`    // Thời điểm tạo
}
type CertificateOnChain struct {
	CertID              string `json:"cert_id" bson:"cert_id"`                           // ID của VBCC
	CertHash            string `json:"cert_hash" bson:"cert_hash"`                       // Mã băm các thông tin chính
	HashFile            string `json:"hash_file" bson:"hash_file"`                       // Mã băm file
	UniversitySignature string `json:"university_signature" bson:"university_signature"` // Chữ ký số của trường
	DateOfIssuing       string `json:"date_of_issuing" bson:"date_of_issuing"`           // Ngày cấp
	SerialNumber        string `bson:"serial_number" json:"serial_number"`               // Số hiệu
	RegNo               string `bson:"registration_number" json:"registration_number"`   // Số vào sổ gốc
	Version             int    `json:"version" bson:"version"`                           // Phiên bản VBCC
	UpdatedDate         string `json:"updated_date" bson:"updated_date"`                 // Ngày sửa đổi
}

type CreateCertificateRequest struct {
	StudentCode     string    `json:"student_code" binding:"required"`
	IsDegree        bool      `json:"is_degree"`                 // true: văn bằng, false: chứng chỉ
	CertificateType string    `json:"certificate_type"`          // Bắt buộc nếu là văn bằng
	Course          string    `json:"course,omitempty"`          // Bắt buộc nếu là văn bằng
	GraduationRank  string    `json:"graduation_rank,omitempty"` // Optional
	EducationType   string    `json:"education_type,omitempty"`  // Optional
	Description     string    `json:"description,omitempty"`
	Name            string    `json:"name" binding:"required"`          // Tên văn bằng / chứng chỉ
	SerialNumber    string    `json:"serial_number" binding:"required"` // Số hiệu
	RegNo           string    `json:"reg_no" binding:"required"`        // Số vào sổ
	Major           string    `json:"major,omitempty"`                  // Bắt buộc nếu là văn bằng
	GPA             float64   `json:"gpa,omitempty"`                    // Optional
	IssueDate       time.Time `json:"issue_date" binding:"required"`    // Ngày cấp
}

type CertificateResponse struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"user_id"`
	StudentCode        string  `json:"student_code,omitempty"`
	StudentName        string  `json:"student_name,omitempty"`
	CertificateType    string  `json:"certificate_type,omitempty"`
	Name               string  `json:"name,omitempty"`
	SerialNumber       string  `json:"serial_number,omitempty"`
	RegNo              string  `json:"reg_no,omitempty"`
	Path               string  `json:"path,omitempty"`
	FacultyCode        string  `json:"faculty_code,omitempty"`
	FacultyName        string  `json:"faculty_name,omitempty"`
	UniversityCode     string  `json:"university_code,omitempty"`
	UniversityName     string  `json:"university_name,omitempty"`
	HashFile           string  `json:"hash_file,omitempty"`
	Major              string  `json:"major,omitempty"`           // Ngành đào tạo
	Course             string  `json:"course,omitempty"`          // Khóa học (VD: AT18)
	GPA                float64 `json:"gpa,omitempty"`             // GPA toàn khóa
	GraduationRank     string  `json:"graduation_rank,omitempty"` // Hạng tốt nghiệp: Xuất sắc, Giỏi, Khá...
	EducationType      string  `json:"education_type,omitempty"`  // Hệ đào tạo
	Signed             bool    `json:"signed"`
	PhysicalCopyIssued bool    `json:"physical_copy_issued"`  // Đã phát hành bản giấy
	OnBlockchain       bool    `json:"on_blockchain"`         // Đã đẩy lên blockchain
	IssueDate          string  `json:"issue_date,omitempty"`  // Định dạng ISO hoặc "02/01/2006"
	Description        string  `json:"description,omitempty"` // Mô tả thêm

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchCertificateParams struct {
	StudentCode     string `form:"student_code"`
	FacultyCode     string `form:"faculty_code"`
	Course          string `form:"course"`
	Signed          *bool  `form:"signed"`
	CertificateType string `form:"certificate_type"`
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=10"`
	SortOrder       string `form:"sort_order"` // "asc" | "desc"

}

func ValidateCreateCertificateRequest(sl validator.StructLevel) {
	req := sl.Current().Interface().(CreateCertificateRequest)

	if req.IsDegree {
		if req.CertificateType == "" {
			sl.ReportError(req.CertificateType, "certificate_type", "CertificateType", "required_if_degree", "")
		}
		if req.SerialNumber == "" {
			sl.ReportError(req.SerialNumber, "serial_number", "SerialNumber", "required_if_degree", "")
		}
		if req.RegNo == "" {
			sl.ReportError(req.RegNo, "reg_no", "RegNo", "required_if_degree", "")
		}
		if req.IssueDate.IsZero() {
			sl.ReportError(req.IssueDate, "issue_date", "IssueDate", "required_if_degree", "")
		}
	}
}

type CertificateSimpleResponse struct {
	ID   string `json:"id"`
	Name string `json:"certificate_name"`
}

func NewCertificate(req *CreateCertificateRequest, user *User, universityID primitive.ObjectID) *Certificate {
	now := time.Now()
	return &Certificate{
		ID:              primitive.NewObjectID(),
		UserID:          user.ID,
		FacultyID:       user.FacultyID,
		UniversityID:    universityID,
		StudentCode:     user.StudentCode,
		IsDegree:        req.IsDegree,
		Name:            req.Name,
		CertificateType: req.CertificateType,
		SerialNumber:    req.SerialNumber,
		RegNo:           req.RegNo,
		IssueDate:       req.IssueDate,
		Major:           req.Major,
		Course:          req.Course,
		GPA:             req.GPA,
		GraduationRank:  req.GraduationRank,
		EducationType:   req.EducationType,
		Signed:          false,
		Description:     req.Description,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
