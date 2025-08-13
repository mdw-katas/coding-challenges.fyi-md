package markdown

import (
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
		content, ok := strings.CutPrefix(line, "# ")
		if ok {
			//content, _, _ := strings.Cut(content, "#") // TODO: trailing pound signs are ignored
			results = append(results, Element{
				Line:         l + 1,
				OpeningTag:   "<h1>",
				InnerContent: content,
				ClosingTag:   "</h1>",
			})
		}

	}
	return results
}
