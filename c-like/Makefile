
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

	gcc -DYYDEBUG -I . -Wall -o $@ $^

clean:
	rm build/bin/moss-cc

