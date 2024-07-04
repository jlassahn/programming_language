
#include "tests/unit/unit_test.h"
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

/* FIXME test these
typedef struct DirectorySearch DirectorySearch;

DirectorySearch *DirectorySearchStart(const char *path);
const char *DirectorySearchNextFile(DirectorySearch *dir);
void DirectorySearchEnd(DirectorySearch *dir);
*/

// FIXME test this
// bool IsValidPath(const char *txt);

int main(int argc, const char *argv[])
{
	RUN_TEST(TestBasicRead);
	RUN_TEST(TestBadRead);
	RUN_TEST(TestFileExist);

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

