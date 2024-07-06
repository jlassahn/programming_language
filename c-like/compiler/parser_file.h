
#ifndef INCLUDED_PARSER_FILE_H
#define INCLUDED_PARSER_FILE_H

#include <stdint.h>
#include <stdbool.h>

typedef struct FilePosition FilePosition;
typedef struct FilePositionRange FilePositionRange;
typedef struct ParserFile ParserFile;

struct FilePosition
{
	long offset;
	long line;
	long byte_in_line;
};

struct FilePositionRange
{
	ParserFile *file;
	FilePosition start;
	FilePosition end; // one past the final character in the range
};

struct ParserFile
{
	const char *filename;
	const char *data;
	long length;

	FilePosition current_pos;
	int parser_result;
};


bool ParserFileRead(ParserFile *file, const char *filename);
void FileFree(ParserFile *file);

bool FileMatchAndConsume(ParserFile *file, const char *text);
bool FileMatch(ParserFile *file, const char *text);
void FileConsume(ParserFile *file, int n);
const char *FileGet(ParserFile *file);

#endif

