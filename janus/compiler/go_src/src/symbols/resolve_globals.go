
package symbols

import (
	"fmt"

	"output"
	"parser"
)


//FIXME cleaner organization???
//FIXME do we actually want to keep declarations in the final symbols?

type uninitializedSymbol struct {
	name string
	declarations []parser.ParseElement
	isConst bool
	elementType *parser.Tag

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

		name := ""
		isConst := false

		switch el.ElementType() {
		case parser.DEF:
			name = el.Children()[0].TokenString()
			isConst = false

		case parser.CONST:
			name = el.Children()[0].TokenString()
			isConst = true

		//FIXME implement
		//struct
		//interface
		//method (can be handled in a second pass unless we allow const methods)
		//alias
		//operator

		default:
			continue
		}

		//FIXME getSymbol only called here, maybe inline it.
		sym := getSymbol(name, file, mod)
		if sym == nil {
			parser.Error(el.FilePos(),
				"symbol %v collides with import name", name)
			continue
		}

		if sym.elementType == nil {
			sym.elementType = el.ElementType()
			sym.isConst = isConst
		} else {
			if sym.elementType != el.ElementType() {
				parser.Error(el.FilePos(),
					"name declared with different types: %v", name)
				continue
			}
		}

		sym.declarations = append(sym.declarations, el)
		output.FIXMEDebug("    def for %v %p", name, sym)
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

	//FIXME Symbols needs to be sorted!
	for _,value := range mod.ExportedSymbols.Symbols {
		resolveConstValue(value)
	}
	// FIXME mod.LocalSymbols.Symbols
	// FIXME operators?
}

