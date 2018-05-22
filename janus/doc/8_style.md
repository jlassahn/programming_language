
# Documentation and Style


Recommended style for identifiers:

* variables are lowercase with underscores, a_local_variable
* functions and methods are capitalized, AFunctionOrMethod
* symbolic constants are all caps with underscores, A_SYMBOLIC_CONSTANT
* module names are all lower case, some.module

File organization for reusable code

each module has a .janus file in /interfaces named for the module name.
It contains only public definitions, and should be readable as brief
documentation for the module.

each module has private implementation in a .jsrc file in /source named after
the module.

If a module needs enough code that splitting it into multiple files is helpful
other files are in a subdirectory of /source with .janus files to be imported
bu the main file, and possibly related .jsrc files.

'''
interfaces/
	module.janus
source/
	module.jsrc
	module/
		other_stuff.janus
		other_stuff.jsrc
'''


## Inline Documentation

'''
def SomeFunction(thing Thing, stuff Stuff) -> MoreStuff;
#{
 Does a Thing with some Stuff
 }


def SomeFunction(thing Thing, stuff Stuff) -> MoreStuff;
#
# Does a Thing with some Stuff
#

'''

