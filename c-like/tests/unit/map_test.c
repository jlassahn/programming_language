
#include "compiler/types.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define CHECK(x) DoCheck(x, __LINE__, __FILE__)

int fail_count = 0;

void DoCheck(bool x, int line, const char *file)
{
	if (!x)
	{
		printf("FAILED line %d file %s\n", line ,file);
		fail_count ++;
	}
}

void MakeString(String *s, const char *txt)
{
	s->data = txt;
	s->length = strlen(txt);
}

int main(void)
{
	static Map map;
	String s;

	CHECK(true);

	MakeString(&s, "Hello");
	CHECK(MapInsert(&map, &s, (void *)s.data));
	CHECK(!MapInsert(&map, &s, (void *)s.data));
	MakeString(&s, "Hello!");
	CHECK(MapInsert(&map, &s, (void *)s.data));
	CHECK(map.count == 2);

	MakeString(&s, "xxx");
	CHECK(NULL == MapFind(&map, &s));
	MakeString(&s, "Hello");
	CHECK(NULL != MapFind(&map, &s));
	CHECK(strcmp(MapFind(&map, &s), "Hello") == 0);

	for (int i=0; i<1000; i++)
	{
		char *data = malloc(5);
		sprintf(data, "%.4d", i);
		MakeString(&s, data);
		CHECK(MapInsert(&map, &s, NULL));
	}
	MakeString(&s, "Hello");
	CHECK(NULL != MapFind(&map, &s));
	CHECK(strcmp(MapFind(&map, &s), "Hello") == 0);

	MakeString(&s, "Later...");
	CHECK(MapInsert(&map, &s, "LATER"));
	CHECK(NULL != MapFind(&map, &s));
	CHECK(strcmp(MapFind(&map, &s), "LATER") == 0);

	CHECK(map.count == 1003);

	MapDestroyAll(&map);
	CHECK(map.count == 0);

	return fail_count;
}

