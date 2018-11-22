
/* janusc: the Janus language compiler

options:
	janusc [files]  compile to an executable program, named from the first file
	janusc -lib output.jlib [files]  FIXME what interfaces are public?
	janusc -tokens [file]  output the token list from parsing a file
	janusc -parse [file] output the parse tree from a file
	FIXME cross-reference generator, showing imports
	FIXME dump symbol tables
*/

package main

import (
	"os"
	"fmt"
	"output"
	"parser"
	"symbols"
)


func PrintHelp() {
	fmt.Println("FIXME help")
}

type parameters struct {
	Files []string
	Flags map[string]bool
}

func parseArgs() *parameters {

	var files []string
	var flags = map[string]bool {
		"help": false,
		"lib": false,
		"show-tokens": false,
		"show-parse": false,
		"show-header": false,
		"parse-only": false,
	}

	for _, arg := range os.Args[1:] {

		if arg[0] == '-' {
			flag := arg[1:]
			_, ok := flags[flag]
			if !ok {
				output.FatalNoFile("unknown option: "+arg)
				flags["help"] = true
			}
			flags[flag] = true
		} else {
			files = append(files, arg)
		}
	}

	if files == nil && !flags["help"] {
		output.FatalNoFile("no source files specified")
	}

	ret := &parameters {
		files,
		flags,
	}

	return ret
}

func main() {
	args := parseArgs()
	if args.Flags["help"] {
		PrintHelp()
		os.Exit(1)
	}

	file_set := symbols.NewFileSet()

	if args.Flags["show-tokens"] {
		output.EnableTokens()
	}

	for _, file := range args.Files {
		file_set.AddByFileName(file)
	}

	if args.Flags["show-parse"] {
		for _, file := range file_set.FileList {
			parser.EmitParseTree(file.ParseTree)
		}
	}

	if args.Flags["show-header"] {
		for _, file := range file_set.FileList {
			file.Options.Emit()
		}
	}

	if args.Flags["parse-only"] {
		return
	}

	fmt.Println("FIXME continuing after parse")

	symbols.ResolveImports(file_set)

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

