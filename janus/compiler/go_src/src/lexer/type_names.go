
package lexer

type Tag struct { string }

func (t *Tag) String() string {
	return t.string
}

var ERROR = &Tag{"ERROR"}
var EOF = &Tag{"EOF"}
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
var INTERFACE = &Tag{"INTERFACE"}
var METHOD = &Tag{"METHOD"}
var OPERATOR_DEF = &Tag{"OPERATOR_DEF"}


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

var SuffixOperators = map[string]bool { }


//FIXME complete list
var Keywords = map[string]bool {
	"import" : true,
	"def" : true,
	"struct" : true,
	"m_struct" : true,
	"interface" : true,
	"method" : true,
	"operator" : true,
	"function" : true,
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

