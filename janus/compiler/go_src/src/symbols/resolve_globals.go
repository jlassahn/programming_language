
package symbols

import (
	"fmt"

	"parser"
)


type uninitializedSymbol struct {
	name string
	isConst bool
	typeParseTrees []parser.ParseElement
	valueParseTrees []parser.ParseElement //functions can have multiple
	needs *uninitializedSymbol
	initialized Symbol
}

func (self *uninitializedSymbol) Name() string { return self.name; }
func (self *uninitializedSymbol) Type() DataType { return nil; }
func (self *uninitializedSymbol) InitialValue() DataValue { return nil; }
func (self *uninitializedSymbol) IsConst() bool { return self.isConst; }

func (self *uninitializedSymbol) String() string {
	return self.name + ":uninitialized"
}

func ResolveGlobals(fileSet *FileSet) {
	//FIXME implement

	//fill module symbol tables with uninitializedSymbol
	findSymbolsForModule(fileSet.RootModule)

	//copy symbols from module symbol tables to file tables based on imports
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

func getSymbol(name string, file *SourceFile,
	mod *Module) *uninitializedSymbol {

	//FIXME  create or find in locals, move to ExportedSymbols if needed
	//sym := mod.LocalSymbols.Symbols[name].(*uninitializedSymbol)

	return nil
}

