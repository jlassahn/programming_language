
#ifndef INCLUDED_FAKE_NODES_H
#define INCLUDED_FAKE_NODES_H

#include "compiler/parser_node.h"

ParserNode *MakeNodeFakeValue(ParserSymbol *symbol, const char *value);
int PushNodeStack(ParserNode *node);
int MakeNodeOnStack(ParserSymbol *symbol, int count);
ParserNode *GetNodeStackTop(void);
void FreeFakeNodeValues(void);

#endif

