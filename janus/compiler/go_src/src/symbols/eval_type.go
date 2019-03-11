
package symbols

import (
	"output"
	"parser"
)


func evalType(el parser.ParseElement, ctx *EvalContext) DataValue {

	tp := loopHandler(el.Children()[0], ctx)

	if tp.Type() == CTypeType {
		if len(el.Children()) != 1 {
			parser.Error(el.FilePos(),
				"%s does not take parameters", el.TokenString())
		}
		return tp

	} else if tp.Type() == MetaTypeType {

		if len(el.Children()) != 2 {
			parser.Error(el.FilePos(),
				"%s needs parameters", el.TokenString())
			return nil
		}

		// FIXME handle special case for initializer with
		// implicit size param.

		mtype := tp.(TypeDataValue).AsDataType()
		paramEls := el.Children()[1].Children()
		if len(paramEls) != len(mtype.SubTypes()) {
			parser.Error(el.FilePos(),
				"%s needs %d parameters",
				el.TokenString(), len(mtype.SubTypes()))
			return nil
		}

		params := make([]DataValue, len(paramEls))
		for i, el := range paramEls {
			p := loopHandler(el, ctx)
			if p == nil {
				return nil
			}
			params[i] = p
		}

		dtype := insertParameters(mtype, params)
		if dtype == nil {
			return nil
		}
		return &typeDV{CTypeType, dtype}

	} else {
		parser.Error(el.FilePos(), "%s is not a type", el.TokenString())
		return nil
	}

	return tp
}

func insertParameters(mtype DataType, paramsIn []DataValue) DataType {

	pdt, ok := mtype.(*paramDT)
	if !ok {
		return nil
	}

	// FIXME handle struct types, etc
	// FIXME handle members

	paramsOut := make([]DTypeParameter, len(paramsIn))

	for i := range paramsIn {
		mparam := pdt.params[i].DType.(*typevarDT)
		if mparam.numeric {
			dv := ConvertConstant(paramsIn[i], Int64Type)
			if dv == nil {
				return nil
			}
			paramsOut[i].DType = nil
			paramsOut[i].Number = dv.(SignedDataValue).AsSigned64()

		} else {
			dv := ConvertConstant(paramsIn[i], CTypeType)
			if dv == nil {
				return nil
			}
			paramsOut[i].DType = dv.(TypeDataValue).AsDataType()
			paramsOut[i].Number = 0
		}
	}

	if pdt.members != nil {
		output.FatalError("FIXME unimplemented parameterized type with members")
		return nil
	}

	return &paramDT{
		tag: pdt.tag,
		params: paramsOut,
		members: nil,
	}
}
