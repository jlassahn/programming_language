
package output

import (
	"os"
	"fmt"
)

type ErrorLogger interface {
	FatalError(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warning(msg string, arg ...interface{})
	Info(msg string, arg ...interface{})
	Emit(msg string, arg ...interface{})
}


type LLVMFile struct {
	//FIXME implement
}

var CurrentLogger ErrorLogger = DefaultLogger{}
var ErrorCount int


type NilLogger struct {}
func (self NilLogger) FatalError(msg string, args ...interface{}) { }
func (self NilLogger) Error(msg string, args ...interface{}) { }
func (self NilLogger) Warning(msg string, arg ...interface{}) { }
func (self NilLogger) Info(msg string, arg ...interface{}) { }
func (self NilLogger) Emit(msg string, arg ...interface{}) { }

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

func FatalError(msg string, args ...interface{}) {
	ErrorCount ++
	CurrentLogger.FatalError(msg, args...)
}
func Error(msg string, args ...interface{}) {
	ErrorCount ++
	CurrentLogger.Error(msg, args...)
}

func Warning(msg string, args ...interface{}) {
	CurrentLogger.Warning(msg, args...)
}

func Info(msg string, args ...interface{}) {
	CurrentLogger.Info(msg, args...)
}

func Emit(msg string, args ...interface{}) {
	CurrentLogger.Emit(msg, args...)
}

func EmitIndented(depth int, msg string, args ...interface{}) {
	for i:=0; i<depth; i++ {
		msg = "\t"+msg
	}
	CurrentLogger.Emit(msg, args...)
}

//FIXME separate path for this?
func FIXMEDebug(msg string, args ...interface{}) {
	CurrentLogger.Emit(msg, args...)
}

