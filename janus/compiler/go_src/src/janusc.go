
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
	"parser"
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
			processFile(fp, args)
			fp.Close()
		} else {
			fmt.Println(err)
		}
	}
}

func processFile(fp io.Reader, args *parameters ) {
	lex := lexer.MakeLexer(fp)
	parser := parser.NewParser(lex)

	if args.Mode == MODE_TOKEN {
		printTokens(lex)
	} else {
		printParseTree(parser)
	}


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

func printParseTree(parse parser.Parser) {

	for {
		el := parse.GetElement()

		printElementTree(el, 0, false)

		/*
		tok := el.Token()

		if tok != nil {
			line, col := el.Position()
			fmt.Printf("type = %d Position = (%d, %d) tt=%d txt=%s\n",
				el.ElementType(),
				line, col,
				tok.TokenType,
				tok.Text)
		} else {
			//FIXME handle nonterminals
		}
		*/

		if el.ElementType() <= lexer.EOF {
			break
		}
	}

	/*
	Children() []ParseElement
	Comments() []ParseElement
	ElementType() int
	Position() (int, int)
	Token() *lexer.Token
	*/
}

func printElementTree(el parser.ParseElement, depth int, cmt bool) {

	line, col := el.Position()
	fmt.Printf("(%3d, %2d) ", line, col)

	for i:=0; i<depth; i++ {
		fmt.Print("\t")
	}

	if cmt {
		fmt.Print("* ")
	}

	fmt.Printf("type = %s txt=%s\n",
		lexer.TypeNames[el.ElementType()],
		el.TokenString())

	for _, child := range el.Comments() {
		printElementTree(child, depth+1, true)
	}
	for _, child := range el.Children() {
		printElementTree(child, depth+1, cmt)
	}
}

