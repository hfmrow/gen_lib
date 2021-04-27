// file_func.c

#include "file_func.h"

/*****************
 * Miscellaneous
 *****************/
// 2d char array memory allocation
char **calloc_2d_array(char **array, int rows, int cols)
{
	array = (char **)calloc(rows, sizeof(char*));
	for (int i = 0; i < rows; i++) {
		array[i] = (char*)malloc(cols * sizeof(char));
	}
	return array;
}

// ask for [ENTER] to continue
void press_enter_continue()
{
	printf("Press [Enter] to continue.\n");
	//while(getchar()!='\n');  /* clean stdin */
	getchar(); // wait for ENTER
}

timespec time_spec_get()
{
	timespec ts = {0,0};
	int r = clock_gettime(CLOCK_REALTIME, &ts);
	if (r == -1) {
		ts.tv_nsec = -1;
		ts.tv_sec = -1;
	}
	return ts;
}

void sleep_milisec(long int milisec)
{
	long int tmp_ms = milisec;
	timespec ts = {0,0};
	if (milisec >= 1000) {
		tmp_ms = (milisec % 1000);
		sleep((unsigned int)milisec / 1000);
	}
	ts.tv_nsec = tmp_ms * 1000000;
	nanosleep(&ts, NULL);
}

double *time_spec_seconds(timespec* ts, double *value)
{
	//	1000000000 ns = 1 s
	*value = (double) ts->tv_sec + (double) ts->tv_nsec * 1.0e-9;
	return value;
}

double *time_spec_sec_sub(timespec* ts1, timespec* ts2, double *value)
{
	double v1, v2;
	*value = *time_spec_seconds(ts1, &v1) - *time_spec_seconds(ts2, &v2);
	return value;
}

/*******************
 * Handling ERRORS
 *******************/
// Retrieve and clean error information if there is.
char *internal_error_get(char *error)
{
	sprintf(error, "%s", ERROR_MESSAGE);
	internal_error_clear();
	return error;
}

// clear local error information.
void internal_error_clear()
{
	memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
	ERROR_IS_SET = false;
	ERROR_IS_EOF = false;
}

// set error information.
void internal_error_set(char *error)
{
	memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
	sprintf(ERROR_MESSAGE, "%s", error);
	ERROR_IS_SET = true;
}

/*********************
 * Strings functions
 *********************/
// regular expression word matching
int reg_match(char *string, char *pattern)
{
	int    status;
	regex_t    re;

	//char *tofind = calloc(256, sizeof(char));
	//sprintf(tofind, "\\b%s\\b", pattern);

	if (regcomp(&re, pattern, REG_EXTENDED|REG_NOSUB) != 0) {
		return(0);      /* Report error. */
	}
	status = regexec(&re, (const char*)string, (size_t) 0, NULL, 0);
	regfree(&re);
	if (status != 0) {
		return(0);      /* Report error. */
	}
	return(1);
}

// defining comparator function as per the requirement
int compare_func(const void* a, const void* b)
{
	// setting up rules for comparison
	return strcmp(*(char**)a, *(char**)b);
}

// function to sort the array
void sort(char* arr[], int n)
{
	// calling qsort function to sort the array
	// with the help of Comparator
	qsort((const char**)arr, n, sizeof(const char*), compare_func);
}

// set 'char_array' to 0
char *clean_char_array(char *char_array)
{
	int size = sizeof(char_array);
	memset(char_array, 0, size);
	return char_array;
}

// cut at 1st ':'
char *cut_at_first(char *str, char *cut_at)
{
	str[strcspn(str, cut_at )] = '\0';
	return str;
}

// remove trailing '\n'
char *remove_lf(char *str)
{
	str[strcspn ( str, "\n" )] = '\0';
	return str;
}

// Trim white space of given string
char *trim_space(char *str)
{
	char *end;
	// skip leading whitespace
	while (isspace(*str)) {
		str++;
	}
	// remove trailing whitespace
	end = str + strlen(str) - 1;
	while (end > str && isspace(*end)) {
		end--;
	}
	// write null character
	*(end+1) = '\0';
	return str;
}

/*******************
 * Files functions
 *******************/
// get symlink endpoint
char *symlink_endpoint(char *path, char *endpoint)
{
	char buf[1024];
	ssize_t len;
	if ((len = readlink(path, buf, sizeof(buf)-1)) != -1)
		buf[len] = '\0';
	else
		return NULL;

	strcpy(endpoint, buf);
	return endpoint;
}

// Lists all files and sub-directories at given path.
// and fill 'files_list' 2 dimensional array with.
bool list_files(const char *path, char **files_list, int *count_out)
{
	struct dirent *dp;
	DIR *dir = opendir(path);
	int count = 0;

	if (!dir)
		return false;

	while ((dp = readdir(dir)) != NULL) {
		if (!strcmp(dp->d_name, ".") || !strcmp(dp->d_name, ".."))
			continue;
		sprintf(files_list[count], "%s/%s", path, dp->d_name);
		count++;
	}
	closedir(dir);
	*count_out = count;
	return true;
}

// get directory name from path
char *dirname_get(char *path)
{
	return dirname(strdup(path));
}

// get base name from path
char *basename_get(char *path)
{
	return basename(strdup(path));
}

