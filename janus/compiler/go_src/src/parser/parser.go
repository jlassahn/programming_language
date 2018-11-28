
package parser

import (
	"fmt"
)

type ParseElement interface {
	Children() []ParseElement
	Comments() []ParseElement
	ElementType() *Tag
	FilePos() *FilePosition
	Token() *Token
	TokenString() string
}

type Parser interface {
	GetElement() ParseElement
}

type parseElement struct {
	children []ParseElement
	comments []ParseElement
	elementType *Tag
	pos *FilePosition
	token *Token
}

func (pe *parseElement) String() string {
	return pe.elementType.String() + "(" + pe.TokenString() + ")"
}

func (pe *parseElement) Children() []ParseElement {
	return pe.children
}

func (pe *parseElement) Comments() []ParseElement {
	return pe.comments
}

func (pe *parseElement) ElementType() *Tag {
	return pe.elementType
}

func (pe *parseElement) Position() (int, int) {
	return pe.pos.Line, pe.pos.Column
}

func (pe *parseElement) FilePos() *FilePosition {
	return pe.pos
}

func (pe *parseElement) Token() *Token {
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
	lex *Lexer
}

func (tw *tokenWrapper) GetElement() ParseElement {
	tok := tw.lex.NextToken()

	ret := parseElement{
		children: nil,
		comments: nil,
		elementType: tok.TokenType,
		pos: &tok.Position,
		token: tok,
		}

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

		if el.ElementType() == EOF {
			if cm.depth > 0 {
				pos := cm.comments[0].FilePos()
				FatalError(pos, "EOF inside block comment")
			}
			if cm.comments != nil {
				comments := cm.comments
				cm.comments = nil
				pos := el.FilePos()
				return &parseElement {
					children: el.Children(),
					comments: comments,   //FIXME merge existing comments
					elementType: el.ElementType(),
					pos: pos,
					token: el.Token(),
				}
			} else {
				return el
			}
		}

		if cm.depth == 0 && el.ElementType() != COMMENT {
			if el.ElementType() == PUNCTUATION && el.TokenString() == "#{" {
				cm.comments = append(cm.comments, el)
				cm.depth = 1
			} else if cm.comments != nil {
				comments := cm.comments
				cm.comments = nil
				pos := el.FilePos()
				return &parseElement {
					children: el.Children(),
					comments: comments,   //FIXME merge existing comments
					elementType: el.ElementType(),
					pos: pos,
					token: el.Token(),
				}
			} else {
				return el
			}
		} else {
			cm.comments = append(cm.comments, el)
			if el.ElementType() == PUNCTUATION {
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

	if ret.ElementType() == EOF {
		FatalError(ret.FilePos(), "inernal error, unexpected EOF")
	}

	return ret
}

func (mp *mainParser) peek(pos int, etype *Tag, txt string) bool {
	el := mp.queue[pos]
	if etype != el.ElementType() {
		return false
	}

	if txt != "" && txt != el.TokenString() {
		return false;
	}

	return true
}

func (mp *mainParser) match(etype *Tag, txt string) ParseElement {
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

func (mp *mainParser) checkProgress() bool {
	if mp.progress {
		mp.progress = false
		return true
	} else if mp.peek(0, EOF, "") {
		Error(mp.queue[0].FilePos(), "unexpected EOF")
		return false
	} else {
		err := fmt.Sprintf("unexpected %s %s",
			mp.queue[0].ElementType(),
			mp.queue[0].TokenString())
		mp.error(err)
		mp.consume()
		return true
	}
}

func (mp *mainParser) tryMatch(etype *Tag, txt string) bool {
	if mp.peek(0, etype, txt) {
		mp.resync = false
		mp.consume()
		return true
	}
	return false
}

func (mp *mainParser) tryOperator(oplist map[string]bool) bool {
	if mp.queue[0].ElementType() != OPERATOR {
		return false
	}
	return oplist[mp.queue[0].TokenString()]
}


func (mp *mainParser) startElement(etype *Tag) *parseElement {

	pos := mp.queue[0].FilePos()

	comments := mp.queue[0].Comments()

	ret := parseElement {
		children: nil,
		comments: comments,
		elementType: etype,
		pos: pos,
		token: nil,
	}

	return &ret
}

func (mp *mainParser) error(txt string) {

	if !mp.resync {
		pos := mp.queue[0].FilePos()
		Error(pos, txt)
		mp.resync = true
	}
}

func (mp *mainParser) GetElement() ParseElement {

	return mp.parseFile()
}

func NewParser(lex *Lexer) Parser {

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

/**********
file:
	header
	| file  file_declaration
	;
******/

func (mp *mainParser) parseFile() ParseElement {

	ret := mp.startElement(SOURCE_FILE)

	ret.addChild(mp.parseHeader())

	for !mp.peek(0, EOF, "") {
		ret.addChild(mp.parseFileDeclaration())
		if !mp.checkProgress() { break }
	}

	return ret
}

/*********
header:
	JANUS NUMBER_TOKEN ';'
	| JANUS NUMBER_TOKEN '{' header_options '}'
	;
*****/

func (mp *mainParser) parseHeader() ParseElement {

	mp.match(KEYWORD, "janus")
	ret := mp.startElement(HEADER)
	ret.addChild(mp.match(NUMBER, ""))
	if mp.tryMatch(PUNCTUATION, "{") {
		ret.addChild(mp.parseHeaderOptions())
		mp.match(PUNCTUATION, "}")
	} else {
		mp.match(PUNCTUATION, ";")
	}

	return ret
}

/**********
header_options:
	// empty
	| header_options header_option
	;
******/

func (mp *mainParser) parseHeaderOptions() ParseElement {
	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseHeaderOption())
		if !mp.checkProgress() { break }
	}
	return ret
}

/******
header_option:
	expression_dot '=' expression ';'
	;
*****/

func (mp *mainParser) parseHeaderOption() ParseElement {
	ret := mp.startElement(ASSIGNMENT)
	ret.addChild(mp.parseExpressionDot())
	mp.match(OPERATOR, "=")
	ret.addChild(mp.parseExpression())
	mp.match(PUNCTUATION, ";")
	return ret
}

/******
file_declaration:
	';'
	| import_statement
	| struct_declaration
	| interface_declaration
	| method_declaration
	| alias_declaration
	| operator_declaration
	| def_statement
	;
*****/

func (mp *mainParser) parseFileDeclaration() ParseElement {

	if mp.tryMatch(PUNCTUATION, ";") {
		// empty statement
		return nil
	}

	if mp.peek(0, KEYWORD, "import") {
		return mp.parseImportStatement()
	}

	if mp.peek(0, KEYWORD, "struct") {
		return mp.parseStructDeclaration()
	}
	if mp.peek(0, KEYWORD, "m_struct") {
		return mp.parseStructDeclaration()
	}

	if mp.peek(0, KEYWORD, "interface") {
		return mp.parseInterfaceDeclaration()
	}

	if mp.peek(0, KEYWORD, "method") {
		return mp.parseMethodDeclaration()
	}

	if mp.peek(0, KEYWORD, "alias") {
		return mp.parseAliasDeclaration()
	}

	if mp.peek(0, KEYWORD, "operator") {
		return mp.parseOperatorDeclaration()
	}

	if mp.peek(0, KEYWORD, "def") {
		return mp.parseDefStatement()
	}
	if mp.peek(0, KEYWORD, "const") {
		return mp.parseDefStatement()
	}

	return nil
}

/*****
import_statement:
	IMPORT expression_dot ';'
	| IMPORT '.' '=' expression_dot ';'
	| IMPORT expression_dot '=' expression_dot ';'
	;
*****/

func (mp *mainParser) parseImportStatement() ParseElement {

	ret := mp.startElement(IMPORT)
	mp.match(KEYWORD, "import")

	if mp.peek(0, OPERATOR, ".") {
		ret.addChild(mp.match(OPERATOR, "."))
		mp.match(OPERATOR, "=")
		ret.addChild(mp.parseExpressionDot())
	} else {
		ret.addChild(mp.parseExpressionDot())
		if mp.tryMatch(OPERATOR, "=") {
			ret.addChild(mp.parseExpressionDot())
		}
	}
	mp.match(PUNCTUATION, ";")
	return ret
}

/*****
struct_declaration:
	STRUCT_OR_MSTRUCT type_name struct_options ';'
	| STRUCT_OR_MSTRUCT type_name struct_options '{' struct_content '}'
	;
*****/

func (mp *mainParser) parseStructDeclaration() ParseElement {
	ret := mp.startElement(STRUCT_DEF)

	ret.addChild(mp.match(KEYWORD, ""))
	ret.addChild(mp.parseTypeName())

	ret.addChild(mp.parseStructOptions())

	if mp.tryMatch(PUNCTUATION, "{") {
		ret.addChild(mp.parseStructContent())
		mp.match(PUNCTUATION, "}")
	} else {
		mp.match(PUNCTUATION, ";")
	}
	return ret
}

/*****
struct_options:
	// empty
	| struct_options_
	;

struct_options_:
	struct_option
	| struct_options_ ',' struct_option
	;
*****/

func (mp *mainParser) parseStructOptions() ParseElement {

	ret := mp.startElement(LIST)
	if mp.peek(0, PUNCTUATION, "{") {
		return ret
	}
	if mp.peek(0, PUNCTUATION, ";") {
		return ret
	}

	for {
		ret.addChild(mp.parseStructOption())

		if !mp.tryMatch(PUNCTUATION, ",") {
			return ret
		}

		if !mp.checkProgress() {
			return ret
		}
	}
}

/*******
struct_option:
	EXTENDS type
	| IMPLEMENTS type
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN
	FIXME do we want   | SIZE expression
	;
*****/

func (mp *mainParser) parseStructOption() ParseElement {

	if mp.tryMatch(KEYWORD, "extends") {
		ret := mp.startElement(EXTENDS_DEF)
		ret.addChild(mp.parseType())
		return ret
	}

	if mp.tryMatch(KEYWORD, "implements") {
		ret := mp.startElement(IMPLEMENTS_DEF)
		ret.addChild(mp.parseType())

		if mp.tryMatch(KEYWORD, "alias") {
			ret.addChild(mp.match(SYMBOL, ""))
		}
		return ret
	}

	/* FIXME do we want
	if mp.tryMatch(KEYWORD, "size") {
		ret := mp.startElement(SIZE_DEF)
		ret.addChild(mp.parseExpression())
		return ret
	}
	*/

	return nil
}

/******
struct_content:
	// empty
	| struct_content struct_element
	;
******/

func (mp *mainParser) parseStructContent() ParseElement {
	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseStructElement())
		if !mp.checkProgress() { break }
	}
	return ret
}

