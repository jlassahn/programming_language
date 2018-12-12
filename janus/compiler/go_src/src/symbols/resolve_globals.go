
package symbols

import (
	"output"
	"parser"
)


type uninitDeclaration struct {
	parseTree parser.ParseElement
	file *SourceFile
}

type uninitializedSymbol struct {
	name string
	declarations []*uninitDeclaration
	needs *uninitializedSymbol
	initialized Symbol
}

func (self *uninitializedSymbol) Name() string { return self.name; }
func (self *uninitializedSymbol) Type() DataType { return nil; }
func (self *uninitializedSymbol) InitialValue() DataValue { return nil; }
func (self *uninitializedSymbol) IsConst() bool { return false }
func (self *uninitializedSymbol) SetGenVal(val interface{}) { }
func (self *uninitializedSymbol) GetGenVal() interface{} { return nil }

func (self *uninitializedSymbol) String() string {
	return self.name + ":uninitialized"
}

func ResolveGlobals(fileSet *FileSet) {

	//create namespaces for imports in file symbol tables
	//FIXME this could happen earlier, when the imports are first parsed

	for _,file := range fileSet.FileList {
		addImportNamespaces(file)
	}

	//fill module symbol tables with uninitializedSymbol
	//fill file symbol tables with symbols defined in that file
	findSymbolsForModule(fileSet.RootModule)

	//copy imported symbols from module symbol tables to file tables
	for _,file := range fileSet.FileList {
		resolveImportedSymbols(file, fileSet)
	}

	
	//evaluate types and initial values, storing in initialized subsymbol
	//    symbol table Lookup knows about uninitialized symbols
	//    follows link to initialized if it exists
	//    propagates evaluation if not

	resolveConstValues(fileSet.RootModule, fileSet)
	resolveVariableTypes(fileSet.RootModule, fileSet)
	replaceUninitialized(fileSet)
	resolveVariableValues(fileSet.RootModule, fileSet)

}

func findSymbolsForModule(mod *Module) {

	output.FIXMEDebug("searching module: %v", mod.Name)
	for _, file := range mod.FileList {
		output.FIXMEDebug("  searching file: %v", file.FileName)
		findSymbolsForFile(file, mod)
	}

	for _,child := range mod.Children {
		findSymbolsForModule(child)
	}
}

func findSymbolsForFile(file *SourceFile, mod *Module) {
	for _, el := range file.ParseTree.Children() {

		switch el.ElementType() {

		case parser.DEF:
			/*
				KEYWORD def or const
				SYMBOL name
				TYPE or FUNCTION_TYPE dtype
				EXPRESSION, etc initializer
			*/
			name := el.Children()[1].TokenString()
			sym := getSymbol(name, file, mod)
			if sym == nil {
				parser.Error(el.FilePos(),
					"symbol %v collides with import name", name)
				continue
			}

			dec := &uninitDeclaration {
				parseTree: el,
				file: file,
			}
			sym.declarations = append(sym.declarations, dec)
			output.FIXMEDebug("    def for %v %p", name, sym)

		//FIXME implement
		//struct
		//interface
		//method (can be handled in a second pass unless we allow const methods)
		//alias
		//operator

		default:
			continue
		}
	}
}

func resolveImportedSymbols(file *SourceFile, fileSet *FileSet) {

	output.FIXMEDebug("resolve imports for file %v", file.FileName)

	for _,imp := range file.Imports {

		mod := imp.ModuleData

		table := file.FileSymbols

		// operators always import into the main file table

		output.FIXMEDebug("importing operators from %v", table.Name)
		for key,value := range mod.ExportedSymbols.Operators {
			output.FIXMEDebug("  importing %v %v", key, value)
			if table.Operators[key] == nil {
				table.Operators[key] = value
			} else {
				output.Error("operator import collision %v %v",
					ToDotString(imp.ImportName), key)
			}
		}

		// find the table for the module namespace
		for _,name := range imp.ImportName {
			sym := table.Symbols[name]
			table = sym.InitialValue().(NamespaceDataValue).AsSymbolTable()
		}

		// import symbols into the module namespace
		output.FIXMEDebug("importing symbols from %v", table.Name)
		for key,value := range mod.ExportedSymbols.Symbols {
			output.FIXMEDebug("  importing %v %v", key, value)
			if table.Symbols[key] == nil {
				table.Symbols[key] = value
			} else {
				output.Error("import collision %v.%v",
					ToDotString(imp.ImportName), key)
			}
		}

	}
}

