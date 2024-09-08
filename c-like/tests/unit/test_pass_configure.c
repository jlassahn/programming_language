
#include "tests/unit/unit_test.h"
#include "tests/unit/utils.h"
#include "compiler/passes.h"
#include "compiler/compile_state.h"
#include "compiler/commandargs.h"
#include "compiler/compiler_file.h"
#include "compiler/namespace.h"
#include <string.h>

/*
PassConfigure:
	sets up basedirs from envronment, -I args
	set up input files and input modules
		decides which is a file and which a module, anything that's
		syntactically a module is treated as one.

	files become:
		List input_files; // List of CompilerFile*
	the input_files only have path value set.

	modules become:
		List input_modules; // List of Namespace*
	the modules are installed in the namespace tree
	no files or symbols yet

	No file IO needed, checking of whether the files actually exist in
	later pass PassSearchAndParse.

	returns false for fatal errors, true otherwise.
*/

static void CheckStringEntry(ListEntry **entry, const char *match)
{
	CHECK(*entry != NULL);
	if (*entry == NULL)
		return;

	StringBuffer *sb;
	sb = (*entry)->item;
	CHECK(strcmp(sb->buffer, match) == 0);

	*entry = (*entry)->next;
}

static void CheckFileEntry(ListEntry **entry, const char *match)
{
	CHECK(*entry != NULL);
	if (*entry == NULL)
		return;

	CompilerFile *cf = (*entry)->item;
	CHECK(strcmp(cf->path->buffer, match) == 0);
	CHECK(cf->flags == FILE_FROM_INPUT);
	CHECK(cf->root == NULL);
	CHECK(cf->namespace == NULL);
	CHECK(cf->parser_file.data == NULL);

	*entry = (*entry)->next;
}

static void CheckModuleEntry(ListEntry **entry,
		const char *path, const char *stem)
{
	CHECK(*entry != NULL);
	if (*entry == NULL)
		return;

	Namespace *ns = (*entry)->item;
	CHECK(strcmp(ns->path->buffer, path) == 0);

	*entry = (*entry)->next;
}

static void CheckBuiltins(CompileState *cs)
{
	Map *builtins = &cs->builtins;
	CHECK(MapFind(builtins, TempString("void")) != NULL);
}

void TestPassConfigure(void)
{
	CompileState compile_state;

	const char *argv[] =
	{
		"moss-cc",
		"-I", "path1/",
		"-I", "path2",
		"-I", "path3/",
		"./file.moss",
		"system.clib",
	};
	int argc = sizeof(argv)/sizeof(const char *);

	const char *env = "env1;env2";

	CompileStateInit(&compile_state);
	const CompilerArgs *args = ParseArgs(argc, argv);
	CHECK(args != NULL);

	CHECK(PassConfigure(&compile_state, args, env));

	ListEntry *entry;
	entry = compile_state.basedirs.first;
	CheckStringEntry(&entry, "./");
	CheckStringEntry(&entry, "env1/");
	CheckStringEntry(&entry, "env2/");
	CheckStringEntry(&entry, "path1/");
	CheckStringEntry(&entry, "path2/");
	CheckStringEntry(&entry, "path3/");
	CHECK(entry == NULL);

	entry = compile_state.input_files.first;
	CheckFileEntry(&entry, "./file.moss");
	CHECK(entry == NULL);

	entry = compile_state.input_modules.first;
	CheckModuleEntry(&entry, "system/clib/", "clib");
	CHECK(entry == NULL);

	CheckBuiltins(&compile_state);

	FreeArgs(args);
	CompileStateFree(&compile_state);
}

