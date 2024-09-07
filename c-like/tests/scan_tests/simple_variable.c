
#include "tests/unit/unit_test.h"
#include "tests/unit/utils.h"
#include "compiler/errors.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"
#include "compiler/data_value.h"
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

	Namespace *ns = NamespaceGetChild(&compile_state.root_namespace, TempString("local"));
	CHECK(ns != NULL);
	ns = NamespaceGetChild(ns, TempString("variable"));
	CHECK(ns != NULL);

	Symbol *sym = NamespaceFindSymbol(ns, TempString("var_int32"));
	CHECK(sym != NULL);
	CHECK(StringEqualsCString(&sym->name, "var_int32"));

	Symbol *sym2 = NamespaceFindPrivateSymbol(ns, TempString("var_int32"));
	CHECK(sym2 != NULL);
	CHECK(StringEqualsCString(&sym2->name, "var_int32"));
	CHECK(sym->associated == sym2);
	CHECK(sym2->associated == sym);
	CHECK(sym2 != sym);

	SymbolTable *syms;
	syms = &ns->public_syms.symbol_table;
	CHECK(SymbolTableFind(syms, TempString("var_int32")) == sym);

	syms = &ns->private_syms.symbol_table;
	CHECK(SymbolTableFind(syms, TempString("var_int32")) == sym2);

	Symbol *dtype_sym = SymbolTableFind(syms, TempString("int32"));
	CHECK(dtype_sym != NULL);
	DataValue *dtype_dval = SymbolGetDValue(dtype_sym);
	CHECK(dtype_dval != NULL);
	CHECK(dtype_dval->value_type == VTYPE_DTYPE);

	CHECK(ErrorCount() == 0);
	CompileStateFree(&compile_state);
	FreeArgs(args);
}

