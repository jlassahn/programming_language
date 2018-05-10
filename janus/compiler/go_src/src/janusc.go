
package main

import (
	"io"
	"os"
	"fmt"
	"lexer"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		fp, err := os.Open(os.Args[1])
		if err == nil {
			print_tokens(fp)
			fp.Close()
		} else {
			print_tokens(fp) // FIXME testing error path
			fmt.Println(err)
		}
	} else {
		src := strings.NewReader("janus 1.0;\n #comment")
		print_tokens(src)
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

