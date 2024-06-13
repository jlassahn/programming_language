
#ifndef INCLUDED_TYPES_H
#define INCLUDED_TYPES_H

#include <stdbool.h>

typedef struct String String;
struct String
{
	const char *data;
	int length;
};

bool StringEquals(const String *a, const String *b);
// int StringCompare(const String *a, const String *b); // FIXME do we need?

typedef struct Map Map;
typedef struct HashEntry HashEntry;
typedef struct HashBin HashBin;

struct HashEntry
{
	HashEntry *prev;
	HashEntry *next;
	String key;
	void *value;
};

struct HashBin
{
	int count;
	// either list or subtable is NULL
	HashEntry *list;
	Map *subtable;
};

struct Map
{
	int count;
	int shift;
	HashBin bins[256];
};


// copies the String struct, but not the backing char array.
bool MapInsert(Map *map, const String *key, void *value);
// void *MapRemove(Map *map, const String *key); // FIXME do we need this?
void *MapFind(Map *map, const String *key);
void MapDestroyAll(Map *map);

void MapPrint(Map *map);

#endif

