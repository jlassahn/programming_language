
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
	"log"
	"output"
	"parser"
	"symbols"
)

const (
	MODE_COMPILE = 0
	MODE_LIB = 1
	MODE_SYMBOL = 2
	MODE_PARSE = 3
	MODE_TOKEN = 4
)

type parameters struct {
	Files []string
	Mode int
}

func parseArgs() *parameters {
	ret := &parameters {
		nil,
		MODE_COMPILE }

	for _, arg := range os.Args[1:] {

		if arg[0] == '-' {
			switch arg[1:] {
				case "lib":
					ret.Mode = MODE_LIB
				case "tokens":
					ret.Mode = MODE_TOKEN
				case "parse":
					ret.Mode = MODE_PARSE
				case "symbols":
					ret.Mode = MODE_SYMBOL

				default:
					log.Fatal("unknown option: "+arg) //FIXME better handling
			}

		} else {
			ret.Files = append(ret.Files, arg)
		}
	}

	if ret.Files == nil {
		log.Fatal("no source files specified")
	}

	return ret
}

func main() {
	args := parseArgs()
	if args == nil {
		os.Exit(1)
	}

	file_set := symbols.NewFileSet()

	if args.Mode == MODE_TOKEN {
		output.EnableTokens()
	}

	for _, file := range args.Files {
		file_set.AddByFileName(file)
	}

	if args.Mode == MODE_PARSE {
		for _, file := range file_set.FileList {
			parser.EmitParseTree(file.ParseTree)
		}
	}

	if args.Mode >= MODE_PARSE {
		return
	}

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

