// sys_partitions.c

#include <string.h>
#include <stdio.h>
#include "sys_partitions.h"
#include "file_func.h"

/***********
* Diskstats
************/

bool diskstats_update(diskstats *curr_diskstats)
{
	bool ret = false;
	char blank_string[128];

	int sizeof_line = 512*sizeof(char);
	char *line = (char *)malloc(sizeof_line);

	FILE *fp;
	char *path = "/proc/diskstats";
	//char *path = "/tmp/diskstats";

	fp = fopen(path, "r");
	if (fp == NULL)
		goto end;

	for (int i = 0; i < curr_diskstats->count; i++) {

		//if (line_contain(fp, curr_diskstats[i].device) == NULL){
		if (line_contain(fp, curr_diskstats[i].device, line, sizeof_line) == NULL) {
			// nothing was found, the device seems to be no longer
			// connected or something else.
			strcpy(curr_diskstats[i].device, UNAVAILABLE_STR);
			strcpy(curr_diskstats[i].dev_type, UNAVAILABLE_STR);
			continue;
		}

		// prefer reading line instead of file stream
		int expected_minimal_read = 14;
		//int read = fscanf(fp, "%s%s%s%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu",
		int read = sscanf(line, "%s%s%s%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu",
		                  blank_string,
		                  blank_string,
		                  blank_string,
		                  &curr_diskstats[i].reads_completed,
		                  &curr_diskstats[i].reads_merged,
		                  &curr_diskstats[i].sectors_read,
		                  &curr_diskstats[i].time_reading_ms,
		                  &curr_diskstats[i].writes_completed,
		                  &curr_diskstats[i].writes_merged,
		                  &curr_diskstats[i].sectors_written,
		                  &curr_diskstats[i].time_writing_ms,
		                  &curr_diskstats[i].ios_in_progress,
		                  &curr_diskstats[i].time_doing_ios_ms,
		                  &curr_diskstats[i].weighted_time_ios_ms, // 14
		                  // the following values will be read if they are
		                  // present otherwise they will be ignored.
		                  &curr_diskstats[i].discards_completed,
		                  &curr_diskstats[i].discards_merged,
		                  &curr_diskstats[i].sectors_discarded,
		                  &curr_diskstats[i].spent_discarding_ms,
		                  &curr_diskstats[i].undefined1,
		                  &curr_diskstats[i].undefined2,
		                  &curr_diskstats[i].undefined3);
		if (read <=  expected_minimal_read) {
			ERROR_IS_SET = true;
			sprintf(ERROR_MESSAGE, "Error reading file [%s], read: %d, minimal expected: %d variable(s)\n", path, read, expected_minimal_read);
			ret = false;
			goto end;
		}

		//printf("DEVICE: %s, READ: %d\n", curr_diskstats[i].device, read);
	}

	ret = true;

end:
	if (line != NULL)
		free(line);
	if (fp != NULL)
		fclose(fp);
	if (ret == false) {
		if ( !ERROR_IS_SET) {
			sprintf(ERROR_MESSAGE, "Unable to retrieve information from: %s", path);
			ERROR_IS_SET = true;
		}
	}

	return ret;
}

// Create and initialize a new 'diskstats' structure.
diskstats *diskstats_new()
{
	// get partitions count and names
	partitions *curr_partitions = partitions_get();
	if (curr_partitions == NULL)
		return NULL;

	diskstats *new_diskstats = malloc(sizeof(diskstats) * curr_partitions->count);
	new_diskstats->count =  curr_partitions->count;
	for (int i = 0; i < new_diskstats->count; i++) {
		strcpy(new_diskstats[i].device, curr_partitions[i].name);
		strcpy(new_diskstats[i].dev_type, curr_partitions[i].class_block.dev_type);
	}
	return new_diskstats;
}

// Free a 'diskstats' structure.
void diskstats_free(diskstats *diskstats_to_free)
{
	free(diskstats_to_free);
	diskstats_to_free = NULL;
}

diskstats *diskstats_get()
{
	diskstats *curr_diskstats = diskstats_new();
	if ( !diskstats_update(curr_diskstats)){
		diskstats_free(curr_diskstats);
		return NULL;
	}

	return curr_diskstats;
}

// only used for binding structure (with golang for example)
diskstats *diskstats_get_single(diskstats *curr_diskstats, int ask_for)
{
	return &curr_diskstats[ask_for];
}

/************
* Partitions
*************/
// base dire: '/sys/class/block/[dev]/*'
// following files are not present when device is a
// partition instead of a disk (except 'removable' and 'size')
#define CLASS_BLOCK_FILES_SIZE 8
static const char *CLASS_BLOCK_FILES_LIST[CLASS_BLOCK_FILES_SIZE] = {
	"queue/hw_sector_size",
	"queue/logical_block_size",
	"queue/max_hw_sectors_kb",
	"queue/physical_block_size",
	"queue/read_ahead_kb",
	"queue/write_cache",
	"removable",
	"size",
};

