
package driver

import (
	"os/exec"
	"path/filepath"

	"output"
)

func runLLVM(basePath string, llvmName string, asmName string ) {

	cmd := filepath.Join(basePath, "llc")
	if filepath.Base(cmd) == cmd {
		cmd = "./"+cmd
	}

	err := exec.Command(cmd, llvmName, "-o", asmName).Run()
	if err != nil {
		output.Error("LLVM code generation error: %v", err)
	}
}

func runAssembleLink(basePath string, asmName string, name string) {

	cmd := "gcc"
	err := exec.Command(cmd, asmName, "-o", name).Run()
	if err != nil {
		output.Error("GCC asm/link error: %v", err)
	}
}

