
package symbols

import (
	"lexer"
	"parser"
	"output"
	)

type GenerateNode struct {
	EvaluateConstExpression func(parser.ParseElement, *SymbolTable, *DataType) *DataValue
	GenerateLLVM func(parser.ParseElement, *SymbolTable, *DataType, output.LLVMFile)
}

var Handlers = map[*lexer.Tag]GenerateNode {
	lexer.SOURCE_FILE: {nil, nil} }


