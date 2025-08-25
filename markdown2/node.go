package markdown2

import (
	"strings"

	"github.com/mdw-go/printing"
	"github.com/mdw-katas/coding-challenges.fyi-md/util/list"
)

type Node struct {
	Token    Token
	Closed   bool
	Scanner  Scanner
	Text     string
	Parent   *Node
	Children *list.List[*Node]
}

func NewNode(token Token, scanner Scanner) *Node {
	return &Node{
		Token:    token,
		Scanner:  scanner,
		Children: list.Of[*Node](),
	}
}

func (this *Node) AddChild(child *Node) *Node {
	child.Parent = this
	this.Children.Add(child)
	return child
}

func (this *Node) LocateOpenLeaf() *Node {
	if this.Children.Len() == 0 && !this.Closed {
		return this
	}
	if !this.Closed && this.Children.Len() > 0 {
		for child := range this.Children.All() {
			n := child.LocateOpenLeaf()
			if n != nil {
				return n
			}

		}
	}
	return nil
}

func (this *Node) Render(printer printing.Printer, level int) {
	printer.Print(strings.Repeat("  ", level))
	printer.Println(this.Token.String(), strings.ReplaceAll(this.Text, "\n", "\\n"))
	for child := range this.Children.All() {
		child.Render(printer, level+1)
	}
}

type Scanner func(line string, node *Node) *Node

func ScanBlocks(line string, node *Node) *Node {
	if len(line) == 0 {
		return node
	}
	if strings.HasPrefix(line, "# ") {
		heading := node.AddChild(NewNode(TokenH1, ScanH1))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "> ") {
		quote := node.AddChild(NewNode(TokenBlockQuote, ScanBlockQuote))
		return ScanBlockQuote(line, quote)
	}
	panic(line)
}
func ScanH1(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "# ")
	return node.Parent
}
func ScanBlockQuote(line string, node *Node) *Node {
	if len(line) == 0 {
		return node.Parent
	}
	content, _ := strings.CutPrefix(line, "> ")
	if len(node.Text) > 0 {
		node.Text += "\n"
	}
	node.Text += content
	return node
}
