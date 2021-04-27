// embeddingFiles.go

/*
	Source file auto-generated on Sat, 07 Nov 2020 06:45:03 using Gotk3ObjHandler v1.6.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2020 H.F.M - github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

/*
	This library is designed to write directory tree with binary file(s) as json text format.
	This way, that permit to transfer/copy and reorganize files tree and use them as embedded data.
*/

package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FilesListEmbedding struct {
	Root,
	RootOut string
	Files []fileData
	Verbose,
	GetData bool
}

type fileData struct {
	Filename string
	Data     []byte
	IsDir    bool
	Size     int64
	FileMode os.FileMode
}

// FilesListEmbeddingStructNew: Read file(s) and store tree and content
// To a Json file. Structure contain all needed information about
// converted file(s) to re-create it as it has been before his conversion,
// this include file/directory tree.
func FilesListEmbeddingStructNew(root string) (fls *FilesListEmbedding) {
	fls = new(FilesListEmbedding)
	fls.Root = root
	fls.GetData = true
	return fls
}

func (fls *FilesListEmbedding) GetSize() (size int64, err error) {
	tmpGetData := fls.GetData
	fls.GetData = false
	if err = fls.StoreFiles(); err == nil {
		for _, file := range fls.Files {
			size += file.Size
		}
	}
	fls.GetData = tmpGetData
	return
}

func (fls *FilesListEmbedding) StoreFile(path string, info os.FileInfo, errIn error) error {
	var (
		rel    string
		errRel error
		fd     fileData
	)
	if errIn != nil {
		return errIn
	}
	if rel, errRel = filepath.Rel(fls.Root, path); errRel == nil && rel != "." {
		if !info.IsDir() && fls.GetData {

			fd.Data, errRel = ioutil.ReadFile(path)
		} else {
			fd.IsDir = true
		}
		fd.Size = info.Size()
		fd.Filename = rel
		fd.FileMode = info.Mode()
		fls.Files = append(fls.Files, fd)

		if fls.Verbose {
			fmt.Println("store", " < ", filepath.Join(fls.Root, fd.Filename))
		}
	}
	return errRel
}

func (fls *FilesListEmbedding) SortFilesListToRem() {

	// Sort string preserving order descendant
	sort.SliceStable(fls.Files, func(i, j int) bool {
		return strings.ToLower(fls.Files[i].Filename) > strings.ToLower(fls.Files[j].Filename)
	})
}

func (fls *FilesListEmbedding) SortFilesListToCpy() {

	// Sort string preserving order ascendant
	sort.SliceStable(fls.Files, func(i, j int) bool {
		return strings.ToLower(fls.Files[i].Filename) < strings.ToLower(fls.Files[j].Filename)
	})
}

func (fls *FilesListEmbedding) StoreFiles() (err error) {

	var walkFunc filepath.WalkFunc
	walkFunc = fls.StoreFile

	err = filepath.Walk(fls.Root, walkFunc)
	return
}

func (fls *FilesListEmbedding) RestoreFiles(callback ...func(path string, err error) error) (err error) {

	var filename string
	for _, file := range fls.Files {

		filename = filepath.Join(fls.RootOut, file.Filename)
		if file.IsDir {
			err = os.MkdirAll(filename, file.FileMode)
		} else {
			err = ioutil.WriteFile(filename, file.Data, file.FileMode)
		}
		if err != nil {
			break
		}
		if fls.Verbose {
			fmt.Println("restore", " > ", filepath.Join(fls.RootOut, file.Filename))
		}
		if len(callback) > 0 {
			if err = callback[0](filename, err); err != nil {
				break
			}
		}
	}
	return
}

// Read: Options from file.
func (fls *FilesListEmbedding) Read(filename string) (err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(bytes, &fls)
	}
	return
}

// Write: Options to file
func (fls *FilesListEmbedding) Write(filename string) (err error) {
	var jsonData []byte
	var out bytes.Buffer
	if jsonData, err = json.Marshal(&fls); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return
}
