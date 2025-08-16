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
		result.Print(strings.Join(element.InnerLines, "\n"))
		// TODO: or recurse into InnerElements
		result.Print(element.ClosingTag)
		result.Println()
	}
	return result.Inner().String()
}

type Element struct {
	OpeningTag    string
	InnerLines    []string
	InnerElements []Element
	ClosingTag    string
}

type Scanner struct {
	lines  []string
	cursor int
}

func NewScanner(text string) *Scanner {
	return &Scanner{lines: strings.Split(text, "\n"), cursor: -1}
}
func (this *Scanner) LineNumber() int {
	return this.cursor + 1
}
func (this *Scanner) LineText(offset int) string {
	return this.lines[this.cursor+offset]
}
func (this *Scanner) Line(offset int) (int, string) {
	return this.LineNumber() + offset, this.LineText(offset)
}
func (this *Scanner) IsEOF() bool {
	return this.cursor >= len(this.lines)
}
func (this *Scanner) Advance() bool {
	this.cursor++
	return this.cursor < len(this.lines)
}

func parse(md string) (results []Element) {
	p := makeTag(paragraph)
	scanner := NewScanner(md)
	for scanner.Advance() {
		line, text := scanner.Line(0)
		atxHeader, ok := parseATXHeader(line, text)
		if ok {
			results = append(results, atxHeader)
			continue
		}

		if text != "" {
			p.InnerLines = append(p.InnerLines, text)
		} else {
			if len(p.InnerLines) > 0 {
				results = append(results, p)
			}
			p = makeTag(paragraph)
		}
	}
	if len(p.InnerLines) > 0 {
		results = append(results, p)
	}
	return results
}
func parseATXHeader(lineNumber int, line string) (Element, bool) {
	for range 3 {
		line = strings.TrimPrefix(line, space)
	}
	if content, ok := strings.CutPrefix(line, atx1Prefix); ok {
		return makeATXHeader(lineNumber, h1, content), true
	}
	if content, ok := strings.CutPrefix(line, atx2Prefix); ok {
		return makeATXHeader(lineNumber, h2, content), true
	}
	if content, ok := strings.CutPrefix(line, atx3Prefix); ok {
		return makeATXHeader(lineNumber, h3, content), true
	}
	if content, ok := strings.CutPrefix(line, atx4Prefix); ok {
		return makeATXHeader(lineNumber, h4, content), true
	}
	if content, ok := strings.CutPrefix(line, atx5Prefix); ok {
		return makeATXHeader(lineNumber, h5, content), true
	}
	if content, ok := strings.CutPrefix(line, atx6Prefix); ok {
		return makeATXHeader(lineNumber, h6, content), true
	}
	return Element{}, false
}
func makeATXHeader(n int, tag, line string) Element {
	line, _, _ = strings.Cut(line, atxSuffix)
	result := makeTag(tag)
	result.InnerLines = append(result.InnerLines, strings.TrimSpace(line))
	return result
}
func makeTag(name string) Element {
	return Element{
		OpeningTag: fmt.Sprintf(openingTagTemplate, name),
		ClosingTag: fmt.Sprintf(closingTagTemplate, name),
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

	paragraph = "p"
)
