package markdown

import (
	"fmt"
	"strings"

	"github.com/mdw-go/printing"
)

func ConvertToHTML(md string) string {
	result := printing.NewBuilder()
	for _, element := range parse(md) {
		result.Print(element.OpeningTag)
		result.Print(element.InnerContent)
		result.Print(element.ClosingTag)
		result.Println()
	}
	return result.Inner().String()
}

type Element struct {
	Line          int
	OpeningTag    string
	InnerContent  string
	InnerElements []Element
	ClosingTag    string
}

func parse(md string) (results []Element) {
	for l, line := range strings.Split(md, "\n") {
		header, ok := parseATXHeader(l, line)
		if ok {
			results = append(results, header)
		}
	}
	return results
}
func parseATXHeader(lineNumber int, line string) (Element, bool) {
	for range 3 {
		line = strings.TrimPrefix(line, space)
	}
	if content, ok := strings.CutPrefix(line, atx1Prefix); ok {
		return atxHeader(lineNumber, h1, content), true
	}
	if content, ok := strings.CutPrefix(line, atx2Prefix); ok {
		return atxHeader(lineNumber, h2, content), true
	}
	if content, ok := strings.CutPrefix(line, atx3Prefix); ok {
		return atxHeader(lineNumber, h3, content), true
	}
	if content, ok := strings.CutPrefix(line, atx4Prefix); ok {
		return atxHeader(lineNumber, h4, content), true
	}
	if content, ok := strings.CutPrefix(line, atx5Prefix); ok {
		return atxHeader(lineNumber, h5, content), true
	}
	if content, ok := strings.CutPrefix(line, atx6Prefix); ok {
		return atxHeader(lineNumber, h6, content), true
	}
	return Element{}, false
}
func atxHeader(n int, tag, line string) Element {
	line, _, _ = strings.Cut(line, atxSuffix)
	return Element{
		Line:         n + 1,
		OpeningTag:   fmt.Sprintf(openingTagTemplate, tag),
		InnerContent: strings.TrimSpace(line),
		ClosingTag:   fmt.Sprintf(closingTagTemplate, tag),
	}
}

const (
	space = " "

	openingTagTemplate = "<%s>"
	closingTagTemplate = "</%s>"

	h1 = "h1"
	h2 = "h2"
	h3 = "h3"
	h4 = "h4"
	h5 = "h5"
	h6 = "h6"

	atx1Prefix = "# "
	atx2Prefix = "## "
	atx3Prefix = "### "
	atx4Prefix = "#### "
	atx5Prefix = "##### "
	atx6Prefix = "###### "
	atxSuffix  = " #"
)
