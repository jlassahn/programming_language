
#ifndef INCLUDED_SEARCH_H
#define INCLUDED_SEARCH_H

/*!
  @file
  @brief Functions for searching for imports, libraries ,etc.
*/

#include "compiler/types.h"

/*!
  @brief An opaque handle for source tree search state.
*/
typedef struct SearchFiles SearchFiles;

/*!
  @brief Start searching for files in a set of search paths.

  Searches multiple base directories for files of the form
  basedir/part/path/filename
  checks if the beginning of the filename matches the string in head,
  and the end of the filename matches one of the string list in tails
  and returns all matches.

  @param basedirs A List of StringBuffer* representing the search paths.
  @param part A string such as "import/" or "lib/" for which subtree to search.
  @param path A string describing the module path inside the tree.
  @param head A string to filter the file against.
  @param tails An array of strings to filter the file against, ending in NULL.
  @return A handle to use for file searches.
*/
SearchFiles *SearchFilesStart(
		List *basedirs,
		const char *part,
		const char *path,
		const char *head,
		const char *tails[]);

/*!
  @brief Return the next file found by the search search.

  The returned StringBuffer is owned by the caller, and must be freed
  when it is no longer needed.  It will contain the complete path
  built up from the basedir.

  @param sf A handle holding the search state.
  @return The next file found, or NULL if no more are available.
*/
StringBuffer *SearchFilesNext(SearchFiles *sf);

/*!
  @brief End the search and free the associated handle.

  @param ds The handle to free.
*/
void SearchFilesEnd(SearchFiles *sf);

#endif

