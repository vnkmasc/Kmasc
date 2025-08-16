package mkv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// KeyManager là singleton để quản lý khóa mã hóa
type KeyManager struct {
	encryptionKey []byte
	password      string
	initialized   bool
}

var (
	keyManagerInstance *KeyManager
	keyManagerOnce     sync.Once
	keyManagerMutex    sync.RWMutex
)

// readPasswordFromFile đọc password từ file password.txt
func readPasswordFromFile() (string, error) {
	// Thử đọc từ thư mục hiện tại trước
	passwordPath := "password.txt"
	if _, err := os.Stat(passwordPath); os.IsNotExist(err) {
		// Nếu không có, thử đọc từ /tmp/mkv
		passwordPath = "/tmp/mkv/password.txt"
	}

	passwordBytes, err := ioutil.ReadFile(passwordPath)
	if err != nil {
		return "", fmt.Errorf("failed to read password file %s: %v", passwordPath, err)
	}

	// Loại bỏ whitespace và newline
	password := strings.TrimSpace(string(passwordBytes))
	if password == "" {
		return "", fmt.Errorf("password file is empty")
	}

	return password, nil
}

// GetKeyManager trả về instance duy nhất của KeyManager
func GetKeyManager() *KeyManager {
	keyManagerOnce.Do(func() {
		// Đọc password từ file
		password, err := readPasswordFromFile()
		if err != nil {
			logToFileMKV("KEY_MANAGER_INIT", "", "", "WARN", "Failed to read password file: "+err.Error()+", using default password")
			// Fallback to default password if file read fails
			password = "kmasc"
		}

		keyManagerInstance = &KeyManager{
			password: password,
		}
		// Tự động khởi tạo khi lần đầu được gọi
		if err := keyManagerInstance.initialize(); err != nil {
			// Log lỗi nhưng không panic, sử dụng fallback key
			logToFileMKV("KEY_MANAGER_INIT", "", "", "ERROR", err.Error())
		}
	})
	return keyManagerInstance
}

// initialize khởi tạo hệ thống quản lý khóa
func (km *KeyManager) initialize() error {
	keyManagerMutex.Lock()
	defer keyManagerMutex.Unlock()

	if km.initialized {
		return nil
	}

	logToFileMKV("KEY_MANAGER_INIT", "", "", "START", "Initializing key management system")

	// Kiểm tra xem đã có key file chưa
	if km.hasExistingKeys() {
		// Load key hiện có
		if err := km.loadExistingKeys(); err != nil {
			logToFileMKV("KEY_MANAGER_INIT", "", "", "WARN", "Failed to load existing keys: "+err.Error())
			// Fallback: tạo key mới
			return km.createNewKeySystem()
		}
		logToFileMKV("KEY_MANAGER_INIT", "", "", "SUCCESS", "Loaded existing keys")
	} else {
		// Tạo hệ thống key mới
		if err := km.createNewKeySystem(); err != nil {
			return err
		}
		logToFileMKV("KEY_MANAGER_INIT", "", "", "SUCCESS", "Created new key system")
	}

	km.initialized = true
	return nil
}

// hasExistingKeys kiểm tra xem đã có key files chưa
func (km *KeyManager) hasExistingKeys() bool {
	searchPaths := km.getSearchPaths()
	for _, path := range searchPaths {
		encryptedKeyPath := filepath.Join(path, "encrypted_k1.key")
		if _, err := os.Stat(encryptedKeyPath); err == nil {
			return true
		}
	}
	return false
}

