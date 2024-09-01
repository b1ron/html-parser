package parser

import (
	"container/list"
	"io"
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
	currentToken list.List
	lex          *lexer
	state        state
	returnState  *state
}

func newParser(r io.Reader) *parser {
	p := parser{}
	p.lex = &lexer{}
	p.lex.Init(r)
	p.lex.Whitespace ^= 1 << ' '

	p.state = data // inital state
	p.currentToken = *list.New()
	return &p
}

// holds a DOCTYPE, start tag or end tag token
type tokBuf struct {
	name []rune
}

func (p *parser) create(token rune) {
	p.currentToken.PushBack(tokBuf{name: []rune{token}})
}

func (p *parser) append(token rune) {
	e := p.currentToken.Back()
	curr := p.currentToken.Remove(e).(tokBuf)
	curr.name = append(curr.name, token)
	p.currentToken.PushBack(curr)
}

func (p *parser) emit() {
	// TODO emit the current DOCTYPE token
}

func (p *parser) reconsume() {
	// TODO reconsume the current input character in the DOCTYPE name state
}

// TODO figure out how to handle the return state in the parser
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
			// create a new start tag token, set its tag name to the empty string. reconsume in the tag name state
			if unicode.IsLetter(token) {
				p.create('0')
				p.state = tagName
			}
		case tagName:
			switch token {
			case ' ':
				p.state = beforeAttributeName
			case '/':
				// TODO switch to the self-closing start tag state
			case '>':
				// TODO emit the current tag token
				p.state = data
			default:
				// anything else append the current input character to the current tag token's tag name
				p.append(token)
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
				// reconsume in the before DOCTYPE name state
			}
		case beforeDOCTYPEName:
			// create a new DOCTYPE token. set the token's name to the current input character. switch to the DOCTYPE name state
			p.create(token)
			p.state = DOCTYPEName
		case afterDOCTYPEName:
			p.state = bogusDOCTYPE
		case bogusDOCTYPE:
			switch token {
			case '>':
				// switch to the data state. emit the DOCTYPE token
				p.state = data
			case scanner.EOF:
				// TODO emit an end-of-file token
			default:
				// anything else ignore the character
				p.lex.Scan()
			}
		case DOCTYPEName:
			switch token {
			case ' ':
				p.state = afterDOCTYPEName
			case '>':
				// emit the current DOCTYPE token. switch to the data state
				p.state = data
			default:
				// anything else append the current input character to the current tag token's tag name
				p.append(token)
			}
		}
	}
}

// var simple = `
// <!DOCTYPE html>
//   <html>
//   </html>
// `
