package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/vnkmasc/Kmasc/app/backend/internal/models"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func seedAdminAccount(db *mongo.Database) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Fatal("ADMIN_EMAIL hoặc ADMIN_PASSWORD chưa được cấu hình trong .env")
	}

	collection := db.Collection("accounts")

	// Kiểm tra nếu admin đã tồn tại
	count, err := collection.CountDocuments(context.TODO(), bson.M{"personal_email": adminEmail})
	if err != nil {
		log.Fatalf("Lỗi khi kiểm tra tài khoản admin: %v", err)
	}
	if count > 0 {
		log.Println("Tài khoản admin đã tồn tại, không cần tạo thêm.")
		return
	}

	passwordHash, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("Lỗi khi hash mật khẩu: %v", err)
	}

	admin := models.Account{
		ID:            primitive.NewObjectID(),
		StudentID:     primitive.NilObjectID,
		StudentEmail:  "",
		PersonalEmail: adminEmail,
		PasswordHash:  passwordHash,
		CreatedAt:     time.Now(),
		Role:          "admin",
	}

	_, err = collection.InsertOne(context.TODO(), admin)
	if err != nil {
		log.Fatalf("Lỗi khi tạo tài khoản admin: %v", err)
	}

	log.Println("Tài khoản admin đã được tạo thành công.")
}

func SeedTemplateSamples(ctx context.Context, repo *repository.TemplateSampleRepo) error {
	// Kiểm tra đã tồn tại chưa
	count, err := repo.Count(ctx, bson.M{"university_id": primitive.NilObjectID})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Lấy thư mục hiện tại
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	readHTML := func(fileName string) string {
		path := filepath.Join(cwd, "pkg", "database", fileName)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read template file %s: %v", path, err)
		}
		return string(data)
	}

	samples := []models.TemplateSample{
		{Name: "Mẫu 1", HTMLContent: readHTML("template1.html"), UniversityID: primitive.NilObjectID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Mẫu 2", HTMLContent: readHTML("template2.html"), UniversityID: primitive.NilObjectID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Mẫu 3", HTMLContent: readHTML("template3.html"), UniversityID: primitive.NilObjectID, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, s := range samples {
		if _, err := repo.Create(ctx, &s); err != nil {
			return err
		}
	}

	return nil
}
