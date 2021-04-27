// sys_pid_status.c

#include <sys/types.h>
#include <sys/stat.h>
#include <errno.h>
#include <unistd.h>
#include <dirent.h>
#include <regex.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <err.h>
#include <stdlib.h>
#include <assert.h>
#include "sys_pid_status.h"
#include "file_func.h"


/*********************************************
 * This part retrive '/proc/PID/status' data
 *********************************************/
// read '/proc/PID/status' file and fill corresponding 'status_file' struct.
bool status_file_read(status_file *curr_status_file, uint pid)
{
	char path[PATH_LENGTH];
	char temp_read[256] = {0};
	bool ret = false;

	FILE *fp;
	sprintf(path, "/proc/%d/status", pid);
	fp = fopen(path, "r");
	if (fp == NULL) {
		memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
		sprintf(ERROR_MESSAGE, "Unable to open: %s", path);
		ERROR_IS_SET = true;
		goto end;
	}
	if ( !is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%s%s%s",
	                temp_read, curr_status_file->name,
	                temp_read, curr_status_file->umask), 4))
		goto end;
	char *tmp_values = get_paired_name_value(fp, ":", "State");
	if (tmp_values == NULL) {
		internal_error_set("Unable to read 'State' field");
		goto end;
	}
	strcpy(curr_status_file->state, tmp_values);
	if ( !is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%u%s%u%s%u%s%u%s%u",
	                temp_read, &curr_status_file->tgid,
	                temp_read, &curr_status_file->ngid,
	                temp_read, &curr_status_file->pid,
	                temp_read, &curr_status_file->ppid,
	                temp_read, &curr_status_file->tracer_pid), 10))
		goto end;
	tmp_values = get_paired_name_value(fp, ":", "Uid");
	if (tmp_values == NULL) {
		internal_error_set("Unable to read 'Uid' field");
		goto end;
	}
	sscanf(tmp_values, "%u%u%u%u",
	       &curr_status_file->uid.real,
	       &curr_status_file->uid.effective,
	       &curr_status_file->uid.saved_set,
	       &curr_status_file->uid.file_system);
	tmp_values = get_paired_name_value(fp, ":", "Gid");
	if (tmp_values == NULL) {
		internal_error_set("Unable to read 'Gid' field");
		goto end;
	}
	sscanf(tmp_values, "%u%u%u%u",
	       &curr_status_file->gid.real,
	       &curr_status_file->gid.effective,
	       &curr_status_file->gid.saved_set,
	       &curr_status_file->gid.file_system);
	if ( !is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%llu",
	                temp_read, &curr_status_file->fd_size), 2))
		goto end;
	curr_status_file->groups = malloc(sizeof(uint_array));
	if ( !uint_array_fill(fp, curr_status_file->groups, ":", " ", "Groups")) {
		free(curr_status_file->groups);
		curr_status_file->groups = NULL;
		goto end;
	}
	curr_status_file->ns_tgid = malloc(sizeof(uint_array));
	if ( !uint_array_fill(fp, curr_status_file->ns_tgid, ":", " ", "NStgid")) {
		free(curr_status_file->ns_tgid);
		curr_status_file->ns_tgid = NULL;
		goto end;
	}
	curr_status_file->ns_pid = malloc(sizeof(uint_array));
	if ( !uint_array_fill(fp, curr_status_file->ns_pid, ":", " ", "NSpid")) {
		free(curr_status_file->ns_pid);
		curr_status_file->ns_pid = NULL;
		goto end;
	}
	curr_status_file->ns_pgid = malloc(sizeof(uint_array));
	if ( !uint_array_fill(fp, curr_status_file->ns_pgid, ":", " ", "NSpgid")) {
		free(curr_status_file->ns_pgid);
		curr_status_file->ns_pgid = NULL;
		goto end;
	}
	curr_status_file->ns_sid = malloc(sizeof(uint_array));
	if ( !uint_array_fill(fp, curr_status_file->ns_sid, ":", " ", "NSsid")) {
		free(curr_status_file->ns_sid);
		curr_status_file->ns_sid = NULL;
		goto end;
	}
	// read data (if thy exist), to fill 'status_file_vmem' structure,
	// the 5 lines below, check for 'VmPeak' field exist before doing anything.
	tmp_values = get_paired_name_value(fp, ":", "VmPeak");
	if (tmp_values != NULL) {
		curr_status_file->vm = malloc(sizeof(status_file_vmem));
		if (is_correctly_read(
		        NULL,
		        path,
		        sscanf(tmp_values, "%llu%s", &curr_status_file->vm->vm_peak, temp_read), 2)) {
			if (is_correctly_read(
			        fp,
			        path,
			        fscanf(fp, "%s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%llu%s %s%u %s%u",
			               temp_read, &curr_status_file->vm->vm_size, temp_read,
			               temp_read, &curr_status_file->vm->vm_lck, temp_read,
			               temp_read, &curr_status_file->vm->vm_pin, temp_read,
			               temp_read, &curr_status_file->vm->vm_hwm, temp_read,
			               temp_read, &curr_status_file->vm->vm_rss, temp_read,
			               temp_read, &curr_status_file->vm->rss_anon, temp_read,
			               temp_read, &curr_status_file->vm->rss_file, temp_read,
			               temp_read, &curr_status_file->vm->rss_shmem, temp_read,
			               temp_read, &curr_status_file->vm->vm_data, temp_read,
			               temp_read, &curr_status_file->vm->vm_stk, temp_read,
			               temp_read, &curr_status_file->vm->vm_exe, temp_read,
			               temp_read, &curr_status_file->vm->vm_lib, temp_read,
			               temp_read, &curr_status_file->vm->vm_pte, temp_read,
			               temp_read, &curr_status_file->vm->vm_swap, temp_read,
			               temp_read, &curr_status_file->vm->hugetlb_pages, temp_read,
			               temp_read, &curr_status_file->vm->core_dumping,
			               temp_read, &curr_status_file->vm->thp_enabled), 49))
				goto next_fields;
		}
		/* DEBUG */
		printf("%s\n", ERROR_MESSAGE);

		// the 'VmPeak' field cannot be read correctly, so free the memory
		// allocated for the structure and set its pointer to NULL
		ERROR_IS_SET = false;
		free(curr_status_file->vm);
		curr_status_file->vm = NULL;
	} else {
		// the 'VmPeak' field seems not to exist, so set its pointer to NULL
		curr_status_file->vm = NULL;
	}
