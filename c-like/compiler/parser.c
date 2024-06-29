
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/exit_codes.h"
#include "compiler/parser_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

ParserSymbol SYM_UNDEF = { "UNDEFINED", 0 };
ParserSymbol SYM_IDENTIFIER = { "IDENTIFIER", PRINT_CONTENT };
ParserSymbol SYM_NUMBER = { "NUMBER", PRINT_CONTENT };
ParserSymbol SYM_CHARCONST = { "CHARCONST", PRINT_CONTENT };
ParserSymbol SYM_STRINGCONST = { "STRINGCONST", PRINT_CONTENT };
ParserSymbol SYM_PUNCTUATION = { "PUNCTUATION", PRINT_CONTENT | SYM_DISCARD };
ParserSymbol SYM_KEYWORD = { "KEYWORD", PRINT_CONTENT | SYM_DISCARD };
ParserSymbol SYM_OPERATOR = { "OPERATOR", PRINT_CONTENT | SYM_DISCARD };

ParserSymbol SYM_FIXME = { "FIXME", 0 };

ParserSymbol SYM_EMPTY = { "EMPTY", 0 };
ParserSymbol SYM_LIST = { "LIST", 0 };
ParserSymbol SYM_DOT_OP = { "DOT_OP", 0 };
ParserSymbol SYM_IMPORT = { "IMPORT", 0 };
ParserSymbol SYM_IMPORT_PRIVATE = { "IMPORT_PRIVATE", 0 };
ParserSymbol SYM_PROTOTYPE = { "PROTOTYPE", 0 };
ParserSymbol SYM_FUNC = { "FUNC", 0 };
ParserSymbol SYM_DECLARATION = { "DECLARATION", 0 };
ParserSymbol SYM_USING = { "USING", 0 };
ParserSymbol SYM_USING_AS = { "USING_AS", 0 };
ParserSymbol SYM_TRAILING_COMMA = { "TRAILING_COMMA", 0 };
ParserSymbol SYM_ELLIPSIS = { "ELLIPSIS", 0 };
ParserSymbol SYM_PARAMETER = { "PARAMETER", 0 };
ParserSymbol SYM_PARAM_TYPE = { "PARAM_TYPE", 0 };
ParserSymbol SYM_DECL_TYPE = { "DECL_TYPE", 0 };
ParserSymbol SYM_INITIALIZE = { "INITIALIZE", 0 };
ParserSymbol SYM_INIT_STRUCT = { "INIT_STRUCT", 0 };
ParserSymbol SYM_INIT_ARRAY = { "INIT_ARRAY", 0 };
ParserSymbol SYM_STATEMENT_LIST = { "STATEMENT_LIST", 0 };
ParserSymbol SYM_STRUCT_LIST = { "STRUCT_LIST", 0 };
ParserSymbol SYM_ENUM_LIST = { "ENUM_LIST", 0 };
ParserSymbol SYM_ENUM_ELEMENT = { "ENUM_ELEMENT", 0 };
ParserSymbol SYM_STRUCT_DEC = { "STRUCT_DEC", 0 };
ParserSymbol SYM_STRUCT_DEF = { "STRUCT_DEF", 0 };
ParserSymbol SYM_UNION_DEC = { "UNION_DEC", 0 };
ParserSymbol SYM_UNION_DEF = { "UNION_DEF", 0 };
ParserSymbol SYM_ENUM_DEC = { "ENUM_DEC", 0 };
ParserSymbol SYM_ENUM_DEF = { "ENUM_DEF", 0 };
ParserSymbol SYM_EXPRESSION_STATEMENT = { "EXPRESSION_STATEMENT", 0 };
ParserSymbol SYM_EMPTY_STATEMENT = { "EMPTY_STATEMENT", 0 };
ParserSymbol SYM_LABEL_STATEMENT = { "LABEL_STATEMENT", 0 };
ParserSymbol SYM_FOR_STATEMENT = { "FOR_STATEMENT", 0 };
ParserSymbol SYM_WHILE_STATEMENT = { "WHILE_STATEMENT", 0 };
ParserSymbol SYM_DO_STATEMENT = { "DO_STATEMENT", 0 };
ParserSymbol SYM_IF_STATEMENT = { "IF_STATEMENT", 0 };
ParserSymbol SYM_IF_ELSE = { "IF_ELSE", 0 };
ParserSymbol SYM_SWITCH_STATEMENT = { "SWITCH_STATEMENT", 0 };
ParserSymbol SYM_BREAK_STATEMENT = { "BREAK_STATEMENT", 0 };
ParserSymbol SYM_CONTINUE_STATEMENT = { "CONTINUE_STATEMENT", 0 };
ParserSymbol SYM_GOTO_STATEMENT = { "GOTO_STATEMENT", 0 };
ParserSymbol SYM_RETURN_STATEMENT = { "RETURN_STATEMENT", 0 };
ParserSymbol SYM_RETURN_VOID = { "RETURN_VOID", 0 };
ParserSymbol SYM_CASE_ELEMENT = { "CASE_ELEMENT", 0 };
ParserSymbol SYM_CASE_END_ELEMENT = { "CASE_END_ELEMENT", 0 };
ParserSymbol SYM_DEFAULT_LABEL = { "DEFAULT_LABEL", 0 };
ParserSymbol SYM_CASE_LABEL = { "CASE_LABEL", 0 };
ParserSymbol SYM_CONSTANT = { "CONSTANT", 0 };
ParserSymbol SYM_ASSIGN_OP = { "ASSIGN_OP", 0 };
ParserSymbol SYM_ASSIGN_MULT_OP = { "ASSIGN_MULT_OP", 0 };
ParserSymbol SYM_ASSIGN_DIV_OP = { "ASSIGN_DIV_OP", 0 };
ParserSymbol SYM_ASSIGN_MOD_OP = { "ASSIGN_MOD_OP", 0 };
ParserSymbol SYM_ASSIGN_ADD_OP = { "ASSIGN_ADD_OP", 0 };
ParserSymbol SYM_ASSIGN_SUB_OP = { "ASSIGN_SUB_OP", 0 };
ParserSymbol SYM_ASSIGN_SHR_OP = { "ASSIGN_SHR_OP", 0 };
ParserSymbol SYM_ASSIGN_SHL_OP = { "ASSIGN_SHL_OP", 0 };
ParserSymbol SYM_ASSIGN_AND_OP = { "ASSIGN_AND_OP", 0 };
ParserSymbol SYM_ASSIGN_OR_OP = { "ASSIGN_OR_OP", 0 };
ParserSymbol SYM_ASSIGN_XOR_OP = { "ASSIGN_XOR_OP", 0 };
ParserSymbol SYM_CONDITIONAL = { "CONDITIONAL", 0 };
ParserSymbol SYM_LOG_OR_OP = { "LOG_OR_OP", 0 };
ParserSymbol SYM_LOG_AND_OP = { "LOG_AND_OP", 0 };
ParserSymbol SYM_OR_OP = { "OR_OP", 0 };
ParserSymbol SYM_AND_OP = { "AND_OP", 0 };
ParserSymbol SYM_ADDR_OP = { "ADDR_OP", 0 };
ParserSymbol SYM_XOR_OP = { "XOR_OP", 0 };
ParserSymbol SYM_EQUAL_OP = { "EQUAL_OP", 0 };
ParserSymbol SYM_NEQUAL_OP = { "NEQUAL_OP", 0 };
ParserSymbol SYM_LESS_OP = { "LESS_OP", 0 };
ParserSymbol SYM_GREATER_OP = { "GREATER_OP", 0 };
ParserSymbol SYM_LESSEQ_OP = { "LESSEQ_OP", 0 };
ParserSymbol SYM_GREATEREQ_OP = { "GREATER_EQ_OP", 0 };
ParserSymbol SYM_SHL_OP = { "SHL_OP", 0 };
ParserSymbol SYM_SHR_OP = { "SHR_OP", 0 };
ParserSymbol SYM_ADD_OP = { "ADD_OP", 0 };
ParserSymbol SYM_SUB_OP = { "SUB_OP", 0 };
ParserSymbol SYM_DIV_OP = { "DIV_OP", 0 };
ParserSymbol SYM_MOD_OP = { "MOD_OP", 0 };
ParserSymbol SYM_MULT_OP = { "MULT_OP", 0 };
ParserSymbol SYM_PTR_OP = { "PTR_OP", 0 };
ParserSymbol SYM_NOT_OP = { "NOT_OP", 0 };
ParserSymbol SYM_BITNOT_OP = { "BITNOT_OP", 0 };
ParserSymbol SYM_PREINC_OP = { "PREINC_OP", 0 };
ParserSymbol SYM_PREDEC_OP = { "PREDEC_OP", 0 };
ParserSymbol SYM_POSTINC_OP = { "POSTINC_OP", 0 };
ParserSymbol SYM_POSTDEC_OP = { "POSTDEC_OP", 0 };
ParserSymbol SYM_NEG_OP = { "NEG_OP", 0 };
ParserSymbol SYM_POS_OP = { "POS_OP", 0 };
ParserSymbol SYM_SIZEOF_OP = { "SIZEOF_OP", 0 };
ParserSymbol SYM_ARRAY_OP = { "ARRAY_OP", 0 };
ParserSymbol SYM_CALL_OP = { "CALL_OP", 0 };
ParserSymbol SYM_INIT_OP = { "INIT_OP", 0 };
ParserSymbol SYM_PAREN_EXPRESSION = { "PAREN_EXPRESSION", 0 };
ParserSymbol SYM_STRING = { "STRING", 0 };
ParserSymbol SYM_TYPE_EXPRESSION = { "TYPE_EXPRESSION", 0 };
ParserSymbol SYM_TYPE_ARRAY = { "TYPE_ARRAY", 0 };
ParserSymbol SYM_TYPE_ARRAY_MATCH = { "TYPE_ARRAY_MATCH", 0 };
ParserSymbol SYM_TYPE_BITFIELD = { "TYPE_BITFIELD", 0 };
ParserSymbol SYM_TYPE_LINKAGE = { "TYPE_LINKAGE", 0 };
ParserSymbol SYM_TYPE_LINKNAME = { "TYPE_LINKNAME", 0 };

// hacky global state for finding the top of the parse tree in Bison
static ParserNode *last_node = NULL;

static int allocated_nodes = 0;

// FIXME strategy for freeing allocated nodes?
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
	if (count != 0)
	{
		node->position.file = params[0]->position.file;
		node->position.start = params[0]->position.start;
		node->position.end = params[count-1]->position.end;
	}

	int child_count = 0;
	for (int i=0; i<count; i++)
	{
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
	printf("ERROR: %s yylval = %p\n", s, yylval);
	if (yylval)
	{
		printf("error at [%s:%ld:%ld]%s\n",
				yylval->position.file->filename,
				yylval->position.start.line+1,
				yylval->position.start.byte_in_line+1,
				yylval->symbol->rule_name);
	}
	else
	{
		printf("error at EOF\n");
	}
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

