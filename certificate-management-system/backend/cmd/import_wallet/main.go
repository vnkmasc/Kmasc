package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	// Lấy đường dẫn đến thư mục MSP từ biến môi trường
	credPath := os.Getenv("FABRIC_ADMIN_CRED_PATH")
	if credPath == "" {
		panic("biến môi trường FABRIC_ADMIN_CRED_PATH chưa được thiết lập")
	}

	// Lấy MSP ID từ biến môi trường
	mspID := os.Getenv("FABRIC_MSP_ID")
	if mspID == "" {
		panic("biến môi trường FABRIC_MSP_ID chưa được thiết lập")
	}

	// Lấy tên định danh trong ví
	identityLabel := os.Getenv("FABRIC_IDENTITY")
	if identityLabel == "" {
		identityLabel = "admin" // fallback nếu không thiết lập
	}

	// Đọc file cert
	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc cert: %w", err))
	}

	// Đọc file private key từ thư mục keystore
	keyDir := filepath.Join(credPath, "keystore")
	keyFiles, err := os.ReadDir(keyDir)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc thư mục keystore: %w", err))
	}
	if len(keyFiles) == 0 {
		panic("không tìm thấy file private key trong keystore")
	}
	keyPath := filepath.Join(keyDir, keyFiles[0].Name())
	key, err := os.ReadFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("lỗi đọc private key: %w", err))
	}

	// Tạo ví và import identity
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		panic(fmt.Errorf("lỗi tạo wallet: %w", err))
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))
	if err = wallet.Put(identityLabel, identity); err != nil {
		panic(fmt.Errorf("lỗi import identity: %w", err))
	}

	fmt.Printf(" Đã import %s (%s) vào ví 'wallet'\n", identityLabel, mspID)
}
