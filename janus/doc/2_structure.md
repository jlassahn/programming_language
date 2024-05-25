
# Program Structure

Very short Janus programs may be written as a single source file.  Larger
programs, or code that is part of a reusable library, should be organized
into multiple files which follow specific conventions about names and
locations.

FIXME what about libraries???  Is there a .jlib file to describe linkable
objects?

FIXME what about accessing symbols that aren't exported, e.g. for tests?
maybe allow
```
def some.other.module.FunctionToTest();
```
but only in a non-exported context.  Can we also use that to inject methods
into other modules?

## The Namespace

Named symbols exist in a hierarchical namespace.
FIXME the root of the namespace should be called something like "global"
	we're using . or other special notations for the global namespace but
	we should probably change everything to use global instead.

When each source file is processed, it builds up its own global namespace.
Things usually tend to be in the same place in different files' namespaces,
but things like import module name overrides can make different files
have things in different places.

When the main program is compiled it puts its top-level symbols into global.
FIXME how to tell which files are the "main program"?
FIXME would it be better to put all symbols into nested namespaces?
	maybe any file listed on the command line of the compiler that
	doesn't explicitly call out a module name is in global?

When a file is imported, it can only put things into the namespace it is
importing into. so
```
import containers.vector;
```
will only add things to global.containers.vector, while
```
import global = containers.vector;
```
will add the same things to global instead.

FIXME what happens to imports inside of imported files?  Do they get
added to the top-level file's namespace?  If so, where?
FIXME are two copies of a type which get added to the namespace in different
places still the same type?
```
#application
import module1;

# Can I access module2.Thing here?
# or is it module1.module2.Thing?
# what if I did import alias = module1; ?
```

```
#module1
import module2;
def doStuff(module2.Thing a, Int32 b) -> Module2.Thing;
```

Probably there's a "canonical" namespace, where each symbol has it's full
unaliased name.  Things that have the same canonical name refer to the same
object.  It's an error for the canonical namespace to be inconsistent.
e.g. having multiple files that claim to describe the same full module path
but don't have consistent definitions for everything.

Probably imports inside of imports don't populate the top-level caller's
namespace, which means the compile has to deal with symbols that are used
in the code but aren't active in the namespace.  E.g. a function in module1
takes a parameter type declared in module2.  The module1 code must import
module2, but an application that imports module1 has access to the function
even though it hasn't imported module2.  So the compiler has to reason about
types in module2 which aren't in the application namespace.

FIXME consider making imports always put things in the canonical spot, then
having a separate alias feature that makes links to things.

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

Normally these symbols are named using the full module name.  If that is not
convenient, the module name can be overridden.  For example doing
`import handlers = experimental.web_handlers;` makes a symbol normally named
`experimental.web_handlers.BaseHandler` instead be `handlers.BaseHandler`.
Doing `import . = experimental.web_handlers` causes all the symbols in that
module to be visible with _no_ module name prefix at all.

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

import . = example;

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

Modules typically have either a single .janus file, or one .janus file and
one .jsrc file.  It is sometimes possible to have multiple files in the
same module.  This is permitted, and follows the symbol resolution rules
below, but is bad style.  Compilers may generate a warning about this.

Frequently, the .jsrc file for a module will import it's own module name.
This gives the .jsrc file access to everything declared in the .janus file.
If a file does not import it's own module, it will only have access to
symbols it explicitly declares.

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

Expressions in a file may only reference symbols explicitly declared in or
imported by that file.  In particular if a module has a source file which
declares an unexported symbol, other source files in that same module can
only access the symbol if they also declare it.


## Modules With Multiple Source Files

Sometimes a module is large enough that splitting it into multiple files is
useful.  The usual style for doing this is to have a single .janus file
which declares all the external symbols, a single .jsrc file which has
definitiions for those symbols, and submodules with only .jsrc files which
contain additional code.  For example:

```janus
# example.janus

janus 1.0;

def SomeComplicatedFunction();
def AnotherComplicatedFuncion();
```

```janus
# example.jsrc

janus 1.0;

# import some subpackages from
# example/utilities.jsrc
# and
# example/another_function.jsrc
import example.utilities;
import example.another_function;

def SomeComplicatedFunction()
{
	# call some stuff in the subpackages
	def x = example.utilities.DoSomething();
	example.utilities.DoMoreStuff(x);
}

# alias a complete function definition from a subpackage
def AnotherComplicatedFunction() =
	example.another_function.AnotherComplicatedFunction;
```


# The Main Function

One file in any program should declare a function at file scope named Main.

The Main function must either have no parameters and Void return type:

```janus
def Main()
{
}
```

or

```janus
def Main() -> Void
{
}
```

