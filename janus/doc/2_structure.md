
# Program Structure

Very short Janus programs may be written as a single source file.  Larger
programs, or code that is part of a reusable library, should be organized
into multiple files which follow specific conventions about names and
locations.

FIXME what about libraries???  Is there a .jlib file to describe linkable
objects?

## Modules, Names, and Paths

Symbols from one source file can be made visible in another using the
import statement.  The code made visible by a single import statement is
called a "module".

Module names are tokens separated by dots.  The compiler searches for
source files for a module by converting these to directory and file names.
For example,

```janus
import experimental.web_handlers;
```

causes the compiler to look for these files:
* interfaces/experimental/web_handlers.janus
* source/experimental/web_handlers.janus
* source/experimental/web_handlers.jsrc

It looks for these under a set of base directories which depend on how
the compiler is configured, but typically includes the current directory.

Symbols declared in files ending in .janus are externally visible to
programs which import the .janus file.  Symbols from files ending in .jsrc
are local to the source file.  For reusable code, it is good style
to declare visible variables and functions in a .janus file without
defintions, and provide definitions for them in a .jsrc file of the same
name.

```janus

# example.janus
janus 1.0;

def global_variable Float64;

def GetGlobal() -> Float64;
def SetGlobal(x Float64);
```

```janus

# example.jsrc
janus 1.0;

def global_variable Float64 = 0.0;

def GetGlobal() -> Float64
{
	InternalCall();
	return global_variable;
}

def SetGlobal(x Float64)
{
	InternalCall();
	global_variable = x;
}

def InternalCall()
{
	io.stdout.print("a global access has happened`lf`");
}

```

## Overriding Module Name Defaults

Almost always, it is best to use the compiler's default mappings between
file names and module names.  They can be overridden using janus header
options, however.

```janus

janus 1.0
{
	module_name = example.module;
	export_symbols = True;
}

```

In general, the compiler will not be able to automatically find imported
files which have set module_name, and they must be explicitly specified.


## Detailed Rules for Symbol Resolution

Every file has a module name and an export state.

The module name is the  module_name header option, if one is specified.

Otherwise, the module name is derived from the path to the file.  It is an
error if the path name is not a valid module name (e.g. if the file name
contains spaces).

The export state is the export_symbols header option, if one is specified.

Otherwise if the filename ends in .janus the export state is True, if not
it is False.

Any symbol declared at file scope in a file with export state True is
externally visible.

Any symbol declared at file scope only in files with export state False is
not externally visible.

All symbols declared at file scope with the same name in all files with the
same module name refer to the same object.  All these declarations must
have the same type.

At most one of the declarations of a given object may initialize it.  Multiple
intializations are an error even if they specify the same value.  If none
of an object's declarations initialize it, it is intialized with the
default value for that data type.


# The Main Function

One file in any program should declare a function at file scope named Main.
The Main function is never exported, and can not be accessed through an
import.  Only one file in a program can contain a Main function.

The Main function must either have no return value, or return any integer
type. It must take either no parameters or a single parameter of one
of these types:
* array(String)
* m_slice(MString)

FIXME need to work out string types in more detail

It is possible to export the Main function indirectly by aliasing it:
```janus

# MainAlias is exported using the same rules as other functions.
def MainAlias() = Main;

def Main()
{
}
```

