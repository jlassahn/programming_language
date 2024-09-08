
#include "compiler/passes.h"
#include "compiler/compile_state.h"
#include "compiler/eval.h"
#include "compiler/errors.h"
#include <stdio.h>

typedef struct ResolveGlobalsCtx ResolveGlobalsCtx;
struct ResolveGlobalsCtx
{
	CompileState *state;
	bool ret;
};


static void ResolveSymbols(const String *key, void *value, void *vctx)
{
	Symbol *sym = value;
	ResolveGlobalsCtx *ctx = vctx;

	if (sym->definitions.first->next != NULL)
	{
		ParserNode *node = sym->definitions.first->next->item;
		ErrorAtNode(ERROR_INTERNAL, node, "FIXME multiple definitions not handled.");
	}

	EvalContext eval_ctx;
	if (sym->flags & SYM_PRIVATE)
		eval_ctx.symbol_table = &sym->exported_from->private_syms.symbol_table;
	else
		eval_ctx.symbol_table = &sym->exported_from->public_syms.symbol_table;

	ParserNode *node = sym->definitions.first->item;
	ParserNode *dtype_node = node->children[0];
	// ParserNode *props_node = node->children[2];
	// ParserNode *value_node = node->children[3];

	DataValue dv;
	if (!EvalExpression(&dv, dtype_node, &eval_ctx))
	{
		// FIXME figure out how not to get two errors for the public and private versions of the symbol
		ctx->ret = false;
		return;
	}

	if (dv.value_type != VTYPE_DTYPE)
	{
		ErrorAtNode(ERROR_DEFINITION, dtype_node, "Expression is not a data type.");
		ctx->ret = false;
		return;
	}

	sym->dtype = DTypeCopy(dv.value.dtype);
	DValueClear(&dv);
}

static void ResolveGlobals(const String *key, void *value, void *vctx)
{
	Namespace *ns = value;
	ResolveGlobalsCtx *ctx = vctx;

	MapIterate(&ns->children, ResolveGlobals, vctx);

	MapIterate(&ns->public_syms.exports, ResolveSymbols, ctx);
	MapIterate(&ns->private_syms.exports, ResolveSymbols, ctx);
}


bool PassResolveGlobals(CompileState *state)
{
	ResolveGlobalsCtx ctx;
	ctx.ret = true;
	ctx.state = state;

	MapIterate(&state->root_namespace.children, ResolveGlobals, &ctx);

	return ctx.ret;
}

