
package symbols

import (
	"parser"
)



func evalSymbol(el parser.ParseElement, ctx *EvalContext) DataValue {

	pos := el.FilePos()
	name := el.TokenString()
	symbol := ctx.Lookup(name)

	if symbol == nil {
		parser.Error(pos, "undefined symbol: %v", name)
		return nil
	}

	if symbol.IsConst() {
		return symbol.InitialValue()
	}

	parser.Error(pos, "symbol must be const: %v", name)
	return nil
}

