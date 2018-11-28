
package parser

import (
	"fmt"
	"io"
	"bytes"
)

type FilePosition struct {
	Line int
	Column int
	File string
}

type Token struct {
	Text []byte
	TokenType *Tag
	Position FilePosition
}

func (tok *Token) String() string {
	return fmt.Sprintf("(%d, %d) %s %s",
		tok.Position.Line, tok.Position.Column,
		tok.TokenType.string,
		string(tok.Text))
}


func newToken(txt []byte, tt *Tag) *Token {
	return &Token {
		Text: txt,
		TokenType: tt,
		Position: FilePosition {
			Line: 0,
			Column: 0,
			File: "",
		},
	}

}

type Lexer struct {
	reader io.Reader
	charbuf []byte
	bufcount int
	fileError error

	pos FilePosition

	//FIXME remove
	//line int
	//column int
	//filename string
}

//FIXME NewLexer!
func MakeLexer(src io.Reader, filename string) *Lexer {

	ret := &Lexer{
		reader: src,
		charbuf: make([]byte, 16),
		bufcount: 0,
		fileError: nil,
		pos: FilePosition {
			Line: 1,
			Column: 1,
			File: filename,
		},
	}

	ret.fill()
	return ret
}

func (lex *Lexer) fill() {

	if lex.bufcount == 16 {
		return
	}

	if lex.fileError != nil {
		return
	}

	n, err := lex.reader.Read(lex.charbuf[lex.bufcount:])
	lex.bufcount += n
	if err != nil {
		lex.fileError = err
	}
}

func (lex *Lexer) consume(n int) {

	for i:=0; i<n; i++ {
		if lex.charbuf[0] == 10 {
			lex.pos.Line ++
			lex.pos.Column = 1
		} else {
			lex.pos.Column ++
		}
		copy(lex.charbuf[:15], lex.charbuf[1:16])
		lex.charbuf[15] = 0
	}
	lex.bufcount -= n

	lex.fill()
}

func (lex *Lexer) isEOF() bool {
	return lex.bufcount == 0
}

func (lex *Lexer) match(x string) bool {
	return bytes.HasPrefix(lex.charbuf, []byte(x))
}

func (lex *Lexer) matchByte(x byte) bool {
	return lex.charbuf[0] == x
}

func isSpace(x byte) bool {
	if x == ' ' || x == '\t' || x == 10 || x == 13 {
		return true
	}
	return false
}

func isDigit(x byte) bool {
	return x >= '0' && x <= '9'
}

func isNumberChar(x byte) bool {
	return ((x >= '0' && x <= '9') ||
		(x >= 'a' && x <= 'z') ||
		(x >= 'A' && x <= 'Z') ||
		(x=='_') || (x =='.'))
}

func isIdentifierStart(x byte) bool {
	return ((x >= 'a' && x <= 'z') ||
		(x >= 'A' && x <= 'Z') ||
		(x=='_') ||
		(x >= 128))
}

func isIdentifier(x byte) bool {
	return isDigit(x) || isIdentifierStart(x)
}

func (lex *Lexer) skipSpace() {
	for isSpace(lex.charbuf[0]) {
		lex.consume(1)
	}
}

func (lex *Lexer)  NextToken() *Token {

	lex.skipSpace()

	/* FIXME remove
	line := lex.line
	col := lex.column
	*/

	pos := lex.pos
	tok := lex.readToken()
	tok.Position = pos

	/* FIXME remove
	tok.Position.Line = line
	tok.Position.Column = col
	tok.Position.File = lex.filename
	*/

	EmitToken(tok.String())
	return tok
}

func (lex *Lexer) readToken() *Token {

	if lex.isEOF() {
		if lex.fileError != io.EOF {
			FatalError(&lex.pos, lex.fileError.Error())
		}
		return newToken([]byte("EOF"), EOF)
	}

	if lex.match("#{") {
		lex.consume(2)
		return newToken([]byte("#{"), PUNCTUATION)
	}

	for _, op := range Operators {
		if lex.match(op) {
			lex.consume(len(op))
			return newToken([]byte(op), OPERATOR)
		}
	}

	if lex.matchByte('#') {
		return lex.getComment()
	}

	if lex.match("\"\"\"") {
		return lex.getLongString()
	}

	if lex.matchByte('"') {
		return lex.getString()
	}

	if lex.matchByte('`') {
		return lex.getChar()
	}

	if isDigit(lex.charbuf[0]) {
		return lex.getNumber()
	}

	if isIdentifierStart(lex.charbuf[0]) {
		return lex.getSymbol()
	}

	ret := make([]byte, 1)
	ret[0] = lex.charbuf[0]
	lex.consume(1)
	return newToken(ret, PUNCTUATION)
}

func (lex *Lexer) getComment() *Token {

	buf := make([]byte, 0)

	for {
		if lex.matchByte(10) {
			break
		}

		if lex.isEOF() {
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return newToken(buf, COMMENT)
}

func (lex *Lexer) getString() *Token {
	buf := make([]byte, 0)

	lex.consume(1)

	for {
		if lex.matchByte(10) || lex.matchByte(13) || lex.isEOF() {
			FatalError(&lex.pos, "Newline in string constant")
		}

		if lex.matchByte('"') {
			lex.consume(1)
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return newToken(buf, STRING)
}

func (lex *Lexer) getLongString() *Token {
	buf := make([]byte, 0)

	lex.consume(3)

	for {
		if lex.isEOF() {
			FatalError(&lex.pos, "EOF in string constant")
		}

		if lex.match("\"\"\"") {
			lex.consume(3)
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return newToken(buf, STRING)
}

func (lex *Lexer) getChar() *Token {
	buf := make([]byte, 0)

	lex.consume(1)

	for {
		if lex.matchByte(10) || lex.matchByte(13) || lex.isEOF() {
			FatalError(&lex.pos, "Newline in character constant")
		}

		if lex.matchByte('`') {
			lex.consume(1)
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return newToken(buf, CHARACTER)
}

func (lex *Lexer) getNumber() *Token {
	var buf []byte = nil

	for isNumberChar(lex.charbuf[0]) {
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return newToken(buf, NUMBER)
}

func (lex *Lexer) getSymbol() *Token {
	var buf []byte = nil

	for isIdentifier(lex.charbuf[0]) {
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	if Keywords[string(buf)] {
		return newToken(buf, KEYWORD)
	} else {
		return newToken(buf, SYMBOL)
	}
}

func IsValidIdentifier(name string) bool {
	if len(name) == 0 {
		return false
	}

	if !isIdentifierStart(name[0]) {
		return false
	}

	for _,x := range([]byte(name)) {
		if !isIdentifier(x) {
			return false
		}
	}

	return true
}

