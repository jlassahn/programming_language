
#
# General targets
#
all: make_compiler make_tests

test: make_compiler make_tests run_tests

clean:
	rm -Rf build/*


#
# Outputs by category
#
make_compiler: \
	build/bin/moss-cc

make_tests: \
	build/bin/all_tests \
	build/bin/scan_tests \
	build/bin/test_fileio \


run_tests:
	build/bin/test_fileio
	build/bin/all_tests
	build/bin/scan_tests

#
# Generated grammar files
#
build/generated/grammar.c: grammar.bison
	@mkdir -p $(@D)
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
	build/obj/compiler/stringtypes.o \
	build/obj/compiler/errors.o \
	build/obj/compiler/fileio.o \
	build/obj/compiler/token_defs.o \
	build/obj/compiler/tokenizer.o \
	build/obj/compiler/parser_file.o \
	build/obj/compiler/parser_node.o \
	build/obj/compiler/parser_symbols.o \
	build/obj/compiler/parser.o \
	build/obj/compiler/namespace.o \
	build/obj/compiler/commandargs.o \
	build/obj/compiler/compile_state.o \
	build/obj/compiler/compiler_file.o \
	build/obj/compiler/search.o \
	build/obj/compiler/pass_configure.o \
	build/obj/compiler/pass_search_and_parse.o \
	build/obj/compiler/moss-cc.o \

	@mkdir -p $(@D)
	gcc -o $@ $^


#
# Individual Tests
#
build/bin/test_fileio: \
	build/obj/tests/unit/unit_test.o \
	build/obj/compiler/memory.o \
	build/obj/compiler/fileio.o \
	build/obj/tests/test_fileio.o \

	@mkdir -p $(@D)
	gcc -o $@ $^

build/bin/all_tests: \
	build/obj/compiler/types.o \
	build/obj/compiler/memory.o \
	build/obj/compiler/stringtypes.o \
	build/obj/compiler/namespace.o \
	build/obj/compiler/parser_file.o \
	build/obj/compiler/parser_node.o \
	build/obj/compiler/parser_symbols.o \
	build/obj/compiler/compiler_file.o \
	build/obj/compiler/compile_state.o \
	build/obj/compiler/search.o \
	build/obj/compiler/commandargs.o \
	build/obj/compiler/pass_configure.o \
	build/obj/compiler/pass_search_and_parse.o \
	build/obj/tests/unit/unit_test.o \
	build/obj/tests/unit/fake_errors.o \
	build/obj/tests/unit/fake_directory.o \
	build/obj/tests/unit/fake_parser.o \
	build/obj/tests/unit/fake_nodes.o \
	build/obj/tests/unit/all_tests.o \
	build/obj/tests/unit/test_fake_dir.o \
	build/obj/tests/unit/test_map.o \
	build/obj/tests/unit/test_stringtypes.o \
	build/obj/tests/unit/test_compiler_file.o \
	build/obj/tests/unit/test_search.o \
	build/obj/tests/unit/test_pass_configure.o \
	build/obj/tests/unit/test_pass_search_and_parse.o \

	@mkdir -p $(@D)
	gcc -o $@ $^

build/bin/scan_tests: \
	build/generated/grammar.o \
	build/obj/compiler/fileio.o \
	build/obj/compiler/types.o \
	build/obj/compiler/memory.o \
	build/obj/compiler/stringtypes.o \
	build/obj/compiler/search.o \
	build/obj/compiler/namespace.o \
	build/obj/compiler/parser_file.o \
	build/obj/compiler/parser_node.o \
	build/obj/compiler/parser_symbols.o \
	build/obj/compiler/parser.o \
	build/obj/compiler/token_defs.o \
	build/obj/compiler/tokenizer.o \
	build/obj/compiler/compiler_file.o \
	build/obj/compiler/compile_state.o \
	build/obj/compiler/commandargs.o \
	build/obj/compiler/pass_configure.o \
	build/obj/compiler/pass_search_and_parse.o \
	build/obj/tests/unit/unit_test.o \
	build/obj/tests/unit/fake_errors.o \
	build/obj/tests/scan_tests/scan_tests.o \
	build/obj/tests/scan_tests/good_filenames.o \
	build/obj/tests/scan_tests/bad_filenames.o \
	build/obj/tests/scan_tests/relative_paths.o \


	@mkdir -p $(@D)
	gcc -o $@ $^

#
# Patterns
#
headers := $(wildcard compiler/*.h)

build/obj/%.o: %.c $(headers)
	@mkdir -p $(@D)
	gcc -c -I . -Wall -o $@ $<


