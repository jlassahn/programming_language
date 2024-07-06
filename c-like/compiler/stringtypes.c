
#include "compiler/stringtypes.h"

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

