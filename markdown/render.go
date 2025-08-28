package markdown

import (
	"io"

	"github.com/mdw-go/printing"
)

func RenderHTML(w io.Writer, node *Node) {
	printer := printing.NewPrinter(w)
	switch node.Token {
	case TokenNone: // no-op
	case TokenDocument:
		for child := range node.Children.All() {
			RenderHTML(w, child)
		}
	case TokenHeader1:
		printer.Printf("<h1>%s</h1>\n", node.Text)
	case TokenHeader2:
		printer.Printf("<h2>%s</h2>\n", node.Text)
	case TokenHeader3:
		printer.Printf("<h3>%s</h3>\n", node.Text)
	case TokenHeader4:
		printer.Printf("<h4>%s</h4>\n", node.Text)
	case TokenHeader5:
		printer.Printf("<h5>%s</h5>\n", node.Text)
	case TokenHeader6:
		printer.Printf("<h6>%s</h6>\n", node.Text)
	case TokenThematicBreak:
		printer.Println("<hr>")
	case TokenOrderedList:
		printer.Println("<ol>")
		for child := range node.Children.All() {
			RenderHTML(w, child)
		}
		printer.Println("</ol>")
	case TokenUnorderedList:
		printer.Println("<ul>")
		for child := range node.Children.All() {
			RenderHTML(w, child)
		}
		printer.Println("</ul>")
	case TokenListItem:
		printer.Printf("<li>%s</li>\n", node.Text)
	case TokenBlockQuote:
		printer.Printf("<blockquote><p>%s</p></blockquote>\n", node.Text)
	case TokenParagraph:
		printer.Printf("<p>%s</p>\n", node.Text)
	case TokenPreCode:
		if node.FencedMeta.Info == "" {
			printer.Printf("<pre><code>%s</code></pre>\n", node.Text)
		} else {
			printer.Printf(`<pre><code class="language-%s">%s</code></pre>`+"\n",
				node.FencedMeta.Info, node.Text)
		}
	case TokenHTML:
	case TokenStrong:
	case TokenEmphasis:
	case TokenImage:
	case TokenLink:
	}
}
