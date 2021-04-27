// sys_meminfo.h

#ifndef FUNCTIONS_MEMINFO_INCLUDED
#define FUNCTIONS_MEMINFO_INCLUDED

#include <stdbool.h>

// structure to hold '/proc/meminfo' data
typedef struct meminfo {
	unsigned long mem_total;
	unsigned long mem_free;
	unsigned long mem_available;
	unsigned long buffers;
	unsigned long cached;
	unsigned long swap_cached;
	unsigned long active;
	unsigned long inactive;
	unsigned long active_anon_;
	unsigned long inactive_anon_;
	unsigned long active_file_;
	unsigned long inactive_file_;
	unsigned long unevictable;
	unsigned long mlocked;
	unsigned long swap_total;
	unsigned long swap_free;
	unsigned long dirty;
	unsigned long writeback;
	unsigned long anon_pages;
	unsigned long mapped;
	unsigned long shmem;
	unsigned long kreclaimable;
	unsigned long slab;
	unsigned long sreclaimable;
	unsigned long sunreclaim;
	unsigned long kernel_stack;
	unsigned long page_tables;
	unsigned long nfs_unstable;
	unsigned long bounce;
	unsigned long writeback_tmp;
	unsigned long commit_limit;
	unsigned long committed_as;
	unsigned long vmalloc_total;
	unsigned long vmalloc_used;
	unsigned long vmalloc_chunk;
	unsigned long percpu;
	unsigned long hardware_corrupted;
	unsigned long anon_huge_pages;
	unsigned long shmem_huge_pages;
	unsigned long shmem_pmd_mapped;
	unsigned long file_huge_pages;
	unsigned long file_pmd_mapped;
	unsigned long huge_pages_total;
	unsigned long huge_pages_free;
	unsigned long huge_pages_rsvd;
	unsigned long huge_pages_surp;
	unsigned long hugepagesize;
	unsigned long hugetlb;
	unsigned long direct_map4k;
	unsigned long direct_map2_m;
	unsigned long direct_map1_g;
} meminfo;

// Create and initialize a new meminfo structure.
meminfo *meminfo_new();

// Free a meminfo structure.
void meminfo_free(meminfo *meminfo_to_free);

bool meminfo_update(meminfo *curr_meminfo);

meminfo *meminfo_get();

#endif
