

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

