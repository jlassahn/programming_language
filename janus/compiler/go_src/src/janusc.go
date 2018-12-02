
package main

import (
	"os"
	"strings"

	"driver"
)


func main() {

	args := os.Args[1:]

	envList := os.Environ()
	env := map[string]string {}
	for _,x := range envList {
		toks := strings.SplitN(x, "=", 2)
		env[toks[0]] = toks[1]
	}

	ret := driver.Compile(args, env)

	os.Exit(ret)
}

