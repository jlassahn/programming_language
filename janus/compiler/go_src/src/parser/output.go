
package parser

import (
	"fmt"
)


func EmitElementTree(el ParseElement, depth int, cmt bool) {

	line, col := el.Position()
	fmt.Printf("(%3d, %2d) ", line, col)

	for i:=0; i<depth; i++ {
		fmt.Print("\t")
	}

	if cmt {
		fmt.Print("* ")
	}

	fmt.Printf("[%s] %s\n",
		el.ElementType(),
		el.TokenString())

	for _, child := range el.Comments() {
		EmitElementTree(child, depth+1, true)
	}
	for _, child := range el.Children() {
		EmitElementTree(child, depth+1, cmt)
	}
}

func EmitParseTree(el ParseElement) {
	EmitElementTree(el, 0, false)
}

