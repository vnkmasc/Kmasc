package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateRepository interface {
	FindByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	UpdateDiplomaTemplateByID(ctx context.Context, template *models.DiplomaTemplate) error
	UpdateStatusAndMinEduSignatureByID(ctx context.Context, id primitive.ObjectID, newStatus string, signatureOfMinEdu string) error
	FindByIDAndUniversity(ctx context.Context, templateID, universityID primitive.ObjectID) (*models.DiplomaTemplate, error)
	Create(ctx context.Context, template *models.DiplomaTemplate) error
	UpdateIfNotLocked(ctx context.Context, id primitive.ObjectID, updated *models.DiplomaTemplate) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.DiplomaTemplate, error)
	FindByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindByFacultyIDs(ctx context.Context, facultyIDs []primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindPendingByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindPendingByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	UpdateStatusAndSignatureByID(ctx context.Context, id primitive.ObjectID, newStatus string, signatureOfUni string) error
	Update(ctx context.Context, template *models.DiplomaTemplate) error
	LockTemplate(ctx context.Context, templateID primitive.ObjectID) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
	FindByTemplateSampleID(ctx context.Context, sampleID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
}

type templateRepository struct {
	collection *mongo.Collection
}

func NewTemplateRepository(db *mongo.Database) TemplateRepository {
	return &templateRepository{
		collection: db.Collection("diploma_templates"),
	}
}

func (r *templateRepository) FindByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{"university_id": universityID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*models.DiplomaTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *templateRepository) UpdateDiplomaTemplateByID(
	ctx context.Context,
	template *models.DiplomaTemplate,
) error {
	filter := bson.M{"_id": template.ID, "is_locked": false}
	update := bson.M{
		"$set": bson.M{
			"name":               template.Name,
			"description":        template.Description,
			"template_sample_id": template.TemplateSampleID,
			"faculty_id":         template.FacultyID,
			"updated_at":         template.UpdatedAt,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("template is locked or not found")
	}
	return nil
}

func (r *templateRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return r.collection.DeleteOne(ctx, filter)
}

func (r *templateRepository) Update(ctx context.Context, template *models.DiplomaTemplate) error {
	filter := bson.M{"_id": template.ID}
	update := bson.M{
		"$set": bson.M{
			"name":          template.Name,
			"description":   template.Description,
			"hash_template": template.HashTemplate,
			"updated_at":    template.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *templateRepository) LockTemplate(ctx context.Context, templateID primitive.ObjectID) error {
	filter := bson.M{"_id": templateID}
	update := bson.M{
		"$set": bson.M{
			"is_locked":  true,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *templateRepository) UpdateStatusAndMinEduSignatureByID(
	ctx context.Context,
	id primitive.ObjectID,
	newStatus string,
	signatureOfMinEdu string,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":              newStatus,
		"signature_of_minedu": signatureOfMinEdu,
		"updated_at":          time.Now(),
	}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *templateRepository) FindByFacultyIDs(ctx context.Context, facultyIDs []primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{"faculty_id": bson.M{"$in": facultyIDs}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*models.DiplomaTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *templateRepository) FindByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{
		"university_id": universityID,
		"faculty_id":    facultyID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*models.DiplomaTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

func (r *templateRepository) Create(ctx context.Context, template *models.DiplomaTemplate) error {
	_, err := r.collection.InsertOne(ctx, template)
	return err
}

func (r *templateRepository) FindByIDAndUniversity(ctx context.Context, templateID, universityID primitive.ObjectID) (*models.DiplomaTemplate, error) {
	filter := bson.M{"_id": templateID, "university_id": universityID}
	var template models.DiplomaTemplate
	err := r.collection.FindOne(ctx, filter).Decode(&template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.DiplomaTemplate, error) {
	var template models.DiplomaTemplate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("template with id %s not found", id.Hex())
		}
		return nil, err
	}
	return &template, nil
}

func (r *templateRepository) UpdateIfNotLocked(ctx context.Context, id primitive.ObjectID, updated *models.DiplomaTemplate) error {
	filter := bson.M{"_id": id, "is_locked": false}
	update := bson.M{
		"$set": bson.M{
			"name":        updated.Name,
			"description": updated.Description,
			"updated_at":  time.Now(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("template not found or already locked")
	}
	return nil
}

func (r *templateRepository) UpdateStatusAndSignatureByID(
	ctx context.Context,
	id primitive.ObjectID,
	newStatus string,
	signatureOfUni string,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":           newStatus,
		"signature_of_uni": signatureOfUni,
		"updated_at":       time.Now(),
	}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Lấy tất cả template đang Pending theo khoa
func (r *templateRepository) FindPendingByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{
		"university_id": universityID,
		"faculty_id":    facultyID,
		"status":        "PENDING",
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*models.DiplomaTemplate
	for cursor.Next(ctx) {
		var tmpl models.DiplomaTemplate
		if err := cursor.Decode(&tmpl); err != nil {
			return nil, err
		}
		result = append(result, &tmpl)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Lấy tất cả template Pending toàn trường
func (r *templateRepository) FindPendingByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{
		"university_id": universityID,
		"status":        "PENDING",
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*models.DiplomaTemplate
	for cursor.Next(ctx) {
		var tmpl models.DiplomaTemplate
		if err := cursor.Decode(&tmpl); err != nil {
			return nil, err
		}
		result = append(result, &tmpl)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *templateRepository) FindByTemplateSampleID(ctx context.Context, sampleID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{
		"template_sample_id": sampleID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find templates: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*models.DiplomaTemplate
	for cursor.Next(ctx) {
		var tmpl models.DiplomaTemplate
		if err := cursor.Decode(&tmpl); err != nil {
			return nil, fmt.Errorf("failed to decode template: %w", err)
		}
		templates = append(templates, &tmpl)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return templates, nil
}
