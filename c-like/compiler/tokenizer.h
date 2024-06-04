
#ifndef INCLUDED_TOKENIZER_H
#define INCLUDED_TOKENIZER_H

#include "compiler/parser_file.h"

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

extern const TokenType TOKEN_EOF;
extern const TokenType TOKEN_IDENTIFIER;
extern const TokenType TOKEN_NUMBER;
extern const TokenType TOKEN_CHARCONST;
extern const TokenType TOKEN_STRINGCONST;

extern const TokenType TOKEN_UNKNOWN;

// extern const TokenType TOKEN_PUNCTUATION;

// extern const TokenType TOKEN_OPERATOR; // FIXME remove

// keywords
extern const TokenType TOKEN_ALIAS;
extern const TokenType TOKEN_ALLIGNMENT; // FIXME maybe not real
extern const TokenType TOKEN_AS;
extern const TokenType TOKEN_ARRAY;
extern const TokenType TOKEN_AUTO;
extern const TokenType TOKEN_BITFIELD;
extern const TokenType TOKEN_BREAK;
extern const TokenType TOKEN_CASE;
extern const TokenType TOKEN_CONSTANT;
extern const TokenType TOKEN_CONTINUE;
extern const TokenType TOKEN_DEFAULT;
extern const TokenType TOKEN_DO;
extern const TokenType TOKEN_ELSE;
extern const TokenType TOKEN_ENUM;
extern const TokenType TOKEN_FOR;
extern const TokenType TOKEN_GOTO;
extern const TokenType TOKEN_IF;
extern const TokenType TOKEN_IMPORT;
extern const TokenType TOKEN_INLINE;
extern const TokenType TOKEN_LINKAGE;
extern const TokenType TOKEN_LINKNAME;
extern const TokenType TOKEN_POINTER;
extern const TokenType TOKEN_PRIVATE;
extern const TokenType TOKEN_READONLY;
extern const TokenType TOKEN_REGISTER;
extern const TokenType TOKEN_RESTRICT;
extern const TokenType TOKEN_RETURN;
extern const TokenType TOKEN_SIZEOF;
extern const TokenType TOKEN_STATIC;
extern const TokenType TOKEN_STRUCT;
extern const TokenType TOKEN_SWITCH;
extern const TokenType TOKEN_TYPEDEF;
extern const TokenType TOKEN_UNION;
extern const TokenType TOKEN_USING;
extern const TokenType TOKEN_VOLATILE;
extern const TokenType TOKEN_WHILE;
// reserve TEMPLATE and CLASS for future...

// operators
extern const TokenType TOKEN_ASSIGN_OP;      /* = */
extern const TokenType TOKEN_ASSIGN_MULT_OP; /* *= */
extern const TokenType TOKEN_ASSIGN_DIV_OP;  /* /= */
extern const TokenType TOKEN_ASSIGN_MOD_OP;  /* %= */
extern const TokenType TOKEN_ASSIGN_ADD_OP;  /* += */
extern const TokenType TOKEN_ASSIGN_SUB_OP;  /* -= */
extern const TokenType TOKEN_ASSIGN_SHR_OP;  /* >>= */
extern const TokenType TOKEN_ASSIGN_SHL_OP;  /* <<= */
extern const TokenType TOKEN_ASSIGN_AND_OP;  /* &= */
extern const TokenType TOKEN_ASSIGN_OR_OP;   /* |= */
extern const TokenType TOKEN_ASSIGN_XOR_OP;  /* ^= */
extern const TokenType TOKEN_LOG_OR_OP;      /* || */
extern const TokenType TOKEN_LOG_AND_OP;     /* && */
extern const TokenType TOKEN_OR_OP;          /* | */
extern const TokenType TOKEN_AND_ADDR_OP;    /* & */
extern const TokenType TOKEN_XOR_OP;         /* ^ */
extern const TokenType TOKEN_EQUAL_OP;       /* == */
extern const TokenType TOKEN_NEQUAL_OP;      /* != */
extern const TokenType TOKEN_LESS_OP;        /* < */
extern const TokenType TOKEN_GREATER_OP;     /* > */
extern const TokenType TOKEN_LESSEQ_OP;      /* <= */
extern const TokenType TOKEN_GREATEREQ_OP;   /* >= */
extern const TokenType TOKEN_SHL_OP;         /* << */
extern const TokenType TOKEN_SHR_OP;         /* >> */
extern const TokenType TOKEN_ADD_OP;         /* + */
extern const TokenType TOKEN_SUB_OP;         /* - */
extern const TokenType TOKEN_DIV_OP;         /* / */
extern const TokenType TOKEN_MOD_OP;         /* % */
extern const TokenType TOKEN_MULT_PTR_OP;    /* * */
extern const TokenType TOKEN_NOT_OP;         /* ! */
extern const TokenType TOKEN_BITNOT_OP;      /* ~ */
extern const TokenType TOKEN_INC_OP;         /* ++ */
extern const TokenType TOKEN_DEC_OP;         /* -- */

extern const TokenType TOKEN_ELIPSIS;

extern const TokenType TOKEN_SEMICOLON;
extern const TokenType TOKEN_LCURLY;
extern const TokenType TOKEN_RCURLY;
extern const TokenType TOKEN_COMMA;
extern const TokenType TOKEN_DOT;
extern const TokenType TOKEN_LSQUARE;
extern const TokenType TOKEN_RSQUARE;
extern const TokenType TOKEN_LPAREN;
extern const TokenType TOKEN_RPAREN;
extern const TokenType TOKEN_COLON;
extern const TokenType TOKEN_QUESTION;

#endif

