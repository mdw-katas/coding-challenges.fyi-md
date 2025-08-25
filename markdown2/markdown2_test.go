package markdown2

import (
	"strings"
	"testing"

	"github.com/mdw-go/printing"
	"github.com/mdw-katas/coding-challenges.fyi-md/util/str"
)

func Test(t *testing.T) {
	parser := NewPhase1Parser()
	for line := range str.IterateLines(strings.NewReader(input)) {
		parser.Feed(line)
	}
	parser.Finalize()
	parser.root.Render(printing.NewPrinter(t.Output()), 0)
}

const input = `# Main Heading
> Lorem ipsum dolor
sit amet.
> - Qui *quodsi iracundia*
> - aliquando id`
