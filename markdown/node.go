package markdown

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

	FencedMeta struct {
		Info    string
		Started bool
	}
	OrderedMeta struct {
		StartingNumber int
		Loose          bool
	}
	UnorderedMeta struct {
		Bullet rune
		Loose  bool
	}
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
	if strings.HasPrefix(line, "- ") {
		unordered := node.AddChild(NewNode(TokenUnorderedList, ScanUnorderedList("", '-')))
		return unordered.Scanner(line, unordered)
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

func ScanUnorderedList(indent string, bullet rune) Scanner {
	return func(line string, node *Node) *Node {
		prefix, _, _ := str.CutIndent(line)
		if len(prefix) > len(indent) {
			sublist := node.AddChild(NewNode(TokenUnorderedList, ScanUnorderedList(prefix, bullet)))
			return sublist.Scanner(line, sublist)
		}
		if !strings.HasPrefix(line, indent+string(bullet)+" ") {
			// TODO: we are losing the first item of new lists that come right after the last item in the current list
			return node.Parent
		}
		item := node.AddChild(NewNode(TokenListItem, ScanDashedListItem(indent, bullet)))
		item.UnorderedMeta.Bullet = bullet
		item.Scanner(line, item)
		return node
	}
}

func ScanDashedListItem(indent string, bullet rune) Scanner {
	return func(line string, node *Node) *Node {
		node.Text = strings.TrimPrefix(line, indent+string(bullet)+" ")
		return node.Parent
	}
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
func ScanThematicBreak(line string, node *Node) *Node {
	node.Closed = true
	return node.Parent
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
