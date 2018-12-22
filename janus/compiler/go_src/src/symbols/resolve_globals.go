
package symbols

import (
	"fmt"

	"output"
	"parser"
)


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
	resolveVariableValues(fileSet.RootModule, fileSet)
	replaceUninitialized(fileSet)

}

func findSymbolsForModule(mod *Module) {

	for _, file := range mod.FileList {
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
	}
}

func resolveImportedSymbols(file *SourceFile, fileSet *FileSet) {

	for _,imp := range file.Imports {

		mod := imp.ModuleData

		table := file.FileSymbols

		// operators always import into the main file table

		for key,value := range mod.ExportedSymbols.Operators {
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
		for key,value := range mod.ExportedSymbols.Symbols {
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

	for _,x := range mod.Children {
		resolveConstValues(x, fileSet)
	}

	//FIXME Symbols needs to be sorted!
	for _,value := range mod.ExportedSymbols.Symbols {
		resolveConstValue(value)
	}

	for _,value := range mod.ExportedSymbols.Operators {
		resolveConstValue(value)
	}

	for _,value := range mod.LocalSymbols.Symbols {
		resolveConstValue(value)
	}

	for _,value := range mod.LocalSymbols.Operators {
		resolveConstValue(value)
	}
}

func resolveConstValue(value Symbol) Symbol {

	uninit, ok := value.(*uninitializedSymbol)
	if !ok {
		return value
	}

	if uninit.needs != nil {
		output.Error("definition loop")
		x := uninit
		for {
			output.Error("  definition %v -> %v", x.name, x.needs.name)
			x = x.needs
			if x == nil { break }
			if x == uninit { break }
		}
		return nil
	}

	if !uninit.isConst {
		return nil
	}

	if uninit.initialized != nil {
		return uninit.initialized
	}

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

		dval := resolveSymbolValue(uninit, dtype, resolveConstValue)

		if dtype == nil && dval != nil {
			dtype = dval.Type()
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

			sym := &baseSymbol {
				name: uninit.Name(),
				dtype: dtype,
				initialValue: dval,
				isConst: uninit.isConst,
				genVal: nil,
			}
			if funcSymbol.Add(sym) != nil {
				parser.Error(el.FilePos(), "delcaration mismatch for %v", value)
				return nil
			}

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

func resolveVariableValues(mod *Module, fileSet *FileSet) {

	for _,x := range mod.Children {
		resolveVariableValues(x, fileSet)
	}

	for _,value := range mod.ExportedSymbols.Symbols {
		resolveVariableValue(value)
	}

	for _,value := range mod.ExportedSymbols.Operators {
		resolveVariableValue(value)
	}

	for _,value := range mod.LocalSymbols.Symbols {
		resolveVariableValue(value)
	}

	for _,value := range mod.LocalSymbols.Operators {
		resolveVariableValue(value)
	}
}

func resolveVariableValue(value Symbol) Symbol {

	uninit, ok := value.(*uninitializedSymbol)
	if !ok {
		return value
	}

	if uninit.initialized != nil {
		return uninit.initialized
	}

	//FIXME below is nearly identical to resolveConstValues
	var dtypeMatch DataType
	var dvalMatch DataValue
	var funcSymbol *functionChoiceSymbol

	for _,el := range uninit.declarations {

		dtype, err := resolveSymbolType(uninit, el, nil)
		if err != nil {
			return nil
		}

		//dtype can be nil when the type is inferred from the initializer

		dval := resolveSymbolValue(uninit, dtype, resolveVariableValue)

		if dtype == nil && dval != nil {
			dtype = dval.Type()
		}

		if dtype == nil {
			parser.Error(el.FilePos(), "no data type for %v", value)
			return nil
		}

		if dtype.Base() == FUNCTION_TYPE {
			if dtypeMatch != nil {
				parser.Error(el.FilePos(), "data type mismatch for %v", value)
				return nil
			}

			if funcSymbol == nil {
				funcSymbol = &functionChoiceSymbol {
					name: uninit.Name(),
					choices: []Symbol { },
				}
			}

			sym := &baseSymbol {
				name: uninit.Name(),
				dtype: dtype,
				initialValue: dval,
				isConst: uninit.isConst,
				genVal: nil,
			}
			if funcSymbol.Add(sym) != nil {
				parser.Error(el.FilePos(), "delcaration mismatch for %v", value)
				return nil
			}

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

func replaceUninitialized(fileSet *FileSet) {
	replaceUninitializedInModule(fileSet.RootModule)

	for _,file := range fileSet.FileList {
		replaceUninitializedInMap(file.FileSymbols.Symbols)
		// FIXME operators?
	}
}

func replaceUninitializedInModule(mod *Module) {

	for _,x := range mod.Children {
		replaceUninitializedInModule(x)
	}
	replaceUninitializedInMap(mod.ExportedSymbols.Symbols)
	replaceUninitializedInMap(mod.LocalSymbols.Symbols)
	//FIXME operators
	//replaceUninitializedInMap(mod.ExportedSymbols.Operators)
	//replaceUninitializedInMap(mod.LocalSymbols.Operators)
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