/******
struct_element:
	DEF SYMBOL_TOKEN type ';'
	| DEF SYMBOL_TOKEN function_type '{' function_content '}'
	| DEF SYMBOL_TOKEN type '{' function_content '}'
	| extends_declaration
	| implements_declaration
	;

******/

func (mp *mainParser) parseStructElement() ParseElement {

	if mp.tryMatch(KEYWORD, "def") {
		el := mp.startElement(DEF)
		el.addChild(mp.match(SYMBOL, ""))
		if mp.peek(0, PUNCTUATION, "(") {
			el.addChild(mp.parseFunctionType())
		} else {
			el.addChild(mp.parseType())
		}
		if mp.tryMatch(PUNCTUATION, "{") {
			el.addChild(mp.parseFunctionContent())
			mp.match(PUNCTUATION, "}")
		} else {
			mp.match(PUNCTUATION, ";")
		}

		return el
	}

	if mp.peek(0, KEYWORD, "extends") {
		return mp.parseExtendsDeclaration()
	}

	if mp.peek(0, KEYWORD, "implements") {
		return mp.parseImplementsDeclaration()
	}

	return nil
}


/******
extends_declaration:
	EXTENDS type ';'
	| EXTENDS type '{' extends_content '}'
	;
*****/

func (mp *mainParser) parseExtendsDeclaration() ParseElement {

	mp.match(KEYWORD, "extends")
	el := mp.startElement(EXTENDS_DEF)

	el.addChild(mp.parseType())
	if mp.peek(0, PUNCTUATION, "{") {
		el.addChild(mp.parseExtendsContent())
		mp.match(PUNCTUATION, "}")
	} else {
		mp.match(PUNCTUATION, ";")
	}

	return el
}

