
package symbols

import (
	"parser"
)


type TypeEval struct {}
func (*TypeEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	//FIXME is this always right?
	return EvaluateConstExpression(el.Children()[0], ctx)
}

