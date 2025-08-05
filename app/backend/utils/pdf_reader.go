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
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(html))))
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}
	return pdfg.Bytes(), nil
}
