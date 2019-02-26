
package symbols

import (
	"parser"
)


func evalType(el parser.ParseElement, ctx *EvalContext) DataValue {

	//FIXME handle parameterized types
	//      Children()[0] is base symbol
	//      Children()[1] is LIST of parameters
	return loopHandler(el.Children()[0], ctx)
}

