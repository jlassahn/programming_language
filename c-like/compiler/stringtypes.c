
#include "compiler/stringtypes.h"
#include "compiler/types.h"
#include "compiler/errors.h"

bool IsSpace(char x)
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

bool IsLetter(char x)
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

bool IsDigit(char x)
{
	if ((x >= '0') && (x <= '9'))
		return true;
	return false;
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

bool IsValidPath(const char *txt)
{
	int prevchar = -1;

	if (*txt == 0)
		return false;

	// checks for first character in first filename
	if (txt[0] == ' ')
		return false;

	while (*txt != 0)
	{
		int c = *txt;

		if (c < 0x20)
			return false;
		if (c == 0x7F)
			return false;
		if (c == '\"')
			return false;
		if (c == '\'')
			return false;

		// checks for first character in filename
		if ((prevchar == '/') || (prevchar == '\\'))
		{
			if (c == '/')
				return false;
			if (c == '\\')
				return false;
			if (c == ' ')
				return false;
		}

		// checks for last character in internal filenames
		if ((c == '/') || (c == '\\'))
		{
			if (prevchar == ' ')
				return false;
		}

		prevchar = c;
		txt ++;
	}

	// checks for last character in final filename
	if (prevchar == ' ')
		return false;

	return true;
}

USE_RESULT
StringBuffer *NormalizePath(StringBuffer *path)
{
	for (int i=0; i<path->string.length; i++)
	{
		if (path->buffer[i] == '\\')
			path->buffer[i] = '/';
	}
	if (path->buffer[path->string.length-1] != '/')
		path = StringBufferAppendChars(path, "/");

	return path;
}

static bool AppendPathItem(List *list, String *item)
{
	StringBuffer *sb = StringBufferFromString(item);
	if (IsValidPath(sb->buffer))
	{
		sb = NormalizePath(sb);
		StringBufferLock(sb);
		ListInsertLast(list, sb);
		return true;
	}
	else
	{
		Error(ERROR_FILE, "Invalid name in path list '%s'", sb->buffer);
		StringBufferFree(sb);
		return false;
	}
}

// FIXME maybe only needed in pass_configure?
bool AppendPathList(List *list, const char *env)
{
	bool ret = true;
	int start = 0;
	int end = 0;
	for (end=0; env[end]!=0; end++)
	{

		if (env[end] == ';')
		{
			String s;
			s.data = &env[start];
			s.length = end-start;

			if (!AppendPathItem(list, &s))
				ret = false;
			start = end+1;
		}
	}
	if (end > start)
	{
		String s;
		s.data = &env[start];
		s.length = end-start;

		if (!AppendPathItem(list, &s))
			ret = false;
	}

	return ret;
}

