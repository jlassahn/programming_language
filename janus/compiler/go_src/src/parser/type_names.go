
package parser

type Tag struct { string }

func (t *Tag) String() string {
	return t.string
}

var EOF = &Tag{"EOF"}
var EMPTY = &Tag{"EMPTY"}
var COMMENT = &Tag{"COMMENT"}
var NUMBER = &Tag{"NUMBER"}
var STRING = &Tag{"STRING"}
var CHARACTER = &Tag{"CHARACTER"}
var SYMBOL = &Tag{"SYMBOL"}
var OPERATOR = &Tag{"OPERATOR"}
var PUNCTUATION = &Tag{"PUNCTUATION"}
var KEYWORD = &Tag{"KEYWORD"}
	// end of token types
var SOURCE_FILE = &Tag{"SOURCE_FILE"}
var HEADER = &Tag{"HEADER"}
var LIST = &Tag{"LIST"}
var IMPORT = &Tag{"IMPORT"}
var DOT_LIST = &Tag{"DOT_LIST"}
var DEF = &Tag{"DEF"}
var TYPE = &Tag{"TYPE"}
var FUNCTION_TYPE = &Tag{"FUNCTION_TYPE"}
var PARAMETER_LIST = &Tag{"PARAMETER_LIST"}
var PARAMETER = &Tag{"PARAMETER"}
var TYPE_LIST = &Tag{"TYPE_LIST"}
var FUNCTION_CONTENT = &Tag{"FUNCTION_CONTENT"}
var ASSIGNMENT = &Tag{"ASSIGNMENT"}
var EXPRESSION = &Tag{"EXPRESSION"}
var INDEX = &Tag{"INDEX"}
var CALL = &Tag{"CALL"}
var STRUCT_DEF = &Tag{"STRUCT_DEF"}
var TYPE_NAME = &Tag{"TYPE_NAME"}
var INTERFACE_MAP = &Tag{"INTERFACE_MAP"}
var IF = &Tag{"IF"}
var WHILE = &Tag{"WHILE"}
var FOR = &Tag{"FOR"}
var WITH = &Tag{"WITH"}
var INTERFACE = &Tag{"INTERFACE"}
var METHOD = &Tag{"METHOD"}
var OPERATOR_DEF = &Tag{"OPERATOR_DEF"}
var EXTENDS_DEF = &Tag{"EXTENDS_DEF"}
var IMPLEMENTS_DEF = &Tag{"IMPLEMENTS_DEF"} //FIXME INTERFACE_MAP
// FIXME do we want var SIZE_DEF = &Tag("SIZE_DEF"}
var ALIAS_DEF = &Tag{"ALIAS_DEF"}


//FIXME stuff that can be defined by the "operator" syntax should be an
//      Operator, other stuff should be Punctuation?
//FIXME complete list
//this must have longer operators first
var Operators = []string {
	"//=",
	"+/=",
	"-/=",
	"%%=",
	"+%=",
	"-%=",
	"^|=",
	">>=",
	"<<=",
	"->",
	"//",
	"+/",
	"-/",
	"%%",
	"+%",
	"-%",
	"^|",
	">>",
	"<<",
	"==",
	"!=",
	"<=",
	">=",
	"~~",
	"!~",
	"++",
	"--",
	"+=",
	"-=",
	"*=",
	"/=",
	"&=",
	"|=",
	"^",
	"*",
	"/",
	"+",
	"-",
	"&",
	"|",
	"!",
	"=",
	">",
	"<",
	":",
	"." }


var AssignmentOperators = map[string]bool {
	"//=": true,
	"+/=": true,
	"-/=": true,
	"%%=": true,
	"+%=": true,
	"-%=": true,
	"^|=": true,
	">>=": true,
	"<<=": true,
	"+=": true,
	"-=": true,
	"*=": true,
	"/=": true,
	"&=": true,
	"|=": true,
	"=": true }

var OrOperators = map[string]bool {
	"|": true,
	"^|": true }

var AndOperators = map[string]bool {
	"&": true }

var CompareOperators = map[string]bool {
	":": true,
	"<": true,
	">": true,
	"==": true,
	"!=": true,
	"<=": true,
	">=": true,
	"~~": true,
	"!~": true }

var AddOperators = map[string]bool {
	"+": true,
	"-": true }

var MultOperators = map[string]bool {
	">>": true,
	"<<": true,
	"//": true,
	"+/": true,
	"-/": true,
	"%%": true,
	"+%": true,
	"-%": true,
	"/": true,
	"*": true }

var ExpOperators = map[string]bool {
	"^": true }

var PrefixOperators = map[string]bool {
	"!": true }

var SuffixOperators = map[string]bool {
	"++": true,
	"--": true,
}


//FIXME complete list
var Keywords = map[string]bool {
	"import" : true,
	"def" : true,
	"const" : true,
	"struct" : true,
	"m_struct" : true,
	"interface" : true,
	"method" : true,
	"alias" : true,
	"operator" : true,
	"function" : true,
	"implements" : true,
	"extends" : true,
	// FIXME do we want "size" : true,
	"if" : true,
	"else" : true,
	"while" : true,
	"for" : true,
	"with" : true,
	"return" : true,
	"continue" : true,
	"break" : true,
	"label" : true,
	"goto" : true,
	"janus": true }