// Read line at current offset and check whether 'contain'
// is present or not, if there is, true is returned.
// Otherwise, file is processed again from start offset.
// If nothing is found at the second pass, offset is restored
// to initial position and false is returned.
char *line_contain(FILE *fp, char *contain, char *line, int size)
{
	bool first_pass_done = false;
	char *res;
	fpos_t initial_offset, last_line_offset;

	// build regexp to match word boundary
	regex_t regx;
	char reg_str[64];// = malloc(256 * sizeof(char));
	sprintf(reg_str, "\\b%s\\b", contain);
	if (regcomp(&regx, reg_str, REG_EXTENDED|REG_NOSUB) != 0)
		printf("Error regexp pattern: %s\n", reg_str);

	// store current offset
	fgetpos(fp, &initial_offset);
	while (true) {
		fgetpos(fp, &last_line_offset);
		if ((res = fgets(line, size, fp)) != NULL) {
			if (regexec(&regx, line, (size_t) 0, NULL, 0) == 0) {
				fsetpos(fp, &last_line_offset);
				break;
			}
		} else if (!first_pass_done) {
			// nothing found, so mark the first pass was already
			// done and rewind to start of file.
			first_pass_done = true;
			rewind(fp);
		} else {
			// set offset to the initial position and return NULL.
			fsetpos(fp, &initial_offset);
			res = NULL;
			break;
		}
	}
	regfree(&regx);
	return res;
}

// Read error checking. On error, internal flag 'ERROR_IS_SET'
// is toggled and 'ERROR_MESSAGE' is filled with information.
// get_internal_error() must be called to reset 'ERROR_IS_SET'
// flag. The file is closed before returning.
bool is_correctly_read(FILE *fp, char *path, int was_read, int expected)
{
	ERROR_IS_SET = false;
	ERROR_IS_EOF = false;
	if (was_read == expected)
		return true;

	if (was_read == EOF) {
		ERROR_IS_EOF = true;
		return false;
	}
	if (fp != NULL) {
		fclose(fp);
		fp = NULL;
	}
	ERROR_IS_SET = true;
	sprintf(ERROR_MESSAGE, "Error reading file [%s], read: %d, expected: %d variable(s)\n", path, was_read, expected);
	return false;
}

// Read name value pairs from a file.
// preserve file offset on exit if nothing was found.
// usage for
// numeric:
// 			sscanf(get_paired_name_value(fp, ":", "Pid"), "%lu", &tmp_store_file->pid);
// string with dynamic allocation:
//			char *tmp_str = get_paired_name_value(fp, ":", "Name");
//			tmp_store_file->name = (char*)malloc(sizeof(char) * (strlen(tmp_str) + 1 ));
//			strcpy(tmp_store_file->name, tmp_str);
char *get_paired_name_value(FILE *fp, char *sep, char *val_name)
{
	fpos_t initial_offset;
	int line_length = 4096;
	char line[line_length];
	bool first_pass_done = false;

	// store current file offset
	fgetpos(fp, &initial_offset);


	while ( !first_pass_done) { // walk file
		while (fgets(line, line_length, fp) != NULL)
			if ( !strcmp(val_name, trim_space(strtok(line, sep))))
				return trim_space(strtok(NULL, sep));
		first_pass_done = true;
	}
	// restore file offset on exit if nothing was found
	fsetpos(fp, &initial_offset);
	return NULL;
}

// read file that contain single line and retrieve value
// to 'dest', 'fmt' means format like 'fscanf' command.
bool fgets_void(char *path, char *fmt, void *dest)
{
	bool ret = false;
	int line_length = 128;
	char line[line_length];
	FILE *fp;
	fp = fopen(path, "r");
	if (fp == NULL) { // the file does not exist, so use the defaults
		if (strstr(fmt, "d"))
			sscanf(UNAVAILABLE_STR_INT, fmt, dest);
		else if (strstr(fmt, "f"))
			sscanf(UNAVAILABLE_STR_INT, fmt, dest);
		else if (strstr(fmt, "u"))
			sscanf(UNAVAILABLE_STR_UNS, fmt, dest);
		else if (strstr(fmt, "s"))
			strcpy(dest, UNAVAILABLE_STR);
		else
			printf("fgets_void: Unable to parse, %s\n", fmt);
	} else {
		if (!strcmp(fmt, "%s")) {
			if (fgets(line, line_length, fp)) {
				line[strcspn ( line, "\n" )] = '\0'; // remove '\n'
				strcpy(dest, line);
				ret = true;
			}
		} else if (fscanf(fp, fmt, dest) == 1)
			ret = true;
		fclose(fp);
	}
	return ret;
}

// walk the given stream to the next line (\n)
bool seek_to_next_line(FILE *fp)
{
	char ch;
	while ((ch = fgetc(fp)) != EOF)
		if(ch == '\n')
			return true;

	return false;
}

/********************
 * NOT USED ANYMORE
 *******************/
// Read error/EOF checking. Break on error, return EOF
// wether file is ended.
int is_correctly_read_or_eof(FILE *fp, char path[], int read, int expected)
{
	if (read == expected)
		return true;
	if  (read == EOF)
		return EOF;

	fclose(fp);
	printf("Error, reading: [%s], read: %d, expected: %d variable(s)\n", path, read, expected);
	exit(EXIT_FAILURE);

	return false;
}

// Will fill 'start_end_off' with start, end offsets.
// Lines count start at 0.
void file_seek_line(FILE *fp, int line, long *start_end_off[2])
{
	char ch;
	int curr_line = 0;
	long last_off;

	rewind(fp);
	while ((ch = fgetc(fp)) != EOF) {

		if(ch == '\n') {
			last_off = ftell(fp);
			curr_line++;
			if (curr_line == line) {
				start_end_off[1] = &last_off;
				break;
			}
			last_off++;
			start_end_off[0] = &last_off;
		}
	}
}
