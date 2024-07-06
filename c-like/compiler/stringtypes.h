
#ifndef INCLUDED_STRINGTYPES_H
#define INCLUDED_STRINGTYPES_H

#include "compiler/types.h"
#include <stdbool.h>

bool IsSpace(char x);
bool IsLetter(char x);
bool IsDigit(char x);

// check if strings are syntactically valid names
bool IsValidNamespace(const char *txt);
bool IsValidNamespaceName(String *str);

#endif

