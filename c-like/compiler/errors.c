
#include "compiler/errors.h"
#include <stdarg.h>
#include <stdio.h>

static int error_count;
static int warning_count;

void Error(ErrorCategory cat, const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "ERROR: ");
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);

	error_count ++;
}

void ErrorAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "ERROR[%s:%ld:%ld]: ", filename, pos->line+1, pos->byte_in_line+1);
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);

	error_count ++;
}

void Warning(ErrorCategory cat, const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "WARNING: ");
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);

	warning_count ++;
}

void WarningAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "WARNING[%s:%ld:%ld]: ", filename, pos->line+1, pos->byte_in_line+1);
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);

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


