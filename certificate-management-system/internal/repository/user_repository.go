package repository

import (
	"context"
	"log"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error
	SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	FindByStudentCode(ctx context.Context, studentID string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByStudentCodeAndUniversityID(ctx context.Context, studentCode string, universityID primitive.ObjectID) (bool, error)
	FindUsersByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]*models.User, error)
	FindByStudentCodeAndUniversityID(ctx context.Context, studentCode string, universityID primitive.ObjectID) (*models.User, error)
}
type userRepository struct {
	col        *mongo.Collection
	facultyCol *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	col := db.Collection("users")
	facultyCol := db.Collection("faculties")

	repo := &userRepository{
		col:        col,
		facultyCol: facultyCol,
	}
	return repo
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (r *userRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) SearchUsers(ctx context.Context, params models.SearchUserParams) ([]*models.User, int64, error) {
	log.Printf(" SearchUserParams received: %+v\n", params)

	filter := bson.M{}

	if !params.UniversityID.IsZero() {
		filter["university_id"] = params.UniversityID
	}

	if params.StudentCode != "" {
		filter["student_code"] = bson.M{"$regex": params.StudentCode, "$options": "i"}
	}
	if params.FullName != "" {
		filter["full_name"] = bson.M{"$regex": params.FullName, "$options": "i"}
	}
	if params.Email != "" {
		filter["email"] = bson.M{"$regex": params.Email, "$options": "i"}
	}
	if params.Status != 0 {
		filter["status"] = params.Status
	}

	if params.Faculty != "" {
		log.Println("Filtering by faculty_code:", params.Faculty)

		var faculty struct {
			ID primitive.ObjectID `bson:"_id"`
		}

		facultyFilter := bson.M{
			"faculty_code": bson.M{"$regex": params.Faculty, "$options": "i"},
		}

		if !params.UniversityID.IsZero() {
			facultyFilter["university_id"] = params.UniversityID
		}

		err := r.facultyCol.FindOne(ctx, facultyFilter).Decode(&faculty)
		if err != nil {
			log.Println(" Faculty not found with code:", params.Faculty, "err:", err)

			return []*models.User{}, 0, nil
		} else {
			filter["faculty_id"] = faculty.ID
			log.Println(" Found faculty_id:", faculty.ID.Hex())
		}
	}

	if params.Course != "" {
		filter["course"] = bson.M{"$regex": params.Course, "$options": "i"}
	}

	skip := int64((params.Page - 1) * params.PageSize)
	limit := int64(params.PageSize)

	log.Printf(" Final MongoDB filter: %+v\n", filter)

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

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.col.InsertOne(ctx, user)
	return err
}
func (r *userRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindByStudentCode(ctx context.Context, studentID string) (*models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"student_code": studentID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *userRepository) ExistsByStudentCodeAndUniversityID(ctx context.Context, studentCode string, universityID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"student_code":  studentCode,
		"university_id": universityID,
	}
	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *userRepository) FindByStudentCodeAndUniversityID(ctx context.Context, studentCode string, universityID primitive.ObjectID) (*models.User, error) {
	filter := bson.M{
		"student_code":  studentCode,
		"university_id": universityID,
	}
	var user models.User
	err := r.col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindUsersByFacultyID(ctx context.Context, facultyID primitive.ObjectID) ([]*models.User, error) {
	filter := bson.M{"faculty_id": facultyID}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
