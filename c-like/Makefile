
all: build/bin/moss-cc

build/generated/grammar.c: grammar.bison
	bison -d -o $@ $^

build/bin/moss-cc: \
	build/generated/grammar.c \
	compiler/moss-cc.c \
	compiler/parser_file.c \
	compiler/tokenizer.c \
	compiler/token_defs.c \
	compiler/parser.c \
	compiler/errors.c \
	compiler/types.c \
	compiler/commandargs.c \
	compiler/fileio.c \
	compiler/compile_state.c \


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
