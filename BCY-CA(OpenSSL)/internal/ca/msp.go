package ca

import (
	"os"
	"path/filepath"
)

func ExportMSP(mspDir, caCert, userCert, userKey string) error {
	dirs := []string{"cacerts", "signcerts", "keystore"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(mspDir, d), 0755); err != nil {
			return err
		}
	}
	if err := copyFile(caCert, filepath.Join(mspDir, "cacerts", "ca-cert.pem")); err != nil {
		return err
	}
	if err := copyFile(userCert, filepath.Join(mspDir, "signcerts", "cert.pem")); err != nil {
		return err
	}
	if err := copyFile(userKey, filepath.Join(mspDir, "keystore", "key.pem")); err != nil {
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}
