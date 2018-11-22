
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
			Symbols: PredefinedSymbols(),
		}


		dotval := getUndefinedDotlist(opt.Children()[1], ctx)

		var value DataValue
		if dotval == nil {
			value = EvaluateConstExpression(opt.Children()[1], ctx)
		}

		file.Options.ByName[name] = &HeaderOption {
			ParseTree: opt.Children()[1],
			Name: name,
			Value: value,
			DotName: dotval,
		}
	}
}

func (self *HeaderOptions) Emit() {

	for k, v := range self.ByName {
		if v.DotName == nil {
			fmt.Printf("option [%v] = [%v]\n", k, v.Value)
		} else {
			fmt.Printf("option [%v] = [%v]\n", k, v.DotName)
		}
	}
}

func getUndefinedDotlist(el parser.ParseElement, ctx *EvalContext) []string {
		dotval := DotListAsStrings(el)
		if dotval == nil {
			return nil
		}
		if ctx.Symbols.Lookup(dotval[0]) == nil {
			return dotval
		}
		return nil
}

