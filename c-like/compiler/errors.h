
#ifndef INCLUDED_ERRORS_H
#define INCLUDED_ERRORS_H

#include "parser_file.h"

void Error(const char *text, ...);
void ErrorAt(const char *filename, FilePosition *pos, const char *text, ...);

#endif

