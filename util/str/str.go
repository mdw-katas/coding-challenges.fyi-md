package str

import "strings"

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

const space = " "
const minorIndentWidth = 3
