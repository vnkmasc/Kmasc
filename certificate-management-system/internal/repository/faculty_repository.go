package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/common"
	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacultyRepository interface {
	FindByCodeAndUniversityID(ctx context.Context, facultyCode string, universityID primitive.ObjectID) (*models.Faculty, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Faculty, error)
	Create(ctx context.Context, faculty *models.Faculty) error
	FindAllByUniversityID(ctx context.Context, universityID primitive.ObjectID) ([]*models.Faculty, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	UpdateFaculty(ctx context.Context, id primitive.ObjectID, update bson.M) error
	FindByFacultyCode(ctx context.Context, code string) (*models.Faculty, error)
}

type facultyRepository struct {
	col *mongo.Collection
}

func NewFacultyRepository(db *mongo.Database) FacultyRepository {
	return &facultyRepository{
		col: db.Collection("faculties"),
	}
}

func (r *facultyRepository) FindByCodeAndUniversityID(ctx context.Context, code string, universityID primitive.ObjectID) (*models.Faculty, error) {
	filter := bson.M{
		"faculty_code":  bson.M{"$regex": code, "$options": "i"},
		"university_id": universityID,
	}
	var faculty models.Faculty
	err := r.col.FindOne(ctx, filter).Decode(&faculty)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *facultyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&faculty)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &faculty, nil
}

func (r *facultyRepository) Create(ctx context.Context, faculty *models.Faculty) error {
	_, err := r.col.InsertOne(ctx, faculty)
	return err
}
func (r *facultyRepository) FindAllByUniversityID(ctx context.Context, universityID primitive.ObjectID) ([]*models.Faculty, error) {
	filter := bson.M{"university_id": universityID}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var faculties []*models.Faculty
	for cursor.Next(ctx) {
		var faculty models.Faculty
		if err := cursor.Decode(&faculty); err != nil {
			return nil, err
		}
		faculties = append(faculties, &faculty)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return faculties, nil
}
func (r *facultyRepository) UpdateFaculty(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}
	_, err := r.col.UpdateOne(ctx, filter, updateDoc)
	return err
}

func (r *facultyRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return common.ErrFacultyNotFound
	}
	return nil
}
func (r *facultyRepository) FindByFacultyCode(ctx context.Context, code string) (*models.Faculty, error) {
	filter := bson.M{"faculty_code": code}
	var faculty models.Faculty
	err := r.col.FindOne(ctx, filter).Decode(&faculty)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &faculty, nil
}
