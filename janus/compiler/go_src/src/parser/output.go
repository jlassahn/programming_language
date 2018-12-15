
package parser

import (
	"fmt"

	"output"
)


func EmitElementTree(el ParseElement, depth int, cmt bool) {

	str := ""

	pos := el.FilePos()
	str = str + fmt.Sprintf("(%3d, %2d) ", pos.Line, pos.Column)

	for i:=0; i<depth; i++ {
		str = str + "\t"
	}

	if cmt {
		str = str + "* "
	}

	str = str + fmt.Sprintf("[%s] %s", el.ElementType(), el.TokenString())
	output.Emit(str)

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


var enable_tokens bool = false

func EnableTokens() {
	enable_tokens = true
}

func EmitToken(x string) {
	if enable_tokens {
		output.Emit(x)
	}
}



func posTag(pos *FilePosition) string {
	return fmt.Sprintf("%v:%v,%v: ", pos.File.Name(), pos.Line, pos.Column)
}

func FatalError(pos *FilePosition, msg string, args ...interface{}) {
	output.FatalError(posTag(pos) + msg, args...)
}

func Error(pos *FilePosition, msg string, args ...interface{}) {
	output.Error(posTag(pos) + msg, args...)
}

func Warning(pos *FilePosition, msg string, args ...interface{}) {
	output.Warning(posTag(pos) + msg, args...)
}

func Info(pos *FilePosition, msg string, args ...interface{}) {
	output.Info(posTag(pos) + msg, args...)
}

func Emit(msg string, args ...interface{}) {
	output.Emit(msg, args...)
}

