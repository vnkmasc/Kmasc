package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	StudentID     primitive.ObjectID `bson:"student_id"`
	UniversityID  primitive.ObjectID `bson:"university_id,omitempty"`
	StudentEmail  string             `bson:"student_email"`  // Email @actvn.edu.vn
	PersonalEmail string             `bson:"personal_email"` // Email Gmail, dùng để login sau này
	PasswordHash  string             `bson:"password_hash"`
	CreatedAt     time.Time          `bson:"created_at"`
	Role          string             `bson:"role"`
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
