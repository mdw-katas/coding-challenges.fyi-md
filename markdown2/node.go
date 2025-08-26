package markdown2

import (
	"strings"

	"github.com/mdw-katas/coding-challenges.fyi-md/util/list"
	"github.com/mdw-katas/coding-challenges.fyi-md/util/str"
)

type Node struct {
	Token    Token
	Closed   bool
	Scanner  Scanner
	Text     string
	Parent   *Node
	Children *list.List[*Node]

	FencedMeta FencedCodeBlockMeta
}
type FencedCodeBlockMeta struct {
	Info    string
	Started bool
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

type Scanner func(line string, node *Node) *Node

func ScanBlocks(line string, node *Node) *Node {
	if len(line) == 0 {
		return node
	}
	if strings.HasPrefix(line, "# ") {
		heading := node.AddChild(NewNode(TokenH1, ScanH1))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "## ") {
		heading := node.AddChild(NewNode(TokenH2, ScanH2))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "### ") {
		heading := node.AddChild(NewNode(TokenH2, ScanH3))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "#### ") {
		heading := node.AddChild(NewNode(TokenH2, ScanH4))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "##### ") {
		heading := node.AddChild(NewNode(TokenH2, ScanH5))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "###### ") {
		heading := node.AddChild(NewNode(TokenH2, ScanH6))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "> ") {
		quote := node.AddChild(NewNode(TokenBlockQuote, ScanBlockQuote))
		return quote.Scanner(line, quote)
	}
	if strings.HasPrefix(line, "```") {
		codeBlock := node.AddChild(NewNode(TokenPreCode, ScanFencedCodeBlock))
		return codeBlock.Scanner(line, codeBlock)
	}
	if str.IsOnly(line, '-') && strings.Count(line, "-") >= 3 {
		thematicBreak := node.AddChild(NewNode(TokenHR, ScanThematicBreak))
		return thematicBreak.Scanner(line, thematicBreak)
	}
	paragraph := node.AddChild(NewNode(TokenParagraph, ScanParagraph))
	return paragraph.Scanner(line, paragraph)
}
func ScanH1(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "# ")
	return node.Parent
}
func ScanH2(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "## ")
	return node.Parent
}
func ScanH3(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "### ")
	return node.Parent
}
func ScanH4(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "#### ")
	return node.Parent
}
func ScanH5(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "##### ")
	return node.Parent
}
func ScanH6(line string, node *Node) *Node {
	node.Text, _ = strings.CutPrefix(line, "###### ")
	return node.Parent
}
func ScanThematicBreak(line string, node *Node) *Node {
	node.Closed = true
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
func ScanFencedCodeBlock(line string, node *Node) *Node {
	if !node.FencedMeta.Started {
		info, ok := strings.CutPrefix(line, "```")
		if ok {
			node.FencedMeta.Info = info
		}
		node.FencedMeta.Started = true
		return node
	}
	if node.FencedMeta.Started && strings.HasPrefix(line, "```") {
		node.Closed = true
		return node.Parent
	}
	node.Text += line + "\n"
	return node
}
func ScanParagraph(line string, node *Node) *Node {
	if len(line) == 0 {
		return node.Parent
	}
	if len(node.Text) > 0 {
		node.Text += "\n"
	}
	node.Text += line
	return node
}