/*****
extends_content:
	// empty
	| extends_content extends_item
	;
*****/

func (mp *mainParser) parseExtendsContent() ParseElement {

	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseExtendsItem())
		if !mp.checkProgress() { break }
	}
	return ret
}

/*****
extends_item:
	SYMBOL_TOKEN '=' SYMBOL_TOKEN ';'
	;
******/

func (mp *mainParser) parseExtendsItem() ParseElement {

	ret := mp.startElement(EXTENDS_DEF)
	ret.addChild(mp.match(SYMBOL, ""))
	mp.match(OPERATOR, "=")
	ret.addChild(mp.match(SYMBOL, ""))
	mp.match(PUNCTUATION, ";")

	return ret
}


/******
implements_declaration:
	IMPLEMENTS type ';'
	| IMPLEMENTS type '{' implements_content '}'
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN ';'
	| IMPLEMENTS type ALIAS SYMBOL_TOKEN '{' implements_content '}'
	;

*****/

func (mp *mainParser) parseImplementsDeclaration() ParseElement {

	mp.match(KEYWORD, "implements")
	el := mp.startElement(INTERFACE_MAP)

	el.addChild(mp.parseType())

	if mp.tryMatch(KEYWORD, "alias") {
		el.addChild(mp.match(SYMBOL, ""))
	}

	if mp.peek(0, PUNCTUATION, "{") {
		el.addChild(mp.parseImplementsContent())
		mp.match(PUNCTUATION, "}")
	} else {
		mp.match(PUNCTUATION, ";")
	}

	return el
}

