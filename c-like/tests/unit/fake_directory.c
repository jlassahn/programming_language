
#include "tests/unit/fake_directory.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/fileio.h"
#include <string.h>

typedef struct FakeDir FakeDir;
struct FakeDir
{
	const char *path;
	List files;
};

struct DirectorySearch
{
	FakeDir *dir;
	ListEntry *file_entry;
};

static List dirs;

void FakeDirectoryAdd(const char *path)
{
	FakeDir *dir = Alloc(sizeof(FakeDir));
	dir->path = path;
	ListInsertLast(&dirs, dir);
}

void FakeDirectoryAddFile(const char *filename)
{
	FakeDir *dir = dirs.last->item;
	ListInsertLast(&dir->files, (void *)filename);
}

void FakeDirectoryFree(void)
{
	while (dirs.first != NULL)
	{
		FakeDir *dir = ListRemoveFirst(&dirs);
		while (ListRemoveFirst(&dir->files) != NULL)
			;
		Free(dir);
	}
}

DirectorySearch *DirectorySearchStart(const char *path)
{
	FakeDir *dir = NULL;
	for (ListEntry *entry=dirs.first; entry!=NULL; entry=entry->next)
	{
		FakeDir *d = entry->item;
		if (strcmp(d->path, path) == 0)
		{
			dir = d;
			break;
		}
	}
	if (dir == NULL)
		return NULL;

	DirectorySearch *ds = Alloc(sizeof(DirectorySearch));
	ds->dir = dir;
	ds->file_entry = dir->files.first;
	return ds;
}

const char *DirectorySearchNextFile(DirectorySearch *ds)
{
	if (ds->file_entry == NULL)
		return NULL;

	const char *filename = ds->file_entry->item;
	ds->file_entry = ds->file_entry->next;
	return filename;
}

void DirectorySearchEnd(DirectorySearch *ds)
{
	Free(ds);
}
