// csvStruct.go

/*
	©2019 H.F.M
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This library allow to facilitate C.S.V operations
*/

package csvHandler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html/charset"
	iconv "gopkg.in/iconv.v1"
)

type CsvStructure struct {
	// Public
	Ready               bool // true when csv data are present.
	Loaded              bool // true when loaded successfully.
	Comma               rune
	Comment             rune
	RowWithColName      int
	LazyQuotes          bool
	FieldsPerRecord     int
	TrimLeadingSpace    bool
	UseCRLF             bool
	Charset             string
	KeepOriginalCharset bool
	CsvLines            [][]string
	RawData             []string

	// private
	analyseRowCount  int
	originalFilename string
	// encodingCharset encoding.Encoding
}

type rowStore struct {
	Idx int
	Cnt int
	Tot int
	Str string
}

func CsvStructureNew() *CsvStructure {
	cs := new(CsvStructure)
	cs.TrimLeadingSpace = true
	cs.Comma = ','
	cs.Comment = '#'
	cs.Charset = "utf-8"
	cs.KeepOriginalCharset = true
	cs.analyseRowCount = 50
	return cs
}

// Read: open a reader and read entire file to memory and store datas into CsvStructure
func (cs *CsvStructure) Read(filename string, autoOptions bool) (err error) {
	var file *os.File
	var csvReader *csv.Reader
	var stats os.FileInfo
	var byteRead []byte
	var bytesBuffer = &bytes.Buffer{}
	var text string
	var bytesReaded int
	cs.Loaded = false
	if file, err = os.Open(filename); err == nil {
		defer func() {
			file.Close()
		}()
		// Read file to Check for CRLF
		if stats, err = file.Stat(); err == nil {
			byteRead = make([]byte, stats.Size())
			if bytesReaded, err = file.Read(byteRead); err == nil {
				if bytesReaded < 10 {
					return errors.New(fmt.Sprintf("File size is: %d, unable to parse this csv content.", bytesReaded))
				}
				cs.UseCRLF = bytes.Contains(byteRead, []byte{0x0D, 0x0A})
				byteRead = cs.checkForConvertToUtf8(byteRead)
				if cs.UseCRLF {
					cs.RawData = strings.Split(string(byteRead), string([]byte{0x0D, 0x0A}))
				} else {
					cs.RawData = strings.Split(string(byteRead), string([]byte{0x0A}))
				}
				if autoOptions { // Try to identify csv parameters
					cs.autoSetOptions()
				}
				tmpRawData := cs.RawData[cs.RowWithColName:]
				if cs.UseCRLF {
					text = strings.Join(tmpRawData, string([]byte{0x0D, 0x0A}))
				} else {
					text = strings.Join(tmpRawData, string([]byte{0x0A}))
				}
				if _, err = bytesBuffer.WriteString(text); err == nil {
					csvReader = csv.NewReader(bytesBuffer)
					csvReader.LazyQuotes = cs.LazyQuotes
					csvReader.Comma = cs.Comma
					csvReader.Comment = cs.Comment
					csvReader.FieldsPerRecord = cs.FieldsPerRecord
					if cs.CsvLines, err = csvReader.ReadAll(); err == nil {
						cs.Ready = true
						cs.Loaded = true
					}
				}
			}
		}
	}
	return err
}

// ReadRawData: Open CSV file as raw data (without csv formatting).
func (cs *CsvStructure) ReadRawData(filename string, autoOptions bool) (err error) {
	var file *os.File
	var outBytes []byte
	var stats os.FileInfo
	if file, err = os.Open(filename); err == nil {
		defer file.Close()
		if stats, err = file.Stat(); err == nil { // Read file and Check for CRLF
			outBytes = make([]byte, stats.Size())
			if _, err = file.Read(outBytes); err == nil {
				cs.UseCRLF = bytes.Contains(outBytes, []byte{0x0D, 0x0A})
				outBytes = cs.checkForConvertToUtf8(outBytes)
				if cs.UseCRLF {
					cs.RawData = strings.Split(string(outBytes), string([]byte{0x0D, 0x0A}))
				} else {
					cs.RawData = strings.Split(string(outBytes), string([]byte{0x0A}))
				}
				if autoOptions { // Try to identify csv parameters
					cs.autoSetOptions()
				}
				cs.Ready = true
			}
		}
	}
	return
}

// FormatRow: Formatting a Raw line to csv compliant. (usually used to preview column name)
func (cs *CsvStructure) FormatRawRow(row string) (outSlice []string, err error) {
	var csvFormatted [][]string
	bytesBuffer := &bytes.Buffer{}
	if _, err = bytesBuffer.WriteString(row); err == nil {
		csvReader := csv.NewReader(bytesBuffer)
		csvReader.LazyQuotes = cs.LazyQuotes
		csvReader.Comma = cs.Comma
		csvReader.Comment = cs.Comment
		csvReader.FieldsPerRecord = cs.FieldsPerRecord
		//		csvReader.FieldsPerRecord = cs.FieldsPerRecord
		if csvFormatted, err = csvReader.ReadAll(); err == nil {
			if len(csvFormatted) > 0 {
				outSlice = csvFormatted[0]
			}
		}
	}
	return outSlice, err
}

