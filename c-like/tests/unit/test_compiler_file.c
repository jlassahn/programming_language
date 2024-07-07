

#include "tests/unit/unit_test.h"
#include "tests/unit/fake_directory.h"
#include "compiler/compiler_file.h"

void TestCompilerFile(void)
{
	FakeFileSet("test.moss", "some file data");

	StringBuffer *path = StringBufferFromChars("test.moss");
	StringBufferLock(path);

	CompilerFile *cf = CompilerFileCreate(path);
	CHECK(cf != NULL);
	CHECK(cf->path == path);
	CHECK(ParserFileRead(&cf->parser_file, path->buffer));
	CompilerFileFree(cf);

	FakeFilesFree();
}