// continue with the 'Threads' field and the following ones.
next_fields:
	if ( !is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%s%d %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%s %s%d %s%d %s%s %s%s %s%s %s%s %s%s %s%llu %s%llu",
	                temp_read, &curr_status_file->threads,
	                temp_read, curr_status_file->sig_q,
	                temp_read, curr_status_file->sig_pnd,
	                temp_read, curr_status_file->shd_pnd,
	                temp_read, curr_status_file->sig_blk,
	                temp_read, curr_status_file->sig_ign,
	                temp_read, curr_status_file->sig_cgt,
	                temp_read, curr_status_file->cap_inh,
	                temp_read, curr_status_file->cap_prm,
	                temp_read, curr_status_file->cap_eff,
	                temp_read, curr_status_file->cap_bnd,
	                temp_read, curr_status_file->cap_amb,
	                temp_read, &curr_status_file->no_new_privs,
	                temp_read, &curr_status_file->seccomp,
	                temp_read, curr_status_file->speculation_Store_Bypass,
	                temp_read, curr_status_file->cpus_allowed,
	                temp_read, curr_status_file->cpus_allowed_list,
	                temp_read, curr_status_file->mems_allowed,
	                temp_read, curr_status_file->mems_allowed_list,
	                temp_read, &curr_status_file->voluntary_ctxt_switches,
	                temp_read, &curr_status_file->nonvoluntary_ctxt_switches), 42))
		goto end;
	ret = true;

end:
	if (fp != NULL)
		fclose(fp);
	return ret;
}

/**/
/* status_file_vmem */
/**/
void status_file_vmem_free(status_file_vmem *status_file_vmem_to_free)
{
	if (status_file_vmem_to_free != NULL) {
		free(status_file_vmem_to_free);
		status_file_vmem_to_free = NULL;
	}
}

