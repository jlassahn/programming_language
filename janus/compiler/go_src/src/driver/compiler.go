
/* The main program for executing the compiler
*/

package driver

import (
	"os"
	"strings"
	"path/filepath"

	"output"
	"parser"
	"symbols"
	"generator"
)


type parameters struct {
	Files []string
	Flags map[string]bool
	StringOpts map[string]string
	ListOpts map[string] []string
}

func Compile(basePath string, argsIn []string, envIn map[string]string) int {

	output.ErrorCount = 0

	args := parseArgs(argsIn)
	if args.Flags["help"] {
		EmitHelp()
		return 1
	}

	sourcePaths := args.ListOpts["source"]
	sourcePaths = append(sourcePaths,
		filepath.SplitList(envIn["JANUS_SOURCE_PATH"])...)
	sourcePaths = append(sourcePaths, ".", "./source",
		filepath.Join(basePath, "library", "source"))
	sourcePaths = normalize(sourcePaths)

	interfacePaths := args.ListOpts["interface"]
	interfacePaths = append(interfacePaths,
		filepath.SplitList(envIn["JANUS_INTERFACE_PATH"])...)
	interfacePaths = append(interfacePaths, ".", "./interfaces",
		filepath.Join(basePath, "library", "interfaces"))
	interfacePaths = normalize(interfacePaths)

	if args.Flags["show-paths"] {
		EmitPaths(sourcePaths, interfacePaths)
	}

	if output.ErrorCount > 0 {
		return 1
	}

	symbols.InitializeTypes()
	symbols.PredefinedSymbols() //initialize symbols

	if args.Flags["show-predefined"] {
		symbols.PredefinedSymbols().Emit(true)
	}

	fileSet := symbols.NewFileSet()

	if args.Flags["show-tokens"] {
		parser.EnableTokens()
	}

	for _, file := range args.Files {
		fileSet.AddByFileName(file)
	}

	if args.Flags["show-parse"] {
		for _, file := range fileSet.FileList {
			parser.EmitParseTree(file.ParseTree)
		}
	}

	if output.ErrorCount > 0 {
		return 1
	}

	if args.Flags["parse-only"] {
		return 0
	}

	for _, file := range fileSet.FileList {
		symbols.InterpretHeaderOptions(file)
	}

	if args.Flags["show-header"] {
		for _, file := range fileSet.FileList {
			file.Options.Emit()
		}
	}

	for _, file := range fileSet.FileList {
		fileSet.AddFileToModules(file)
	}


	symbols.ResolveImports(fileSet, interfacePaths, sourcePaths,
		args.Flags["show-imports"])

	if args.Flags["show-modules"] {
		fileSet.EmitModuleTree()
	}

	if output.ErrorCount > 0 {
		return 1
	}

	symbols.ResolveGlobals(fileSet)

	if args.Flags["show-globals"] {

		EmitGlobals(fileSet)
	}

	if output.ErrorCount > 0 {
		return 1
	}

	if args.Flags["imports-only"] {
		return 0
	}

	name := args.StringOpts["name"]
	if name == "" {
		name = filepath.Base(args.Files[0])
		name = strings.Split(name, ".")[0]
	}
	llvmName := name + ".ll"
	asmName := name + ".s"

	generateCode(llvmName, fileSet)
	if output.ErrorCount > 0 {
		return 1
	}

	if args.Flags["llvm-only"] {
		return 0
	}

	runLLVM(basePath, llvmName, asmName)
	if output.ErrorCount > 0 {
		return 1
	}

	if args.Flags["asm-only"] {
		return 0
	}

	runAssembleLink(basePath, asmName, name)
	if output.ErrorCount > 0 {
		return 1
	}

	return 0
}

func generateCode(llvmName string, fileSet *symbols.FileSet) {

	fp, err := os.Create(llvmName)
	if err != nil {
		output.Error("can't create output file %v: %v", llvmName, err)
		return
	}
	outfile := output.NewObjectFile(fp)
	generator.GenerateCode(fileSet, outfile)
	fp.Close()
}

