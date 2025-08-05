package repository

import (
	"context"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EDiplomaRepository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.EDiploma, error)
	Save(ctx context.Context, ediploma *models.EDiploma) error
}

type eDiplomaRepository struct {
	db *mongo.Collection
}

func NewEDiplomaRepository(db *mongo.Database) EDiplomaRepository {
	return &eDiplomaRepository{
		db: db.Collection("ediplomas"),
	}
}

func (r *eDiplomaRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.EDiploma, error) {
	var diploma models.EDiploma
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&diploma)
	if err != nil {
		return nil, err
	}
	return &diploma, nil
}

func (r *eDiplomaRepository) Save(ctx context.Context, ediploma *models.EDiploma) error {
	_, err := r.db.InsertOne(ctx, ediploma)
	return err
}