func addImportNamespaces(file *SourceFile) {

	for _,imp := range file.Imports {
		table := file.FileSymbols
		for _,name := range imp.ImportName {
			if table.Symbols[name] == nil {

				//FIXME better table name
				newTable := NewSymbolTable(name, nil)

				val := &namespaceDV {
					value: newTable,
				}

				table.Symbols[name] = &baseSymbol {
					name: name,
					dtype: NamespaceType,
					initialValue: val,
					isConst: true,
				}
			}

			sym := table.Symbols[name]
			table = sym.InitialValue().(NamespaceDataValue).AsSymbolTable()
		}
	}
}

func getSymbol(name string, file *SourceFile,
	mod *Module) *uninitializedSymbol {

	if mod.LocalSymbols.Symbols[name] == nil {
		sym := &uninitializedSymbol {
			name: name,
			declarations: nil,
			needs: nil,
			initialized: nil,
		}
		mod.LocalSymbols.Symbols[name] = sym
	}

	sym := mod.LocalSymbols.Symbols[name].(*uninitializedSymbol)
	if file.Options.ExportSymbols {
		mod.ExportedSymbols.Symbols[name] = sym
	}

	if file.FileSymbols.Symbols[name] == nil {
		file.FileSymbols.Symbols[name] = sym
	}

	if file.FileSymbols.Symbols[name] != sym {
		return nil
	}

	return sym
}

/* 
rules for global type and value resolution:
const values can only depend on other const values (not variables)
consts can't have circular dependencies, even the obscure cases
	where consts could resolve because each only requires partial
	info about the other are forbidden.
the type of a variable can only depend on consts
a variable with implicit type has its type based on const values
	and the types of other variables, but not the values of variables.
variables with implicit types can't have circular dependencies with
	each other.
the only non-const initializers allowed for data are
	references to data values
	references to struct members of data values
	maybe references to functions or methods?
variable initializers can form cyclic references, always safe
	because initializers only reference types, never other initializers.
types are a kind of const value (they're of data type CType)

ways in which normal function declarations are like const declarations:
	can't be written to, only initialized
	can be used as an initializer with implicit type
ways they're different:
	can't assign one to a const
	can make a reference to one

const function declarations can be _called_ from const initializers.
	const functions _can_ have references
*/

//FIXME struct content definitions in local tables shouldn't
//      leak into exported tables.
//      each file can see a different subset of the struct's
//      members, methods, size and extensions

func resolveConstValues(mod *Module, fileSet *FileSet) {

	output.FIXMEDebug("resolveConstValues for %v", mod.Name)
	for _,x := range mod.Children {
		resolveConstValues(x, fileSet)
	}

	for _,value := range mod.ExportedSymbols.Symbols {
		resolveConstValue(value)
	}
	// FIXME mod.LocalSymbols.Symbols
	// FIXME operators?
}

func resolveConstValue(value Symbol) Symbol {
	return nil
}

func resolveVariableTypes(mod *Module, fileSet *FileSet) {

	output.FIXMEDebug("resolveVariableTypes for %v", mod.Name)
	for _,x := range mod.Children {
		resolveVariableTypes(x, fileSet)
	}

	for _,value := range mod.ExportedSymbols.Symbols {
		resolveVariableType(value)
	}
	// FIXME mod.LocalSymbols.Symbols
	// FIXME operators?
}

