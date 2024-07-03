
#ifndef INCLUDED_FILEIO_H
#define INCLUDED_FILEIO_H

#include <stdbool.h>

#ifdef _WIN32
#define PATH_SEPARATOR '\\'
#define PATH_SEPARATOR_STRING "\\"
#else
#define PATH_SEPARATOR '/'
#define PATH_SEPARATOR_STRING "/"
#endif

typedef struct DirectorySearch DirectorySearch;

DirectorySearch *DirectorySearchStart(const char *path);
const char *DirectorySearchNextFile(DirectorySearch *dir);
void DirectorySearchEnd(DirectorySearch *dir);

bool DoesDirectoryExist(const char *path);
bool DoesFileExist(const char *path);

// actually returns false for some weird things that are valid Unix
// path names, but I don't want to deal with.
bool IsValidPath(const char *txt);

#endif

