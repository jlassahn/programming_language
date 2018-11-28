
package parser

import (
	"os"
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


var enable_tokens bool = false

func EnableTokens() {
	enable_tokens = true
}

func EmitToken(x string) {
	if enable_tokens {
		fmt.Println(x)
	}
}


type ErrorLogger interface {
	FatalError(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warning(msg string, arg ...interface{})
	Info(msg string, arg ...interface{})
	Emit(msg string, arg ...interface{})
}


var CurrentLogger ErrorLogger = DefaultLogger{}

type DefaultLogger struct { }


func (self DefaultLogger) FatalError(msg string, args ...interface{}) {
	fmt.Printf("FATAL " + msg + "\n", args...)
	os.Exit(1)
}

func (self DefaultLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("ERROR " + msg + "\n", args...)
}

func (self DefaultLogger) Warning(msg string, args ...interface{}) {
	fmt.Printf("WARNING " + msg + "\n", args...)
}

func (self DefaultLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("INFO " + msg + "\n", args...)
}

func (self DefaultLogger) Emit(msg string, args ...interface{}) {
	fmt.Printf(msg + "\n", args...)
}

func posTag(pos *FilePosition) string {
	return fmt.Sprintf("%v:%v,%v: ", pos.File, pos.Line, pos.Column)
}

func FatalError(pos *FilePosition, msg string, args ...interface{}) {
	CurrentLogger.FatalError(posTag(pos) + msg, args...)
}

func Error(pos *FilePosition, msg string, args ...interface{}) {
	CurrentLogger.Error(posTag(pos) + msg, args...)
}

func Warning(pos *FilePosition, msg string, args ...interface{}) {
	CurrentLogger.Warning(posTag(pos) + msg, args...)
}

func Info(pos *FilePosition, msg string, args ...interface{}) {
	CurrentLogger.Info(posTag(pos) + msg, args...)
}

func Emit(msg string, args ...interface{}) {
	CurrentLogger.Emit(msg, args...)
}

