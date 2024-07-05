
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
	const char *head;
	const char **tails;
	ListEntry *basedir;
	StringBuffer *pathbuf;
	DirectorySearch *ds;
};


SearchFiles *SearchFilesStart(
		List *basedirs,
		const char *part,
		const char *path,
		const char *head,
		const char *tails[])
{
	SearchFiles *sf = Alloc(sizeof(SearchFiles));
	sf->basedirs = basedirs;
	sf->part = part;
	sf->path = path;
	sf->head = head;
	sf->tails = tails;

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

		int hlength = strlen(sf->head);
		if (strncmp(sf->head, filename, hlength) != 0)
			continue;

		int length = strlen(filename);
		for (int i=0; sf->tails[i]!=NULL; i++)
		{
			int flength = strlen(sf->tails[i]);
			if (strcmp(sf->tails[i], filename +length-flength) == 0)
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

