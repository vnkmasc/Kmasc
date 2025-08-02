package mkv

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncryptDecryptValueMKV(t *testing.T) {
	plaintext := []byte("Hello, this is a test for MKV encryption!")
	key := []byte("1234567890abcdef1234567890abcdef") // 32 bytes = 256 bit
	if len(key) != 32 {
		t.Fatalf("Key length is not 32 bytes: got %d", len(key))
	}
	t.Logf("Key (hex): %s", hex.EncodeToString(key))

	ciphertext := EncryptValueMKV(plaintext, key)
	if ciphertext == nil {
		t.Fatal("Encryption failed")
	}
	t.Logf("Ciphertext (hex): %s", hex.EncodeToString(ciphertext))

	decrypted := DecryptValueMKV(ciphertext, key)
	if decrypted == nil {
		t.Fatal("Decryption failed")
	}
	t.Logf("Decrypted: %s", string(decrypted))

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text does not match original.\nGot: %s\nWant: %s", decrypted, plaintext)
	}
}
