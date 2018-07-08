
package lexer

const (
	ERROR = iota
	EOF
	COMMENT
	NUMBER
	STRING
	CHARACTER
	SYMBOL
	OPERATOR
	PUNCTUATION
	KEYWORD
	// end of token types
	SOURCE_FILE
	HEADER
	LIST
	IMPORT
	DOT_LIST
	DEF
	TYPE
	FUNCTION_TYPE
	PARAMETER_LIST
	PARAMETER
	TYPE_LIST
	FUNCTION_CONTENT
	ASSIGNMENT
	EXPRESSION
	INDEX
	CALL
	STRUCT_DEF
	TYPE_NAME
	INTERFACE_MAP
	IF
	INTERFACE
	METHOD
	OPERATOR_DEF
)


var TypeNames = map[int]string {
	ERROR: "ERROR",
	EOF: "EOF",
	COMMENT: "COMMENT",
	NUMBER: "NUMBER",
	STRING: "STRING",
	CHARACTER: "CHARACTER",
	SYMBOL: "SYMBOL",
	OPERATOR: "OPERATOR",
	PUNCTUATION: "PUNCTUATION",
	KEYWORD: "KEYWORD",
	SOURCE_FILE: "SOURCE_FILE",
	HEADER: "HEADER",
	LIST: "LIST",
	IMPORT: "IMPORT",
	DOT_LIST: "DOT_LIST",
	DEF: "DEF",
	TYPE: "TYPE",
	FUNCTION_TYPE: "FUNCTION_TYPE",
	PARAMETER_LIST: "PARAMETER_LIST",
	PARAMETER: "PARAMETER",
	TYPE_LIST: "TYPE_LIST",
	FUNCTION_CONTENT: "FUNCTION_CONTENT",
	ASSIGNMENT: "ASSIGNMENT",
	EXPRESSION: "EXPRESSION",
	INDEX: "INDEX",
	CALL: "CALL",
	STRUCT_DEF: "STRUCT_DEF",
	TYPE_NAME: "TYPE_NAME",
	INTERFACE_MAP: "INTERFACE_MAP",
	IF: "IF",
	INTERFACE: "INTERFACE",
	METHOD: "METHOD",
	OPERATOR_DEF: "OPERATOR_DEF" }

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

