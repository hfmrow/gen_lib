// sys_smaps.c

#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include "sys_smaps.h"
#include "file_func.h"

/********
 * smaps
 ********/
// update/get all values from 'smaps' file.
bool smaps_update(int i_pid, smaps *store_smaps, bool convert_to_bytes)
{
	FILE *fp;

	// retrieve '/proc/[pid]/smaps' data
	char path[PATH_LENGTH];
	sprintf(path, "/%s/%d/%s", "proc", i_pid, "smaps");
	if ((fp = fopen(path, "r")) != NULL) {

		store_smaps->count = 0;
		while(smap_get(fp, path, &store_smaps->details[store_smaps->count], convert_to_bytes)) {

			// increase memory allocation if limit has been reached.
			if (store_smaps->count == store_smaps->internal_alloc) {
				store_smaps->details =
				    realloc(store_smaps->details, (sizeof(smap) + sizeof(map_header)) *
				            (store_smaps->count + store_smaps->internal_alloc_step));

				store_smaps->internal_alloc += store_smaps->internal_alloc_step;
			}
			store_smaps->count++;
		}

		if (ERROR_IS_SET)
			return false;

		store_smaps->count--;
		return true;
	}

	sprintf(ERROR_MESSAGE, "Error opening file [%s]", path);
	ERROR_IS_SET = true;
	return false;
}

// Create and initialize a new smaps structure.
smaps *smaps_new(int count)
{
	smaps *store_smaps = malloc(sizeof(smaps));
	assert(store_smaps != NULL);
	store_smaps->internal_alloc = count;
	store_smaps->internal_alloc_step = count;

	store_smaps->details = malloc((sizeof(smap) + sizeof(map_header)) * count);
	assert(store_smaps->details != NULL);

	return store_smaps;
}

// Free a smaps structure.
void smaps_free(smaps *smaps_to_free)
{
	free(smaps_to_free->details);
	smaps_to_free->details = NULL;
	free(smaps_to_free);
	smaps_to_free = NULL;
}

// init and get smaps (used on 1st time request)
smaps *smaps_get(int pid, int size, bool convert_to_bytes)
{
	smaps *store_smaps = smaps_new(size);
	internal_error_clear();

	if (smaps_update(pid, store_smaps, convert_to_bytes)) {
		// reallocate memory to fit the real count of entries.
		store_smaps->details =
		    realloc(store_smaps->details, (sizeof(smap) + sizeof(map_header)) * store_smaps->count);
		return store_smaps;
	}

	return NULL;
}

// only used for binding structure (with golang for example)
// retrieve specific smap from array.
smap *smaps_get_smap(smaps *smaps_to_smap, int ask_for)
{
	return &smaps_to_smap->details[ask_for];
}

