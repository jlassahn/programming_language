
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

/*!
  @brief Checks whether a string is a reasonable path name.

  Returns false for some weird things that are valid Unix
  path names, but I don't want to deal with.

  @param path The string to check.
  @return true if the string is a usable path.
*/
bool IsValidPath(const char *path);

/*!
  @brief Cleans up a path.

  The path string is reallocated and returned.
  The result will have forward slash as the directory separator, and
  will always end with a slash.

  @param path The path the be normalized.
  @return A new pointer to the path.
*/
USE_RESULT
StringBuffer *NormalizePath(StringBuffer *path);

bool AppendPathList(List *list, const char *env);

#endif

