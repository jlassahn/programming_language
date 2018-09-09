
package parser

import (
	"fmt"
	"lexer"
	"output"
)

type ParseElement interface {
	Children() []ParseElement
	Comments() []ParseElement
	ElementType() *lexer.Tag
	Position() (int, int)
	Token() *lexer.Token
	TokenString() string
}

type Parser interface {
	GetElement() ParseElement
}

type parseElement struct {
	children []ParseElement
	comments []ParseElement
	elementType *lexer.Tag
	line, column int
	token *lexer.Token
}

func (pe *parseElement) Children() []ParseElement {
	return pe.children
}

func (pe *parseElement) Comments() []ParseElement {
	return pe.comments
}

func (pe *parseElement) ElementType() *lexer.Tag {
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
				line, col := cm.comments[0].Position()
				output.FatalError(line, col, "EOF inside block comment")
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
	progress bool
	resync bool
}

func (mp *mainParser) consume() ParseElement {
	ret := mp.queue[0]
	copy(mp.queue[:15], mp.queue[1:16])
	mp.queue[15] = mp.upstream.GetElement()
	mp.progress = true

	return ret
}

func (mp *mainParser) peek(pos int, etype *lexer.Tag, txt string) bool {
	el := mp.queue[pos]
	if etype != el.ElementType() {
		return false
	}

	if txt != "" && txt != el.TokenString() {
		return false;
	}

	return true
}

func (mp *mainParser) match(etype *lexer.Tag, txt string) ParseElement {
	if mp.peek(0, etype, txt) {
		mp.resync = false
		return mp.consume()
	} else {
		var err string
		if txt == "" {
			err = fmt.Sprintf("expected %s got %s %s",
				etype,
				mp.queue[0].ElementType(),
				mp.queue[0].TokenString())
		} else {
			err = fmt.Sprintf("expected %s %s got %s %s",
				etype, txt,
				mp.queue[0].ElementType(),
				mp.queue[0].TokenString())
		}
		mp.error(err)
		return nil
	}
}

func (mp *mainParser) checkProgress() {
	if mp.progress {
		mp.progress = false
	} else {
		err := fmt.Sprintf("unexpected %s %s",
			mp.queue[0].ElementType(),
			mp.queue[0].TokenString())
		mp.error(err)
		mp.consume()
	}
}

func (mp *mainParser) tryMatch(etype *lexer.Tag, txt string) bool {
	if mp.peek(0, etype, txt) {
		mp.resync = false
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


func (mp *mainParser) startElement(etype *lexer.Tag) *parseElement {

	line, col := mp.queue[0].Position()
	comments := mp.queue[0].Comments()

	ret := parseElement {
		nil,
		comments,
		etype,
		line, col,
		nil }

	return &ret
}

func (mp *mainParser) error(txt string) {

	if !mp.resync {
		line, col := mp.queue[0].Position()
		output.Error(line, col, txt)
		mp.resync = true
	}
}

func (mp *mainParser) GetElement() ParseElement {

	return mp.parseFile()
}

func NewParser(lex *lexer.Lexer) Parser {

	ret := &mainParser{
		&commentMerger{ &tokenWrapper{ lex }, nil, 0 },
		make([]ParseElement, 16),
		true,
		false }

	for i:=0; i<16; i++ {
		ret.queue[i] = ret.upstream.GetElement()
	}

	return ret
}

func (mp *mainParser) parseFile() ParseElement {

	ret := mp.startElement(lexer.SOURCE_FILE)

	ret.addChild(mp.parseHeader())

	for !mp.peek(0, lexer.EOF, "") {
		ret.addChild(mp.parseFileDeclaration())
		mp.checkProgress()
	}

	return ret
}

func (mp *mainParser) parseHeader() ParseElement {

	mp.match(lexer.KEYWORD, "janus")
	ret := mp.startElement(lexer.HEADER)
	ret.addChild(mp.match(lexer.NUMBER, ""))
	if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret.addChild(mp.parseHeaderOptions())
		mp.match(lexer.PUNCTUATION, "}")
	} else {
		mp.match(lexer.PUNCTUATION, ";")
	}

	return ret
}

func (mp *mainParser) parseHeaderOptions() ParseElement {
	ret := mp.startElement(lexer.LIST)
	for !mp.peek(0, lexer.PUNCTUATION, "}") {
		ret.addChild(mp.parseHeaderOption())
		mp.checkProgress()
	}
	return ret
}

func (mp *mainParser) parseHeaderOption() ParseElement {
	ret := mp.startElement(lexer.ASSIGNMENT)
	ret.addChild(mp.parseExpressionDot())
	mp.match(lexer.OPERATOR, "=")
	ret.addChild(mp.parseExpression())
	mp.match(lexer.PUNCTUATION, ";")
	return ret
}

func (mp *mainParser) parseFileDeclaration() ParseElement {

	if mp.peek(0, lexer.KEYWORD, "import") {
		return mp.parseImport()
	}

	if mp.peek(0, lexer.KEYWORD, "def") {
		return mp.parseDef()
	}

	if mp.peek(0, lexer.KEYWORD, "struct") {
		return mp.parseStruct()
	}
	if mp.peek(0, lexer.KEYWORD, "m_struct") {
		return mp.parseStruct()
	}

	if mp.peek(0, lexer.KEYWORD, "interface") {
		return mp.parseInterface()
	}

	if mp.peek(0, lexer.KEYWORD, "method") {
		return mp.parseMethod()
	}

	if mp.peek(0, lexer.KEYWORD, "operator") {
		return mp.parseOperator()
	}

	if mp.tryMatch(lexer.PUNCTUATION, ";") {
		// empty statement
		return nil
	}

	return nil
}

func (mp *mainParser) parseImport() ParseElement {

	ret := mp.startElement(lexer.IMPORT)
	mp.match(lexer.KEYWORD, "import")
	ret.addChild(mp.parseExpressionDot())

	if mp.tryMatch(lexer.OPERATOR, "=") {
		if mp.peek(0, lexer.OPERATOR, ".") {
			ret.addChild(mp.match(lexer.OPERATOR, "."))
		} else {
			ret.addChild(mp.parseExpressionDot())
		}
	}
	mp.match(lexer.PUNCTUATION, ";")
	return ret
}

func (mp *mainParser) parseDef() ParseElement {

	ret := mp.startElement(lexer.DEF)
	mp.match(lexer.KEYWORD, "def")
	ret.addChild(mp.match(lexer.SYMBOL, ""))

	if mp.peek(0, lexer.PUNCTUATION, "(") {
		ret.addChild(mp.parseFunctionType())
	} else if !mp.peek(0, lexer.OPERATOR, "=") {
		ret.addChild(mp.parseType())
	}

	if mp.tryMatch(lexer.OPERATOR, "=") {
		ret.addChild(mp.parseExpression())
		mp.match(lexer.PUNCTUATION, ";")
	} else if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret.addChild(mp.parseFunctionContent())
		mp.match(lexer.PUNCTUATION, "}")
	} else {
		mp.match(lexer.PUNCTUATION, ";")
	}

	return ret
}

