package utils

import (
	"bytes"
	"fmt"
)

// AttachDigitalSignatureToPDF nhúng chữ ký số CMS vào cuối file PDF như một "embedded attachment".
func AttachDigitalSignatureToPDF(pdfData []byte, cmsSignature []byte) ([]byte, error) {
	if !bytes.HasPrefix(pdfData, []byte("%PDF")) {
		return nil, fmt.Errorf("file không phải là PDF hợp lệ")
	}

	// Tạo một comment chứa signature (bạn có thể thay bằng object / annotation tùy yêu cầu phức tạp)
	comment := []byte("\n%--CMS-Signature-Start--\n")
	comment = append(comment, cmsSignature...)
	comment = append(comment, []byte("\n%--CMS-Signature-End--\n")...)

	// Gắn vào cuối file PDF
	signedPDF := append(pdfData, comment...)

	return signedPDF, nil
}
