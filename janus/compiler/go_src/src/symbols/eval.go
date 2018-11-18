
package symbols

import (
	"lexer"
	"parser"
	"fmt"
)


type Evaluator interface {
	EvaluateConstExpression(el parser.ParseElement, ctx *EvalContext) DataValue
}

type EvalContext struct {
	Symbols SymbolTable
}

var evaluators = map[lexer.Tag] Evaluator {
	*lexer.NUMBER: &NumberEval {},
	*lexer.DOT_LIST: &DotListEval {},
	*lexer.EXPRESSION: &ExpressionEval {},
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

