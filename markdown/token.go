package markdown

type Token int

const (
	TokenNone Token = iota
	TokenDocument
	TokenHeader1
	TokenHeader2
	TokenHeader3
	TokenHeader4
	TokenHeader5
	TokenHeader6
	TokenThematicBreak
	TokenOrderedList
	TokenUnorderedList
	TokenListItem
	TokenBlockQuote
	TokenPreCode
	TokenParagraph // TODO: parse inlines
	TokenHTML      // TODO
	TokenStrong    // TODO: parse
	TokenEmphasis  // TODO: parse
	TokenImage     // TODO: parse
	TokenLink      // TODO: parse
)

func (t Token) String() string {
	switch t {
	case TokenDocument:
		return "Document"
	case TokenHeader1:
		return "H1"
	case TokenHeader2:
		return "H2"
	case TokenHeader3:
		return "H3"
	case TokenHeader4:
		return "H4"
	case TokenHeader5:
		return "H5"
	case TokenHeader6:
		return "H6"
	case TokenThematicBreak:
		return "HR"
	case TokenOrderedList:
		return "OrderedList"
	case TokenUnorderedList:
		return "UnorderedList"
	case TokenListItem:
		return "ListItem"
	case TokenBlockQuote:
		return "BlockQuote"
	case TokenParagraph:
		return "Paragraph"
	case TokenPreCode:
		return "PreCode"
	case TokenHTML:
		return "HTML"
	case TokenStrong:
		return "Strong"
	case TokenEmphasis:
		return "Emphasis"
	case TokenImage:
		return "Image"
	case TokenLink:
		return "Link"
	default:
		return "<unknown token>"
	}
}
