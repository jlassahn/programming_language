
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/errors.h"
#include "compiler/fileio.h"
#include "compiler/parser_file.h"
#include <stdint.h>
#include <stdbool.h>
//#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

bool ParserFileRead(ParserFile *file, const char *filename)
{
	OSFile *fp = OSFileOpenRead(filename);
	if (!fp)
	{
		Error(ERROR_FILE, "Unable to open file. Filename(%s) Reason(%s)",
				filename, strerror(errno));
		return false;
	}

	// FIXME consider making a way to turn stdin or other non-seekable streams
	// into source files.
	long length = OSFileGetSize(fp);

	// add some trailing zeros onto the end of the file buffer so we can
	// do fixed-length memory compares looking for keywords and such without
	// overrunning at End of File.
	int padding = 16;

	long buffer_size = strlen(filename) + 1 + length + padding;

	char *p = Alloc(buffer_size);

	strcpy(p, filename);
	file->filename = p;
	p += strlen(filename) + 1;

	long read_length = OSFileRead(fp, p, length);
	OSFileClose(fp);

	if (read_length != length)
	{
		Error(ERROR_FILE, "Unable to read file. Filename(%s) Reason(%s)",
				filename, strerror(errno));
		Free(p);
		file->filename = NULL;
		return false;
	}

	file->data = p;
	file->length = length;
	file->parser_result = -1;

	return true;
}

void FileFree(ParserFile *file)
{
	void *p = (void *)file->filename;
	if (p != NULL)
		Free(p);
	file->data = NULL;
	file->filename = NULL;
	file->length = 0;
	file->parser_result = -1;
}

bool FileMatchAndConsume(ParserFile *file, const char *text)
{
	const char *cur =  &file->data[file->current_pos.offset];
	int len = strlen(text);
	if (strncmp(cur, text, len) == 0)
	{
		FileConsume(file, len);
		return true;
	}
	return false;
}

bool FileMatch(ParserFile *file, const char *text)
{
	const char *cur =  &file->data[file->current_pos.offset];
	int len = strlen(text);
	return (strncmp(cur, text, len) == 0);
}

void FileConsume(ParserFile *file, int n)
{
	const char *cur =  &file->data[file->current_pos.offset];
	for (int i=0; i<n; i++)
	{
		if (cur[i] == 0)
			return;
		file->current_pos.offset ++;
		file->current_pos.byte_in_line ++;
		if (cur[i] == '\n')
		{
			file->current_pos.line ++;
			file->current_pos.byte_in_line = 0;
		}
	}
}

const char *FileGet(ParserFile *file)
{
	return &file->data[file->current_pos.offset];
}


