
#include "tests/unit/utils.h"
#include <string.h>

String *TempString(const char *cstr)
{
	static String s;
	s. length = strlen(cstr);
	s.data = cstr;
	return &s;
}

