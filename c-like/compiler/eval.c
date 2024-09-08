
#include "compiler/eval.h"
#include "compiler/parser_symbols.h"
#include "compiler/errors.h"


bool EvalExpression(DataValue *dv_out, ParserNode *node, EvalContext *ctx)
{
	String name;
	if (node->symbol == &SYM_IDENTIFIER)
	{
		if (!ParserNodeGetValue(node, &name))
			return false; // FIXME maybe can't happen?

		Symbol *sym = SymbolTableFind(ctx->symbol_table, &name);
		if (sym == NULL)
		{
			ErrorAtNode(ERROR_DEFINITION, node, "Undefined symbol: %.*s", name.length, name.data);
			return false;
		}

		// FIXME check if symbol is fully resolved, and recurse if not.

		DataValue *dv = SymbolGetDValue(sym);
		DValueCopy(dv_out, dv);
		return true;
	}
	ErrorAtNode(ERROR_INTERNAL, node, "Unsupported node type: %s\n", node->symbol->rule_name);
	return false;
}

