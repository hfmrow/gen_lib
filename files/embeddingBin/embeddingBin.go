// embeddingBin.go

/*
	©2019 H.F.M
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	The source-code below is derived from the work of:
	[github.com/jteeuwen/go-bindata], his work is subject to the CC0 1.0 Universal
	(CC0 1.0) Public Domain Dedication. http://creativecommons.org/publicdomain/zero/1.0/
	which I thank the author (jteeuwen) for his great work.
*/

package embedding

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	glsg "github.com/hfmrow/gen_lib/strings"
)

const lowerHex = "0123456789abcdef"

type EmbedFileStruct struct {
	Root,
	RootOut,
	OutFilename,

	GoFileMatrix string

	Files []embedFileData

	Append,
	RecurseScan,
	VarNameLowerAtFirst,
	Verbose,
	IgnoreExt,
	StoreDataInStruct,
	StoreDataToGoFile bool

	outFile *os.File
}

type embedFileData struct {
	VarName,
	Filename,
	Path string
	Data     []byte
	IsDir    bool
	Size     int64
	FileMode os.FileMode
}

// EmbedFileStructNew: Read file(s) and convert to embedded data in
// golang source code. Output file is golang code compatible and can
// be used as []byte variable(s) in the source code. Structure contain
// all needed information about converted file(s) to re-create it as
// it has been before his conversion, this include file/directory tree.
// The stored datas are compressed with gzip.
func EmbedFileStructNew(root, outFilename string) (efs *EmbedFileStruct) {
	efs = new(EmbedFileStruct)
	efs.Root = root
	efs.OutFilename = outFilename
	efs.GoFileMatrix = newFileMatrix
	efs.RecurseScan = true
	return
}

// StoreFiles: get data from files, start at Root
func (efs *EmbedFileStruct) StoreFiles() (err error) {

	var fi os.FileInfo

	if _, err = os.Stat(efs.OutFilename); os.IsNotExist(err) || !efs.Append {
		if efs.StoreDataToGoFile {
			if err = ioutil.WriteFile(efs.OutFilename, []byte(newFileMatrix), 0644); err != nil {
				return
			}
		}
	}

	// Check for single file
	if fi, err = os.Stat(efs.Root); err != nil {
		return
	} else if !fi.IsDir() {
		filename := efs.Root
		efs.Root = filepath.Dir(efs.Root)
		err = efs.storeFile(filename, fi, err)
	} else {
		// Not a single file (dir)
		if err = filepath.Walk(efs.Root, efs.storeFile); err == nil {
			// Sort string preserving order ascendant
			efs.SortFilesListToCpy()
		}
	}
	if err == nil {
		err = efs.binFilesToHexFile()
	}
	return
}

// storeFile: get data from a file
func (efs *EmbedFileStruct) storeFile(path string, info os.FileInfo, errIn error) error {
	var (
		rel    string
		errRel error
		fd     embedFileData
		num    int

		checkExists = func(in string) bool {
			for _, file := range efs.Files {
				if file.VarName == in {
					return true
				}
			}
			return false
		}
	)

	if errIn != nil {
		return errIn
	}

	// check if we need to recurse files
	if !efs.RecurseScan {
		currRoot := filepath.Dir(path)
		if fi, err := os.Stat(path); err == nil {
			if efs.Root != currRoot || fi.IsDir() {
				return nil
			}
		}
	}

	if rel, errRel = filepath.Rel(efs.Root, path); errRel == nil && rel != "." {

		if !info.IsDir() {
			name := filepath.Base(rel)

			// Remove Ext if requested
			if efs.IgnoreExt {
				name = strings.TrimSuffix(name, filepath.Ext(name))
			}

			fd.VarName = glsg.ToCamel(name, efs.VarNameLowerAtFirst)

			for checkExists(fd.VarName) {
				num++
				fd.VarName = glsg.ToCamel(name, efs.VarNameLowerAtFirst) + fmt.Sprint(num)
			}

		} else {
			fd.IsDir = true
		}

		fd.Path = path
		fd.Size = info.Size()
		fd.Filename = rel
		fd.FileMode = info.Mode()
		efs.Files = append(efs.Files, fd)

		if efs.Verbose {
			fmt.Println("store", " < ", filepath.Join(efs.Root, fd.Filename))
		}
	}
	return errRel
}

// RestoreFiles: Write to file(s) tree structure and data contained in Json file.
// callback if exists, is called for each file.
func (efs *EmbedFileStruct) RestoreFiles(callback ...func(path string, err error) error) (err error) {

	var f func(path string, err error) error

	if len(callback) > 0 {
		f = callback[0]
	}

	for _, file := range efs.Files {

		filename := filepath.Join(efs.RootOut, file.Filename)
		if file.IsDir {
			err = os.MkdirAll(filename, file.FileMode)
		} else {
			err = ioutil.WriteFile(
				filename,
				efs.HexToBytes(file.Filename, file.Data),
				file.FileMode)
		}
		if err != nil {
			break
		}
		if efs.Verbose {
			fmt.Println("restore", " > ", filepath.Join(efs.RootOut, file.Filename))
		}
		if f != nil {
			if err = f(filename, err); err != nil {
				break
			}
		}
	}
	return
}

