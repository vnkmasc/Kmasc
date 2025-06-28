package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID    primitive.ObjectID `bson:"student_id" json:"student_id"`
	UniversityID primitive.ObjectID `bson:"university_id,omitempty" json:"university_id,omitempty"`

	StudentEmail  string `bson:"student_email" json:"student_email"`   // Email sinh viên (VD: @actvn.edu.vn)
	PersonalEmail string `bson:"personal_email" json:"personal_email"` // Email cá nhân dùng đăng nhập
	PasswordHash  string `bson:"password_hash" json:"-"`               // Không trả ra ngoài

	Role        string `bson:"role" json:"role"`                                   // VD: student / admin / staff
	Description string `bson:"description,omitempty" json:"description,omitempty"` // Mô tả thêm (nếu có)

	CreatedAt time.Time `bson:"created_at" json:"created_at"` // Ngày tạo
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"` // Ngày cập nhật cuối
}

type AccountResponse struct {
	ID            primitive.ObjectID  `json:"id"`
	StudentID     *primitive.ObjectID `json:"student_id,omitempty"`
	UniversityID  *primitive.ObjectID `json:"university_id,omitempty"`
	StudentEmail  string              `json:"student_email,omitempty"`
	PersonalEmail string              `json:"personal_email"`
	CreatedAt     string              `json:"created_at"`
	Role          string              `json:"role"`
}

type OTP struct {
	Email     string    `bson:"email"`
	Code      string    `bson:"code"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type RequestOTPInput struct {
	StudentEmail string `json:"student_email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	StudentEmail string `json:"student_email" binding:"required,email"`
	OTP          string `json:"otp" binding:"required,len=6"`
}

type VerifyOTPResponse struct {
	UserID string `json:"user_id"`
}

type RegisterRequest struct {
	UserID        string `json:"user_id" binding:"required"`
	PersonalEmail string `json:"personal_email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
