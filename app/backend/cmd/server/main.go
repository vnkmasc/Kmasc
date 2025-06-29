package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vnkmasc/Kmasc/app/backend/internal/handlers"
	"github.com/vnkmasc/Kmasc/app/backend/internal/repository"
	"github.com/vnkmasc/Kmasc/app/backend/internal/service"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/blockchain"
	"github.com/vnkmasc/Kmasc/app/backend/pkg/database"
	"github.com/vnkmasc/Kmasc/app/backend/routes"
	"github.com/vnkmasc/Kmasc/app/backend/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Không tìm thấy file .env, đang dùng biến môi trường hệ thống")
	}

	if err := database.ConnectMongo(); err != nil {
		log.Fatalf("Lỗi khi kết nối MongoDB: %v", err)
	}
	db := database.DB
	fabricCfg := blockchain.NewFabricConfigFromEnv()

	InitValidator()
	seedAdminAccount(db)

	emailSender := utils.NewSMTPSender(
		os.Getenv("EMAIL_FROM"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
		os.Getenv("EMAIL_PORT"),
	)

	useSSL := false
	if strings.ToLower(os.Getenv("MINIO_USE_SSL")) == "true" {
		useSSL = true
	}

	minioClient, err := database.NewMinioClient(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_BUCKET"),
		useSSL,
	)
	if err != nil {
		log.Fatalf("Không thể khởi tạo MinIO client: %v", err)
	}
	fabricClient, err := blockchain.NewFabricClient(fabricCfg)
	if err != nil {
		log.Fatalf("khởi tạo FabricClient thất bại: %v", err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	universityRepo := repository.NewUniversityRepository(db)
	certificateRepo := repository.NewCertificateRepository(db)
	facultyRepo := repository.NewFacultyRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)
	rewardDisciplineRepo := repository.NewRewardDisciplineRepository(db)

	// Services
	userService := service.NewUserService(userRepo, universityRepo, facultyRepo)
	authService := service.NewAuthService(authRepo, userRepo, emailSender)
	universityService := service.NewUniversityService(universityRepo, authRepo, emailSender)
	certificateService := service.NewCertificateService(certificateRepo, userRepo, facultyRepo, universityRepo, minioClient)
	facultyService := service.NewFacultyService(universityRepo, facultyRepo)
	verificationService := service.NewVerificationService(verificationRepo, certificateService)
	rewardDisciplineService := service.NewRewardDisciplineService(rewardDisciplineRepo, userRepo)
	blockchainSvc := service.NewBlockchainService(
		certificateRepo, userRepo, facultyRepo, universityRepo, fabricClient,
	)

	// Handlers
	facultyHandler := handlers.NewFacultyHandler(facultyService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, universityService, userService, facultyService)
	universityHandler := handlers.NewUniversityHandler(universityService)
	certificateHandler := handlers.NewCertificateHandler(certificateService, universityService, facultyService, userService, minioClient)
	verificationHandler := handlers.NewVerificationHandler(
		verificationService,
		userService,
		certificateService,
		minioClient,
	)
	rewardDisciplineHandler := handlers.NewRewardDisciplineHandler(rewardDisciplineService)

	fileHandler := handlers.NewFileHandler(minioClient)
	blockchainHandler := handlers.NewBlockchainHandler(blockchainSvc)

	// Setup router
	r := routes.SetupRouter(
		userHandler,
		authHandler,
		certificateHandler,
		universityHandler,
		facultyHandler,
		fileHandler,
		verificationHandler,
		rewardDisciplineHandler,
		blockchainHandler,
	)

	// Xử lý tín hiệu dừng
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Đang tắt server...")
		if err := database.CloseMongo(); err != nil {
			log.Printf("Lỗi khi đóng kết nối MongoDB: %v", err)
		}
		os.Exit(0)
	}()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