bool read_class_block_file(char *device, const char *suffix, void *dest)
{
	char path[PATH_LENGTH];
	sprintf(path, "/sys/class/block/%s/%s", device, suffix);
	fgets_void(path, "%lu", dest);
	return true;
}

bool class_block_fill(char *device, class_block *curr_class_block)
{
	bool ret = false;
	unsigned long values[CLASS_BLOCK_FILES_SIZE] = {0};

	char **files_list = calloc_2d_array(files_list, 128, 128);
	int files_count = 0;

	FILE *fp;
	char path[PATH_LENGTH];
	char *sep = "=";

	sprintf(path, "/sys/class/block/%s/%s", device, "uevent");
	fp = fopen(path, "r");
	if (fp == NULL) {
		goto end;
	}
	char *tmp_string = get_paired_name_value(fp, sep, "DEVTYPE");
	if (tmp_string != NULL) {
		sprintf(curr_class_block->dev_type, "%s", tmp_string);

		if ( !strcmp(tmp_string, "disk")) {
			for (int i = 0; i < CLASS_BLOCK_FILES_SIZE; i++) {
				read_class_block_file(device, CLASS_BLOCK_FILES_LIST[i], &values[i]);
			}
			memcpy(&curr_class_block->hw_sector_size, &values, sizeof(values));
		}
	} else
		sprintf(curr_class_block->dev_type, "%s", UNAVAILABLE_STR);

	tmp_string = get_paired_name_value(fp, sep, "PARTNAME");
	if (tmp_string != NULL)
		sprintf(curr_class_block->part_name, "%s", tmp_string);
	else
		sprintf(curr_class_block->part_name, "%s", UNAVAILABLE_STR);

	sprintf(path, "%s", "/dev/disk/by-uuid");
	if (list_files(path, files_list, &files_count))
		for (int i = 0; i < files_count; i++) {
			if (symlink_endpoint(files_list[i], path) != NULL) {
				strcpy(path, basename_get(path));
				if ( !strcmp(path, device)) {
					strcpy(curr_class_block->uuid, basename_get(files_list[i]));
					goto pre_end;
				}
			}
		}
	strcpy(curr_class_block->uuid, UNAVAILABLE_STR);

pre_end:
	free(files_list);
	ret = true;
end:
	if (fp != NULL)
		fclose(fp);
	if (ret == false) {
		sprintf(ERROR_MESSAGE, "Unable to retrieve information from: %s", path);
		ERROR_IS_SET = true;
	}
	return ret;
}

bool partitions_update(partitions *curr_partitions)
{
	bool ret = false;
	char *path = "/proc/partitions";
	FILE *fp = fopen(path, "r");
	if (fp == NULL) {
		goto end;
	}
	// skip 2 lines
	seek_to_next_line(fp);
	seek_to_next_line(fp);
	curr_partitions->count = 0;

	while (is_correctly_read(
	           fp,
	           path,
	           fscanf(fp, "%d%d%lu%s",
	                  &curr_partitions[curr_partitions->count].major,
	                  &curr_partitions[curr_partitions->count].minor,
	                  &curr_partitions[curr_partitions->count].blocks,
	                  curr_partitions[curr_partitions->count].name), 4)) {

		if ( !class_block_fill(
		         curr_partitions[curr_partitions->count].name,
		         &curr_partitions[curr_partitions->count].class_block))
			goto end;

		curr_partitions->count++;
	}

	if (ERROR_IS_EOF) {
		ret = true;
		long size = sizeof(partitions) * curr_partitions->count;
		curr_partitions = (partitions*)realloc(curr_partitions, size);
		if (curr_partitions == NULL)
			ret = false;
	}

end:
	if (fp != NULL)
		fclose(fp);
	if (ret == false) {
		partitions_free(curr_partitions);
		sprintf(ERROR_MESSAGE, "Unable to retrieve information from: %s", path);
		ERROR_IS_SET = true;
	}
	return ret;
}

// Create and initialize a new 'partitions' structure.
partitions *partitions_new()
{
	partitions *new_partitions = malloc(sizeof(partitions) * 64);
	return new_partitions;
}

// Free a 'partitions' structure.
void partitions_free(partitions *partitions_to_free)
{
	free(partitions_to_free);
	partitions_to_free = NULL;
}

partitions *partitions_get()
{
	partitions *curr_partitions = partitions_new();
	if ( !partitions_update(curr_partitions)){
		partitions_free(curr_partitions);
		return NULL;
	}
	return curr_partitions;
}

// only used for binding structure (with golang for example)
partitions *partitions_get_single(partitions *curr_partitions, int ask_for)
{
	return &curr_partitions[ask_for];
}
