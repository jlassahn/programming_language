
# The Janus Programming Language

Janus is a language with a split personality.  It can be used both as an
applications programming language which provides moderately high performance
along with the expressiveness and safety needed for large programming projects.
It can also be used for embedded and systems programming tasks which require
low overhead, high performance, and do not necessarily have access to the
operating system services which applications programmers are accustomed to.

The language supports these two programming domains by providing two
slightly different programming paradigms with different data types, along
with a well-defined set of rules for interactions between them.  A given
program may use either or both of these depending on its requirements.

## Machine Janus
Machine Janus is a subset of the language which requires little system
support at runtime.  Data structures written in Machine Janus have little
or no hidden overhead and a predictable layout in memory.  Access to system
services is through explicit function calls.  Dynamic memory allocation
is done with explicit allocate and free calls, with the programmer responsible
for managing the lifetime of allocated memory.

## Object Janus
Object Janus is a subset of the language which provides convenient support
for large programs using object oriented programming practices.  Objects
are typically dynamically allocated, with automatic garbage collection to
manage object lifetimes.

## A Simple Program

'''
janus 1.0;

import io;

def main() -> Void
{
	io.stdout.print("Hello, World`lf`");
	return 0;
}
'''