func (efs *EmbedFileStruct) SortFilesListToRem() {

	// Sort string preserving order descendant
	sort.SliceStable(efs.Files, func(i, j int) bool {
		return strings.ToLower(efs.Files[i].Filename) > strings.ToLower(efs.Files[j].Filename)
	})
}

func (efs *EmbedFileStruct) SortFilesListToCpy() {

	// Sort string preserving order ascendant
	sort.SliceStable(efs.Files, func(i, j int) bool {
		return strings.ToLower(efs.Files[i].Filename) < strings.ToLower(efs.Files[j].Filename)
	})
}

// binToHexString: Convert binary file to gzipped []byte data
func (efs *EmbedFileStruct) binToHexString(filename string) (out interface{}, err error) {

	// var byteToString = func(data []byte) (outString string) {
	// 	var inByte byte
	// 	buffer := []byte(`\x00`)
	// 	for _, inByte = range data {
	// 		buffer[2] = lowerHex[inByte/16]
	// 		buffer[3] = lowerHex[inByte%16]
	// 		outString += string(buffer)
	// 	}
	// 	return
	// }

	fdIn, err := os.Open(filename)
	if err != nil {
		return out, err
	}
	fi, _ := fdIn.Stat()
	buff := make([]byte, fi.Size())
	_, err = fdIn.Read(buff)
	if err != nil {
		return
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(buff)
	err = w.Close()
	if err != nil {
		return out, err
	}

	return b.Bytes(), nil
	// return byteToString(b.Bytes()), nil
}

// GetSize: retrieve files size
func (efs *EmbedFileStruct) GetSize() (size int64, err error) {

	// Save vars state
	str := efs.StoreDataInStruct
	gof := efs.StoreDataToGoFile
	efs.StoreDataInStruct = false
	efs.StoreDataToGoFile = false

	if err = efs.StoreFiles(); err == nil {
		for _, file := range efs.Files {
			size += file.Size
		}
	}

	// Restore vars state
	efs.StoreDataInStruct = str
	efs.StoreDataToGoFile = gof
	return
}

// binFilesToHexFile: Convert binary files to gzipped []byte in specific file.
// Much faster than the version below that only deals with one file at a time.
func (efs *EmbedFileStruct) binFilesToHexFile() (err error) {

	var (
		fileRead *os.File
		value    interface{}
		w        *bufio.Writer
	)

	// Create a buffered writer for better performance.
	if efs.StoreDataToGoFile {

		efs.outFile, err = os.OpenFile(efs.OutFilename, os.O_APPEND|os.O_WRONLY, 0644)
		defer efs.outFile.Close()

		w = bufio.NewWriter(efs.outFile)
		defer w.Flush()
	}

	for idx, file := range efs.Files {

		if file.IsDir {
			continue
		}

		// Store data to Go file
		if efs.StoreDataToGoFile {
			if _, err = fmt.Fprintf(w, `var %s = HexToBytes("%s", []byte("`, file.VarName, file.VarName); err == nil {

				// Read file content
				fileRead, err = os.Open(file.Path)
				if err == nil {

					// Compress and write
					gz := gzip.NewWriter(&stringWriter{Writer: w})
					_, err = io.Copy(gz, fileRead)
					gz.Close()
					if err == nil {
						_, err = fmt.Fprintf(w, `"))
`)
					} else {
						break
					}
				}
			}
		}
		// Store data to Json file
		if efs.StoreDataInStruct && err == nil {
			value, err = efs.binToHexString(file.Path)
			efs.Files[idx].Data = (value).([]byte)
		}
	}

	return
}

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func (efs *EmbedFileStruct) HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		if _, err = io.Copy(&bBuffer, r); err == nil {
			if err = r.Close(); err == nil {
				return bBuffer.Bytes()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occurred while reading: %s\n%v\n", varPath, err.Error())
	}
	return outByte
}

// GetBytesFromVarAsset: Get []byte representation from file or asset, depending on type
func (efs *EmbedFileStruct) GetBytesFromVarAsset(varPath interface{}) (outBytes []byte, err error) {
	var rBytes []byte
	switch reflect.TypeOf(varPath).String() {
	case "string":
		rBytes, err = ioutil.ReadFile(varPath.(string))
	case "[]uint8":
		rBytes = varPath.([]byte)
	}
	return rBytes, err
}

// Read: structure from file.
func (efs *EmbedFileStruct) Read(filename string) (err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(bytes, &efs)
	}
	return
}

// Write: structure to file
func (efs *EmbedFileStruct) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(&efs); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}

type stringWriter struct {
	io.Writer
	c int
}

func (w *stringWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	buf := []byte(`\x00`)
	var b byte
	for n, b = range p {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		w.Writer.Write(buf)
		w.c++
	}
	n++
	return
}

var newFileMatrix = `
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		if _, err = io.Copy(&bBuffer, r); err == nil {
			if err = r.Close(); err == nil {
				return bBuffer.Bytes()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occurred while reading: %s\n%v\n", varPath, err.Error())
	}
	return outByte
}

`
