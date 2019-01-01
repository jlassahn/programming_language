
package generator

import (
	"parser"
	"symbols"
)


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

