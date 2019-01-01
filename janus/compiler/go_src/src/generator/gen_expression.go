
package generator

import (
	"output"
	"parser"
	"symbols"
)

func genExpression(genFunc GeneratedFunction,
	ctx *symbols.EvalContext, el parser.ParseElement) Result {

	children := el.Children()
	opElement := children[0]

	if opElement.ElementType() == parser.OPERATOR {

		opName := opElement.TokenString()
		opChoices := ctx.LookupOperator(opName)
		if opChoices == nil {
			parser.Error(el.FilePos(), "no definition for operator %v", opName)
			return nil
		}

		opResult := NewFunctionChoiceResult(opChoices)
		argList := children[1:]
		return genInvokeFunction(genFunc, ctx, opResult, argList)
	}

	//FIXME implement  (what non-operator expressions are there???)
	output.FIXMEDebug("expression with non-operator %v", opElement)
	return nil
}

