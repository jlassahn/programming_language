
package symbols

import (
	"lexer"
	"parser"
)


type DotListEval struct {}
func (*DotListEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

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

	if el.ElementType() != lexer.DOT_LIST {
		return nil
	}

	var ret []string
	for _, x := range(el.Children()) {
		ret = append(ret, x.TokenString())
	}
	return ret
}

