
package driver

import (
	"os"
	"os/exec"
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
	"expressions",
	"structs",
	"statements",
}

var definitionOKTests = []string {
	"consts",
	"const_exp",
}

var codegenOKTests = []string {
	"simple",
	"main",
	"factorial",
}

var executionOKTests = []string {
	"factorial",
	"while_loop",
}

var basePath = ".."

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

	runInfoCompare(t, name, args, infoExp)
}

func TestDefinitions(t *testing.T) {

	for _,name := range definitionOKTests {
		t.Run(name, func(t *testing.T) { runDefinitionTest(t, name) })
	}
}

func runDefinitionTest(t *testing.T, name string) {

	args := []string{
		"-imports-only",
		"-show-globals",
		name+".janus",
	}

	infoExp,_ := ioutil.ReadFile(name+".globals")

	runInfoCompare(t, name, args, infoExp)
}

func TestCodeGen(t *testing.T) {

	for _,name := range codegenOKTests {
		t.Run(name, func(t *testing.T) { runCodeGenTest(t, name) })
	}
}

func runCodeGenTest(t *testing.T, name string) {

	args := []string{
		name+".janus",
		"-asm-only",
	}

	genExp,_ := ioutil.ReadFile(name+".generated")

	logger.errors.Reset()
	logger.warns.Reset()
	logger.info.Reset()

	ret := Compile(basePath, args, nil)

	if ret != 0 {
		t.Errorf("compiler returned failure")
	}

	genAct,_ := ioutil.ReadFile(name+".ll")

	if !bytes.Equal(genAct, genExp) {
		t.Errorf("output doesn't match")
	}

	if logger.errors.Len() != 0 {
		t.Errorf("compiler generated errors")
	}

	if logger.warns.Len() != 0 {
		t.Errorf("compiler generated warnings")
	}

	/* FIXME do we care?
	if logger.info.Len() != 0 {
		t.Errorf("compiler generated info messages")
	}
	*/

	fp,_ := os.Create(name+".out")
	fp.Write(logger.errors.Bytes())
	fp.Write(logger.warns.Bytes())
	fp.Write(logger.info.Bytes())
	fp.Close()
}

func TestExecution(t *testing.T) {

	for _,name := range executionOKTests {
		t.Run(name, func(t *testing.T) { runExecutionTest(t, name) })
	}
}

func runExecutionTest(t *testing.T, name string) {

	args := []string{
		name+".janus",
		"-name", name+".exe",
	}

	runExp,_ := ioutil.ReadFile(name+".run")

	logger.errors.Reset()
	logger.warns.Reset()
	logger.info.Reset()

	ret := Compile(basePath, args, nil)

	if ret != 0 {
		t.Errorf("compiler returned failure")
	}

	if logger.errors.Len() != 0 {
		t.Errorf("compiler generated errors")
	}

	if logger.warns.Len() != 0 {
		t.Errorf("compiler generated warnings")
	}

	/* FIXME do we care?
	if logger.info.Len() != 0 {
		t.Errorf("compiler generated info messages")
	}
	*/

	runAct, err := exec.Command("./"+name+".exe").CombinedOutput()
	if err != nil {
		t.Errorf("execution error: %v", err)
	}

	if !bytes.Equal(runAct, runExp) {
		t.Errorf("output doesn't match")
	}

	fp,_ := os.Create(name+".buildout")
	fp.Write(logger.errors.Bytes())
	fp.Write(logger.warns.Bytes())
	fp.Write(logger.info.Bytes())
	fp.Close()

	fp,_ = os.Create(name+".runout")
	fp.Write(runAct)
	fp.Close()
}

func runInfoCompare(t *testing.T, name string, args []string, infoExp []byte) {

	logger.errors.Reset()
	logger.warns.Reset()
	logger.info.Reset()

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

