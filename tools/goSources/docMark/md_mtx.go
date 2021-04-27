// md_mtx.go

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
	"fmt"
)

var ()

type optMD int

const (
	_MD_DEFAULT optMD = 1
	_MD_FUNC    optMD = 1 << 1
	_MD_TYPE    optMD = 1 << 2
	_MD_CONST   optMD = 1 << 3
	_MD_SUBLIST optMD = 1 << 4

	UNEXPORTED = "// contains filtered or unexported fields"
)

func (dm *DocMark) DocMarkBuild(_ interface{}) {
	fmt.Printf(dm.dispDecl("ProcNetDev", `type ProcNetDev struct {

    // Suffix, default: 'iB' !!! max char length = 15
    Suffix string
    // Unit, default: '/s' !!! max char length = 15
    Unit       string
    Interfaces []iface
    // contains filtered or unexported fields
}`, "", _MD_TYPE))

	fmt.Printf(dm.dispDecl("ProcNetDevNew",
		`func ProcNetDevNew(pid ...uint32) (*ProcNetDev, error)`,
		"", _MD_FUNC))

	fmt.Printf(dm.dispDecl("ProcNetDevNew",
		`func ProcNetDevNew(pid ...uint32) (*ProcNetDev, error)`,
		`ProcNetDevNew: Create and initialize the "C" structure. If a "pid" is given, the statistics relate to the process. Otherwise, it's the overall flow`, _MD_FUNC))
}

func (dm *DocMark) indexDisp(name string, opt optMD) string {
	var (
		sub    string
		prefix = `func`
	)
	switch {
	case opt&_MD_FUNC != 0:
		prefix = `func`
	case opt&_MD_TYPE != 0:
		prefix = `type`
	case opt&_MD_CONST != 0:
		prefix = `const`
	case opt&_MD_SUBLIST != 0:
		sub = "  "
	}
	return sub + `- [` + prefix + ` ` + name + `](<#` + prefix + `-` + toKebab(name) + `>)`

}

func (dm *DocMark) dispDecl(name, decl, comment string, opt optMD) string {

	var prefix = `### func`

	switch {
	case opt&_MD_FUNC != 0:
		prefix = `### func `
	case opt&_MD_TYPE != 0:
		prefix = `## type `
	}
	return prefix + name + `

` + dm.startCode() + `
` + decl + `
` + dm.endCode() + `
` + dm.commentDisp(comment)
}

func (dm *DocMark) commentDisp(comment string) string {

	if len(comment) > 0 {
		return `
` + comment + `

`
	}
	return `
`
}

func (dm *DocMark) startCode() string {
	return string([]byte{0x60, 0x60, 0x60}) + "go"
}
func (dm *DocMark) endCode() string {
	return string([]byte{0x60, 0x60, 0x60})
}
