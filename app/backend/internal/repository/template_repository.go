package repository

import (
	"context"
	"errors"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *models.DiplomaTemplate) error
	UpdateIfNotLocked(ctx context.Context, id primitive.ObjectID, updated *models.DiplomaTemplate) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.DiplomaTemplate, error)
	LockTemplate(ctx context.Context, id primitive.ObjectID) error
	FindByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindByFacultyIDs(ctx context.Context, facultyIDs []primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindPendingByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	FindPendingByUniversityAndFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	UpdateStatusAndSignatureByID(ctx context.Context, id primitive.ObjectID, newStatus string, signatureOfUni string) error
	UpdateStatusAndMinEduSignatureByID(ctx context.Context, id primitive.ObjectID, status, signature string) error
	FindSignedByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error)
	VerifyTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) error
	UpdateIsLocked(ctx context.Context, id primitive.ObjectID, isLocked bool) error
	Update(ctx context.Context, template *models.DiplomaTemplate) error
}

type templateRepository struct {
	collection *mongo.Collection
}

func NewTemplateRepository(db *mongo.Database) TemplateRepository {
	return &templateRepository{
		collection: db.Collection("diploma_templates"),
	}
}

func (r *templateRepository) Update(ctx context.Context, template *models.DiplomaTemplate) error {
	filter := bson.M{"_id": template.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        template.Name,
			"description": template.Description,
			"file_link":   template.FileLink,
			"hash":        template.Hash,
			"updated_at":  template.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *templateRepository) UpdateIsLocked(ctx context.Context, id primitive.ObjectID, isLocked bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"isLocked": isLocked}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *templateRepository) VerifyTemplatesByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) error {
	filter := bson.M{
		"university_id": universityID,
		"faculty_id":    facultyID,
		"status":        "SIGNED_BY_UNI",
	}
	update := bson.M{
		"$set": bson.M{
			"status":              "SIGNED_BY_MINEDU",
			"signature_of_minedu": "SIMULATED_SIGNATURE_MINEDU",
			"updated_at":          time.Now(),
		},
	}
	_, err := r.collection.UpdateMany(ctx, filter, update)
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

func (r *templateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.DiplomaTemplate, error) {
	var template models.DiplomaTemplate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&template)
	if err != nil {
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

func (r *templateRepository) LockTemplate(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"is_locked":  true,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateByID(ctx, id, update)
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
func (r *templateRepository) UpdateStatusAndMinEduSignatureByID(ctx context.Context, id primitive.ObjectID, status, signature string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":              status,
			"signature_of_minedu": signature,
			"updated_at":          time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *templateRepository) FindSignedByUniversity(ctx context.Context, universityID primitive.ObjectID) ([]*models.DiplomaTemplate, error) {
	filter := bson.M{
		"university_id":    universityID,
		"status":           "SIGNED_BY_UNI",
		"signature_of_uni": bson.M{"$ne": nil},
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
