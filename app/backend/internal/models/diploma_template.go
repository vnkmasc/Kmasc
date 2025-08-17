package models

import (
	"bytes"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiplomaTemplate struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Description       string             `bson:"description" json:"description"`
	TemplateSampleID  primitive.ObjectID `bson:"template_sample_id" json:"template_sample_id"`
	HashTemplate      string             `bson:"hash_template" json:"hash_template"`
	SignatureOfUni    string             `bson:"signature_of_uni,omitempty" json:"signatureOfUni,omitempty"`
	SignatureOfMinEdu string             `bson:"signature_of_minedu,omitempty" json:"signatureOfMinEdu,omitempty"`
	Status            string             `bson:"status" json:"status"`
	IsLocked          bool               `bson:"is_locked" json:"isLocked"`
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
	UniversityID      primitive.ObjectID `bson:"university_id" json:"universityId"`
	FacultyID         primitive.ObjectID `bson:"faculty_id" json:"facultyId"`
}

type UpdateDiplomaTemplateRequest struct {
	Name             string             `json:"name" binding:"required"`
	Description      string             `json:"description" binding:"required"`
	TemplateSampleID primitive.ObjectID `json:"template_sample_id" binding:"required"`
	FacultyID        primitive.ObjectID `json:"faculty_id" binding:"required"`
}

type TemplateEngine struct{}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

func (te *TemplateEngine) Render(templateContent string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("diploma").Parse(templateContent)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
