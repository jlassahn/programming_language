
#ifndef INCLUDED_UNIT_TEST_H
#define INCLUDED_UNIT_TEST_H

#include <stdbool.h>

#define CHECK(x) DoCheck(x, #x, __LINE__, __FILE__)
#define CHECK_AND_EXIT(x) DoCheckAndExit(x, #x, __LINE__, __FILE__)

void DoCheck(bool x, const char *x_text, int line, const char *file);
void DoCheckAndExit(bool x, const char *x_text, int line, const char *file);

#define RUN_TEST(x) RunTest(x, #x)

typedef void (*UnitTest)(void);
void RunTest(UnitTest test, const char *name);

int TotalTests(void);
int TotalErrors(void);
int TestsFailed(void);

#endif

