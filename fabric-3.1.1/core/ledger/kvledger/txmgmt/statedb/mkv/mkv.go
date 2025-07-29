//go:build !test
// +build !test

package mkv

/*
#cgo LDFLAGS: -L. -lmkv
#include "mkv.h"
*/
import "C"
import (
	"unsafe"
)

// EncryptValueMKV mã hóa value bằng MKV256
func EncryptValueMKV(value []byte, key []byte) []byte {
	if value == nil || len(value) == 0 || key == nil || len(key) == 0 {
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
		return nil
	}
	return ciphertext[:int(cCiphertextLen)]
}

// DecryptValueMKV giải mã value bằng MKV256
func DecryptValueMKV(value []byte, key []byte) []byte {
	if value == nil || len(value) == 0 || key == nil || len(key) == 0 {
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
		return nil
	}
	return plaintext[:int(cPlaintextLen)]
}
