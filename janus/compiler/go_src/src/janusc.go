
/* janusc: the Janus language compiler
*/

package main

import (
	"os"
	"fmt"
	"path/filepath"
	"parser"
	"symbols"
)


func PrintHelp() { fmt.Print(
`janusc [options] [files]

compile a Janus program.  Unless the -lib option is used the output will be
an executable program named after the first source file.

Options:
 -help        : print this help text
 -lib         : output a .jlib library file

 -parse-only  : stop after parsing the source files

 -show-paths  : print the search paths used by the compiler
 -show-tokens : print the tokenized files to stdout
 -show-parse  : print a parse tree to stdout
 -show-header : print the janus header options to stdout
 -show-modules: print included source files and module names
 -show-imports: print files found through import statements

 -name *      : set the name of the output file
 -source *    : add directory to source search path
 -interface * : add directour to interface search path

  FIXME cross-reference generator, showing imports
  FIXME dump symbol tables
`) }

type parameters struct {
	Files []string
	Flags map[string]bool
	StringOpts map[string]string
	ListOpts map[string] []string
}

func parseArgs() *parameters {

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
		"parse-only": false,
	}

	var stringOpts = map[string]string {
		"name": "",
	}

	var listOpts  = map[string] []string {
		"interface": nil,
		"source": nil,
	}

	args := os.Args[1:]

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
					parser.CurrentLogger.FatalError(
						"option %v can only be used once", arg)
				}
				continue
			} else {
				parser.CurrentLogger.FatalError("missing parameter: %v", arg)
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
				parser.CurrentLogger.FatalError("missing parameter: %v", arg)
			}
		}

		parser.CurrentLogger.FatalError("unknown option: %v", arg)
	}

	if files == nil && !flags["help"] {
		parser.CurrentLogger.FatalError("no source files specified")
	}

	ret := &parameters {
		Files: files,
		Flags: flags,
		StringOpts: stringOpts,
		ListOpts: listOpts,
	}

	return ret
}

func main() {

	args := parseArgs()
	if args.Flags["help"] {
		PrintHelp()
		os.Exit(1)
	}

	sourcePaths := args.ListOpts["source"]
	sourcePaths = append(sourcePaths,
		filepath.SplitList(os.Getenv("JANUS_SOURCE_PATH"))...)
	sourcePaths = append(sourcePaths, ".", "./source")

	interfacePaths := args.ListOpts["interface"]
	interfacePaths = append(interfacePaths,
		filepath.SplitList(os.Getenv("JANUS_INTERFACE_PATH"))...)
	interfacePaths = append(interfacePaths, ".", "./interfaces")

	if args.Flags["show-paths"] {
		//FIXME clean up
		fmt.Println(sourcePaths)
		fmt.Println(interfacePaths)
	}

	file_set := symbols.NewFileSet()

	if args.Flags["show-tokens"] {
		parser.EnableTokens()
	}

	for _, file := range args.Files {
		file_set.AddByFileName(file)
	}

	if args.Flags["show-parse"] {
		for _, file := range file_set.FileList {
			parser.EmitParseTree(file.ParseTree)
		}
	}

	if args.Flags["parse-only"] {
		return
	}

	for _, file := range file_set.FileList {
		symbols.InterpretHeaderOptions(file)
	}

	if args.Flags["show-header"] {
		for _, file := range file_set.FileList {
			file.Options.Emit()
		}
	}

	for _, file := range file_set.FileList {
		file_set.AddFileToModules(file)
	}


	symbols.ResolveImports(file_set, interfacePaths, sourcePaths,
		args.Flags["show-imports"])

	if args.Flags["show-modules"] {
		file_set.EmitModuleTree()
	}

	symbols.ResolveGlobals(file_set)

	//FIXME resolve global symbols
	//  pass 1, create symbol table entries for all files without types
	//  pass 2, resolve types
	//  FIXME lots of weird corner cases here, type can be
	//        an explicit type (which might be defined in another file)
	//        the default type of an initializer (which might be a const expr)
	//        a constant of type CType?
	//  so pass 2 is probably a graph walk which visits expressions as needed
	//  and detects cycles.
}


/* FIXME sequence for compile

Load all files:
	* load files from the command line
	* recursively for all imports
		* only load files that haven't already been loaded
		* load interface file
		* load library file
		* load source file

Generate global symbol tables for all files
	* can't resolve all types since imports haven't been mapped...
Map imoprts into symbol tables
*/

