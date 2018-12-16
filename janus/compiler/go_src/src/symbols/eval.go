
package symbols

import (
	"output"
	"parser"
)


type ConstEvaluator func(el parser.ParseElement, ctx *EvalContext) DataValue


type EvalContext struct {
	Symbols SymbolTable
	SymbolPreprocessor func(Symbol)Symbol
	CycleDetectSymbol *uninitializedSymbol
	InitializerType DataType
}

func (self *EvalContext) Lookup(name string) Symbol {

	ret := self.Symbols.Lookup(name)
	if self.SymbolPreprocessor != nil {
		uninit, ok := ret.(*uninitializedSymbol)
		if ok {
			output.FIXMEDebug("lookup bouncing to preproc for %v", ret)
			if self.CycleDetectSymbol != nil {
				self.CycleDetectSymbol.needs = uninit
			}
			ret = self.SymbolPreprocessor(ret)
			if self.CycleDetectSymbol != nil {
				self.CycleDetectSymbol.needs = nil
			}
			output.FIXMEDebug("lookup got from preproc %v", ret)
		}
	}

	return ret
}

func (self *EvalContext) LookupOperator(name string) FunctionChoiceSymbol {
	//FIXME fake
	return self.Symbols.LookupOperator(name)
}


var handlers = map[*parser.Tag] ConstEvaluator {
	parser.NUMBER: evalNumber,
	parser.SYMBOL: evalSymbol,
	parser.DOT_LIST: evalDotList,
	parser.EXPRESSION: evalExpression,
	parser.FUNCTION_TYPE: evalFunctionType,
	parser.TYPE: evalType,
}


// indirect call to EvalConstExpression, to avoid a circular dependency
var loopHandler ConstEvaluator

func EvaluateConstExpression(
	el parser.ParseElement, ctx *EvalContext) DataValue {

	loopHandler = EvaluateConstExpression

	//eval := evaluators[*el.ElementType()]
	eval := handlers[el.ElementType()]
	if eval == nil {
		//FIXME implement
		output.FIXMEDebug("no evaluator for %v\n", el.ElementType())
		parser.EmitParseTree(el)
		ctx.Symbols.Emit()
		return nil
	} else {
		return eval(el, ctx)
	}
}

