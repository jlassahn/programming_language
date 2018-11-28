
package symbols

import (
	"parser"
	"fmt"
)


type Evaluator interface {
	EvaluateConstExpression(el parser.ParseElement, ctx *EvalContext) DataValue
}

type EvalContext struct {
	Symbols SymbolTable
}

var evaluators = map[parser.Tag] Evaluator {
	*parser.NUMBER: &NumberEval {},
	*parser.SYMBOL: &SymbolEval {},
	*parser.DOT_LIST: &DotListEval {},
	*parser.EXPRESSION: &ExpressionEval {},
}

func EvaluateConstExpression(
	el parser.ParseElement, ctx *EvalContext) DataValue {

	eval := evaluators[*el.ElementType()]
	if eval == nil {
		//FIXME implement
		fmt.Printf("no evaluator for %v\n", el.ElementType())
		parser.EmitParseTree(el)
		ctx.Symbols.Emit()
		return nil
	} else {
		return eval.EvaluateConstExpression(el, ctx)
	}
}

