package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Certificate struct {
	ID              primitive.ObjectID `bson:"_id"`
	UserID          primitive.ObjectID `bson:"user_id"`
	ScoreID         primitive.ObjectID `bson:"score_id"`
	FacultyID       primitive.ObjectID `bson:"faculty_id"`
	UniversityID    primitive.ObjectID `bson:"university_id"`
	StudentCode     string             `bson:"student_code"`
	CertificateType string             `bson:"certificate_type"`
	IsDegree        bool               `bson:"is_degree"`
	Name            string             `bson:"name"`
	SerialNumber    string             `bson:"serial_number"`
	RegNo           string             `bson:"registration_number"`
	Path            string             `bson:"path"`
	IssueDate       time.Time          `bson:"issue_date"`

	Signed   bool      `bson:"signed"`
	SignedAt time.Time `bson:"signed_at,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

	VerificationCode string    `bson:"verification_code,omitempty"`
	CodeExpiredAt    time.Time `bson:"code_expired_at,omitempty"`
}

type CreateCertificateRequest struct {
	IsDegree        bool      `json:"is_degree"`
	StudentCode     string    `json:"student_code" binding:"required"`
	CertificateType string    `json:"certificate_type" binding:"certtype"`
	Name            string    `json:"name"`
	SerialNumber    string    `json:"serial_number"`
	RegNo           string    `json:"reg_no"`
	IssueDate       time.Time `json:"issue_date"`
}

type CertificateResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	StudentCode     string    `json:"student_code,omitempty"`
	StudentName     string    `json:"student_name,omitempty"`
	CertificateType string    `json:"certificate_type,omitempty"`
	Name            string    `json:"name,omitempty"`
	SerialNumber    string    `json:"serial_number,omitempty"`
	RegNo           string    `json:"reg_no,omitempty"`
	Path            string    `json:"path,omitempty"`
	FacultyCode     string    `json:"faculty_code,omitempty"`
	FacultyName     string    `json:"faculty_name,omitempty"`
	UniversityCode  string    `json:"university_code,omitempty"`
	UniversityName  string    `json:"university_name,omitempty"`
	Signed          bool      `json:"signed"`
	IssueDate       string    `json:"issue_date,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SearchCertificateParams struct {
	StudentCode     string `form:"student_code"`
	FacultyCode     string `form:"faculty_code"`
	Course          string `form:"course"`
	Signed          *bool  `form:"signed"`
	CertificateType string `form:"certificate_type"`
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=10"`
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
	} else {
		if req.Name == "" {
			sl.ReportError(req.Name, "name", "Name", "required_if_certificate", "")
		}
		if req.IssueDate.IsZero() {
			sl.ReportError(req.IssueDate, "issue_date", "IssueDate", "required_if_certificate", "")
		}
	}
}

type CertificateSimpleResponse struct {
	ID   string `json:"id"`
	Name string `json:"certificate_name"`
}
