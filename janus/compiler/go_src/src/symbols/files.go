
package symbols

import (
	"os"
	"fmt"
	"lexer"
	"parser"
)


type SourceFile struct {
	Name string
	ModulePath []string
	ParseTree parser.ParseElement
}

func NewSourceFile() *SourceFile {
	return &SourceFile {
		"",
		nil,
		nil }
}


type FileSet struct {
	FileList []*SourceFile
}

/* FIXME when a file is specified on the command line,
   its import name defaults to the file name without
   extension or path, if it's .janus.  If it's .jsrc
   it has no default import name */

func (fs *FileSet) AddByFileName(name string) {
	/* FIXME check for duplicates */

	ret := NewSourceFile()
	fp, err := os.Open(name)
	if err == nil {
		lex := lexer.MakeLexer(fp)
		ret.ParseTree = parser.NewParser(lex).GetElement()
		fp.Close()
	} else {
		fmt.Println(err)
	}
	ret.Name = name
	InterpretHeaderOptions(ret)
	fs.FileList = append(fs.FileList, ret)
}

func (fs *FileSet) AddByImportName(name []string) {
	/* FIXME implement */
}

func NewFileSet() *FileSet {
	return &FileSet {
		nil }
}

func printTokens(lex *lexer.Lexer) {
	for {
		tok := lex.NextToken()
		fmt.Println(tok)

		//FIXME halt on first error, or try to recover?
		if tok.TokenType == lexer.ERROR {
			break
		}
		if tok.TokenType == lexer.EOF {
			break
		}
	}
}

func ResolveImports(file_set *FileSet) {

	// not a for range loop, because files are added to the
	// list as the loop progresses
	i := 0

	for i < len(file_set.FileList) {
		file := file_set.FileList[i]
		fmt.Println(file.Name)
		for _, el := range file.ParseTree.Children() {
			if el.ElementType() != lexer.IMPORT {
				continue
			}

			fmt.Println(lexer.TypeNames[el.ElementType()])
		}

		i ++
	}

}

