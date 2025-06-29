package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RewardDiscipline struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	DecisionNumber  string             `bson:"decision_number" json:"decision_number"`
	Description     string             `bson:"description" json:"description"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	IsDiscipline    bool               `bson:"is_discipline" json:"is_discipline"`
	DisciplineLevel *int               `bson:"discipline_level,omitempty" json:"discipline_level,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateRewardDisciplineRequest struct {
	Name            string `json:"name" binding:"required"`
	DecisionNumber  string `json:"decision_number" binding:"required"`
	Description     string `json:"description"`
	StudentCode     string `json:"student_code" binding:"required"`
	IsDiscipline    bool   `json:"is_discipline"`
	DisciplineLevel *int   `json:"discipline_level" binding:"omitempty,disciplinelevel"`
}

type UpdateRewardDisciplineRequest struct {
	Name            *string `json:"name" binding:"omitempty"`
	DecisionNumber  *string `json:"decision_number" binding:"omitempty"`
	Description     *string `json:"description" binding:"omitempty"`
	StudentCode     *string `json:"student_code" binding:"omitempty"`
	IsDiscipline    *bool   `json:"is_discipline" binding:"omitempty"`
	DisciplineLevel *int    `json:"discipline_level" binding:"omitempty,disciplinelevel"`
}

type RewardDisciplineResponse struct {
	ID              primitive.ObjectID `json:"id"`
	Name            string             `json:"name"`
	DecisionNumber  string             `json:"decision_number"`
	Description     string             `json:"description"`
	StudentCode     string             `json:"student_code"`
	StudentName     string             `json:"student_name"`
	IsDiscipline    bool               `json:"is_discipline"`
	DisciplineLevel *int               `json:"discipline_level,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type SearchRewardDisciplineParams struct {
	Name           string `form:"name"`
	DecisionNumber string `form:"decision_number"`
	StudentCode    string `form:"student_code"`
	IsDiscipline   *bool  `form:"is_discipline"`
	Page           int    `form:"page,default=1"`
	PageSize       int    `form:"page_size,default=10"`
}
