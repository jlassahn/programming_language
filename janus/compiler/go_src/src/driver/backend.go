
package driver

import (
	"os/exec"

	"output"
)

//FIXME don't use current directory, use relative to janusc

func runLLVM(llvmName string, asmName string ) {

	err := exec.Command("./llc", llvmName, "-o", asmName).Run()
	if err != nil {
		output.Error("LLVM code generation error: %v", err)
	}
}

func runAssembleLink(asmName string, name string) {
	output.FIXMEDebug("IMPLEMENT GCC PASS")
}

