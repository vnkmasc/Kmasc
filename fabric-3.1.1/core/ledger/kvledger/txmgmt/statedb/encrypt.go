//go:build !test
// +build !test

package statedb

import (
	"fmt"
	"os"
	"sync"
	"time"

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
		logFile, logFileErr = os.OpenFile("/root/state_encryption_disabled.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	})
	if logFileErr != nil || logFile == nil {
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

// EncryptValue - DISABLED: trả về nguyên input (không mã hóa)
func EncryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("ENCRYPT_DISABLED", ns, key, "SKIP_EMPTY", "")
		return value
	}
	logToFile("ENCRYPT_DISABLED", ns, key, "SUCCESS", "Encryption disabled - returning original data")
	return value
}

// DecryptValue - DISABLED: trả về nguyên input (không giải mã)
func DecryptValue(value []byte, ns, key string) []byte {
	if value == nil || len(value) == 0 {
		logToFile("DECRYPT_DISABLED", ns, key, "SKIP_EMPTY", "")
		return value
	}
	logToFile("DECRYPT_DISABLED", ns, key, "SUCCESS", "Decryption disabled - returning original data")
	return value
}
