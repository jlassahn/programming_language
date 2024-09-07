

#include "tests/unit/unit_test.h"
#include "tests/unit/fake_directory.h"
#include "tests/unit/fake_errors.h"
#include "tests/unit/fake_parser.h"
#include "tests/unit/fake_nodes.h"
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
	CHECK(ns->private_syms.files.first != NULL);
	CHECK(ns->private_syms.files.first->item == cf);
	CHECK(ns->public_syms.files.first == NULL);
	CHECK(ns->flags == NAMESPACE_HAS_INFILE);
	CHECK(cf->namespace == ns);
	CHECK(strcmp(cf->parser_file.filename, "./test.moss") == 0);
	CHECK(strcmp(cf->parser_file.data, "abcdef") == 0);
	CHECK(cf->parser_file.length == 6);

	CHECK(cf->root == NULL);
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
	ns = NamespaceMakeChild(ns, &name);

	name.data = "child2";
	name.length = 6;
	ns = NamespaceMakeChild(ns, &name);

	ListInsertFirst(&compile_state.input_modules, ns);

	StringBuffer *sb = StringBufferFromChars("basedir/");
	ListInsertFirst(&compile_state.basedirs, sb);

	FakeFileSet("basedir/import/child1/child2.moss", "abcdef");
	FakeDirectoryAdd("basedir/import/child1/");
	FakeDirectoryAddFile("child2.moss");

	CHECK(PassSearchAndParse(&compile_state));
	CHECK(ns->flags == NAMESPACE_SCANNED);

	CHECK(ns->private_syms.files.first == NULL);
	CHECK(ns->public_syms.files.first != NULL);
	
	CompilerFile *cf = ns->public_syms.files.first->item;
	CHECK(strcmp(cf->path->buffer, "basedir/import/child1/child2.moss") == 0);

	CHECK(ErrorCount() == 0);

	FakeFilesFree();
	FakeDirectoryFree();
	CompileStateFree(&compile_state);
}

static void CheckImportEntry(ListEntry **entry, bool is_private, Namespace *ns)
{
	CHECK(*entry != NULL);
	if (*entry == NULL)
		return;

	ImportLink *import = (*entry)->item;
	CHECK(is_private == import->is_private);
	CHECK(ns == import->namespace);

	(*entry) = (*entry)->next;
}

void TestImportList(void)
{
	CompileState compile_state;

	CompileStateInit(&compile_state);

	CompilerFile *cf = CompilerFileCreate(
			StringBufferFromChars("./test.moss"));
	cf->flags = FILE_FROM_INPUT;

	ListInsertFirst(&compile_state.input_files, cf);

	PushNodeStack(MakeNode(&SYM_EMPTY, 0, NULL));

	PushNodeStack(MakeNodeFakeValue(&SYM_IDENTIFIER, "child1_1"));
	PushNodeStack(MakeNodeFakeValue(&SYM_IDENTIFIER, "child2_1"));
	MakeNodeOnStack(&SYM_DOT_OP, 2);
	PushNodeStack(MakeNodeFakeValue(&SYM_IDENTIFIER, "child3_1"));
	MakeNodeOnStack(&SYM_DOT_OP, 2);
	MakeNodeOnStack(&SYM_IMPORT_PRIVATE, 1);
	MakeNodeOnStack(&SYM_LIST, 2);

	PushNodeStack(MakeNodeFakeValue(&SYM_IDENTIFIER, "child1_2"));
	MakeNodeOnStack(&SYM_IMPORT, 1);
	MakeNodeOnStack(&SYM_LIST, 2);

	PushNodeStack(MakeNodeFakeValue(&SYM_IDENTIFIER, "child1_3"));
	MakeNodeOnStack(&SYM_IMPORT, 1);
	MakeNodeOnStack(&SYM_LIST, 2);

	ParserNode *list = GetNodeStackTop();
	// PrintNodeTree(stdout, list);

	FakeParserSet("./test.moss", list, 0);
	FakeFileSet("./test.moss", "abcdef");

	String name;
	Namespace *ns_root = &compile_state.root_namespace;

	name.data = "child1_1";
	name.length = strlen(name.data);
	Namespace *ns_m1 = NamespaceMakeChild(ns_root, &name);
	name.data = "child2_1";
	name.length = strlen(name.data);
	ns_m1 = NamespaceMakeChild(ns_m1, &name);
	name.data = "child3_1";
	name.length = strlen(name.data);
	ns_m1 = NamespaceMakeChild(ns_m1, &name);

	name.data = "child1_2";
	name.length = strlen(name.data);
	Namespace *ns_m2 = NamespaceMakeChild(ns_root, &name);

	name.data = "child1_3";
	name.length = strlen(name.data);
	Namespace *ns_m3 = NamespaceMakeChild(ns_root, &name);

	// prevent the code from trying to search fake namespaces for files
	// by marking them as already scanned.
	// ns_m1->flags |= NAMESPACE_SCANNED; // let this one get scanned
	ns_m2->flags |= NAMESPACE_SCANNED;
	ns_m3->flags |= NAMESPACE_SCANNED;

	StringBuffer *sb = StringBufferFromChars("basedir/");
	ListInsertFirst(&compile_state.basedirs, sb);
	FakeFileSet("basedir/import/child1_1/child2_1/child3_1.moss", "abcdef");
	FakeDirectoryAdd("basedir/import/child1_1/child2_1/");
	FakeDirectoryAddFile("child3_1.moss");

	CHECK(PassSearchAndParse(&compile_state));

	CHECK(cf->root == list);
	CHECK(cf->flags ==  FILE_FROM_INPUT);

	ListEntry *entry = cf->imports.first;
	CheckImportEntry(&entry, true, ns_m1);
	CheckImportEntry(&entry, false, ns_m2);
	CheckImportEntry(&entry, false, ns_m3);
	CHECK(entry == NULL);

	// check that ns_m1 was recursively scanned
	CHECK(ns_m1->public_syms.files.first != NULL);

	FreeFakeNodeValues();
	FakeDirectoryFree();
	FakeFilesFree();
	FakeParserFree();
	CompileStateFree(&compile_state);
}

void TestPassSearchAndParse(void)
{
	TestSimpleInputFile();
	TestSimpleModule();
	TestImportList();
}