/*****
implements_content:
	// empty
	| implements_content implements_item
	;
******/

func (mp *mainParser) parseImplementsContent() ParseElement {

	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseImplementsItem())
		if !mp.checkProgress() { break }
	}
	return ret
}

/*****
implements_item:
	SYMBOL_TOKEN '=' SYMBOL_TOKEN ';'
	;
*****/

func (mp *mainParser) parseImplementsItem() ParseElement {

	ret := mp.startElement(IMPLEMENTS_DEF)
	ret.addChild(mp.match(SYMBOL, ""))
	mp.match(OPERATOR, "=")
	ret.addChild(mp.match(SYMBOL, ""))
	mp.match(PUNCTUATION, ";")

	return ret
}

/*****
type_name:
	SYMBOL_TOKEN
	| SYMBOL_TOKEN '(' parameter_list ')'
	;
*****/

func (mp *mainParser) parseTypeName() ParseElement {
	ret := mp.startElement(TYPE_NAME)
	ret.addChild(mp.match(SYMBOL, ""))
	if mp.tryMatch(PUNCTUATION, "(") {
		ret.addChild(mp.parseParameterList())
		mp.match(PUNCTUATION, ")")
	}
	return ret
}

/*****
interface_declaration:
	INTERFACE type_name interface_options '{' interface_content '}'
	;
*****/

func (mp *mainParser) parseInterfaceDeclaration() ParseElement {

	ret := mp.startElement(INTERFACE)
	
	mp.match(KEYWORD, "interface")
	ret.addChild(mp.parseTypeName())
	ret.addChild(mp.parseInterfaceOptions())
	mp.match(PUNCTUATION, "{")
	ret.addChild(mp.parseInterfaceContent())
	mp.match(PUNCTUATION, "}")

	return ret
}

/*****
interface_options:
	// empty
	| interface_options_
	;

interface_options_:
	interface_option
	| interface_options_ ',' interface_option
	;
*****/

func (mp *mainParser) parseInterfaceOptions() ParseElement {

	ret := mp.startElement(LIST)
	if mp.peek(0, PUNCTUATION, "{") {
		return ret
	}

	for {
		ret.addChild(mp.parseInterfaceOption())

		if !mp.tryMatch(PUNCTUATION, ",") {
			return ret
		}

		if !mp.checkProgress() {
			return ret
		}
	}
}

/*****
interface_option:
	EXTENDS type
	;
*****/

func (mp *mainParser) parseInterfaceOption() ParseElement {

	if mp.tryMatch(KEYWORD, "extends") {
		ret := mp.startElement(EXTENDS_DEF)
		ret.addChild(mp.parseType())
		return ret
	}

	return nil
}

/******
interface_content:
	// empty
	| interface_content interface_element
	;
******/

func (mp *mainParser) parseInterfaceContent() ParseElement {
	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseInterfaceElement())
		if !mp.checkProgress() { break }
	}
	return ret
}

/*****
interface_element:
	DEF SYMBOL_TOKEN type ';'
	| DEF SYMBOL_TOKEN function_type ';'
	| extends_declaration
	;
*****/

func (mp *mainParser) parseInterfaceElement() ParseElement {

	if mp.tryMatch(KEYWORD, "def") {
		el := mp.startElement(DEF)
		el.addChild(mp.match(SYMBOL, ""))
		if mp.peek(0, PUNCTUATION, "(") {
			el.addChild(mp.parseFunctionType())
		} else {
			el.addChild(mp.parseType())
		}
		mp.match(PUNCTUATION, ";")
		return el
	}

	if mp.peek(0, KEYWORD, "extends") {
		return mp.parseExtendsDeclaration()
	}

	return nil
}

/*****
method_declaration:
	METHOD type SYMBOL_TOKEN function_type '{' function_content '}'
	| METHOD type SYMBOL_TOKEN function_type '=' expression ';'
	| METHOD type SYMBOL_TOKEN function_type ';'
	| METHOD type SYMBOL_TOKEN type '{' function_content '}'
	| METHOD type SYMBOL_TOKEN type '=' expression ';'
	| METHOD type SYMBOL_TOKEN type ';'
	;
*****/

