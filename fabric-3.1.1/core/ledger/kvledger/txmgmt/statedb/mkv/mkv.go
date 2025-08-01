//go:build !test
// +build !test

package mkv

/*
#cgo LDFLAGS: -L. -lmkv
#include "mkv.h"
*/
import "C"
import (
	"fmt"
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
