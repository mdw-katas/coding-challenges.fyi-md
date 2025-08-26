package markdown2

type Phase1Parser struct {
	root *Node
	node *Node
}

func NewPhase1Parser() *Phase1Parser {
	root := NewNode(TokenDocument, ScanBlocks)
	return &Phase1Parser{
		root: root,
		node: root,
	}
}

func (this *Phase1Parser) Feed(line string) {
	this.node = this.node.Scanner(line, this.node)
}

func (this *Phase1Parser) Finalize() {
	// TODO: further parsing of container nodes (block quotes, lists)
	// TODO: close all open nodes
}
