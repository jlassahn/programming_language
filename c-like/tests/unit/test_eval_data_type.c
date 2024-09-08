
#include "tests/unit/unit_test.h"
#include "tests/unit/fake_nodes.h"
#include "compiler/parser_symbols.h"
#include "compiler/symbol_table.h"
#include "compiler/builtins.h"
#include "compiler/eval.h"
#include <stdbool.h>
#include <string.h>

void TestEvalDataType(void)
{
	SymbolTable syms;
	EvalContext ctx;
	DataValue dv;

	SymbolTableInit(&syms);
	InitBuiltins(&syms.root, 64);

	ctx.symbol_table = &syms;

	memset(&dv, 0, sizeof(DataValue));

	ParserNode *node = MakeNodeFakeValue(&SYM_IDENTIFIER, "int32");

	CHECK(EvalExpression(&dv, node, &ctx));
	CHECK(dv.value_type == VTYPE_DTYPE);
	CHECK(dv.value.dtype->base_type == DTYPE_INT32);

	DValueClear(&dv);

	FreeNode(node);
	FreeFakeNodeValues();

	FreeBuiltins(&syms.root);
	SymbolTableDestroy(&syms);
}

