
#include "compiler/parser_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

ParserSymbol SYM_UNDEF = { "UNDEFINED" };
ParserSymbol SYM_IDENTIFIER = { "IDENTIFIER" };
ParserSymbol SYM_NUMBER = { "NUMBER" };
ParserSymbol SYM_CHARCONST = { "CHARCONST" };
ParserSymbol SYM_STRINGCONST = { "STRINGCONST" };
ParserSymbol SYM_PUNCTUATION = { "PUNCTUATION" };
ParserSymbol SYM_KEYWORD = { "KEYWORD" };
ParserSymbol SYM_OPERATOR = { "OPERATOR" };

ParserSymbol SYM_EMPTY = { "EMPTY" };
ParserSymbol SYM_FIXME = { "FIXME" };

// hacky global state for finding the top of the parse tree in Bison
static ParserNode *last_node = NULL;

// FIXME strategy for freeing allocated nodes?
ParserNode *MakeNode(ParserSymbol *kind, int count, ParserNode **params)
{
	if (count > MAX_CHILDREN)
	{
		printf("ERROR: too many children %d\n", count);
		exit(-1);
	}

	ParserNode *node = malloc(sizeof(ParserNode));
	memset(node, 0, sizeof(ParserNode));

	node->symbol = kind;
	printf("%s <-", kind->rule_name);

	node->count = count;
	for (int i=0; i<count; i++)
	{
		node->children[i] = params[i];
		if (params[i])
			printf(" %s", params[i]->symbol->rule_name);
		else
			printf(" NULL");
	}
	if (count == 0)
	{
		printf(" (empty)");
	}
	else
	{
		node->position.file = node->children[0]->position.file;
		node->position.start = node->children[0]->position.start;
		node->position.end = node->children[count-1]->position.end;
	}
	printf("\n");

	last_node = node;
	return node;
}

static void PrintNodeTreeDepth(FILE *fp, ParserNode *node, int depth)
{
	for (int i=0; i<depth; i++)
		fprintf(fp, "  ");
	if (node == NULL)
	{
		fprintf(fp, "(null)\n");
		return;
	}

	ParserFile *file = node->position.file;
	uint64_t start = node->position.start.offset;
	uint64_t end = node->position.end.offset;
	fprintf(fp, "%s", node->symbol->rule_name);
	if (file)
		fprintf(fp," [%.*s]", (int)(end - start), &file->data[start]);
	if (node->count > 0)
		fprintf(fp, ":");
	fprintf(fp, "\n");

	for (int i=0; i<node->count; i++)
		PrintNodeTreeDepth(fp, node->children[i], depth + 1);
}

void PrintNodeTree(FILE *fp, ParserNode *node)
{
	PrintNodeTreeDepth(fp, node, 0);
}

// Bison parser interface
int yyparse (void);

typedef struct BisonConnector BisonConnector;
struct BisonConnector
{
	Tokenizer tokenizer;
	ParserFile *file;
};

static BisonConnector bison_connector;
ParserNode * yylval = NULL;

ParserNode *ParseFile(ParserFile *file, ParserContext *context)
{
	bison_connector.file = file;
	TokenizerStart(&bison_connector.tokenizer, file);

	int ret = yyparse();
	printf("parser return = %d\n", ret);

	// for Bison, the last node allocated is always the top level symbol
	return last_node;
}

int yylex(void)
{
	if (TokenizerIsEOF(&bison_connector.tokenizer))
	{
		yylval = NULL;
		return 0;
	}
	else
	{
		Token token;
		GetCurrentToken(&bison_connector.tokenizer, &token);
		TokenizerConsume(&bison_connector.tokenizer);

		// FIXME fiind a better way than using flags for this
		//       maybe token types have a ParserSymbol* member?
		int token_sym = token.token_type->flags & 15;
		static ParserSymbol *syms[] =
		{
			&SYM_UNDEF,
			&SYM_IDENTIFIER,
			&SYM_NUMBER,
			&SYM_CHARCONST,
			&SYM_STRINGCONST,
			&SYM_PUNCTUATION,
			&SYM_KEYWORD,
			&SYM_OPERATOR
		};

		ParserNode *node = MakeNode(syms[token_sym], 0, NULL);
		node->position = token.position;
		yylval = node;

		int id = token.token_type->id;

		return id;
	}
}

void yyerror(const char *s)
{
	// FIXME hook up error handling system
	printf("ERROR: %s\n", s);
}


// writing a recursive descent parser by hand would look like...
#if 0
/*
STATEMENT:
	STRUCT_DECLARATION    // struct ...
	UNION_DECLARATION     // union ...
	ENUM_DECLARATION      // enum ...
	COMPOUND_STATEMENT    // { ...
	FOR_STATEMENT
	WHILE_STATEMENT
	DO_STATEMENT
	IF_STATEMENT
	SWITCH_STATEMENT
	WHILE_STATEMENT
	BREAK_STATEMENT
	CONTINUE_STATEMENT
	GOTO_STATEMENT
	RETURN_STATEMENT
	;
	EXPRESSION_STATEMENT  // EXPRESSION ;
	LABELED_STATEMENT     // EXPRESSION : ...
	SYMBOL_DEFINITION     // EXPRESSION IDENTIFIER ...
	                      // STORAGE_SPECIFIER ...
*/
ParserNode *ParseStatement(ParserContext *context)
{
	ParserNode *ret = NULL;

	ret = ParseStructDeclaration(context);
	if (!NoMatch(ret))
		return ret;

	ret = ParseUnionDeclaration(context);
	if (!NoMatch(ret))
		return ret;

	// ...

	ParserNode *expression_lookahead = ParseExpression();

	ret = ParseSymbolDefinitionL1(context, expression_lookahead);
	if (!NoMatch(ret))
		return ret;

	ret = ParseLabeledStatementL1(context, expression_lookahead);
	if (!NoMatch(ret))
		return ret;

	ret = ParseExpressionStatementL1(context, expression_lookahead);
	if (!NoMatch(ret))
		return ret;
	return ErrorNode(context, "Expected semicolon after expression");
}
#endif

