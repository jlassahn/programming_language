
package parser

import (
	"os"
	"fmt"
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

func (pe *parseElement) addChild(child ParseElement) {
	if child != nil {
		pe.children = append(pe.children, child)
	}
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

func (mp *mainParser) peek(pos int, etype int, txt string) bool {
	el := mp.queue[pos]
	if etype != el.ElementType() {
		return false
	}

	if txt != "" && txt != el.TokenString() {
		return false;
	}

	return true
}

func (mp *mainParser) match(etype int, txt string) ParseElement {
	if mp.peek(0, etype, txt) {
		return mp.consume()
	} else {
		var err string
		if txt == "" {
			err = fmt.Sprintf("expected %s got %s",
				lexer.TypeNames[etype],
				lexer.TypeNames[mp.queue[0].ElementType()])
		} else {
			err = fmt.Sprintf("expected %s %s got %s %s",
				lexer.TypeNames[etype], txt,
				lexer.TypeNames[mp.queue[0].ElementType()],
				mp.queue[0].TokenString())
		}
		mp.error(err)
		return nil
	}
}

func (mp *mainParser) tryMatch(etype int, txt string) bool {
	if mp.peek(0, etype, txt) {
		mp.consume()
		return true
	}
	return false
}

func (mp *mainParser) tryOperator(oplist map[string]bool) bool {
	if mp.queue[0].ElementType() != lexer.OPERATOR {
		return false
	}
	return oplist[mp.queue[0].TokenString()]
}


func (mp *mainParser) startElement(etype int) *parseElement {

	line, col := mp.queue[0].Position()

	ret := parseElement {
		nil,
		nil,
		etype,
		line, col,
		nil }

	return &ret
}

func (mp *mainParser) error(txt string) {
	line, col := mp.queue[0].Position()
	txt = fmt.Sprintf("at (%d, %d) %s", line, col, txt)
	EmitError(txt)
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

	ret := mp.startElement(lexer.SOURCE_FILE)

	if mp.peek(0, lexer.KEYWORD, "janus") {
		ret.addChild(mp.parseHeader())
	} else {
		mp.error("expected janus declaration as the first statement")
	}

	skipping_fail := false
	for !mp.peek(0, lexer.EOF, "") {
		el := mp.parseFileDeclaration()
		ret.addChild(el)

		if el == nil {
			if !skipping_fail {
				mp.error("invalid statement")
				skipping_fail = true
			}
			mp.consume()
		}
	}

	return ret
}

func (mp *mainParser) parseHeader() ParseElement {

	ret := mp.startElement(lexer.HEADER)
	ret.addChild(mp.match(lexer.KEYWORD, "janus"))
	ret.addChild(mp.match(lexer.NUMBER, ""))
	ret.addChild(mp.match(lexer.PUNCTUATION, ";"))

	return ret
}

func (mp *mainParser) parseFileDeclaration() ParseElement {

	if mp.peek(0, lexer.KEYWORD, "import") {
		return mp.parseImport()
	}

	if mp.peek(0, lexer.KEYWORD, "def") {
		return mp.parseDef()
	}

	//FIXME add
	//STRUCT
	//INTERFACE
	//METHOD
	//OPERATOR_DEF

	return nil
}

func (mp *mainParser) parseImport() ParseElement {

	ret := mp.startElement(lexer.IMPORT)
	mp.match(lexer.KEYWORD, "import")
	ret.addChild(mp.parseImportName())
	if mp.tryMatch(lexer.KEYWORD, "as") {
		ret.addChild(mp.parseImportName())
	}
	mp.match(lexer.PUNCTUATION, ";")
	return ret
}

func (mp *mainParser) parseImportName() ParseElement {

	if mp.peek(0, lexer.OPERATOR, ".") {
		return mp.match(lexer.OPERATOR, ".")
	}

	ret := mp.match(lexer.SYMBOL, "")

	for mp.peek(0, lexer.OPERATOR, ".") {
		ex := mp.startElement(lexer.EXPRESSION)
		ex.addChild(mp.match(lexer.OPERATOR, "."))
		ex.addChild(ret)
		ex.addChild(mp.match(lexer.SYMBOL, ""))
		ret = ex
	}
	return ret
}

func (mp *mainParser) parseDef() ParseElement {

	ret := mp.startElement(lexer.DEF)
	mp.match(lexer.KEYWORD, "def")
	ret.addChild(mp.match(lexer.SYMBOL, ""))

	if mp.tryMatch(lexer.OPERATOR, "=") {
		ret.addChild(mp.parseExpression())
		mp.match(lexer.PUNCTUATION, ";")
	} else {
		ret.addChild(mp.parseType())
		if mp.peek(0, lexer.PUNCTUATION, "{") {
			ret.addChild(mp.parseFunctionContent())
		} else {
			if mp.tryMatch(lexer.OPERATOR, "=") {
				ret.addChild(mp.parseExpression())
			}
			mp.match(lexer.PUNCTUATION, ";")
		}
	}
	return ret
}


func (mp *mainParser) parseType() ParseElement {

	if mp.tryMatch(lexer.PUNCTUATION, "(") {
		ret:= mp.startElement(lexer.FUNCTION_TYPE)
		ret.addChild(mp.parseParameterList())
		mp.match(lexer.PUNCTUATION, ")")
		mp.match(lexer.OPERATOR, "->")
		ret.addChild(mp.parseType())
		return ret
	} else {
		ret:= mp.startElement(lexer.TYPE)
		ret.addChild(mp.match(lexer.SYMBOL, ""))
		if mp.tryMatch(lexer.PUNCTUATION, "(") {
			ret.addChild(mp.parseTypeList())
			mp.match(lexer.PUNCTUATION, ")")
		}
		return ret
	}
}

func (mp *mainParser) parseParameterList() ParseElement {

	ret := mp.startElement(lexer.PARAMETER_LIST)
	for mp.peek(0, lexer.SYMBOL, "") {
		el := mp.startElement(lexer.PARAMETER)
		ret.addChild(el)
		el.addChild(mp.match(lexer.SYMBOL, ""))
		el.addChild(mp.parseType())

		if !mp.tryMatch(lexer.PUNCTUATION, ",") {
			break
		}
	}
	return ret
}

func (mp *mainParser) parseTypeList() ParseElement {

	ret := mp.startElement(lexer.TYPE_LIST)
	for {
		if mp.peek(0, lexer.NUMBER, "") ||
			mp.peek(0, lexer.CHARACTER, "") ||
			mp.peek(0, lexer.STRING, "") {
			ret.addChild(mp.consume())
		} else if mp.peek(0, lexer.SYMBOL, "") ||
			mp.peek(0, lexer.PUNCTUATION, "(") {
			ret.addChild(mp.parseType())
		} else {
			break
		}

		if !mp.tryMatch(lexer.PUNCTUATION, ",") {
			break
		}
	}
	return ret
}

func (mp *mainParser) parseFunctionContent() ParseElement {

	ret := mp.startElement(lexer.FUNCTION_CONTENT)

	mp.match(lexer.PUNCTUATION, "{")

	for !mp.peek(0, lexer.PUNCTUATION, "}") {
		el := mp.parseFunctionStatement()
		ret.addChild(el)
		if el == nil {
			break
		}
	}

	mp.match(lexer.PUNCTUATION, "}")
	return ret
}

func (mp *mainParser) parseFunctionStatement() ParseElement {

	if mp.peek(0, lexer.KEYWORD, "def") {
		return mp.parseDef() //FIXME this allows nested functions!
	}

	//FIXME maybe refactor into mp.parseReturn
	if mp.peek(0, lexer.KEYWORD, "return") {
		ret := mp.startElement(lexer.EXPRESSION)
		ret.addChild(mp.consume())
		if !mp.peek(0, lexer.PUNCTUATION, ";") {
			ret.addChild(mp.parseExpression())
		}
		mp.match(lexer.PUNCTUATION, ";")
		return ret
	}

	//FIXME more keyword statements

	return mp.parseAssignmentStatement()
}

func (mp *mainParser) parseAssignmentStatement() ParseElement {
	lhs := mp.parseExpression()
	if mp.tryMatch(lexer.OPERATOR, "=") { //FIXME other assignment operators
		rhs := mp.parseExpression()
		mp.match(lexer.PUNCTUATION, ";")

		ret := mp.startElement(lexer.ASSIGNMENT)
		ret.addChild(lhs)
		ret.addChild(rhs)
		return ret
	}
	mp.match(lexer.PUNCTUATION, ";")
	return lhs
}

func (mp *mainParser) parseExpression() ParseElement {
	ret := mp.parseExpressionAnd()
	for mp.tryOperator(lexer.OrOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionAnd())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionAnd() ParseElement {
	ret := mp.parseExpressionCompare()
	for mp.tryOperator(lexer.AndOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionCompare())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionCompare() ParseElement {
	ret := mp.parseExpressionAdd()
	for mp.tryOperator(lexer.CompareOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionAdd())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionAdd() ParseElement {
	ret := mp.parseExpressionMult()
	for mp.tryOperator(lexer.AddOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionMult())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionMult() ParseElement {
	ret := mp.parseExpressionExp()
	for mp.tryOperator(lexer.MultOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionExp())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionExp() ParseElement {
	ret := mp.parseExpressionPrefix()
	for mp.tryOperator(lexer.ExpOperators) {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionPrefix())
		ret = el
	}
	return ret
}

func (mp *mainParser) parseExpressionPrefix() ParseElement {
	//FIXME implement
	return mp.parseExpressionSuffix()
}

func (mp *mainParser) parseExpressionSuffix() ParseElement {

	var ret ParseElement

	if mp.peek(0, lexer.SYMBOL, "") {
		ret = mp.parseExpressionDot()
	} else if mp.peek(0, lexer.NUMBER, "") {
		ret = mp.consume()
	} else if mp.peek(0, lexer.STRING, "") {
		ret = mp.consume()
	} else if mp.peek(0, lexer.CHARACTER, "") {
		ret = mp.consume()
	} else if mp.tryMatch(lexer.PUNCTUATION, "(") {
		ret = mp.parseExpression()
		mp.match(lexer.PUNCTUATION, ")")
	}

	for {
		if mp.tryOperator(lexer.SuffixOperators) {
			el := mp.startElement(lexer.EXPRESSION)
			el.addChild(mp.consume())
			el.addChild(ret)
			ret = el
		} else if mp.tryMatch(lexer.PUNCTUATION, "[") {
			el := mp.startElement(lexer.INDEX)
			el.addChild(ret)
			mp.match(lexer.PUNCTUATION, "]")
		} else if mp.tryMatch(lexer.PUNCTUATION, "(") {
			el := mp.startElement(lexer.CALL)
			el.addChild(ret)
			ret = el
			if !mp.tryMatch(lexer.PUNCTUATION, ")") {
				el.addChild(mp.parseExpression())
				for mp.tryMatch(lexer.PUNCTUATION, ",") {
					el.addChild(mp.parseExpression())
				}
				mp.match(lexer.PUNCTUATION, ")")
			}
		} else {
			break
		}
	}

	return ret
}

func (mp *mainParser) parseExpressionDot() ParseElement {

	ret := mp.match(lexer.SYMBOL, "")
	//FIXME why is . an operator?
	for mp.peek(0, lexer.OPERATOR, ".") {
		el := mp.startElement(lexer.EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.match(lexer.SYMBOL, ""))
		ret = el
	}
	return ret
}


/*
EXPRESSION_EXP:
EXPRESSION_PREFIX:
EXPRESSION_SUFFIX:
EXPRESSION_DOT:
EXPRESSION_LIST:

// func (mp *mainParser) tryOperator(oplist map[string]bool) bool {
*/

