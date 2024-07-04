
#include "unit_test.h"
#include <stdio.h>

#ifdef TEST_TESTS
void TestFail(void)
{
	CHECK(false);
}

void TestPass(void)
{
	CHECK(true);
}
#endif



void TestMap(void);

int main(void)
{
#ifdef TEST_TESTS
	RUN_TEST(TestFail);
	RUN_TEST(TestPass);
#endif
	RUN_TEST(TestMap);

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
