
package generator

import (
	"parser"
	"symbols"
)


func genDotList(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	lhsEl := el.Children()[0]
	rhsEl := el.Children()[1]

	lhs := loopHandler(genFunc, ctx, lhsEl)
	if lhs == nil {
		return nil
	}

	if lhs.Type() == symbols.NamespaceType {
		table := lhs.ConstVal().(symbols.NamespaceDataValue).AsSymbolTable()

		impCtx := symbols.EvalContext { }
		impCtx.Symbols = table

		ret := loopHandler(genFunc,&impCtx, rhsEl)

		return ret
	}

	//FIXME implement member access, etc
	parser.Error(el.FilePos(), "FIXME unimplemented dot operator %v, rhsEl")
	return nil
}

