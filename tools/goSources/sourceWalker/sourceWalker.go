// sourceWalker.go

/*
	Copyright ©2018-21 hfmrow - sourceWalker library v1.2 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This package allows you to navigate inside go source code (package or file (s)), and
	to retrieve information (exported or not) on contained functions, methods, structure,
	variables, imports, comments ... All this information is stored in a single structure
	that contains methods to manage them.
*/

package sourceWalker

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	glfsfs "github.com/hfmrow/gen_lib/files/filesOperations"
	glss "github.com/hfmrow/gen_lib/slices"
	glsg "github.com/hfmrow/gen_lib/strings"
	gltsgsge "github.com/hfmrow/gen_lib/tools/goSources/goEnvGet"
)

var (
	// Lib mapping
	gsExistSlIface   = glss.IsExistSlIface
	filesOpStructNew = glfsfs.FilesOpStructNew
	getTextEOL       = glsg.GetTextEOL
	removeDupSpace   = glsg.RemoveDupSpace
	GetGoEnv         = gltsgsge.GetGoEnv
)

var KeyWordsList = []string{
	"break",
	"default",
	"func",
	"select",
	"case",
	"defer",
	"go",
	"struct",
	"else",
	"goto",
	"package",
	"switch",
	"const",
	"fallthrough",
	"if",
	"range",
	"type",
	"continue",
	"for",
	"import",
	"return",
	"var",
	"append",
	"cap",
	"close",
	"copy",
	"delete",
	"imag",
	"len",
	"make",
	"new",
	"panic",
	"recover",
}
var TypesList = []string{
	"chan",
	"map",
	"interface",
	"uint",
	"int",
	"uintptr",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"float32",
	"float64",
	"complex64",
	"complex128",
	"byte",
	"rune",
	"true",
	"false",
	"bool",
	"error",
	"string",
	"iota",
	"nil",
	"real",
	"complex",
	// "Bool",
	// "Int",
	// "Int8",
	// "Int16",
	// "Int32",
	// "Int64",
	// "Uint",
	// "Uint8",
	// "Uint16",
	// "Uint32",
	// "Uint64",
	// "Uintptr",
	// "Float32",
	// "Float64",
	// "Complex64",
	// "Complex128",
	// "String",
	// "UnsafePointer",
}

// GoSourceFileStruct: contain AST file information
type GoSourceFileStruct struct {
	Filename,
	Eol, // End of line used in the input file
	Package,
	GoSourcePath string

	PackageLineIdx int

	Imports []imported
	Func    []Function
	Struct  []Structure
	Var     []Variable

	// Unexported
	data          []byte        // File content
	astOut        *bytes.Buffer // AST representation of the input file
	linesIndexes  [][]int
	offset        int // Define if we start at 0 or 1  when counting lines and offsets positions.
	astFile       *ast.File
	fset          *token.FileSet // Positions are relative to fset.
	tmpMethods    []Function
	varInsideFunc bool

	fos *glfsfs.FilesOpStruct
}

type imported struct {
	Name,
	NameFromSrc,
	File string

	Content content
}

type Function struct {
	Ident   identObj
	Content content
	File,
	Name,
	Recv,
	RecvFromSrc,
	Body string

	Exported   bool
	ParamReslt paramReslt
}

type paramReslt struct {
	Params,
	Result []string
}

type Structure struct {
	Ident   identObj
	Content content
	Fields  []field
	Methods []Function

	File,
	Name string

	Exported bool
}

type Variable struct {
	Objects  field
	Content  content
	File     string
	Found    identObj
	Exported bool
}

type field struct {
	List []identObj

	Type,
	Name string
}

type identObj struct {
	Name,
	Kind, // Struc or Func or Var
	Type,
	Value string

	Idx      int
	Exported bool
}

