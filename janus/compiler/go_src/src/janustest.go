
package main

import (
	"fmt"
	"os/exec"
)


func main() {

	runParseTests()

	if false {
	out, err := exec.Command("./janusc", "tests/hello.janus").CombinedOutput()
	fmt.Println(string(out))
	fmt.Println(err)
	}
}

var valid_parse_tests = []string {
	"./tests/hello",
}

func runParseTests() {
	for _, file := range(valid_parse_tests) {
		out, err := exec.Command("./janusc", "-parse", file).CombinedOutput()
		if err != nil {
			fmt.Println("FIXME FAILED")
		} else {
			fmt.Println(string(out))
		}
		
		
	}

}

