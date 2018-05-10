
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
	"io"
	"os"
	"fmt"
	"log"
	"lexer"
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

	for _, file := range args.Files {
		fp, err := os.Open(file)
		if err == nil {
			print_tokens(fp)
			fp.Close()
		} else {
			fmt.Println(err)
		}
	}
}

func print_tokens(fp io.Reader) {
	lex := lexer.MakeLexer(fp)
	for {
		tok := lex.NextToken()
		fmt.Println(tok)

		if tok.TokenType == lexer.ERROR {
			break
		}
		if tok.TokenType == lexer.EOF {
			break
		}
	}
}

