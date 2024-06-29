
#include "compiler/memory.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>

static int allocation_count = 0;

int AllocCount(void)
{
	return allocation_count;
}

void *Alloc(int size)
{
	void *buf = malloc(size);
	if (buf == NULL)
	{
		fprintf(stderr, "OUT OF MEMORY\n");
		exit(ENOMEM);
	}
	memset(buf, 0, size);
	allocation_count ++;
	return buf;
}

void Free(void *p)
{
	free(p);
	allocation_count --;
}

void *Realloc(void *p, int size)
{
	void *buf = realloc(p, size);
	if (buf == NULL)
	{
		fprintf(stderr, "OUT OF MEMORY\n");
		exit(ENOMEM);
	}
	return buf;
}

