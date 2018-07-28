
package symbols

import (
	"fmt"
	"lexer"
)

func InterpretHeaderOptions(file *SourceFile) {
	if len(file.ParseTree.Children()) == 0 {
		return
	}

	header := file.ParseTree.Children()[0]
	if header.ElementType() != lexer.HEADER {
		return
	}

	version := header.Children()[0].TokenString()
	fmt.Println(version)
	if len(header.Children()) < 2 {
		return
	}
	options := header.Children()[1]
	for _, opt := range options.Children() {
		name := DotListAsStrings(opt.Children()[0])
		value := EvaluateConstExpression(opt.Children()[1], PredefinedSymbols)
		fmt.Println(name)
		fmt.Println(value)
	}
}

