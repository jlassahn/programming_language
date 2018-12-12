
package symbols

import (
	"parser"
	"output"
	)

type GenerateNode struct {
	EvaluateConstExpression func(parser.ParseElement, *SymbolTable, DataType) DataValue
	GenerateLLVM func(parser.ParseElement, *SymbolTable, DataType, output.ObjectFile)
}

var Handlers = map[*parser.Tag]GenerateNode {
	parser.SOURCE_FILE: {nil, nil},
}


