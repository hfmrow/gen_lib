// gZip.go

/*
	©2021 https://github/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gZip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
)

// GzipNew: create a new gzip file with given files. -2 >= compLvl <= 9
func GzipNew(gzipFilename string, files []string, compLvl ...int) error {

	var (
		errorStr string
		cLvl     = gzip.DefaultCompression
	)

	if len(compLvl) > 0 {
		cLvl = compLvl[0]
	}

	f, err := os.Create(gzipFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := gzip.NewWriterLevel(f, cLvl)
	if err != nil {
		return err
	}
	defer w.Close()

	for _, file := range files {
		if err = writeGz(w, file); err == nil {
			err = w.Flush()
		}
		if err != nil {
			errorStr += fmt.Sprintf("Error: %s, %v\n", file, err)
		}
	}

	if len(errorStr) > 0 {
		return fmt.Errorf("Some adverse events have occurred:\n%s", errorStr)
	}
	return nil
}

// writeGz:
func writeGz(w *gzip.Writer, filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	wrote, err := w.Write(data)
	if wrote != len(data) {
		return fmt.Errorf("Bytes written: %d, expected: %d", wrote, len(data))
	}
	return err
}

// func main() {

// 	name_of_file := "Gfg.txt"

// 	f, _ := os.Open("C://ProgramData//" + name_of_file)

// 	read := bufio.NewReader(f)

// 	data, _ := ioutil.ReadAll(read)

// 	name_of_file = strings.Replace(name_of_file, ".txt", ".gz", -1)

// 	f, _ = os.Create("C://ProgramData//" + name_of_file)

// 	w := gzip.NewWriter(f)

// 	w.Write(data)

// 	w.Close()
// }
