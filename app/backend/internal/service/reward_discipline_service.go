package service

import (
	"context"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/common"
	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RewardDisciplineService interface {
	CreateRewardDiscipline(ctx context.Context, req *models.CreateRewardDisciplineRequest) (*models.RewardDisciplineResponse, error)
	GetRewardDisciplineByID(ctx context.Context, id primitive.ObjectID) (*models.RewardDisciplineResponse, error)
	GetAllRewardDisciplines(ctx context.Context) ([]models.RewardDisciplineResponse, error)
	UpdateRewardDiscipline(ctx context.Context, id primitive.ObjectID, req *models.UpdateRewardDisciplineRequest) error
	DeleteRewardDiscipline(ctx context.Context, id primitive.ObjectID) error
	SearchRewardDisciplines(ctx context.Context, params models.SearchRewardDisciplineParams) ([]models.RewardDisciplineResponse, int64, error)
	GetMyRewardDisciplines(ctx context.Context) ([]models.RewardDisciplineResponse, error)
}

type rewardDisciplineService struct {
	rdRepo   repository.RewardDisciplineRepository
	userRepo repository.UserRepository
}

func NewRewardDisciplineService(
	rdRepo repository.RewardDisciplineRepository,
	userRepo repository.UserRepository,
) RewardDisciplineService {
	return &rewardDisciplineService{
		rdRepo:   rdRepo,
		userRepo: userRepo,
	}
}

func (s *rewardDisciplineService) CreateRewardDiscipline(ctx context.Context, req *models.CreateRewardDisciplineRequest) (*models.RewardDisciplineResponse, error) {
	// Check if DecisionNumber already exists
	exists, err := s.rdRepo.ExistsByDecisionNumber(ctx, req.DecisionNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrDecisionNumberExists
	}

	// Check if user exists
	user, err := s.userRepo.FindByStudentCode(ctx, req.StudentCode)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrUserNotExisted
		}
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotExisted
	}

	// Validate discipline level
	if req.IsDiscipline && (req.DisciplineLevel == nil || *req.DisciplineLevel < 1 || *req.DisciplineLevel > 4) {
		return nil, common.NewValidationError("DisciplineLevel", "Mức độ kỷ luật phải từ 1 đến 4 khi IsDiscipline=true")
	}
	if !req.IsDiscipline && req.DisciplineLevel != nil {
		req.DisciplineLevel = nil // Clear discipline level if not a discipline
	}

	rd := &models.RewardDiscipline{
		ID:              primitive.NewObjectID(),
		Name:            req.Name,
		DecisionNumber:  req.DecisionNumber,
		Description:     req.Description,
		UserID:          user.ID,
		IsDiscipline:    req.IsDiscipline,
		DisciplineLevel: req.DisciplineLevel,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.rdRepo.Create(ctx, rd); err != nil {
		return nil, err
	}

	return &models.RewardDisciplineResponse{
		ID:              rd.ID,
		Name:            rd.Name,
		DecisionNumber:  rd.DecisionNumber,
		Description:     rd.Description,
		StudentCode:     user.StudentCode,
		StudentName:     user.FullName,
		IsDiscipline:    rd.IsDiscipline,
		DisciplineLevel: rd.DisciplineLevel,
		CreatedAt:       rd.CreatedAt,
		UpdatedAt:       rd.UpdatedAt,
	}, nil
}

func (s *rewardDisciplineService) GetRewardDisciplineByID(ctx context.Context, id primitive.ObjectID) (*models.RewardDisciplineResponse, error) {
	rd, err := s.rdRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, rd.UserID)
	if err != nil {
		return nil, err
	}

	return &models.RewardDisciplineResponse{
		ID:              rd.ID,
		Name:            rd.Name,
		DecisionNumber:  rd.DecisionNumber,
		Description:     rd.Description,
		StudentCode:     user.StudentCode,
		StudentName:     user.FullName,
		IsDiscipline:    rd.IsDiscipline,
		DisciplineLevel: rd.DisciplineLevel,
		CreatedAt:       rd.CreatedAt,
		UpdatedAt:       rd.UpdatedAt,
	}, nil
}

