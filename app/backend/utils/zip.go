package utils

import (
	"archive/zip"
	"os"
	"path/filepath"
)

func CreateZipFromFiles(zipPath string, files []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		fName := filepath.Base(file)
		f, err := zipWriter.Create(fName)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return err
		}
	}

	return nil
}
