
#include "compiler/errors.h"
#include "compiler/tokenizer.h"
#include "compiler/parser_file.h"
#include <string.h>


typedef struct TokenInfo TokenInfo;
struct TokenInfo
{
	const char *text;
	const TokenType *token_type;
};

// shorter tokens must follow longer tokens which have the same prefix
TokenInfo operator_list[] =
{
	{"...", &TOKEN_ELLIPSIS },
	{">>=", &TOKEN_ASSIGN_SHR_OP },
	{"<<=", &TOKEN_ASSIGN_SHL_OP },
	{"*=", &TOKEN_ASSIGN_MULT_OP },
	{"/=", &TOKEN_ASSIGN_DIV_OP },
	{"%=", &TOKEN_ASSIGN_MOD_OP },
	{"+=", &TOKEN_ASSIGN_ADD_OP },
	{"-=", &TOKEN_ASSIGN_SUB_OP },
	{"&=", &TOKEN_ASSIGN_AND_OP },
	{"|=", &TOKEN_ASSIGN_OR_OP },
	{"^=", &TOKEN_ASSIGN_XOR_OP },
	{"||", &TOKEN_LOG_OR_OP },
	{"&&", &TOKEN_LOG_AND_OP },
	{"==", &TOKEN_EQUAL_OP },
	{"!=", &TOKEN_NEQUAL_OP },
	{"<=", &TOKEN_LESSEQ_OP },
	{">=", &TOKEN_GREATEREQ_OP },
	{"<<", &TOKEN_SHL_OP },
	{">>", &TOKEN_SHR_OP },
	{"++", &TOKEN_INC_OP },
	{"--", &TOKEN_DEC_OP },
	{"=", &TOKEN_ASSIGN_OP },
	{"|", &TOKEN_OR_OP },
	{"&", &TOKEN_AND_ADDR_OP },
	{"^", &TOKEN_XOR_OP },
	{"<", &TOKEN_LESS_OP },
	{">", &TOKEN_GREATER_OP },
	{"+", &TOKEN_ADD_OP },
	{"-", &TOKEN_SUB_OP },
	{"/", &TOKEN_DIV_OP },
	{"%", &TOKEN_MOD_OP },
	{"*", &TOKEN_MULT_PTR_OP },
	{"!", &TOKEN_NOT_OP },
	{"~", &TOKEN_BITNOT_OP },

	{";", &TOKEN_SEMICOLON },
	{"{", &TOKEN_LCURLY },
	{"}", &TOKEN_RCURLY },
	{",", &TOKEN_COMMA },
	{".", &TOKEN_DOT },
	{"[", &TOKEN_LSQUARE },
	{"]", &TOKEN_RSQUARE },
	{"(", &TOKEN_LPAREN },
	{")", &TOKEN_RPAREN },
	{":", &TOKEN_COLON },
	{"?", &TOKEN_QUESTION },

	{ NULL, NULL }
};

TokenInfo keyword_list[] =
{
	{ "alias", &TOKEN_ALIAS },
	{ "allignment", &TOKEN_ALLIGNMENT }, // FIXME maybe not real
	{ "as", &TOKEN_AS },
	{ "array", &TOKEN_ARRAY },
	{ "auto", &TOKEN_AUTO },
	{ "bitfield", &TOKEN_BITFIELD },
	{ "break", &TOKEN_BREAK },
	{ "case", &TOKEN_CASE },
	{ "constant", &TOKEN_CONSTANT },
	{ "continue", &TOKEN_CONTINUE },
	{ "default", &TOKEN_DEFAULT },
	{ "do", &TOKEN_DO },
	{ "else", &TOKEN_ELSE },
	{ "enum", &TOKEN_ENUM },
	{ "for", &TOKEN_FOR },
	{ "goto", &TOKEN_GOTO },
	{ "if", &TOKEN_IF },
	{ "import", &TOKEN_IMPORT },
	{ "inline", &TOKEN_INLINE },
	{ "linkage", &TOKEN_LINKAGE },
	{ "linkname", &TOKEN_LINKNAME },
	{ "pointer", &TOKEN_POINTER },
	{ "private", &TOKEN_PRIVATE },
	{ "readonly", &TOKEN_READONLY },
	{ "register", &TOKEN_REGISTER },
	{ "restrict", &TOKEN_RESTRICT },
	{ "return", &TOKEN_RETURN },
	{ "sizeof", &TOKEN_SIZEOF },
	{ "static", &TOKEN_STATIC },
	{ "struct", &TOKEN_STRUCT },
	{ "switch", &TOKEN_SWITCH },
	{ "typedef", &TOKEN_TYPEDEF },
	{ "union", &TOKEN_UNION },
	{ "using", &TOKEN_USING },
	{ "volatile", &TOKEN_VOLATILE },
	{ "while", &TOKEN_WHILE },
	{ NULL, NULL }
};

