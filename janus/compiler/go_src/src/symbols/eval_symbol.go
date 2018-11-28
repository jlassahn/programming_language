
package symbols

import (
	"parser"
	"output"
)


type SymbolEval struct {}
func (*SymbolEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	pos := el.FilePos()
	name := el.TokenString()
	symbol := ctx.Symbols.Lookup(name)

	if symbol == nil {
		output.Error(pos.Line, pos.Column, "undefined symbol: "+name)
		return nil
	}

	if !symbol.IsConst() {
		output.Error(pos.Line, pos.Column, "symbol must be const: "+name)
		return nil
	}
	return symbol.InitialValue()
}

