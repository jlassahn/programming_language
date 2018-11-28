
package output

import (
	"os"
	"fmt"
)

var enable_tokens bool = false


// FIXME report filename
// FIXME handle error in file open, etc

func FatalError(line int, col int, txt string) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR at (%d, %d) %s\n",
		line, col, txt)
	os.Exit(1)
}

func FatalNoFile(txt string) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR %s\n", txt)
	os.Exit(1)
}

func Error(line int, col int, txt string) {
	fmt.Fprintf(os.Stderr, "ERROR at (%d, %d) %s\n",
		line, col, txt)
}

func Warning(line int, col int, txt string) {
	fmt.Fprintf(os.Stderr, "WARNING at (%d, %d) %s\n",
		line, col, txt)
}


func EnableTokens() {
	enable_tokens = true
}

func EmitToken(x string) {
	if enable_tokens {
		fmt.Println(x)
	}
}

type LLVMFile struct {
	//FIXME implement
}