static bool IsSpace(char x);
static bool IsLetter(char x);
static bool IsDigit(char x);
static void SkipSpaceAndComments(Tokenizer *tokenizer);
static bool SeekTo(ParserFile *file, const char *text);
static void ConsumePPIdentifier(Tokenizer *tokenizer);
static void ConsumePPNumber(Tokenizer *tokenizer);
static void ConsumeStringLiteral(Tokenizer *tokenizer);
static void ConsumeCharLiteral(Tokenizer *tokenizer);
static bool MatchAndConsumeList(ParserFile *file, TokenInfo *list, const TokenType **tt_out);
static const TokenType *MatchTokenList(TokenInfo *list, const char *text, int length);

void TokenizerStart(Tokenizer *tokenizer, ParserFile* file)
{
	tokenizer->file = file;

	TokenizerConsume(tokenizer);
}

void TokenizerConsume(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;

	SkipSpaceAndComments(tokenizer);

	// FIXME don't support \U unicode escapes outside of char and string consts.
	// FIXME don't support wide character strings.

	tokenizer->current_token.position.start = file->current_pos;
	tokenizer->current_token.position.file = file;

	const TokenType *tt = NULL;
	const char *cur = FileGet(file);
	if (cur[0] == 0)
	{
		tokenizer->current_token.token_type = &TOKEN_EOF;
	}
	else if (IsLetter(cur[0]))
	{
		ConsumePPIdentifier(tokenizer);

		const char *text = cur;
		int length = file->current_pos.offset -
			tokenizer->current_token.position.start.offset;
		tt = MatchTokenList(keyword_list, text, length);
		if (tt)
			tokenizer->current_token.token_type = tt;
		else
			tokenizer->current_token.token_type = &TOKEN_IDENTIFIER;
	}
	else if (IsDigit(cur[0]) || (cur[0] == '.' && IsDigit(cur[1])))
	{
		tokenizer->current_token.token_type = &TOKEN_NUMBER;
		ConsumePPNumber(tokenizer);
	}
	else if (cur[0] == '\"')
	{
		tokenizer->current_token.token_type = &TOKEN_STRINGCONST;
		ConsumeStringLiteral(tokenizer);
	}
	else if (cur[0] == '\'')
	{
		tokenizer->current_token.token_type = &TOKEN_CHARCONST;
		ConsumeCharLiteral(tokenizer);
	}
	else if (MatchAndConsumeList(file, operator_list, &tt))
	{
		tokenizer->current_token.token_type = tt;
	}
	else
	{
		tokenizer->current_token.token_type = &TOKEN_UNKNOWN;
		FileConsume(file, 1);
	}

	tokenizer->current_token.position.end = file->current_pos;
}

void GetCurrentToken(Tokenizer *tokenizer, Token *token_out)
{
	memcpy(token_out, &tokenizer->current_token, sizeof(Token));
}

bool TokenizerIsEOF(Tokenizer *tokenizer)
{
	return tokenizer->current_token.token_type == &TOKEN_EOF;
}

void TokenPrint(FILE *fp, Token *token)
{
	fprintf(fp, "[%s:%ld:%ld]%s(%d)[%ld:%ld] = [%.*s]\n",
			token->position.file->filename,
			token->position.start.line+1,
			token->position.start.byte_in_line+1,
			token->token_type->name,
			token->token_type->id,
			token->position.end.line+1,
			token->position.end.byte_in_line+1,
			(int)(token->position.end.offset - token->position.start.offset),
			&token->position.file->data[token->position.start.offset]);
}

/*
   Token rules
   white space (comments, space htab, vtab, newline, formfeed)
   identifiers ( _ letters, followed by _ letters digits)
   ppnumbers digit(. digit char e+ e-) or . digit (. digit char e+ e-)
   character constants ' text '  or L' text ' handle \x
   string literals   " text " or L" text " handle \x
   punctuators (up to two chars long)
   other characters
   Unicode escapes can occur anywhere a normal character can, they
   are up to 10 bytes long: \U12345678

   Error tokens for EOF in the middle of a token.
*/

static void SkipSpaceAndComments(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;
	while (true)
	{
		FilePosition pos = file->current_pos;
		if (FileMatchAndConsume(file, "#"))
		{
			if (!SeekTo(file, "\n"))
			{
				ErrorAt(ERROR_PARSER, file->filename, &pos,
						"No end of line after #");
				return;
			}
		}
		else if (FileMatchAndConsume(file, "//"))
		{
			if (!SeekTo(file, "\n"))
			{
				ErrorAt(ERROR_PARSER, file->filename, &pos,
						"No end of line after //");
				return;
			}
		}
		else if (FileMatchAndConsume(file, "/*"))
		{
			if (!SeekTo(file, "*/"))
			{
				ErrorAt(ERROR_PARSER, file->filename, &pos,
						"End of file inside /*");
				return;
			}
		}
		else if (IsSpace(*FileGet(file)))
		{
			FileConsume(file, 1);
		}
		else
		{
			return;
		}
	}
}

