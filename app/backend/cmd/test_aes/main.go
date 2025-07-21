package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$yljkp7BGUe6ARkRb5E39p.tJKr5meWvDD0Wfi6ZTvSQRxz4Smreea"
	password := "1"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ Sai mật khẩu:", err)
	} else {
		fmt.Println("✅ Mật khẩu đúng")
	}
}
