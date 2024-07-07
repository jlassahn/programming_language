
#include "tests/unit/fake_directory.h"
#include "tests/unit/unit_test.h"
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

typedef struct FakeFile FakeFile;
struct FakeFile
{
	const char *path;
	const char *data;
	bool is_open;
	int length;
	int index;
};


static List files;

OSFile *FakeFileSet(const char *path, const char *data)
{
	FakeFile *file = Alloc(sizeof(FakeFile));
	file->path = path;
	file->data = data;
	file->length = strlen(data);

	ListInsertFirst(&files, file);

	return (OSFile *)file;
}

void FakeFilesFree(void)
{
	FakeFile *file;
	while (true)
	{
		file = ListRemoveFirst(&files);
		if (file == NULL)
			break;

		CHECK(!file->is_open);
		Free(file);
	}
}

OSFile *OSFileOpenRead(const char *path)
{
	for (ListEntry *entry=files.first; entry!=NULL; entry=entry->next)
	{
		FakeFile *file = entry->item;

		if (strcmp(path, file->path) == 0)
		{
			CHECK(!file->is_open);
			file->is_open = true;
			file->index = 0;
			return (OSFile *)file;
		}
	}

	return NULL;
}

void OSFileClose(OSFile *fp)
{
	FakeFile *file = (FakeFile *)fp;

	CHECK(file->is_open);
	file->is_open = false;
}

long OSFileGetSize(OSFile *fp)
{
	FakeFile *file = (FakeFile *)fp;
	return file->length;
}

long OSFileRead(OSFile *fp, void *data_out, long max_bytes)
{
	FakeFile *file = (FakeFile *)fp;

	CHECK(file->is_open);
	long length = file->length - file->index;
	if (length > max_bytes)
		length = max_bytes;

	memcpy(data_out, file->data+file->index, length);
	file->index += length;
	return length;
}

