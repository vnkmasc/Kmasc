package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vnkmasc/Kmasc/app/backend/internal/handlers"
	"github.com/vnkmasc/Kmasc/app/backend/internal/middleware"
)

func SetupRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	certificateHandler *handlers.CertificateHandler,
	universityHandler *handlers.UniversityHandler,
	facultyHandler *handlers.FacultyHandler,
	fileHandler *handlers.FileHandler,
	verificationHandler *handlers.VerificationHandler,
	rewardDisciplineHandler *handlers.RewardDisciplineHandler,
	blockchainHandler *handlers.BlockchainHandler,
	majorHandler *handlers.MajorHandler,
	templateHandler *handlers.TemplateHandler,
	ediplomaHandler *handlers.EDiplomaHandler,
	templateSampleHandler *handlers.TemplateSampleHandler,

) *gin.Engine {
	r := gin.Default()

	// CORS setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api/v1")

	// ===== Auth routes =====
	authPublic := api.Group("/auth")
	authPublic.POST("/login", authHandler.Login)
	authPublic.POST("/request-otp", authHandler.RequestOTP)
	authPublic.POST("/verify-otp", authHandler.VerifyOTP)
	authPublic.POST("/register", authHandler.Register)
	authPublic.POST("/verification", verificationHandler.VerifyCode)

	authPrivate := api.Group("/auth")
	authPrivate.Use(middleware.JWTAuthMiddleware())
	authPrivate.GET("/accounts", authHandler.GetAllAccounts)
	authPrivate.DELETE("/accounts", authHandler.DeleteAccount)
	authPrivate.POST("/change-password", authHandler.ChangePassword)
	authPrivate.GET("/university-admin-info", middleware.JWTAuthMiddleware(), authHandler.GetUniversityAdmins)
	authPrivate.GET("/students-info", middleware.JWTAuthMiddleware(), authHandler.GetStudentAccounts)

	// ===== User routes =====
	userGroup := api.Group("/users")
	userGroup.Use(middleware.JWTAuthMiddleware())
	userGroup.POST("/import-excel", userHandler.ImportUsersFromExcel)
	userGroup.GET("", userHandler.GetAllUsers)
	userGroup.POST("", userHandler.CreateUser)
	userGroup.GET("/:id", userHandler.GetUserByID)
	userGroup.PUT("/:id", userHandler.UpdateUser)
	userGroup.GET("/search", userHandler.SearchUsers)
	userGroup.GET("/me", userHandler.GetMyProfile)
	userGroup.DELETE("/:id", userHandler.DeleteUser)
	userGroup.GET("/faculty/:faculty_code", userHandler.GetUsersByFacultyCode)

	// ===== Certificate routes =====
	certificateGroup := api.Group("/certificates")
	certificateGroup.Use(middleware.JWTAuthMiddleware())
	certificateGroup.GET("", certificateHandler.GetAllCertificates)
	certificateGroup.POST("", certificateHandler.CreateCertificate)
	certificateGroup.GET("/:id", certificateHandler.GetCertificateByID)
	certificateGroup.POST("/upload-pdf", certificateHandler.UploadCertificateFile)
	certificateGroup.GET("/file/:id", certificateHandler.GetCertificateFile)
	certificateGroup.GET("/student/:id", certificateHandler.GetCertificatesByStudentID)
	certificateGroup.GET("/search", certificateHandler.SearchCertificates)
	certificateGroup.GET("/my-certificate", certificateHandler.GetMyCertificates)
	certificateGroup.DELETE("/:id", certificateHandler.DeleteCertificate)
	certificateGroup.GET("/simple", certificateHandler.GetMyCertificateNames)
	certificateGroup.POST("/import-excel", certificateHandler.ImportCertificatesFromExcel)

	// ===== University routes =====
	universityGroup := api.Group("/universities")
	universityGroup.POST("", universityHandler.CreateUniversity)
	universityGroup.POST("/approve-or-reject", universityHandler.ApproveOrRejectUniversity)
	universityGroup.GET("", universityHandler.GetAllUniversities)
	universityGroup.GET("/status", universityHandler.GetUniversities)

	//Faculty
	facultyGroup := api.Group("/faculties")
	facultyGroup.Use(middleware.JWTAuthMiddleware())
	facultyGroup.POST("", facultyHandler.CreateFaculty)
	facultyGroup.GET("", facultyHandler.GetAllFaculties)
	facultyGroup.PUT("/:id", facultyHandler.UpdateFaculty)
	facultyGroup.DELETE("/:id", facultyHandler.DeleteFaculty)
	facultyGroup.GET("/:id", facultyHandler.GetFacultyByID)
	facultyGroup.GET("/university/:university_id", facultyHandler.GetFacultiesByUniversity)

	//temp
	api.POST("/upload", fileHandler.UploadFile)

	//verification
	auth := api.Group("/verification").Use(middleware.JWTAuthMiddleware())
	auth.POST("/create", verificationHandler.CreateVerificationCode)
	auth.GET("/my-codes", verificationHandler.GetMyCodes)

	// Reward/Discipline routes
	rdGroup := api.Group("/reward-disciplines")
	rdGroup.Use(middleware.JWTAuthMiddleware())
	rdGroup.POST("", rewardDisciplineHandler.CreateRewardDiscipline)
	rdGroup.GET("", rewardDisciplineHandler.GetAllRewardDisciplines)
	rdGroup.GET("/:id", rewardDisciplineHandler.GetRewardDisciplineByID)
	rdGroup.PUT("/:id", rewardDisciplineHandler.UpdateRewardDiscipline)
	rdGroup.DELETE("/:id", rewardDisciplineHandler.DeleteRewardDiscipline)
	rdGroup.GET("/search", rewardDisciplineHandler.SearchRewardDisciplines)
	rdGroup.GET("/my-reward-disciplines", rewardDisciplineHandler.GetMyRewardDisciplines)
	rdGroup.POST("/import-excel", rewardDisciplineHandler.ImportRewardDisciplinesFromExcel)

	//blockchain
	blockchainGroup := api.Group("/blockchain")
	blockchainGroup.POST("/push-chain/:id", blockchainHandler.PushCertificateToChain)
	blockchainGroup.GET("/certificate-on-chain/:id", blockchainHandler.GetCertificateByID)
	blockchainGroup.GET("/verify/:id", blockchainHandler.VerifyCertificateIntegrity)
	blockchainGroup.GET("/verify-file/:id", blockchainHandler.VerifyCertificateFile)
	blockchainGroup.POST("/push-ediploma", blockchainHandler.PushEDiplomasToBlockchain)

	// ===== Major routes =====
	majorGroup := api.Group("/majors")
	majorGroup.Use(middleware.JWTAuthMiddleware())
	majorGroup.POST("", majorHandler.CreateMajor)
	majorGroup.GET("/faculty/:faculty_id", majorHandler.GetMajorsByFaculty)
	majorGroup.DELETE("/:id", majorHandler.DeleteMajor)

	templateGroup := api.Group("/templates")
	templateGroup.Use(middleware.JWTAuthMiddleware())
	templateGroup.POST("", templateHandler.CreateTemplate)
	templateGroup.GET("/faculty/:faculty_id", templateHandler.GetTemplatesByFaculty)
	templateGroup.GET("/university/:university_id/faculty/:faculty_id", templateHandler.GetTemplatesByFacultyAndUniversity)
	templateGroup.POST("/sign/faculty/:faculty_id", templateHandler.SignTemplatesByFaculty)
	templateGroup.POST("/sign/university", templateHandler.SignAllPendingTemplatesOfUniversity)
	templateGroup.POST("/sign/minedu/:university_id", templateHandler.SignTemplatesByMinEdu)
	templateGroup.POST("/verify/faculty/:faculty_id", templateHandler.VerifyTemplatesByFaculty)
	// templateGroup.GET("/:id/file", templateHandler.GetTemplateFile)
	// templateGroup.PUT("/:id", templateHandler.UpdateTemplate)
	templateGroup.GET("/:id", templateHandler.GetTemplateByID)
	templateGroup.POST("/:template_id/sign", templateHandler.SignTemplateByID)

	templateSampleGroup := api.Group("/template-samples")
	templateSampleGroup.Use(middleware.JWTAuthMiddleware())
	templateSampleGroup.POST("", templateSampleHandler.CreateTemplateSample)
	templateSampleGroup.GET("/:id", templateSampleHandler.GetTemplateSampleByID)
	templateSampleGroup.PUT("/:id", templateSampleHandler.UpdateTemplateSample)
	templateSampleGroup.GET("", templateSampleHandler.GetAllTemplateSamples)
	templateSampleGroup.GET("/view/:id", templateSampleHandler.GetTemplateSampleView)

	ediplomaGroup := api.Group("/ediplomas")
	ediplomaGroup.Use(middleware.JWTAuthMiddleware())
	ediplomaGroup.POST("/generate", ediplomaHandler.GenerateEDiploma)
	ediplomaGroup.POST("/generate-bulk", ediplomaHandler.GenerateBulkEDiplomas)
	ediplomaGroup.POST("/generate-bulk-zip", ediplomaHandler.GenerateBulkEDiplomasZip)
	ediplomaGroup.POST("/upload-zip", ediplomaHandler.UploadEDiplomasZip)
	ediplomaGroup.GET("/search", ediplomaHandler.SearchEDiplomas)
	ediplomaGroup.GET("/file/:id", ediplomaHandler.ViewEDiploma)
	return r
}
