package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type University struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UniversityName string             `bson:"university_name"`
	UniversityCode string             `bson:"university_code"`
	Address        string             `bson:"address"`
	EmailDomain    string             `bson:"email_domain"`
	Status         string             `bson:"status"` // "pending", "approved", "rejected"
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
type CreateUniversityRequest struct {
	UniversityName string `json:"university_name" binding:"required"`
	UniversityCode string `json:"university_code" binding:"required"`
	Address        string `json:"address" binding:"required"`
	EmailDomain    string `json:"email_domain" binding:"required,email"`
}
type UniversityResponse struct {
	ID             string `json:"id"`
	UniversityName string `json:"university_name"`
	UniversityCode string `json:"university_code"`
	EmailDomain    string `json:"email_domain"`
	Address        string `json:"address"`
	Status         string `json:"status"`
}

type ApproveOrRejectUniversityRequest struct {
	UniversityID string `json:"university_id" binding:"required"`
	Action       string `json:"action" binding:"required,oneof=approve reject"`
}
