package aeslib

/*
#cgo CFLAGS: -I../../crypto
#cgo LDFLAGS: -L../../crypto -laes_encryptor
#include "aes_encryptor.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	AESKeySize = 32
	AESIVSize  = 16
)

func GenerateAESKey() []byte {
	key := make([]byte, AESKeySize)
	C.generate_aes_key((*C.uchar)(unsafe.Pointer(&key[0])))
	return key
}

func GenerateIV() []byte {
	iv := make([]byte, AESIVSize)
	C.generate_iv((*C.uchar)(unsafe.Pointer(&iv[0])))
	return iv
}

// EncryptAES dùng module C để mã hóa data với key và IV
func EncryptAES(data, key, iv []byte) ([]byte, error) {
	var outLen C.size_t

	cipher := C.aes_encrypt(
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&key[0])),
		(*C.uchar)(unsafe.Pointer(&iv[0])),
		&outLen,
	)

	if cipher == nil {
		return nil, fmt.Errorf("encryption failed")
	}
	defer C.free(unsafe.Pointer(cipher))

	encrypted := C.GoBytes(unsafe.Pointer(cipher), C.int(outLen))
	return encrypted, nil
}

// DecryptAES dùng module C để giải mã data với key và IV
func DecryptAES(encData, key, iv []byte) ([]byte, error) {
	var outLen C.size_t

	plain := C.aes_decrypt(
		(*C.uchar)(unsafe.Pointer(&encData[0])),
		C.size_t(len(encData)),
		(*C.uchar)(unsafe.Pointer(&key[0])),
		(*C.uchar)(unsafe.Pointer(&iv[0])),
		&outLen,
	)

	if plain == nil {
		return nil, fmt.Errorf("decryption failed")
	}
	defer C.free(unsafe.Pointer(plain))

	decrypted := C.GoBytes(unsafe.Pointer(plain), C.int(outLen))
	return decrypted, nil
}
