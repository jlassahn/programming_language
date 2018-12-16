
package symbols

import (
	"parser"
)


func evalFunctionContent(el parser.ParseElement, ctx *EvalContext) DataValue {


	return &codeDV {
		dtype: CodeType,
		element: el,
		file: el.FilePos().File.(*SourceFile),
	}
}

