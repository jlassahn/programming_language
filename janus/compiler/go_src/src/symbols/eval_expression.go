
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

	//FIXME operators have to be treated specially in symbol tables
	//      because they aren't qualified by namespaces like other
	//      symbols.  e.g. if a module defines +, you can't disambiguate it
	//      by saying 123 module.+ 456

	op := ctx.Symbols.Lookup(opName)
	if op == nil {
		output.Error(line, col, "No definition for operator "+opName)
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

