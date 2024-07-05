
#include "compiler/search.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/fileio.h"
#include <string.h>

struct SearchFiles
{
	List *basedirs;
	const char *part;
	const char *path;
	const char **filter;
	ListEntry *basedir;
	StringBuffer *pathbuf;
	DirectorySearch *ds;
};


SearchFiles *SearchFilesStart(
		List *basedirs,
		const char *part,
		const char *path,
		const char *filter[])
{
	SearchFiles *sf = Alloc(sizeof(SearchFiles));
	sf->basedirs = basedirs;
	sf->part = part;
	sf->path = path;
	sf->filter = filter;

	sf->basedir = basedirs->first;
	sf->pathbuf = StringBufferCreateEmpty(0);
	sf->ds = NULL;

	return sf;
}

StringBuffer *SearchFilesNext(SearchFiles *sf)
{
	while (true)
	{
		if (sf->basedir == NULL)
			return NULL;

		if (sf->ds == NULL)
		{
			StringBuffer *path = sf->pathbuf;
			StringBufferClear(path);
			path = StringBufferAppendBuffer(path, sf->basedir->item);
			path = StringBufferAppendChars(path, sf->part);
			path = StringBufferAppendChars(path, sf->path);
			sf->pathbuf = path;

			sf->ds = DirectorySearchStart(path->buffer);
			if (sf->ds == NULL)
			{
				sf->basedir = sf->basedir->next;
				continue;
			}
		}

		const char *filename = DirectorySearchNextFile(sf->ds);
		if (filename == NULL)
		{
			DirectorySearchEnd(sf->ds);
			sf->ds = NULL;
			sf->basedir = sf->basedir->next;
			continue;
		}

		int length = strlen(filename);
		for (int i=0; sf->filter[i]!=NULL; i++)
		{
			int flength = strlen(sf->filter[i]);
			if (strcmp(sf->filter[i], filename +length-flength) == 0)
			{
				StringBuffer *ret;
				ret = StringBufferFromString(&sf->pathbuf->string);
				ret = StringBufferAppendChars(ret, filename);
				return ret;
			}
		}
	}
}

void SearchFilesEnd(SearchFiles *sf)
{
	if (sf->ds)
		DirectorySearchEnd(sf->ds);
	StringBufferFree(sf->pathbuf);
	Free(sf);
}

