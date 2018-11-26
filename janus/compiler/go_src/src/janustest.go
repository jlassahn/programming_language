
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)


var err_fp *os.File
var env []string

func main() {

	path := filepath.Dir(os.Args[0])
	path = filepath.Join(path, "tests")
	env = []string {
		"JANUS_SOURCE_PATH="+filepath.Join(path, "source"),
		"JANUS_INTERFACE_PATH="+filepath.Join(path, "interfaces"),
	}

	err_fp, _ = os.Create("janustest.log")
	runParseTests()
	err_fp.Close()
}

var error_count int
var test_name string
var test_file string

func StartTest(name string, file string) {
	error_count = 0
	test_name = name
	test_file = file
	fmt.Printf("%s %s ...", name, file)
	fmt.Fprintf(err_fp, "START %s %s\n", name, file)
}

func EndTest() {
	if error_count > 0 {
		fmt.Printf("FAILED (%d errors)\n", error_count)
		fmt.Fprintf(err_fp, "END   %s %s FAILED\n", test_name, test_file)
	} else {
		fmt.Printf("PASSED\n")
		fmt.Fprintf(err_fp, "END   %s %s PASSED\n", test_name, test_file)
	}
}

func Error(msg string) {
	fmt.Fprintln(err_fp, msg)
	error_count ++
}

func runParseTests() {
	for _, file := range(valid_parse_tests) {

		StartTest("valid parse", file)

		input_file := file + ".janus"
		compare_file := file + ".parse"

		cmd := exec.Command(
			"./janusc", "-parse-only", "-show-parse",
			input_file)
		cmd.Env = env
		out, err := cmd.CombinedOutput()

		{
			fp, _ := os.Create(file + ".out")
			fp.Write(out)
			fp.Close()
		}

		if err != nil {
			Error("Compiler returned failure")
		} else {
			cmp := make([]byte, len(out) + 1)

			fp, _ := os.Open(compare_file)
			n, _ := fp.Read(cmp)
			fp.Close()

			if n != len(out) {
				Error("output doesn't match")
			} else {
				for i :=0; i<len(out); i++ {
					if out[i] != cmp[i] {
						Error("output doesn't match")
						break
					}
				}
			}
		}

		EndTest()
	}

}

var valid_parse_tests = []string {
	"./tests/hello",
	"./tests/options",
}

