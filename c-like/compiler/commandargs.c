
#include "compiler/commandargs.h"
#include "compiler/memory.h"
#include "compiler/types.h"
#include "compiler/errors.h"
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdarg.h>

static CompilerArgs args;

typedef struct OptionList OptionList;
struct OptionList
{
	const char *name;
	const char *help;
};

typedef struct CommandArg CommandArg;
struct CommandArg
{
	const char *flag;
	ArgStringList **dest;
	OptionList *opts;
	const char *help;
};

static OptionList warn_opts[] =
{
	{"all", "Enable all generally useful warnings."},
	{NULL, NULL}
};

static OptionList opt_opts[] =
{
	{"speed", "Optimize for speed."},
	{"size", "Optimize for code size."},
	{NULL, NULL}
};

static OptionList gen_opts[] =
{
	{"debug", "Generate debyg symbols."},
	{"mosscall", "Use normal Moss calling conventions."},
	{"isource", "Prefer to compile imports from source."},
	{"istatic", "Prefer to link imports from static libraries."},
	{"ishared", "Prefer to link imports from shared libraries."},
	{"exe", "Generate an executable (the default)."},
	{"obj", "Generate object files."},
	{"static", "Generate a static library."},
	{"shared", "Generate a shared library."},
	{NULL, NULL}
};

static CommandArg arg_list[] =
{
	{"-W", &args.warnings,      warn_opts, " [suboption]  Control warnings."},
	{"-O", &args.optimizations, opt_opts,  " [suboption]  Control optimization."},
	{"-Z", &args.generation,    gen_opts,  " [suboption]  Code generation options."},
	{"-I", &args.basedirs,      NULL,      " [path]       Add a search path for input files."},
	{"-D", &args.defines,       NULL,      " [name=value] Define an option variable value."},
	{"-V", &args.versions,      NULL,      " [version]    Select language version."},
	{"-FO", &args.outfile,      NULL,       "[filename]   Change the main output filename."},
	{"-FI", &args.outdir,       NULL,       "[path]       Path to the output file tree.."},
	{"-FT", &args.treefile,     NULL,       "[filename]   Filename for import tree output."},
	{NULL, NULL, NULL, NULL} // end of list
};

static void PrintHelp(void)
{
	printf(
		"moss-cc [options] [inputs]\n"
		"Compiler for Moss code.\n"
		"Inputs can be:\n"
		"    A module name e.g. user.test.hello_world\n"
		"    A Moss source file e.g. ./hello.moss\n"
		"    A library file e.g. ./libhello.a, ./libhello.so\n"
		"\n"
		"Options can be:\n"
	);

	for (int i=0; arg_list[i].dest != NULL; i++)
	{
		printf("    %s %s\n", arg_list[i].flag, arg_list[i].help);
		if (arg_list[i].opts)
		{
			OptionList *opts = arg_list[i].opts;
			for (int j=0; opts[j].name != NULL; j++)
				printf("            %s -- %s\n", opts[j].name, opts[j].help);
		}
	}
}

static void ArgError(const char *txt, ...)
{
	va_list ap;
	va_start(ap, txt);
	fprintf(stderr, "ERROR: ");
	vfprintf(stderr, txt, ap);
	va_end(ap);
	fprintf(stderr, "Use -h for command parameter help.\n");
}

static void AddArg(ArgStringList **list, const char *value)
{
	ArgStringList *item = Alloc(sizeof(ArgStringList));
	memset(item, 0, sizeof(ArgStringList));

	item->next = NULL;
	item->arg = value;
	if (*list == NULL)
	{
		*list = item;
		return;
	}

	ArgStringList *tail = *list;
	while (tail->next)
		tail = tail->next;
	tail->next = item;
}

const CompilerArgs *ParseArgs(int argc, const char *argv[])
{
	ArgStringList **dest;
	const char *param;
	const char *flag;

	for (int i=1; i<argc; i++)
	{
		const char *arg = argv[i];
		if (arg[0] == '@')
		{
			if (arg[1] == 0)
			{
				if (i == argc-1)
				{
					ArgError("@ must be followed by a filename\n");
					FreeArgs(&args);
					return NULL;
				}
				i ++;
				param = argv[i];
			}
			else
			{
				param = arg + 1;
			}
			printf("@file = %s\n", param);
		}
		else if (arg[0] == '-')
		{
			if ((strcmp(arg, "-h") == 0) ||
				(strcmp(arg, "-?") == 0) ||
				(strcmp(arg, "--help") == 0))
			{
				PrintHelp();
				FreeArgs(&args);
				return NULL;
			}

			dest = NULL;
			flag = NULL;
			param = NULL;
			for (int j=0; arg_list[j].dest != NULL; j++)
			{
				flag = arg_list[j].flag;
				dest = arg_list[j].dest;
				int len = strlen(flag);
				if (strcmp(arg, flag) == 0)
				{
					if (i == argc-1)
					{
						ArgError("option %s must be followed by a filename\n", flag);
						FreeArgs(&args);
						return NULL;
					}
					i ++;
					param = argv[i];
					break;
				}
				else if (strncmp(arg, flag, len) == 0)
				{
					param = arg + len;
					break;
				}
			}

			if (param)
			{
				AddArg(dest, param);
			}
			else
			{
				ArgError("Unrecognized option: %s\n", arg);
				FreeArgs(&args);
				return NULL;
			}
		}
		else
		{
			AddArg(&args.inputs, arg);
		}
	}

	return &args;
}

static void FreeArgList(ArgStringList *list)
{
	ArgStringList *next;
	while (list)
	{
		next = list->next;
		Free(list);
		list = next;
	}
}

void FreeArgs(const CompilerArgs *args_in)
{
	CompilerArgs *args = (CompilerArgs *)args_in;
	FreeArgList(args->inputs);
	FreeArgList(args->warnings);
	FreeArgList(args->optimizations);
	FreeArgList(args->generation);
	FreeArgList(args->defines);
	FreeArgList(args->basedirs);
	FreeArgList(args->versions);
	FreeArgList(args->outfile);
	FreeArgList(args->outdir);
	FreeArgList(args->treefile);
	memset(args, 0, sizeof(CompilerArgs));
}

void PrintArgList(ArgStringList *list)
{
	ArgStringList *item;
	for (item = list; item != NULL; item = item->next)
		printf("    %s\n", item->arg);
}

