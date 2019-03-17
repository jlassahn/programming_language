
package generator

import (
	"parser"
	"symbols"
)


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

	sym, err := ctx.Symbols.AddVar(name, dtype)
	if err != nil {
		parser.Error(el.FilePos(), "%v", err)
		return nil
	}

	genVal := NewLocalVal(genFunc.File(), dtype, name)
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

	ret:= loopHandler(genFunc, ctx, el)

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

