
#ifndef INCLUDED_ERRORS_H
#define INCLUDED_ERRORS_H

#include "compiler/parser_file.h"

typedef enum
{
	ERROR_FILE,
	ERROR_PARSER,
}
ErrorCategory;

void Error(ErrorCategory cat, const char *text, ...);

void Warning(ErrorCategory cat, const char *text, ...);

void ErrorAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...);

void WarningAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...);

int ErrorCount(void);
int WarningCount(void);

#endif

