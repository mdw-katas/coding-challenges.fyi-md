package str

import (
	"strings"
	"unicode"
)

func TrimMinorLeadingIndent(line string) string {
	for range minorIndentWidth {
		line = strings.TrimPrefix(line, space)
	}
	return line
}
func TrimTrailingSpace(line string) string {
	return strings.TrimSuffix(line, space)
}
func IsOnly(text string, char rune) bool {
	text = TrimMinorLeadingIndent(text)
	text = TrimTrailingSpace(text)
	count := strings.Count(text, string(char))
	return count > 0 && count == len(text)
}

func TrimIndent(line string) string {
	for len(line) > 0 && unicode.IsSpace(rune(line[0])) {
		line = line[1:]
	}
	return line
}
func CutIndent(line string) (indent, content string, ok bool) {
	x := 0
	for ; x < len(line); x++ {
		if !unicode.IsSpace(rune(line[x])) {
			break
		}
	}
	return line[:x], line[x:], x > 0
}

const space = " "
const minorIndentWidth = 3
