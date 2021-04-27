// sys_stat_x.h

#ifndef FUNCTIONS_STATE_INCLUDED
#define FUNCTIONS_STATE_INCLUDED
/* ^^ these are the include guards */

#include <stdbool.h>
#include "sys_cpu.h"
#include "sys_pid_stat.h"

// Structure to hold previously defined ones.
typedef struct cpu_percent_pid {
	float cpu_percent;
	long long page_size;
	long long memory_rss;
	stat_cpu *cpu;
	proc_pid_stat *pid;
} cpu_percent_pid;

// initialize a new 'cpu_percent_pid' structure.
cpu_percent_pid *cpu_percent_pid_new();

// Free a cpu_percent_pid structure.
void cpu_percent_pid_free(cpu_percent_pid *cpu_percent_pid_to_free);

// retrieving/updating values and cpu% calculation.
bool cpu_percent_pid_update(cpu_percent_pid *last_store);

// initialize/fill and get a new 'cpu_percent_pid'
// structure with error control.
cpu_percent_pid *cpu_percent_pid_get(unsigned int pid);

#endif
