
package symbols

import (
	"parser"
)


func evalDotList(el parser.ParseElement, ctx *EvalContext) DataValue {

	vals := DotListAsStrings(el)
	if vals == nil {
		return nil
	}

	return evaluateConstDotList(vals, ctx)
}

func evaluateConstDotList(vals []string, ctx *EvalContext) DataValue {

	if len(vals) == 0 {
		return nil
	}

	symbol := ctx.Symbols.Lookup(vals[0])
	if symbol == nil {
		return nil
	}

	//FIXME fake
	return symbol.InitialValue()
}

func DotListAsStrings(el parser.ParseElement) []string {

	if el.ElementType() == parser.SYMBOL {
		return []string{ el.TokenString() }
	}

	if el.ElementType() != parser.DOT_LIST {
		return nil
	}

	var ret []string
	for _, x := range(el.Children()) {
		ret = append(ret, DotListAsStrings(x)...)
	}
	return ret
}

