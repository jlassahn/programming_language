
#include "tests/unit/unit_test.h"
#include "compiler/errors.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"
#include <string.h>

void SimpleVariable(void)
{

	CompileState compile_state;
	CompileStateInit(&compile_state);

	const char *env = "tests";

	const char *argv[] =
	{
		"moss-cc",
		"local.variable"
	};
	int argc = sizeof(argv)/sizeof(const char *);

	const CompilerArgs *args = ParseArgs(argc, argv);
	CHECK(args);
	CHECK(PassConfigure(&compile_state, args, env));
	CHECK(PassSearchAndParse(&compile_state));

	String string; // FIXME consider a StringSet function?
	string.data = "local";
	string.length = strlen(string.data);
	Namespace *ns = NamespaceGetChild(&compile_state.root_namespace, &string);
	CHECK(ns != NULL);

	string.data = "variable";
	string.length = strlen(string.data);
	ns = NamespaceGetChild(ns, &string);
	CHECK(ns != NULL);

	string.data = "var_int32";
	string.length = strlen(string.data);
	Symbol *sym = NamespaceFindSymbol(ns, &string);
	CHECK(sym != NULL);
	CHECK(StringEqualsCString(&sym->name, "var_int32"));

	CHECK(ErrorCount() == 0);
	CompileStateFree(&compile_state);
	FreeArgs(args);
}

