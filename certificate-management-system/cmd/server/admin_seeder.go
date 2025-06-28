package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tuyenngduc/certificate-management-system/internal/models"
	"github.com/tuyenngduc/certificate-management-system/utils"
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
