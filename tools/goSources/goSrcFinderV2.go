// goSrcFinder.go

/*
	©2019 H.F.M. MIT license
*/

// Parse go source file and retrieve information about function, variables, structures ...

package goSources

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	glfsfo "github.com/hfmrow/gen_lib/files/filesOperations"
	glsg "github.com/hfmrow/gen_lib/strings"
)

var OS = glfsfo.PermsStructNew()

type GoSourceFileStructure struct {
	Filename           string
	WriteAstToFilename string // Leave blank, do not save the AST file. To do before initialisation if needed.
	Package            string
	Imports            []Declaration
	Func               []Declaration
	CallExpr           []Declaration
	Comments           []Declaration
	Struct             []Declaration
	Var                []Declaration
	UnImplemented      []Declaration
	Unknown            []Declaration
	Eol                string        // End of line type of the input file
	AstOut             *bytes.Buffer // AST representation of the input file

	// Used to get a new empty structure when the library is not directly in the "import" section
	// but declared as a new type in the source of the end user (that's how I use it in most cases).
	EmptyDeclStruct      Declaration
	EmptySliceDeclStruct []Declaration

	// Unexported
	data         []byte // File content
	linesIndexes [][]int
	offset       token.Pos // Define if we start at 0 or 1  when counting lines and offsets positions.
	astFile      *ast.File
	fset         *token.FileSet // Positions are relative to fset.
}

type Declaration struct {
	From      string
	Params    []Declaration
	Name      string
	Value     string
	Type      string // string, int, bool ...
	Kind      string // var, const, :=, =
	Content   []byte
	OfstStart int
	OfstEnd   int
	LineStart int
	LineEnd   int
	// Prarams secific
	TypeAst     string
	Desc        string
	NameFromSrc string
}

// Local var declaration for GoSourceFileStructure's (declaration) fields. Used by functions
// that does not a part of the main structure (like ast.walk).
var localGsfs *GoSourceFileStructure // bridge to GoSourceFileStructure.Variables

// GoSourceFileStructureSetup: setup and retieve information for designed file.
// Notice: the lines numbers and offsets start at 0. Set "zero" at false to start at 1.
func (gsfs *GoSourceFileStructure) GoSourceFileStructureSetup(filename string, zero ...bool) (err error) {
	gsfs.offset = 1 // lines start at 0 (substract 1 for each offsets position)
	if len(zero) > 0 {
		if zero[0] {
			gsfs.offset = 0
		}
	}

	// Declare local variables to GoSourceFileStructure.
	localGsfs = gsfs

	// Loading data (file)
	gsfs.Filename = filename
	gsfs.fset = token.NewFileSet()
	if gsfs.astFile, err = parser.ParseFile(gsfs.fset, gsfs.Filename, nil, parser.ParseComments); err == nil {
		if gsfs.data, err = ioutil.ReadFile(gsfs.Filename); err == nil {
			if gsfs.AstOut, err = gsfs.astPrintToBuf(gsfs.WriteAstToFilename); err == nil {
				err = gsfs.fillDeclaration()
			}
		}
	}
	return
}

// GetPosFunc: get line start and line end of the specified function
func (gsfs *GoSourceFileStructure) fillDeclaration() (err error) {

	// Setting internal variables
	gsfs.Eol = glsg.GetTextEOL(gsfs.data)
	gsfs.data = append(gsfs.data, []byte(gsfs.Eol)...) // Add an eol to avoid a f..k..g issue where the last line wasn't analysed
	eolRegx := regexp.MustCompile(gsfs.Eol)
	eolPositions := eolRegx.FindAllIndex(gsfs.data, -1)
	// Creating lines indexes
	gsfs.linesIndexes = append(gsfs.linesIndexes, []int{0, eolPositions[0][0]})
	for idx := 1; idx < len(eolPositions); idx++ {
		gsfs.linesIndexes = append(gsfs.linesIndexes, []int{eolPositions[idx-1][1], eolPositions[idx][0]})
	}

	gsfs.Package = gsfs.astFile.Name.String() // get package name
	gsfs.goInspect()
	return
}

// GetPosFunc: get line start and line end of the specified function
func (gsfs *GoSourceFileStructure) GetFuncPos(fName string) (lStart, lEnd int) {
	for _, fct := range gsfs.Func {
		if fct.Name == fName {
			return fct.LineStart, fct.LineEnd
		}
	}
	return -1, -1
}

