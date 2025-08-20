package utils

import (
	"bytes"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (pg *PDFGenerator) ConvertHTMLToPDF(html string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// Tạo page với UTF-8 encoding
	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(html)))
	page.Encoding.Set("utf-8")

	pdfg.AddPage(page)

	// Tạo PDF
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