func (mp *mainParser) parseMethodDeclaration() ParseElement {

	ret := mp.startElement(METHOD)

	mp.match(KEYWORD, "method")

	ret.addChild(mp.parseType())
	ret.addChild(mp.match(SYMBOL, ""))

	if mp.peek(0, PUNCTUATION, "(") {
		ret.addChild(mp.parseFunctionType())
	} else {
		ret.addChild(mp.parseType())
	}

	if mp.tryMatch(PUNCTUATION, "{") {
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
	} else if mp.tryMatch(OPERATOR, "=") {
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, ";")
	} else {
		mp.match(PUNCTUATION, ";")
	}

	return ret
}

/*****
alias_declaration:
	ALIAS SYMBOL_TOKEN type ';'
	;
*****/

func (mp *mainParser) parseAliasDeclaration() ParseElement {

	ret := mp.startElement(ALIAS_DEF)

	mp.match(KEYWORD, "alias")
	ret.addChild(mp.match(SYMBOL, ""))
	ret.addChild(mp.parseType())
	mp.match(PUNCTUATION, ";")

	return ret
}

/*****
operator_declaration:
	OPERATOR ANY_OP function_type '{' function_content '}'
	| OPERATOR ANY_OP function_type '=' expression ';'
	| OPERATOR ANY_OP function_type ';'
	| OPERATOR ANY_OP type '{' function_content '}'
	| OPERATOR ANY_OP type '=' expression ';'
	| OPERATOR ANY_OP type ';'
	;
*****/

func (mp *mainParser) parseOperatorDeclaration() ParseElement {

	//FIXME can any operator be defined or is it a subset?
	//      definitely as set up now "." and "=" should be forbidden

	ret := mp.startElement(OPERATOR_DEF)
	ret.addChild(mp.match(OPERATOR, ""))

	if mp.peek(0, PUNCTUATION, "(") {
		ret.addChild(mp.parseFunctionType())
	} else {
		ret.addChild(mp.parseType())
	}

	if mp.tryMatch(PUNCTUATION, "{") {
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
	} else if mp.tryMatch(OPERATOR, "=") {
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, ";")
	} else {
		mp.match(PUNCTUATION, ";")
	}

	return ret
}

/*****
def_statement:
	def_or_const SYMBOL_TOKEN function_type initializer
	| def_or_const SYMBOL_TOKEN function_type '{' function_content '}'
	| def_or_const SYMBOL_TOKEN function_type ';'
	| def_or_const SYMBOL_TOKEN type initializer
	| def_or_const SYMBOL_TOKEN type '{' function_content '}'
	| def_or_const SYMBOL_TOKEN type ';'
	| def_or_const SYMBOL_TOKEN initializer
	;

def_or_const:
	DEF
	| CONST
	;
*****/

func (mp *mainParser) parseDefStatement() ParseElement {

	ret := mp.startElement(DEF)

	//def or const
	ret.addChild(mp.match(KEYWORD, ""))

	ret.addChild(mp.match(SYMBOL, ""))

	if mp.peek(0, PUNCTUATION, "(") {
		ret.addChild(mp.parseFunctionType())
	} else if !mp.peek(0, OPERATOR, "=") {
		ret.addChild(mp.parseType())
	} else {
		ret.addChild(mp.startElement(EMPTY))
	}

	if mp.peek(0, OPERATOR, "=") {
		ret.addChild(mp.parseInitializer())
	} else if mp.tryMatch(PUNCTUATION, "{") {
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
	} else {
		ret.addChild(mp.startElement(EMPTY))
		mp.match(PUNCTUATION, ";")
	}

	return ret
}

/*****
initializer:
	'=' expression ';'
	| '=' '{' map_content '}'
	| '=' '[' list_content ']'
	;
*****/

func (mp *mainParser) parseInitializer() ParseElement {

	var ret ParseElement = nil
	mp.match(OPERATOR, "=")

	if mp.tryMatch(PUNCTUATION, "{") {
		ret = mp.parseMapContent()
		mp.match(PUNCTUATION, "}")
	} else if mp.tryMatch(PUNCTUATION, "[") {
		ret = mp.parseListContent()
		mp.match(PUNCTUATION, "]")
	} else {
		ret = mp.parseExpression()
		mp.match(PUNCTUATION, ";")
	}

	return ret
}

/*****
function_content:
	// empty
	| function_content function_statement
	;
*****/

