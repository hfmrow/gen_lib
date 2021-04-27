// sys_smaps.h

#ifndef FUNCTIONS_SMAPS_INCLUDED
#define FUNCTIONS_SMAPS_INCLUDED
/* ^^ these are the include guards */

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include "sys_cpu_percent.h"

// structure to hold 'map_header' values.
// Information 'man procfs' search '/maps' then press
// 'n' until '/proc/[pid]/map_header'.
typedef struct map_header {
	unsigned long start;
	unsigned long end;
	char flags[4];
	unsigned long long offset;
	unsigned long dev_maj;
	unsigned long dev_min;
	unsigned long inode;
	char pathname[256];
} map_header;

// Create and initialize a new map_header structure.
map_header *map_header_new();

// Free a map_header structure.
void map_header_free(map_header *map_header_to_free);

// Used to know from where the 'get_map_header' function
// was called.
typedef enum {
	MAPS,
	SMAPS_ROLLUP,
	SMAPS,
} T_CALLER_MMAP;

// retrieve map_header header information as structure.
bool get_map_header(FILE *fp, char path[], T_CALLER_MMAP caller, map_header *curr_map_header);

// Return string representation of values (like the
// one available in: 'maps', 'smaps', 'smaps_rollup'.
void map_header_to_string(char *out_str, map_header *curr_map_header);

// structure to hold 'smaps_rollup' values.
// Information 'man procfs' search '/smaps_rollup'
// then press 'n' until '/proc/[pid]/smaps_rollup'.
typedef struct smaps_rollup {
	map_header header;
	unsigned long long rss;
	unsigned long long pss;
	unsigned long long pss_anon;
	unsigned long long pss_file;
	unsigned long long pss_shmem;
	unsigned long long shared_clean;
	unsigned long long shared_dirty;
	unsigned long long private_clean;
	unsigned long long private_dirty;
	unsigned long long referenced;
	unsigned long long anonymous;
	unsigned long long lazy_free;
	unsigned long long anon_huge_pages;
	unsigned long long shmem_pmd_mapped;
	unsigned long long file_pmd_mapped;
	unsigned long long shared_hugetlb;
	unsigned long long private_hugetlb;
	unsigned long long swap;
	unsigned long long swap_pss;
	unsigned long long locked;
} smaps_rollup;

// Create and initialize a new smaps_rollup structure.
smaps_rollup *smaps_rollup_new();

// Free a smaps_rollup structure.
void smaps_rollup_free(smaps_rollup *smaps_rollup_to_free);

// Read values starting from 'smaps_rollup' files
bool get_smaps_rollup(int i_pid, smaps_rollup *store_smaps_rollup, bool convert_to_bytes);

// structure to hold 'smaps' values.
// Information 'man procfs' search '/smaps' then press
// 'n' until '/proc/[pid]/smaps'.
typedef struct smap {
	map_header header;
	unsigned long long size;
	unsigned long long kernel_page_size;
	unsigned long long mmupage_size;
	unsigned long long rss;
	unsigned long long pss;
	unsigned long long shared_clean;
	unsigned long long shared_dirty;
	unsigned long long private_clean;
	unsigned long long private_dirty;
	unsigned long long referenced;
	unsigned long long anonymous;
	unsigned long long lazy_free;
	unsigned long long anon_huge_pages;
	unsigned long long shmem_pmd_mapped;
	unsigned long long file_pmd_mapped;
	unsigned long long shared_hugetlb;
	unsigned long long private_hugetlb;
	unsigned long long swap;
	unsigned long long swap_pss;
	unsigned long long locked;
	int protection_key;
	char vm_flags[128];
} smap;

// Create and initialize a new smap structure.
smap *smap_new();

// Free a smap structure.
void smap_free(smap *smap_to_free);

// Get single element values from 'smaps' file.
bool smap_get(FILE *fp, char path[], smap *store_smap, bool convert_to_bytes);

// Structure to hold 'smaps' items.
typedef struct smaps {
	int count;
	int internal_alloc;
	int internal_alloc_step;
	smap *details;
} smaps;

// update/get all values from 'smaps' file.
bool smaps_update(int i_pid, smaps *store_smaps, bool convert_to_bytes);

// Create and initialize a new smaps structure.
smaps *smaps_new(int count);

// Free a smaps structure.
void smaps_free(smaps *smaps_to_free);

// init and get smaps (used on 1st time request)
smaps *smaps_get(int pid, int size, bool convert_to_bytes);

// only used for binding structure (with golang for example)
// retrieve specific smap from array.
smap *smaps_get_smap(smaps *smaps_to_smap, int ask_for);

#endif
