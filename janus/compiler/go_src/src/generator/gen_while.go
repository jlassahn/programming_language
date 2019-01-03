
package generator

import (
	"parser"
	"symbols"
)


func genWhile(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	testTree := el.Children()[0]
	bodyTree := el.Children()[1]

	startLabel := NewTempVal(genFunc.File(), symbols.LabelType)
	bodyLabel := NewTempVal(genFunc.File(), symbols.LabelType)
	endLabel := NewTempVal(genFunc.File(), symbols.LabelType)

	genFunc.AddBody("\tbr label %v", startLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", startLabel.Name(), startLabel.ID())

	testVal := loopHandler(genFunc, ctx, testTree)
	if testVal == nil {
		return nil
	}
	if !symbols.CanConvert(testVal.Type(), symbols.BoolType) {
		parser.Error(testTree.FilePos(), "not a boolean expression")
		return nil
	}
	testVal = ConvertParameter(genFunc, testVal, symbols.BoolType)

	genFunc.AddBody("\tbr i1 %v, label %v, label %v",
		testVal.LLVMVal(),
		bodyLabel.LLVMVal(),
		endLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", bodyLabel.Name(), bodyLabel.ID())

	loopHandler(genFunc, ctx, bodyTree)

	genFunc.AddBody("\tbr label %v", startLabel.LLVMVal())
	genFunc.AddBody("%s_%d:", endLabel.Name(), endLabel.ID())

	return nil
}

