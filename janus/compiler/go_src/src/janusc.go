
package main

import (
	"os"
	"strings"
	"path/filepath"

	"driver"
)


func main() {

	//fmt.Println(runtime.GOOS) //darwin, freebsd, linux, windows, ...

	basePath := filepath.Dir(os.Args[0])

	args := os.Args[1:]

	envList := os.Environ()
	env := map[string]string {}
	for _,x := range envList {
		toks := strings.SplitN(x, "=", 2)
		env[toks[0]] = toks[1]
	}

	ret := driver.Compile(basePath, args, env)

	os.Exit(ret)
}

