
package generator

import (
	"output"
	"parser"
	"symbols"
)

type ExpressionGenerator func(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result

//FIXME organize
var handlers = map[*parser.Tag] ExpressionGenerator {
	parser.EXPRESSION: genExpression,
	parser.NUMBER: genNumber,
	parser.SYMBOL: genSymbol,
	parser.RETURN: genReturn,
	parser.DEF: genDef,
	parser.CALL: genCall,
	parser.DOT_LIST: genDotList,
	parser.IF: genIf,
	parser.FUNCTION_CONTENT: genFunctionContent,
	parser.ASSIGNMENT: genAssignment,
	parser.WHILE: genWhile,
}

// indirect call to GenerateExpression, to avoid a circular dependency
var loopHandler ExpressionGenerator



func GenerateExpression(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	output.FIXMEDebug("GenerateExpression: %v", el)

	loopHandler = GenerateExpression

	handler := handlers[el.ElementType()]
	if handler == nil {
		output.FatalError("no expression generator for type %v", el.ElementType())
		return nil
	}

	return handler(genFunc, ctx, el)
}


func genNumber(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	dv := symbols.EvaluateConstExpression(el, ctx)
	if dv == nil {
		return nil
	}
	return NewDataVal(dv)
}

//FIXME rename -- this converts and loads any data access

func ConvertParameter(genFunc GeneratedFunction,
	arg Result, dtype symbols.DataType) Result {

	if arg.IsVariableRef() || arg.IsGlobalRef() {
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
	ret := NewTempVal(fp, to)

	opString := MakeLLVMConvert(from, ret)
	if opString != "" {
		genFunc.AddBody("%v", opString)
		return ret
	}

	//FIXME handle composite types, etc

	//FIXME internal error, should have been checked before getting here
	output.Error("no conversion from %v to %v", from, to)
	return nil 
}

