
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
	opElement := children[0]
	opName := opElement.TokenString()
	line, col := opElement.Position()

	if opElement.ElementType() != lexer.OPERATOR {
		output.Error(line, col, "FIXME not an operator: "+opName)
		return nil
	}

	args := make([]DataValue, len(children) -1)
	for i, x := range(children[1:]) {
		args[i] = EvaluateConstExpression(x, ctx)
	}

	op := ctx.Symbols.LookupOperator(opName)
	if op == nil {
		output.Error(line, col, "No definition for operator "+opName)
		//FIXME testing
		ctx.Symbols.Emit()
		return nil
	}

	return doConstOp(opElement, args, ctx)
}

func doConstOp(op parser.ParseElement,
	args []DataValue, ctx *EvalContext) DataValue {

	fmt.Printf("FIXME operator: %v args: %v\n", op, args)
	//FIXME
	return nil
}

