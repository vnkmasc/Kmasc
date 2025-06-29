package repository

import (
	"context"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RewardDisciplineRepository interface {
	Create(ctx context.Context, rd *models.RewardDiscipline) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.RewardDiscipline, error)
	GetAll(ctx context.Context) ([]*models.RewardDiscipline, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Search(ctx context.Context, params models.SearchRewardDisciplineParams) ([]*models.RewardDiscipline, int64, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.RewardDiscipline, error)
	ExistsByDecisionNumber(ctx context.Context, decisionNumber string) (bool, error)
	ExistsByDecisionNumberExcludeID(ctx context.Context, decisionNumber string, excludeID primitive.ObjectID) (bool, error)
}

type rewardDisciplineRepository struct {
	col *mongo.Collection
}

func NewRewardDisciplineRepository(db *mongo.Database) RewardDisciplineRepository {
	return &rewardDisciplineRepository{
		col: db.Collection("reward_disciplines"),
	}
}

func (r *rewardDisciplineRepository) Create(ctx context.Context, rd *models.RewardDiscipline) error {
	_, err := r.col.InsertOne(ctx, rd)
	return err
}

func (r *rewardDisciplineRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.RewardDiscipline, error) {
	var rd models.RewardDiscipline
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&rd)
	if err != nil {
		return nil, err
	}
	return &rd, nil
}

func (r *rewardDisciplineRepository) GetAll(ctx context.Context) ([]*models.RewardDiscipline, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rds []*models.RewardDiscipline
	if err := cursor.All(ctx, &rds); err != nil {
		return nil, err
	}
	return rds, nil
}

func (r *rewardDisciplineRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *rewardDisciplineRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *rewardDisciplineRepository) Search(ctx context.Context, params models.SearchRewardDisciplineParams) ([]*models.RewardDiscipline, int64, error) {
	filter := bson.M{}

	if params.Name != "" {
		filter["name"] = bson.M{"$regex": params.Name, "$options": "i"}
	}
	if params.DecisionNumber != "" {
		filter["decision_number"] = bson.M{"$regex": params.DecisionNumber, "$options": "i"}
	}
	if params.IsDiscipline != nil {
		filter["is_discipline"] = *params.IsDiscipline
	}

	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var rds []*models.RewardDiscipline
	if err := cursor.All(ctx, &rds); err != nil {
		return nil, 0, err
	}

	return rds, total, nil
}

func (r *rewardDisciplineRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.RewardDiscipline, error) {
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rds []*models.RewardDiscipline
	if err := cursor.All(ctx, &rds); err != nil {
		return nil, err
	}
	return rds, nil
}

func (r *rewardDisciplineRepository) ExistsByDecisionNumber(ctx context.Context, decisionNumber string) (bool, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"decision_number": decisionNumber})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *rewardDisciplineRepository) ExistsByDecisionNumberExcludeID(ctx context.Context, decisionNumber string, excludeID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"decision_number": decisionNumber,
		"_id":             bson.M{"$ne": excludeID},
	}
	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
