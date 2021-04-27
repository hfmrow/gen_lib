// sys_meminfo.c

#include <stdio.h>
#include <stdlib.h>
#include "sys_meminfo.h"
#include "file_func.h"

// retrieving memory information from '/proc/meminfo' file.
bool meminfo_update(meminfo *curr_meminfo)
{
	bool ret = false;
	char *path = "/proc/meminfo";
	FILE *fp = fopen(path, "r");
	if (fp == NULL) {
		goto end;
	}
	char blank_string[64];
	if ( ! is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s%lu%s%lu%s%lu%s%lu%s%s%lu%s%s%lu%s%s%lu%s%s%lu%s",
	                blank_string, &curr_meminfo->mem_total, blank_string,
	                blank_string, &curr_meminfo->mem_free, blank_string,
	                blank_string, &curr_meminfo->mem_available, blank_string,
	                blank_string, &curr_meminfo->buffers, blank_string,
	                blank_string, &curr_meminfo->cached, blank_string,
	                blank_string, &curr_meminfo->swap_cached, blank_string,
	                blank_string, &curr_meminfo->active, blank_string,
	                blank_string, &curr_meminfo->inactive, blank_string,
	                blank_string, &curr_meminfo->active_anon_, blank_string,
	                blank_string, &curr_meminfo->inactive_anon_, blank_string,
	                blank_string, &curr_meminfo->active_file_, blank_string,
	                blank_string, &curr_meminfo->inactive_file_, blank_string,
	                blank_string, &curr_meminfo->unevictable, blank_string,
	                blank_string, &curr_meminfo->mlocked, blank_string,
	                blank_string, &curr_meminfo->swap_total, blank_string,
	                blank_string, &curr_meminfo->swap_free, blank_string,
	                blank_string, &curr_meminfo->dirty, blank_string,
	                blank_string, &curr_meminfo->writeback, blank_string,
	                blank_string, &curr_meminfo->anon_pages, blank_string,
	                blank_string, &curr_meminfo->mapped, blank_string,
	                blank_string, &curr_meminfo->shmem, blank_string,
	                blank_string, &curr_meminfo->kreclaimable, blank_string,
	                blank_string, &curr_meminfo->slab, blank_string,
	                blank_string, &curr_meminfo->sreclaimable, blank_string,
	                blank_string, &curr_meminfo->sunreclaim, blank_string,
	                blank_string, &curr_meminfo->kernel_stack, blank_string,
	                blank_string, &curr_meminfo->page_tables, blank_string,
	                blank_string, &curr_meminfo->nfs_unstable, blank_string,
	                blank_string, &curr_meminfo->bounce, blank_string,
	                blank_string, &curr_meminfo->writeback_tmp, blank_string,
	                blank_string, &curr_meminfo->commit_limit, blank_string,
	                blank_string, &curr_meminfo->committed_as, blank_string,
	                blank_string, &curr_meminfo->vmalloc_total, blank_string,
	                blank_string, &curr_meminfo->vmalloc_used, blank_string,
	                blank_string, &curr_meminfo->vmalloc_chunk, blank_string,
	                blank_string, &curr_meminfo->percpu, blank_string,
	                blank_string, &curr_meminfo->hardware_corrupted, blank_string,
	                blank_string, &curr_meminfo->anon_huge_pages, blank_string,
	                blank_string, &curr_meminfo->shmem_huge_pages, blank_string,
	                blank_string, &curr_meminfo->shmem_pmd_mapped, blank_string,
	                blank_string, &curr_meminfo->file_huge_pages, blank_string,
	                blank_string, &curr_meminfo->file_pmd_mapped, blank_string,
	                blank_string, &curr_meminfo->huge_pages_total,
	                blank_string, &curr_meminfo->huge_pages_free,
	                blank_string, &curr_meminfo->huge_pages_rsvd,
	                blank_string, &curr_meminfo->huge_pages_surp,
	                blank_string, &curr_meminfo->hugepagesize, blank_string,
	                blank_string, &curr_meminfo->hugetlb, blank_string,
	                blank_string, &curr_meminfo->direct_map4k, blank_string,
	                blank_string, &curr_meminfo->direct_map2_m, blank_string,
	                blank_string, &curr_meminfo->direct_map1_g, blank_string), 149))
		return false;

	ret=true;
end:
	if (fp != NULL)
		fclose(fp);
	if (ret == false) {
		meminfo_free(curr_meminfo);
		sprintf(ERROR_MESSAGE, "Unable to retrieve information from: %s", path);
		ERROR_IS_SET = true;
	}
	return ret;
}

meminfo *meminfo_get()
{
	meminfo *curr_meminfo = meminfo_new();
	meminfo_update(curr_meminfo);
	return curr_meminfo;
}

// Create and initialize a new meminfo structure.
meminfo *meminfo_new()
{
	meminfo *new_meminfo = malloc(sizeof(meminfo));
	return new_meminfo;
}

// Free a meminfo structure.
void meminfo_free(meminfo *meminfo_to_free)
{
	if (meminfo_to_free != NULL) {
		free(meminfo_to_free);
		meminfo_to_free = NULL;
	}
}
