package parser

import (
	"io"
	"strings"
	"text/scanner"
	"unicode"
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

type lexer struct {
	scanner.Scanner
}

type parser struct {
	b           []byte // holds the current token in certain states
	lex         *lexer
	state       state
	returnState *state
}

func NewParser(r io.Reader) *parser {
	p := parser{}
	p.b = make([]byte, 0)
	p.lex = &lexer{}
	p.lex.Init(r)
	p.state = data // inital state
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
				p.state = beforeDOCTYPEName
			case '>':
				// reconsume in the before DOCTYPE name state.
			}
		case beforeDOCTYPEName:
			for strings.ContainsRune(" \n\t", token) {
				// ignore the character
				token = p.lex.Scan()
			}
			switch {
			case unicode.IsUpper(token):
				// create a new DOCTYPE token, set its name to the lowercase version of the current input character
				p.state = DOCTYPEName
			default: // to handle anything else
				// create a new DOCTYPE token. set the token's name to the current input character. switch to the DOCTYPE name state.
				p.b = append(p.b, byte(token))
				p.state = DOCTYPEName
			}
		case DOCTYPEName:
			switch token {
			case ' ':
				p.state = afterDOCTYPEName
			case '>':
				p.state = data
			default:
				p.b = append(p.b, byte(token))
			}
		}
	}
}
