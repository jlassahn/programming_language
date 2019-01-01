
package generator

import (
	"parser"
	"symbols"
)


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

