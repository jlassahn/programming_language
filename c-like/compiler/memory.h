
#ifndef INCLUDED_MEMORY_H
#define INCLUDED_MEMORY_H

void *Alloc(int size);
void Free(void *p);
void *Realloc(void *p, int size);
int AllocCount(void);

#endif