// loadExistingKeys load key hiện có từ file
func (km *KeyManager) loadExistingKeys() error {
	// Tìm file encrypted key
	encryptedKeyPath, err := km.findKeyFile("encrypted_k1.key")
	if err != nil {
		return fmt.Errorf("encrypted key file not found: %v", err)
	}

	// Load encrypted key
	encryptedKey, err := LoadEncryptedK1FromFile(encryptedKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load encrypted key: %v", err)
	}

	// Tạo K0 từ password
	k0, err := GenerateK0FromPassword(km.password)
	if err != nil {
		return fmt.Errorf("failed to generate K0: %v", err)
	}

	// Giải mã K1
	k1 := DecryptK1WithK0(encryptedKey, k0)
	if k1 == nil {
		return fmt.Errorf("failed to decrypt K1")
	}

	km.encryptionKey = k1
	logToFileMKV("KEY_MANAGER_LOAD", "", encryptedKeyPath, "SUCCESS", "")
	return nil
}

// createNewKeySystem tạo hệ thống key mới
func (km *KeyManager) createNewKeySystem() error {
	logToFileMKV("KEY_MANAGER_CREATE", "", "", "START", "Creating new key system")

	// Tạo K1 ngẫu nhiên
	k1, err := GenerateK1()
	if err != nil {
		return fmt.Errorf("failed to generate K1: %v", err)
	}

	// Tạo salt ngẫu nhiên
	salt, err := GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	// Tạo K0 từ password với salt
	k0, err := GenerateK0FromPasswordWithSalt(km.password, salt)
	if err != nil {
		return fmt.Errorf("failed to generate K0: %v", err)
	}

	// Mã hóa K1 bằng K0
	encryptedK1 := EncryptK1WithK0(k1, k0)
	if encryptedK1 == nil {
		return fmt.Errorf("failed to encrypt K1")
	}

	// Lưu tất cả keys vào các vị trí
	if err := km.saveAllKeys(k1, k0, salt, encryptedK1); err != nil {
		return fmt.Errorf("failed to save keys: %v", err)
	}

	km.encryptionKey = k1
	logToFileMKV("KEY_MANAGER_CREATE", "", "", "SUCCESS", "New key system created")
	return nil
}

// saveAllKeys lưu tất cả keys vào các vị trí
func (km *KeyManager) saveAllKeys(k1, k0, salt, encryptedK1 []byte) error {
	savePaths := km.getSearchPaths()

	for _, path := range savePaths {
		// Tạo thư mục nếu chưa có
		if path != "." {
			if err := os.MkdirAll(path, 0755); err != nil {
				logToFileMKV("KEY_MANAGER_SAVE", "", path, "FAIL", "Failed to create directory: "+err.Error())
				continue
			}
		}

		// Lưu K1 plaintext
		k1Path := filepath.Join(path, "k1.key")
		if err := SaveK1ToFile(k1, k1Path); err != nil {
			logToFileMKV("KEY_MANAGER_SAVE", "", k1Path, "FAIL", err.Error())
		}

		// Lưu salt
		saltPath := filepath.Join(path, "k0_salt.key")
		if err := SaveSaltToFile(salt, saltPath); err != nil {
			logToFileMKV("KEY_MANAGER_SAVE", "", saltPath, "FAIL", err.Error())
		}

		// Lưu K0
		k0Path := filepath.Join(path, "k0.key")
		if err := ioutil.WriteFile(k0Path, k0, 0600); err != nil {
			logToFileMKV("KEY_MANAGER_SAVE", "", k0Path, "FAIL", err.Error())
		}

		// Lưu K1 đã mã hóa
		encryptedK1Path := filepath.Join(path, "encrypted_k1.key")
		if err := SaveEncryptedK1ToFile(encryptedK1, encryptedK1Path); err != nil {
			logToFileMKV("KEY_MANAGER_SAVE", "", encryptedK1Path, "FAIL", err.Error())
		}

		// Lưu password (chỉ để tham khảo, trong thực tế nên bảo mật hơn)
		passwordPath := filepath.Join(path, "password.txt")
		if err := ioutil.WriteFile(passwordPath, []byte(km.password), 0600); err != nil {
			logToFileMKV("KEY_MANAGER_SAVE", "", passwordPath, "FAIL", err.Error())
		}

		logToFileMKV("KEY_MANAGER_SAVE", "", path, "SUCCESS", "")
	}

	return nil
}

