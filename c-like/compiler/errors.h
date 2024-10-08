
#ifndef INCLUDED_ERRORS_H
#define INCLUDED_ERRORS_H

#include "compiler/parser_file.h"
#include "compiler/parser_node.h"

typedef enum
{
	ERROR_INTERNAL,
	ERROR_FILE,
	ERROR_PARSER,
	ERROR_DEFINITION,
}
ErrorCategory;

void Error(ErrorCategory cat, const char *text, ...);

void Warning(ErrorCategory cat, const char *text, ...);

void ErrorAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...);

void WarningAt(ErrorCategory cat, const char *filename, FilePosition *pos,
		const char *text, ...);

void ErrorAtNode(ErrorCategory cat, ParserNode *node, const char *text, ...);

void WarningAtNode(ErrorCategory cat, ParserNode *node, const char *text, ...);

int ErrorCount(void);
int WarningCount(void);

#endif

