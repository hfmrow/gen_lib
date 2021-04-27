// types.go

package types

import (
	"fmt"
	"strconv"
	"strings"
)

// IsFloat: Check if string is float
func IsFloat(inString string) bool {
	_, err := strconv.ParseFloat(strings.Replace(strings.Replace(inString, " ", "", -1), ",", ".", -1), 64)
	if err == nil {
		return true
	}
	return false
}

// IsDate: Check if string is date
func IsDate(inString string) bool {
	dateFormats := NewDateFormat()
	for _, dteFmt := range dateFormats {
		if len(FindDate(inString, dteFmt+" %H:%M:%S")) != 0 {
			return true
		}
	}
	return false
}

// ByteCountDecimal: Format byte size to human readable format
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// ByteCountBinary: Format byte size to human readable format
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
