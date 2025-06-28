package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerificationCode struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	Code         string             `bson:"code" json:"code"`
	CanViewScore bool               `bson:"can_view_score" json:"can_view_score"`
	CanViewData  bool               `bson:"can_view_data" json:"can_view_data"`
	CanViewFile  bool               `bson:"can_view_file" json:"can_view_file"`

	ExpiredAt time.Time `bson:"expired_at" json:"expired_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type CreateVerificationCodeRequest struct {
	DurationMinutes int  `json:"duration_minutes" binding:"required,min=1"`
	CanViewScore    bool `json:"can_view_score"`
	CanViewData     bool `json:"can_view_data"`
	CanViewFile     bool `json:"can_view_file"`
}

type VerificationCodeResponse struct {
	ID               primitive.ObjectID `json:"id"`
	Code             string             `json:"code"`
	CanViewScore     bool               `json:"can_view_score"`
	CanViewData      bool               `json:"can_view_data"`
	CanViewFile      bool               `json:"can_view_file"`
	ViewedScore      bool               `json:"viewed_score"`
	ViewedData       bool               `json:"viewed_data"`
	ViewedFile       bool               `json:"viewed_file"`
	ExpiredInMinutes int64              `json:"expired_in_minutes"`
	CreatedAt        time.Time          `json:"created_at"`
}

type VerifyCodeRequest struct {
	Code     string `json:"code" binding:"required"`
	ViewType string `json:"view_type" binding:"required,oneof=score data file"`
}
