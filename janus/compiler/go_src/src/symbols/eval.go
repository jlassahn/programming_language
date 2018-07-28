
package symbols

import (
	"parser"
)

func DotListAsStrings(el parser.ParseElement) []string {
	var ret []string
	for _, x := range(el.Children()) {
		ret = append(ret, x.TokenString())
	}
	return ret
}

func EvaluateConstExpression(
	el parser.ParseElement, ctx *SymbolTable) *DataValue {

	parser.EmitParseTree(el)
	//FIXME implement
	return nil
}

