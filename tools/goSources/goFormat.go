// goFormat.go

package goSources

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
)

// GoFormatBytes: Format bytes like gofmt function.
func GoFormatBytes(inBytes []byte) (outByte []byte, err error) {
	return format.Source(inBytes)
}

// GoFormatFile: Format file like gofmt function.
func GoFormatFile(filename string) (err error) {
	var data []byte
	var fi os.FileInfo

	if fi, err = os.Stat(filename); err == nil {
		if data, err = ioutil.ReadFile(filename); err == nil {
			if data, err = GoFormatBytes(data); err == nil {
				err = ioutil.WriteFile(filename, data, os.ModePerm&fi.Mode())
			}
		}
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("Issue occured while GoFormat the source file,"+
			" it may contain some errors, please check this out:\n[%s]\n%s\n", filename, err.Error()))
	}
	return
}

/*

var r io.Reader
var err error
r, err = os.Open("file.txt")

You can also make a Reader from a normal string using `strings.NewReader`:

var r io.Reader
r = strings.NewReader("Read will return these bytes")

A bytes.Buffer is a Reader:

var r io.Reader
var buf bytes.Buffer
r = &buf

*/