// Get single 'smap' values from 'smaps' file.
bool smap_get(FILE *fp, char path[], smap *store_smap, bool convert_to_bytes)
{
	// there are 23 elements to read, 20 of which will only be read by the loop.
	const int lines_to_read = 20;
	unsigned long long values[lines_to_read];
	char fake_store[50] = {0};
	int count = 0;

	// Getting 'mmap' header.
	if ( !get_map_header(fp, path, SMAPS, &store_smap->header))
		return false;

	// get the second argument of the line which contains 3 arguments.
	if ( !convert_to_bytes) {
		while (is_correctly_read(
		           fp,
		           path,
		           fscanf(fp, "%s%llu%[^\n]", fake_store, &values[count], fake_store), 3))
			if (++count == lines_to_read)
				break;

		if (ERROR_IS_SET)
			return false;

	} else { /*  Values converted to bytes instead of Kb. */
		unsigned long long tmp_value;
		while (is_correctly_read(
		           fp,
		           path,
		           fscanf(fp, "%s%llu%[^\n]", fake_store, &tmp_value, fake_store), 3)) {
			values[count] = tmp_value * 1024;
			if (++count == lines_to_read)
				break;
		}

		if (ERROR_IS_SET)
			return false;
	}
	// assignate value to structure (quick method).
	memcpy(&store_smap->size, &values, sizeof(store_smap->size) * lines_to_read);

	// read 'THPeligible' & 'VmFlags' separatly.
	if ( ! is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%d", fake_store, &store_smap->protection_key), 2))
		return false;

	if ( ! is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%[^\n]", fake_store, store_smap->vm_flags), 2))
		return false;

	return true;
}

// Create and initialize a new smap structure.
smap *smap_new()
{
	smap *new_smap = malloc(sizeof(smap));
	assert(new_smap != NULL);
	return new_smap;
}

// Free a smap structure.
void smap_free(smap *smap_to_free)
{
	free(smap_to_free);
	smap_to_free = NULL;
}

/***************
 * smaps_rollup
 ***************/
// Read values starting from 'smaps_rollup' file.
bool get_smaps_rollup(int i_pid, smaps_rollup *store_smaps_rollup, bool convert_to_bytes)
{
	FILE *fp;
	char path[PATH_LENGTH];
	unsigned long long values[30] = {0};
	char fake_store[50] = {0};
	int count = 0;

	// Getting '/proc/[pid]/smaps_rollup' data
	sprintf(path, "/%s/%d/%s", "proc", i_pid, "smaps_rollup");
	if ((fp = fopen(path, "r")) != NULL) {

		// Getting 'mmap' header.
		if ( ! get_map_header(fp, path, SMAPS_ROLLUP, &store_smaps_rollup->header))
			return false;

		// get the second argument of the line which contains 3 arguments.
		if ( ! convert_to_bytes) {
			while (is_correctly_read(
			           fp,
			           path,
			           fscanf(fp, "%s%llu%s", fake_store, &values[count], fake_store), 3))
				count++;

			if (ERROR_IS_SET)
				return false;

		} else { /*  Values converted to bytes instead of Kb. */
			unsigned long long tmp_value;
			while (is_correctly_read(
			           fp,
			           path,
			           fscanf(fp, "%s%llu%s", fake_store, &tmp_value, fake_store), 3)) {
				values[count] = tmp_value * 1024;
				count++;
			}

			if (ERROR_IS_SET)
				return false;
		}

		// assignate value to structure (quick method)
		memcpy(&store_smaps_rollup->rss, &values, sizeof(store_smaps_rollup->rss) * count);
		fclose(fp);
		return true;
	}

	sprintf(ERROR_MESSAGE, "Error opening file [%s]", path);
	ERROR_IS_SET = true;
	return false;
}

// Create and initialize a new smaps_rollup structure.
smaps_rollup *smaps_rollup_new()
{
	smaps_rollup *new_smaps_rollup = malloc(sizeof(smaps_rollup));
	assert(new_smaps_rollup != NULL);
	return new_smaps_rollup;
}

// Free a smaps_rollup structure.
void smaps_rollup_free(smaps_rollup *smaps_rollup_to_free)
{
	free(smaps_rollup_to_free);
	smaps_rollup_to_free = NULL;
}

/***********************************************************
 * map_header, structure to hold header type values.
 * this format is used in: 'maps', 'smaps', 'smaps_rollup'.
 ***********************************************************/
// retrieve map_header header information as structure.
bool get_map_header(FILE *fp, char path[], T_CALLER_MMAP caller, map_header *curr_map_header)
{
	char *item;

	// ogiginal format: "%08lx-%08lx%c%c%c%c%08llx%02lx:%02lx%lu"
	if ( ! is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%lx-%lx %s %llx %lx:%lx %lu %s",
	                &curr_map_header->start,
	                &curr_map_header->end,
	                curr_map_header->flags,
	                &curr_map_header->offset,
	                &curr_map_header->dev_maj,
	                &curr_map_header->dev_min,
	                &curr_map_header->inode,
	                curr_map_header->pathname), 8))
		return false;

	// caller selection.
	switch (caller) {
	case MAPS: // not yet implemented
		break;
	case SMAPS_ROLLUP:
		return true;
	case SMAPS:
		item = "Size:";
	}

	char tmp_pathname[256];
	// Check the last argument retrieved, if its equal to 'item'
	// this means that the next line is reached and we have to go back
	// to previous value for the next 'fscanf'.
	if (strcmp(curr_map_header->pathname, item)) {
		if ( ! is_correctly_read(
		         fp,
		         path,
		         fscanf(fp, "%s", tmp_pathname), 1))
			return false;

		while (strcmp(tmp_pathname, item)) {
			strcat(strcat(curr_map_header->pathname, " "),tmp_pathname);
			if ( ! is_correctly_read(
			         fp,
			         path,
			         fscanf(fp, "%s", tmp_pathname), 1))
				return false;
		}
	} else {
		//clear pathname
		memset(curr_map_header->pathname, 0, sizeof(curr_map_header->pathname));
	}
	fseek(fp, - strlen(item), SEEK_CUR);

	return true;
}

// retrieve map_header header information as structure.
bool get_map_header_N(FILE *fp, char path[], map_header *curr_map_header)
{
	char *item = "Size:";

	// ogiginal format: "%08lx-%08lx%c%c%c%c%08llx%02lx:%02lx%lu"
	if ( ! is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%lx-%lx %s %llx %lx:%lx %lu %s",
	                &curr_map_header->start,
	                &curr_map_header->end,
	                curr_map_header->flags,
	                &curr_map_header->offset,
	                &curr_map_header->dev_maj,
	                &curr_map_header->dev_min,
	                &curr_map_header->inode,
	                curr_map_header->pathname), 8))
		return false;

	// Check the last argument retrieved, if its equal to 'item'
	// this means that the next line is reached and we have to go back
	// to previous value for the next 'fscanf'.
	if (strcmp(curr_map_header->pathname, item) == 0) {
		fseek(fp, - strlen(curr_map_header->pathname), SEEK_CUR);
		// clear pathname
		memset(curr_map_header->pathname, 0, sizeof(curr_map_header->pathname));
	}
	return true;
}

// Create and initialize a new map_header structure.
map_header *map_header_new()
{
	map_header *new_map_header = malloc(sizeof(map_header));
	assert(new_map_header != NULL);
	return new_map_header;
}

// Free a map_header structure.
void map_header_free(map_header *map_header_to_free)
{
	free(map_header_to_free);
	map_header_to_free = NULL;
}

// Return string representation of values (like the
// one available in: 'maps', 'smaps', 'smaps_rollup'.
void map_header_to_string(char *out_str, map_header *curr_map_header)
{
	sprintf(out_str, "%08lx-%08lx %s %08llx %02lx:%02lx %lu %s",
	        curr_map_header->start,
	        curr_map_header->end,
	        curr_map_header->flags,
	        curr_map_header->offset,
	        curr_map_header->dev_maj,
	        curr_map_header->dev_min,
	        curr_map_header->inode,
	        curr_map_header->pathname);
}
