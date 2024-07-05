
#include "tests/unit/unit_test.h"
#include "tests/unit/fake_directory.h"
#include "compiler/search.h"
#include "compiler/types.h"
#include <string.h>

void TestSearch(void)
{

	FakeDirectoryAdd("base1/import/system/");
	FakeDirectoryAddFile("clib.part1.moss");
	FakeDirectoryAddFile("clib.part2.moss");
	FakeDirectoryAddFile("filelib.moss");
	FakeDirectoryAddFile("filelib.package");

	FakeDirectoryAdd("base2/import/system/");
	FakeDirectoryAddFile("clib.part3.moss");

	FakeDirectoryAdd("base2/lib/system/");
	FakeDirectoryAddFile("clib.a");
	FakeDirectoryAddFile("clib.so");
	FakeDirectoryAddFile("clib.lib");
	FakeDirectoryAddFile("clib.dll");
	FakeDirectoryAddFile("clib.package");

	List basedirs = {NULL, NULL};
	ListInsertLast(&basedirs, StringBufferFromChars("base1/"));
	ListInsertLast(&basedirs, StringBufferFromChars("base2/"));

	const char *filter[] =
	{
		".moss",
		NULL
	};

	SearchFiles *sf = SearchFilesStart(&basedirs, "import/", "system/", filter);
	CHECK(sf != NULL);
	if (sf != NULL)
	{
		StringBuffer *file;

		file = SearchFilesNext(sf);
		CHECK(file != NULL);
		CHECK(0 == strcmp(file->buffer,
					"base1/import/system/clib.part1.moss"));
		StringBufferFree(file);

		file = SearchFilesNext(sf);
		CHECK(file != NULL);
		CHECK(0 == strcmp(file->buffer,
					"base1/import/system/clib.part2.moss"));
		StringBufferFree(file);

		file = SearchFilesNext(sf);
		CHECK(file != NULL);
		CHECK(0 == strcmp(file->buffer,
					"base1/import/system/filelib.moss"));
		StringBufferFree(file);

		file = SearchFilesNext(sf);
		CHECK(file != NULL);
		CHECK(0 == strcmp(file->buffer,
					"base2/import/system/clib.part3.moss"));
		StringBufferFree(file);

		file = SearchFilesNext(sf);
		CHECK(file == NULL);

		SearchFilesEnd(sf);
	}

	while (basedirs.first)
	{
		StringBufferFree(ListRemoveFirst(&basedirs));
	}

	FakeDirectoryFree();
}

