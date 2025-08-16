package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TemplateSample struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	HTMLContent  string             `bson:"html_content"`
	UniversityID primitive.ObjectID `bson:"university_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}
