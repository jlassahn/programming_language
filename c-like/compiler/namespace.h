
#ifndef INCLUDED_NAMESPACE_H
#define INCLUDED_NAMESPACE_H

#include "compiler/types.h"
#include "compiler/symbol.h"
#include <stdint.h>

typedef struct Namespace Namespace;
struct Namespace
{
	uint32_t flags;
	Namespace *parent;
	StringBuffer *path;
	String stem;

	Map children;  // Map of Namespace
	List public_files; // List of CompilerFile
	List private_files; // List of CompilerFile

	Map symbols; // Map of Symbol
};

typedef enum
{
	NAMESPACE_HAS_INFILE = 0x0001,
	NAMESPACE_SCANNED = 0x0002,
}
NamespaceFlags;

Namespace *NamespaceMakeChild(Namespace *parent, String *name);
Namespace *NamespaceGetChild(Namespace *parent, String *name);
void NamespaceFree(Namespace *root);
Symbol *NamespaceFindSymbol(Namespace *ns, String *name);

void NamespacePrinter(const String *key, void *value, void *depth);

#endif

