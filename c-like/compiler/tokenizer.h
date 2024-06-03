
#ifndef INCLUDED_TOKENIZER_H
#define INCLUDED_TOKENIZER_H

#include "parser_file.h"

#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>

typedef struct TokenType TokenType;
typedef struct Token Token;
typedef struct Tokenizer Tokenizer;

/* FIXME what information to carry with tokens?
	one TokenType per operator, keyword or punctuator
		TokenType.name is the string value of the token in these cases
	A TokenType for NUMBER, one for STRING, etc.
		TokenType.name is "PPNUMBER", etc
		the value of these can be computed from the contents.
		Don't precompute values during parse.
	For YACC TokenType also needs to be assigned the YACC token class.
	Maybe set some flags for whether something's got a value, is an operator...
	
	How to convert a token to a ParseNode with a semantic action?
		cannot put an action on the TokenType because some tokens have
		different behavior at different times.  (e.g. unary vs binary -)
*/

struct TokenType
{
	const char *name;
	uint32_t flags;
	uint32_t id; // this holds the YACC token class when parsing with YACC
};

struct Token
{
	const TokenType *token_type;
	FilePositionRange position;
};

struct Tokenizer
{
	ParserFile *file;
	Token current_token;
};


void TokenizerStart(Tokenizer *tokenizer, ParserFile* file);
void TokenizerConsume(Tokenizer *tokenizer);
void GetCurrentToken(Tokenizer *tokenizer, Token *token_out);
bool TokenizerIsEOF(Tokenizer *tokenizer);

void TokenPrint(FILE *fp, Token *token);

const TokenType TOKEN_EOF;
const TokenType TOKEN_PPIDENTIFIER;
const TokenType TOKEN_PPNUMBER;
const TokenType TOKEN_CHARCONST;
const TokenType TOKEN_STRINGCONST;
const TokenType TOKEN_OPERATOR;
const TokenType TOKEN_OTHER;

#endif

