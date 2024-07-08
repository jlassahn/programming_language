
#include "compiler/parser_node.h"
#include "compiler/exit_codes.h"
#include "compiler/memory.h"

#include <stdlib.h>
#include <string.h>

// hacky global state for finding the top of the parse tree in Bison
static ParserNode *last_node = NULL;

static int allocated_nodes = 0;

ParserNode *GetLastNode(void)
{
	return last_node;
}

ParserNode *MakeNode(ParserSymbol *kind, int count, ParserNode **params)
{
	if (count > MAX_CHILDREN)
	{
		fprintf(stderr, "Internal parser error: too many children %d\n", count);
		exit(EXIT_SOFTWARE);
	}

	ParserNode *node = Alloc(sizeof(ParserNode));
	memset(node, 0, sizeof(ParserNode));

	node->symbol = kind;

	int child_count = 0;
	for (int i=0; i<count; i++)
	{
		if (params[i]->position.file)
		{
			if (node->position.file == NULL)
			{
				node->position.file = params[i]->position.file;
				node->position.start = params[i]->position.start;
			}
			node->position.end = params[i]->position.end;
		}

		if (params[i]->symbol->flags & SYM_DISCARD)
		{
			FreeNode(params[i]);
		}
		else
		{
			node->children[child_count] = params[i];
			child_count ++;
		}
	}
	node->count = child_count;
	last_node = node; // FIXME hack for tracking Bison results

	allocated_nodes ++;
	return node;
}

void FreeNode(ParserNode *node)
{
	for(int i=0; i<node->count; i++)
		FreeNode(node->children[i]);
	Free(node);
	allocated_nodes --;
}

bool ParserNodeGetValue(ParserNode *node, String *name_out)
{
	if (node->position.file == NULL)
		return false;

	ParserFile *file = node->position.file;
	uint64_t start = node->position.start.offset;
	uint64_t end = node->position.end.offset;

	int length = (int)(end - start);
	if (length <= 0)
		return false;

	name_out->data = file->data + start;
	name_out->length = length;
	return true;
}

int GetNodeCount(void)
{
	return allocated_nodes;
}

typedef struct Indent Indent;
struct Indent
{
	int node_count;
	Indent *parent;
};

static void PrintIndent(Indent *indent, bool top)
{
	if (!indent)
		return;

	PrintIndent(indent->parent, false);
	if (indent->node_count > 0)
	{
		if (top)
			printf("+-");
		else
			printf("| ");
	}
	else
	{
		printf("  ");
	}
}

static void PrintNodeTreeDepth(FILE *fp, ParserNode *node, Indent *parent)
{
	PrintIndent(parent, true);
	if (node == NULL)
	{
		fprintf(fp, "(null)\n");
		return;
	}

	ParserFile *file = node->position.file;
	uint64_t start = node->position.start.offset;
	uint64_t end = node->position.end.offset;
	fprintf(fp, "%s", node->symbol->rule_name);
	if (file && (node->symbol->flags & PRINT_CONTENT))
		fprintf(fp," [%.*s]", (int)(end - start), &file->data[start]);
	if (node->count > 0)
		fprintf(fp, ":");
	fprintf(fp, "\n");

	Indent indent;
	indent.parent = parent;
	indent.node_count = node->count;
	if (parent && (parent->node_count > 0))
		parent->node_count --;

	for (int i=0; i<node->count; i++)
		PrintNodeTreeDepth(fp, node->children[i], &indent);
}

void PrintNodeTree(FILE *fp, ParserNode *node)
{
	PrintNodeTreeDepth(fp, node, NULL);
}

