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
	"log"
	"regexp"
	"strings"
)

// paramsResults:
type paramsResults struct {
	Params,
	Results []string
	InLine string
}

// TODO to replace using AST... WIP
func (c *content) GetParamsResults() (out paramsResults) {

	var (
		reVar      = regexp.MustCompile(`(?m)(\(.*\))(.*)`)
		reVarInOut = regexp.MustCompile(`(?m)(\)\s)`)
	)
	out.InLine = removeDupSpace(strings.Join(strings.Split(string(c.Head), c.eol), " "))
	tmpVar := reVar.FindString(out.InLine)
	tmpStrs := reVarInOut.Split(tmpVar, -1)
	for idx, v := range tmpStrs {
		tmpStr := strings.TrimPrefix(strings.TrimSuffix(v, ")"), "(") // TODO replace with regexp
		switch idx {
		case 0:
			out.Params = append(out.Params, c.splitArgs(tmpStr)...)
		case 1:
			out.Results = append(out.Results, c.splitArgs(tmpStr)...)
		}
	}
	return
}

// splitArgs: separate parameters and results arguments from inline declaration.
func (c *content) splitArgs(in string) (out []string) {

	splitted := strings.Split(in, ",")
	for _, s := range splitted {
		spcSplitted := strings.Split(s, " ")
		if len(spcSplitted) > 0 {
			out = append(out, spcSplitted[len(spcSplitted)-1])
			continue
		}
		// I don't think there is any point in emphasizing this as a function
		// might not have parameters and results. Doing it just in debug mode.
		log.Printf("docMark/splitArgs: A mistake occure on getting arg from " + in)
	}
	return
}
