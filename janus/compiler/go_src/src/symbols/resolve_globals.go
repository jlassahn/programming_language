
package symbols

import (
	"fmt"

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

	//FIXME implement

	//evaluate types and initial values, storing in intiialized subsymbol
	//    symbol table Lookup knows about uninitialized symbols
	//    follows link to initialized if it exists
	//    propagates evaluation if not
	//replace all uninitialized Symbols with initialized in all symbol tables

}

func findSymbolsForModule(mod *Module) {

	fmt.Printf("searching module: %v\n", mod.Name)
	for _, file := range mod.FileList {
		fmt.Printf("  searching file: %v\n", file.FileName)
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
			fmt.Printf("    def for %v %p\n", name, sym)

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

	fmt.Printf("FIXME resolve imports for file %v\n", file.FileName)

	for _,imp := range file.Imports {

		mod := imp.ModuleData

		table := file.FileSymbols

		// operators always import into the main file table

		fmt.Printf("FIXME importing operators from %v\n", table.Name)
		for key,value := range mod.ExportedSymbols.Operators {
			fmt.Printf("  importing %v %v\n", key, value)
			if table.Operators[key] == nil {
				table.Operators[key] = value
			} else {
				parser.CurrentLogger.Error("operator import collision %v %v",
					ToDotString(imp.ImportName), key)
			}
		}

		// find the table for the module namespace
		for _,name := range imp.ImportName {
			sym := table.Symbols[name]
			table = sym.InitialValue().(NamespaceDataValue).AsSymbolTable()
		}

		// import symbols into the module namespace
		fmt.Printf("FIXME importing symbols from %v\n", table.Name)
		for key,value := range mod.ExportedSymbols.Symbols {
			fmt.Printf("  importing %v %v\n", key, value)
			if table.Symbols[key] == nil {
				table.Symbols[key] = value
			} else {
				parser.CurrentLogger.Error("import collision %v.%v",
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

