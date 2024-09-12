// package parser implements an HTML parser and state machine

package parser

import (
	"bufio"
	"fmt"
	"io"
)

type state int

// HTML states https://html.spec.whatwg.org/#data-state
const (
	data state = iota
	RCDATA
	RAWTEXT
	scriptData
	PLAINTEXT
	tagOpen
	endTagOpen
	tagName
	RCDATALessThanSign
	RCDATAEndTagOpen
	RCDATAEndTagName
	RAWTEXTLessThanSign
	RAWTEXTEndTagOpen
	RAWTEXTEndTagName
	scriptDataLessThanSign
	scriptDataEndTagOpen
	scriptDataEndTagName
	scriptDataEscapeStart
	scriptDataEscapeStartDash
	scriptDataEscaped
	scriptDataEscapedDash
	scriptDataEscapedDashDash
	scriptDataEscapedLessThanSign
	scriptDataEscapedEndTagOpen
	scriptDataEscapedEndTagName
	scriptDataDoubleEscapeStart
	scriptDataDoubleEscaped
	scriptDataDoubleEscapedDash
	scriptDataDoubleEscapedDashDash
	scriptDataDoubleEscapedLessThanSign
	scriptDataDoubleEscapeEnd
	beforeAttributeName
	attributeName
	afterAttributeName
	beforeAttributeValue
	attributeValueDoubleQuoted
	attributeValueSingleQuoted
	attributeValueUnquoted
	afterAttributeValueQuoted
	selfClosingStartTag
	bogusComment
	markupDeclarationOpen
	commentStart
	commentStartDash
	comment
	commentEndDash
	commentEnd
	commentEndBang
	DOCTYPE
	beforeDOCTYPEName
	DOCTYPEName
	afterDOCTYPEName
	afterDOCTYPEPublicKeyword
	beforeDOCTYPEPublicIdentifier
	DOCTYPEPublicIdentifierDoubleQuoted
	DOCTYPEPublicIdentifierSingleQuoted
	afterDOCTYPEPublicIdentifier
	betweenDOCTYPEPublicAndSystemIdentifiers
	afterDOCTYPESystemKeyword
	beforeDOCTYPESystemIdentifier
	DOCTYPESystemIdentifierDoubleQuoted
	DOCTYPESystemIdentifierSingleQuoted
	afterDOCTYPESystemIdentifier
	bogusDOCTYPE
	CDATASection
	CDATASectionBracket
	CDATASectionEnd
	characterReference
	namedCharacterReference
	ambiguousAmpersand
	numericCharacterReference
	hexadecimalCharacterReferenceStart
	decimalCharacterReferenceStart
	hexadecimalCharacterReference
	decimalCharacterReference
	numericCharacterReferenceEnd
)

type mode int

// HTML insertion modes https://html.spec.whatwg.org/#insertion-mode
const (
	initial mode = iota
	beforeHTML
	beforeHead
	inHead
	inHeadNoscript
	afterHead
	inBody
	text
	inTable
	inTableText
	inCaption
	inColumnGroup
	inTableBody
	inRow
	inCell
	inSelect
	inSelectInTable
	inTemplate
	afterBody
	inFrameset
	afterFrameset
	afterAfterBody
	afterAfterFrameset
)

const EOF = -1

// tree represents an HTML document's DOM tree
type tree struct {
	root *listNode
}

// scanner represents a lexical scanner
type scanner struct {
	r *bufio.Reader
}

type parser struct {
	s     *scanner
	state state
	mode  mode
}

func (s *scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return EOF
	}
	return r
}

func (s *scanner) unread() { _ = s.r.UnreadRune() }

func (s *scanner) scan() (tok rune) { return s.read() }

// scanIdent consumes the current rune and all contiguous ident runes
func (s *scanner) scanIdent() (lit string) {
	for {
		ch := s.read()
		if ch == EOF {
			break
		}
		if isDelim(ch) {
			s.unread()
			break
		}
		lit += string(ch)
	}
	return
}

// isDelim returns true for tokens corresponding to state transitions or parsing actions
func isDelim(ch rune) bool {
	return ch == '<' || ch == '>' || ch == '/' || ch == '&' || ch == '!' || ch == ' '
}

// newParser returns a new instance of parser
func newParser(r io.Reader) *parser {
	return &parser{s: &scanner{r: bufio.NewReader(r)}, state: data, mode: initial}
}

// parse parses the input
func (p *parser) parse() *tree {
	// ...
	t := &tree{}
	t.root = t.newList()
	for {
		token := p.s.scan()
		if token == EOF {
			break
		}

		switch p.state {
		case data:
			switch token {
			case '<':
				p.state = tagOpen
			}
		case tagOpen:
			switch token {
			case '!':
				p.state = markupDeclarationOpen
			case '/':
				p.state = endTagOpen
			}
		case markupDeclarationOpen:
			p.s.unread()
			switch p.s.scanIdent() {
			case "DOCTYPE":
				p.state = DOCTYPE
			}
		case DOCTYPE:
			switch token {
			case ' ':
				p.state = beforeDOCTYPEName
			case '>':
				// TODO reconsume in beforeDOCTYPEName state
				p.state = beforeDOCTYPEName
			}
		case beforeDOCTYPEName:
			switch token {
			case ' ':
				// ignore
				continue
			}
			p.s.unread()
			p.state = DOCTYPEName
		case DOCTYPEName:
			switch token {
			case ' ':
				p.state = afterDOCTYPEName
			case '>':
				p.state = data
				continue
			}
			p.s.unread()
			t.root.append(&elementNode{data: p.s.scanIdent()})
		}
	}
	fmt.Println(t.root.next().Data())
	return nil
}
