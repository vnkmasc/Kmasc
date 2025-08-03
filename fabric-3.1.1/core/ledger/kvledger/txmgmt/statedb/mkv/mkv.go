//go:build !test
// +build !test

package mkv

/*
#cgo LDFLAGS: -L. -lmkv
#include "mkv.h"
*/
import "C"
import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
	"unsafe"
)

var (
	logFileOnce sync.Once
	logFile     *os.File
	logFileErr  error
	logFileMu   sync.Mutex
)

// logToFileMKV ghi log các thao tác mã hóa/giải mã MKV vào file /root/state_mkv.log
func logToFileMKV(op, ns, key, status, errMsg string) {
	logFileOnce.Do(func() {
		logFile, logFileErr = os.OpenFile("/tmp/state_mkv.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	})
	if logFileErr != nil || logFile == nil {
		fmt.Fprintf(os.Stderr, "MKV log error: %v\n", logFileErr)
		return
	}
	logFileMu.Lock()
	defer logFileMu.Unlock()
	now := time.Now().UTC()
	timestamp := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%06dZ",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond()/1000)
	msg := timestamp + " " + op + " ns=" + ns + " key=" + key + " " + status
	if errMsg != "" {
		msg += " ERROR: " + errMsg
	}
	logFile.WriteString(msg + "\n")
}

// EncryptValueMKV mã hóa value bằng MKV256
func EncryptValueMKV(value []byte, key []byte) []byte {
	if value == nil || len(value) == 0 || key == nil || len(key) == 0 {
		logToFileMKV("ENCRYPT", "", "", "SKIP_EMPTY", "")
		return value
	}
	ciphertextLen := len(value) + 32 // padding tối đa 1 block
	ciphertext := make([]byte, ciphertextLen)
	var cPlaintext *C.uchar
	if len(value) > 0 {
		cPlaintext = (*C.uchar)(unsafe.Pointer(&value[0]))
	}
	cCiphertext := (*C.uchar)(unsafe.Pointer(&ciphertext[0]))
	cCiphertextLen := C.int(0)
	cKey := (*C.uchar)(unsafe.Pointer(&key[0]))
	keyLen := C.int(len(key) * 8) // bit
	ret := C.mkv_encrypt(cPlaintext, C.int(len(value)), cCiphertext, &cCiphertextLen, cKey, keyLen)
	if ret != 0 {
		logToFileMKV("ENCRYPT", "", "", "FAIL", "EncryptValueMKV error")
		return nil
	}
	logToFileMKV("ENCRYPT", "", "", "SUCCESS", "")
	return ciphertext[:int(cCiphertextLen)]
}

// DecryptValueMKV giải mã value bằng MKV256
func DecryptValueMKV(value []byte, key []byte) []byte {
	if value == nil || len(value) == 0 || key == nil || len(key) == 0 {
		logToFileMKV("DECRYPT", "", "", "SKIP_EMPTY", "")
		return value
	}
	plaintext := make([]byte, len(value))
	var cCiphertext *C.uchar
	if len(value) > 0 {
		cCiphertext = (*C.uchar)(unsafe.Pointer(&value[0]))
	}
	cPlaintext := (*C.uchar)(unsafe.Pointer(&plaintext[0]))
	cPlaintextLen := C.int(0)
	cKey := (*C.uchar)(unsafe.Pointer(&key[0]))
	keyLen := C.int(len(key) * 8)
	ret := C.mkv_decrypt(cCiphertext, C.int(len(value)), cPlaintext, &cPlaintextLen, cKey, keyLen)
	if ret != 0 {
		logToFileMKV("DECRYPT", "", "", "FAIL", "DecryptValueMKV error")
		return nil
	}
	logToFileMKV("DECRYPT", "", "", "SUCCESS", "")
	return plaintext[:int(cPlaintextLen)]
}

// GenerateK1 tạo khóa K1 ngẫu nhiên 32 bytes
func GenerateK1() ([]byte, error) {
	k1 := make([]byte, 32)
	_, err := rand.Read(k1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate K1: %v", err)
	}
	logToFileMKV("GENERATE_K1", "", "", "SUCCESS", "")
	return k1, nil
}

// GenerateK0FromPassword tạo khóa K0 từ password bằng PBKDF2 (đơn giản hóa bằng SHA256)
func GenerateK0FromPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	k0 := hash[:]
	logToFileMKV("GENERATE_K0", "", "", "SUCCESS", "")
	return k0
}

// EncryptK1WithK0 mã hóa K1 bằng K0 sử dụng MKV
func EncryptK1WithK0(k1 []byte, k0 []byte) []byte {
	return EncryptValueMKV(k1, k0)
}

// DecryptK1WithK0 giải mã K1 bằng K0 sử dụng MKV
func DecryptK1WithK0(encryptedK1 []byte, k0 []byte) []byte {
	return DecryptValueMKV(encryptedK1, k0)
}

// SaveK1ToFile lưu K1 vào file
func SaveK1ToFile(k1 []byte, filename string) error {
	err := ioutil.WriteFile(filename, k1, 0600)
	if err != nil {
		logToFileMKV("SAVE_K1", "", filename, "FAIL", err.Error())
		return fmt.Errorf("failed to save K1: %v", err)
	}
	logToFileMKV("SAVE_K1", "", filename, "SUCCESS", "")
	return nil
}

// LoadK1FromFile đọc K1 từ file
func LoadK1FromFile(filename string) ([]byte, error) {
	k1, err := ioutil.ReadFile(filename)
	if err != nil {
		logToFileMKV("LOAD_K1", "", filename, "FAIL", err.Error())
		return nil, fmt.Errorf("failed to load K1: %v", err)
	}
	logToFileMKV("LOAD_K1", "", filename, "SUCCESS", "")
	return k1, nil
}

// SaveEncryptedK1ToFile lưu K1 đã mã vào file
func SaveEncryptedK1ToFile(encryptedK1 []byte, filename string) error {
	err := ioutil.WriteFile(filename, encryptedK1, 0600)
	if err != nil {
		logToFileMKV("SAVE_ENCRYPTED_K1", "", filename, "FAIL", err.Error())
		return fmt.Errorf("failed to save encrypted K1: %v", err)
	}
	logToFileMKV("SAVE_ENCRYPTED_K1", "", filename, "SUCCESS", "")
	return nil
}

// LoadEncryptedK1FromFile đọc K1 đã mã từ file
func LoadEncryptedK1FromFile(filename string) ([]byte, error) {
	encryptedK1, err := ioutil.ReadFile(filename)
	if err != nil {
		logToFileMKV("LOAD_ENCRYPTED_K1", "", filename, "FAIL", err.Error())
		return nil, fmt.Errorf("failed to load encrypted K1: %v", err)
	}
	logToFileMKV("LOAD_ENCRYPTED_K1", "", filename, "SUCCESS", "")
	return encryptedK1, nil
}

// GetCurrentK1 lấy K1 hiện tại (giải mã từ file nếu cần)
func GetCurrentK1(password string) ([]byte, error) {
	// Thử đọc K1 đã mã trước
	encryptedK1, err := LoadEncryptedK1FromFile("encrypted_k1.key")
	if err != nil {
		// Nếu không có, thử đọc K1 plaintext
		k1, err := LoadK1FromFile("k1.key")
		if err != nil {
			return nil, fmt.Errorf("no K1 found: %v", err)
		}
		return k1, nil
	}

	// Giải mã K1 bằng password
	k0 := GenerateK0FromPassword(password)
	k1 := DecryptK1WithK0(encryptedK1, k0)
	if k1 == nil {
		return nil, fmt.Errorf("failed to decrypt K1 with password")
	}
	return k1, nil
}

// InitializeKeyManagement khởi tạo hệ thống quản lý khóa
func InitializeKeyManagement(password string) error {
	// Tạo K1 ngẫu nhiên
	k1, err := GenerateK1()
	if err != nil {
		return err
	}

	// Lưu K1 plaintext
	err = SaveK1ToFile(k1, "k1.key")
	if err != nil {
		return err
	}

	// Tạo K0 từ password
	k0 := GenerateK0FromPassword(password)

	// Mã K1 bằng K0
	encryptedK1 := EncryptK1WithK0(k1, k0)
	if encryptedK1 == nil {
		return fmt.Errorf("failed to encrypt K1 with K0")
	}

	// Lưu K1 đã mã
	err = SaveEncryptedK1ToFile(encryptedK1, "encrypted_k1.key")
	if err != nil {
		return err
	}

	logToFileMKV("INIT_KEYS", "", "", "SUCCESS", "")
	return nil
}

// ChangePassword thay đổi password (giải mã và mã lại K1)
func ChangePassword(oldPassword, newPassword string) error {
	// Lấy K1 hiện tại
	k1, err := GetCurrentK1(oldPassword)
	if err != nil {
		return fmt.Errorf("failed to get current K1: %v", err)
	}

	// Tạo K0 mới từ password mới
	newK0 := GenerateK0FromPassword(newPassword)

	// Mã lại K1 bằng K0 mới
	encryptedK1 := EncryptK1WithK0(k1, newK0)
	if encryptedK1 == nil {
		return fmt.Errorf("failed to encrypt K1 with new K0")
	}

	// Lưu K1 đã mã mới
	err = SaveEncryptedK1ToFile(encryptedK1, "encrypted_k1.key")
	if err != nil {
		return err
	}

	logToFileMKV("CHANGE_PASSWORD", "", "", "SUCCESS", "")
	return nil
}
