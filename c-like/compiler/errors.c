
#include "compiler/errors.h"
#include <stdarg.h>
#include <stdio.h>

void Error(const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "ERROR: ");
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);
}

void ErrorAt(const char *filename, FilePosition *pos, const char *text, ...)
{
	va_list args;
	va_start(args, text);
	fprintf(stderr, "ERROR[%s:%ld:%ld]: ", filename, pos->line+1, pos->byte_in_line+1);
	vfprintf(stderr, text, args);
	fprintf(stderr, "\n");
	va_end(args);
}

