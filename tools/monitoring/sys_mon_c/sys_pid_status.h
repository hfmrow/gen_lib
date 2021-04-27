// sys_pid_status.h

#ifndef FUNCTIONS_WALK_DIR_INCLUDED
#define FUNCTIONS_WALK_DIR_INCLUDED
/* ^^ these are the include guards */

#include <regex.h>
#include <stdio.h>
#include <sys/stat.h>
#include <stdbool.h>

enum {
	WALK_OK = 0,
	WALK_BADPATTERN,
	WALK_NAMETOOLONG,
	WALK_BADIO,
};

typedef struct resf_id {
	unsigned int real;
	unsigned int effective;
	unsigned int saved_set;
	unsigned int file_system;
} resf_id ;

typedef struct status_file_vmem {
	unsigned long long vm_peak;
	unsigned long long vm_size;
	unsigned long long vm_lck;
	unsigned long long vm_pin;
	unsigned long long vm_hwm;
	unsigned long long vm_rss;
	unsigned long long rss_anon;
	unsigned long long rss_file;
	unsigned long long rss_shmem;
	unsigned long long vm_data;
	unsigned long long vm_stk;
	unsigned long long vm_exe;
	unsigned long long vm_lib;
	unsigned long long vm_pte;
	unsigned long long vm_swap;
	unsigned long long hugetlb_pages;
	int core_dumping; // 0/1 like bool
	int thp_enabled; // 0/1 like bool
} status_file_vmem ;

void status_file_vmem_free(status_file_vmem *status_file_vmem_to_free);

typedef struct uint_array {
	int count;
	unsigned int *values;
} uint_array ;

void uint_array_free(uint_array *uint_array_to_free);

// retrieve uint[ask_for] single value.
uint *uint_array_pick_value(uint_array *uint_array_to_pick, int ask_for);

// fill 'uint_array' structure with paired values starting at current file offset
// on error, FILE is closed, on success, FILE offset is pointing to next line.
bool uint_array_fill(FILE *fp, uint_array *uint_array_to_fill, char *sep_name, char *sep_vals, char *val_name);

/*********************************************
 * Structure to hold '/proc/PID/status' data
 *********************************************/
//typedef struct status_file status_file;
typedef struct status_file {
	char name[17];
	char umask[4];
	char state[64];
	unsigned int tgid;
	unsigned int ngid;
	unsigned int pid;
	unsigned int ppid;
	unsigned int tracer_pid;
	resf_id uid;
	resf_id gid;
	unsigned long long fd_size;
	// start list array
	uint_array *groups;
	uint_array *ns_tgid;
	uint_array *ns_pid;
	uint_array *ns_pgid;
	uint_array *ns_sid;
	// end list array
	status_file_vmem *vm; // when NULL means N/A
	int threads;
	char sig_q[16]; // number of signals queued/max.
	char sig_pnd[16];
	char shd_pnd[16];
	char sig_blk[16];
	char sig_ign[16];
	char sig_cgt[16];
	char cap_inh[16];
	char cap_prm[16];
	char cap_eff[16];
	char cap_bnd[16];
	char cap_amb[16];
	int no_new_privs;
	int seccomp;
	char speculation_Store_Bypass[16];
	char cpus_allowed[2];
	char cpus_allowed_list[16];
	char mems_allowed[32 * (8 + 1)]; // 00000000,00000000,00000000, ... (32 * (8 bytes + 1 byte for comma))
	char mems_allowed_list[16];
	unsigned long long voluntary_ctxt_switches;
	unsigned long long nonvoluntary_ctxt_switches;
} status_file ;

// read '/proc/PID/status' file and fill corresponding 'status_file' struct.
bool status_file_read(status_file *curr_status_file, uint pid);

void status_file_free(status_file *status_file_to_free);

// create 'status_file' struct and initialize it using
// '/proc/PID/status' file that matches the given 'pid'.
status_file *status_file_new(uint pid);

/*************************************
 * This part retrive all '/proc/PID'
 *************************************/
//typedef struct store_file store_file;
typedef struct store_file {
	unsigned int pid;
	unsigned int ppid;
	resf_id uid;
	resf_id gid;
	char name[17];
	char state[64];
	char filename[512];
	char dirname[128];
} store_file ;

//typedef struct store_files store_files;
typedef struct store_files {
	unsigned int count;		// number of files found
	unsigned int size;		// the real allocated size of the array
	unsigned int inc_size;	// memory allocation increase step
	store_file *details;	// files details
} store_files ;

store_files *store_files_new(int count);

void store_files_free(store_files *store_files_to_free);

// only used for binding structure (with golang for example)
store_file *store_files_get_single(store_files *curr_store_files, int ask_for);

int walk_dir(store_files *current_store_files, char *dname, regex_t *reg);

// return store_files structure that contain information
// retrieved in '/proc/PID/...', on ERROR, NULL is returned
// and the internal error can be checked to know from where
store_files *get_pid_infos();

// retrieve pid owned by the program 'name'
unsigned int *get_pid_by_name(char *name, unsigned int *pid);

// retrieve pid owned by the program 'filename'
unsigned int *get_pid_by_filename(char *filename, unsigned int *pid);

#endif
