
#include "tests/unit/unit_test.h"
#include "tests/unit/fake_directory.h"
#include "compiler/fileio.h"
#include <string.h>

void TestFakeDir(void)
{
	FakeDirectoryAdd("d1/d2/d3/");
	FakeDirectoryAddFile("f1");
	FakeDirectoryAddFile("f2");

	DirectorySearch *ds = DirectorySearchStart("d1/d2/d3/");
	CHECK(ds != NULL);
	if (ds != NULL)
	{
		CHECK(strcmp(DirectorySearchNextFile(ds), "f1") == 0);
		CHECK(strcmp(DirectorySearchNextFile(ds), "f2") == 0);
		CHECK(DirectorySearchNextFile(ds) == NULL);
		DirectorySearchEnd(ds);
	}

	ds = DirectorySearchStart("d1/d2/d3/x");
	CHECK(ds == NULL);

	FakeDirectoryFree();
}

