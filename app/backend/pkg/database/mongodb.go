package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Không thể load file .env hoặc không tồn tại, dùng biến môi trường hệ thống")
	}
}

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("biến môi trường MONGODB_URI chưa được thiết lập")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "vbcc_data"
	}

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("kết nối MongoDB thất bại: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping MongoDB thất bại: %w", err)
	}

	Client = client
	DB = client.Database(dbName)

	log.Println("MongoDB kết nối thành công với DB:", dbName)
	return nil
}

func CloseMongo() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("đóng kết nối MongoDB lỗi: %w", err)
	}
	log.Println("MongoDB đã đóng kết nối thành công")
	return nil
}
