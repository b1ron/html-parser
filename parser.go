package parser

import (
	"io"
	"text/scanner"
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

// holds a DOCTYPE token
type tokBuf struct {
	name []byte
}

type lexer struct {
	scanner.Scanner
}

type parser struct {
	tokBuf      tokBuf
	lex         *lexer
	state       state
	returnState *state
}

func NewParser(r io.Reader) *parser {
	p := parser{}
	p.lex = &lexer{}
	p.lex.Init(r)
	p.lex.Whitespace ^= 1 << ' '

	p.state = data // inital state
	p.tokBuf.name = make([]byte, 0)
	return &p
}

// Certain states also use a temporary buffer to track progress, and the character
// reference state uses a return state to return to the state it was invoked from.
func (p *parser) setReturnState(s state) {
	p.returnState = &s
}

func (p *parser) parse() {
	for {
		token := p.lex.Scan()
		if token == scanner.EOF {
			break
		}

		switch p.state {
		case data:
			switch token {
			case '&':
				p.setReturnState(data)
				p.state = characterReference
			case '<':
				p.state = tagOpen
			case '0':
				// TODO emit the current input character as a character token
			case scanner.EOF:
				// TODO emit an end-of-file token
			default:
				// TODO emit the current input character as a character token
			}
		case tagOpen:
			switch token {
			case '!':
				p.state = markupDeclarationOpen
			case '/':
				p.state = endTagOpen
			}
		case markupDeclarationOpen:
			// match for the word "DOCTYPE"
			switch p.lex.TokenText() {
			case "DOCTYPE":
				p.state = DOCTYPE
			}
		case DOCTYPE:
			switch token {
			case ' ':
				p.state = afterDOCTYPEName
			case '>':
				// reconsume in the before DOCTYPE name state.
			}
		case beforeDOCTYPEName:
			// consume the next input character
			// create a new DOCTYPE token. set the token's name to the current input character. switch to the DOCTYPE name state.
			p.tokBuf.name = append(p.tokBuf.name, byte(token))
			p.state = DOCTYPEName
		case afterDOCTYPEName:
			p.state = bogusDOCTYPE
		case bogusDOCTYPE:
			switch token {
			case '>':
				// switch to the data state. emit the DOCTYPE token.
				p.state = data
			}
		case DOCTYPEName:
			switch token {
			case ' ':
				p.state = afterDOCTYPEName
			case '>':
				// emit the current DOCTYPE token. switch to the data state.
				p.state = data
			default:
				p.tokBuf.name = append(p.tokBuf.name, byte(token))
			}
		}
	}
}

// var simple = `
// <!DOCTYPE html>
//   <html>
//   </html>
// `
