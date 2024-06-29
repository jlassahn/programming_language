
#ifndef INCLUDED_EXIT_CODES_H
#define INCLUDED_EXIT_CODES_H

// These kind of match BSD exit codes from sysexits.h
typedef enum
{
	EXIT_OK = 0,
	EXIT_USAGE = 64,    // bad command line args
	EXIT_DATAERR = 65,  // syntax errors, general compile errors
	EXIT_NOINPUT = 66,  // can't open input file
	EXIT_SOFTWARE = 70, // out of memory, assert failed, compiler bugs
	EXIT_IOERR = 74     // can't open output, file I/O failure, etc
}
ExitCode;

#endif

