
package symbols

import (
	"output"
	"parser"
)


type FunctionTypeEval struct {}
func (*FunctionTypeEval) EvaluateConstExpression(
		el parser.ParseElement, ctx *EvalContext) DataValue {

	paramList := el.Children()[0]
	retType := el.Children()[1]

	ret := &functionDT{}

	for _,param := range paramList.Children() {
		paramName := param.Children()[0]
		paramType := param.Children()[1]
		name := paramName.TokenString()

		if ctx.Symbols.Lookup(name) != nil {
			parser.Error(paramName.FilePos(),
				"parameter name already defined: %v", name)
			return nil
		}

		dtypeExp := EvaluateConstExpression(paramType, ctx)
		if dtypeExp == nil {
			parser.Error(paramName.FilePos(), "unknown data type")
			return nil
		}
		dtype := dtypeExp.(TypeDataValue).AsDataType()

		autoconv := false //FIXME get autoconvert flag!

		ret.parameters = append(ret.parameters,  FunctionParameter {
			Name: name,
			DType: dtype,
			AutoConvert: autoconv,
		})

		output.FIXMEDebug("param: %v %v", name, paramType)
	}

	output.FIXMEDebug("return: %v", retType)

	var dtypeExp DataValue
	if retType.ElementType() == parser.EMPTY {
		dtypeExp = &typeDV{CTypeType, VoidType}
	} else {
		dtypeExp = EvaluateConstExpression(retType, ctx)
	}

	if dtypeExp == nil {
		parser.Error(retType.FilePos(), "unknown data type")
		return nil
	}
	dtype := dtypeExp.(TypeDataValue).AsDataType()
	ret.returnType = dtype

	return &typeDV {
		dtype: CTypeType,
		value: ret,
	}
}

