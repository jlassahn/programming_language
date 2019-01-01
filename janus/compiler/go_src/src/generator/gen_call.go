
package generator

import (
	"fmt"

	"output"
	"parser"
	"symbols"
)


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

	if opResult.IsFunctionChoice() {
		functionChoice := opResult.FunctionChoice()
		opName := functionChoice.Name()
		fnSym := symbols.SelectFunctionChoice(functionChoice, argTypes)
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

	//FIXME clean up naming
	op := opResult

	convertedArgs := make([]Result, len(args))
	dtype := op.Type().(symbols.FunctionDataType)
	for i, dest := range dtype.Parameters() {
		convertedArgs[i] = ConvertParameter(genFunc, args[i], dest.DType)
	}

	retType := dtype.ReturnType()

	ret := NewTempVal(fp, retType)

	if op.IsConst() && op.ConstVal().Tag() == symbols.INTRINSIC_VALUE {
		opName := op.ConstVal().(symbols.IntrinsicDataValue).ValueAsString()

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

