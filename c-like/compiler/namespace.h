
#ifndef INCLUDED_NAMESPACE_H
#define INCLUDED_NAMESPACE_H

#include <stdint.h>

/*
CompilerState:
	root_namespace
		child_namespaces...
			flags (initialized, private, etc)
			file_list
				file...
					symbol_declarations
					using_statements
					import_statements
					file_namespace
						using_namespace_stuff
			symbols
*/

typedef struct CompilerFile CompilerFile;

typedef struct FileList FileList;
struct FileList
{
	CompilerFile *head;
	CompilerFile *tail;
};

typedef struct Namespace Namespace;
struct Namespace
{
	uint32_t flags;
	FileList files;
	MapNamespace namespaces;
	MapSymbol symbols;
};

struct SymbolTable
{
	Namespace global_namespace;
	Namespace file_scope;
	NamespaceStack local_scope;
};

#endif

