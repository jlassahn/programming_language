
#include "tests/unit/unit_test.h"
#include "compiler/errors.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"

void TestFile(const char *name)
{
	CompileState compile_state;
	CompileStateInit(&compile_state);

	const char *env = "tests/scan_tests/";

	const char *argv[] =
	{
		"moss-cc",
		name
	};
	int argc = sizeof(argv)/sizeof(const char *);

	const CompilerArgs *args = ParseArgs(argc, argv);
	CHECK(args);
	CHECK(PassConfigure(&compile_state, args, env));
	CHECK(!PassSearchAndParse(&compile_state));

	CHECK(ErrorCount() > 0);
	CompileStateFree(&compile_state);
	FreeArgs(args);
}

void BadFilenames(void)
{
	TestFile("tests/scan_tests/bad_files/123.moss");
	TestFile("tests/scan_tests/bad_files/_.moss");
}

