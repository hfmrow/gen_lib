package typeMimeHttp

import (
	"net/http"
	"os"
)

// CheckMime: return the mime type of a file
func CheckMime(filename string) (mime string, err error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return mime, err
	}
	buff := make([]byte, 512)
	if _, err = file.Read(buff); err != nil {
		return mime, err
	}
	return http.DetectContentType(buff), err
}
