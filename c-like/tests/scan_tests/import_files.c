
#include "tests/unit/unit_test.h"
#include "tests/unit/utils.h"
#include "compiler/errors.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"
#include "compiler/compiler_file.h"

#include <string.h>

bool CompilerFileListContains(const List *list, const char *path)
{
	for (ListEntry *entry = list->first; entry != NULL; entry = entry->next)
	{
		CompilerFile *cf = entry->item;
		if (strcmp(cf->path->buffer, path) == 0)
			return true;
	}
	return false;
}

int ListLength(const List *list)
{
	int n = 0;
	for (ListEntry *entry = list->first; entry != NULL; entry = entry->next)
		n ++;
	return n;
}

void ImportFiles(void)
{

	CompileState compile_state;
	CompileStateInit(&compile_state);

	const char *env = "tests";

	const char *argv[] =
	{
		"moss-cc",
		"import_test.base"
	};
	int argc = sizeof(argv)/sizeof(const char *);

	const CompilerArgs *args = ParseArgs(argc, argv);
	CHECK(args);
	CHECK(PassConfigure(&compile_state, args, env));
	CHECK(PassSearchAndParse(&compile_state));

	Namespace *group_ns = NamespaceGetChild(&compile_state.root_namespace, TempString("import_test"));
	CHECK(group_ns != NULL);
	Namespace *ns;
	ns = NamespaceGetChild(group_ns, TempString("base"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/base/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/base.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 0);
	CHECK(ListLength(&ns->public_syms.imports) == 3);
	CHECK(ListLength(&ns->private_syms.imports) == 3);

	ns = NamespaceGetChild(group_ns, TempString("lib_a"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/lib_a/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/lib_a.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 0);
	CHECK(ListLength(&ns->public_syms.imports) == 1);
	CHECK(ListLength(&ns->private_syms.imports) == 1);

	ns = NamespaceGetChild(group_ns, TempString("lib_b"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/lib_b/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/lib_b.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 0);
	CHECK(ListLength(&ns->public_syms.imports) == 1);
	CHECK(ListLength(&ns->private_syms.imports) == 1);

	ns = NamespaceGetChild(group_ns, TempString("lib_c"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/lib_c/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/lib_c.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 0);
	CHECK(ListLength(&ns->public_syms.imports) == 1);
	CHECK(ListLength(&ns->private_syms.imports) == 1);

	ns = NamespaceGetChild(group_ns, TempString("lib_d"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/lib_d/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/lib_d.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->private_syms.files, "tests/source/import_test/lib_d.moss"));
	CHECK(ListLength(&ns->public_syms.imports) == 0);
	CHECK(ListLength(&ns->private_syms.imports) == 1);

	ns = NamespaceGetChild(group_ns, TempString("lib_e"));
	CHECK(ns != NULL);
	CHECK(strcmp("import_test/lib_e/", ns->path->buffer) == 0);
	CHECK(ListLength(&ns->public_syms.files) == 1);
	CHECK(CompilerFileListContains(&ns->public_syms.files, "tests/import/import_test/lib_e.moss"));
	CHECK(ListLength(&ns->private_syms.files) == 2);
	CHECK(CompilerFileListContains(&ns->private_syms.files, "tests/source/import_test/lib_e.part1.moss"));
	CHECK(CompilerFileListContains(&ns->private_syms.files, "tests/source/import_test/lib_e.part2.moss"));
	CHECK(ListLength(&ns->public_syms.imports) == 0);
	CHECK(ListLength(&ns->private_syms.imports) == 0);

	CHECK(ErrorCount() == 0);
	CompileStateFree(&compile_state);
	FreeArgs(args);
}

