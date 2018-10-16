
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

		file.Options.ByName[name] = &HeaderOption {
			ParseTree: opt.Children()[1],
			Name: name }

		value := EvaluateConstExpression(opt.Children()[1], PredefinedSymbols)
		fmt.Println(value)
	}

	for k, v := range file.Options.ByName {
		fmt.Println(k)
		fmt.Println(v)
	}
}

