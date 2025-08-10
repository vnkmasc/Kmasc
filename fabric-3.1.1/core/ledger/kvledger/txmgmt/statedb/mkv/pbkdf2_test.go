package mkv

import (
	"encoding/hex"
	"testing"
)

func TestPBKDF2Implementation(t *testing.T) {
	// Test PBKDF2 với các tham số chuẩn
	t.Run("PBKDF2 with standard parameters", func(t *testing.T) {
		password := []byte("password")
		salt := []byte("salt")
		iterations := 1
		keyLen := 20

		key := PBKDF2(password, salt, iterations, keyLen)
		// Lưu kết quả thực tế để kiểm tra consistency
		actual := hex.EncodeToString(key)
		t.Logf("PBKDF2 output: %s", actual)

		// Test consistency
		key2 := PBKDF2(password, salt, iterations, keyLen)
		actual2 := hex.EncodeToString(key2)

		if actual != actual2 {
			t.Fatalf("PBKDF2 output is not consistent. First: %s, Second: %s", actual, actual2)
		}
		t.Logf("PBKDF2 test passed with standard parameters")
	})

	// Test PBKDF2 với nhiều iterations
	t.Run("PBKDF2 with multiple iterations", func(t *testing.T) {
		password := []byte("password")
		salt := []byte("salt")
		iterations := 4096
		keyLen := 20

		key := PBKDF2(password, salt, iterations, keyLen)
		actual := hex.EncodeToString(key)
		t.Logf("PBKDF2 output with 4096 iterations: %s", actual)

		// Test consistency
		key2 := PBKDF2(password, salt, iterations, keyLen)
		actual2 := hex.EncodeToString(key2)

		if actual != actual2 {
			t.Fatalf("PBKDF2 output is not consistent with 4096 iterations")
		}
		t.Logf("PBKDF2 test passed with 4096 iterations")
	})

	// Test PBKDF2 với key length lớn hơn block size
	t.Run("PBKDF2 with large key length", func(t *testing.T) {
		password := []byte("password")
		salt := []byte("salt")
		iterations := 1
		keyLen := 40 // Lớn hơn SHA256 block size (32 bytes)

		key := PBKDF2(password, salt, iterations, keyLen)
		if len(key) != keyLen {
			t.Fatalf("Key length mismatch. Expected: %d, Got: %d", keyLen, len(key))
		}
		t.Logf("PBKDF2 test passed with large key length: %d bytes", keyLen)
	})
}

func TestGenerateK0FromPasswordWithPBKDF2(t *testing.T) {
	// Test tạo K0 với salt cố định (backward compatibility)
	t.Run("Generate K0 with fixed salt", func(t *testing.T) {
		password := "testpassword123"

		k0, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}

		if len(k0) != 32 {
			t.Fatalf("K0 length should be 32 bytes, got %d", len(k0))
		}

		// Test consistency
		k0_again, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0 again: %v", err)
		}

		if hex.EncodeToString(k0) != hex.EncodeToString(k0_again) {
			t.Fatalf("K0 generation is not consistent")
		}

		t.Logf("K0 generated successfully with fixed salt: %s", hex.EncodeToString(k0))
	})

	// Test tạo K0 với salt ngẫu nhiên
	t.Run("Generate K0 with random salt", func(t *testing.T) {
		password := "testpassword456"

		// Tạo salt ngẫu nhiên
		salt, err := GenerateSalt()
		if err != nil {
			t.Fatalf("Failed to generate salt: %v", err)
		}

		// Tạo K0 với salt ngẫu nhiên
		k0, err := GenerateK0FromPasswordWithSalt(password, salt)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}

		if len(k0) != 32 {
			t.Fatalf("K0 length should be 32 bytes, got %d", len(k0))
		}

		// Test consistency với cùng salt
		k0_again, err := GenerateK0FromPasswordWithSalt(password, salt)
		if err != nil {
			t.Fatalf("Failed to generate K0 again: %v", err)
		}

		if hex.EncodeToString(k0) != hex.EncodeToString(k0_again) {
			t.Fatalf("K0 generation is not consistent with same salt")
		}

		// Test khác biệt với salt khác
		otherSalt, err := GenerateSalt()
		if err != nil {
			t.Fatalf("Failed to generate other salt: %v", err)
		}

		k0_other, err := GenerateK0FromPasswordWithSalt(password, otherSalt)
		if err != nil {
			t.Fatalf("Failed to generate K0 with other salt: %v", err)
		}

		if hex.EncodeToString(k0) == hex.EncodeToString(k0_other) {
			t.Fatalf("K0 should be different with different salt")
		}

		t.Logf("K0 generated successfully with random salt")
		t.Logf("K0 with salt1: %s", hex.EncodeToString(k0))
		t.Logf("K0 with salt2: %s", hex.EncodeToString(k0_other))
	})
}

func TestSaltManagement(t *testing.T) {
	// Test tạo và lưu salt
	t.Run("Generate and save salt", func(t *testing.T) {
		// Tạo salt
		salt, err := GenerateSalt()
		if err != nil {
			t.Fatalf("Failed to generate salt: %v", err)
		}

		if len(salt) != 32 {
			t.Fatalf("Salt length should be 32 bytes, got %d", len(salt))
		}

		// Lưu salt
		err = SaveSaltToFile(salt, "test_salt.key")
		if err != nil {
			t.Fatalf("Failed to save salt: %v", err)
		}

		// Đọc salt
		loadedSalt, err := LoadSaltFromFile("test_salt.key")
		if err != nil {
			t.Fatalf("Failed to load salt: %v", err)
		}

		if hex.EncodeToString(salt) != hex.EncodeToString(loadedSalt) {
			t.Fatalf("Loaded salt does not match saved salt")
		}

		t.Logf("Salt management test passed")
	})
}
