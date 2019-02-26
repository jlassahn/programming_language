
package generator

import (
	"output"
	"parser"
	"symbols"
)


func genAssignment(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	lhsEl := el.Children()[0]
	rhsEl := el.Children()[1]

	lhs := loopHandler(genFunc, ctx, lhsEl)
	if lhs == nil {
		return nil
	}

	dtype := lhs.Type()
	rhs := genRHS(genFunc, ctx, dtype, rhsEl)
	if rhs == nil {
		return nil
	}


	if lhs.IsVariableRef() {
		output.FIXMEDebug("lhs is variable ref")
	}
	if lhs.IsGlobalRef() {
		output.FIXMEDebug("lhs is global def")
	}
	output.FIXMEDebug("lhs = %v", lhs.LLVMVal())
	output.FIXMEDebug("rhs = %v", rhs.LLVMVal())

	genFunc.AddBody("\tstore %v %v, %v* %v",
		rhs.LLVMType(),
		rhs.LLVMVal(),
		lhs.LLVMType(),
		lhs.LLVMVal())

	return nil
}

