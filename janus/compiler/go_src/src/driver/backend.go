
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

	ret, err := exec.Command(cmd, llvmName, "-o", asmName).CombinedOutput()
	if err != nil {
		output.Error("LLVM code generation error: %v", string(ret))
	}
}

func runAssembleLink(basePath string, asmName string, name string) {

	cmd := "gcc"
	lib := filepath.Join(basePath, "library/clib/clib.o")
	ret, err := exec.Command(cmd, asmName, "-o", name, lib).CombinedOutput()
	if err != nil {
		output.Error("GCC asm/link error: %v", string(ret))
	}
}

