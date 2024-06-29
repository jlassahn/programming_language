
all: build/bin/moss-cc

build/generated/grammar.c: grammar.bison
	bison -d -o $@ $^

build/bin/moss-cc: \
	build/generated/grammar.c \
	compiler/memory.c \
	compiler/types.c \
	compiler/errors.c \
	compiler/fileio.c \
	compiler/token_defs.c \
	compiler/tokenizer.c \
	compiler/parser_file.c \
	compiler/parser.c \
	compiler/commandargs.c \
	compiler/compile_state.c \
	compiler/moss-cc.c \


	gcc -DYYDEBUG -I . -Wall -o $@ $^

clean:
	rm build/bin/moss-cc
	rm build/generated/*.c
	rm build/generated/*.h


tests: \
	build/bin/map_test \


build/bin/map_test: \
	tests/unit/map_test.c \
	compiler/types.c \

	gcc -DYYDEBUG -I . -Wall -o $@ $^
