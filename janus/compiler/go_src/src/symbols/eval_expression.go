
package symbols

import (
	"output"
	"parser"
)


type ExpressionEval struct {}
func (*ExpressionEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	children := el.Children()
	opElement := children[0]
	opName := opElement.TokenString()
	pos := opElement.FilePos()

	if opElement.ElementType() != parser.OPERATOR {
		parser.Error(pos, "FIXME not an operator: %v", opName)
		return nil
	}

	args := make([]DataValue, len(children) -1)
	for i, x := range(children[1:]) {
		args[i] = EvaluateConstExpression(x, ctx)
		if args[i] == nil {
			output.FIXMEDebug("FIXME args not available")
			return nil
		}
	}

	opChoices := ctx.Symbols.LookupOperator(opName)
	if opChoices == nil {
		parser.Error(pos, "No definition for operator %v", opName)
		//FIXME testing
		ctx.Symbols.Emit()
		return nil
	}

	op := selectFunctionChoice(opChoices, args)
	if op == nil {
		parser.Error(pos, "Operator %v can't take these parameters", opName)
		return nil
	}

	if !op.IsConst() {
		parser.Error(pos, "Operator %v not const", opName)
		return nil
	}

	return doConstOp(op, args, ctx)
}

func doConstOp(op Symbol,
	args []DataValue, ctx *EvalContext) DataValue {

	val := op.InitialValue()

	return val.(CodeDataValue).EvaluateConst(op, args)
}

func selectFunctionChoice(op FunctionChoiceSymbol, args []DataValue) Symbol {

	for _,choice := range op.Choices() {

		dtype := choice.Type().(FunctionDataType)
		params := dtype.Parameters()

		if !CanConvertArgs(args, params) { continue; }

		return choice
	}

	return nil
}


func CanConvertArgs(args []DataValue, params []FunctionParameter) bool {

	if len(args) != len(params) {
		return false
	}

	for i, param := range params {
		arg := args[i]

		if param.AutoConvert {
			if !CanConvert(arg.Type(), param.DType) {
				return false
			}
		} else {
			if !TypeMatches(arg.Type(), param.DType) {
				return false
			}
		}
	}

	return true
}

//FIXME move these

func CanConvert(argType DataType, paramType DataType) bool {

	//FIXME fake
	return argType == paramType
}

func TypeMatches(a DataType, b DataType) bool {
	//FIXME fake
	return a == b
}

