package markdown

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/mdw-go/printing"
	"github.com/mdw-katas/coding-challenges.fyi-md/util/str"
)

func Test(t *testing.T) {
	input, err := os.ReadFile("testdata/input.md")
	if err != nil {
		t.Fatal(err)
	}
	parser := NewPhase1Parser()
	for line := range str.IterateLines(bytes.NewReader(input)) {
		parser.Feed(line)
	}
	parser.Finalize()
	parser.root.Render(printing.NewPrinter(t.Output()), 0)
	RenderHTML(t.Output(), parser.root)
}

func (this *Node) Render(printer printing.Printer, level int) {
	printer.Print(strings.Repeat("  ", level))
	printer.Println(this.Token.String(), strings.ReplaceAll(this.Text, "\n", "\\n"))
	for child := range this.Children.All() {
		child.Render(printer, level+1)
	}
}
