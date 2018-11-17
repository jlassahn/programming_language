
package symbols

import (
	"fmt"
	"parser"
	"lexer"
	"output"
)


type ExpressionEval struct {}
func (*ExpressionEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	children := el.Children()
	op := children[0]

	if op.ElementType() != lexer.OPERATOR {
		line, col := op.Position()
		output.Error(line, col, "FIXME not an operator: "+op.TokenString())
		return nil
	}

	//FIXME figure out preferred data type stuff here.

	args := make([]DataValue, len(children) -1)
	for i, x := range(children[1:]) {
		args[i] = EvaluateConstExpression(x, ctx)
	}

	return doConstOp(op, args, ctx)
}

func doConstOp(op parser.ParseElement,
	args []DataValue, ctx *EvalContext) DataValue {

	fmt.Printf("operator: %v args: %v\n", op, args)
	//FIXME
	return nil
}

