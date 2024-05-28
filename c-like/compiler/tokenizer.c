
#include "errors.h"
#include "tokenizer.h"
#include "parser_file.h"
#include <string.h>

const TokenType TOKEN_EOF = { "EOF", 0x0000 };
const TokenType TOKEN_PPIDENTIFIER = { "IDENTIFIER", 0x0000 };
const TokenType TOKEN_PPNUMBER = { "PPNUMBER", 0x0000 };
const TokenType TOKEN_CHARCONST = { "CHARCONST", 0x0000 };
const TokenType TOKEN_STRINGCONST = { "STRINGCONST", 0x0000 };
const TokenType TOKEN_OPERATOR = { "OPERATOR", 0x0000 };
const TokenType TOKEN_OTHER = { "OTHER", 0x0000 };

// shorter tokens must follow longer tokens which have the same prefix
const char *operator_list[] =
{
	"...",
	"<<=",
	">>=",
	"->",
	"++",
	"--",
	"<<",
	">>",
	"<=",
	">=",
	"==",
	"!=",
	"&&",
	"||",
	"*=",
	"/=",
	"%=",
	"+=",
	"-=",
	"&=",
	"^=",
	"|=",
	NULL
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
static bool MatchAndConsumeList(ParserFile *file, const char **list);

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

	const char *cur = FileGet(file);
	if (cur[0] == 0)
	{
		tokenizer->current_token.token_type = &TOKEN_EOF;
	}
	else if (IsLetter(cur[0]))
	{
		tokenizer->current_token.token_type = &TOKEN_PPIDENTIFIER;
		ConsumePPIdentifier(tokenizer);
	}
	else if (IsDigit(cur[0]) || (cur[0] == '.' && IsDigit(cur[1])))
	{
		tokenizer->current_token.token_type = &TOKEN_PPNUMBER;
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
	else if (MatchAndConsumeList(file, operator_list))
	{
		tokenizer->current_token.token_type = &TOKEN_OPERATOR;
	}
	else
	{
		tokenizer->current_token.token_type = &TOKEN_OTHER;
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
	fprintf(fp, "[%s:%ld:%ld]%s[%ld:%ld] = [%.*s]\n",
			token->position.file->filename,
			token->position.start.line+1,
			token->position.start.byte_in_line+1,
			token->token_type->name,
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
				ErrorAt(file->filename, &pos, "No end of line after #");
				return;
			}
		}
		else if (FileMatchAndConsume(file, "//"))
		{
			if (!SeekTo(file, "\n"))
			{
				ErrorAt(file->filename, &pos, "No end of line after //");
				return;
			}
		}
		else if (FileMatchAndConsume(file, "/*"))
		{
			if (!SeekTo(file, "*/"))
			{
				ErrorAt(file->filename, &pos, "End of file inside /*");
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
			ErrorAt(file->filename, &pos, "End of file inside string.");
			return;
		}
		else if (c == '\n')
		{
			ErrorAt(file->filename, &pos, "End of line inside string.");
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
			ErrorAt(file->filename, &pos,
					"End of file inside character constant.");
			return;
		}
		else if (c == '\n')
		{
			ErrorAt(file->filename, &pos,
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


static bool MatchAndConsumeList(ParserFile *file, const char **list)
{
	for (int i=0; list[i] != NULL; i++)
	{
		if (FileMatchAndConsume(file, list[i]))
			return true;
	}
	return false;
}

