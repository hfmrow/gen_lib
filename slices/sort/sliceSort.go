// sliceSort.go

package sort

import (
	"fmt"
	"sort"
	"strings"

	glsssn "github.com/hfmrow/gen_lib/strings/strNum"
	glte "github.com/hfmrow/gen_lib/types"
	gltsct "github.com/hfmrow/gen_lib/types/convert"
)

// SliceSortDate: Sort 2d string slice with date inside
func SliceSortDate(slice [][]string, fmtDate string, dateCol, secDateCol int, ascendant bool) [][]string {
	fieldsCount := len(slice[0]) // Get nb of columns
	var firstLine int
	var previous, after string
	var positiveidx, negativeidx int
	// compute unix date using given column numbers
	for idx := firstLine; idx < len(slice); idx++ {
		dateStr := glte.FindDate(slice[idx][dateCol], fmtDate)
		if dateStr != nil { // search for 1st column
			slice[idx] = append(slice[idx], fmt.Sprintf("%d", glte.FormatDate(fmtDate, dateStr[0]).Unix()))
		} else if secDateCol != -1 { // Check for second column if it was given
			dateStr = glte.FindDate(slice[idx][secDateCol], fmtDate)
			if dateStr != nil { // If date was not found in 1st column, search for 2nd column
				slice[idx] = append(slice[idx], fmt.Sprintf("%d", glte.FormatDate(fmtDate, slice[idx][secDateCol]).Unix()))
			} else { //  in case where none of the columns given contain date field, put null string if there is no way to find a date
				slice[idx] = append(slice[idx], ``)
			}
		} else { // put null string if there is no way to find a date
			slice[idx] = append(slice[idx], ``)
		}
	}
	// Ensure we always have a value in sorting field (get previous or next closer)
	for idx := firstLine; idx < len(slice); idx++ {
		if slice[idx][fieldsCount] == `` {
			for idxFind := firstLine + 1; idxFind < len(slice); idxFind++ {
				positiveidx = idx + idxFind
				negativeidx = idx - idxFind
				if positiveidx >= len(slice) { // Check index to avoiding 'out of range'
					positiveidx = len(slice) - 1
				}
				if negativeidx <= 0 {
					negativeidx = 0
				}
				after = slice[positiveidx][fieldsCount] // Get previous or next value
				previous = slice[negativeidx][fieldsCount]
				if previous != `` { // Set value, prioritise the previous one.
					slice[idx][fieldsCount] = previous
					break
				}
				if after != `` {
					slice[idx][fieldsCount] = after
					break
				}
			}
		}
	}
	tmpLines := make([][]string, 0)
	if ascendant != true {
		// Sort by unix date preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] > slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	} else {
		// Sort by unix date preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return slice[i][len(slice[i])-1] < slice[j][len(slice[i])-1] })
		for idx := firstLine; idx < len(slice); idx++ { // Store row count elements - 1
			tmpLines = append(tmpLines, slice[idx][:len(slice[idx])-1])
		}
	}
	return tmpLines
}

// SliceSortString: Sort 2d string slice
func SliceSortString(slice [][]string, col int, ascendant, caseSensitive, numbered bool) {
	if numbered {
		var tmpWordList []string
		for _, wrd := range slice {
			tmpWordList = append(tmpWordList, wrd[col])
		}
		numberedWords := new(glsssn.WordWithDigit)
		numberedWords.Init(tmpWordList)

		if ascendant != true {
			// Sort string preserving order descendant
			sort.SliceStable(slice, func(i, j int) bool {
				return numberedWords.FillWordToMatchMaxLength(slice[i][col]) > numberedWords.FillWordToMatchMaxLength(slice[j][col])
			})
		} else {
			// Sort string preserving order ascendant
			sort.SliceStable(slice, func(i, j int) bool {
				return numberedWords.FillWordToMatchMaxLength(slice[i][col]) < numberedWords.FillWordToMatchMaxLength(slice[j][col])
			})
		}
		return
	}

	toLowerCase := func(inString string) string {
		return inString
	}
	if !caseSensitive {
		toLowerCase = func(inString string) string { return strings.ToLower(inString) }
	}

	if ascendant != true {
		// Sort string preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool { return toLowerCase(slice[i][col]) > toLowerCase(slice[j][col]) })
	} else {
		// Sort string preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool { return toLowerCase(slice[i][col]) < toLowerCase(slice[j][col]) })
	}
}

// SliceSortFloat: Sort 2d string with float value
func SliceSortFloat(slice [][]string, col int, ascendant bool, decimalChar string) {
	if ascendant != true {
		// Sort string (float) preserving order descendant
		sort.SliceStable(slice, func(i, j int) bool {
			return gltsct.StringDecimalSwitchFloat(decimalChar, slice[i][col]) > gcv.StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	} else {
		// Sort string (float) preserving order ascendant
		sort.SliceStable(slice, func(i, j int) bool {
			return gltsct.StringDecimalSwitchFloat(decimalChar, slice[i][col]) < gcv.StringDecimalSwitchFloat(decimalChar, slice[j][col])
		})
	}
}
