
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
	//FIXME implement

	//fill module symbol tables with uninitializedSymbol
	//fill file symbol tables with symbols defined in that file
	findSymbolsForModule(fileSet.RootModule)

	//FIXME
	fileSet.EmitModuleSymbols()

	//copy imported symbols from module symbol tables to file tables
	for _,file := range fileSet.FileList {
		resolveImportedSymbols(file, fileSet)
	}

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
	//FIXME implement
	fmt.Printf("FIXME resolve imports for file %v\n", file.FileName)
	//FIXME create namespaceDV for each import node
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

	file.FileSymbols.Symbols[name] = sym

	return sym
}

