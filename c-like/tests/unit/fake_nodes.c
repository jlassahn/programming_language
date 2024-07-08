
#include "tests/unit/unit_test.h"
#include "tests/unit/fake_nodes.h"
#include "compiler/parser_file.h"
#include "compiler/types.h"
#include <string.h>

static StringBuffer *value_buffer = NULL;
static ParserFile value_file =
{
	.filename = "value_buffer",
	.data = NULL,
	.length = 0,
};

static ParserNode *node_stack[100]; // is this enough?
static int stack_top = 0;

ParserNode *MakeNodeFakeValue(ParserSymbol *symbol, const char *value)
{
	if (value_buffer == NULL)
		value_buffer = StringBufferCreateEmpty(0);

	int start = value_buffer->string.length;
	value_buffer = StringBufferAppendChars(value_buffer, value);
	int end = value_buffer->string.length;

	value_file.data = value_buffer->buffer;
	value_file.length = value_buffer->string.length;

	ParserNode *node = MakeNode(symbol, 0, NULL);

	node->position.file = &value_file;
	node->position.start.offset = start;
	node->position.end.offset = end;
	return node;
}

int PushNodeStack(ParserNode *node)
{
	node_stack[stack_top] = node;
	stack_top ++;
	return stack_top;
}

int MakeNodeOnStack(ParserSymbol *symbol, int count)
{
	stack_top -= count;
	CHECK(stack_top >= 0);
	ParserNode *node = MakeNode(symbol, count, &node_stack[stack_top]);
	node_stack[stack_top] = node;
	stack_top ++;
	return stack_top;
}

ParserNode *GetNodeStackTop(void)
{
	stack_top --;
	return node_stack[stack_top];
}

void FreeFakeNodeValues(void)
{
	CHECK(stack_top == 0);
	if (value_buffer != NULL)
		StringBufferFree(value_buffer);
	value_buffer = NULL;
}

