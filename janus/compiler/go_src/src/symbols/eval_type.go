
package symbols

import (
	"parser"
)


func evalType(el parser.ParseElement, ctx *EvalContext) DataValue {

	//FIXME is this always right?
	return loopHandler(el.Children()[0], ctx)
}

