
#ifndef INCLUDED_NAMESPACE_H
#define INCLUDED_NAMESPACE_H

typedef struct Namespace Namespace;

struct Namespace
{
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

