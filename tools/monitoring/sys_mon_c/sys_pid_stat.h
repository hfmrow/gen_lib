// sys_pid_stat.h

#ifndef FUNCTIONS_PID_STAT_INCLUDED
#define FUNCTIONS_PID_STAT_INCLUDED

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdbool.h>

/****************
 * proc_pid_stat
 ****************/
// Information 'man procfs' search '/stat' then press
// 'n' until '/proc/[pid]/stat'.
typedef struct proc_pid_stat {
	unsigned int pid;
	char comm[17];
	char state[1];
	int ppid;
	int pgrp;
	int session;
	int tty_nr;
	int tpgid;
	unsigned int flags;
	unsigned long minflt;
	unsigned long cminflt;
	unsigned long majflt;
	unsigned long cmajflt;
	unsigned long utime;
	unsigned long stime;
	long cutime;
	long cstime;
	long priority;
	long nice;
	long num_threads;
	long itrealvalue;
	unsigned long long starttime;
	unsigned long vsize;
	long rss;
	unsigned long rsslim;
	unsigned long startcode;
	unsigned long endcode;
	unsigned long startstack;
	unsigned long kstkesp;
	unsigned long kstkeip;
	unsigned long signal;
	unsigned long blocked;
	unsigned long sigignore;
	unsigned long sigcatch;
	unsigned long wchan;
	unsigned long nswap;
	unsigned long cnswap;
	int exit_signal;
	int processor;
	unsigned int rt_priority;
	unsigned int policy;
	unsigned long long delayacct_blkio_ticks;
	unsigned long guest_time;
	long cguest_time;
	unsigned long start_data;
	unsigned long end_data;
	unsigned long start_brk;
	unsigned long arg_start;
	unsigned long arg_end;
	unsigned long env_start;
	unsigned long env_end;
	int exit_code;
} proc_pid_stat;

// update '/proc/[pid]/stat' data
bool proc_pid_stat_update(proc_pid_stat *curr_proc_pid_stat);

// retrieve '/proc/[pid]/stat' data
proc_pid_stat *proc_pid_stat_get(unsigned int pid);

// Initialize a new a proc_pid_stat structure.
proc_pid_stat *proc_pid_stat_new();

// Free a proc_pid_stat structure.
void proc_pid_stat_free(proc_pid_stat *proc_pid_stat_to_free);

#endif