//FIXME does this need a cycle resolver?  All consts should be finished
//      already, and types can't depend on variables
func resolveVariableType(value Symbol) Symbol {

	output.FIXMEDebug("resolveVariableType %v", value)
	uninit, ok := value.(*uninitializedSymbol)
	if !ok {
		output.FIXMEDebug("value %v already finalized", value)
		return value
	}

	if uninit.needs != nil {
		output.Error("definition loop")
		x := uninit.needs
		for {
			if x == nil { break }
			if x == uninit { break }
			output.Error("  definition %v -> %v", x.name, x.needs.name)
			x = x.needs
		}
		return nil
	}

	if uninit.initialized != nil {
		output.FIXMEDebug("value %v already resolved", value)
		return uninit.initialized
	}

	elType := uninit.declarations[0].parseTree.ElementType()
	for _,dec := range uninit.declarations {
		if dec.parseTree.ElementType() != elType {
			parser.Error(dec.parseTree.FilePos(),
				"declaration mismatch with %v",
				uninit.declarations[0].parseTree.FilePos())
			return nil
		}
	}

	var sym Symbol
	switch elType {

	case parser.DEF:
		for _,dec := range uninit.declarations {
			el := dec.parseTree
			file := dec.file
			output.FIXMEDebug("initializing %v %v", file.FileName, el)
			//FIXME is string compare the right thing here?
			isConst := (el.Children()[0].TokenString() == "const")
			typeTree := el.Children()[2]
			valTree := el.Children()[3]

			output.FIXMEDebug("  symbol %v %v %v", isConst, typeTree, valTree)

			//FIXME include unitialized resolver and current symbol in ctx
			ctx := &EvalContext {
				Symbols: file.FileSymbols,
				SymbolPreprocessor: resolveVariableType,
				CycleDetectSymbol: uninit,
			}

			symTypeVal := EvaluateConstExpression(typeTree, ctx)
			if symTypeVal == nil {
				parser.Error(typeTree.FilePos(), "unknown data type")
				return nil
			}
			//FIXME function types include the parameter names
			//      have to make the names consistent!
			//      have to use the names in the actual definition.
			symType := symTypeVal.(TypeDataValue).AsDataType()
			output.FIXMEDebug("  symType %v", symType)

			dval := &codeDV {
				dtype: CodeType,
				element: valTree,
				file: file,
			}

			if sym == nil {
				if symType.Base() == FUNCTION_TYPE {
					sym = &functionChoiceSymbol {
						name: uninit.Name(),
						choices: []Symbol {
							&baseSymbol {
								name: uninit.Name(),
								dtype: symType,
								initialValue: dval,
								isConst: isConst,
								genVal: nil,
							},
						},
					}
				} else {
					output.FIXMEDebug("NOT A FUNCTION")
				}

			} else {
				output.FatalError("FIXME merge multiple defs of symbol")
			}
		}

	//FIXME implement
	//struct
	//interface
	//method (can be handled in a second pass unless we allow const methods)
	//alias
	//operator

	default:
		parser.FatalError(
			uninit.declarations[0].parseTree.FilePos(),
			"Unhandled element: %v", elType)
	}

	uninit.initialized = sym
	return sym
}

func resolveVariableValues(mod *Module, fileSet *FileSet) {

	output.FIXMEDebug("resolveVariableValues for %v", mod.Name)
	for _,x := range mod.Children {
		resolveVariableValues(x, fileSet)
	}
	 //FIXME do something
}

func replaceUninitialized(fileSet *FileSet) {
	replaceUninitializedInModule(fileSet.RootModule)

	for _,file := range fileSet.FileList {
		replaceUninitializedInMap(file.FileSymbols.Symbols)
		// FIXME operators?
	}
}

func replaceUninitializedInModule(mod *Module) {

	output.FIXMEDebug("replaceUninitialized for %v", mod.Name)
	for _,x := range mod.Children {
		replaceUninitializedInModule(x)
	}
	replaceUninitializedInMap(mod.ExportedSymbols.Symbols)
	replaceUninitializedInMap(mod.LocalSymbols.Symbols)
	// FIXME operators?
}

func replaceUninitializedInMap(syms map[string]Symbol) {

	replace := map[string]Symbol {}

	for key,value := range syms {
		uninit, ok := value.(*uninitializedSymbol)
		if !ok {
			continue
		}
		if uninit.initialized == nil {
			output.FatalError("symbol %v remains uninitialized", key)
			continue
		}
		replace[key] = uninit.initialized
	}

	for key,value := range replace {
		syms[key] = value
	}
}

