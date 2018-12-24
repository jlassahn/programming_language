
package generator

import (
	"fmt"

	"output"
	"parser"
	"symbols"
)

type ExpressionGenerator func(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result

var handlers = map[*parser.Tag] ExpressionGenerator {
	parser.EXPRESSION: genExpression,
	parser.NUMBER: genNumber,
	parser.SYMBOL: genSymbol,
	parser.RETURN: genReturn,
}

// indirect call to GenerateExpression, to avoid a circular dependency
var loopHandler ExpressionGenerator



func GenerateExpression(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	output.FIXMEDebug("GenerateExpression: %v", el)

	loopHandler = GenerateExpression

	handler := handlers[el.ElementType()]
	if handler == nil {
		output.FIXMEDebug("no expression generator for type %v", el.ElementType())
		return nil
	}

	return handler(genFunc, ctx, el)
}


func genExpression(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	children := el.Children()
	opElement := children[0]
	opName := opElement.TokenString()

	args := make([]Result, len(children) - 1)
	for i, x := range(children[1:]) {
		args[i] = loopHandler(genFunc, ctx, x)
		if args[i] == nil {
			output.FIXMEDebug("FIXME args not available")
			return nil
		}
	}

	if opElement.ElementType() == parser.OPERATOR {
		return genOperator(genFunc, ctx, opElement, args)
	}

	//FIXME implement
	output.FIXMEDebug("applying %v to %v", opName, args)
	return nil
}

func genOperator(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement,
	args []Result) Result {

	fp := genFunc.File()
	opName := el.TokenString()
	pos := el.FilePos()

	argTypes := make([]symbols.DataType, len(args))
	for i,x := range args {
		argTypes[i] = x.Type()
	}

	//FIXME make this a reusable function
	opChoices := ctx.LookupOperator(opName)
	if opChoices == nil {
		parser.Error(pos, "No definition for operator %v", opName)
		//FIXME testing
		ctx.Symbols.Emit(true)
		return nil
	}

	//FIXME below here is the same as non-operator functions
	op := symbols.SelectFunctionChoice(opChoices, argTypes)
	if op == nil {
		parser.Error(pos, "Operator %v can't take these parameters", opName)
		return nil
	}

	convertedArgs := make([]Result, len(args))
	dtype := op.Type().(symbols.FunctionDataType)
	for i, dest := range dtype.Parameters() {
		convertedArgs[i] = ConvertParameter(genFunc, args[i], dest.DType)
	}

	retType := op.Type().(symbols.FunctionDataType).ReturnType()

	if op.InitialValue().Type() == symbols.IntrinsicType {
		opName := op.InitialValue().(symbols.IntrinsicDataValue).ValueAsString()
		ret := NewTempVal(fp, retType)

		output.FIXMEDebug("applying intrinsic %v to %v", opName, convertedArgs)

		genFunc.AddBody("%v", MakeIntrinsicOp(ret, opName, convertedArgs))
		return ret
	}

	output.FIXMEDebug("no support for operator %v", op)

	//FIXME handle cases...
	//   const intrinsic
	//   var intrinsic
	//   const code
	//   var code

	// FIXME
	// if op is of type CODE
	// cdv = op.InitialValue().(CodeDataValue)
	//  info to build a global reference:
	//  DataType    == op.Type()
	//  name        == op.Name()
	//  module path == cdv.AsSourceFile().Options.ModuleName

	return nil
}

func genSymbol(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	sym := ctx.Lookup(el.TokenString())
	output.FIXMEDebug("looking up %v %v", el.TokenString(), sym)

	ret, ok := sym.GetGenVal().(Result)
	if ok {
		output.FIXMEDebug("found value %v", ret)

		return ret
	}

	output.FIXMEDebug("NO VALUE FOUND")
	return nil
}

func genReturn(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	dtype := genFunc.ReturnType()

	argEl := el.Children()[0]

	if argEl.ElementType() == parser.EMPTY {
		if dtype != symbols.VoidType {
			parser.Error(el.FilePos(),
				"return value is void should be %v", dtype)
		}
		genFunc.AddBody("\tret void")
		return nil
	}

	arg := loopHandler(genFunc, ctx, argEl)
	if arg == nil {
		return nil
	}

	if !symbols.CanConvert(arg.Type(), dtype) {
		parser.Error(el.FilePos(),
				"return value is %v should be %v", arg.Type(), dtype)
		return nil
	}

	convertedArg := ConvertParameter(genFunc, arg, dtype)

	genFunc.AddBody("\tret %v %v",
		convertedArg.LLVMType(),
		convertedArg.LLVMVal())
	return nil
}

func genNumber(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	dv := symbols.EvaluateConstExpression(el, ctx)
	if dv == nil {
		return nil
	}
	return NewDataVal(dv)
}


//FIXME simplify control flow
func ConvertParameter(genFunc GeneratedFunction,
	arg Result, dtype symbols.DataType) Result {

	if arg.IsVariableRef() {
		arg = DereferenceVariable(genFunc, arg)
	}

	arg = ConvertValue(genFunc, arg, dtype)
	if arg == nil {
		output.Error("INTERNAL ERROR")
		return nil
	}

	return arg
}

func DereferenceVariable(genFunc GeneratedFunction, src Result) Result {

	fp := genFunc.File()
	ret := NewTempVal(fp, src.Type())

	genFunc.AddBody("\t%v = load %v, %v* %v",
		ret.LLVMVal(),
		ret.LLVMType(),
		src.LLVMType(),
		src.LLVMVal())

	return ret
}

func ConvertValue(genFunc GeneratedFunction,
	from Result, to symbols.DataType) Result {

	if symbols.TypeMatches(from.Type(), to) {
		return from
	}

	if from.IsConst() {
		dval := from.ConstVal()
		dval = symbols.ConvertConstant(dval, to)
		return NewDataVal(dval)
	}

	fp := genFunc.File()
	opString, ok := baseTypeConvert[tagPair{from.Type().Base(), to.Base()}]
	if ok {
		ret := NewTempVal(fp, to)

		genFunc.AddBody("\t%v = %s %v %v to %v",
			ret.LLVMVal(),
			opString,
			from.LLVMType(),
			from.LLVMVal(),
			ret.LLVMType())

		return ret
	}

	//FIXME handle composite types, etc

	//FIXME internal error, should have been checked before getting here
	output.Error("no conversion from %v to %v", from, to)
	return nil 
}

func MakeIntrinsicOp(ret Result, opName string, args []Result) string {

	op, ok := LLVMOperator[opName]
	if !ok {
		output.FatalError("Unimplemented intrinsic %v", opName)
		op = "UNIMPLEMENTED"
	}

	s := fmt.Sprintf("\t%v = %v %v %v, %v",
		ret.LLVMVal(),
		op,
		args[0].LLVMType(),
		args[0].LLVMVal(),
		args[1].LLVMVal())

	return s
}

var LLVMOperator = map[string]string {
	"add_Int64": "add",
	"add_Int32": "add",
}


//FIXME reorganize
type tagPair struct {
	from *symbols.Tag
	to *symbols.Tag
}

var baseTypeConvert = map[tagPair] string  {
	{symbols.INT8_TYPE, symbols.INT16_TYPE}: "sext",
	{symbols.INT8_TYPE, symbols.INT32_TYPE}: "sext",
	{symbols.INT8_TYPE, symbols.INT64_TYPE}: "sext",
	{symbols.INT16_TYPE, symbols.INT32_TYPE}: "sext",
	{symbols.INT16_TYPE, symbols.INT64_TYPE}: "sext",
}

