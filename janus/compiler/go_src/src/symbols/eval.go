
package symbols

import (
	"output"
	"parser"
)


type ConstEvaluator func(el parser.ParseElement, ctx *EvalContext) DataValue


type EvalContext struct {
	Symbols SymbolTable
	SymbolPreprocessor func(Symbol)Symbol
	CycleDetectSymbol *uninitializedSymbol
}

var handlers = map[*parser.Tag] ConstEvaluator {
	parser.NUMBER: evalNumber,
	parser.SYMBOL: evalSymbol,
	parser.DOT_LIST: evalDotList,
	parser.EXPRESSION: evalExpression,
	parser.FUNCTION_TYPE: evalFunctionType,
	parser.TYPE: evalType,
}


// indirect call to EvalConstExpression, to avoid a circular dependency
var loopHandler ConstEvaluator

func EvaluateConstExpression(
	el parser.ParseElement, ctx *EvalContext) DataValue {

	loopHandler = EvaluateConstExpression

	//eval := evaluators[*el.ElementType()]
	eval := handlers[el.ElementType()]
	if eval == nil {
		//FIXME implement
		output.FIXMEDebug("no evaluator for %v\n", el.ElementType())
		parser.EmitParseTree(el)
		ctx.Symbols.Emit()
		return nil
	} else {
		return eval(el, ctx)
	}
}

