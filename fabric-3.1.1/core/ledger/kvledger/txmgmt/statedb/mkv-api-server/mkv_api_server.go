package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)

// Request structures
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type InitializeRequest struct {
	Password string `json:"password"`
}

type StatusResponse struct {
	Status            string `json:"status"`
	Message           string `json:"message"`
	K1Exists          bool   `json:"k1_exists"`
	SaltExists        bool   `json:"salt_exists"`
	EncryptedK1Exists bool   `json:"encrypted_k1_exists"`
}

type TestPasswordRequest struct {
	Password string `json:"password"`
}

type TestPasswordResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// Global variables
var (
	serverPort = "9876"
	apiKey     = "mkv_api_secret_2025" // In production, use environment variable
)

func main() {
	// Set port from environment variable if available
	if port := os.Getenv("MKV_API_PORT"); port != "" {
		serverPort = port
	}

	// Set API key from environment variable if available
	if key := os.Getenv("MKV_API_KEY"); key != "" {
		apiKey = key
	}

	// Setup routes
	http.HandleFunc("/api/v1/status", authMiddleware(statusHandler))
	http.HandleFunc("/api/v1/initialize", authMiddleware(initializeHandler))
	http.HandleFunc("/api/v1/change-password", authMiddleware(changePasswordHandler))
	http.HandleFunc("/api/v1/test-password", authMiddleware(testPasswordHandler))
	http.HandleFunc("/api/v1/health", healthHandler)

	// Start server
	serverAddr := ":" + serverPort
	log.Printf("🚀 MKV API Server starting on port %s", serverPort)
	log.Printf("📋 Available endpoints:")
	log.Printf("   GET  /api/v1/status")
	log.Printf("   POST /api/v1/initialize")
	log.Printf("   POST /api/v1/change-password")
	log.Printf("   POST /api/v1/test-password")
	log.Printf("   GET  /api/v1/health")

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("🛑 Shutting down MKV API Server...")
		os.Exit(0)
	}()

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Middleware for API key authentication
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check
		if r.URL.Path == "/api/v1/health" {
			next(w, r)
			return
		}

		// Check API key in header
		authHeader := r.Header.Get("X-API-Key")
		if authHeader != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// Status endpoint
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := StatusResponse{
		Status: "success",
	}

	// Check if files exist
	_, err := os.Stat("k1.key")
	response.K1Exists = err == nil

	_, err = os.Stat("k0_salt.key")
	response.SaltExists = err == nil

	_, err = os.Stat("encrypted_k1.key")
	response.EncryptedK1Exists = err == nil

	if response.K1Exists && response.SaltExists && response.EncryptedK1Exists {
		response.Message = "System is initialized and ready"
	} else {
		response.Message = "System needs initialization"
	}

	json.NewEncoder(w).Encode(response)
}

// Initialize endpoint
func initializeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req InitializeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Initialize the system
	err := mkv.InitializeKeyManagement(req.Password)
	if err != nil {
		response := map[string]string{
			"status":  "error",
			"message": fmt.Sprintf("Failed to initialize: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "System initialized successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// Change password endpoint
func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Both old and new passwords are required", http.StatusBadRequest)
		return
	}

	// Change password
	err := mkv.ChangePassword(req.OldPassword, req.NewPassword)
	if err != nil {
		response := map[string]string{
			"status":  "error",
			"message": fmt.Sprintf("Failed to change password: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Password changed successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// Test password endpoint
func testPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req TestPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Test password
	_, err := mkv.GetCurrentK1(req.Password)
	response := TestPasswordResponse{
		Valid: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}
