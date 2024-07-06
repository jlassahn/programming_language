
#include "tests/unit/unit_test.h"
#include "compiler/errors.h"

static int error_count;
static int warning_count;

void Error(ErrorCategory cat, const char *text, ...)
{
	error_count ++;
}


void Warning(ErrorCategory cat, const char *text, ...)
{
	warning_count ++;
}

void ErrorAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...)
{
	error_count ++;
}


void WarningAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...)
{
	warning_count ++;
}

int ErrorCount(void)
{
	return error_count;
}

int WarningCount(void)
{
	return warning_count;
}

