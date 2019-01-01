
package generator

import (
	"output"
	"parser"
	"symbols"
)


func genSymbol(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	sym := ctx.Lookup(el.TokenString())
	if sym == nil {
		parser.Error(el.FilePos(), "undefined symbol %v", el.TokenString())
		return nil
	}

	if sym.Type() == symbols.FunctionChoiceType {
		fn := sym.(symbols.FunctionChoiceSymbol)
		return NewFunctionChoiceResult(fn)
	}

	ret, ok := sym.GetGenVal().(Result)
	if ok {
		return ret
	}

	if sym.IsConst() {
		return NewDataVal(sym.InitialValue())
	}

	output.FIXMEDebug("NO VALUE FOUND")
	return nil
}

