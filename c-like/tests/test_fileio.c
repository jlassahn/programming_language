
#include "tests/unit/unit_test.h"
#include "compiler/stringtypes.h"
#include "compiler/fileio.h"
#include <stdio.h>
#include <string.h>

void TestBasicRead(void)
{
	char buffer[16];

	OSFile *fp = OSFileOpenRead("tests/test_files/subdir1/length10.txt");
	CHECK(fp != NULL);
	if (fp)
	{
		CHECK(10 == OSFileGetSize(fp));
		CHECK(5 == OSFileRead(fp, buffer, 5));
		CHECK(0 == memcmp(buffer, "12345", 5));
		CHECK(10 == OSFileGetSize(fp));
		CHECK(5 == OSFileRead(fp, buffer, 5));
		CHECK(0 == memcmp(buffer, "6789\n", 5));
		OSFileClose(fp);
	}
}

void TestBadRead(void)
{
	char buffer[16];
	OSFile *fp = OSFileOpenRead("tests/test_files/subdir1/notreal.txt");
	CHECK(fp == NULL);
	fp = OSFileOpenRead("tests/test_files/subdir1/length10.txt");
	CHECK(fp != NULL);
	if (fp)
	{
		CHECK(10 == OSFileRead(fp, buffer, 11));
		CHECK(0 == memcmp(buffer, "123456789\n", 10));
		CHECK(0 == OSFileRead(fp, buffer, 11));
		OSFileClose(fp);
	}
}

void TestFileExist(void)
{
	CHECK(DoesDirectoryExist("tests/test_files/subdir1"));
	CHECK(DoesDirectoryExist("tests/test_files/subdir2"));
	CHECK(DoesDirectoryExist("tests/test_files"));
	CHECK(DoesDirectoryExist("tests/test_files/"));
	CHECK(!DoesDirectoryExist("tests/test_files/notreal"));
	CHECK(DoesFileExist("tests/test_files/subdir1/length10.txt"));
	CHECK(!DoesFileExist("tests/test_files/subdir1/notreal.txt"));
	CHECK(!DoesFileExist("tests/test_files/subdir1"));
}

void TestDirectorySearch(void)
{
	DirectorySearch *ds = NULL;

	ds = DirectorySearchStart("tests/test_files/subdir1");
	CHECK(ds != NULL);
	if (ds != NULL)
	{
		CHECK(0 == strcmp(DirectorySearchNextFile(ds), "length10.txt"));
		CHECK(DirectorySearchNextFile(ds) == NULL);
		DirectorySearchEnd(ds);
	}
	ds = DirectorySearchStart("tests/test_files");
	CHECK(ds != NULL);
	if (ds != NULL)
	{
		// should not find subdirectories
		CHECK(DirectorySearchNextFile(ds) == NULL);
		DirectorySearchEnd(ds);
	}
}

int main(int argc, const char *argv[])
{
	RUN_TEST(TestBasicRead);
	RUN_TEST(TestBadRead);
	RUN_TEST(TestFileExist);
	RUN_TEST(TestDirectorySearch);

	int errs = TotalErrors();
	int failed_tests = TestsFailed();
	int test_count = TotalTests();

	if (failed_tests > 0)
	{
		printf("%d/%d tests failed.  %d total errors\n",
				failed_tests, test_count, errs);
	}
	return failed_tests;
}

