package markdown

type Token int

const (
	TokenDocument Token = iota
	TokenH1
	TokenH2
	TokenH3
	TokenH4
	TokenH5
	TokenH6
	TokenHR
	TokenOrderedList
	TokenUnorderedList
	TokenListItem
	TokenBlockQuote
	TokenParagraph
	TokenPreCode
	TokenHTML
	TokenStrong
	TokenEmphasis
	TokenImage
	TokenLink
)

func (t Token) String() string {
	switch t {
	case TokenDocument:
		return "Document"
	case TokenH1:
		return "H1"
	case TokenH2:
		return "H2"
	case TokenH3:
		return "H3"
	case TokenH4:
		return "H4"
	case TokenH5:
		return "H5"
	case TokenH6:
		return "H6"
	case TokenHR:
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
