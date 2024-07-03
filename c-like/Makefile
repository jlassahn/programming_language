
#
# General targets
#
all: make_dirs make_compiler make_tests

test: make_dirs make_tests

clean:
	rm -Rf build/*


#
# Outputs by category
#
make_compiler: \
	build/bin/moss-cc

make_tests: \
	build/bin/map_test \

make_dirs: \
	build/generated/ \
	build/obj/compiler/ \
	build/obj/tests/ \
	build/obj/tests/unit/ \
	build/bin/ \


#
# Generated grammar files
#
build/generated/grammar.c: grammar.bison
	bison -d -o $@ $^

build/generated/grammar.o: build/generated/grammar.c
	gcc -c -DYYDEBUG -I . -Wall -o $@ $<


#
# Compiler
#
build/bin/moss-cc: \
	build/generated/grammar.o \
	build/obj/compiler/memory.o \
	build/obj/compiler/types.o \
	build/obj/compiler/errors.o \
	build/obj/compiler/fileio.o \
	build/obj/compiler/token_defs.o \
	build/obj/compiler/tokenizer.o \
	build/obj/compiler/parser_file.o \
	build/obj/compiler/parser.o \
	build/obj/compiler/namespace.o \
	build/obj/compiler/commandargs.o \
	build/obj/compiler/compile_state.o \
	build/obj/compiler/compiler_file.o \
	build/obj/compiler/moss-cc.o \

	gcc -o $@ $^


#
# Individual Tests
#
build/bin/map_test: \
	build/obj/tests/unit/map_test.o \
	build/obj/compiler/memory.o \
	build/obj/compiler/types.o \

	gcc -o $@ $^


#
# Patterns
#
headers := $(wildcard compiler/*.h)

build/obj/%.o: %.c $(headers)
	gcc -c -I . -Wall -o $@ $<

build/%/:
	mkdir -p $@

