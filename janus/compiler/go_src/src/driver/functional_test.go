
package driver

import (
	"os"
	"fmt"
	"bytes"
	"io/ioutil"
	"testing"

	"output"
)

var parserOKTests = []string {
	"hello",
	"options",
	"comments",
	"literals",
}

func TestMain(m *testing.M) {

	logger = &Logger { }
	output.CurrentLogger = logger

	//fmt.Println(os.Getwd()) // cwd == compiler/go_src/src/driver
	//fmt.Println(os.Args)    // binary is in /var/{generated}

	os.Chdir("../../../tests")

	os.Exit(m.Run())
}


func TestParser(t *testing.T) {

	for _,name := range parserOKTests {
		t.Run(name, func(t *testing.T) { runParseTest(t, name) })
	}
}

func runParseTest(t *testing.T, name string) {

	args := []string{
		"-parse-only",
		"-show-parse",
		name+".janus",
	}

	infoExp,_ := ioutil.ReadFile(name+".parse")

	logger.errors.Reset()
	logger.warns.Reset()
	logger.info.Reset()

	basePath := "."
	ret := Compile(basePath, args, nil)

	if ret != 0 {
		t.Errorf("compiler returned failure")
	}

	if !bytes.Equal(logger.info.Bytes(), infoExp) {
		t.Errorf("output doesn't match")
	}

	if logger.errors.Len() != 0 {
		t.Errorf("compiler generated errors")
	}

	if logger.warns.Len() != 0 {
		t.Errorf("compiler generated warnings")
	}

	fp,_ := os.Create(name+".out")
	fp.Write(logger.errors.Bytes())
	fp.Write(logger.warns.Bytes())
	fp.Write(logger.info.Bytes())
	fp.Close()
}

var logger *Logger

type Logger struct {
	errors bytes.Buffer
	warns bytes.Buffer
	info bytes.Buffer
}

func (self *Logger) FatalError(msg string, args ...interface{}) {
	fmt.Fprintf(&self.errors, "FATAL " + msg + "\n", args...)
	panic("FATAL ERROR")
}

func (self *Logger) Error(msg string, args ...interface{}) {
	fmt.Fprintf(&self.errors, "ERROR " + msg + "\n", args...)
}

func (self *Logger) Warning(msg string, args ...interface{}) {
	fmt.Fprintf(&self.warns, "WARNING " + msg + "\n", args...)
}

func (self *Logger) Info(msg string, args ...interface{}) {
	fmt.Fprintf(&self.info, "INFO " + msg + "\n", args...)
}

func (self *Logger) Emit(msg string, args ...interface{}) {
	fmt.Fprintf(&self.info, msg + "\n", args...)
}

