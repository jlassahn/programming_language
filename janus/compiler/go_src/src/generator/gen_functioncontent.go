
package generator

import (
	"fmt"

	"parser"
	"symbols"
)


func genFunctionContent(genFunc GeneratedFunction,
	ctxBase *symbols.EvalContext, el parser.ParseElement) Result {

	
	symbolTable := symbols.NewSymbolTable(
		fmt.Sprintf("local@%d", el.FilePos().Line),
		ctxBase.Symbols)

	ctx := &symbols.EvalContext {
		Symbols: symbolTable,
	}

	for _,elem := range el.Children() {
		loopHandler(genFunc, ctx, elem)
	}

	return nil
}

