
all: build/bin/moss-cc

build/bin/moss-cc: \
	compiler/moss-cc.c \
	compiler/parser_file.c \
	compiler/tokenizer.c \
	compiler/errors.c \

	gcc -I compiler -Wall -o $@ $^

clean:
	rm build/bin/moss-cc

