
package symbols

import (
	"os"
	"strings"
	"path/filepath"

	"output"
	"parser"
)


type ImportLink struct {
	ModuleName []string
	ImportName []string
	ModuleData *Module
}

type SourceFile struct {
	FileName string
	ParseTree parser.ParseElement
	Options HeaderOptions
	Imports []*ImportLink
	FileSymbols *symbolTable
}

func NewSourceFile() *SourceFile {
	return &SourceFile { }
}

func (self *SourceFile) SetModuleByFileName() {

	filename := filepath.Base(self.FileName)
	parts := strings.Split(filename, ".")

	//don't display errors here, it's OK to have a weird filename as long
	// as header options override it.

	if len(parts) != 2 {
		return
	}
	base := parts[0]
	ext := parts[1]

	if !parser.IsValidIdentifier(base) {
		return
	}

	self.Options.ModuleName = []string{base}
	if ext == "janus" {
		self.Options.ExportSymbols = true
	}
}

func (self *SourceFile) EmitGlobals() {
	self.FileSymbols.Emit()
}

type Module struct {

	Name string

	Children map[string]*Module

	FileList []*SourceFile

	ExportedSymbols *symbolTable
	LocalSymbols *symbolTable
}

func NewModule(name string) *Module {
	return &Module {
		Name: name,
		Children: map[string]*Module {},
		FileList: nil,
		ExportedSymbols: NewSymbolTable(name+"(EXPORTED)", nil),
		LocalSymbols: NewSymbolTable(name+"(LOCAL)", nil),
	}
}

func (self *Module) EmitModuleTree(depth int) {
	for _, x := range self.FileList {

		output.EmitIndented(depth, "file: %v exported: %v",
			x.FileName, x.Options.ExportSymbols)
	}

	var keys []string
	for k := range self.Children {
		keys = append(keys, k)
	}

	for _,k := range keys {
		output.EmitIndented(depth, "%v:", k)

		self.Children[k].EmitModuleTree(depth + 1)
	}

}

func (self *Module) EmitModuleSymbols() {
	self.LocalSymbols.Emit()
	self.ExportedSymbols.Emit()

	for _,x := range self.Children {
		x.EmitModuleSymbols()
	}

}


type FileSet struct {
	FileList []*SourceFile

	RootModule *Module
}

func NewFileSet() *FileSet {
	return &FileSet {
		FileList: nil,
		RootModule: NewModule("@root") }
}


func (self *FileSet) LookupModule(path []string) *Module {

	mod := self.RootModule

	for _,name := range path {
		mod = mod.Children[name]
		if mod == nil {
			return nil
		}
	}
	return mod
}

/* FIXME when a file is specified on the command line,
   its import name defaults to the file name without
   extension or path. */

func (fs *FileSet) AddByFileName(name string) *SourceFile {

	ret := NewSourceFile()

	ret.FileSymbols = NewSymbolTable(name, PredefinedSymbols())

	fp, err := os.Open(name)
	if err != nil {
		output.Error("unable to open file %v", name)
		return nil
	}

	lex := parser.MakeLexer(fp, name)
	ret.ParseTree = parser.NewParser(lex).GetElement()
	fp.Close()


	ret.FileName = name
	fs.FileList = append(fs.FileList, ret)

	ret.SetModuleByFileName()  //might get overwritten by header options

	return ret
}

func (self *FileSet) AddFileToModules(file *SourceFile) {

	if len(file.Options.ModuleName) < 1 {
		output.Error("can't infer a module name for file %v", file.FileName)
		return
	}

	baseMod := self.RootModule
	path := file.Options.ModuleName

	for _,x := range path {
		mod := baseMod.Children[x]
		if mod == nil {
			mod = NewModule(x)
			baseMod.Children[x] = mod
		}

		baseMod = mod
	}

	baseMod.FileList = append(baseMod.FileList, file)
}

func (self *FileSet) EmitModuleTree() {
	self.RootModule.EmitModuleTree(0)
}

func (self *FileSet) EmitModuleSymbols() {
	for _,x := range self.RootModule.Children {
		x.EmitModuleSymbols()
	}
}

func ResolveImports(file_set *FileSet,
	interfacePaths []string,
	sourcePaths[] string,
	showImports bool) {

	// not a for range loop, because files are added to the
	// list as the loop progresses
	i := 0

	for i < len(file_set.FileList) {

		file := file_set.FileList[i]

		for _, el := range file.ParseTree.Children() {

			if el.ElementType() != parser.IMPORT {
				continue
			}

			pos := el.FilePos()

			args := el.Children()

			modname := DotListAsStrings(args[1])
			impname := modname
			if args[0].ElementType() != parser.EMPTY {
				impname = DotListAsStrings(args[0])
			}

			if showImports {
				output.Emit("file %v importing %v as %v",
					file.FileName, ToDotString(modname), ToDotString(impname))
			}

			mod := file_set.LookupModule(modname)
			if mod != nil {
				link := &ImportLink {
					ModuleName: modname,
					ImportName: impname,
					ModuleData: mod,
				}
				file.Imports = append(file.Imports, link)

				if showImports {
					output.Emit("   - %v already imported",
						ToDotString(modname))
				}
				continue
			}

			fileList := getSearchList(modname, interfacePaths, sourcePaths)
			for _,name := range fileList {

				fp, err := os.Open(name)
				if err != nil {
					continue
				}

				newFile := NewSourceFile()
				newFile.FileSymbols = NewSymbolTable(name, PredefinedSymbols())

				lex := parser.MakeLexer(fp, name)
				newFile.ParseTree = parser.NewParser(lex).GetElement()
				fp.Close()

				newFile.FileName = name
				newFile.SetModuleByFileName() //FIXME only set ExportSymbols
				newFile.Options.ModuleName = modname

				InterpretHeaderOptions(newFile)

				//FIXME more efficient compare
				if ToDotString(newFile.Options.ModuleName) !=
					ToDotString(modname) {
					output.Error(
						"module name doesn't match path for file %v", name)
				}
				file_set.FileList = append(file_set.FileList, newFile)
				file_set.AddFileToModules(newFile)

				if showImports {
					output.Emit("   - found file %v", name)
				}
			}

			mod = file_set.LookupModule(modname)
			if mod != nil {
				link := &ImportLink {
					ModuleName: modname,
					ImportName: impname,
					ModuleData: mod,
				}
				file.Imports = append(file.Imports, link)

				if showImports {
					output.Emit("   - module %v done", ToDotString(modname))
				}
				continue
			} else {
				parser.Error(pos, "no file found for import %v", ToDotString(modname))
			}
		}

		i ++
	}

}

func getSearchList(modname []string,
	interfacePaths []string,
	sourcePaths []string) [] string {

	fileset := map[string]bool {}
	var fileList []string

	relPath := filepath.Join(modname...)
	for _,x := range interfacePaths {
		path := filepath.Join(x, relPath)+".janus"
		if !fileset[path] {
			fileset[path] = true
			fileList = append(fileList, path)
		}
	}

	for _,x := range sourcePaths {
		path := filepath.Join(x, relPath)+".janus"
		if !fileset[path] {
			fileset[path] = true
			fileList = append(fileList, path)
		}
		path = filepath.Join(x, relPath)+".jsrc"
		if !fileset[path] {
			fileset[path] = true
			fileList = append(fileList, path)
		}
	}
	return fileList
}

//FIXME put somewhere else
func ToDotString(name []string) string {
	if len(name) == 0 {
		return ""
	}

	ret := name[0]
	for _,x := range name[1:] {
		ret = ret + "."+x
	}

	return ret
}

