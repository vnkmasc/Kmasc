package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "103.124.93.139:9000" // bỏ https:// và dấu /
	accessKey := "admin"
	secretKey := "12345678"
	bucketName := "certificates"
	objectName := "ediplomas/68975e260b5e20fe197fc873.pdf"

	// Kết nối MinIO
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // dùng HTTPS
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Tạo presigned URL có hạn 1 giờ
	presignedURL, err := minioClient.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		time.Hour,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Test download link (valid for 1h):")
	fmt.Println(presignedURL.String())
}
