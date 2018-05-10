
package lexer

import (
	"fmt"
	"io"
	"bytes"
)

const (
	ERROR = iota
	EOF
	COMMENT
	NUMBER
	STRING
	SYMBOL
	OPERATOR
)

var type_names = map[int]string {
	ERROR: "ERROR",
	EOF: "EOF",
	COMMENT: "COMMENT",
	NUMBER: "NUMBER",
	STRING: "STRING",
	SYMBOL: "SYMBOL",
	OPERATOR: "OPERATOR" }

//FIXME complete list
//this must have longer operators first, only multichar operators needed here
var operators = []string {
	"->",
	">>",
	"<<",
	"++",
	"--",
	"#{" }

type Token struct {
	Text []byte
	TokenType int
}

func (tok *Token) String() string {
	return fmt.Sprintf("%s %s", type_names[tok.TokenType], string(tok.Text))
}

type Lexer struct {
	reader io.Reader
	charbuf []byte
	bufcount int
	file_error error
}

func MakeLexer(src io.Reader) *Lexer {

	ret := &Lexer{
		src,
		make([]byte, 16),
		0,
		nil }

	ret.fill()
	return ret
}

func (lex *Lexer) fill() {

	if lex.bufcount == 16 {
		return
	}

	if lex.file_error != nil {
		return
	}

	n, err := lex.reader.Read(lex.charbuf[lex.bufcount:])
	lex.bufcount += n
	if err != nil {
		lex.file_error = err
	}
}

func (lex *Lexer) consume(n int) {
	if n <= 0 {
		return
	}

	copy(lex.charbuf[:lex.bufcount - n], lex.charbuf[n:lex.bufcount])
	lex.bufcount -= n
	lex.charbuf[lex.bufcount] = 0

	lex.fill()
}

func (lex *Lexer) is_eof() bool {
	return lex.bufcount == 0
}

func (lex *Lexer) match(x string) bool {
	return bytes.HasPrefix(lex.charbuf, []byte(x))
}

func (lex *Lexer) match_byte(x byte) bool {
	return lex.charbuf[0] == x
}

func is_space(x byte) bool {
	if x == ' ' || x == '\t' || x == 10 || x == 13 {
		return true
	}
	return false
}

func is_digit(x byte) bool {
	return x >= '0' && x <= '9'
}

func is_number_char(x byte) bool {
	return ((x >= '0' && x <= '9') ||
		(x >= 'a' && x <= 'z') ||
		(x >= 'A' && x <= 'Z') ||
		(x=='_') || (x =='.'))
}

func is_identifier_start(x byte) bool {
	return ((x >= 'a' && x <= 'z') ||
		(x >= 'A' && x <= 'Z') ||
		(x=='_') ||
		(x >= 128))
}

func is_identifier(x byte) bool {
	return is_digit(x) || is_identifier_start(x)
}

func (lex *Lexer) skip_space() {
	for is_space(lex.charbuf[0]) {
		lex.consume(1)
	}
}

func (lex *Lexer)  NextToken() *Token {

	lex.skip_space()

	tok := lex.read_token()

	return tok
}

func (lex *Lexer) read_token() *Token {

	if lex.is_eof() {
		if lex.file_error == io.EOF {
			return &Token{[]byte("EOF"), EOF}
		} else {
			return &Token{ []byte(lex.file_error.Error()), ERROR }
		}
	}

	for _, op := range operators {
		if lex.match(op) {
			lex.consume(len(op))
			return &Token{[]byte(op), OPERATOR}
		}
	}

	if lex.match_byte('#') {
		return lex.get_comment()
	}

	if lex.match("\"\"\"") {
		return lex.get_long_string()
	}

	if lex.match_byte('"') {
		return lex.get_string()
	}

	if is_digit(lex.charbuf[0]) {
		return lex.get_number()
	}

	if is_identifier_start(lex.charbuf[0]) {
		return lex.get_symbol()
	}

	ret := make([]byte, 1)
	ret[0] = lex.charbuf[0]
	lex.consume(1)
	return &Token{ret, OPERATOR}
}

func (lex *Lexer) get_comment() *Token {

	buf := make([]byte, 0)

	for {
		if lex.match_byte(10) {
			break
		}

		if lex.is_eof() {
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return &Token{ buf, COMMENT }
}

func (lex *Lexer) get_string() *Token {
	buf := make([]byte, 0)

	lex.consume(1)

	for {
		if lex.match_byte(10) || lex.match_byte(13) || lex.is_eof() {
			return &Token{ []byte("Newline in string constant"), ERROR }
		}

		if lex.match_byte('"') {
			lex.consume(1)
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return &Token{ buf, STRING }
}

func (lex *Lexer) get_long_string() *Token {
	buf := make([]byte, 0)

	lex.consume(3)

	for {
		if lex.is_eof() {
			return &Token{ []byte("EOF in string constant"), ERROR }
		}

		if lex.match("\"\"\"") {
			lex.consume(3)
			break
		}
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return &Token{ buf, STRING }
}

func (lex *Lexer) get_number() *Token {
	var buf []byte = nil

	for is_number_char(lex.charbuf[0]) {
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return &Token{ buf, NUMBER }
}

func (lex *Lexer) get_symbol() *Token {
	var buf []byte = nil

	for is_identifier(lex.charbuf[0]) {
		buf = append(buf, lex.charbuf[0])
		lex.consume(1)
	}

	return &Token{ buf, SYMBOL }
}