type content struct {
	OfstStart,
	OfstEnd,
	LineStart,
	LineEnd,
	LBrace,
	RBrace int
	// Full declaration
	Content []byte
	// Simplified function declaration (params and results
	Head,
	// Simplified version (all concerned comment are packed in a string)
	Comment string
	// Contain parameters en results arguments.
	PrmsRslt paramsResults
	// Detailled comments version (in separateed lines format)
	Comments []CommentStruct

	eol string
}

type CommentStruct struct {
	// Raw, direct from source.
	Text string
	// stored in separated lines
	Lines []string
	// Generally the comment is contained between '/ *' and '* /'
	IsMultiLines bool
}

func GoSourceFileStructNew() (gsfs *GoSourceFileStruct, err error) {

	gsfs = new(GoSourceFileStruct)
	gsfs.GoSourcePath, err = GetGoEnv("GOPATH")
	if err != nil {
		return nil, fmt.Errorf("GOPATH or GOROOT environment variables seem not to be set correctly. Operation will not be done.")
	}
	gsfs.GoSourcePath = filepath.Join(gsfs.GoSourcePath, "src")
	gsfs.fos, err = filesOpStructNew()
	return
}

// GoSourcePkgStructSetup: Recover all data from package files.
func GoSourcePkgStructSetup(path string) (*GoSourceFileStruct, error) {

	var currPkg string

	gsfs, err := GoSourceFileStructNew()
	if err != nil {
		return nil, err
	}
	gsfs.fos.Masks = []string{"*.go"}
	err = gsfs.fos.GetFilesDetails(path)
	if err != nil {
		return nil, err
	}
	for idx, fdtl := range gsfs.fos.Files {
		gsfs.tmpMethods = []Function{}
		if !fdtl.IsDir {
			if idx == 0 {
				err = gsfs.GoSourceFileStructureSetup(fdtl.FilenameFull)
				if err != nil {
					return nil, err
				}
				// Store pkg name on first file
				currPkg = gsfs.Package
				continue
			}
			if currPkg == gsfs.Package {
				err = gsfs.AppendFile(fdtl.FilenameFull)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return gsfs, nil
}

// GetFuncByName: Optional "unExported": empty = both.
func (gsfs *GoSourceFileStruct) GetFuncByName(fName string, unExported ...bool) (funct *Function) {

	for idx, fnc := range gsfs.Func {
		if fnc.Ident.Name == fName {
			return &gsfs.Func[idx]
		}
	}
	return
}

// GetStructByName:
func (gsfs *GoSourceFileStruct) GetStructByName(sName string) (stru *Structure) {

	for idx, stc := range gsfs.Struct {
		if stc.Ident.Name == sName {
			return &gsfs.Struct[idx]
		}
	}
	return
}

// GetVarByName: "Position" contain the position in "list" and "values" fields.
func (gsfs *GoSourceFileStruct) GetVarByName(vName string) (vari *Variable) {

	for idx, vr := range gsfs.Var {
		for idn, vn := range vr.Objects.List {
			if vn.Name == vName {
				tmpV := gsfs.Var[idx]
				tmpV.Found.Name = vn.Name
				tmpV.Found.Type = vr.Objects.Type
				tmpV.Found.Value = vn.Value
				tmpV.Found.Idx = idn
				return &tmpV
			}
		}
	}
	return
}

// varWalker: called from the walker or others to check for variables tokens
func varWalker(t token.Token) (kind string, ok bool) {
	switch t {
	case token.VAR:
		return "var", true
	case token.CONST:
		return "const", true
	case token.ASSIGN:
		return "=", true
	case token.DEFINE:
		return ":=", true
	}
	return
}

// GetImportsOnly: translation to the internal function, designed fo kick use.
func (gsfs *GoSourceFileStruct) GetImportsOnly(filename string) (err error) {
	gsfs.Filename = filename
	if err = gsfs.loadDataFile(); err == nil {
		err = gsfs.getImportOnly()
	}
	return
}

// getImportOnly: Retrieve all 'imports' contained in current file
func (gsfs *GoSourceFileStruct) getImportOnly() (err error) {

	var exist bool

	if filename, err := filepath.Rel(gsfs.GoSourcePath, gsfs.Filename); err == nil {

		for _, val := range gsfs.astFile.Imports {
			content := gsfs.getContentFromPos(val.Pos(), val.End())

			i := imported{
				Name:        val.Path.Value,
				NameFromSrc: string(gsfs.data[val.Pos()-1 : val.End()-1]),
				Content:     content,
				File:        filename}

			// Avoiding duplicate entry
			for _, imp := range gsfs.Imports {
				if imp.NameFromSrc == i.NameFromSrc {
					exist = true
					break
				}
			}

			if !exist {
				gsfs.Imports = append(gsfs.Imports, i)
			} else {
				exist = false
			}
		}
	}
	return
}
func (gsfs *GoSourceFileStruct) getSelector(xv *ast.StarExpr) (string, error) {
	switch xx := xv.X.(type) {
	case *ast.SelectorExpr:

		return "*" + xx.X.(*ast.Ident).Name + "." +
			xx.Sel.Name, nil
	case *ast.Ident:

		return "*" + xx.Name, nil
	case *ast.InterfaceType:

		return "*" + "interface{}", nil
	case *ast.ArrayType:

		return "*" + gsfs.getArray(xx), nil
	case *ast.StarExpr:

		return gsfs.getSelector(xx)
	}
	return "", fmt.Errorf("getParamReslt/StarExpr/SelectorExpr: %T, not handled actually\n", xv.X)
}

// TODO Use AST to get params and results WIP
func (gsfs *GoSourceFileStruct) getParamReslt(val *ast.FuncDecl) paramReslt {
	var (
		pr paramReslt
		// tmpSelName  string
		// tmpIdenName string
		// err         error
	)

	// // Parameters
	// for _, v := range val.Type.Params.List {
	// 	tmpIdenName = ""
	// 	tmpSelName = ""
	// 	for _, n := range v.Names {
	// 		tmpIdenName += n.Name
	// 		fmt.Printf("name: %s\n", tmpIdenName)
	// 	}

	// 	switch xv := v.Type.(type) {
	// 	case *ast.StarExpr:
	// 		tmpSelName, err = gsfs.getSelector(xv)
	// 		if err != nil {
	// 			log.Printf("%v", err)
	// 		}
	// 	case *ast.Ident:
	// 		tmpSelName = xv.Name
	// 	case *ast.ArrayType:
	// 		tmpSelName = gsfs.getArray(xv)
	// 	case *ast.InterfaceType:
	// 		tmpSelName = "interface{}"
	// 	case *ast.MapType:
	// 		k := xv.Key.(*ast.Ident).Name
	// 		v := xv.Value.(*ast.Ident).Name
	// 		tmpSelName = "map[" + k + "]" + v
	// 	case *ast.ChanType:
	// 		tmpSelName = "chan"
	// 	default:
	// 		log.Printf("getParamReslt/StarExpr: %T, not handled actually\n", v.Type)
	// 	}
	// 	fmt.Printf("name: %s\n", tmpSelName)
	// }

	return pr
}

// goInspect: parse go file and retrieve into structure that was found.
func (gsfs *GoSourceFileStruct) goInspect() {
	filename, _ := filepath.Rel(gsfs.GoSourcePath, gsfs.Filename)
	gsfs.getImportOnly()
	// Get nodes infos
	ast.Inspect(gsfs.astFile, func(node ast.Node) bool {
		switch val := node.(type) {
		case *ast.FuncDecl: // Functions
			exported := val.Name.IsExported()
			content := gsfs.getContentFromPos(val.Pos(), val.End(), val.Doc, int(val.Body.Lbrace), int(val.Body.Rbrace))

			if val.Recv == nil {
				// Functions
				var funct Function
				obj := gsfs.getIdent(val.Name)
				funct.ParamReslt = gsfs.getParamReslt(val)
				funct.Name = val.Name.Name
				funct.Ident.Name = obj.Name
				funct.Ident.Kind = obj.Kind
				funct.Content = content
				funct.File = filename
				funct.Exported = exported
				gsfs.Func = append(gsfs.Func, funct)
			} else {
				// Methods
				var method Function
				method.ParamReslt = gsfs.getParamReslt(val)
				method.Ident.Name = gsfs.getIdent(val.Name).Name
				// Retrieve reciever
				method.RecvFromSrc = string(gsfs.data[val.Recv.Opening : val.Recv.Closing-1])
				for _, v := range val.Recv.List {
					switch xv := v.Type.(type) {
					case *ast.StarExpr:
						if si, ok := xv.X.(*ast.Ident); ok {
							method.Recv = "*" + si.Name
						}
					case *ast.Ident:
						method.Recv = xv.Name
					}
				}
				method.Body = string(gsfs.data[val.Body.Lbrace : val.Body.Rbrace-1])
				method.Ident.Type = gsfs.getFields(val.Recv.List).Ident.Type
				method.Content = content
				method.File = filename
				method.Exported = exported
				gsfs.tmpMethods = append(gsfs.tmpMethods, method)
			}
		case *ast.GenDecl:
			for _, spec := range val.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					exported := s.Name.IsExported()
					stru := gsfs.getStruct(s)
					stru.Name = s.Name.Name
					stru.Content = gsfs.getContentFromPos(val.Pos(), val.End(), val.Doc)
					stru.Ident.Type = ""
					stru.File = filename
					stru.Exported = exported
					gsfs.Struct = append(gsfs.Struct, stru)

				case *ast.ValueSpec:
					if _, ok := varWalker(val.Tok); ok && !gsfs.varInsideFunc {
						if fld := gsfs.getSpecs([]ast.Spec{s}); fld != nil {
							gsfs.Var = append(gsfs.Var, Variable{
								Objects: *fld,
								File:    filename,
								Content: gsfs.getContentFromPos(val.Pos(), val.End(), s.Doc)})
						}
					}
				}
			}
		}
		// rStmt, ok := node.(*ast.Ident)
		// if ok {
		// 	fmt.Printf("return statement found on line %d \t ret: %v\n", gsfs.fset.Position(rStmt.Pos()).Line, rStmt)

		// }
		return true
	})
}

// getIdent:
func (gsfs *GoSourceFileStruct) getIdent(ident *ast.Ident) (obj identObj) {
	obj.Name = ident.Name
	if ident.Obj != nil {
		obj.Kind = ident.Obj.Kind.String()
		if ident.Obj.Type != nil {
			obj.Type = ident.Obj.Type.(*ast.Ident).Name
		}
		if ident.Obj.Data != nil {
			switch iod := ident.Obj.Data.(type) {
			case *ast.BasicLit:
				obj.Value = iod.Value
			}
		}
	}
	return
}

// getBasicLit:
func (gsfs *GoSourceFileStruct) getBasicLit(bl *ast.BasicLit) (outValue, outType string) {
	return bl.Value, strings.ToLower(bl.Kind.String())
}

// getAssignStmt:
func (gsfs *GoSourceFileStruct) getAssignStmt(aStmt *ast.AssignStmt) (fld *field) {

	var obj identObj
	fld = new(field)
	// fld.Type = aStmt.Tok.String()
	for _, lhs := range aStmt.Lhs {
		switch lhst := lhs.(type) {
		case *ast.Ident:
			obj = gsfs.getIdent(lhst)
			obj.Kind = aStmt.Tok.String()
			obj.Type = "func"
		}
		fld.List = append(fld.List, obj)
	}
	for idx, rhs := range aStmt.Rhs {
		switch rhst := rhs.(type) {
		case *ast.BasicLit:
			fld.List[idx].Value, fld.Type = gsfs.getBasicLit(rhst)

		case *ast.Ident:
			fld.List[idx].Value = gsfs.getIdent(rhst).Name
			if fld.List[idx].Value == "true" || fld.List[idx].Value == "false" {
				fld.Type = "bool"
			}
		}
	}

	return
}

// getSpecs:
func (gsfs *GoSourceFileStruct) getSpecs(specs []ast.Spec) (fld *field) {
	for _, spec := range specs {
		var tmpStr string
		switch s := spec.(type) {
		case *ast.ValueSpec:
			fld = new(field)
			for _, idnt := range s.Names {
				obj := gsfs.getIdent(idnt)
				obj.Exported = idnt.IsExported()
				fld.List = append(fld.List, obj)
			}
			if s.Type != nil {
				switch st := s.Type.(type) {
				case *ast.Ident:
					fld.Type = gsfs.getIdent(st).Name
				case *ast.StarExpr:
					if st.X != nil {
						fld.Type = gsfs.getStarX(st.X)
					}
				case *ast.ArrayType:
					fld.Type = gsfs.getArray(st)
				}
			}
			var values []string
			if s.Values != nil {
				for _, value := range s.Values {
					switch v := value.(type) {
					case *ast.FuncLit: // Right arg is a function
						fld.Type = "func"
					case *ast.Ident:
						fld.Type = gsfs.getIdent(v).Kind
						fld.Name = gsfs.getIdent(v).Name
					case *ast.SelectorExpr:
						if v.X != nil {
							switch vXt := v.X.(type) {
							case *ast.Ident:
								tmpStr = gsfs.getIdent(vXt).Name
							}
						}
						if v.Sel != nil {
							tmpStr += "." + gsfs.getIdent(v.Sel).Name
						}
						fld.Type = tmpStr
					case *ast.BasicLit:
						fld.Type = strings.ToLower(v.Kind.String())
						values = append(values, v.Value)
					case *ast.CallExpr:
						if v.Fun != nil {
							switch vt := v.Fun.(type) {
							case *ast.Ident:
								fld.Type = gsfs.getIdent(vt).Name
							}
						}
						if v.Args != nil {
							for _, arg := range v.Args {
								switch va := arg.(type) {
								case *ast.BasicLit:
									values = append(values, va.Value)
								}
							}
						}
					}
				}
			}
			// Fill with values
			for idx := len(fld.List) - 1; idx >= 0; idx-- {
				if len(values) > idx {
					fld.List[idx].Value = values[idx]
				}
			}
			if len(fld.List) == 0 {
				fld = nil
			}
		case *ast.TypeSpec: // Struct
			// (not implemented here. it was done above)
		}
	}
	return
}

// getArray:
func (gsfs *GoSourceFileStruct) getArray(ary *ast.ArrayType, name ...string) (fld string) {

	// TODO work right, and give good count of '[]', just missing name and StarExp
	// var (
	// 	ok bool
	// )
	// fld = "[]"
	// ary, ok = ary.Elt.(*ast.ArrayType)
	// for ok {
	// 	fld += "[]"
	// 	ary, ok = ary.Elt.(*ast.ArrayType)
	// }
	// if len(name) > 0 {
	// 	fld += name[0]
	// }

	switch fvv := ary.Elt.(type) {
	case *ast.ArrayType:

		switch fvvE := fvv.Elt.(type) {
		case *ast.ArrayType:
			switch fvvF := fvvE.Elt.(type) {
			case *ast.Ident:
				fld = "[][][]" + fvvF.Name
			case *ast.StarExpr:
				fld = "[][][]" + gsfs.getStarX(fvvF.X)
			}
		case *ast.Ident:
			fld = "[][]" + fvvE.Name
		case *ast.StarExpr:
			fld = "[][]" + gsfs.getStarX(fvvE.X)
		}
	case *ast.Ident:
		fld = "[]" + fvv.Name
	case *ast.StarExpr:
		fld = "[]" + gsfs.getStarX(fvv.X)
	}

	// switch fvv := ary.Elt.(type) {
	// case *ast.ArrayType:
	// 	switch fvvE := fvv.Elt.(type) {
	// 	case *ast.ArrayType:
	// 		switch fvvF := fvvE.Elt.(type) {
	// 		case *ast.Ident:
	// 			fld = "[][][]" + fvvF.Name
	// 		case *ast.StarExpr:
	// 			fld = "[][][]" + gsfs.getStarX(fvvF.X)
	// 		}
	// 	case *ast.Ident:
	// 		fld = "[][]" + fvvE.Name
	// 	case *ast.StarExpr:
	// 		fld = "[][]" + gsfs.getStarX(fvvE.X)
	// 	}
	// case *ast.Ident:
	// 	fld = "[]" + fvv.Name
	// case *ast.StarExpr:
	// 	fld = "[]" + gsfs.getStarX(fvv.X)
	// }

	return
}

// getStarX:
func (gsfs *GoSourceFileStruct) getStarX(sX ast.Expr) (fld string) {

	switch fvX := sX.(type) {
	case *ast.Ident:
		fld = "*" + fvX.Name
	case *ast.ArrayType:
		fld = "*" + gsfs.getArray(fvX)
	}
	return
}

// getField:
func (gsfs *GoSourceFileStruct) getFields(fields []*ast.Field) (stru Structure) {
	for _, fList := range fields {
		var fld field
		for _, ident := range fList.Names {
			fld.List = append(fld.List, gsfs.getIdent(ident))
		}
		switch fv := fList.Type.(type) {
		case *ast.Ident:
			fld.Type = fv.Name
		case *ast.StarExpr:
			fld.Type = gsfs.getStarX(fv.X)
		case *ast.ArrayType:
			fld.Type = gsfs.getArray(fv)
		}
		stru.Ident.Type = fld.Type
		stru.Fields = append(stru.Fields, fld)
	}
	return
}

// getStruct:
func (gsfs *GoSourceFileStruct) getStruct(s *ast.TypeSpec) (stru Structure) {
	switch st := s.Type.(type) {
	case *ast.StructType:
		stru = gsfs.getFields(st.Fields.List)
	}
	// stru = gsfs.getFields(s.Type.(*ast.StructType).Fields.List)
	obj := gsfs.getIdent(s.Name)
	stru.Ident.Name = obj.Name
	switch t := s.Type.(type) {
	case *ast.ArrayType:
		stru.Ident.Kind = gsfs.getArray(t, obj.Name)
	case *ast.InterfaceType:
		stru.Ident.Kind = "interface{}"
	case *ast.MapType:
		kName, vName := "n\a", "n\a"
		k, ok := t.Key.(*ast.Ident)
		if ok {
			kName = k.Name
		}
		// TODO find a way to make it working right with interface{}
		v, ok := t.Value.(*ast.Ident)
		if ok {
			vName = v.Name
		}
		stru.Ident.Kind = "map[" + kName + "]" + vName
	case *ast.ChanType:
		stru.Ident.Kind = "chan"
	default:
		stru.Ident.Kind = obj.Kind
	}

	return
}

// GoSourceFileStructureSetup: setup and retieve information for designed file.
// Notice: the lines numbers and offsets start at 0. Set "zero" at false to start at 1.
func (gsfs *GoSourceFileStruct) GoSourceFileStructureSetup(filename string, zero ...bool) (err error) {
	gsfs.offset = 1 // lines start at 0 (substract 1 for each offsets position)
	if len(zero) > 0 {
		if zero[0] {
			gsfs.offset = 0
		}
	}
	gsfs.Filename = filename
	if err = gsfs.loadDataFile(); err == nil {
		if err = gsfs.fillDeclaration(); err == nil {
			gsfs.filteringMethods()
		}
	}
	return
}

// AppendFile:
func (gsfs *GoSourceFileStruct) AppendFile(filename string) (err error) {
	gsfs.tmpMethods = []Function{}
	return gsfs.GoSourceFileStructureSetup(filename)
}

func (gsfs *GoSourceFileStruct) loadDataFile() (err error) {
	// Loading data (file)
	gsfs.fset = token.NewFileSet()
	if gsfs.astFile, err = parser.ParseFile(gsfs.fset, gsfs.Filename, nil, parser.ParseComments); err == nil {
		gsfs.data, err = ioutil.ReadFile(gsfs.Filename)
	}
	return
}

// fillDeclaration:
func (gsfs *GoSourceFileStruct) fillDeclaration() (err error) {
	// Setting internal variables
	gsfs.Eol = getTextEOL(gsfs.data)
	gsfs.data = append(gsfs.data, []byte(gsfs.Eol)...) // Add an eol to avoid a f..k..g issue where the last line wasn't analysed
	eolRegx := regexp.MustCompile(gsfs.Eol)
	eolPositions := eolRegx.FindAllIndex(gsfs.data, -1)

	// Define and prepare slice of line indexes
	gsfs.linesIndexes = make([][]int, len(eolPositions)+1)
	gsfs.linesIndexes[0] = []int{0, eolPositions[0][0]}
	// Creating lines indexes
	for idx := 1; idx < len(eolPositions); idx++ {
		gsfs.linesIndexes[idx] = []int{eolPositions[idx-1][1], eolPositions[idx][0]}
	}
	// get package name
	gsfs.Package = gsfs.astFile.Name.String()
	gsfs.PackageLineIdx, _ = gsfs.getLineFromOffsets(int(gsfs.astFile.Package), int(gsfs.astFile.Package))
	// Inspecting
	gsfs.goInspect()
	return
}

// getLineFromOffsets: get the line number corresponding to offsets. Notice, line number start at 0.
func (gsfs *GoSourceFileStruct) getLineFromOffsets(sOfst, eOfst int) (lStart, lEnd int) {
	for lineNb, lineIdxs := range gsfs.linesIndexes {
		switch {
		case sOfst >= lineIdxs[0] && sOfst <= lineIdxs[1]:
			lStart = lineNb
			if eOfst <= lineIdxs[1] { // only one line
				lEnd = lineNb
				return
			}
		case eOfst >= lineIdxs[0] && eOfst <= lineIdxs[1]:
			lEnd = lineNb
			return
		}
	}
	return
}

// getContentFromPos: fill content structure
func (gsfs *GoSourceFileStruct) getContentFromPos(pos, end token.Pos, comment ...interface{}) (cnt content) {
	// Set to relative offset
	sOfst := gsfs.fset.PositionFor(pos, true).Offset
	eOfst := gsfs.fset.PositionFor(end, true).Offset
	// Make content structure
	cnt.OfstStart, cnt.OfstEnd = sOfst-gsfs.offset, eOfst-gsfs.offset
	cnt.LineStart, cnt.LineEnd = gsfs.getLineFromOffsets(sOfst-1, eOfst-1)
	cnt.Content = gsfs.data[sOfst-1 : eOfst]
	cnt.eol = gsfs.Eol

	if len(comment) > 0 {

		// Detailled comments version (in separateed lines format)
		cnt.Comments = gsfs.getComments(comment[0].(*ast.CommentGroup))

		// Simplified version (all concerned comment are packed in a string)
		if len(cnt.Comments) > 0 {
			for idx, c := range cnt.Comments {
				if c.IsMultiLines {
					cnt.Comment += strings.Join(c.Lines, gsfs.Eol)
				} else {
					cnt.Comment += c.Lines[0]
					if len(cnt.Comments)-1 != idx {
						cnt.Comment += gsfs.Eol
					}
				}
			}
		}
	}
	if len(comment) > 1 {
		cnt.LBrace = comment[1].(int)
	}
	if len(comment) > 2 {
		cnt.RBrace = comment[2].(int)
		cnt.Head = string(cnt.Content[:(cnt.LBrace-cnt.OfstStart)-2])
	}
	return
}

func (gsfs *GoSourceFileStruct) GetGlobalComments() (cs []CommentStruct) {

	for _, cg := range gsfs.astFile.Comments {
		cs = append(cs, gsfs.getComments(cg)...)
	}
	return
}

func (gsfs *GoSourceFileStruct) GetBuildConstraints() (bDir []string) {

	// Build Constraints regexp
	buildReg := regexp.MustCompile(`(?m)^(\/\/\s\+build\s)|^(\/\/go:)`)
	for _, cg := range gsfs.astFile.Comments {
		cmt := gsfs.getComments(cg)
		for _, c := range cmt {
			if buildReg.MatchString(c.Text) {
				bDir = append(bDir, c.Text)
			}
		}
	}
	return
}

// getComments: retrieve comments structure from 'CommentGroup'.
func (gsfs *GoSourceFileStruct) getComments(cg *ast.CommentGroup) []CommentStruct {
	var (
		cmnts            []CommentStruct
		reRemAllCmtMarks = regexp.MustCompile(`(?m)^(\s{0,1}\*{1,}\/\s{0,1})|(\/\*{1,})|(\*{1,}\s{0,1}\/)|(\*{1,}\s{0,1})|^(\/{2,}\s{0,1})`)
	)
	if cg == nil {
		return cmnts
	}
	for _, c := range cg.List {

		var cmnt CommentStruct

		if strings.Contains(c.Text, gsfs.Eol) {
			cmnt.IsMultiLines = true
			tStr := fmt.Sprintf("%s", c.Text)
			splittedLines := strings.Split(tStr, gsfs.Eol)
			for idx, l := range splittedLines {
				if idx == 0 || idx == len(splittedLines)-1 && len(l) == 0 {
					continue
				}
				cmnt.Lines = append(cmnt.Lines, reRemAllCmtMarks.ReplaceAllString(l, ""))
			}
			cmnt.Text = c.Text
		} else {
			cmnt.Lines = append(cmnt.Lines, reRemAllCmtMarks.ReplaceAllString(c.Text, ""))
			cmnt.Text = c.Text
		}
		cmnts = append(cmnts, cmnt)
	}
	return cmnts
}

// AstToFileAndByteBuf: to simply display ast content for an overview of declarations. DEBUG purpose ...
func (gsfs *GoSourceFileStruct) AstToFileAndBBuff(saveToFilename ...string) (bytesBuf *bytes.Buffer, err error) {
	var writer io.Writer
	bytesBuf = new(bytes.Buffer)
	writer = bytesBuf

	err = ast.Fprint(writer, gsfs.fset, gsfs.astFile, ast.NotNilFilter)
	if len(saveToFilename) > 0 {
		if len(saveToFilename[0]) != 0 {

			if err = gsfs.fos.WriteFile(saveToFilename[0], bytesBuf.Bytes(), gsfs.fos.Perms.File); err == nil {
				fmt.Printf("AST file saved successfully: %s\n", saveToFilename[0])
			} else {
				fmt.Printf("Unable to save AST file: %s\n", err.Error())
			}
		}
	}
	return bytesBuf, err
}

// filteringMethods: put methods with their respective structures
func (gsfs *GoSourceFileStruct) filteringMethods() {
	for idx, stru := range gsfs.Struct {
		for _, mtd := range gsfs.tmpMethods {

			if stru.Ident.Name == strings.Trim(mtd.Ident.Type, "*") {
				gsfs.Struct[idx].Methods = append(gsfs.Struct[idx].Methods, mtd)
			}
		}
	}
}
