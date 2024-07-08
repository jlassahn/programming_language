
#ifndef INCLUDED_PARSER_H
#define INCLUDED_PARSER_H

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include "compiler/parser_file.h"
#include "compiler/parser_node.h"

#include "compiler/parser_symbols.h" //FIXME maybe don't do this automatically

typedef struct ParserContext ParserContext; // FIXME maybe not needed?

void ParseSetDebug(bool on);

ParserNode *ParseFile(ParserFile *file, ParserContext *context);

#endif