func (mp *mainParser) parseStruct() ParseElement {
	ret := mp.startElement(lexer.STRUCT_DEF)

	ret.addChild(mp.match(lexer.KEYWORD, ""))
	ret.addChild(mp.parseTypeName())
	if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret.addChild(mp.parseStructContent())
		mp.match(lexer.PUNCTUATION, "}")
	} else {
		mp.match(lexer.PUNCTUATION, ";")
	}
	return ret
}

func (mp *mainParser) parseInterface() ParseElement {
	//FIXME implement
	return nil
}

func (mp *mainParser) parseMethod() ParseElement {
	//FIXME implement
	return nil
}

func (mp *mainParser) parseOperator() ParseElement {
	//FIXME implement
	return nil
}

func (mp *mainParser) parseTypeName() ParseElement {
	ret := mp.startElement(lexer.TYPE_NAME)
	ret.addChild(mp.match(lexer.SYMBOL, ""))
	if mp.tryMatch(lexer.PUNCTUATION, "(") {
		ret.addChild(mp.parseParameterList())
		mp.match(lexer.PUNCTUATION, ")")
	}
	return ret
}

func (mp *mainParser) parseFunctionType() ParseElement {

	ret:= mp.startElement(lexer.FUNCTION_TYPE)
	mp.match(lexer.PUNCTUATION, "(")
	ret.addChild(mp.parseParameterList())
	mp.match(lexer.PUNCTUATION, ")")
	if mp.tryMatch(lexer.OPERATOR, "->") {
		ret.addChild(mp.parseType())
	}
	return ret
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

func (mp *mainParser) parseType() ParseElement {

	ret := mp.startElement(lexer.TYPE)
	ret.addChild(mp.parseExpression())
	return ret
}

/* FIXME remove
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
*/

func (mp *mainParser) parseStructContent() ParseElement {
	ret := mp.startElement(lexer.LIST)
	for !mp.peek(0, lexer.PUNCTUATION, "}") {

		//FIXME this is similar to parseDef but not identical
		if mp.tryMatch(lexer.KEYWORD, "def") {
			el := mp.startElement(lexer.DEF)
			el.addChild(mp.match(lexer.SYMBOL, ""))
			if mp.peek(0, lexer.PUNCTUATION, "(") {
				el.addChild(mp.parseFunctionType())
			} else {
				el.addChild(mp.parseType())
			}
			if mp.tryMatch(lexer.PUNCTUATION, "{") {
				el.addChild(mp.parseFunctionContent())
				mp.match(lexer.PUNCTUATION, "}")
			} else {
				mp.match(lexer.PUNCTUATION, ";")
			}
		}

		if mp.tryMatch(lexer.KEYWORD, "implements") {
			el := mp.startElement(lexer.INTERFACE_MAP)
			el.addChild(mp.parseType())
			if mp.peek(0, lexer.PUNCTUATION, "{") {
				el.addChild(mp.parseImplementsContent())
			} else {
				mp.match(lexer.PUNCTUATION, ";")
			}
		}

		mp.checkProgress()
	}

	return ret
}

func (mp *mainParser) parseFunctionContent() ParseElement {

	ret := mp.startElement(lexer.FUNCTION_CONTENT)

	for !mp.peek(0, lexer.PUNCTUATION, "}") {
		ret.addChild(mp.parseFunctionStatement())
		mp.checkProgress()
	}

	return ret
}

func (mp *mainParser) parseFunctionStatement() ParseElement {

	if mp.tryMatch(lexer.PUNCTUATION, ";") {
		//empty statement
		return nil
	}

	if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret := mp.parseFunctionContent()
		mp.match(lexer.PUNCTUATION, "}")
		return ret
	}

	if mp.peek(0, lexer.KEYWORD, "def") {
		return mp.parseDef() //FIXME this allows nested functions!
	}

	if mp.tryMatch(lexer.KEYWORD, "if") {
		ret := mp.startElement(lexer.IF)
		ret.addChild(mp.parseExpression())
		mp.match(lexer.PUNCTUATION, "{")
		ret.addChild(mp.parseFunctionContent())
		mp.match(lexer.PUNCTUATION, "}")

		if mp.peek(0, lexer.KEYWORD, "else") {
			ret.addChild(mp.parseElseClause())
		}

		return ret
	}

	if mp.peek(0, lexer.KEYWORD, "while") {
		//FIXME implement
		return nil
	}

	if mp.peek(0, lexer.KEYWORD, "for") {
		//FIXME implement
		return nil
	}

	if mp.peek(0, lexer.KEYWORD, "with") {
		//FIXME implement
		return nil
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

	if mp.peek(0, lexer.KEYWORD, "continue") {
		//FIXME implement
		return nil
	}

	if mp.peek(0, lexer.KEYWORD, "break") {
		//FIXME implement
		return nil
	}

	if mp.peek(0, lexer.KEYWORD, "label") {
		//FIXME implement
		return nil
	}

	if mp.peek(0, lexer.KEYWORD, "goto") {
		//FIXME implement
		return nil
	}

	//FIXME more keyword statements

	return mp.parseAssignmentStatement()
}

func (mp *mainParser) parseElseClause() ParseElement {
	mp.match(lexer.KEYWORD, "else")
	if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret := mp.parseFunctionContent()
		mp.match(lexer.PUNCTUATION, "}")
		return ret
	}

	mp.match(lexer.KEYWORD, "if")
	ret := mp.startElement(lexer.IF)
	ret.addChild(mp.parseExpression())
	mp.match(lexer.PUNCTUATION, "{")
	ret.addChild(mp.parseFunctionContent())
	mp.match(lexer.PUNCTUATION, "}")

	if mp.peek(0, lexer.KEYWORD, "else") {
		ret.addChild(mp.parseElseClause())
	}

	return ret
}

