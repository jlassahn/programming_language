
#include "compiler/pass_configure.h"
#include "compiler/errors.h"
#include "compiler/types.h"
#include "compiler/stringtypes.h"
#include "compiler/fileio.h"
#include "compiler/compiler_file.h"
#include <string.h>

bool AddBaseDir(CompileState *cs, const char *path)
{
	if (!IsValidPath(path))
	{
		Error(ERROR_FILE,
				"parameter '%s' is not a valid path.", path);
		return false;
	}

	StringBuffer *sb = StringBufferFromChars(path);
	sb = NormalizePath(sb);

	StringBufferLock(sb);
	ListInsertLast(&cs->basedirs, sb);
	return true;
}

bool AddInputModule(CompileState *state, const char *path)
{
	if (!IsValidNamespace(path))
		return false;

	Namespace *ns = &state->root_namespace;

	String name;
	int length = strlen(path);
	int start = 0;
	int end = 0;
	for (int i=0; i<length; i++)
	{
		if (path[i] == '.')
		{
			end = i;
			name.data = &path[start];
			name.length = end-start;
			ns = NamespaceGetChild(ns, &name);
			if (ns == NULL)
				return false;
			start = i+1;
		}
	}
	if (start < length)
	{
		end = length;
		name.data = &path[start];
		name.length = end-start;
		ns = NamespaceGetChild(ns, &name);
		if (ns == NULL)
			return false;
	}

	ListInsertLast(&state->input_modules, ns);

	return true;
}

bool AddInputFile(CompileState *state, const char *name)
{
	StringBuffer *sb = StringBufferFromChars(name);
	StringBufferLock(sb);

	// FIXME if name ends with ".moss" assume it's a source file
	//       otherwise it's a library or something.

	CompilerFile *cf = CompilerFileCreate(sb);
	cf->flags |= FILE_FROM_INPUT;
	ListInsertLast(&state->input_files, cf);

	return true;
}

bool AddInput(CompileState *cs, const char *name)
{
	if (!IsValidPath(name))
	{
		Error(ERROR_FILE,
				"parameter '%s' is not a valid input name.", name);
		return false;
	}

	if (AddInputModule(cs, name))
		return true;
	if (AddInputFile(cs, name))
		return true;

	Error(ERROR_FILE,
			"Input name '%s' must be either a file or a module.", name);
	return false;
}

bool PassConfigure(CompileState *state, const CompilerArgs *args,
		const char *env)
{
	AddBaseDir(state, ".");

	if (!AppendPathList(&state->basedirs, env))
	{
		Error(ERROR_FILE,
				"Moss input path environment variable has bad paths.");
	}

	for (ArgStringList *entry=args->basedirs; entry!=NULL; entry=entry->next)
	{
		AddBaseDir(state, entry->arg);
	}

	for (ArgStringList *entry=args->inputs; entry!=NULL; entry=entry->next)
	{
		AddInput(state, entry->arg);
	}

	// FIXME handle other args
	// check args for validity
	//    warnings;
	//    optimizations;
	//    generation;
	//    defines;
	//    versions;
	//    outfile;
	//    outdir;
	//    treefile;

	// true to continue compilation.  So far no errors are bad enough to abort.
	return true;
}

