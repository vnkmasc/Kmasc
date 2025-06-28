package repository

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UniversityRepository interface {
	CheckUniversityConflicts(ctx context.Context, universityName, emailDomain, universityCode string) (string, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.University, error)
	FindByCode(ctx context.Context, code string) (*models.University, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateUniversity(ctx context.Context, uni *models.University) error
	GetAllUniversities(ctx context.Context) ([]*models.University, error)
	GetUniversitiesByStatus(ctx context.Context, status string) ([]*models.University, error)
	GetUniversityByCode(ctx context.Context, code string) (*models.University, error)
}

type universityRepository struct {
	col *mongo.Collection
}

func NewUniversityRepository(db *mongo.Database) UniversityRepository {
	col := db.Collection("universities")
	return &universityRepository{col: col}
}
func (r *universityRepository) GetUniversityByCode(ctx context.Context, code string) (*models.University, error) {
	var university models.University
	err := r.col.FindOne(ctx, bson.M{"university_code": code}).Decode(&university)
	if err != nil {
		return nil, err
	}
	return &university, nil
}

func (r *universityRepository) CreateUniversity(ctx context.Context, uni *models.University) error {
	_, err := r.col.InsertOne(ctx, uni)
	return err
}

func (r *universityRepository) CheckUniversityConflicts(ctx context.Context, universityName, emailDomain, universityCode string) (string, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"university_name": universityName})
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "university_name", nil
	}

	count, err = r.col.CountDocuments(ctx, bson.M{"email_domain": emailDomain})
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "email_domain", nil
	}

	count, err = r.col.CountDocuments(ctx, bson.M{"university_code": universityCode})
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "university_code", nil
	}

	return "", nil
}

func (r *universityRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.University, error) {
	var university models.University
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&university)
	if err != nil {
		return nil, err
	}
	return &university, nil
}

func (r *universityRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	})
	return err
}

func (r *universityRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *universityRepository) GetAllUniversities(ctx context.Context) ([]*models.University, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var universities []*models.University
	if err := cursor.All(ctx, &universities); err != nil {
		return nil, err
	}
	return universities, nil
}
func (r *universityRepository) GetUniversitiesByStatus(ctx context.Context, status string) ([]*models.University, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var universities []*models.University
	for cursor.Next(ctx) {
		var uni models.University
		if err := cursor.Decode(&uni); err != nil {
			return nil, err
		}
		universities = append(universities, &uni)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return universities, nil
}
func (r *universityRepository) FindByCode(ctx context.Context, code string) (*models.University, error) {
	filter := bson.M{"university_code": code}
	var university models.University
	err := r.col.FindOne(ctx, filter).Decode(&university)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &university, nil
}
