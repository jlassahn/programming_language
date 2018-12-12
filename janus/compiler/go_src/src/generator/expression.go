
package generator

import (
	"fmt"

	"output"
	"parser"
	"symbols"
)

type ExpressionGenerator func(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) GeneratedStatement

var handlers = map[*parser.Tag] ExpressionGenerator {
	parser.EXPRESSION: genExpression,
	parser.SYMBOL: genSymbol,
}

// indirect call to GenerateExpression, to avoid a circular dependency
var loopHandler ExpressionGenerator



func GenerateExpression(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) GeneratedStatement {

	output.FIXMEDebug("GenerateExpression: %v", el)

	loopHandler = GenerateExpression

	handler := handlers[el.ElementType()]
	if handler == nil {
		output.FIXMEDebug("no expression generator for type %v", el.ElementType())
		return nil
	}

	return handler(fp, genFunc, ctx, el)
}


func genExpression(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) GeneratedStatement {

	children := el.Children()
	opElement := children[0]
	opName := opElement.TokenString()

	args := make([]GeneratedStatement, len(children) - 1)
	for i, x := range(children[1:]) {
		args[i] = loopHandler(fp, genFunc, ctx, x)
		if args[i] == nil {
			output.FIXMEDebug("FIXME args not available")
			return nil
		}
	}


	if opElement.ElementType() == parser.OPERATOR {
		return genOperator(fp, genFunc, ctx, opElement, args)
	}

	//FIXME fake, keywords should be handled differently in parser!
	if opElement.ElementType() == parser.KEYWORD {
		if opElement.TokenString() == "return" {
			genFunc.AddBody("\tret %v %v",
				MakeLLVMType(args[0].Result().Type()),
				args[0].String())
		}
	}

	//FIXME implement
	output.FIXMEDebug("applying %v to %v", opName, args)
	return nil
}

func genOperator(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement,
	args []GeneratedStatement) GeneratedStatement {

	opName := el.TokenString()
	pos := el.FilePos()

	argTypes := make([]symbols.DataType, len(args))
	for i,x := range args {
		argTypes[i] = x.Result().Type()
	}

	//FIXME make this a reusable function
	opChoices := ctx.Symbols.LookupOperator(opName)
	if opChoices == nil {
		parser.Error(pos, "No definition for operator %v", opName)
		//FIXME testing
		ctx.Symbols.Emit()
		return nil
	}

	//FIXME below here is the same as non-operator functions
	op := symbols.SelectFunctionChoice(opChoices, argTypes)
	if op == nil {
		parser.Error(pos, "Operator %v can't take these parameters", opName)
		return nil
	}

	//FIXME handle type conversions for args here

	retType := op.Type().(symbols.FunctionDataType).ReturnType()

	if op.InitialValue().Type() == symbols.IntrinsicType {
		opName := op.InitialValue().(symbols.IntrinsicDataValue).ValueAsString()
		ret := NewTempVal(fp, retType)

		output.FIXMEDebug("applying intrinsic %v to %v", opName, args)

		genFunc.AddBody("%v", MakeIntrinsicOp(ret, opName, args))
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

func genSymbol(fp GeneratedFile, genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) GeneratedStatement {

	sym := ctx.Symbols.Lookup(el.TokenString())
	output.FIXMEDebug("looking up %v %v", el.TokenString(), sym)

	src, ok := sym.GetGenVal().(GeneratedStatement)
	if ok {
		output.FIXMEDebug("found value %v", src)

		ret := NewTempVal(fp, sym.Type())

		genFunc.AddBody("\t%v = load %v, %v* %v",
			ret.Result().String(),
			MakeLLVMType(ret.Result().Type()),
			MakeLLVMType(src.Result().Type()),
			src.Result().String())

		return ret
	}

	output.FIXMEDebug("NO VALUE FOUND")
	return nil
}

func MakeIntrinsicOp(ret GeneratedStatement, opName string,
	args []GeneratedStatement) string {

	s := fmt.Sprintf("\t%v = %v %v %v, %v",
		ret.Result().String(),
		LLVMOperator[opName],
		MakeLLVMType(args[0].Result().Type()),
		args[0].String(),
		args[1].String())

	output.FIXMEDebug("%v", s)
	return s
}

var LLVMOperator = map[string]string {
	"add_Int64": "add",
}

