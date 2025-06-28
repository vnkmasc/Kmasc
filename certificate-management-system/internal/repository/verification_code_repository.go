package repository

import (
	"context"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VerificationRepository interface {
	Save(ctx context.Context, code *models.VerificationCode) error
	GetByCode(ctx context.Context, code string) (*models.VerificationCode, error)
	MarkViewed(ctx context.Context, id primitive.ObjectID, viewType string) error
	GetByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int64) ([]models.VerificationCode, int64, error)
}

type verificationRepository struct {
	collection *mongo.Collection
}

func NewVerificationRepository(db *mongo.Database) VerificationRepository {
	return &verificationRepository{
		collection: db.Collection("verification_codes"),
	}
}

func (r *verificationRepository) Save(ctx context.Context, code *models.VerificationCode) error {
	_, err := r.collection.InsertOne(ctx, code)
	return err
}

func (r *verificationRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int64) ([]models.VerificationCode, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip((page - 1) * pageSize).
		SetLimit(pageSize)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var codes []models.VerificationCode
	if err := cursor.All(ctx, &codes); err != nil {
		return nil, 0, err
	}

	return codes, total, nil
}

func (r *verificationRepository) GetByCode(ctx context.Context, code string) (*models.VerificationCode, error) {
	filter := bson.M{"code": code}
	var result models.VerificationCode
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *verificationRepository) MarkViewed(ctx context.Context, id primitive.ObjectID, viewType string) error {
	update := bson.M{}
	switch viewType {
	case "score":
		update = bson.M{"$set": bson.M{"viewed_score": true}}
	case "data":
		update = bson.M{"$set": bson.M{"viewed_data": true}}
	case "file":
		update = bson.M{"$set": bson.M{"viewed_file": true}}
	}
	_, err := r.collection.UpdateByID(ctx, id, update)
	return err
}
