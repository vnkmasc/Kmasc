//go:build !test
// +build !test

package statedb

/*
#include "encrypt.h"
*/
import "C"

import (
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/hyperledger/fabric-lib-go/common/flogging"
)

var encryptLogger = flogging.MustGetLogger("encrypt")

var (
	logFileOnce sync.Once
	logFile     *os.File
	logFileErr  error
	logFileMu   sync.Mutex
)

func logToFile(op, ns, key, status, errMsg string) {
	logFileOnce.Do(func() {
		logFile, logFileErr = os.OpenFile("/root/state_encryption.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	})
	if logFileErr != nil || logFile == nil {
		return
	}
	logFileMu.Lock()
	defer logFileMu.Unlock()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	msg := timestamp + " " + op + " ns=" + ns + " key=" + key + " " + status
	if errMsg != "" {
		msg += " ERROR: " + errMsg
	}
	logFile.WriteString(msg + "\n")
}

// EncryptValue mã hóa giá trị sử dụng hàm C
func EncryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("ENCRYPT", ns, key, "SKIP_EMPTY", "")
		return value
	}
	ciphertextLen := len(value) + 32 // Thêm padding cho AES block size
	ciphertext := make([]byte, ciphertextLen)
	var cPlaintext *C.uchar
	if len(value) > 0 {
		cPlaintext = (*C.uchar)(unsafe.Pointer(&value[0]))
	}
	cCiphertext := (*C.uchar)(unsafe.Pointer(&ciphertext[0]))
	cCiphertextLen := C.int(0)
	result := C.encrypt_aes_cbc(cPlaintext, C.int(len(value)), cCiphertext, &cCiphertextLen)
	if result != 0 {
		logToFile("ENCRYPT", ns, key, "FAIL", "C function error")
		return nil
	}
	logToFile("ENCRYPT", ns, key, "SUCCESS", "")
	encryptedData := ciphertext[:int(cCiphertextLen)]
	return encryptedData
}

// DecryptValue giải mã giá trị sử dụng hàm C
func DecryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("DECRYPT", ns, key, "SKIP_EMPTY", "")
		return value
	}
	plaintextLen := len(value)
	plaintext := make([]byte, plaintextLen)
	var cCiphertext *C.uchar
	if len(value) > 0 {
		cCiphertext = (*C.uchar)(unsafe.Pointer(&value[0]))
	}
	cPlaintext := (*C.uchar)(unsafe.Pointer(&plaintext[0]))
	cPlaintextLen := C.int(0)
	result := C.decrypt_aes_cbc(cCiphertext, C.int(len(value)), cPlaintext, &cPlaintextLen)
	if result != 0 {
		logToFile("DECRYPT", ns, key, "FAIL", "C function error")
		return nil
	}
	logToFile("DECRYPT", ns, key, "SUCCESS", "")
	decryptedData := plaintext[:int(cPlaintextLen)]
	return decryptedData
}
