
#include "compiler/types.h"
#include "compiler/memory.h"
#include "compiler/exit_codes.h"
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h> // for MapPrint


bool StringEquals(const String *a, const String *b)
{
	if (a->length != b->length)
		return false;

	return (memcmp(a->data, b->data, a->length) == 0);
}


static void StringBufferAssertUnlocked(StringBuffer *sb)
{
	if (sb->capacity <= 0)
	{
		fprintf(stderr, "MODIFYING LOCKED StringBuffer: %s\n", sb->buffer);
		exit(EXIT_SOFTWARE);
	}
}

USE_RESULT
static StringBuffer *StringBufferEnsureCapacity(StringBuffer *sb, int capacity)
{
	if (sb->capacity >= capacity)
		return sb;

	// grow by some minimum size rather than byte-by-byte
	if (capacity < sb->capacity + 16)
		capacity = sb->capacity + 16;

	sb = Realloc(sb, sizeof(StringBuffer) + capacity);
	sb->string.data = sb->buffer;
	sb->capacity = capacity;
	return sb;
}

USE_RESULT
StringBuffer *StringBufferCreateEmpty(int capacity)
{
	// default to some reasonable initial size
	if (capacity <= 0)
		capacity = 200;

	StringBuffer *sb = Alloc(sizeof(StringBuffer) + capacity);
	sb->string.data = sb->buffer;
	sb->string.length = 0;
	sb->capacity = capacity;

	return sb;
}

USE_RESULT
StringBuffer *StringBufferFromChars(const char *chars)
{
	int length = strlen(chars);
	int capacity = length + 1;
	StringBuffer *sb = Alloc(sizeof(StringBuffer) + capacity);
	sb->string.data = sb->buffer;
	sb->string.length = length;
	sb->capacity = capacity;
	strcpy(sb->buffer, chars);

	return sb;
}

USE_RESULT
StringBuffer *StringBufferAppendChars(StringBuffer *sb, const char *chars)
{
	StringBufferAssertUnlocked(sb);

	int length = strlen(chars);
	sb = StringBufferEnsureCapacity(sb, sb->string.length + length + 1);

	memcpy(sb->buffer + sb->string.length, chars, length);
	sb->string.length += length;
	sb->buffer[sb->string.length] = 0;

	return sb;
}

USE_RESULT
StringBuffer *StringBufferAppendString(StringBuffer *sb, const String *str)
{
	StringBufferAssertUnlocked(sb);

	int length = str->length;
	sb = StringBufferEnsureCapacity(sb, sb->string.length + length + 1);

	memcpy(sb->buffer + sb->string.length, str->data, length);
	sb->string.length += length;
	sb->buffer[sb->string.length] = 0;

	return sb;
}

USE_RESULT
StringBuffer *StringBufferAppendBuffer(StringBuffer *sb, const StringBuffer *b)
{
	return StringBufferAppendString(sb, &b->string);
}

void StringBufferLock(StringBuffer *sb)
{
	if (sb->capacity > 0)
		sb->capacity = -sb->capacity;
}

void StringBufferClear(StringBuffer *sb)
{
	StringBufferAssertUnlocked(sb);
	sb->buffer[0] = 0;
	sb->string.length = 0;
}

void StringBufferFree(StringBuffer *sb)
{
	Free(sb);
}

void ListInsertFirst(List *list, void *item)
{
	ListEntry *entry = Alloc(sizeof(ListEntry));
	entry->prev = NULL;
	entry->next = list->first;
	entry->item = item;

	if (list->first)
		list->first->prev = entry;

	list->first = entry;
	if (list->last == NULL)
		list->last = entry;
}

void ListInsertLast(List *list, void *item)
{
	ListEntry *entry = Alloc(sizeof(ListEntry));
	entry->prev = list->last;
	entry->next = NULL;
	entry->item = item;

	if (list->last)
		list->last->next = entry;

	list->last = entry;
	if (list->first == NULL)
		list->first = entry;
}

void *ListRemoveFirst(List *list)
{
	if (list->first == NULL)
		return NULL;

	ListEntry *entry = list->first;
	list->first = entry->next;
	if (entry->next)
		entry->next->prev = NULL;

	if (list->first == NULL)
		list->last = NULL;

	void *item = entry->item;
	Free(entry);
	return item;
}

// if a Map bin has more than this, make a subtable
const int MAX_COLLISIONS = 4;

static uint32_t HashString(const String *s)
{
	// This is the DJB2 algorithm, a traditional hash function
	// with OK properties that's easy to compute.

	uint32_t hash = 5381;
	for (int i=0; i<s->length; i++)
		hash = hash*33 + s->data[i];

	return hash;
}


bool MapInsert(Map *map, const String *key, void *value)
{
	uint32_t hash = HashString(key);
	int i = (hash >> map->shift) & 255;
	HashBin *bin = &map->bins[i];
	if (bin->subtable)
	{
		bool ret = MapInsert(bin->subtable, key, value);
		if (ret)
		{
			bin->count ++;
			map->count ++;
		}
		return ret;
	}

	for (HashEntry *entry = bin->list; entry != NULL; entry = entry->next)
	{
		if (StringEquals(&entry->key, key))
			return false;
	}

	if ((bin->count > MAX_COLLISIONS) && (map->shift < 24))
	{
		Map *submap = Alloc(sizeof(Map));
		submap->shift = map->shift + 8;

		HashEntry *list = bin->list;
		bin->list = NULL;
		bin->subtable = submap;

		// this is inefficient -- frees and reallocates list entries
		while (list != NULL)
		{
			HashEntry *next = list->next;
			MapInsert(submap, &list->key, list->value);
			Free(list);
			list = next;
		}

		MapInsert(submap, key, value);
		bin->count ++;
		map->count ++;
		return true;
	}

	HashEntry *entry = Alloc(sizeof(HashEntry));
	entry->key = *key;
	entry->value = value;
	entry->next = bin->list;
	entry->prev = NULL;
	bin->list = entry;
	bin->count ++;
	map->count ++;
	return true;
}

void *MapFind(Map *map, const String *key)
{
	uint32_t hash = HashString(key);

	while (true)
	{
		int i = (hash >> map->shift) & 255;
		HashBin *bin = &map->bins[i];
		if (bin->subtable)
		{
			map = bin->subtable;
			continue;
		}
		for (HashEntry *entry = bin->list; entry != NULL; entry = entry->next)
		{
			if (StringEquals(&entry->key, key))
				return entry->value;
		}
		return NULL;
	}
}

void MapDestroyAll(Map *map)
{
	for (int i=0; i<256; i++)
	{
		HashBin *bin = &map->bins[i];
		HashEntry *entry = bin->list;
		while (entry != NULL)
		{
			HashEntry *next_entry = entry->next;
			Free(entry);
			entry = next_entry;
		}
		bin->list = NULL;

		if (bin->subtable)
		{
			MapDestroyAll(bin->subtable);
			Free(bin->subtable);
		}
		bin->subtable = NULL;
		bin->count = 0;
	}
	map->count = 0;
}


static void Indent(int depth)
{
	for (int i=0; i<depth; i++)
		printf("  ");
}

static void MapPrintDepth(Map *map, int depth)
{
	Indent(depth); printf("map count = %d\n", map->count);
	for (int i=0; i<256; i++)
	{
		HashBin *bin = &map->bins[i];
		if (bin->count == 0)
			continue;
		Indent(depth); printf("bin %d: count %d\n", i, bin->count);
		if (bin->subtable)
			MapPrintDepth(bin->subtable, depth+1);
	}
}

void MapPrint(Map *map)
{
	MapPrintDepth(map, 0);
}