func (mp *mainParser) parseFunctionContent() ParseElement {

	ret := mp.startElement(FUNCTION_CONTENT)

	for !mp.peek(0, PUNCTUATION, "}") {
		ret.addChild(mp.parseFunctionStatement())
		if !mp.checkProgress() { break }
	}

	return ret
}

/*****
function_statement:
	';'
	| '{' function_content '}'
	| def_statement
	| IF expression '{' function_content '}'
	| IF expression '{' function_content '}' else_statement
	| WHILE expression '{' function_content '}'
	| FOR SYMBOL_TOKEN '=' expression '{' function_content '}'
	| WITH SYMBOL_TOKEN '=' expression '{' function_content '}'
	| RETURN expression ';'
	| RETURN ';'
	| CONTINUE ';'
	| BREAK NUMBER_TOKEN ';'
	| BREAK ';'
	| LABEL SYMBOL_TOKEN ';'
	| GOTO SYMBOL_TOKEN ';'
	| assignment_statement
	;
*****/

func (mp *mainParser) parseFunctionStatement() ParseElement {

	if mp.tryMatch(PUNCTUATION, ";") {
		//empty statement
		return nil
	}

	if mp.tryMatch(PUNCTUATION, "{") {
		ret := mp.parseFunctionContent()
		mp.match(PUNCTUATION, "}")
		return ret
	}

	if mp.peek(0, KEYWORD, "def") {
		return mp.parseDefStatement()
	}
	if mp.peek(0, KEYWORD, "const") {
		return mp.parseDefStatement()
	}

	if mp.tryMatch(KEYWORD, "if") {
		ret := mp.startElement(IF)
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, "{")
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")

		if mp.peek(0, KEYWORD, "else") {
			ret.addChild(mp.parseElseStatement())
		}

		return ret
	}

	if mp.tryMatch(KEYWORD, "while") {
		ret := mp.startElement(WHILE)
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, "{")
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
		return ret
	}

	if mp.tryMatch(KEYWORD, "for") {
		ret := mp.startElement(FOR)
		ret.addChild(mp.match(SYMBOL, ""))
		mp.match(OPERATOR, "=")
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, "{")
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
		return ret
	}

	if mp.peek(0, KEYWORD, "with") {
		ret := mp.startElement(WITH)
		ret.addChild(mp.match(SYMBOL, ""))
		mp.match(OPERATOR, "=")
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, "{")
		ret.addChild(mp.parseFunctionContent())
		mp.match(PUNCTUATION, "}")
		return ret
	}

	if mp.peek(0, KEYWORD, "return") {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		if !mp.peek(0, PUNCTUATION, ";") {
			ret.addChild(mp.parseExpression())
		}
		mp.match(PUNCTUATION, ";")
		return ret
	}

	if mp.peek(0, KEYWORD, "continue") {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		mp.match(PUNCTUATION, ";")
		return ret
	}

	if mp.peek(0, KEYWORD, "break") {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		if !mp.peek(0, PUNCTUATION, ";") {
			ret.addChild(mp.match(NUMBER, ""))
		}
		mp.match(PUNCTUATION, ";")
		return ret
	}

	if mp.peek(0, KEYWORD, "label") {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		ret.addChild(mp.match(SYMBOL, ""))
		mp.match(PUNCTUATION, ";")
		return ret
	}

	if mp.peek(0, KEYWORD, "goto") {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		ret.addChild(mp.match(SYMBOL, ""))
		mp.match(PUNCTUATION, ";")
		return ret
	}

	return mp.parseAssignmentStatement()
}

/*****
else_statement:
	ELSE '{' function_content '}'
	| ELSE IF expression '{' function_content '}'
	| ELSE IF expression '{' function_content '}' else_statement
	;
*****/

func (mp *mainParser) parseElseStatement() ParseElement {
	mp.match(KEYWORD, "else")
	if mp.tryMatch(PUNCTUATION, "{") {
		ret := mp.parseFunctionContent()
		mp.match(PUNCTUATION, "}")
		return ret
	}

	mp.match(KEYWORD, "if")
	ret := mp.startElement(IF)
	ret.addChild(mp.parseExpression())
	mp.match(PUNCTUATION, "{")
	ret.addChild(mp.parseFunctionContent())
	mp.match(PUNCTUATION, "}")

	if mp.peek(0, KEYWORD, "else") {
		ret.addChild(mp.parseElseStatement())
	}

	return ret
}

/****
assignment_statement:
	expression ';'
	| expression initializer
	| expression ASSIGNMENT_OP expression ';'
	;
*****/

