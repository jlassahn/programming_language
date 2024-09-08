
#ifndef INCLUDED_EVAL_H
#define INCLUDED_EVAL_H

#include "compiler/symbol_table.h"
#include "compiler/parser_node.h"

typedef struct EvalContext EvalContext;
struct EvalContext
{
	SymbolTable *symbol_table;
};

bool EvalExpression(DataValue *dv_out, ParserNode *node, EvalContext *ctx);

#endif

