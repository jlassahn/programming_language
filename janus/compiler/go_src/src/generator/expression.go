
package generator

import (
	"fmt"

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


func genCall(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	children := el.Children()
	opElement := children[0]
	argList := children[1].Children()

	opResult := loopHandler(genFunc, ctx, opElement)
	if opResult == nil {
		output.FIXMEDebug("FIXME opResult lookup failed")
		return nil
	}

	//FIXME
	//if opResult is method call
	//  twiddle args to form non-method version

	return genInvokeFunction(genFunc, ctx, opResult, argList)
}

//FIXME clean this up

func genInvokeFunction(
	genFunc GeneratedFunction,
	ctx *symbols.EvalContext,
	opResult Result,
	argList []parser.ParseElement,
) Result {

	fp := genFunc.File()

	args := make([]Result, len(argList))
	argTypes := make([]symbols.DataType, len(argList))
	for i, x := range(argList) {
		args[i] = loopHandler(genFunc, ctx, x)
		if args[i] == nil {
			output.FIXMEDebug("FIXME args not available")
			return nil
		}
		argTypes[i] = args[i].Type()
	}

	//FIXME make this a function?
	if opResult.IsFunctionChoice() {
		functionChoice := opResult.FunctionChoice()
		opName := functionChoice.Name()
		fnSym := symbols.SelectFunctionChoice(functionChoice, argTypes)
		output.FIXMEDebug("FIXME resolve function choice for %v = %v", opName, fnSym)
		if fnSym == nil {
			parser.Error(argList[0].FilePos(),
				"Operator %v can't take these parameters", opName)
			return nil
		}


		if fnSym.GetGenVal() != nil {
			opResult = fnSym.GetGenVal().(Result)

		} else if fnSym.InitialValue().Tag() == symbols.INTRINSIC_VALUE {
			//FIXME do we need Typed... ? InitialValue should have the right type.
			opResult = NewTypedDataVal(fnSym.Type(), fnSym.InitialValue())

		} else if fnSym.InitialValue().Tag() == symbols.CODE_VALUE  {

			//FIXME messy, should have something like Symbol.ModulePath()

			modPath := fnSym.InitialValue().(symbols.CodeDataValue).AsSourceFile().Options.ModuleName
			name := MakeSymbolName(modPath, fnSym.Type(), fnSym.Name())
			opResult = NewGlobalVal(fp, fnSym.Type(), name)

		} else {
			output.FIXMEDebug("NO GENVAL FOR FUNCTION %v", fnSym)
			return nil
		}
	}

	output.FIXMEDebug("applying %v to %v", opResult.Name(), args)
	op := opResult

	//FIXME combine with genOperator
	output.FIXMEDebug("dtype = %v", op.Type())

	//fp := genFunc.File()
	convertedArgs := make([]Result, len(args))
	dtype := op.Type().(symbols.FunctionDataType)
	for i, dest := range dtype.Parameters() {
		convertedArgs[i] = ConvertParameter(genFunc, args[i], dest.DType)
	}

	retType := dtype.ReturnType()

	//FIXME how do intrinsics work here?
	//   intrinsics have symbols with DataValue of type IntrinsicDataValue
	//   we can put that DataValue in a Result....

	ret := NewTempVal(fp, retType)

	if op.IsConst() && op.ConstVal().Tag() == symbols.INTRINSIC_VALUE {
		opName := op.ConstVal().(symbols.IntrinsicDataValue).ValueAsString()

		output.FIXMEDebug("applying intrinsic %v to %v", opName, convertedArgs)

		genFunc.AddBody("%v", MakeIntrinsicOp(ret, opName, convertedArgs))
		return ret
	}


	var callStr string
	if ret.Type() == symbols.VoidType {
		callStr = fmt.Sprintf("\tcall %v %v(",
			ret.LLVMType(),
			op.LLVMVal())
	} else {
		callStr = fmt.Sprintf("\t%v = call %v %v(",
			ret.LLVMVal(),
			ret.LLVMType(),
			op.LLVMVal())
	}

	for i,arg := range args {
		if i > 0 {
			callStr = callStr + ", "
		}
		callStr = callStr + fmt.Sprintf("%v %v",
			arg.LLVMType(),
			arg.LLVMVal())
	}
	callStr = callStr + ")"

	genFunc.AddBody("%v", callStr)
	return ret

	//FIXME handle cases...
	//   const intrinsic
	//   var intrinsic
	//   const code
	//   var code

}


func genExpression(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	children := el.Children()
	opElement := children[0]

	if opElement.ElementType() == parser.OPERATOR {

		opName := opElement.TokenString()
		opChoices := ctx.LookupOperator(opName)
		if opChoices == nil {
			parser.Error(el.FilePos(), "no definition for operator %v", opName)
			return nil
		}

		opResult := NewFunctionChoiceResult(opChoices)
		argList := children[1:]
		return genInvokeFunction(genFunc, ctx, opResult, argList)
	}

	//FIXME implement  (what non-operator expressions are there???)
	output.FIXMEDebug("expression with non-operator %v", opElement)
	return nil
}

func genSymbol(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	sym := ctx.Lookup(el.TokenString())
	if sym == nil {
		parser.Error(el.FilePos(), "undefined symbol %v", el.TokenString())
		return nil
	}

	if sym.Type() == symbols.FunctionChoiceType {
		fn := sym.(symbols.FunctionChoiceSymbol)
		return NewFunctionChoiceResult(fn)
	}

	output.FIXMEDebug("looking up %v %v", el.TokenString(), sym)
	output.FIXMEDebug("Symbol %v %v %v", sym.Name(), sym.Type(), sym.InitialValue())
	
	ret, ok := sym.GetGenVal().(Result)
	if ok {
		output.FIXMEDebug("found value %v", ret)

		return ret
	}

	if sym.IsConst() {
		return NewDataVal(sym.InitialValue())
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

	label := NewTempVal(genFunc.File(), symbols.LabelType)
	genFunc.AddBody("%s_%d:", label.Name(), label.ID())

	return label
}

func genNumber(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	dv := symbols.EvaluateConstExpression(el, ctx)
	if dv == nil {
		return nil
	}
	return NewDataVal(dv)
}

func genDef(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	name := el.Children()[0].TokenString()
	typeTree := el.Children()[1]
	valTree := el.Children()[2]

	var dtype symbols.DataType
	if typeTree.ElementType() != parser.EMPTY {
		dval := symbols.EvaluateConstExpression(typeTree, ctx)
		if dval == nil {
			return nil
		}

		if dval.Tag() != symbols.TYPE_VALUE {
			parser.Error(typeTree.FilePos(), "not a data type")
			return nil
		}

		symTypeVal := dval.(symbols.TypeDataValue)
		dtype = symTypeVal.AsDataType()
	}

	var dval Result

	if valTree.ElementType() != parser.EMPTY {
		dval = genRHS(genFunc, ctx, dtype, valTree)
	}

	if dtype == nil && dval != nil {
		dtype = dval.Type()
	}

	if dtype == nil {
		parser.Error(el.FilePos(), "unable to infer data type for %v", name)
		return nil
	}

	output.FIXMEDebug("def: %v %v %v", name, dtype, dval)

	sym, err := ctx.Symbols.AddVar(name, dtype)
	if err != nil {
		parser.Error(el.FilePos(), "%v", err)
		return nil
	}

	genVal := NewNamedVal(genFunc.File(), dtype, name)
	sym.SetGenVal(genVal)

	genFunc.AddPrologue("\t%v = alloca %v", genVal.LLVMVal(), genVal.LLVMType())

	if dval == nil {
		genFunc.AddBody("\tstore %v zeroinitializer, %v* %v",
			genVal.LLVMType(),
			genVal.LLVMType(),
			genVal.LLVMVal())
	} else {
		genFunc.AddBody("\tstore %v %v, %v* %v",
			genVal.LLVMType(),
			dval.LLVMVal(),
			genVal.LLVMType(),
			genVal.LLVMVal())
	}

	return nil
}

func genRHS(
	genFunc GeneratedFunction,
	ctx *symbols.EvalContext,
	dtype symbols.DataType,
	el parser.ParseElement,
) Result {

	//FIXME handle block initializers

	ctx.InitializerType = dtype

	output.FIXMEDebug("starting genRHS %v", el)
	ret:= loopHandler(genFunc, ctx, el)
	output.FIXMEDebug("finished genRHS %v", el)

	ctx.InitializerType = nil

	//do type conversions 
	if ret != nil {

		if dtype == nil {
			dtype = ret.Type()
		} else if !symbols.CanConvert(ret.Type(), dtype) {
			parser.Error(el.FilePos(), "can't convert value to %v", dtype)
			return nil
		}
		ret = ConvertParameter(genFunc, ret, dtype)
	}

	return ret
}

func genDotList(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	lhsEl := el.Children()[0]
	rhsEl := el.Children()[1]

	lhs := loopHandler(genFunc, ctx, lhsEl)
	if lhs == nil {
		return nil
	}

	output.FIXMEDebug("getDotList %v %v", lhsEl, rhsEl)

	if lhs.Type() == symbols.NamespaceType {
		table := lhs.ConstVal().(symbols.NamespaceDataValue).AsSymbolTable()

		impCtx := symbols.EvalContext { }
		impCtx.Symbols = table

		output.FIXMEDebug("push module context %v", table)
		ret := loopHandler(genFunc,&impCtx, rhsEl)
		output.FIXMEDebug("pop module context %v", table)

		return ret
	}

	//FIXME implement member access, etc
	parser.Error(el.FilePos(), "FIXME unimplemented dot operator %v, rhsEl")
	return nil
}

func genIf(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	testTree := el.Children()[0]
	trueTree := el.Children()[1]
	falseTree := el.Children()[2]

	testVal := loopHandler(genFunc, ctx, testTree)
	if testVal == nil {
		return nil
	}
	if !symbols.CanConvert(testVal.Type(), symbols.BoolType) {
		parser.Error(testTree.FilePos(), "not a boolean expression")
		return nil
	}
	testVal = ConvertParameter(genFunc, testVal, symbols.BoolType)

	trueLabel := NewTempVal(genFunc.File(), symbols.LabelType)
	falseLabel := NewTempVal(genFunc.File(), symbols.LabelType)
	endLabel := NewTempVal(genFunc.File(), symbols.LabelType)

	genFunc.AddBody("\tbr i1 %v, label %v, label %v",
		testVal.LLVMVal(),
		trueLabel.LLVMVal(),
		falseLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", trueLabel.Name(), trueLabel.ID())

	loopHandler(genFunc, ctx, trueTree)

	genFunc.AddBody("\tbr label %v", endLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", falseLabel.Name(), falseLabel.ID())

	if falseTree.ElementType() != parser.EMPTY {
		loopHandler(genFunc, ctx, falseTree)
	}

	genFunc.AddBody("\tbr label %v", endLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", endLabel.Name(), endLabel.ID())

	return nil
}

func genFunctionContent(genFunc GeneratedFunction,
	ctxBase *symbols.EvalContext, el parser.ParseElement) Result {

	
	symbolTable := symbols.NewSymbolTable(
		fmt.Sprintf("local@%d", el.FilePos().Line),
		ctxBase.Symbols)

	ctx := &symbols.EvalContext {
		Symbols: symbolTable,
	}

	for _,elem := range el.Children() {
		loopHandler(genFunc, ctx, elem)
	}

	return nil
}

//FIXME rename -- this converts and loads any data access

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
	if ok {
		return  fmt.Sprintf("\t%v = %v %v %v, %v",
			ret.LLVMVal(),
			op,
			args[0].LLVMType(),
			args[0].LLVMVal(),
			args[1].LLVMVal())
	}

	//FIXME this could be done by using the normal function call path
	//   and providing a Result with name = LLVMFunction[opName]
	//   which could be injected into the intrinsic by SetGenVal
	op, ok = LLVMFunction[opName]
	if ok {

		var s string
		if ret.Type() == symbols.VoidType {
			s = fmt.Sprintf("\tcall %v %v(",
				ret.LLVMType(),
				op)
		} else {
			s = fmt.Sprintf("\t%v = call %v %v(",
				ret.LLVMVal(),
				ret.LLVMType(),
				op)
		}

		for i,arg := range args {
			if i > 0 {
				s = s + ", "
			}
			s = s + fmt.Sprintf("%v %v", arg.LLVMType(), arg.LLVMVal())
		}
		s = s + ")"

		return s
	}

	output.FatalError("Unimplemented intrinsic %v", opName)

	return ""
}

var LLVMOperator = map[string]string {
	//FIXME add a bunch of stuff here
	"add_Int64": "add",
	"add_Int32": "add",

	"sub_Int64": "sub",

	"mul_Int64": "mul",

	"cmp_le_Int64": "icmp sle",
}

var LLVMFunction = map[string]string {
	"sqrt_Real64": "@llvm.sqrt.f64",
	"print_Real64": "@clib_print_Real64",
	"print_Int64": "@clib_print_Int64",
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