func (mp *mainParser) parseAssignmentStatement() ParseElement {

	lhs := mp.parseExpression()
	if mp.peek(0, OPERATOR, "=") {
		ret := mp.startElement(ASSIGNMENT)
		ret.addChild(lhs)
		ret.addChild(mp.parseInitializer())
		return ret
	}

	if mp.tryOperator(AssignmentOperators) {
		ret := mp.startElement(ASSIGNMENT)
		ret.addChild(mp.consume())
		ret.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, ";")
		return ret
	}

	mp.match(PUNCTUATION, ";")
	return lhs
}

/*****
expression:
	expression_and
	| expression OR_OP expression_and   //  "|"  "^|"
	;
*****/

func (mp *mainParser) parseExpression() ParseElement {
	ret := mp.parseExpressionAnd()
	for mp.tryOperator(OrOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionAnd())
		ret = el
	}
	return ret
}

/*****
expression_and:
	expression_compare
	| expression_and AND_OP expression_compare  //  "&"
	;
*****/

func (mp *mainParser) parseExpressionAnd() ParseElement {
	ret := mp.parseExpressionCompare()
	for mp.tryOperator(AndOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionCompare())
		ret = el
	}
	return ret
}

/*****
expression_compare:
	expression_add
	| expression_compare COMPARE_OP expression_add //"==" "!=" "~~" "!~" "<=" ">=" ">" "<" ":"
	;
*****/

func (mp *mainParser) parseExpressionCompare() ParseElement {
	ret := mp.parseExpressionAdd()
	for mp.tryOperator(CompareOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionAdd())
		ret = el
	}
	return ret
}

/*****
expression_add:
	expression_mult
	| expression_add ADD_OP expression_mult  //  "+" "-"
	;
******/

func (mp *mainParser) parseExpressionAdd() ParseElement {
	ret := mp.parseExpressionMult()
	for mp.tryOperator(AddOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionMult())
		ret = el
	}
	return ret
}

/*****
expression_mult:
	expression_exp
	| expression_mult MULT_OP expression_exp // "*" "/"  "//" "+/" "-/" "%%" "+%" "-%" "<<" ">>"
	;
******/

func (mp *mainParser) parseExpressionMult() ParseElement {
	ret := mp.parseExpressionExp()
	for mp.tryOperator(MultOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionExp())
		ret = el
	}
	return ret
}

/*****
expression_exp:
	expression_prefix
	| expression_exp EXP_OP expression_prefix  //  "^"
*****/

func (mp *mainParser) parseExpressionExp() ParseElement {
	ret := mp.parseExpressionPrefix()
	for mp.tryOperator(ExpOperators) {
		el := mp.startElement(EXPRESSION)
		el.addChild(mp.consume())
		el.addChild(ret)
		el.addChild(mp.parseExpressionPrefix())
		ret = el
	}
	return ret
}

/*****
expression_prefix:
	expression_suffix
	| PREFIX_OP expression_prefix   // "!"
*****/

func (mp *mainParser) parseExpressionPrefix() ParseElement {

	if mp.tryOperator(PrefixOperators) {
		ret := mp.startElement(EXPRESSION)
		ret.addChild(mp.consume())
		ret.addChild(mp.parseExpressionPrefix())
		return ret
	} else {
		return mp.parseExpressionSuffix()
	}
}

/*****
expression_suffix:
	NUMBER_TOKEN
	| STRING_TOKEN
	| CHARACTER_TOKEN
	| SYMBOL_TOKEN
	| FUNCTION function_type
	| expression_suffix SUFFIX_OP
	| expression_suffix '[' expression ']'
	| expression_suffix '(' list_content ')'
	| expression_suffix '.' SYMBOL_TOKEN
	| '(' expression ')'
	;
*****/

