
package symbols

import (
	"parser"
)


func evalFunctionContent(el parser.ParseElement, ctx *EvalContext) DataValue {

	if ctx.InitializerType == nil {
		parser.Error(el.FilePos(), "function body on value with no type")
		return nil
	}
	if ctx.InitializerType.Base() != FUNCTION_TYPE {
		parser.Error(el.FilePos(), "function body on non-function value")
		return nil
	}

	return &codeDV {
		dtype: ctx.InitializerType,
		element: el,
		file: el.FilePos().File.(*SourceFile),
	}
}

