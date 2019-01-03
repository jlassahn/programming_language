
package generator

import (
	"output"
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

	if rhsEl.ElementType() == parser.SYMBOL {
		name := rhsEl.TokenString()
		sym := lhs.Type().Members()[name]
		if sym != nil {
			output.FIXMEDebug("FIXME dot found member %v", sym)

			if sym.Type() == symbols.FunctionChoiceType {
				fn := sym.(symbols.FunctionChoiceSymbol)
				return NewMethodChoiceResult(fn, lhs)
			}

			output.FatalError("unimplemented dot operator %v", sym)
		}
		parser.Error(rhsEl.FilePos(), "unknown member or method %v", name)
		return nil
	}

	//FIXME implement member access, etc
	parser.Error(el.FilePos(), "FIXME unimplemented dot operator %v", rhsEl)
	return nil
}

