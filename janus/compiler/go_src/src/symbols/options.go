
package symbols

import (
	"fmt"
	"parser"
)

type HeaderOptions struct {
	ByName map[string] *HeaderOption

	Version string
	ModuleName []string
	ExportSymbols bool
	ObjectMode bool
	MachineMode bool
}

type HeaderOption struct {
	Name string
	ParseTree parser.ParseElement
	Value DataValue
	DotName []string
	Recognized bool
}

func InterpretHeaderOptions(file *SourceFile) {
	if len(file.ParseTree.Children()) == 0 {
		return
	}

	header := file.ParseTree.Children()[0]
	if header.ElementType() != parser.HEADER {
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
			Recognized: false,
		}
	}

	val := file.Options.ByName["module_name"]
	if val != nil {
		val.Recognized = true
		if val.DotName != nil {
			file.Options.ModuleName = val.DotName
		} else {
			pos := val.ParseTree.FilePos()
			parser.Error(pos, "invalid value for module_name")
		}
	}

	// FIXME add these
	//ExportSymbols bool
	//ObjectMode bool
	//MachineMode bool

	for k, v := range file.Options.ByName {
		if !v.Recognized {
			pos := v.ParseTree.FilePos()
			parser.Warning(pos, "unrecognized option: %v", k)
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

