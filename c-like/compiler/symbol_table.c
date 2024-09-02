
#include "compiler/symbol_table.h"
#include "compiler/errors.h"
#include <string.h>

void SymbolTableInit(SymbolTable *syms)
{
	memset(syms, 0, sizeof(SymbolTable));
}

void SymbolTableDestroy(SymbolTable *syms)
{
	MapDestroyAll(&syms->root);
}

Symbol *SymbolTableFind(SymbolTable *syms, String *name)
{
	return MapFind(&syms->root, name);
}

bool SymbolTableInsert(SymbolTable *syms, Symbol *sym)
{
	if (SymbolTableFind(syms, &sym->name))
	{
		// FIXME add parse tree info to symbol and do ErrorAt here
		Error(ERROR_DEFINITION, "Redefining an already defined symbol %.*s", sym->name.length, sym->name.data);
		return false;
	}
	MapInsert(&syms->root, &sym->name, sym);
	return true;
}

typedef struct SymbolInsertCtx SymbolInsertCtx;
struct SymbolInsertCtx
{
	SymbolTable *table;
	bool ret;
};

static void SymbolInserter(const String *key, void *value, void *vctx)
{
	SymbolInsertCtx *ctx = vctx;
	Symbol *sym = value;
	if (!SymbolTableInsert(ctx->table, sym))
		ctx->ret = false;
}

bool SymbolTableInsertMap(SymbolTable *syms, Map *map)
{
	SymbolInsertCtx ctx;
	ctx.table = syms;
	ctx.ret = true;

	MapIterate(map, SymbolInserter, &ctx);

	return ctx.ret;
}

