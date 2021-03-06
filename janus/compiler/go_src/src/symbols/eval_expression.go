
package symbols

import (
	"output"
	"parser"
)


func evalExpression(el parser.ParseElement, ctx *EvalContext) DataValue {

	children := el.Children()
	opElement := children[0]
	opName := opElement.TokenString()
	pos := opElement.FilePos()

	if opElement.ElementType() != parser.OPERATOR {
		parser.Error(pos, "FIXME not an operator: %v", opName)
		return nil
	}

	args := make([]DataValue, len(children) - 1)
	argTypes := make([]DataType, len(children) - 1)
	for i, x := range(children[1:]) {
		args[i] = loopHandler(x, ctx)
		if args[i] == nil {
			output.FIXMEDebug("FIXME args not available")
			return nil
		}
		argTypes[i] = args[i].Type()
	}

	opChoices := ctx.LookupOperator(opName)
	if opChoices == nil {
		parser.Error(pos, "No definition for operator %v", opName)
		//FIXME testing
		//ctx.Symbols.Emit()
		return nil
	}

	op := SelectFunctionChoice(opChoices, argTypes)
	if op == nil {
		parser.Error(pos, "Operator %v can't take these parameters", opName)
		return nil
	}

	if !op.IsConst() {
		parser.Error(pos, "Operator %v not const", opName)
		return nil
	}

	ret := doConstOp(op, args, ctx)

	return ret
}

func doConstOp(op Symbol,
	args []DataValue, ctx *EvalContext) DataValue {

	convertedArgs := make([]DataValue, len(args))

	dtype := op.Type().(FunctionDataType)
	for i, dest := range dtype.Parameters() {
		convertedArgs[i] = ConvertConstant(args[i], dest.DType)
		if convertedArgs[i] == nil {
			output.Error("INTERNAL ERROR")
			return nil
		}
	}

	//FIXME check if this is really an intrinsic
	opName := op.InitialValue().(IntrinsicDataValue).ValueAsString()
	return EvaluateIntrinsic(opName, convertedArgs)
}

//FIXME where should this live?
func SelectFunctionChoice(op FunctionChoiceSymbol, args []DataType) Symbol {

	for _,choice := range op.Choices() {

		dtype := choice.Type().(FunctionDataType)
		params := dtype.Parameters()

		if !CanConvertArgs(args, params) { continue; }

		return choice
	}

	return nil
}

