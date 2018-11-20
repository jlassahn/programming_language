
package symbols

import (
	"parser"
	"output"
)


type SymbolEval struct {}
func (*SymbolEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	line, col := el.Position()
	name := el.TokenString()
	symbol := ctx.Symbols.Lookup(name)

	if symbol == nil {
		output.Error(line, col, "undefined symbol: "+name)
		return nil
	}

	if !symbol.IsConst() {
		output.Error(line, col, "symbol must be const: "+name)
		return nil
	}
	return symbol.InitialValue()
}

