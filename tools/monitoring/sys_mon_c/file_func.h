// file_func.h

#ifndef FUNCTIONS_FILE_INCLUDED
#define FUNCTIONS_FILE_INCLUDED

#include <stdbool.h>
#include <libgen.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <dirent.h>
#include <ctype.h>
#include <time.h>
#include <regex.h>
#include <unistd.h>

#define PATH_LENGTH 256
#define UNAVAILABLE_STR "n/a"
#define UNAVAILABLE_STR_INT "-1"
#define UNAVAILABLE_STR_UNS "0"

// Local error storage
char ERROR_MESSAGE[384];
bool ERROR_IS_SET;
bool ERROR_IS_EOF;

/*****************
 * Miscellaneous
 *****************/
// 2d char array memory allocation
char **calloc_2d_array(char **array, int rows, int cols);

// ask for [ENTER] to continue
void press_enter_continue();

/*****************
 * Time functions
 *****************/
typedef struct timespec timespec;

timespec time_spec_get();

void sleep_milisec(long int milisec);

double *time_spec_sec_sub( timespec* ts1,  timespec* ts2, double *value);

double *time_spec_seconds( timespec* ts, double *value);

/*******************
 * Handling ERRORS
 *******************/
// Retrieve and clean error information if there is.
char *internal_error_get(char *error);

// clear local error information.
void internal_error_clear();

// set error information.
void internal_error_set(char *error);

/*********************
 * Strings functions
 *********************/
 // regular expression word matching
int reg_match(char *string, char *pattern);

// defining comparator function as per the requirement
int compare_func(const void* a, const void* b);

// function to sort the array
void sort(char* arr[], int n);

// set 'char_array' to 0
char *clean_char_array(char *char_array);

// cut at 1st ':'
char *cut_at_first(char *str, char *cut_at);

// remove trailing '\n'
char *remove_lf(char *str);

// Trim white space of given string
char *trim_space(char *str);

/*******************
 * Files functions
 *******************/
// get symlink endpoint
char *symlink_endpoint(char *path, char *endpoint);

// Lists all files and sub-directories at given path.
// and fill 'files_list' 2 dimensional array with.
bool list_files(const char *path, char **files_list, int *count);

// get directory name from path (posix)
char *dirname_get(char *path);

// get base name from path (posix)
char *basename_get(char *path);

// Read line at current offset and check whether 'contain'
// is present or not, if there is, true is returned.
// Otherwise, file is processed again from start offset.
// If nothing is found at the second pass, offset is restored
// to initial position and false is returned.
char *line_contain(FILE *fp, char *contain, char *line, int size);

// Read error checking. On error, internal flag 'ERROR_IS_SET'
// is toggled and 'ERROR_MESSAGE' is filled with information.
// get_internal_error() must be called to reset 'ERROR_IS_SET'
// flag. The file is closed before returning.
bool is_correctly_read(FILE *fp, char path[], int was_read, int expected);

// Read name value pairs from a file
char *get_paired_name_value(FILE *fp, char sep[], char *val_name);

// read file that contain single line and retrieve value
// to 'dest', 'fmt' means format like 'printf' command.
bool fgets_void(char *path, char *fmt, void *dest);

// walk the given stream to the next line
bool seek_to_next_line(FILE *fp);

/********************
 * NOT USED ANYMORE
 *******************/
// Read error/EOF checking. Break on error, return EOF
// wether file is ended.
int is_correctly_read_or_eof(FILE *fp, char path[], int read, int expected);

// Will fill 'start_end_off' with start, end offsets.
// Lines start at 0.
void file_seek_line(FILE *fp, int line, long *start_end_off[2]);

#endif
