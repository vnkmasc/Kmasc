package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// removeDiacritics loại bỏ dấu tiếng Việt
func removeDiacritics(str string) string {
	t := norm.NFD.String(str)
	result := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func Slugify(input string) string {
	// Bước 1: chuẩn hóa
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, "chứng chỉ", "")
	input = strings.ReplaceAll(input, "bằng", "")
	input = strings.TrimSpace(input)

	// Bước 2: loại bỏ dấu tiếng Việt
	input = removeDiacritics(input)

	// Bước 3: thay thế ký tự
	input = strings.ReplaceAll(input, "+", "plus")
	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, ":", "")
	input = strings.ReplaceAll(input, " ", "-")

	// Bước 4: giữ lại ký tự hợp lệ
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	input = re.ReplaceAllString(input, "")

	// Bước 5: xóa dấu `-` thừa
	input = strings.Trim(input, "-")

	return input
}
