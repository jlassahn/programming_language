
package parser

import (
	"fmt"
)


func EmitElementTree(el ParseElement, depth int, cmt bool) {

	pos := el.FilePos()
	fmt.Printf("(%3d, %2d) ", pos.Line, pos.Column)

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

