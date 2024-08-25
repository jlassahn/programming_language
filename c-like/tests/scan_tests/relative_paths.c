
#include "tests/unit/unit_test.h"
#include "compiler/errors.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"

void RelativePaths(void)
{

	CompileState compile_state;
	CompileStateInit(&compile_state);

	const char *env = "tests/scan_tests/relative_files";

	const char *argv[] =
	{
		"moss-cc",
		"local.group.part1"
	};
	int argc = sizeof(argv)/sizeof(const char *);

	const CompilerArgs *args = ParseArgs(argc, argv);
	CHECK(args);
	CHECK(PassConfigure(&compile_state, args, env));
	CHECK(PassSearchAndParse(&compile_state));

	CHECK(ErrorCount() == 0);
	CompileStateFree(&compile_state);
	FreeArgs(args);
}

