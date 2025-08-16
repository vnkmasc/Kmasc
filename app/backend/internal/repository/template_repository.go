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

type TemplateSampleRepo struct {
	collection *mongo.Collection
}

func NewTemplateSampleRepo(db *mongo.Database) *TemplateSampleRepo {
	return &TemplateSampleRepo{
		collection: db.Collection("template_samples"),
	}
}

func (r *TemplateSampleRepo) Create(ctx context.Context, sample *models.TemplateSample) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, sample)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to convert InsertedID to ObjectID")
	}
	return id, nil
}

func (r *TemplateSampleRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*models.TemplateSample, error) {
	var sample models.TemplateSample
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sample)
	if err != nil {
		return nil, err
	}
	return &sample, nil
}
func (r *TemplateSampleRepo) Update(ctx context.Context, sample *models.TemplateSample) error {
	if sample.ID.IsZero() {
		return errors.New("invalid template sample ID")
	}

	sample.UpdatedAt = time.Now() // cập nhật thời gian
	update := bson.M{
		"$set": bson.M{
			"name":         sample.Name,
			"html_content": sample.HTMLContent,
			"updated_at":   sample.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateByID(ctx, sample.ID, update)
	return err
}

func (r *TemplateSampleRepo) GetAllVisible(ctx context.Context, universityID primitive.ObjectID) ([]*models.TemplateSample, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"university_id": primitive.NilObjectID}, // mẫu global (cố định)
			{"university_id": universityID},          // mẫu riêng trường
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var samples []*models.TemplateSample
	for cursor.Next(ctx) {
		var s models.TemplateSample
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		samples = append(samples, &s)
	}

	return samples, nil
}
func (r *TemplateSampleRepo) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
