
package parser

import (
	"os"
	"lexer"
)

type ParseElement interface {
	Children() []ParseElement
	Comments() []ParseElement
	ElementType() int
	Position() (int, int)
	Token() *lexer.Token
	TokenString() string
}

type Parser interface {
	GetElement() ParseElement
}

func EmitError(err string) {
	os.Stderr.WriteString(err + "\n")
}


const (
	SOURCE_FILE = iota + lexer.END_OF_TOKEN_TYPES + 1
	HEADER
)

type parseElement struct {
	children []ParseElement
	comments []ParseElement
	elementType int
	line, column int
	token *lexer.Token
}

func (pe *parseElement) Children() []ParseElement {
	return pe.children
}

func (pe *parseElement) Comments() []ParseElement {
	return pe.comments
}

func (pe *parseElement) ElementType() int {
	return pe.elementType
}

func (pe *parseElement) Position() (int, int) {
	return pe.line, pe.column
}

func (pe *parseElement) Token() *lexer.Token {
	return pe.token
}

func (pe *parseElement) TokenString() string {
	if pe.token != nil {
		return string(pe.token.Text)
	}
	return ""
}



type tokenWrapper struct {
	lex *lexer.Lexer
}

func (tw *tokenWrapper) GetElement() ParseElement {
	tok := tw.lex.NextToken()

	ret := parseElement{
		nil, nil,
		tok.TokenType,
		tok.Line,
		tok.Column,
		tok }

	return &ret
}

type commentMerger struct {
	upstream Parser
	comments []ParseElement
	depth int
}

func (cm *commentMerger) GetElement() ParseElement {

	for {
		el := cm.upstream.GetElement()

		if el.ElementType() == lexer.EOF {
			if cm.depth > 0 {
				EmitError("EOF inside block comment")
			}
			if cm.comments != nil {
				comments := cm.comments
				cm.comments = nil
				line, col := el.Position()
				return &parseElement {
					el.Children(),
					comments,   //FIXME merge existing comments
					el.ElementType(),
					line, col,
					el.Token() }
			} else {
				return el
			}
		}

		if cm.depth == 0 && el.ElementType() != lexer.COMMENT {
			if el.ElementType() == lexer.PUNCTUATION && el.TokenString() == "#{" {
				cm.comments = append(cm.comments, el)
				cm.depth = 1
			} else if cm.comments != nil {
				comments := cm.comments
				cm.comments = nil
				line, col := el.Position()
				return &parseElement {
					el.Children(),
					comments,   //FIXME merge existing comments
					el.ElementType(),
					line, col,
					el.Token() }
			} else {
				return el
			}
		} else {
			cm.comments = append(cm.comments, el)
			if el.ElementType() == lexer.PUNCTUATION {
				if el.TokenString() == "#{" {
					cm.depth ++
				}
				if el.TokenString() == "{" {
					cm.depth ++
				}
				if el.TokenString() == "}" {
					cm.depth --
				}
			}
		}
	}
}

type mainParser struct {
	upstream Parser

	queue []ParseElement
}

func (mp *mainParser) consume() ParseElement {
	ret := mp.queue[0]
	copy(mp.queue[:15], mp.queue[1:16])
	mp.queue[15] = mp.upstream.GetElement()

	return ret
}

func (mp *mainParser) GetElement() ParseElement {

	return mp.parseFile()
}

func NewParser(lex *lexer.Lexer) Parser {

	ret := &mainParser{
		&commentMerger{ &tokenWrapper{ lex }, nil, 0 },
		make([]ParseElement, 16) }

	for i:=0; i<16; i++ {
		ret.queue[i] = ret.upstream.GetElement()
	}

	return ret
}

func (mp *mainParser) parseFile() ParseElement {
	return mp.consume() //FIXME fake
}

