
#include "tests/unit/unit_test.h"
#include "tests/unit/fake_errors.h"

#include "compiler/stringtypes.h"
#include <string.h>

void TestIsValidPath(void)
{
	CHECK(!IsValidPath(""));
	CHECK(!IsValidPath("\n"));
	CHECK(!IsValidPath("this/a\ttab/"));
	CHECK(!IsValidPath("a//b"));
	CHECK(!IsValidPath(" "));
	CHECK(!IsValidPath(" x"));
	CHECK(!IsValidPath("x "));
	CHECK(!IsValidPath("a/ /b"));
	CHECK(!IsValidPath("a/ x/b"));
	CHECK(!IsValidPath("a/x /b"));

	CHECK(IsValidPath("."));
	CHECK(IsValidPath("./x"));
	CHECK(IsValidPath("./x/"));
	CHECK(IsValidPath("x"));
	CHECK(IsValidPath("x/y"));
	CHECK(IsValidPath("/y"));
	CHECK(IsValidPath("C:/x"));
	CHECK(IsValidPath("a long file with spaces"));

	CHECK(IsValidPath("C:\\x"));
}

static void RemoveStringAndCheck(List *list, const char *match)
{
	StringBuffer *sb = ListRemoveFirst(list);
	CHECK(sb != NULL);
	CHECK(strcmp(sb->buffer, match) == 0);
	StringBufferFree(sb);
}

void TestAppendPathList(void)
{
	List list;
	list.first = NULL;
	list.last = NULL;

	CHECK(AppendPathList(&list, "a;b;c;d"));
	RemoveStringAndCheck(&list, "a/");
	RemoveStringAndCheck(&list, "b/");
	RemoveStringAndCheck(&list, "c/");
	RemoveStringAndCheck(&list, "d/");
	CHECK(list.first == NULL);

	// Windows style paths with :
	CHECK(AppendPathList(&list, "C:\\moss;.\\build"));
	RemoveStringAndCheck(&list, "C:/moss/");
	RemoveStringAndCheck(&list, "./build/");
	CHECK(list.first == NULL);

	// Paths with spaces
	CHECK(AppendPathList(&list, "some path;some other path"));
	RemoveStringAndCheck(&list, "some path/");
	RemoveStringAndCheck(&list, "some other path/");
	CHECK(list.first == NULL);

	//Invalid paths
	ClearErrorCounts();
	CHECK(!AppendPathList(&list, "good; bad spaces ;good;;weird\nchars\t;double//slash"));
	RemoveStringAndCheck(&list, "good/");
	RemoveStringAndCheck(&list, "good/");
	CHECK(list.first == NULL);
	CHECK(ErrorCount() == 4);
	ClearErrorCounts();

	// empty
	CHECK(AppendPathList(&list, ""));
	CHECK(list.first == NULL);
	CHECK(ErrorCount() == 0);

}

