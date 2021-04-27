// convert.go

package convert

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// ByteToHexStr: Convert []byte to hexString
func ByteToHexStr(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// Convert comma to dot if needed and return 0 if input string is empty.
func StringDecimalSwitchFloat(decimalChar, inString string) float64 {
	if inString == "" {
		inString = "0"
	}
	switch decimalChar {
	case ",":
		f, _ := strconv.ParseFloat(strings.Replace(inString, ",", ".", 1), 64)
		return f
	case ".":
		f, _ := strconv.ParseFloat(inString, 64)
		return f
	}
	return -1
}
