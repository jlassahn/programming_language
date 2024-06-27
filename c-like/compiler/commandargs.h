
#ifndef INCLUDED_COMMANDARGS_H
#define INCLUDED_COMMANDARGS_H

#include <stdint.h>
#include <stdbool.h>

typedef struct ArgStringList ArgStringList;
struct ArgStringList
{
	const char *arg;
	ArgStringList *next;
};

typedef struct CompilerArgs CompilerArgs;
struct CompilerArgs
{
	ArgStringList *inputs;
	ArgStringList *warnings;
	ArgStringList *optimizations;
	ArgStringList *generation;
	ArgStringList *defines;
	ArgStringList *basedirs;
	ArgStringList *versions;
	ArgStringList *outfile;
	ArgStringList *outdir;
	ArgStringList *treefile;
};

const CompilerArgs *ParseArgs(int argc, const char *argv[]);
void FreeArgs(const CompilerArgs *args);
void PrintArgList(ArgStringList *list);


#endif
