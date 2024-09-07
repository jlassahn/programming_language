
#ifndef INCLUDED_NAMESPACE_H
#define INCLUDED_NAMESPACE_H

#include "compiler/types.h"
#include "compiler/symbol.h"
#include "compiler/symbol_table.h"
#include "compiler/parser_node.h"
#include <stdint.h>

typedef struct ImportLink ImportLink;
typedef struct Namespace Namespace;

struct ImportLink
{
	ParserNode *parse;
	bool is_private;
	Namespace *namespace;
};

typedef struct NamespaceSymbols NamespaceSymbols;
struct NamespaceSymbols
{
	List files; // List of CompilerFile
	List imports; // List of ImportLink
	Map exports; // Map of Symbol
	SymbolTable symbol_table;
};

struct Namespace
{
	uint32_t flags;
	Namespace *parent;
	StringBuffer *path;
	String stem;
	Map children;  // Map of Namespace

	NamespaceSymbols public_syms;
	NamespaceSymbols private_syms;
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
Symbol *NamespaceFindPrivateSymbol(Namespace *ns, String *name);

void NamespacePrinter(const String *key, void *value, void *depth);

#endif

