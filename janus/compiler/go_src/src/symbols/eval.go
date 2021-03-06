
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
			if self.CycleDetectSymbol != nil {
				self.CycleDetectSymbol.needs = uninit
			}
			ret = self.SymbolPreprocessor(ret)
			if self.CycleDetectSymbol != nil {
				self.CycleDetectSymbol.needs = nil
			}
		}
	}

	return ret
}

func (self *EvalContext) LookupOperator(name string) FunctionChoiceSymbol {

	//FIXME handle uninitializedSymbol and SymbolPreprocessor
	return self.Symbols.LookupOperator(name)
}

var handlers = map[*parser.Tag] ConstEvaluator {
	parser.NUMBER: evalNumber,
	parser.SYMBOL: evalSymbol,
	parser.DOT_LIST: evalDotList,
	parser.EXPRESSION: evalExpression,
	parser.FUNCTION_TYPE: evalFunctionType,
	parser.TYPE: evalType,
	parser.FUNCTION_CONTENT: evalFunctionContent,
	parser.LIST: evalList,
}


//FIXME rename and cleanup loopHandler!

// indirect call to EvalConstExpression, to avoid a circular dependency
var loopHandler ConstEvaluator
func init() {
	loopHandler = EvaluateConstExpression
}

func EvaluateConstExpression(
	el parser.ParseElement, ctx *EvalContext) DataValue {

	eval := handlers[el.ElementType()]
	if eval == nil {
		//FIXME implement
		output.Error("FIXME no evaluator for %v\n", el.ElementType())
		//parser.EmitParseTree(el)
		//ctx.Symbols.Emit()
		return nil
	} else {
		return eval(el, ctx)
	}
}

func EvaluateConstRHS(el parser.ParseElement, ctx *EvalContext) DataValue {

	initDT := ctx.InitializerType

	ret := loopHandler(el, ctx)
	if ret == nil {
		return nil
	}

	ret = MaskConstant(ret)
	if initDT != nil {
		ret = ConvertConstant(ret, initDT)
	}
	return ret
}