//FIXME parseElseClause
//FIXME parseInterfaceContent

func (mp *mainParser) parseImplementsContent() ParseElement {
	//FIXME implement
	return nil
}

func (mp *mainParser) parseExpressionDot() ParseElement {

	ret := mp.startElement(lexer.DOT_LIST)
	ret.addChild(mp.match(lexer.SYMBOL, ""))

	//FIXME why is . an operator?
	for mp.tryMatch(lexer.OPERATOR, ".") {
		ret.addChild(mp.match(lexer.SYMBOL, ""))
	}
	return ret
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

	if mp.tryMatch(lexer.KEYWORD, "function") {
		return mp.parseFunctionType()
	} else if mp.peek(0, lexer.SYMBOL, "") {
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
	}	else if mp.tryMatch(lexer.PUNCTUATION, "[") {
		ret = mp.parseListContent()
		mp.match(lexer.PUNCTUATION, "]")
	}	else if mp.tryMatch(lexer.PUNCTUATION, "{") {
		ret = mp.parseMapContent()
		mp.match(lexer.PUNCTUATION, "}")
	} else {
		mp.error("missing expression")
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

func (mp *mainParser) parseListContent() ParseElement {
	ret := mp.startElement(lexer.LIST)
	for !mp.peek(0, lexer.PUNCTUATION, "]") {
		ret.addChild(mp.parseExpression());
		if !mp.tryMatch(lexer.PUNCTUATION, ",") {
			break
		}
		mp.checkProgress()
	}
	return ret
}

func (mp *mainParser) parseMapContent() ParseElement {
	ret := mp.startElement(lexer.LIST)
	for !mp.peek(0, lexer.PUNCTUATION, "}") {
		//FIXME organize assignment parsing
		el := mp.startElement(lexer.ASSIGNMENT)
		el.addChild(mp.parseExpression())
		mp.match(lexer.OPERATOR, "=")
		el.addChild(mp.parseExpression())
		mp.match(lexer.PUNCTUATION, ";")
		ret.addChild(el)
		mp.checkProgress()
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