static bool SeekTo(ParserFile *file, const char *text)
{
	while (!FileMatchAndConsume(file, text))
	{
		if (*FileGet(file) == 0)
			return false;

		FileConsume(file, 1);
	}
	return true;
}

static bool IsSpace(char x)
{
	if (x == ' ')
		return true;
	if (x == '\t')
		return true;
	if (x == '\f')
		return true;
	if (x == '\n')
		return true;
	if (x == '\v')
		return true;
	return false;
}

static bool IsLetter(char x)
{
	if ((x >= 'a') && (x <= 'z'))
		return true;
	if ((x >= 'A') && (x <= 'Z'))
		return true;
	if (x == '_')
		return true;

	// FIXME what about Unicode in identifiers?
	// if (x >= 128)
	//    return true;
	return false;
}

static bool IsDigit(char x)
{
	if ((x >= '0') && (x <= '9'))
		return true;
	return false;
}

static void ConsumePPIdentifier(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;
	while (IsLetter(*FileGet(file)) || IsDigit(*FileGet(file)))
		FileConsume(file, 1);
}

static void ConsumePPNumber(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;
	bool is_e = false;
	while (true)
	{
		char c = *FileGet(file);

		if (IsDigit(c) || IsLetter(c) || (c == '.') ||
				(is_e && (c == '+')) || (is_e && (c == '-')))
		{
			FileConsume(file, 1);
		}
		else
		{
			return;
		}

		if ((c=='e') || (c=='E') || (c=='p') || (c=='P'))
			is_e = true;
	}
}

static void ConsumeStringLiteral(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;
	FilePosition pos = file->current_pos;
	FileMatchAndConsume(file, "\"");
	while (true)
	{
		char c = *FileGet(file);
		if (c == 0)
		{
			ErrorAt(ERROR_PARSER, file->filename, &pos,
					"End of file inside string.");
			return;
		}
		else if (c == '\n')
		{
			ErrorAt(ERROR_PARSER, file->filename, &pos,
					"End of line inside string.");
			return;
		}
		else if (FileMatch(file, "\\\""))
		{
			FileConsume(file, 2);
		}
		else if (c == '\"')
		{
			FileConsume(file, 1);
			return;
		}
		else
		{
			FileConsume(file, 1);
		}
	}
}

static void ConsumeCharLiteral(Tokenizer *tokenizer)
{
	ParserFile *file = tokenizer->file;
	FilePosition pos = file->current_pos;
	FileMatchAndConsume(file, "\'");
	while (true)
	{
		char c = *FileGet(file);
		if (c == 0)
		{
			ErrorAt(ERROR_PARSER, file->filename, &pos,
					"End of file inside character constant.");
			return;
		}
		else if (c == '\n')
		{
			ErrorAt(ERROR_PARSER, file->filename, &pos,
					"End of line inside character constant.");
			return;
		}
		else if (FileMatch(file, "\\\'"))
		{
			FileConsume(file, 2);
		}
		else if (c == '\'')
		{
			FileConsume(file, 1);
			return;
		}
		else
		{
			FileConsume(file, 1);
		}
	}
}


static bool MatchAndConsumeList(ParserFile *file, TokenInfo *list, const TokenType **tt_out)
{
	for (int i=0; list[i].text != NULL; i++)
	{
		if (FileMatchAndConsume(file, list[i].text))
		{
			*tt_out = list[i].token_type;
			return true;
		}
	}
	*tt_out = NULL;
	return false;
}

static const TokenType *MatchTokenList(TokenInfo *list, const char *text, int length)
{
	for (int i=0; list[i].text != NULL; i++)
	{
		if ((strncmp(list[i].text, text, length) == 0)
				&& (list[i].text[length] == 0))
		{
			return list[i].token_type;
		}
	}
	return NULL;
}

bool IsValidNamespace(const char *txt)
{
	// FIXME don't allow "_" as a namespace element

	bool word_start = true;
	while (true)
	{
		int c = *txt;

		if (c == 0)
			return !word_start;

		if (word_start)
		{
			if (!IsLetter(c))
				return false;
			word_start = false;
		}
		else
		{
			if (c == '.')
				word_start = true;
			else if (!IsLetter(c) && !IsDigit(c))
				return false;
		}

		txt ++;
	}
}

bool IsValidNamespaceName(String *str)
{
	if (str->length <= 0)
		return false;

	if ((str->length == 1) && (str->data[0] == '_'))
		return false;

	if (!IsLetter(str->data[0]))
		return false;

	for (int i=1; i<str->length; i++)
	{
		if (!IsLetter(str->data[i]) && !IsDigit(str->data[i]))
			return false;
	}

	return true;
}

