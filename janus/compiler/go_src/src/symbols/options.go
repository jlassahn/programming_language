
package symbols

import (
	"fmt"
	"lexer"
	"parser"
)

type HeaderOptions struct {
	ByName map[string] *HeaderOption
	Version string
}

type HeaderOption struct {
	Name string
	ParseTree parser.ParseElement
	Value DataValue
	DotName []string
}

func InterpretHeaderOptions(file *SourceFile) {
	if len(file.ParseTree.Children()) == 0 {
		return
	}

	header := file.ParseTree.Children()[0]
	if header.ElementType() != lexer.HEADER {
		return
	}

	file.Options.Version = header.Children()[0].TokenString()
	file.Options.ByName = make(map[string] *HeaderOption)

	if len(header.Children()) < 2 {
		return
	}
	options := header.Children()[1]
	for _, opt := range options.Children() {

		keys := DotListAsStrings(opt.Children()[0])
		name := ""
		for _, x := range(keys) {
			if name == "" {
				name = x
			} else {
				name = name + "." + x
			}
		}

		//values are either constant expressions or
		// previously undefined dot lists used as e.g. module names
		//
		ctx := &EvalContext {
			Symbols: PredefinedSymbols,
		}
		value := EvaluateConstExpression(opt.Children()[1], ctx)
		var dotval []string
		if value == nil {
			dotval = DotListAsStrings(opt.Children()[1])
		}

		file.Options.ByName[name] = &HeaderOption {
			ParseTree: opt.Children()[1],
			Name: name,
			Value: value,
			DotName: dotval,
		}
	}

	for k, v := range file.Options.ByName {
		if v.Value == nil {
			fmt.Printf("option [%v] = [%v]\n", k, v.DotName)
		} else {
			fmt.Printf("option [%v] = [%v]\n", k, v.Value)
		}
	}
}

