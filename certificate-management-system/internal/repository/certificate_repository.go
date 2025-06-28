package repository

import (
	"context"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CertificateRepository interface {
	GetAllCertificates(ctx context.Context) ([]*models.Certificate, error)
	FindOne(ctx context.Context, filter interface{}) (*models.Certificate, error)
	GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error)
	FindCertificateByStudentCodeAndName(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error)
	ExistsByRegNo(ctx context.Context, universityID primitive.ObjectID, regNo string, isDegree bool) (bool, error)
	ExistsBySerial(ctx context.Context, universityID primitive.ObjectID, serial string, isDegree bool) (bool, error)
	DeleteCertificate(ctx context.Context, id primitive.ObjectID) error
	DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error
	CreateCertificate(ctx context.Context, cert *models.Certificate) error
	UpdateCertificatePath(ctx context.Context, certificateID primitive.ObjectID, path string) error
	FindBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error)
	FindLatestCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Certificate, error)
	FindCertificate(ctx context.Context, filter bson.M, page, pageSize int) ([]*models.Certificate, int64, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Certificate, error)
	ExistsCertificateByStudentCodeAndName(ctx context.Context, studentCode string, universityID primitive.ObjectID, name string) (bool, error)

	ExistsDegreeByStudentCodeAndType(ctx context.Context, studentCode string, universityID primitive.ObjectID, certType string) (bool, error)
	FindBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error)
}
type certificateRepository struct {
	col *mongo.Collection
}

func NewCertificateRepository(db *mongo.Database) CertificateRepository {
	col := db.Collection("certificates")
	return &certificateRepository{col: col}
}

func (r *certificateRepository) CreateCertificate(ctx context.Context, cert *models.Certificate) error {
	_, err := r.col.InsertOne(ctx, cert)
	return err
}

func (r *certificateRepository) GetAllCertificates(ctx context.Context) ([]*models.Certificate, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certs []*models.Certificate
	if err := cursor.All(ctx, &certs); err != nil {
		return nil, err
	}
	return certs, nil
}
func (r *certificateRepository) GetCertificateByID(ctx context.Context, id primitive.ObjectID) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) DeleteCertificate(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *certificateRepository) UpdateCertificatePath(ctx context.Context, certificateID primitive.ObjectID, path string) error {
	filter := bson.M{"_id": certificateID}
	update := bson.M{
		"$set": bson.M{
			"path":       path,
			"updated_at": time.Now(),
		},
	}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
func (r *certificateRepository) FindBySerialNumber(ctx context.Context, serial string) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.col.FindOne(ctx, bson.M{"serial_number": serial}).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
func (r *certificateRepository) FindLatestCertificateByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Certificate, error) {
	filter := bson.M{"user_id": userID}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}) // sắp xếp giảm dần theo created_at để lấy mới nhất
	var certificate models.Certificate
	err := r.col.FindOne(ctx, filter, opts).Decode(&certificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &certificate, nil
}
func (r *certificateRepository) FindOne(ctx context.Context, filter interface{}) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.col.FindOne(ctx, filter).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) FindCertificate(ctx context.Context, filter bson.M, page, pageSize int) ([]*models.Certificate, int64, error) {
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.col.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var certs []*models.Certificate
	if err := cursor.All(ctx, &certs); err != nil {
		return nil, 0, err
	}
	return certs, total, nil
}
func (r *certificateRepository) UpdateVerificationCode(ctx context.Context, id primitive.ObjectID, code string, expired time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"verification_code": code,
			"code_expired_at":   expired,
		},
	}
	_, err := r.col.UpdateByID(ctx, id, update)
	return err
}

func (r *certificateRepository) FindBySerialAndUniversity(ctx context.Context, serial string, universityID primitive.ObjectID) (*models.Certificate, error) {
	filter := bson.M{
		"serial_number": serial,
		"university_id": universityID,
	}
	var cert models.Certificate
	err := r.col.FindOne(ctx, filter).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
func (r *certificateRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Certificate, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certs []*models.Certificate
	if err := cursor.All(ctx, &certs); err != nil {
		return nil, err
	}
	return certs, nil
}
func (r *certificateRepository) DeleteCertificateByID(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	res, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *certificateRepository) ExistsDegreeByStudentCodeAndType(ctx context.Context, studentCode string, universityID primitive.ObjectID, certType string) (bool, error) {
	filter := bson.M{
		"student_code":        studentCode,
		"university_id":       universityID,
		"certificate_type":    certType,
		"serial_number":       bson.M{"$ne": ""},
		"registration_number": bson.M{"$ne": ""},
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *certificateRepository) ExistsCertificateByStudentCodeAndName(ctx context.Context, studentCode string, universityID primitive.ObjectID, name string) (bool, error) {
	filter := bson.M{
		"student_code":  studentCode,
		"university_id": universityID,
		"name":          name,
		"is_degree":     false,
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *certificateRepository) FindCertificateByStudentCodeAndName(ctx context.Context, studentCode, name string, universityID primitive.ObjectID) (*models.Certificate, error) {
	filter := bson.M{
		"student_code":  studentCode,
		"name":          name,
		"university_id": universityID,
		"is_degree":     false,
	}
	var cert models.Certificate
	err := r.col.FindOne(ctx, filter).Decode(&cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) ExistsBySerial(ctx context.Context, universityID primitive.ObjectID, serial string, isDegree bool) (bool, error) {
	filter := bson.M{
		"university_id": universityID,
		"serial_number": serial,
		"is_degree":     isDegree,
	}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *certificateRepository) ExistsByRegNo(ctx context.Context, universityID primitive.ObjectID, regNo string, isDegree bool) (bool, error) {
	filter := bson.M{
		"university_id": universityID,
		"reg_no":        regNo,
		"is_degree":     isDegree,
	}
	count, err := r.col.CountDocuments(ctx, filter)
	return count > 0, err
}
