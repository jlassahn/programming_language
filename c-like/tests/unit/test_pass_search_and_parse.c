

#include "tests/unit/unit_test.h"
#include "compiler/pass_search_and_parse.h"
#include "compiler/compile_state.h"

void TestPassSearchAndParse(void)
{
	CompileState compile_state;
	CompileStateInit(&compile_state);

	CHECK(PassSearchAndParse(&compile_state));

	CompileStateFree(&compile_state);
}

