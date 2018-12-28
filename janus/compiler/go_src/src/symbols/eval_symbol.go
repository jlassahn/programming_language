
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

	if !symbol.IsConst() {
		parser.Error(pos, "symbol must be const: %v", name)
		return nil
	}

	if symbol.Type() == FunctionChoiceType {
		for _,choice := range symbol.(FunctionChoiceSymbol).Choices() {
			if TypeMatches(choice.Type(), ctx.InitializerType) {
				if !choice.IsConst() {
					parser.Error(pos, "symbol must be const: %v", name)
					return nil
				}
				return choice.InitialValue()
			}
		}
		parser.Error(pos, "no matching function type for %v", name)
		return nil
	}

	return symbol.InitialValue()
}

