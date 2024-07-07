
#include "compiler/exit_codes.h"
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/fileio.h"
#include "compiler/errors.h"
#include "compiler/commandargs.h"
#include "compiler/parser_file.h"
#include "compiler/compiler_file.h"
#include "compiler/tokenizer.h"
#include "compiler/parser.h"
#include "compiler/compile_state.h"
#include "compiler/namespace.h"
#include "compiler/search.h"
#include "compiler/stringtypes.h"
#include "compiler/pass_configure.h"
#include "compiler/pass_search_and_parse.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

static CompileState compile_state;

int main(int argc, const char *argv[])
{
	CompileStateInit(&compile_state);

	const CompilerArgs *args = ParseArgs(argc, argv);
	if (args == NULL)
		return EXIT_USAGE;

	const char *env = getenv("MOSS_IMPORT_PATH");
	if (env == NULL)
		env = "";

	if (!PassConfigure(&compile_state, args, env))
		return EXIT_USAGE;

	CompileStatePrint(&compile_state);

	ParseSetDebug(false);
	bool inputs_good = PassSearchAndParse(&compile_state);

	if ((compile_state.input_files.first == NULL)
			&& (compile_state.input_modules.first == NULL))
	{
		Error(ERROR_FILE, "No inputs specified.");
		inputs_good = false;
	}

	if (!inputs_good)
	{
		printf("BAD INPUTS\n");
		// FIXME skip compile steps and exit
	}

	// FIXME PassTranslate(&compile_state);
	// FIXME PassGenerate(&compile_state);
	// FIXME PassLink(&compile_state);

	int depth = 1;
	printf("\nNAMESPACE:\n");
	MapIterate(&compile_state.root_namespace.children,
			NamespacePrinter, &depth);

	// FIXME cleanup starts here

	CompileStateFree(&compile_state);
	FreeArgs(args);

	printf("allocation count = %d\n", AllocCount());

	printf("errors: %d, warnings: %d\n", ErrorCount(), WarningCount());

	if (ErrorCount() > 0)
		return EXIT_DATAERR;
	return EXIT_OK;
}

