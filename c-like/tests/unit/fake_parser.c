
#include "compiler/parser_node.h"
#include "compiler/parser.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include <string.h>

typedef struct FakeParse FakeParse;
struct FakeParse
{
	const char *path;
	ParserNode *root;
	int error_code;
};

static List parses;

void FakeParserSet(const char *path, ParserNode *root, int error_code)
{
	FakeParse *parse = Alloc(sizeof(FakeParse));
	parse->path = path;
	parse->root = root;
	parse->error_code = error_code;

	ListInsertFirst(&parses, parse);
}

void FakeParserFree(void)
{
	while (true)
	{
		FakeParse *parse = ListRemoveFirst(&parses);
		if (parse == NULL)
			break;
		if (parse->root)
			FreeNode(parse->root);
		Free(parse);
	}
}

ParserNode *ParseFile(ParserFile *file, ParserContext *context)
{
	for (ListEntry *entry=parses.first; entry!=NULL; entry=entry->next)
	{
		FakeParse *parse = entry->item;
		if (strcmp(parse->path, file->filename) == 0)
		{
			file->parser_result = parse->error_code;
			ParserNode *root = parse->root;
			// don't keep the nodes around, becuase we don't own them anymore.
			parse->root = NULL;
			parse->error_code = -1;
			return root;
		}
	}
	file->parser_result = 0;
	return NULL;
}


