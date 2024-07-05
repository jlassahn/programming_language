
#include "unit_test.h"
#include "compiler/memory.h"
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>

static int test_errors = 0;
static int test_count = 0;
static int tests_failed = 0;

void DoCheck(bool x, const char *x_text, int line, const char *file)
{
	if (x)
		return;

	test_errors ++;
	printf("    CHECK(%s) failed at line %d of file %s\n",
			x_text, line ,file);
}

void DoCheckAndExit(bool x, const char *x_text, int line, const char *file)
{
	DoCheck(x, x_text, line, file);
	if (!x)
		exit(70);
}

void RunTest(UnitTest test, const char *name)
{
	int old_errors = test_errors;
	test();
	CHECK(AllocCount() == 0);
	int errs = test_errors - old_errors;

	test_count ++;
	if (errs > 0)
	{
		tests_failed ++;
		printf("FAILED: %s  (%d errors)\n", name, errs);
	}
	else
	{
		printf("PASSED: %s\n", name);
	}
}

int TotalTests(void)
{
	return test_count;
}

int TotalErrors(void)
{
	return test_errors;
}

int TestsFailed(void)
{
	return tests_failed;
}

