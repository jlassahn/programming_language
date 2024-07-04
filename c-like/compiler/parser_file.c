
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/errors.h"
#include "compiler/parser_file.h"
#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

ParserFile *ParserFileRead(const char *filename)
{
	FILE *fp = fopen(filename, "rb");
	if (!fp)
	{
		Error(ERROR_FILE, "Unable to open file. Filename(%s) Reason(%s)",
				filename, strerror(errno));
		return NULL;
	}

	// FIXME consider making a way to turn stdin or other non-seekable streams
	// into source files.
	fseek(fp, 0, SEEK_END);
	long length = ftell(fp);
	fseek(fp, 0, SEEK_SET);

	// add some trailing zeros onto the end of the file buffer so we can
	// do fixed-length memory compares looking for keywords and such without
	// overrunning at End of File.
	int padding = 16;

	long buffer_size = 
		sizeof(ParserFile) +
		strlen(filename) + 1 +
		length + padding;

	char *p = Alloc(buffer_size);

	ParserFile *file = (ParserFile *)p;
	p += sizeof(ParserFile);

	strcpy(p, filename);
	file->filename = p;
	p += strlen(filename) + 1;

	long read_length = fread(p, 1, length,  fp);
	if (read_length != length)
	{
		Error(ERROR_FILE, "Unable to read file. Filename(%s) Reason(%s)",
				filename, strerror(errno));
		fclose(fp);
		Free(file);
		return NULL;
	}

	file->data = p;
	file->length = length;
	file->parser_result = -1;

	fclose(fp);
	return file;
}

void FileFree(ParserFile *file)
{
	Free(file);
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


