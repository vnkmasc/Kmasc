package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EDiplomaRepository interface {
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.EDiploma, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error
	UpdateFields(ctx context.Context, id primitive.ObjectID, updates bson.M) error
	FindByStudentCode(ctx context.Context, studentCode string) (*models.EDiploma, error)
	FindByDynamicFilter(ctx context.Context, filter bson.M) ([]*models.EDiploma, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.EDiploma, error)
	Save(ctx context.Context, ediploma *models.EDiploma) error
	GetByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]*models.EDiploma, error)
	SearchByFilters(ctx context.Context, filter models.EDiplomaSearchFilter) ([]*models.EDiploma, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, ed *models.EDiploma) error
	FindByStudentCodeAndFacultyID(ctx context.Context, studentCode string, facultyID primitive.ObjectID) (*models.EDiploma, error)
}

type eDiplomaRepository struct {
	db          *mongo.Collection
	facultyRepo FacultyRepository
}

func NewEDiplomaRepository(db *mongo.Database, facultyRepo FacultyRepository) EDiplomaRepository {
	return &eDiplomaRepository{
		db:          db.Collection("ediplomas"),
		facultyRepo: facultyRepo,
	}
}

func (r *eDiplomaRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.db.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)
	return err
}

func (r *eDiplomaRepository) SearchByFilters(ctx context.Context, filter models.EDiplomaSearchFilter) ([]*models.EDiploma, int64, error) {
	bsonFilter := bson.M{}
	if filter.UniversityID != "" {
		universityID, err := primitive.ObjectIDFromHex(filter.UniversityID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid university_id: %w", err)
		}
		bsonFilter["university_id"] = universityID
	}
	if filter.FacultyID != "" {
		facultyID, err := primitive.ObjectIDFromHex(filter.FacultyID)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid faculty_id: %w", err)
		}
		bsonFilter["faculty_id"] = facultyID
	}

	if filter.CertificateType != "" {
		bsonFilter["certificate_type"] = filter.CertificateType
	}
	if filter.Course != "" {
		bsonFilter["course"] = filter.Course
	}
	if filter.Issued != nil {
		bsonFilter["issued"] = *filter.Issued
	}

	// Đếm tổng số kết quả
	total, err := r.db.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Tính skip/limit
	skip := int64((filter.Page - 1) * filter.PageSize)
	limit := int64(filter.PageSize)

	findOpts := options.Find().SetSkip(skip).SetLimit(limit)

	cursor, err := r.db.Find(ctx, bsonFilter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*models.EDiploma
	for cursor.Next(ctx) {
		var e models.EDiploma
		if err := cursor.Decode(&e); err != nil {
			return nil, 0, err
		}
		results = append(results, &e)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *eDiplomaRepository) UpdateFields(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	updates["updated_at"] = time.Now()
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

func (r *eDiplomaRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.EDiploma, error) {
	var diploma models.EDiploma
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&diploma)
	if err != nil {
		return nil, err
	}
	return &diploma, nil
}

func (r *eDiplomaRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.EDiploma, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var diplomas []*models.EDiploma
	if err := cursor.All(ctx, &diplomas); err != nil {
		return nil, err
	}
	return diplomas, nil
}

func (r *eDiplomaRepository) FindByDynamicFilter(ctx context.Context, filter bson.M) ([]*models.EDiploma, error) {
	var results []*models.EDiploma

	// Nếu không truyền filter thì mặc định lấy tất cả
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ed models.EDiploma
		if err := cursor.Decode(&ed); err != nil {
			return nil, err
		}
		results = append(results, &ed)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *eDiplomaRepository) FindByStudentCodeAndFacultyID(ctx context.Context, studentCode string, facultyID primitive.ObjectID) (*models.EDiploma, error) {
	filter := bson.M{
		"student_code": studentCode,
		"faculty_id":   facultyID,
	}

	var ed models.EDiploma
	err := r.db.FindOne(ctx, filter).Decode(&ed)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ed, nil
}

func (r *eDiplomaRepository) Update(ctx context.Context, id primitive.ObjectID, ed *models.EDiploma) error {
	update := bson.M{
		"$set": bson.M{
			"template_id":         ed.TemplateID,
			"signature_of_uni":    ed.SignatureOfUni,
			"signature_of_minedu": ed.SignatureOfMinEdu,
			"issue_date":          ed.IssueDate,
			"issued":              ed.Issued,
			"on_blockchain":       ed.OnBlockchain,
			"updated_at":          time.Now(),
		},
	}
	_, err := r.db.UpdateByID(ctx, id, update)
	return err
}

func (r *eDiplomaRepository) FindByStudentCode(ctx context.Context, studentCode string) (*models.EDiploma, error) {
	var ediploma models.EDiploma

	filter := bson.M{
		"student_code": studentCode,
	}

	err := r.db.FindOne(ctx, filter).Decode(&ediploma)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy
		}
		return nil, err
	}

	return &ediploma, nil
}

func (r *eDiplomaRepository) Save(ctx context.Context, ediploma *models.EDiploma) error {
	_, err := r.db.InsertOne(ctx, ediploma)
	return err
}
func (r *eDiplomaRepository) GetByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]*models.EDiploma, error) {
	filter := bson.M{"faculty_id": facultyID}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ediplomas []*models.EDiploma
	for cursor.Next(ctx) {
		var ed models.EDiploma
		if err := cursor.Decode(&ed); err != nil {
			return nil, err
		}
		ediplomas = append(ediplomas, &ed)
	}

	return ediplomas, nil
}
