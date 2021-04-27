// md_src.go

/*
	Source file auto-generated on Mon, 12 Apr 2021 03:57:03 using Gotk3 Objects Handler v1.7.5 ©2018-21 hfmrow
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2021 hfmrow - docMark library v1.0 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package docMark

import (
	glfs "github.com/hfmrow/gen_lib/files"
	glsg "github.com/hfmrow/gen_lib/strings"
	gltsgssw "github.com/hfmrow/gen_lib/tools/goSources/sourceWalker"
)

var (
	Debug bool

	// Lib mapping
	goSourcePkgStructSetup = gltsgssw.GoSourcePkgStructSetup

	// Files
	truncatePath = glfs.TruncatePath

	// Misc
	toKebab        = glsg.ToKebab
	removeDupSpace = glsg.RemoveDupSpace
)

func DocMarkNew(root string) (*DocMark, error) {
	var err error
	dm := new(DocMark)
	dm.Root = root
	dm.gsf, err = goSourcePkgStructSetup(root)
	if err != nil {
		return nil, err
	}
	// TODO debug
	_, err = dm.gsf.AstToFileAndBBuff("/media/syndicate/storage/Documents/dev/go/src/github.com/hfmrow/go-doc-mark/ast.txt")
	if err != nil {
		return nil, err
	}

	dm.declDataNew()
	return dm, nil
}

type DocMark struct {
	Root,
	Pkge string

	Impt []string
	Decl []*declData

	gsf *gltsgssw.GoSourceFileStruct
}

type declData struct {
	Name,
	Head,
	Decl,
	Comm,
	Kind,
	File,
	Recv string

	PrmsRslt prmsRslt
	Expt     bool
	Coms     []gltsgssw.CommentStruct
}

// prmsRslt:
type prmsRslt struct {
	Params,
	Results []string
	InLine string
}

func (dm *DocMark) declDataNew() {

	dm.Pkge = dm.gsf.Package
	// Imports
	for _, imp := range dm.gsf.Imports {
		dm.Impt = append(dm.Impt, imp.Name)
	}

	// Functions
	dm.Decl = append(dm.Decl, dm.getDeclData(dm.gsf.Func)...)
	// Structures
	dm.Decl = append(dm.Decl, dm.getDeclData(dm.gsf.Struct)...)
}

// getDeclData:
func (dm *DocMark) getDeclData(g interface{}) []*declData {
	var arDeclData []*declData
	switch t := g.(type) {
	case []gltsgssw.Function:
		for _, f := range t {
			decl := new(declData)
			pr := f.Content.GetParamsResults()
			decl = &declData{
				Name: f.Ident.Name,
				Head: string(f.Content.Head),
				Decl: string(f.Content.Content),
				Coms: f.Content.Comments,
				Kind: f.Ident.Kind,
				File: truncatePath(f.File, 2),
				Recv: f.RecvFromSrc,
				PrmsRslt: prmsRslt{
					Params:  pr.Params,
					Results: pr.Results,
					InLine:  pr.InLine},
				Expt: f.Exported,
			}
			arDeclData = append(arDeclData, decl)
		}
	case []gltsgssw.Structure:
		for _, f := range t {
			decl := new(declData)
			pr := f.Content.GetParamsResults()
			decl = &declData{
				Name: f.Ident.Name,
				Head: string(f.Content.Head),
				Decl: string(f.Content.Content),
				Coms: f.Content.Comments,
				Kind: f.Ident.Kind,
				File: truncatePath(f.File, 2),
				PrmsRslt: prmsRslt{
					Params:  pr.Params,
					Results: pr.Results,
					InLine:  pr.InLine},
				Expt: f.Exported,
			}
			arDeclData = append(arDeclData, decl)
			if len(f.Methods) > 0 {
				arDeclData = append(arDeclData, dm.getDeclData(f.Methods)...)
			}
		}
	}
	return arDeclData
}
