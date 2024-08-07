
#include "compiler/fileio.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#ifdef _WIN32

#error "FIXME implement Windows directory operations"
struct DirectorySearch
{
	HANDLE handle;
};

DirectorySearch *DirectorySearchStart(const char *path)
{
	return NULL;
}

const char *DirectorySearchNextFile(DirectorySearch *dir)
{
	return NULL;
}

void DirectorySearchEnd(DirectorySearch *dir)
{
}

bool DoesDirectoryExist(const char *path);
bool DoesFileExist(const char *path);

	/* in Windows
	   pattern = path + "\\*"
	   // if pattern is path it checks if that specific file exists
	   handle = FindFirstFile(pattern, &ffd_out);
	   const char * name = ffd_out.cFileName;
	   handle = FindNextFile(handle, &ffd_out);
	   FindClose(handle);

	   // to check if the file is a directory
	   ffd_out.dwFileAttributes & FILE_ATTRIBUTE_DIRECTORY
	*/

#else

#include <dirent.h>
#include <sys/stat.h>

struct DirectorySearch
{
	DIR *dir;
};

DirectorySearch *DirectorySearchStart(const char *path)
{
	DirectorySearch *ds = Alloc(sizeof(DirectorySearch));
	memset(ds, 0, sizeof(DirectorySearch));

	ds->dir = opendir(path);
	if (!ds->dir)
	{
		Free(ds);
		return NULL;
	}

	return ds;
}

const char *DirectorySearchNextFile(DirectorySearch *dir)
{
	struct dirent *entry;

	while (1)
	{
		entry = readdir(dir->dir);
		if (entry == NULL)
			return NULL;

		if (entry->d_type != DT_REG)
			continue;

		return entry->d_name;
	}
}

void DirectorySearchEnd(DirectorySearch *dir)
{
	closedir(dir->dir);
	memset(dir, 0, sizeof(DirectorySearch));
	Free(dir);
}

bool DoesDirectoryExist(const char *path)
{
	struct stat statbuf;
	if (stat(path, &statbuf) != 0)
		return false;
	if ((statbuf.st_mode & S_IFMT) == S_IFDIR)
		return true;
	return false;
}

bool DoesFileExist(const char *path)
{
	struct stat statbuf;
	if (stat(path, &statbuf) != 0)
		return false;
	if ((statbuf.st_mode & S_IFMT) == S_IFREG)
		return true;
	return false;
}

#endif

OSFile *OSFileOpenRead(const char *path)
{
	FILE *file = fopen(path, "rb");
	return (OSFile *)file;
}

void OSFileClose(OSFile *fp)
{
	FILE *file = (FILE *)fp;
	fclose(file);
}

long OSFileGetSize(OSFile *fp)
{
	FILE *file = (FILE *)fp;

	long original = ftell(file);
	fseek(file, 0, SEEK_END);
	long length = ftell(file);
	fseek(file, original, SEEK_SET);

	return length;
}

long OSFileRead(OSFile *fp, void *data_out, long max_bytes)
{
	FILE *file = (FILE *)fp;

	long read_length = fread(data_out, 1, max_bytes, file);
	return read_length;
}