func (s *rewardDisciplineService) GetAllRewardDisciplines(ctx context.Context) ([]models.RewardDisciplineResponse, error) {
	rds, err := s.rdRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []models.RewardDisciplineResponse
	for _, rd := range rds {
		user, err := s.userRepo.GetUserByID(ctx, rd.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		responses = append(responses, models.RewardDisciplineResponse{
			ID:              rd.ID,
			Name:            rd.Name,
			DecisionNumber:  rd.DecisionNumber,
			Description:     rd.Description,
			StudentCode:     user.StudentCode,
			StudentName:     user.FullName,
			IsDiscipline:    rd.IsDiscipline,
			DisciplineLevel: rd.DisciplineLevel,
			CreatedAt:       rd.CreatedAt,
			UpdatedAt:       rd.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *rewardDisciplineService) UpdateRewardDiscipline(ctx context.Context, id primitive.ObjectID, req *models.UpdateRewardDisciplineRequest) error {
	// Check if record exists
	existing, err := s.rdRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return common.ErrNotFound
		}
		return err
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.DecisionNumber != nil {
		// Check if DecisionNumber already exists (excluding current record)
		exists, err := s.rdRepo.ExistsByDecisionNumberExcludeID(ctx, *req.DecisionNumber, id)
		if err != nil {
			return err
		}
		if exists {
			return common.ErrDecisionNumberExists
		}
		update["decision_number"] = *req.DecisionNumber
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.StudentCode != nil {
		user, err := s.userRepo.FindByStudentCode(ctx, *req.StudentCode)
		if err != nil || user == nil {
			return common.ErrUserNotExisted
		}
		update["user_id"] = user.ID
	}
	if req.IsDiscipline != nil {
		update["is_discipline"] = *req.IsDiscipline

		// If changing to discipline, validate level
		if *req.IsDiscipline {
			if req.DisciplineLevel == nil || *req.DisciplineLevel < 1 || *req.DisciplineLevel > 4 {
				return common.NewValidationError("DisciplineLevel", "Mức độ kỷ luật phải từ 1 đến 4 khi IsDiscipline=true")
			}
			update["discipline_level"] = *req.DisciplineLevel
		} else {
			// If changing to reward, remove discipline level
			update["discipline_level"] = nil
		}
	} else if req.DisciplineLevel != nil {
		// Only update discipline level if it's currently a discipline
		if existing.IsDiscipline {
			if *req.DisciplineLevel < 1 || *req.DisciplineLevel > 4 {
				return common.NewValidationError("DisciplineLevel", "Mức độ kỷ luật phải từ 1 đến 4")
			}
			update["discipline_level"] = *req.DisciplineLevel
		}
	}

	return s.rdRepo.Update(ctx, id, update)
}

func (s *rewardDisciplineService) DeleteRewardDiscipline(ctx context.Context, id primitive.ObjectID) error {
	return s.rdRepo.Delete(ctx, id)
}

func (s *rewardDisciplineService) SearchRewardDisciplines(ctx context.Context, params models.SearchRewardDisciplineParams) ([]models.RewardDisciplineResponse, int64, error) {
	// Get university context for filtering
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, 0, common.ErrUnauthorized
	}

	// If searching by student code, find the user first
	if params.StudentCode != "" {
		universityID, err := primitive.ObjectIDFromHex(claims.UniversityID)
		if err != nil {
			return nil, 0, common.ErrInvalidToken
		}

		user, err := s.userRepo.FindByStudentCodeAndUniversityID(ctx, params.StudentCode, universityID)
		if err != nil || user == nil {
			return []models.RewardDisciplineResponse{}, 0, nil
		}

		// Get reward disciplines for this user
		rds, err := s.rdRepo.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, 0, err
		}

		var responses []models.RewardDisciplineResponse
		for _, rd := range rds {
			responses = append(responses, models.RewardDisciplineResponse{
				ID:              rd.ID,
				Name:            rd.Name,
				DecisionNumber:  rd.DecisionNumber,
				Description:     rd.Description,
				StudentCode:     user.StudentCode,
				StudentName:     user.FullName,
				IsDiscipline:    rd.IsDiscipline,
				DisciplineLevel: rd.DisciplineLevel,
				CreatedAt:       rd.CreatedAt,
				UpdatedAt:       rd.UpdatedAt,
			})
		}

		return responses, int64(len(responses)), nil
	}

	rds, total, err := s.rdRepo.Search(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.RewardDisciplineResponse
	for _, rd := range rds {
		user, err := s.userRepo.GetUserByID(ctx, rd.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		responses = append(responses, models.RewardDisciplineResponse{
			ID:              rd.ID,
			Name:            rd.Name,
			DecisionNumber:  rd.DecisionNumber,
			Description:     rd.Description,
			StudentCode:     user.StudentCode,
			StudentName:     user.FullName,
			IsDiscipline:    rd.IsDiscipline,
			DisciplineLevel: rd.DisciplineLevel,
			CreatedAt:       rd.CreatedAt,
			UpdatedAt:       rd.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (s *rewardDisciplineService) GetMyRewardDisciplines(ctx context.Context) ([]models.RewardDisciplineResponse, error) {
	// Get user from context
	claimsVal := ctx.Value(utils.ClaimsContextKey)
	claims, ok := claimsVal.(*utils.CustomClaims)
	if !ok || claims == nil {
		return nil, common.ErrUnauthorized
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	// Get user to verify existence
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, common.ErrUserNotExisted
	}

	// Get reward disciplines for this user
	rds, err := s.rdRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []models.RewardDisciplineResponse
	for _, rd := range rds {
		responses = append(responses, models.RewardDisciplineResponse{
			ID:              rd.ID,
			Name:            rd.Name,
			DecisionNumber:  rd.DecisionNumber,
			Description:     rd.Description,
			StudentCode:     user.StudentCode,
			StudentName:     user.FullName,
			IsDiscipline:    rd.IsDiscipline,
			DisciplineLevel: rd.DisciplineLevel,
			CreatedAt:       rd.CreatedAt,
			UpdatedAt:       rd.UpdatedAt,
		})
	}

	return responses, nil
}