func normalize(paths []string) []string {

	var ret []string

	dups := map[string]bool { }

	for _,x := range paths {
		x = filepath.Clean(x)
		full, err := filepath.Abs(x)
		if err == nil {
			if dups[full] {
				continue
			}
			dups[full] = true
		}
		ret = append(ret, x)
	}

	return ret
}


func EmitPaths(sourcePaths []string, interfacePaths []string) {
	output.Emit("Source Search Paths:")
	for _,x := range sourcePaths {
		output.Emit("  %v", x)
	}
	output.Emit("")
	output.Emit("Interface Search Paths:")
	for _,x := range interfacePaths {
		output.Emit("  %v", x)
	}
	output.Emit("")
}

func EmitGlobals(fileSet *symbols.FileSet) {

	output.Emit("MODULE SYMBOLS")

	fileSet.EmitModuleSymbols()

	output.Emit("FILE SYMBOLS")

	for _,file := range fileSet.FileList {
		file.EmitGlobals()
	}

}

func EmitHelp() { output.Emit(
`janusc [options] [files]

compile a Janus program.  Unless the -lib option is used the output will be
an executable program named after the first source file.

Options:
 -help        : print this help text
 -lib         : output a .jlib library file

 -parse-only  : stop after parsing the source files
 -imports-only: stop after resolving imported symbols
 -llvm-only   : stop after generating LLVM
 -asm-only    : stop after generating assembly

 -show-paths  : print the search paths used by the compiler
 -show-tokens : print the tokenized files to stdout
 -show-parse  : print a parse tree to stdout
 -show-header : print the janus header options to stdout
 -show-modules: print included source files and module names
 -show-imports: print files found through import statements
 -show-globals: print file scope symbol tables
 -show-predefined: print predefined built-in symbol tables

 -name *      : set the name of the output file
 -source *    : add directory to source search path
 -interface * : add directour to interface search path

  FIXME cross-reference generator, showing imports
  FIXME dump symbol tables
`) }

func parseArgs(args []string) *parameters {

	var files []string
	var flags = map[string]bool {
		"help": false,
		"lib": false,
		"show-paths": false,
		"show-tokens": false,
		"show-parse": false,
		"show-header": false,
		"show-modules": false,
		"show-imports": false,
		"show-globals": false,
		"show-predefined": false,
		"parse-only": false,
		"imports-only": false,
		"llvm-only": false,
		"asm-only": false,
	}

	var stringOpts = map[string]string {
		"name": "",
	}

	var listOpts  = map[string] []string {
		"interface": nil,
		"source": nil,
	}

	for i:=0; i<len(args); i++ {
		arg := args[i]

		if arg[0] != '-' {
			files = append(files, arg)
			continue
		}

		opt := arg[1:]
		_, ok := flags[opt]
		if ok {
			flags[opt] = true
			continue
		}

		_, ok = stringOpts[opt]
		if ok {
			i ++
			if len(args) > i {
				val := args[i]
				if stringOpts[opt] == ""{
					stringOpts[opt] = val
				} else {
					output.Error("option %v can only be used once", arg)
				}
				continue
			} else {
				output.Error("missing parameter: %v", arg)
			}
		}

		_, ok = listOpts[opt]
		if ok {
			i ++
			if len(args) > i {
				val := args[i]
				listOpts[opt] = append(listOpts[opt], val)
				continue
			} else {
				output.Error("missing parameter: %v", arg)
			}
		}

		output.Error("unknown option: %v", arg)
	}

	if files == nil && !flags["help"] {
		output.Error("no source files specified")
	}

	ret := &parameters {
		Files: files,
		Flags: flags,
		StringOpts: stringOpts,
		ListOpts: listOpts,
	}

	return ret
}

