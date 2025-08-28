package markdown

import (
	"strconv"
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
	ListMeta struct {
		Bullet string
		Start  int
		Index  int
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
	if scanner, token := ScanHeader(line); scanner != nil {
		heading := node.AddChild(NewNode(token, scanner))
		return heading.Scanner(line, heading)
	}
	if strings.HasPrefix(line, "> ") {
		quote := node.AddChild(NewNode(TokenBlockQuote, ScanBlockQuote))
		return quote.Scanner(line, quote)
	}
	if meta := parseListMeta(line); meta.Bullet != "" && meta.Indent == "" {
		list_ := node.AddChild(NewNode(meta.Token, ScanList(meta)))
		list_.ListMeta.Bullet = meta.Bullet
		return list_.Scanner(line, list_)
	}
	if strings.HasPrefix(line, "```") {
		codeBlock := node.AddChild(NewNode(TokenPreCode, ScanFencedCodeBlock))
		return codeBlock.Scanner(line, codeBlock)
	}
	if str.IsOnly(line, '-') && strings.Count(line, "-") >= 3 {
		thematicBreak := node.AddChild(NewNode(TokenThematicBreak, ScanThematicBreak))
		return thematicBreak.Scanner(line, thematicBreak)
	}
	paragraph := node.AddChild(NewNode(TokenParagraph, ScanParagraph))
	return paragraph.Scanner(line, paragraph)
}
func ScanHeader(line string) (Scanner, Token) {
	prefix2Header := map[string]Token{
		"# ":      TokenHeader1,
		"## ":     TokenHeader2,
		"### ":    TokenHeader3,
		"#### ":   TokenHeader4,
		"##### ":  TokenHeader5,
		"###### ": TokenHeader6,
	}
	for prefix, token := range prefix2Header {
		if strings.HasPrefix(line, prefix) {
			return func(line string, node *Node) *Node {
				node.Text, _ = strings.CutPrefix(line, prefix)
				return node.Parent
			}, token
		}
	}
	return nil, TokenNone
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
	// TODO: scan text for other containers (nested blockquotes, lists, etc.)
	return node
}
func ScanList(meta ListMeta) Scanner {
	return func(line string, node *Node) *Node {
		meta2 := parseListMeta(line)
		if meta2.Token == TokenNone {
			return node.Parent
		}
		if len(meta2.Indent) > len(meta.Indent) {
			sublist := node.AddChild(NewNode(meta2.Token, ScanList(meta2)))
			return sublist.Scanner(line, sublist)
		}
		if len(meta2.Indent) < len(meta.Indent) {
			node = node.Parent
			meta.Indent = meta2.Indent
		}
		if meta2.Token != meta.Token {
			newList := node.Parent.AddChild(NewNode(meta2.Token, ScanList(meta2)))
			newList.ListMeta.Bullet = meta.Bullet
			item := newList.AddChild(NewNode(TokenListItem, ScanListItem(meta2.Indent, meta2.Bullet)))
			return item.Scanner(line, item)
		}
		item := node.AddChild(NewNode(TokenListItem, ScanListItem(meta.Indent, meta.Bullet)))
		item.ListMeta.Bullet = meta.Bullet
		item.Scanner(line, item)
		return node
	}
}
func ScanListItem(indent, bullet string) Scanner {
	return func(line string, node *Node) *Node {
		node.Text = strings.TrimPrefix(line, indent+bullet) // Trim prefix (including numeric prefixes)
		return node.Parent
	}
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

type ListMeta struct {
	Indent string
	Bullet string
	Start  int
	Token  Token
}

func parseListMeta(line string) (result ListMeta) {
	result.Indent, line, _ = str.CutIndent(line)
	if strings.HasPrefix(line, "- ") {
		result.Bullet = "- "
		result.Token = TokenUnorderedList
		return result
	}
	beforeDot, _, ok := strings.Cut(line, ".")
	if !ok {
		return ListMeta{}
	}
	n, err := strconv.Atoi(beforeDot)
	if err != nil {
		return ListMeta{}
	}
	if n < 1 {
		return ListMeta{}
	}
	result.Start = n
	result.Bullet = ". "
	result.Token = TokenOrderedList
	return result
}
