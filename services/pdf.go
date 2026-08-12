package services

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(path string) (string, error) {

	file, reader, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer

	totalPage := reader.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {

		page := reader.Page(pageIndex)

		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}

		buf.WriteString(text)
	}

	return buf.String(), nil
}