// getDeclByName: get content Declaration of name from kindOf declarations list
// kind parameter try to match the kind property of the Declaration to search
func (gsfs *GoSourceFileStructure) getDeclByName(kindOf *[]Declaration, dName string, kind ...string) (decl *Declaration) {
	var foundDecl = new([]Declaration)
	decl = new(Declaration)

	for _, d := range *kindOf {
		if d.Name == dName {
			*foundDecl = append(*foundDecl, d)
			// fmt.Println(d)
		}
	}
	switch len(*foundDecl) {
	case 0:
		// 0 found
		return nil
	case 1:
		// only 1 found
		return &(*foundDecl)[0]
	default:
		// more than 1 found
		if len(kind) > 0 {
			for idx, d := range *foundDecl {
				for _, k := range kind {
					if d.Kind == k {
						return &(*foundDecl)[idx]
					}
				}
			}
		} else {
			// kind argument not present, get the 1st found
			return &(*foundDecl)[0]
		}
	}
	return nil
}

// GetFuncDeclByName: get function declaration by name
func (gsfs *GoSourceFileStructure) GetFuncByName(dName string) (decl *Declaration) {
	return gsfs.getDeclByName(&gsfs.Func, dName)
}

// GetFuncDeclByName:  get structure declaration by name
func (gsfs *GoSourceFileStructure) GetStructByName(dName string) (decl *Declaration) {
	return gsfs.getDeclByName(&gsfs.Struct, dName)
}

// GetFuncDeclByName:  get variable declaration by name
// "kind" means: "var" ,"const" ,":=" ,"="
func (gsfs *GoSourceFileStructure) GetVarByName(dName string, kind ...string) (decl *Declaration) {
	return gsfs.getDeclByName(&gsfs.Var, dName, kind...)
}

// goInspect: parse go file and retrieve into structure that was found.
func (gsfs *GoSourceFileStructure) goInspect() {
	// getImports
	for _, v := range gsfs.astFile.Imports {
		startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
		gsfs.Imports = append(gsfs.Imports,
			Declaration{
				Name:      v.Path.Value,
				Type:      "import",
				OfstStart: startOffset,
				OfstEnd:   endOffset,
				LineStart: lStrt,
				LineEnd:   lEnd,
				Content:   gsfs.data[startOffset:endOffset]})
	}
	// Inspect node
	ast.Inspect(gsfs.astFile, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncDecl: // Functions
			params := []Declaration{} //  parameters
			if paramList := v.Type.Params.List; len(paramList) > 0 {
				for _, param := range paramList {
					startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
					params = append(params,
						Declaration{
							Name:        param.Names[0].String(),
							Type:        "params",
							OfstStart:   startOffset,
							OfstEnd:     endOffset,
							LineStart:   lStrt,
							LineEnd:     lEnd,
							TypeAst:     fmt.Sprintf("%T", param.Type),
							Desc:        fmt.Sprintf("%+v", param.Type),
							NameFromSrc: string(gsfs.data[startOffset:endOffset])})
				}
			}
			startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
			gsfs.Func = append(gsfs.Func,
				Declaration{
					Name:      v.Name.String(),
					Type:      "func",
					OfstStart: startOffset,
					OfstEnd:   endOffset,
					LineStart: lStrt,
					LineEnd:   lEnd,
					Content:   gsfs.data[startOffset:endOffset],
					Params:    params})

		case *ast.CallExpr: // Inside functions
			startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
			gsfs.CallExpr = append(gsfs.CallExpr,
				Declaration{
					Name:      fmt.Sprintf("%v", v.Fun),
					Type:      "inside_func",
					OfstStart: startOffset,
					OfstEnd:   endOffset,
					LineStart: lStrt,
					LineEnd:   lEnd,
					Content:   gsfs.data[startOffset:endOffset]})

		case *ast.TypeSpec: // Retrieve structures
			startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
			gsfs.Struct = append(gsfs.Struct,
				Declaration{
					Name:      v.Name.String(),
					Type:      "type",
					OfstStart: startOffset,
					OfstEnd:   endOffset,
					LineStart: lStrt,
					LineEnd:   lEnd,
					Content:   gsfs.data[startOffset:endOffset]})

		case *ast.GenDecl: // Retrieve Variables
			tmp := parseGenDeclVar(v)
			localGsfs.Var = append(localGsfs.Var, tmp...)

		case *ast.AssignStmt: // Retrieve assigned ariables
			localGsfs.Var = append(localGsfs.Var, parseAssignStmtVar(v)...)

		default: // Not implemented tokens, stored for analysis
			if vType := reflect.TypeOf(v); vType != nil {
				startOffset, endOffset, lStrt, lEnd := gsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
				gsfs.Unknown = append(gsfs.Unknown,
					Declaration{
						Name:      vType.String(),
						OfstStart: startOffset,
						OfstEnd:   endOffset,
						LineStart: lStrt,
						LineEnd:   lEnd,
						Content:   gsfs.data[startOffset:endOffset]})
			}
		}
		return true
	})
	return
}

