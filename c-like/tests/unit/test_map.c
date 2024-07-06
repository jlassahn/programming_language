
#include "compiler/types.h"
#include "compiler/memory.h"
#include "tests/unit/unit_test.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

static void MakeString(String *s, const char *txt)
{
	s->data = txt;
	s->length = strlen(txt);
}

void TestMap(void)
{
	static Map map;
	String s;

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

	char *values[1000];
	for (int i=0; i<1000; i++)
	{
		char *data = Alloc(5);
		sprintf(data, "%.4d", i);
		MakeString(&s, data);
		CHECK(MapInsert(&map, &s, NULL));
		values[i] = data;
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

	for (int i=0; i<1000; i++)
		Free(values[i]);
}

