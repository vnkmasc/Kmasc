package mkv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestKeyManagerSingleton(t *testing.T) {
	// Clean up trước khi test
	cleanupTestFiles(t)

	// Lấy instance đầu tiên
	km1 := GetKeyManager()
	if km1 == nil {
		t.Fatal("KeyManager instance should not be nil")
	}

	// Lấy instance thứ hai
	km2 := GetKeyManager()
	if km2 == nil {
		t.Fatal("Second KeyManager instance should not be nil")
	}

	// Kiểm tra xem có phải cùng một instance không
	if km1 != km2 {
		t.Fatal("KeyManager should return the same instance")
	}

	// Kiểm tra trạng thái
	status := km1.GetStatus()
	if !status["initialized"].(bool) {
		t.Error("KeyManager should be initialized")
	}

	if !status["has_key"].(bool) {
		t.Error("KeyManager should have encryption key")
	}

	if status["key_length"].(int) != 32 {
		t.Error("Encryption key should be 32 bytes")
	}

	cleanupTestFiles(t)
}



func TestKeyManagerEncryptionDecryption(t *testing.T) {
	// Clean up trước khi test
	cleanupTestFiles(t)

	// Lấy KeyManager instance
	_ = GetKeyManager() // Sử dụng _ để bỏ qua biến không sử dụng

	// Đợi khởi tạo hoàn tất
	time.Sleep(100 * time.Millisecond)

	// Test data
	testData := []byte("Hello, MKV encryption test!")

	// Mã hóa
	encrypted := EncryptValueMKV(testData)
	if encrypted == nil {
		t.Fatal("Encryption should succeed")
	}

	// Giải mã
	decrypted := DecryptValueMKV(encrypted)
	if decrypted == nil {
		t.Fatal("Decryption should succeed")
	}

	// So sánh kết quả
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data mismatch. Expected: %s, Got: %s", string(testData), string(decrypted))
	}

	cleanupTestFiles(t)
}

func TestKeyManagerChangePassword(t *testing.T) {
	// Clean up trước khi test
	cleanupTestFiles(t)

	// Lấy KeyManager instance
	km := GetKeyManager()

	// Đợi khởi tạo hoàn tất
	time.Sleep(100 * time.Millisecond)

	// Test data
	testData := []byte("Test data for password change")

	// Mã hóa với password cũ
	encrypted1 := EncryptValueMKV(testData)
	if encrypted1 == nil {
		t.Fatal("Initial encryption should succeed")
	}

	// Thay đổi password
	newPassword := "new_kmasc_password"
	if err := km.ChangePassword(newPassword); err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// Mã hóa lại với password mới
	encrypted2 := EncryptValueMKV(testData)
	if encrypted2 == nil {
		t.Fatal("Encryption with new password should succeed")
	}

	// Giải mã với password mới
	decrypted := DecryptValueMKV(encrypted2)
	if decrypted == nil {
		t.Fatal("Decryption with new password should succeed")
	}

	// So sánh kết quả
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data mismatch after password change. Expected: %s, Got: %s", string(testData), string(decrypted))
	}

	cleanupTestFiles(t)
}

func TestKeyManagerRefreshKeys(t *testing.T) {
	// Clean up trước khi test
	cleanupTestFiles(t)

	// Lấy KeyManager instance
	km := GetKeyManager()

	// Đợi khởi tạo hoàn tất
	time.Sleep(100 * time.Millisecond)

	// Test data
	testData := []byte("Test data for key refresh")

	// Mã hóa với key cũ
	encrypted1 := EncryptValueMKV(testData)
	if encrypted1 == nil {
		t.Fatal("Initial encryption should succeed")
	}

	// Làm mới keys
	if err := km.RefreshKeys(); err != nil {
		t.Fatalf("Failed to refresh keys: %v", err)
	}

	// Mã hóa với key mới
	encrypted2 := EncryptValueMKV(testData)
	if encrypted2 == nil {
		t.Fatal("Encryption with new keys should succeed")
	}

	// Giải mã với key mới
	decrypted := DecryptValueMKV(encrypted2)
	if decrypted == nil {
		t.Fatal("Decryption with new keys should succeed")
	}

	// So sánh kết quả
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data mismatch after key refresh. Expected: %s, Got: %s", string(testData), string(decrypted))
	}

	cleanupTestFiles(t)
}

func TestKeyManagerConcurrentAccess(t *testing.T) {
	// Clean up trước khi test
	cleanupTestFiles(t)

	// Lấy KeyManager instance
	km := GetKeyManager()

	// Đợi khởi tạo hoàn tất
	time.Sleep(100 * time.Millisecond)

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Lấy key
			key := km.GetEncryptionKey()
			if len(key) != 32 {
				t.Errorf("Goroutine %d: Expected key length 32, got %d", id, len(key))
			}

			// Test data
			testData := []byte(fmt.Sprintf("Test data from goroutine %d", id))

			// Mã hóa
			encrypted := EncryptValueMKV(testData)
			if encrypted == nil {
				t.Errorf("Goroutine %d: Encryption failed", id)
				return
			}

			// Giải mã
			decrypted := DecryptValueMKV(encrypted)
			if decrypted == nil {
				t.Errorf("Goroutine %d: Decryption failed", id)
				return
			}

			// So sánh
			if string(decrypted) != string(testData) {
				t.Errorf("Goroutine %d: Data mismatch", id)
			}
		}(i)
	}

	// Đợi tất cả goroutines hoàn thành
	for i := 0; i < 10; i++ {
		<-done
	}

	cleanupTestFiles(t)
}

// Helper function để cleanup test files
func cleanupTestFiles(t *testing.T) {
	searchPaths := []string{".", "/tmp/mkv", "/opt/mkv", "/home/chaincode/mkv", "/root/mkv"}

	for _, path := range searchPaths {
		files := []string{"k1.key", "k0.key", "encrypted_k1.key", "k0_salt.key", "password.txt"}

		for _, file := range files {
			filePath := filepath.Join(path, file)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: Failed to remove %s: %v", filePath, err)
			}
		}

		// Xóa thư mục nếu rỗng (trừ current directory)
		if path != "." {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: Failed to remove directory %s: %v", path, err)
			}
		}
	}
}
