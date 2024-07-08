
#ifndef INCLUDED_FAKE_PARSER_H
#define INCLUDED_FAKE_PARSER_H

#include "compiler/parser_node.h"

void FakeParserSet(const char *path, ParserNode *root, int error_code);
void FakeParserFree(void);

#endif

