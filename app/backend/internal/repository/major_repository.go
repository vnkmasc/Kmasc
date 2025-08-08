package repository

import (
	"context"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MajorRepository interface {
	Insert(ctx context.Context, major *models.Major) error
	FindByCodeAndFacultyID(ctx context.Context, majorCode string, facultyID primitive.ObjectID) (*models.Major, error)
	GetByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.Major, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Major, error)
}

type majorRepository struct {
	collection *mongo.Collection
}

func NewMajorRepository(db *mongo.Database) MajorRepository {
	return &majorRepository{
		collection: db.Collection("majors"),
	}
}

func (r *majorRepository) Insert(ctx context.Context, major *models.Major) error {
	_, err := r.collection.InsertOne(ctx, major)
	return err
}
func (r *majorRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Major, error) {
	var major models.Major
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&major)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // not found
		}
		return nil, err
	}
	return &major, nil
}
func (r *majorRepository) FindByCodeAndFacultyID(ctx context.Context, majorCode string, facultyID primitive.ObjectID) (*models.Major, error) {
	filter := bson.M{
		"major_code": majorCode,
		"faculty_id": facultyID,
	}
	var major models.Major
	err := r.collection.FindOne(ctx, filter).Decode(&major)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &major, nil
}

func (r *majorRepository) GetByFaculty(ctx context.Context, universityID, facultyID primitive.ObjectID) ([]*models.Major, error) {
	filter := bson.M{
		"university_id": universityID,
		"faculty_id":    facultyID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var majors []*models.Major
	if err := cursor.All(ctx, &majors); err != nil {
		return nil, err
	}

	return majors, nil
}
func (r *majorRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
