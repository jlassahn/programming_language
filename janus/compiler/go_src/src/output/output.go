
package output

import (
	"os"
	"fmt"
)

var enable_tokens bool = false


func EnableTokens() {
	enable_tokens = true
}

func EmitError(x string) {
	os.Stderr.WriteString(x + "\n")
}

func EmitWarning(x string) {
	os.Stderr.WriteString(x + "\n")
}

func EmitToken(x string) {
	if enable_tokens {
		fmt.Println(x)
	}
}

type LLVMFile struct {
	//FIXME implement
}

