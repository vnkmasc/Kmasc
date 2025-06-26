package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDate(input string) (time.Time, error) {
	input = strings.TrimSpace(input)

	if serial, err := strconv.ParseFloat(input, 64); err == nil {
		// Excel serial: bắt đầu từ 1899-12-30 (Excel bug năm nhuận)
		base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		return base.AddDate(0, 0, int(serial)), nil
	}

	// Parse theo các định dạng thông dụng
	layouts := []string{"02/01/2006", "2/1/2006", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %v", input)
}

func GetSafe(row []string, index int) string {
	if index < len(row) {
		return strings.TrimSpace(row[index])
	}
	return ""
}
