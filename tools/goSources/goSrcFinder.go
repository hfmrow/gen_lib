// goSrcFinder.go

/*
	©2019 H.F.M. MIT license
*/

package goSources

import (
	"io/ioutil"
	"regexp"
	"strings"

	gst "github.com/hfmrow/gen_lib/strings"
)

type storeDecl struct {
	RowNb      int
	Definition string
}

type GoDeclarations struct {
	Filename       string
	rows           []string
	PackageName    string
	DataPresent    bool
	Functions      []storeDecl
	FunctNoComment []storeDecl
	FullFunc       []storeDecl
	Types          []storeDecl
	TypesNoComment []storeDecl
	Imports        []storeDecl
	FoundRows      []storeDecl
	Var            []storeDecl
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) GoSourceGetLines(filename string, wholeWord bool, terms ...string) (exist bool, err error) {
	var ww string
	var notMatch bool
	var regs []regexp.Regexp

	if err = d.readFile(filename); err == nil {
		if wholeWord {
			ww = "\b"
		}

		for _, term := range terms {
			term = regexp.QuoteMeta(term)
			regs = append(regs, *regexp.MustCompile(`(` + ww + term + ww + `)`))
		}

		for idxRow, row := range d.rows {
			for _, reg := range regs {
				if !reg.MatchString(row) {
					notMatch = true
				}
			}
			if !notMatch {
				d.FoundRows = append(d.FoundRows, storeDecl{idxRow, row})
				exist = true
			}
			notMatch = false
		}
	}
	return exist, err
}

// readFile: read file and make slice of strings
func (d *GoDeclarations) readFile(filename string) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err == nil {
		d.rows = strings.Split(string(data), gst.GetTextEOL(data))
	}
	return err
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) GoSourceGetInfos(filename string, funcName ...string) (err error) {
	var toFind string
	switch len(funcName) {
	case 1:
		toFind = funcName[0]
	}
	if err = d.readFile(filename); err == nil {
		d.Types = d.getDecl("type", "}", "", true, toFind)
		d.TypesNoComment = d.getDecl("type", "}", "", false, toFind)
		d.FullFunc = d.getDecl("func", "}", "", false, toFind)
		d.Functions = d.getDecl("func", "", "{", true, toFind)
		d.FunctNoComment = d.getDecl("func", "", "{", false, toFind)
		d.getImports()
		d.getPackageName()
		d.Var = d.getDecl("var", "", "", false, toFind)
	}
	return err
}

// getPackageName: retrieve package name of Go source file
func (d *GoDeclarations) getPackageName() (err error) {
	var pkgReg *regexp.Regexp
	if pkgReg, err = regexp.Compile(`^(\bpackage\b)`); err == nil {
		for idxRow := 0; idxRow < len(d.rows); idxRow++ {
			if pkgReg.MatchString(d.rows[idxRow]) {
				d.PackageName = strings.Split(d.rows[idxRow], " ")[1]
				break
			}
		}
	}
	return err
}

// getImports: Scan go source file and return all requested imports.
func (d *GoDeclarations) getImports() {
	var tempStrings []string

	importReg := regexp.MustCompile(`^(import .*)`)
	startMulti := regexp.MustCompile(`(.*\()$`)
	endMulti := regexp.MustCompile(`^(\))`)

	for idxRow := 0; idxRow < len(d.rows); idxRow++ {
		if importReg.MatchString(d.rows[idxRow]) {
			if !startMulti.MatchString(d.rows[idxRow]) {
				tempStrings = append(tempStrings, d.rows[idxRow])
				break
			}
			for !endMulti.MatchString(d.rows[idxRow]) {
				idxRow++
			}
			for idx := idxRow; idx >= 0; idx-- {
				if len(d.rows[idx]) != 0 {
					tempStrings = append(tempStrings, d.rows[idx])
				} else {
					tempStrings = append(tempStrings, "")
					break
				}
			}
		}
	}
	for idxRow := len(tempStrings) - 1; idxRow >= 0; idxRow-- {
		d.Imports = append(d.Imports, storeDecl{idxRow, tempStrings[idxRow]})
	}
}

// getDecl: Scan go source file and return declarations.
func (d *GoDeclarations) getDecl(
	startDecl,
	endDeclAtStart,
	endDeclAtEnd string,
	wantComments bool,
	funcName ...string) (tmpStore []storeDecl) {

	var stop bool
	var toFind string

	endDeclAtStart = regexp.QuoteMeta(endDeclAtStart)
	endDeclAtEnd = regexp.QuoteMeta(endDeclAtEnd)

	if len(funcName) != 0 {
		toFind = funcName[0]
	}

	toFindReg := regexp.MustCompile(`(\b` + toFind + `\b)`)
	startDeclReg := regexp.MustCompile(`^(\b` + startDecl + `\b)`)
	endDeclReg := regexp.MustCompile(`^(` + endDeclAtStart + `)`)
	if len(endDeclAtEnd) != 0 {
		endDeclReg = regexp.MustCompile(`(` + endDeclAtEnd + `)$`)
	}

	for idxRow := 0; idxRow < len(d.rows); idxRow++ {

		if startDeclReg.MatchString(d.rows[idxRow]) && toFindReg.MatchString(d.rows[idxRow]) {

			for len(d.rows[idxRow]) != 0 && wantComments && idxRow != 0 { //Go back to get comments
				idxRow--
			}

			if endDeclReg.MatchString(d.rows[idxRow]) {
				tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
				break
			} else {
				for !endDeclReg.MatchString(d.rows[idxRow]) {
					tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
					idxRow++
					stop = true
				}
			}
		}
		if stop {
			tmpStore = append(tmpStore, storeDecl{idxRow, d.rows[idxRow]})
			tmpStore = append(tmpStore, storeDecl{-1, ""})

			stop = false
		}
	}
	if len(tmpStore) != 0 {
		d.DataPresent = true
	}
	return tmpStore
}