// Save: use only os.file for sipmple operation or buffered
// if charset must be changed from utf-8 to the one used by original file.
func (cs *CsvStructure) Save(filename string, doBackup bool) (err error) {
	if cs.Ready {
		var tmpCsvLines [][]string
		var tmpSlice []string
		var tmpDblSlice [][]string
		var csvWriter *csv.Writer
		var file *os.File
		if doBackup { // make backup if requested
			if err = os.Rename(filename, filename+"~"); err != nil {
				return err
			}
		} // Open file for writing
		if file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664); err == nil {
			defer func() {
				file.Close()
			}()
			tmpCsvLines = cs.CsvLines
			// If original csv data have some previous lines, add it.
			if cs.RowWithColName != 0 {
				tmpSlice = cs.RawData[:cs.RowWithColName]
				for _, line := range tmpSlice {
					tmpDblSlice = append(tmpDblSlice, []string{line})
				}
				tmpCsvLines = append(tmpDblSlice, tmpCsvLines...)
			}
			// Charset must be changed ?
			if cs.KeepOriginalCharset && cs.Charset != "utf-8" {
				bytesBuffer := &bytes.Buffer{}         // creates IO Writer
				csvWriter = csv.NewWriter(bytesBuffer) // creates a csv writer that uses the io buffer.
				csvWriter.Comma = cs.Comma
				csvWriter.UseCRLF = cs.UseCRLF
				if err = csvWriter.WriteAll(tmpCsvLines); err != nil { // No need to flush with WriteAll
					return err
				}
				if csvWriter.Error() != nil {
					return err
				}
				if cd, err := iconv.Open(cs.Charset, "utf-8"); err == nil { // Convert utf-8 to source file charset.
					defer cd.Close()
					iconvWriter := iconv.NewWriter(cd, file, 0, false) // Write data to file
					if _, err = fmt.Fprintln(iconvWriter, string(bytesBuffer.Bytes())); err != nil {
						return err
					}
					iconvWriter.Sync() // if autoSync = false, you need call Sync() by yourself
				}
			} else {
				cs.Charset = "utf-8"
				csvWriter = csv.NewWriter(file)
				csvWriter.Comma = cs.Comma
				csvWriter.UseCRLF = cs.UseCRLF
				err = csvWriter.WriteAll(tmpCsvLines)
			}
		}
	}
	return err
}

// AutoSetOptions: Try to find good options for csv file
func (cs *CsvStructure) autoSetOptions() (err error) {
	var colCount [][]int
	var countedCols int
	var formattedLine, tmpSliceRawData []string
	var tmpStringRawData string
	var sortSlice = func(slice [][]int) { // Sort preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool {
			return slice[i][2] > slice[j][2]
		})
	}
	var frequencyNumber = func(sl [][]int, nb int) { // Count occurence of number
		for idx, val := range sl {
			if val[1] == nb {
				colCount[idx][2]++
				break
			}

		}
	}
	// Try to find separator (comma)
	if len(cs.RawData) >= cs.analyseRowCount {
		tmpSliceRawData = cs.RawData[:cs.analyseRowCount]
	} else {
		tmpSliceRawData = cs.RawData
	}
	if cs.UseCRLF { // Keep only "cs.analyseRowCount" number of rows
		tmpStringRawData = strings.Join(tmpSliceRawData, string([]byte{0x0D, 0x0A}))
	} else {
		tmpStringRawData = strings.Join(tmpSliceRawData, string([]byte{0x0A}))
	}
	tmpComma := findCountStr(tmpStringRawData, `[^-: \r\n\(\)"'._/\\\*#[:alnum:]+]`) // Regex exclusion list
	if len(tmpComma) != 0 {
		cs.Comma = []rune(tmpComma[0].Str)[0]
	}
	// Try to find columns count
	for idx, line := range cs.RawData { // Get occurences for nbr of col/row
		if formattedLine, err = cs.FormatRawRow(line); err == nil {
			countedCols = len(formattedLine)
			if countedCols > 0 {
				colCount = append(colCount, []int{idx, countedCols, 1})
				frequencyNumber(colCount, countedCols)
			}
		}
		if idx == cs.analyseRowCount {
			break
		}
	}
	sortSlice(colCount) // sort result and get 1st of theses
	cs.FieldsPerRecord = colCount[0][1]
	cs.RowWithColName = colCount[0][0]

	return
}

// checkForConvertToUtf8: Get charset and convert it to utf-8 if needed
func (cs *CsvStructure) checkForConvertToUtf8(textBytes []byte) (outBytes []byte) {
	var strByte []byte
	var err error
	outBytes = make([]byte, len(textBytes))
	copy(outBytes, textBytes)
	_, cs.Charset, _ = charset.DetermineEncoding(textBytes, "")
	// fmt.Println(cs.Charset)
	if cs.Charset != "utf-8" { // convert to UTF-8
		if reader, err := charset.NewReader(strings.NewReader(string(textBytes)), cs.Charset); err == nil {
			if strByte, err = ioutil.ReadAll(reader); err == nil {
				outBytes = make([]byte, len(strByte))
				copy(outBytes, strByte)
			}
		}
	}
	if err != nil {
		fmt.Println("Error while trying to convert %s to utf-8: ", cs.Charset)
	}
	return outBytes
}

// Get count of non alphanum chars in a string in struct(rowStore) format
func findCountStr(str, regExpression string) []rowStore {
	storeElements := []rowStore{}
	re := regexp.MustCompile(regExpression)
	submatchall := re.FindAllString(str, -1)
	for _, element := range submatchall {
		counted := strings.Count(str, element)
		storeElements = appendIfMissing(storeElements, rowStore{counted, 0, 0, element}) // fmt.Sprintf(`%q`, element)})
	}
	// Sort slice higher to lower
	sort.Slice(storeElements, func(i, j int) bool {
		return storeElements[i].Idx > storeElements[j].Idx
	})
	return storeElements
}

// Append to slice if not already exist (rowStore structure version)
func appendIfMissing(inputSlice []rowStore, input rowStore) []rowStore {
	for _, element := range inputSlice {
		if element == input {
			return inputSlice
		}
	}
	return append(inputSlice, input)
}
