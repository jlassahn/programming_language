
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

typedef struct StringBuffer StringBuffer;
struct StringBuffer
{
	String string; // always points to buffer
	int capacity;
	char buffer[];
};

StringBuffer *StringBufferCreateEmpty(int capacity);
StringBuffer *StringBufferFromChars(const char *chars);
StringBuffer *StringBufferFromString(const String *str);
StringBuffer *StringBufferAppendChars(StringBuffer *sb, const char *chars);
StringBuffer *StringBufferAppendString(StringBuffer *sb, const String *str);
StringBuffer *StringBufferAppendBuffer(StringBuffer *sb, const StringBuffer *b);
void StringBufferClear(StringBuffer *sb);
void StringBufferLock(StringBuffer *sb);
void StringBufferFree(StringBuffer *sb);

typedef struct ListEntry ListEntry;
struct ListEntry
{
	ListEntry *prev;
	ListEntry *next;
	void *item;
};

typedef struct List List;
struct List
{
	ListEntry *first;
	ListEntry *last;
};

void ListInsertFirst(List *list, void *item);
void ListInsertLast(List *list, void *item);
void *ListRemoveFirst(List *list);

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

// FIXME way to destroy by iterating through everything, removing and returning

void MapPrint(Map *map);

#endif

