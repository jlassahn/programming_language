
package generator

import (
	"parser"
	"symbols"
)


func genList(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	baseType := ctx.InitializerType
	if baseType == nil {
		parser.Error(el.FilePos(), "list initializer needs a type")
		return nil
	}

	// FIXME handle other list types, like Array or Slice
	if baseType.Base() != symbols.MARRAY_TYPE {
		parser.Error(el.FilePos(), "list initializer needs an array")
		return nil
	}

	subtype := baseType.SubTypes()[0].DType
	length := baseType.SubTypes()[1].Number

	children := el.Children()
	if int64(len(children)) > length {
		parser.Error(el.FilePos(), "too many initializers")
		return nil
	}

	ret := NewZeroVal(baseType)

	ctx.InitializerType = subtype
	for i, child := range children {

		childVal := loopHandler(genFunc, ctx, child)
		if childVal == nil {
			return nil
		}

		if !symbols.CanConvert(childVal.Type(), subtype) {
			parser.Error(el.FilePos(), "can't convert value to %v", subtype)
			return nil
		}
		val := ConvertParameter(genFunc, childVal, subtype)

		newRet := NewTempVal(genFunc.File(), baseType)
		genFunc.AddBody("\t%v = insertvalue %v %v, %v %v, %v",
			newRet.LLVMVal(),
			ret.LLVMType(),
			ret.LLVMVal(),
			val.LLVMType(),
			val.LLVMVal(),
			i)

		ret = newRet
	}
	ctx.InitializerType = baseType

	return ret
}

