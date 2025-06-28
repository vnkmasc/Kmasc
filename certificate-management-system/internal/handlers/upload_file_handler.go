package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuyenngduc/certificate-management-system/pkg/database"
)

type FileHandler struct {
	MinioClient *database.MinioClient
}

func NewFileHandler(minioClient *database.MinioClient) *FileHandler {
	return &FileHandler{MinioClient: minioClient}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không có file để upload"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể mở file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi đọc file"})
		return
	}

	objectName := file.Filename
	contentType := file.Header.Get("Content-Type")

	err = h.MinioClient.UploadFile(c.Request.Context(), objectName, fileData, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể upload file"})
		return
	}

	fileURL := h.MinioClient.GetFileURL(objectName)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload thành công",
		"file_url": fileURL,
	})
}