/**/
/* uint_array */
/**/
// fill 'uint_array' structure with paired values starting at current file offset
// on error, FILE is closed, on success, FILE offset is pointing to next line.
bool uint_array_fill(FILE *fp, uint_array *uint_array_to_fill, char *sep_name, char *sep_vals, char *val_name)
{
	uint values[128] = {0};
	bool ret = false;
	int count = 0;
	/* get values */
	char *tmp_values = get_paired_name_value(fp, sep_name, val_name);
	if (tmp_values == NULL) {
		sprintf(ERROR_MESSAGE, "Unable to retrieve '%s' field", val_name);
		ERROR_IS_SET = true;
		goto end;
	}
	/* fill with values */
	char *ptr = strtok(tmp_values, sep_vals);
	while (ptr != NULL) {
		sscanf(ptr, "%u", &values[count++]);
		ptr = strtok(NULL, sep_vals);
	}
	if (count > 0) {
		uint_array_to_fill->count = count;
		uint_array_to_fill->values = malloc(sizeof(uint) * count);
		for (int i=0; i<count; i++) {
			uint_array_to_fill->values[i] = values[i];
		}
		ret = true;
		goto end;
	} else {
		sprintf(ERROR_MESSAGE, "There is no value paired with '%s' field", val_name);
		ERROR_IS_SET = true;
		goto end;
	}

end:
	if ( !ret)
		fclose(fp);
	return ret;
}

void uint_array_free(uint_array *uint_array_to_free)
{
	if (uint_array_to_free != NULL) {
		if (uint_array_to_free->count > 0) {
			free(uint_array_to_free->values);
			uint_array_to_free->values = NULL;
		}
		free(uint_array_to_free);
		uint_array_to_free = NULL;
	}
}

// retrieve uint[ask_for] single value from 'uint_array' structure
uint *uint_array_pick_value(uint_array *uint_array_to_pick, int ask_for)
{
	return &uint_array_to_pick->values[ask_for];
}

/**/
/* status_file */
/**/
void status_file_free(status_file *status_file_to_free)
{
	// freeing groups
	uint_array_free(status_file_to_free->groups);
	// freeing 'ns_tgid'
	uint_array_free(status_file_to_free->ns_tgid);
	// freeing 'ns_pid'
	uint_array_free(status_file_to_free->ns_pid);
	// freeing 'ns_pgid'
	uint_array_free(status_file_to_free->ns_pgid);
	// freeing 'ns_sid'
	uint_array_free(status_file_to_free->ns_sid);
	// freeing 'vm'
	status_file_vmem_free(status_file_to_free->vm);

	// freeing main structure
	free(status_file_to_free);
	status_file_to_free = NULL;
}

// create 'status_file' struct and initialize it using
// '/proc/PID/status' file that matches the given 'pid'.
status_file *status_file_new(uint pid)
{
	status_file *curr_status_file = malloc(sizeof(status_file));
	if (status_file_read(curr_status_file, pid))
		return curr_status_file;

	if ( !curr_status_file && ERROR_IS_SET == false) {
		memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
		sprintf(ERROR_MESSAGE, "Unable to retrieve information for pid: %d", pid);
		ERROR_IS_SET = true;
	}
	return NULL;
}

/**************************************************
 * This part retrive all '/proc/PID' information
 **************************************************/
int walk_dir(store_files *current_store_files, char *dname, regex_t *reg)
{
	FILE *fp;
	char path[256];
	char path_sub[512] = {0};

	struct dirent *dent;
	DIR *dir;
	struct stat st;
	char fn[128];
	int res = WALK_OK;
	int len = strlen(dname);
	if (len >= 128 - 1)
		return WALK_NAMETOOLONG;

	strcpy(fn, dname);
	fn[len++] = '/';
	if (!(dir = opendir(dname))) {
		warn("Warning: unable to open %s", dname);
		return WALK_BADIO;
	}

	current_store_files->count = 0;
	errno = 0;
	while ((dent = readdir(dir))) {
		if (!strcmp(dent->d_name, ".") || !strcmp(dent->d_name, ".."))
			continue;
		strncpy(fn + len, dent->d_name, 128 - len);
		if (stat(fn, &st) == -1) {
			errno = 0;
			warn("Warning: unable to stat %s", fn);
			continue;
		}
		if (S_ISDIR(st.st_mode)) {
			// pattern match then record entry
			if (!regexec(reg, fn, 0, 0, 0)) {
				// retrieve '/proc/PID/status' data
				sprintf(path, "%s/%s", fn, "status");
				fp = fopen(path, "r");
				if (fp != NULL) {
					strcpy(current_store_files->details[current_store_files->count].name, get_paired_name_value(fp, ":", "Name"));
					strcpy(current_store_files->details[current_store_files->count].state, get_paired_name_value(fp, ":", "State"));

					// retrieving some IDs information
					sscanf(get_paired_name_value(fp, ":", "Pid"), "%u", &current_store_files->details[current_store_files->count].pid);
					sscanf(get_paired_name_value(fp, ":", "PPid"), "%u", &current_store_files->details[current_store_files->count].ppid);

					sscanf(get_paired_name_value(fp, ":", "Uid"), "%u%u%u%u",
					       &current_store_files->details[current_store_files->count].uid.real,
					       &current_store_files->details[current_store_files->count].uid.effective,
					       &current_store_files->details[current_store_files->count].uid.saved_set,
					       &current_store_files->details[current_store_files->count].uid.file_system);

					sscanf(get_paired_name_value(fp, ":", "Gid"), "%u%u%u%u",
					       &current_store_files->details[current_store_files->count].gid.real,
					       &current_store_files->details[current_store_files->count].gid.effective,
					       &current_store_files->details[current_store_files->count].gid.saved_set,
					       &current_store_files->details[current_store_files->count].gid.file_system);
					fclose(fp);

					// allocate memory and copy directory name to structure
					strcpy(current_store_files->details[current_store_files->count].dirname, fn);

					// check for filename exist
					strcat(fn, "/exe");

					// clean temporary path container
					memset(path_sub, 0, sizeof(path_sub));

					if (readlink(fn, path_sub, sizeof(path_sub)) <= 0) {
						strcpy(current_store_files->details[current_store_files->count].filename, UNAVAILABLE_STR);
						errno = 0;
					} else
						strcpy(current_store_files->details[current_store_files->count].filename, path_sub);

					current_store_files->count++;

					// more memory need to be allocated ?
					if (current_store_files->size == current_store_files->count) {
						current_store_files->size += current_store_files->inc_size;
						// increase memory allocation
						current_store_files->details = realloc(current_store_files->details,
						                                       (sizeof(store_file) * current_store_files->size));
					}
				}
			}
		}
	}

	if (dir) closedir(dir);
	return res ? res : errno ? WALK_BADIO : WALK_OK;
}