func resolveConstValue(value Symbol) Symbol {

	uninit, ok := value.(*uninitializedSymbol)
	if !ok {
		output.FIXMEDebug("  value %v already finalized", value)
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

	if !uninit.isConst {
		output.FIXMEDebug("  value %v not constant", value)
		return nil
	}

	if uninit.initialized != nil {
		output.FIXMEDebug("  value %v already resolved", value)
		return uninit.initialized
	}

	output.FIXMEDebug("  resolving constant %v", value)

	var dtypeMatch DataType
	var dvalMatch DataValue
	var funcSymbol *functionChoiceSymbol

	for _,el := range uninit.declarations {

		dtype, err := resolveSymbolType(uninit, el, resolveConstValue)
		if err != nil {
			output.FIXMEDebug("  resolving type failed for  %v: %v", value, err)
			return nil
		}

		//dtype can be nil when the type is inferred from the initializer
		output.FIXMEDebug("  resolving type for %v to %v", value, dtype)

		//FIXME what happens when dval points to a function body?
		dval := resolveSymbolValue(uninit, dtype, resolveConstValue)
		output.FIXMEDebug("  resolving value for %v to %v", value, dval)

		if dtype == nil && dval != nil {
			dtype = dval.Type()
			output.FIXMEDebug("  inferring type for %v to %v", value, dtype)
		}

		if dtype == nil {
			parser.Error(el.FilePos(), "no data type for %v", value)
			return nil
		}

		if dtype.Base() == FUNCTION_TYPE {
			if dtypeMatch != nil {
				parser.Error(el.FilePos(),
					"data type mismatch for %v", value)
				return nil
			}

			if funcSymbol == nil {
				funcSymbol = &functionChoiceSymbol {
					name: uninit.Name(),
					choices: []Symbol { },
				}
			}

			//FIXME merge functions with the same types
			sym := &baseSymbol {
				name: uninit.Name(),
				dtype: dtype,
				initialValue: dval,
				isConst: uninit.isConst,
				genVal: nil,
			}
			funcSymbol.choices = append(funcSymbol.choices, sym)

		} else {

			if funcSymbol != nil {
				parser.Error(el.FilePos(),
					"data type mismatch for %v", value)
				return nil
			}

			if dtypeMatch == nil {
				dtypeMatch = dtype
			} else if !TypeMatches(dtypeMatch, dtype) {
				parser.Error(el.FilePos(),
					"multiple declaration with different types for %v", value)
				return nil
			}

			if dvalMatch == nil {
				dvalMatch = dval
			} else {
				parser.Error(el.FilePos(),
					"multiple definitions for %v", value)
				return nil
			}
		}
	}


	var sym Symbol
	if funcSymbol != nil {
		sym = funcSymbol
	} else {
		if dtypeMatch == nil {
			output.Error("symbol with no type definition: %v", uninit.Name())
			return nil
		}

		sym = &baseSymbol {
			name: uninit.Name(),
			dtype: dtypeMatch,
			initialValue: dvalMatch,
			isConst: uninit.isConst,
			genVal: nil,
		}
	}

	uninit.initialized = sym
	return sym
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

	elType := uninit.declarations[0].ElementType()
	for _,dec := range uninit.declarations {
		if dec.ElementType() != elType {
			parser.Error(dec.FilePos(),
				"declaration mismatch with %v",
				uninit.declarations[0].FilePos())
			return nil
		}
	}

	//FIXME use resolveSymbolType here

	var sym Symbol
	switch elType {

	case parser.DEF:
		for _,el := range uninit.declarations {
			file := el.FilePos().File.(*SourceFile)
			output.FIXMEDebug("initializing %v %v", file.FileName, el)
			isConst := false //FIXME do we need to track this
			typeTree := el.Children()[1]
			valTree := el.Children()[2]

			output.FIXMEDebug("  symbol %v %v %v", isConst, typeTree, valTree)

			//FIXME include unitialized resolver and current symbol in ctx
			ctx := &EvalContext {
				Symbols: file.FileSymbols,
				SymbolPreprocessor: resolveVariableType,
				CycleDetectSymbol: uninit,
			}

			symTypeVal := EvaluateConstExpression(typeTree, ctx)
			if symTypeVal == nil {
				parser.Error(typeTree.FilePos(), "unknown data type in expr")
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
			uninit.declarations[0].FilePos(),
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
			//FIXME is error here always redundant?
			output.Error("symbol %v remains uninitialized", key)
			continue
		}
		replace[key] = uninit.initialized
	}

	for key,value := range replace {
		syms[key] = value
	}
}

func resolveSymbolType(sym *uninitializedSymbol, el parser.ParseElement,
	handler func(Symbol)Symbol) (DataType, error) {

	var typeTree parser.ParseElement

	switch sym.elementType {
	case parser.DEF:
		typeTree = el.Children()[1]

	case parser.CONST:
		typeTree = el.Children()[1]

	//FIXME and others...

	default:
		output.FatalError("bad resolveSymbolType for %v", sym.elementType)
		return nil, fmt.Errorf("unimplemented")
	}

	ctx := &EvalContext {
		Symbols: typeTree.FilePos().File.(*SourceFile).FileSymbols,
		SymbolPreprocessor: handler,
		CycleDetectSymbol: sym,
	}

	if typeTree.ElementType() == parser.EMPTY {
		return nil, nil
	}

	dval := EvaluateConstExpression(typeTree, ctx)
	if dval == nil {
		//FIXME error should probaly come from deeper..
		parser.Error(typeTree.FilePos(), "unknown data type in expr")
		return nil, fmt.Errorf("undefined data type")
	}
	symTypeVal, ok := dval.(TypeDataValue)
	if !ok {
		parser.Error(typeTree.FilePos(), "not a data type")
		return nil, fmt.Errorf("not a data type")
	}
	symType := symTypeVal.AsDataType()

	output.FIXMEDebug("  symType %v", symType)
	return symType, nil
}


func resolveSymbolValue(sym *uninitializedSymbol, initDT DataType, handler func(Symbol)Symbol) DataValue {

	for _,el := range sym.declarations {

		var valTree parser.ParseElement

		switch sym.elementType {
		case parser.DEF:
			valTree = el.Children()[2]

		case parser.CONST:
			valTree = el.Children()[2]

		//FIXME and others...

		default:
			output.FatalError("bad resolveSymbolValue for %v", sym.elementType)
			return nil
		}

		if valTree.ElementType() == parser.EMPTY {
			continue
		}

		ctx := &EvalContext {
			Symbols: valTree.FilePos().File.(*SourceFile).FileSymbols,
			SymbolPreprocessor: handler,
			CycleDetectSymbol: sym,
			InitializerType: initDT,
		}

		return EvaluateConstExpression(valTree, ctx)
	}

	return nil
}

