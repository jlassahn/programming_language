
#ifndef INCLUDED_FAKE_DIRECTORY_H
#define INCLUDED_FAKE_DIRECTORY_H

#include "compiler/fileio.h"

void FakeDirectoryAdd(const char *path);
void FakeDirectoryAddFile(const char *filename);
void FakeDirectoryFree(void);

OSFile *FakeFileSet(const char *path, const char *data);
void FakeFilesFree(void);

#endif

