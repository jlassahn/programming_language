
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

	
	//evaluate types and initial values, storing in intiialized subsymbol
	//    symbol table Lookup knows about uninitialized symbols
	//    follows link to initialized if it exists
	//    propagates evaluation if not

	resolveSymbolValues(fileSet.RootModule, fileSet)

	//FIXME implement
	//replace all uninitialized Symbols with initialized in all symbol tables

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

func resolveSymbolValues(mod *Module, fileSet *FileSet) {

	for _,x := range mod.Children {
		resolveSymbolValues(x, fileSet)
	}

	//FIXME struct content definitions in local tables shouldn't
	//      leak into exported tables.
	//      each file can see a different subset of the struct's
	//      members, methods, size and extensions
	//FIXME CType variables have similar problems.
	/* FIXME implement
	for key, value :=  mod.ExportedSymbols.Symbols {
	}
	*/

	/* rules for global type and value resolution:
	   const values can only depend on other const values (not variables)
	   consts can't have circular dependencies, even the obscure cases
	     where consts could resolve because each only requires partial
	     info about the other are forbidden.
	   the type of a variable can only depend on consts
	   a variable with implicit type can ony be initialized by
	     consts or function declarations
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

	//self.LocalSymbols.Emit()
	//self.ExportedSymbols.Emit()
}

