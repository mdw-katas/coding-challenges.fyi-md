package markdown

import (
	"fmt"
	"strings"

	"github.com/mdw-go/printing"
	"github.com/mdw-katas/coding-challenges.fyi-md/util/str"
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
func (this *Scanner) LineText(offset int) string {
	return this.lines[this.cursor+offset]
}
func (this *Scanner) IsEOF() bool {
	return this.cursor >= len(this.lines)
}
func (this *Scanner) Advance() bool {
	this.cursor++
	return !this.IsEOF()
}

func parse(md string) (results []Element) {
	p := makeTag(paragraph)
	scanner := NewScanner(md)
	for scanner.Advance() {
		text := scanner.LineText(0)

		atxHeader, ok := parseATXHeader(text)
		if ok {
			results = append(results, atxHeader)
			continue
		}

		setextHeader, ok := parseSetextHeader(text, p.InnerLines)
		if ok {
			results = append(results, setextHeader)
			p = makeTag(paragraph)
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

func parseATXHeader(line string) (Element, bool) {
	line = str.TrimMinorLeadingIndent(line)
	for _, attempt := range atxAttempts {
		prefix, tag := attempt[0], attempt[1]
		if content, ok := strings.CutPrefix(line, prefix); ok {
			return makeATXHeader(tag, content), true
		}
	}
	return Element{}, false
}
func makeATXHeader(tag, line string) Element {
	line, _, _ = strings.Cut(line, atxSuffix)
	result := makeTag(tag)
	result.InnerLines = append(result.InnerLines, strings.TrimSpace(line))
	return result
}

func parseSetextHeader(text string, precedingLines []string) (Element, bool) {
	if str.IsOnly(text, equal) {
		return makeSetextHeader(h1, precedingLines), true
	}
	if str.IsOnly(text, dash) {
		return makeSetextHeader(h2, precedingLines), true
	}
	return Element{}, false
}
func makeSetextHeader(tag string, lines []string) Element {
	result := makeTag(tag)
	result.InnerLines = lines
	return result
}

func makeTag(name string) Element {
	return Element{
		OpeningTag: fmt.Sprintf(openingTagTemplate, name),
		ClosingTag: fmt.Sprintf(closingTagTemplate, name),
	}
}

const (
	dash  = '-'
	equal = '='

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

var atxAttempts = [][]string{
	{atx1Prefix, h1},
	{atx2Prefix, h2},
	{atx3Prefix, h3},
	{atx4Prefix, h4},
	{atx5Prefix, h5},
	{atx6Prefix, h6},
}
