// sys_partitions.h

#ifndef FUNCTIONS_PARTITIONS_INCLUDED
#define FUNCTIONS_PARTITIONS_INCLUDED

#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>

// structure to hold '/proc/diskstats' data
typedef struct diskstats {
	char device[64];
	char dev_type[32];
	int count;
	unsigned long reads_completed;
	unsigned long reads_merged;
	unsigned long sectors_read;
	unsigned long time_reading_ms;
	unsigned long writes_completed;
	unsigned long writes_merged;
	unsigned long sectors_written;
	unsigned long time_writing_ms;
	unsigned long ios_in_progress;
	unsigned long time_doing_ios_ms;
	unsigned long weighted_time_ios_ms;
	// the following values will be read if they are present
	// otherwise they will be ignored
	unsigned long discards_completed;
	unsigned long discards_merged;
	unsigned long sectors_discarded;
	unsigned long spent_discarding_ms;
	unsigned long undefined1;
	unsigned long undefined2;
	unsigned long undefined3;
} diskstats;

bool diskstats_update(diskstats *curr_partitions);

// Create and initialize a new 'diskstats' structure.
diskstats *diskstats_new();

// Free a 'diskstats' structure.
void diskstats_free(diskstats *partitions_to_free);

// only used for binding structure (with golang for example)
diskstats *diskstats_get_single(diskstats *curr_partitions, int ask_for);

diskstats *diskstats_get();

// structure to hold '/sys/class/block/[dev]/queue/' data
// uuid comes from: '/dev/disk/by-uuid/'
typedef struct class_block {
	unsigned long hw_sector_size;
	unsigned long logical_block_size;
	unsigned long max_hw_sectors_kb;
	unsigned long physical_block_size;
	unsigned long read_ahead_kb;
	unsigned long write_cache;
	unsigned long removable;
	unsigned long size;
	char dev_type[32];
	char part_name[256];
	char uuid[128];
} class_block;

// structure to hold '/proc/partitions' data
typedef struct partitions {
	int count;
	int major;
	int minor;
	unsigned long blocks;
	char name[32];
	class_block class_block;
} partitions;

bool partitions_update(partitions *curr_partitions);

// Create and initialize a new 'partitions' structure.
partitions *partitions_new();

// Free a 'partitions' structure.
void partitions_free(partitions *partitions_to_free);

// only used for binding structure (with golang for example)
partitions *partitions_get_single(partitions *curr_partitions, int ask_for);

partitions *partitions_get();

#endif
