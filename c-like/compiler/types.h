
#ifndef INCLUDED_TYPES_H
#define INCLUDED_TYPES_H

#include <stdbool.h>

#ifdef _WIN32
#define USE_RESULT _Check_return_
#else
#define USE_RESULT __attribute__((warn_unused_result))
#endif

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

// The StringBufferAppend functions return a new pointer to the modified
// buffer (like realloc does for memory).

USE_RESULT
StringBuffer *StringBufferCreateEmpty(int capacity);

USE_RESULT
StringBuffer *StringBufferFromChars(const char *chars);

USE_RESULT
StringBuffer *StringBufferFromString(const String *str);

USE_RESULT
StringBuffer *StringBufferAppendChars(StringBuffer *sb, const char *chars);

USE_RESULT
StringBuffer *StringBufferAppendString(StringBuffer *sb, const String *str);

USE_RESULT
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
void *MapFind(Map *map, const String *key);
void *MapRemoveFirst(Map *map);
void MapDestroyAll(Map *map);

// void *MapRemove(Map *map, const String *key); // FIXME do we need this?

void MapPrint(Map *map);

#endif

