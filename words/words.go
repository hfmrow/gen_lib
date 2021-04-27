// words.go

package words

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	gsl "github.com/hfmrow/gen_lib/slices"
)

// WordsFrequency: Counts the frequency with which words appear in a text.
func WordsFrequency(text string) (list [][]string) {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`) //  Remove all non alpha-numeric char
	tmpText := nonAlNum.ReplaceAllString(text, "")
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`) //	Trim [[:space:]] and clean multi [[:space:]] inside
	tmpText = strings.TrimSpace(remInside.ReplaceAllString(tmpText, " "))
	tmpLines := strings.Fields(tmpText)

	for _, word := range tmpLines {
		if !gsl.IsExist2dCol(list, word, 0) {
			list = append(list, []string{word, ""})
		}
	}
	tmpText = strings.Join(tmpLines, " ")
	for idx, word := range list {
		// Compile for matching whole word
		regX, err := regexp.Compile(`\b(` + word[0] + `)\b`)
		if err == nil {
			list[idx][1] = fmt.Sprintf("%d", len(regX.FindAllString(tmpText, -1)))
		}
	}
	// Sorting result ...
	sort.SliceStable(list, func(i, j int) bool {
		return fmt.Sprintf("%d", list[i][1]) > fmt.Sprintf("%d", list[j][1])
	})
	return list
}