func (mp *mainParser) parseExpressionSuffix() ParseElement {

	var ret ParseElement

	if mp.peek(0, NUMBER, "") {
		ret = mp.consume()
	} else if mp.peek(0, STRING, "") {
		ret = mp.consume()
	} else if mp.peek(0, CHARACTER, "") {
		ret = mp.consume()
	} else if mp.peek(0, SYMBOL, "") {
		ret = mp.consume()
	} else if mp.tryMatch(KEYWORD, "function") {
		ret = mp.parseFunctionType()
	} else if mp.tryMatch(PUNCTUATION, "(") {
		ret = mp.parseExpression()
		mp.match(PUNCTUATION, ")")
	} else {
		mp.error("missing expression")
	}


	for {
		if mp.tryOperator(SuffixOperators) {
			el := mp.startElement(EXPRESSION)
			el.addChild(mp.consume())
			el.addChild(ret)
			ret = el
		} else if mp.tryMatch(PUNCTUATION, "[") {
			el := mp.startElement(INDEX)
			el.addChild(ret)
			el.addChild(mp.parseExpression())
			mp.match(PUNCTUATION, "]")
			ret = el
		} else if mp.tryMatch(PUNCTUATION, "(") {
			el := mp.startElement(CALL)
			el.addChild(ret)
			el.addChild(mp.parseListContent())
			mp.match(PUNCTUATION, ")")
			ret = el
		} else if mp.tryMatch(OPERATOR, ".") {
			el := mp.startElement(DOT_LIST)
			el.addChild(ret)
			el.addChild(mp.match(SYMBOL, ""))
			ret = el
		} else {
			break
		}
	}

	return ret
}

/*****
expression_dot:
	SYMBOL_TOKEN
	| expression_dot '.' SYMBOL_TOKEN
	;
*****/

func (mp *mainParser) parseExpressionDot() ParseElement {

	ret := mp.match(SYMBOL, "")

	for mp.tryMatch(OPERATOR, ".") {
		el := mp.startElement(DOT_LIST)
		el.addChild(ret)
		el.addChild(mp.match(SYMBOL, ""))
		ret = el
	}
	return ret
}

/*****
type:
	expression_dot
	expression_dot '(' list_content ')'
	FUNCTION function_type
	;
*****/

func (mp *mainParser) parseType() ParseElement {

	if mp.tryMatch(KEYWORD, "function") {
		return mp.parseFunctionType()
	} else {
		ret := mp.startElement(TYPE)
		ret.addChild(mp.parseExpressionDot())
		if mp.tryMatch(PUNCTUATION, "(") {
			ret.addChild(mp.parseListContent())
			mp.match(PUNCTUATION, ")")
		}
		return ret
	}
}

/****
function_type:
	'(' parameter_list ')'
	| '(' parameter_list ')' "->" type
	;
****/

func (mp *mainParser) parseFunctionType() ParseElement {

	ret:= mp.startElement(FUNCTION_TYPE)
	mp.match(PUNCTUATION, "(")
	ret.addChild(mp.parseParameterList())
	mp.match(PUNCTUATION, ")")
	if mp.tryMatch(OPERATOR, "->") {
		ret.addChild(mp.parseType())
	}
	return ret
}

/*****
list_content:
	// empty
	| list_content_
	;

list_content_:
	expression
	| list_content_ ',' expression
	;
*****/

func (mp *mainParser) parseListContent() ParseElement {
	//FIXME supports optional final comma, not in grammar
	ret := mp.startElement(LIST)
	for {
		if mp.peek(0, PUNCTUATION, "]") {
			break
		}
		if mp.peek(0, PUNCTUATION, ")") {
			break
		}

		ret.addChild(mp.parseExpression());
		if !mp.tryMatch(PUNCTUATION, ",") {
			break
		}
		if !mp.checkProgress() { break }
	}
	return ret
}

/*****
map_content:
	// empty
	| expression '=' expression ';' map_content
*****/

func (mp *mainParser) parseMapContent() ParseElement {
	ret := mp.startElement(LIST)
	for !mp.peek(0, PUNCTUATION, "}") {
		el := mp.startElement(ASSIGNMENT)
		el.addChild(mp.parseExpression())
		mp.match(OPERATOR, "=")
		el.addChild(mp.parseExpression())
		mp.match(PUNCTUATION, ";")
		ret.addChild(el)
		if !mp.checkProgress() { break }
	}
	return ret
}

/*****
parameter_list:
	// empty
	| parameter_list_
	;

parameter_list_:
	SYMBOL_TOKEN type
	| SYMBOL_TOKEN '>' type
	| parameter_list_ ',' SYMBOL_TOKEN type
	| parameter_list_ ',' SYMBOL_TOKEN '>' type
	;

****/


func (mp *mainParser) parseParameterList() ParseElement {

	ret := mp.startElement(PARAMETER_LIST)
	for mp.peek(0, SYMBOL, "") {
		el := mp.startElement(PARAMETER)
		ret.addChild(el)
		el.addChild(mp.match(SYMBOL, ""))

		if mp.tryMatch(OPERATOR, ">") {
			//FIXME how to represent in the syntax tree?
		}

		el.addChild(mp.parseType())

		if !mp.tryMatch(PUNCTUATION, ",") {
			break
		}
	}
	return ret
}


