package database

import (
	"bytes"
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
		log.Printf("Created bucket: %s\n", bucketName)
	}

	return &MinioClient{
		Client: client,
		Bucket: bucketName,
	}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, objectName string, fileData []byte, contentType string) error {
	reader := bytes.NewReader(fileData)
	_, err := m.Client.PutObject(ctx, m.Bucket, objectName, reader, int64(len(fileData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}
}

func NewMinioClientFromEnv() (*MinioClient, error) {
	LoadEnv()

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	useSSLStr := os.Getenv("MINIO_USE_SSL")

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		useSSL = false
	}

	return NewMinioClient(endpoint, accessKey, secretKey, bucket, useSSL)
}

func (m *MinioClient) GetFileURL(objectName string) string {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	useSSL := os.Getenv("MINIO_USE_SSL")
	scheme := "http"
	if useSSL == "true" {
		scheme = "https"
	}
	return scheme + "://" + endpoint + "/" + m.Bucket + "/" + objectName
}
