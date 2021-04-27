// genericReader.go

/*
	©2018-21 https://github/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package compress

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"strings"
)

// GenReaderCompressed: return io.reader for stdIn, .bz2, .gz, standard files.
func GenReaderCompressed(filename string) (io.Reader, *os.File, error) {
	if filename == "" {
		return bufio.NewReader(os.Stdin), nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(filename, "bz2") {
		return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err
	}

	if strings.HasSuffix(filename, "gz") {
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, nil, err
		}
		return bufio.NewReader(reader), file, err
	}
	return bufio.NewReader(file), file, err
}
