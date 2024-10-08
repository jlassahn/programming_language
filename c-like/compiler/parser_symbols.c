
#include "compiler/parser_symbols.h"

ParserSymbol SYM_UNDEF = { "UNDEFINED", 0 };
ParserSymbol SYM_IDENTIFIER = { "IDENTIFIER", PRINT_CONTENT };
ParserSymbol SYM_NUMBER = { "NUMBER", PRINT_CONTENT };
ParserSymbol SYM_CHARCONST = { "CHARCONST", PRINT_CONTENT };
ParserSymbol SYM_STRINGCONST = { "STRINGCONST", PRINT_CONTENT };
ParserSymbol SYM_PUNCTUATION = { "PUNCTUATION", PRINT_CONTENT | SYM_DISCARD };
ParserSymbol SYM_KEYWORD = { "KEYWORD", PRINT_CONTENT | SYM_DISCARD };
ParserSymbol SYM_OPERATOR = { "OPERATOR", PRINT_CONTENT | SYM_DISCARD };

ParserSymbol SYM_EMPTY = { "EMPTY", 0 };
ParserSymbol SYM_LIST = { "LIST", 0 };
ParserSymbol SYM_DOT_OP = { "DOT_OP", 0 };
ParserSymbol SYM_IMPORT = { "IMPORT", 0 };
ParserSymbol SYM_IMPORT_PRIVATE = { "IMPORT_PRIVATE", 0 };
ParserSymbol SYM_PROTOTYPE = { "PROTOTYPE", 0 };
ParserSymbol SYM_FUNC = { "FUNC", 0 };
ParserSymbol SYM_DECLARATION = { "DECLARATION", 0 };
ParserSymbol SYM_USING = { "USING", 0 };
ParserSymbol SYM_USING_AS = { "USING_AS", 0 };
ParserSymbol SYM_TRAILING_COMMA = { "TRAILING_COMMA", 0 };
ParserSymbol SYM_ELLIPSIS = { "ELLIPSIS", 0 };
ParserSymbol SYM_PARAMETER = { "PARAMETER", 0 };
ParserSymbol SYM_PARAM_TYPE = { "PARAM_TYPE", 0 };
ParserSymbol SYM_DECL_TYPE = { "DECL_TYPE", 0 };
ParserSymbol SYM_INITIALIZE = { "INITIALIZE", 0 };
ParserSymbol SYM_INIT_STRUCT = { "INIT_STRUCT", 0 };
ParserSymbol SYM_INIT_ARRAY = { "INIT_ARRAY", 0 };
ParserSymbol SYM_STATEMENT_LIST = { "STATEMENT_LIST", 0 };
ParserSymbol SYM_STRUCT_LIST = { "STRUCT_LIST", 0 };
ParserSymbol SYM_ENUM_LIST = { "ENUM_LIST", 0 };
ParserSymbol SYM_ENUM_ELEMENT = { "ENUM_ELEMENT", 0 };
ParserSymbol SYM_STRUCT_DEC = { "STRUCT_DEC", 0 };
ParserSymbol SYM_STRUCT_DEF = { "STRUCT_DEF", 0 };
ParserSymbol SYM_UNION_DEC = { "UNION_DEC", 0 };
ParserSymbol SYM_UNION_DEF = { "UNION_DEF", 0 };
ParserSymbol SYM_ENUM_DEC = { "ENUM_DEC", 0 };
ParserSymbol SYM_ENUM_DEF = { "ENUM_DEF", 0 };
ParserSymbol SYM_EXPRESSION_STATEMENT = { "EXPRESSION_STATEMENT", 0 };
ParserSymbol SYM_EMPTY_STATEMENT = { "EMPTY_STATEMENT", 0 };
ParserSymbol SYM_LABEL_STATEMENT = { "LABEL_STATEMENT", 0 };
ParserSymbol SYM_FOR_STATEMENT = { "FOR_STATEMENT", 0 };
ParserSymbol SYM_WHILE_STATEMENT = { "WHILE_STATEMENT", 0 };
ParserSymbol SYM_DO_STATEMENT = { "DO_STATEMENT", 0 };
ParserSymbol SYM_IF_STATEMENT = { "IF_STATEMENT", 0 };
ParserSymbol SYM_IF_ELSE = { "IF_ELSE", 0 };
ParserSymbol SYM_SWITCH_STATEMENT = { "SWITCH_STATEMENT", 0 };
ParserSymbol SYM_BREAK_STATEMENT = { "BREAK_STATEMENT", 0 };
ParserSymbol SYM_CONTINUE_STATEMENT = { "CONTINUE_STATEMENT", 0 };
ParserSymbol SYM_GOTO_STATEMENT = { "GOTO_STATEMENT", 0 };
ParserSymbol SYM_RETURN_STATEMENT = { "RETURN_STATEMENT", 0 };
ParserSymbol SYM_RETURN_VOID = { "RETURN_VOID", 0 };
ParserSymbol SYM_CASE_ELEMENT = { "CASE_ELEMENT", 0 };
ParserSymbol SYM_CASE_END_ELEMENT = { "CASE_END_ELEMENT", 0 };
ParserSymbol SYM_DEFAULT_LABEL = { "DEFAULT_LABEL", 0 };
ParserSymbol SYM_CASE_LABEL = { "CASE_LABEL", 0 };
ParserSymbol SYM_CONSTANT = { "CONSTANT", 0 };
ParserSymbol SYM_ASSIGN_OP = { "ASSIGN_OP", 0 };
ParserSymbol SYM_ASSIGN_MULT_OP = { "ASSIGN_MULT_OP", 0 };
ParserSymbol SYM_ASSIGN_DIV_OP = { "ASSIGN_DIV_OP", 0 };
ParserSymbol SYM_ASSIGN_MOD_OP = { "ASSIGN_MOD_OP", 0 };
ParserSymbol SYM_ASSIGN_ADD_OP = { "ASSIGN_ADD_OP", 0 };
ParserSymbol SYM_ASSIGN_SUB_OP = { "ASSIGN_SUB_OP", 0 };
ParserSymbol SYM_ASSIGN_SHR_OP = { "ASSIGN_SHR_OP", 0 };
ParserSymbol SYM_ASSIGN_SHL_OP = { "ASSIGN_SHL_OP", 0 };
ParserSymbol SYM_ASSIGN_AND_OP = { "ASSIGN_AND_OP", 0 };
ParserSymbol SYM_ASSIGN_OR_OP = { "ASSIGN_OR_OP", 0 };
ParserSymbol SYM_ASSIGN_XOR_OP = { "ASSIGN_XOR_OP", 0 };
ParserSymbol SYM_CONDITIONAL = { "CONDITIONAL", 0 };
ParserSymbol SYM_LOG_OR_OP = { "LOG_OR_OP", 0 };
ParserSymbol SYM_LOG_AND_OP = { "LOG_AND_OP", 0 };
ParserSymbol SYM_OR_OP = { "OR_OP", 0 };
ParserSymbol SYM_AND_OP = { "AND_OP", 0 };
ParserSymbol SYM_ADDR_OP = { "ADDR_OP", 0 };
ParserSymbol SYM_XOR_OP = { "XOR_OP", 0 };
ParserSymbol SYM_EQUAL_OP = { "EQUAL_OP", 0 };
ParserSymbol SYM_NEQUAL_OP = { "NEQUAL_OP", 0 };
ParserSymbol SYM_LESS_OP = { "LESS_OP", 0 };
ParserSymbol SYM_GREATER_OP = { "GREATER_OP", 0 };
ParserSymbol SYM_LESSEQ_OP = { "LESSEQ_OP", 0 };
ParserSymbol SYM_GREATEREQ_OP = { "GREATER_EQ_OP", 0 };
ParserSymbol SYM_SHL_OP = { "SHL_OP", 0 };
ParserSymbol SYM_SHR_OP = { "SHR_OP", 0 };
ParserSymbol SYM_ADD_OP = { "ADD_OP", 0 };
ParserSymbol SYM_SUB_OP = { "SUB_OP", 0 };
ParserSymbol SYM_DIV_OP = { "DIV_OP", 0 };
ParserSymbol SYM_MOD_OP = { "MOD_OP", 0 };
ParserSymbol SYM_MULT_OP = { "MULT_OP", 0 };
ParserSymbol SYM_PTR_OP = { "PTR_OP", 0 };
ParserSymbol SYM_NOT_OP = { "NOT_OP", 0 };
ParserSymbol SYM_BITNOT_OP = { "BITNOT_OP", 0 };
ParserSymbol SYM_PREINC_OP = { "PREINC_OP", 0 };
ParserSymbol SYM_PREDEC_OP = { "PREDEC_OP", 0 };
ParserSymbol SYM_POSTINC_OP = { "POSTINC_OP", 0 };
ParserSymbol SYM_POSTDEC_OP = { "POSTDEC_OP", 0 };
ParserSymbol SYM_NEG_OP = { "NEG_OP", 0 };
ParserSymbol SYM_POS_OP = { "POS_OP", 0 };
ParserSymbol SYM_SIZEOF_OP = { "SIZEOF_OP", 0 };
ParserSymbol SYM_ARRAY_OP = { "ARRAY_OP", 0 };
ParserSymbol SYM_CALL_OP = { "CALL_OP", 0 };
ParserSymbol SYM_INIT_OP = { "INIT_OP", 0 };
ParserSymbol SYM_PAREN_EXPRESSION = { "PAREN_EXPRESSION", 0 };
ParserSymbol SYM_STRING = { "STRING", 0 };
ParserSymbol SYM_TYPE_EXPRESSION = { "TYPE_EXPRESSION", 0 };
ParserSymbol SYM_TYPE_ARRAY = { "TYPE_ARRAY", 0 };
ParserSymbol SYM_TYPE_ARRAY_MATCH = { "TYPE_ARRAY_MATCH", 0 };
ParserSymbol SYM_TYPE_BITFIELD = { "TYPE_BITFIELD", 0 };
ParserSymbol SYM_TYPE_FUNCTION = { "TYPE_FUNCTION", 0 };
ParserSymbol SYM_TYPE_LINKAGE = { "TYPE_LINKAGE", 0 };
ParserSymbol SYM_TYPE_LINKNAME = { "TYPE_LINKNAME", 0 };

