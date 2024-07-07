
#ifndef INCLUDED_FILEIO_H
#define INCLUDED_FILEIO_H

/*! @file
  @brief Functions for accessing files and directories.

  Several of these functions have different implementations for different OSes.

  Some functions are thin wrappers around C standard library functions, but
  are used instead of calling the standard library so that fake implementations
  can be used for testing.
*/

#include <stdbool.h>

// FIXME change of plan: normalize all paths to contain forward slashes
//       when reading configuration, make Windows file IO functions convert
//       back as needed.

#ifdef _WIN32
#define PATH_SEPARATOR '\\'
#define PATH_SEPARATOR_STRING "\\"
#else
#define PATH_SEPARATOR '/'
#define PATH_SEPARATOR_STRING "/"
#endif

/*!
  @brief An opaque handle for directory search state.
  Will be created by a successful call to DirectorySearchStart()
  and destroyed by DirectorySearchEnd().
*/
typedef struct DirectorySearch DirectorySearch;

/*!
  @brief Start searching for the files contained in the directory path.
  Only regular files, will be found, not subdirectories or special
  files.

  @param path The path of the directory to search file files.
  @return A handle to query for files.
*/
DirectorySearch *DirectorySearchStart(const char *path);

/*!
  @brief Return the next file found by a directory search.

  The returned string is only valid until the next call to a function
  which uses the handle.

  The return is only the file name, not the complete path.

  @param ds A handle holding the search state.
  @return The name of the file.
*/
const char *DirectorySearchNextFile(DirectorySearch *ds);

/*!
  @brief End the directory search and free the associated handle.

  @param ds The handle to free.
*/
void DirectorySearchEnd(DirectorySearch *ds);

/*!
  @brief Check whether the path exists and is a directory.

  @param path The path to check.
  @return true if the directory exists.
*/
bool DoesDirectoryExist(const char *path);

/*!
  @brief Check whether the path exists and is a regular file.

  @param path The path to check.
  @return true if the file exists.
*/
bool DoesFileExist(const char *path);

/*!
  @brief An opaque handle to an open file.
*/
typedef struct OSFile OSFile;

/*!
  @brief Open a file for reading.

  @param path The path to the file to open.
  @return An open file handle.
*/
OSFile *OSFileOpenRead(const char *path);

/*!
  @brief Close a file.

  @param fp The handle of the file to close.
*/
void OSFileClose(OSFile *fp);

/*!
  @brief Get the size of an open file in bytes.

  @param fp The handle of the file.
  @return The length of the file in bytes.
*/
long OSFileGetSize(OSFile *fp);

/*!
  @brief Read bytes from a file.

  @param fp The handle of the file.
  @param data_out A buffer to receive the data.
  @param max_bytes The maximum number of bytes to read.
  @return The actual number of bytes read.
*/
long OSFileRead(OSFile *fp, void *data_out, long max_bytes);

// maybe we'll want write calls as well soon...
// OSFile *OSFileOpenWrite(const char *path);
// long OSFileWrite(OSFile *fp, void *data, long bytes);

#endif

