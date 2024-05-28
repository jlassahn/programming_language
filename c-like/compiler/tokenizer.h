
#ifndef INCLUDED_TOKENIZER_H
#define INCLUDED_TOKENIZER_H

#include "parser_file.h"

#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>

typedef struct TokenType TokenType;
typedef struct Token Token;
typedef struct Tokenizer Tokenizer;

struct TokenType
{
	const char *name;
	uint32_t flags;
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

