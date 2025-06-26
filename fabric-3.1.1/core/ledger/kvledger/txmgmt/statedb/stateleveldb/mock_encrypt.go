//go:build test
// +build test

package statedb

// EncryptValue mock: trả về nguyên input khi test
func EncryptValue(value []byte) []byte {
	return value
}

// DecryptValue mock: trả về nguyên input khi test
func DecryptValue(value []byte) []byte {
	return value
}
