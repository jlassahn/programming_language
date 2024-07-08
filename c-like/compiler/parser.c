
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/exit_codes.h"
#include "compiler/parser_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>


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
extern int yydebug;

void ParseSetDebug(bool on)
{
	yydebug = (int)on;
}

ParserNode *ParseFile(ParserFile *file, ParserContext *context)
{
	bison_connector.file = file;
	TokenizerStart(&bison_connector.tokenizer, file);

	file->parser_result = yyparse();

	// for Bison, the last node allocated is always the top level symbol
	return GetLastNode();
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

		// FIXME find a better way than using flags for this
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

