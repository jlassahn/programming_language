
package symbols

import (
	"output"
	"parser"
)


func evalList(el parser.ParseElement, ctx *EvalContext) DataValue {

	baseType := ctx.InitializerType
	if baseType == nil {
		parser.Error(el.FilePos(), "list initializer needs a type")
		return nil
	}

	// FIXME handle other list types, like Array or Slice
	if baseType.Base() != MARRAY_TYPE {
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

	ret := &listDV{
		dtype: baseType,
		value: make([]DataValue, length),
	}

	ctx.InitializerType = subtype
	for i, x := range children {
		dval := EvaluateConstRHS(x, ctx)
		output.FIXMEDebug("got element %d: %v", i, dval)
		ret.value[i] = dval
	}
	ctx.InitializerType = baseType

	//FIXME
	output.FIXMEDebug("FIXME implement array literal value")
	return ret
}
