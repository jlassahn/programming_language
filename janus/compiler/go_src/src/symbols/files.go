
package symbols

import (
	"os"
	"fmt"
	"lexer"
	"parser"
	"output"
)


type SourceFile struct {
	FileName string
	ModulePath []string
	ParseTree parser.ParseElement
	Options HeaderOptions
}

func NewSourceFile() *SourceFile {
	return &SourceFile { }
}


type Module struct {
	Children map[string]*Module

	FileList []*SourceFile

	ExportedSymbols SymbolTable
	LocalSymbols SymbolTable
}

type FileSet struct {
	FileList []*SourceFile

	RootModule *Module
}

/* FIXME when a file is specified on the command line,
   its import name defaults to the file name without
   extension or path. */

func (fs *FileSet) AddByFileName(name string) {

	ret := NewSourceFile()

	fp, err := os.Open(name)
	if err != nil {
		output.FatalError(0,0, "unable to open file "+name)
	}

	lex := lexer.MakeLexer(fp)
	ret.ParseTree = parser.NewParser(lex).GetElement()
	fp.Close()


	ret.FileName = name
	fs.FileList = append(fs.FileList, ret)

	InterpretHeaderOptions(ret)

	//FIXME add to module tree
}

func (fs *FileSet) AddByImportName(name []string) {
	/* FIXME implement */
}

func NewFileSet() *FileSet {
	return &FileSet {
		FileList: nil,
		RootModule: nil }
}

func ResolveImports(file_set *FileSet) {

	// not a for range loop, because files are added to the
	// list as the loop progresses
	i := 0

	for i < len(file_set.FileList) {
		file := file_set.FileList[i]
		fmt.Println(file.FileName)
		for _, el := range file.ParseTree.Children() {
			if el.ElementType() != lexer.IMPORT {
				continue
			}

			fmt.Println("---")
			// FIXME implement
			// fmt.Println(lexer.TypeNames[el.ElementType()])
			parser.EmitParseTree(el)
			args := el.Children()

			modname := DotListAsStrings(args[0])
			impname := modname
			if len(args) == 2 {
				modname = DotListAsStrings(args[1])
			}

			fmt.Println(modname)
			fmt.Println(impname)
		}

		i ++
	}

}