// getSearchPaths trả về danh sách các đường dẫn tìm kiếm
func (km *KeyManager) getSearchPaths() []string {
	return []string{
		".",                   // Current directory
		"/tmp",                // Temp directory
		"/tmp/mkv",            // Temp MKV directory
		"/opt/mkv",            // Opt MKV directory
		"/home/chaincode/mkv", // Chaincode MKV directory
		"/root/mkv",           // Root MKV directory
	}
}

// findKeyFile tìm file key trong các đường dẫn tìm kiếm
func (km *KeyManager) findKeyFile(filename string) (string, error) {
	searchPaths := km.getSearchPaths()

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, filename)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("key file %s not found in any search path", filename)
}

// GetEncryptionKey trả về khóa mã hóa hiện tại
func (km *KeyManager) GetEncryptionKey() []byte {
	keyManagerMutex.RLock()
	defer keyManagerMutex.RUnlock()

	if !km.initialized {
		// Nếu chưa khởi tạo, tự động khởi tạo
		go func() {
			if err := km.initialize(); err != nil {
				logToFileMKV("KEY_MANAGER_AUTO_INIT", "", "", "ERROR", err.Error())
			}
		}()

		// Trả về fallback key trong khi chờ khởi tạo
		return []byte("1234567890abcdef1234567890abcdef")
	}

	return km.encryptionKey
}

// ChangePassword thay đổi password và mã lại key
func (km *KeyManager) ChangePassword(newPassword string) error {
	keyManagerMutex.Lock()
	defer keyManagerMutex.Unlock()

	if !km.initialized {
		return fmt.Errorf("key manager not initialized")
	}

	// Tạo K0 mới từ password mới
	newK0, err := GenerateK0FromPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to generate new K0: %v", err)
	}

	// Mã lại K1 bằng K0 mới
	encryptedK1 := EncryptK1WithK0(km.encryptionKey, newK0)
	if encryptedK1 == nil {
		return fmt.Errorf("failed to encrypt K1 with new K0")
	}

	// Lưu K1 đã mã mới vào tất cả các vị trí
	savePaths := km.getSearchPaths()
	for _, path := range savePaths {
		encryptedK1Path := filepath.Join(path, "encrypted_k1.key")
		if err := SaveEncryptedK1ToFile(encryptedK1, encryptedK1Path); err != nil {
			logToFileMKV("KEY_MANAGER_CHANGE_PASS", "", encryptedK1Path, "FAIL", err.Error())
		}

		// Cập nhật password file
		passwordPath := filepath.Join(path, "password.txt")
		if err := ioutil.WriteFile(passwordPath, []byte(newPassword), 0600); err != nil {
			logToFileMKV("KEY_MANAGER_CHANGE_PASS", "", passwordPath, "FAIL", err.Error())
		}
	}

	km.password = newPassword
	logToFileMKV("KEY_MANAGER_CHANGE_PASS", "", "", "SUCCESS", "")
	return nil
}

// GetStatus trả về trạng thái hiện tại của key manager
func (km *KeyManager) GetStatus() map[string]interface{} {
	keyManagerMutex.RLock()
	defer keyManagerMutex.RUnlock()

	status := map[string]interface{}{
		"initialized":  km.initialized,
		"has_key":      len(km.encryptionKey) > 0,
		"key_length":   len(km.encryptionKey),
		"password_set": km.password != "",
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
	}

	return status
}

// RefreshKeys làm mới hệ thống key (tạo key mới)
func (km *KeyManager) RefreshKeys() error {
	keyManagerMutex.Lock()
	defer keyManagerMutex.Unlock()

	logToFileMKV("KEY_MANAGER_REFRESH", "", "", "START", "Refreshing key system")

	// Tạo key mới
	if err := km.createNewKeySystem(); err != nil {
		return fmt.Errorf("failed to refresh keys: %v", err)
	}

	logToFileMKV("KEY_MANAGER_REFRESH", "", "", "SUCCESS", "Keys refreshed successfully")
	return nil
}
