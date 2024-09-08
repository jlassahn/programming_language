
#include "tests/unit/unit_test.h"
#include <stdio.h>

void SimpleVariable(void);

int main(void)
{
	RUN_TEST(SimpleVariable);

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