// retrieve pid owned by the program 'filename'
unsigned int *get_pid_by_filename(char *filename, unsigned int *pid)
{
	store_files *curr_store_files = get_pid_infos();
	for (int i=0; i < curr_store_files->count; i++) {
		if ( !strcmp(basename_get(strdup(curr_store_files->details[i].filename)), filename)) {
			*pid = curr_store_files->details[i].pid;
			goto end;
		}
	}
	pid = NULL;

end:
	store_files_free(curr_store_files);
	return pid;
}

// retrieve pid owned by the program 'name'
unsigned int *get_pid_by_name(char *name, unsigned int *pid)
{
	store_files *curr_store_files = get_pid_infos();
	for (int i=0; i < curr_store_files->count; i++) {
		if ( !strcmp(curr_store_files->details[i].name, name)) {
			*pid = curr_store_files->details[i].pid;
			goto end;
		}
	}
	pid = NULL;

end:
	store_files_free(curr_store_files);
	return pid;
}

// return store_files structure that contain information
// retrieved in '/proc/PID/...', on ERROR, NULL is returned
// and the internal error can be checked to know from where
store_files *get_pid_infos()
{
	regex_t regx;
	if (regcomp(&regx, "[[:digit:]]+", REG_EXTENDED | REG_NOSUB)) {
		internal_error_set("Bad pattern");
		return NULL;
	}
	store_files *current_store_files = store_files_new(255);
	int res = walk_dir(current_store_files, "/proc", &regx);
	regfree(&regx);

	// resize to fit real count of entries
	current_store_files->details = realloc(current_store_files->details,
	                                       (sizeof(store_file) * current_store_files->count));
	if (current_store_files->details == NULL) {
		internal_error_set("Memory reallocation failed!.");
		return NULL;
	}
	switch(res) {
	case WALK_OK:
		return current_store_files;
	case WALK_BADIO:
		internal_error_set("IO error");
		break;
	case WALK_BADPATTERN:
		internal_error_set("Bad pattern");
		break;
	case WALK_NAMETOOLONG:
		internal_error_set("Filename too long");
		break;
	default:
		internal_error_set("Unknown error?");
	}
	return NULL;
}

// only used for binding structure (with golang for example)
store_file *store_files_get_single(store_files *curr_store_files, int ask_for)
{
	return &curr_store_files->details[ask_for];
}

store_files *store_files_new(int inc_size)
{
	store_files *new_files = malloc(sizeof(store_files));
	new_files->inc_size = inc_size;
	new_files->size = inc_size;
	new_files->details = malloc(sizeof(store_file) * new_files->size);

	return new_files;
}

void store_files_free(store_files *store_files_to_free)
{
	free(store_files_to_free->details);
	store_files_to_free->details = NULL;
	free(store_files_to_free);
	store_files_to_free = NULL;
}
