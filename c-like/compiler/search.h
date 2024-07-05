
#ifndef INCLUDED_SEARCH_H
#define INCLUDED_SEARCH_H

#include "compiler/types.h"

typedef struct SearchFiles SearchFiles;

SearchFiles *SearchFilesStart(
		List *basedirs,
		const char *part,
		const char *path,
		const char *filter[]);

StringBuffer *SearchFilesNext(SearchFiles *sf);
void SearchFilesEnd(SearchFiles *sf);

#endif

