package mkv

import (
	"encoding/hex"
	"testing"
)

func TestKeyManagementSystem(t *testing.T) {
	// Test 1: Generate K1
	t.Run("Generate K1", func(t *testing.T) {
		k1, err := GenerateK1()
		if err != nil {
			t.Fatalf("Failed to generate K1: %v", err)
		}
		if len(k1) != 32 {
			t.Fatalf("K1 length should be 32 bytes, got %d", len(k1))
		}
		t.Logf("K1 generated successfully: %s", hex.EncodeToString(k1))
	})

	// Test 2: Generate K0 from password
	t.Run("Generate K0 from password", func(t *testing.T) {
		password := "mysecretpassword123"
		k0, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}
		if len(k0) != 32 {
			t.Fatalf("K0 length should be 32 bytes, got %d", len(k0))
		}
		t.Logf("K0 generated from password: %s", hex.EncodeToString(k0))

		// Test consistency
		k0_again, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0 again: %v", err)
		}
		if hex.EncodeToString(k0) != hex.EncodeToString(k0_again) {
			t.Fatalf("K0 generation is not consistent")
		}
	})

	// Test 3: Encrypt and decrypt K1 with K0
	t.Run("Encrypt/Decrypt K1 with K0", func(t *testing.T) {
		// Generate K1
		k1, err := GenerateK1()
		if err != nil {
			t.Fatalf("Failed to generate K1: %v", err)
		}

		// Generate K0 from password
		password := "testpassword456"
		k0, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}

		// Encrypt K1 with K0
		encryptedK1 := EncryptK1WithK0(k1, k0)
		if encryptedK1 == nil {
			t.Fatalf("Failed to encrypt K1 with K0")
		}
		t.Logf("K1 encrypted with K0, encrypted length: %d", len(encryptedK1))

		// Decrypt K1 with K0
		decryptedK1 := DecryptK1WithK0(encryptedK1, k0)
		if decryptedK1 == nil {
			t.Fatalf("Failed to decrypt K1 with K0")
		}

		// Verify decryption
		if hex.EncodeToString(k1) != hex.EncodeToString(decryptedK1) {
			t.Fatalf("Decrypted K1 does not match original K1")
		}
		t.Logf("K1 decrypted successfully, matches original")
	})

	// Test 4: File operations
	t.Run("File operations", func(t *testing.T) {
		// Generate test data
		k1, err := GenerateK1()
		if err != nil {
			t.Fatalf("Failed to generate K1: %v", err)
		}

		// Save K1 to file
		err = SaveK1ToFile(k1, "test_k1.key")
		if err != nil {
			t.Fatalf("Failed to save K1 to file: %v", err)
		}

		// Load K1 from file
		loadedK1, err := LoadK1FromFile("test_k1.key")
		if err != nil {
			t.Fatalf("Failed to load K1 from file: %v", err)
		}

		// Verify loaded data
		if hex.EncodeToString(k1) != hex.EncodeToString(loadedK1) {
			t.Fatalf("Loaded K1 does not match saved K1")
		}
		t.Logf("File operations successful")
	})

	// Test 5: Encrypted file operations
	t.Run("Encrypted file operations", func(t *testing.T) {
		// Generate test data
		k1, err := GenerateK1()
		if err != nil {
			t.Fatalf("Failed to generate K1: %v", err)
		}

		password := "filetestpassword"
		k0, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}

		// Encrypt K1
		encryptedK1 := EncryptK1WithK0(k1, k0)
		if encryptedK1 == nil {
			t.Fatalf("Failed to encrypt K1")
		}

		// Save encrypted K1
		err = SaveEncryptedK1ToFile(encryptedK1, "test_encrypted_k1.key")
		if err != nil {
			t.Fatalf("Failed to save encrypted K1: %v", err)
		}

		// Load encrypted K1
		loadedEncryptedK1, err := LoadEncryptedK1FromFile("test_encrypted_k1.key")
		if err != nil {
			t.Fatalf("Failed to load encrypted K1: %v", err)
		}

		// Decrypt loaded K1
		decryptedK1 := DecryptK1WithK0(loadedEncryptedK1, k0)
		if decryptedK1 == nil {
			t.Fatalf("Failed to decrypt loaded K1")
		}

		// Verify
		if hex.EncodeToString(k1) != hex.EncodeToString(decryptedK1) {
			t.Fatalf("Decrypted K1 does not match original")
		}
		t.Logf("Encrypted file operations successful")
	})

	// Test 6: Initialize key management
	t.Run("Initialize key management", func(t *testing.T) {
		password := "initpassword789"

		err := InitializeKeyManagement(password)
		if err != nil {
			t.Fatalf("Failed to initialize key management: %v", err)
		}

		// Verify files were created
		k1, err := LoadK1FromFile("k1.key")
		if err != nil {
			t.Fatalf("Failed to load k1.key: %v", err)
		}

		encryptedK1, err := LoadEncryptedK1FromFile("encrypted_k1.key")
		if err != nil {
			t.Fatalf("Failed to load encrypted_k1.key: %v", err)
		}

		// Verify we can decrypt with password
		k0, err := GenerateK0FromPassword(password)
		if err != nil {
			t.Fatalf("Failed to generate K0: %v", err)
		}
		decryptedK1 := DecryptK1WithK0(encryptedK1, k0)
		if decryptedK1 == nil {
			t.Fatalf("Failed to decrypt K1 with password")
		}

		if hex.EncodeToString(k1) != hex.EncodeToString(decryptedK1) {
			t.Fatalf("Decrypted K1 does not match original K1")
		}

		t.Logf("Key management initialization successful")
	})

	// Test 7: Change password
	t.Run("Change password", func(t *testing.T) {
		oldPassword := "oldpassword123"
		newPassword := "newpassword456"

		// First initialize with old password
		err := InitializeKeyManagement(oldPassword)
		if err != nil {
			t.Fatalf("Failed to initialize with old password: %v", err)
		}

		// Get original K1
		originalK1, err := GetCurrentK1(oldPassword)
		if err != nil {
			t.Fatalf("Failed to get current K1: %v", err)
		}

		// Change password
		err = ChangePassword(oldPassword, newPassword)
		if err != nil {
			t.Fatalf("Failed to change password: %v", err)
		}

		// Verify we can get K1 with new password
		newK1, err := GetCurrentK1(newPassword)
		if err != nil {
			t.Fatalf("Failed to get K1 with new password: %v", err)
		}

		// Verify K1 is the same
		if hex.EncodeToString(originalK1) != hex.EncodeToString(newK1) {
			t.Fatalf("K1 changed after password change")
		}

		// Verify old password no longer works
		_, err = GetCurrentK1(oldPassword)
		if err == nil {
			t.Fatalf("Old password still works after password change")
		}

		t.Logf("Password change successful")
	})
}

// Test data encryption with K1
func TestDataEncryptionWithK1(t *testing.T) {
	// Initialize key management
	password := "datatestpassword"
	err := InitializeKeyManagement(password)
	if err != nil {
		t.Fatalf("Failed to initialize key management: %v", err)
	}

	// Get K1
	k1, err := GetCurrentK1(password)
	if err != nil {
		t.Fatalf("Failed to get K1: %v", err)
	}

	// Test data encryption
	testData := []byte("This is test data for MKV encryption with K1")

	// Encrypt data with K1
	encryptedData := EncryptValueMKV(testData, k1)
	if encryptedData == nil {
		t.Fatalf("Failed to encrypt data with K1")
	}

	// Decrypt data with K1
	decryptedData := DecryptValueMKV(encryptedData, k1)
	if decryptedData == nil {
		t.Fatalf("Failed to decrypt data with K1")
	}

	// Verify
	if string(testData) != string(decryptedData) {
		t.Fatalf("Decrypted data does not match original")
	}

	t.Logf("Data encryption with K1 successful")
	t.Logf("Original data: %s", string(testData))
	t.Logf("Encrypted length: %d", len(encryptedData))
	t.Logf("Decrypted data: %s", string(decryptedData))
}
