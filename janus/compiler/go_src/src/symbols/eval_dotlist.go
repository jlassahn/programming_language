
package symbols

import (
	"parser"
)


func evalDotList(el parser.ParseElement, ctx *EvalContext) DataValue {

	lhsEl := el.Children()[0]
	rhsEl := el.Children()[1]

	lhs := loopHandler(lhsEl, ctx)
	if lhs == nil {
		return nil
	}

	if lhs.Type() == NamespaceType {
		table := lhs.(NamespaceDataValue).AsSymbolTable()

		impCtx := EvalContext { }
		impCtx = *ctx
		impCtx.Symbols = table

		ret := loopHandler(rhsEl, &impCtx)

		return ret
	}

	//FIXME implement member access, etc
	parser.Error(el.FilePos(), "FIXME unimplemented dot operator %v, rhsEl")
	return nil

	/* FIXME remove
	vals := DotListAsStrings(el)
	if vals == nil {
		return nil
	}

	return evaluateConstDotList(vals, ctx)
	*/
}

func evaluateConstDotList(vals []string, ctx *EvalContext) DataValue {

	if len(vals) == 0 {
		return nil
	}

	symbol := ctx.Lookup(vals[0])
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

