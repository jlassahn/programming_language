

#include "tests/unit/unit_test.h"
#include "tests/unit/fake_directory.h"
#include "tests/unit/fake_errors.h"
#include "compiler/pass_search_and_parse.h"
#include "compiler/compile_state.h"
#include "compiler/compiler_file.h"
#include <string.h>

void TestSimpleInputFile(void)
{
	CompileState compile_state;

	CompileStateInit(&compile_state);

	CompilerFile *cf = CompilerFileCreate(
			StringBufferFromChars("./test.moss"));
	cf->flags = FILE_FROM_INPUT;

	ListInsertFirst(&compile_state.input_files, cf);

	FakeFileSet("./test.moss", "abcdef");

	CHECK(PassSearchAndParse(&compile_state));

	String s;
	s.data = "test";
	s.length = 4;

	Namespace *ns = MapFind(&compile_state.root_namespace.children, &s);
	CHECK(ns != NULL);
	CHECK(ns->private_files.first != NULL);
	CHECK(ns->private_files.first->item == cf);
	CHECK(ns->public_files.first == NULL);
	CHECK(ns->flags == NAMESPACE_HAS_INFILE);
	CHECK(cf->namespace == ns);
	CHECK(strcmp(cf->parser_file.filename, "./test.moss") == 0);
	CHECK(strcmp(cf->parser_file.data, "abcdef") == 0);
	CHECK(cf->parser_file.length == 6);

	CHECK(cf->root == NULL); // FIXME should match fake parser return
	CHECK(cf->flags ==  FILE_FROM_INPUT);
	// FIXME --  FILE_PARSE_FAILED gets set on parse errors

	CHECK(ErrorCount() == 0);

	FakeFilesFree();
	CompileStateFree(&compile_state);
}

void TestSimpleModule()
{
	CompileState compile_state;

	CompileStateInit(&compile_state);

	String name;
	Namespace *ns = &compile_state.root_namespace;

	name.data = "child1";
	name.length = 6;
	ns = NamespaceGetChild(ns, &name);

	name.data = "child2";
	name.length = 6;
	ns = NamespaceGetChild(ns, &name);

	ListInsertFirst(&compile_state.input_modules, ns);

	StringBuffer *sb = StringBufferFromChars("basedir/");
	ListInsertFirst(&compile_state.basedirs, sb);

	FakeFileSet("basedir/import/child1/child2.moss", "abcdef");
	FakeDirectoryAdd("basedir/import/child1/");
	FakeDirectoryAddFile("child2.moss");

	CHECK(PassSearchAndParse(&compile_state));
	CHECK(ns->flags == NAMESPACE_SCANNED);

	CHECK(ns->private_files.first == NULL);
	CHECK(ns->public_files.first != NULL);
	
	CompilerFile *cf = ns->public_files.first->item;
	CHECK(strcmp(cf->path->buffer, "basedir/import/child1/child2.moss") == 0);

	CHECK(ErrorCount() == 0);

	FakeFilesFree();
	FakeDirectoryFree();
	CompileStateFree(&compile_state);
}

void TestPassSearchAndParse(void)
{
	TestSimpleInputFile();
	TestSimpleModule();
}