// getLineFro:mIdx: get the line number corresponding to indexes. Notice, line number start at 0.
func (gsfs *GoSourceFileStructure) getLineFromOffsets(sOfst, eOfst int) (lStart, lEnd int) {
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

// return position values-offset as int and Lines numbers for the given node.
func (gsfs *GoSourceFileStructure) getOffsetPosAndLinesAsInt(pos, end token.Pos) (sPos, ePos, lStart, lEnd int) {
	sPos, ePos = int(pos-gsfs.offset), int(end-gsfs.offset)
	lStart, lEnd = gsfs.getLineFromOffsets(sPos, ePos)
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

func parseGenDeclVar(v *ast.GenDecl) (decl []Declaration) {
	for _, s := range v.Specs {
		if kind, ok := varWalker(v.Tok); ok { // Variables scanner
			startOffset, endOffset, lStrt, lEnd := localGsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
			var value, name, svType string

			for _, vsn := range s.(*ast.ValueSpec).Names {
				name = vsn.Obj.Name
				break // TODO better than Only one is handled actually
			}
			for _, vsv := range s.(*ast.ValueSpec).Values {
				switch vsv.(type) {
				case *ast.BasicLit:
					value = vsv.(*ast.BasicLit).Value
					svType = strings.ToLower(vsv.(*ast.BasicLit).Kind.String())
				default:
					localGsfs.UnImplemented = append(localGsfs.UnImplemented,
						Declaration{
							Name:      name,
							Kind:      "parseGenDeclVar1",
							Value:     "not implemented",
							TypeAst:   reflect.TypeOf(vsv).String(),
							OfstStart: startOffset,
							OfstEnd:   endOffset,
							LineStart: lStrt,
							LineEnd:   lEnd,
							Content:   localGsfs.data[startOffset:endOffset]})
				}
				break // TODO better than Only one is handled actually
			}

			if s.(*ast.ValueSpec).Type != nil {
				switch s.(*ast.ValueSpec).Type.(type) {
				case *ast.Ident:
					svType += s.(*ast.ValueSpec).Type.(*ast.Ident).Name
				default:
					localGsfs.UnImplemented = append(localGsfs.UnImplemented,
						Declaration{
							Name:      name,
							Kind:      "parseGenDeclVar2",
							Value:     "not implemented",
							TypeAst:   reflect.TypeOf(s.(*ast.ValueSpec).Type).String(),
							OfstStart: startOffset,
							OfstEnd:   endOffset,
							LineStart: lStrt,
							LineEnd:   lEnd,
							Content:   localGsfs.data[startOffset:endOffset]})
				}
			}
			decl = append(decl, Declaration{
				Name:      name,
				Kind:      kind,   // VAR, CONST, ...
				Value:     value,  // var content
				Type:      svType, // var type
				OfstStart: startOffset,
				OfstEnd:   endOffset,
				LineStart: lStrt,
				LineEnd:   lEnd,
				Content:   localGsfs.data[startOffset:endOffset]})
		}
	}
	return
}

func parseAssignStmtVar(v *ast.AssignStmt) (decl []Declaration) {
	if kind, ok := varWalker(v.Tok); ok { // Variables scanner
		startOffset, endOffset, lStrt, lEnd := localGsfs.getOffsetPosAndLinesAsInt(v.Pos(), v.End())
		var value, name, svType string

		for _, lhs := range v.Lhs {
			switch lhs.(type) {
			case *ast.Ident:
				name = lhs.(*ast.Ident).Name // var name
			default:
				localGsfs.UnImplemented = append(localGsfs.UnImplemented,
					Declaration{
						Name:      name,
						Kind:      "parseAssignStmtVar3",
						Value:     "not implemented",
						TypeAst:   reflect.TypeOf(lhs).String(),
						OfstStart: startOffset,
						OfstEnd:   endOffset,
						LineStart: lStrt,
						LineEnd:   lEnd,
						Content:   localGsfs.data[startOffset:endOffset]})
			}
			break // TODO better than Only one is handled actually
		}
		for _, rhs := range v.Rhs {

			switch rhs.(type) {

			case *ast.BasicLit:
				svType = strings.ToLower(rhs.(*ast.BasicLit).Kind.String()) // type string, int ...
				value = rhs.(*ast.BasicLit).Value                           // var content

			case *ast.BinaryExpr:
				if _, ok := varWalker(v.Tok); ok { // Variables scanner
					switch rhs.(*ast.BinaryExpr).X.(type) {
					case *ast.Ident:
						// name = rhs.(*ast.BinaryExpr).X.(*ast.Ident).Name // Name of the sub object to use when implement concatenate obj
					default:
						localGsfs.UnImplemented = append(localGsfs.UnImplemented,
							Declaration{
								Name:      name,
								Kind:      "parseAssignStmtVar1",
								Value:     "not implemented",
								TypeAst:   reflect.TypeOf(rhs.(*ast.BinaryExpr).X).String(),
								OfstStart: startOffset,
								OfstEnd:   endOffset,
								LineStart: lStrt,
								LineEnd:   lEnd,
								Content:   localGsfs.data[startOffset:endOffset]})
					}

					switch rhs.(*ast.BinaryExpr).Y.(type) {
					case *ast.BasicLit:
						svType = strings.ToLower(rhs.(*ast.BinaryExpr).Y.(*ast.BasicLit).Kind.String())
						value = rhs.(*ast.BinaryExpr).Y.(*ast.BasicLit).Value // var content
					default:
						localGsfs.UnImplemented = append(localGsfs.UnImplemented,
							Declaration{
								Name:      name,
								Kind:      "parseAssignStmtVar2",
								Value:     "not implemented",
								TypeAst:   reflect.TypeOf(rhs.(*ast.BinaryExpr).Y).String(),
								OfstStart: startOffset,
								OfstEnd:   endOffset,
								LineStart: lStrt,
								LineEnd:   lEnd,
								Content:   localGsfs.data[startOffset:endOffset]})
						fmt.Printf("parseAssignStmtVar2, "+name+", lines %d,%d, not implemented: %s\n", lStrt, lEnd, reflect.TypeOf(rhs.(*ast.BinaryExpr).Y).String())
					}
				}
			}
			decl = append(decl, Declaration{
				Name:      name,
				Kind:      kind,   // type var, const ...
				Value:     value,  // var content
				Type:      svType, //  type string, int ...
				OfstStart: startOffset,
				OfstEnd:   endOffset,
				LineStart: lStrt,
				LineEnd:   lEnd,
				Content:   localGsfs.data[startOffset:endOffset]})
		}
	}
	return
}

// AstPrint: Simply display ast content for an overview of declarations.
func (gsfs *GoSourceFileStructure) astPrintToBuf(saveToFilename ...string) (bytesBuf *bytes.Buffer, err error) {

	var writer io.Writer
	bytesBuf = new(bytes.Buffer)
	writer = bytesBuf
	err = ast.Fprint(writer, gsfs.fset, gsfs.astFile, ast.NotNilFilter)
	if len(saveToFilename) > 0 {
		if len(saveToFilename[0]) != 0 {
			if err = ioutil.WriteFile(saveToFilename[0], bytesBuf.Bytes(), os.ModePerm&(OS.USER_RW|OS.GROUP_RW|OS.OTH_R)); err == nil {
				fmt.Printf("AST file saved successfully: %s\n", saveToFilename[0])
			} else {
				fmt.Printf("Unable to save AST file: %s\n", err.Error())
			}
		}
	}
	return bytesBuf, err
}

func AstWalk(filename string) {
	v := visitor{}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	ast.Walk(&v, f)
}

type visitor struct {
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n != nil {
		switch v := n.(type) {
		case *ast.GenDecl:
			localGsfs.Var = append(localGsfs.Var, parseGenDeclVar(v)...)
		case *ast.FuncDecl:

		case *ast.AssignStmt:
			localGsfs.Var = append(localGsfs.Var, parseAssignStmtVar(v)...)
		}
	}
	return v
}

// displayDecl: Simply display retrieved content for debug purpose
func DisplayDecl(title string, decl []Declaration, value string, pDesc ...bool) {

	var dispDesc bool
	if len(pDesc) > 0 {
		dispDesc = pDesc[0]
	}
	fmt.Printf("\t** %s (%d) **\n", title, len(decl))
	for _, elem := range decl {

		fmt.Printf("ofst: %d,%d\tlines: %d,%d \tname: %s, value: %s, type: %s, kind: %s, ast: %s\n",
			elem.OfstStart, elem.OfstEnd, elem.LineStart, elem.LineEnd, elem.Name, elem.Value, elem.Type, elem.Kind, elem.TypeAst)

		if dispDesc && len(elem.Params) > 0 {
			for _, params := range elem.Params {
				fmt.Printf("      %d,%d\tlines: %d,%d, \tname: %s, typeAst:%s, desc: %s, frmSrc: %s\n",
					params.OfstStart, params.OfstEnd, params.LineStart, params.LineEnd,
					params.Name, params.TypeAst, params.Desc, strings.Join(strings.Fields(params.NameFromSrc), " "))
			}
		}
	}
}
