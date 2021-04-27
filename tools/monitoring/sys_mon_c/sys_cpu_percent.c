// sys_stat_x.c

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include "file_func.h"
#include "sys_cpu.h"
#include "sys_pid_stat.h"
#include "sys_cpu_percent.h"

// retrieving/updating values and cpu % calculation.
bool cpu_percent_pid_update(cpu_percent_pid *last_store)
{
	cpu_percent_pid *curr_store = cpu_percent_pid_new(last_store->pid->pid);

	if (stat_cpu_update(curr_store->cpu))
		if (proc_pid_stat_update(curr_store->pid)) {
			//Compute cpu % for this pid
			long long int user, nice, system, idle, utime, stime;
			utime = curr_store->pid->utime - last_store->pid->utime;
			stime = curr_store->pid->stime - last_store->pid->stime;
			user = curr_store->cpu->user - last_store->cpu->user;
			nice = curr_store->cpu->nice - last_store->cpu->nice;
			system = curr_store->cpu->system - last_store->cpu->system;
			idle = curr_store->cpu->idle - last_store->cpu->idle;

			// Store cpu percent value
			last_store->cpu_percent = (utime + stime) * 100.0 / (user + nice + system + idle);
			// store real memory rss value (page_size * rss)
			last_store->memory_rss = curr_store->page_size * curr_store->pid->rss;

			// Swap current to last values
			memmove(last_store->pid, curr_store->pid, sizeof(proc_pid_stat));
			memmove(last_store->cpu, curr_store->cpu, sizeof(stat_cpu));
			cpu_percent_pid_free(curr_store);

			return true;
		}

	internal_error_set("Unable to retrieve data for cpu% calculation");
	return false;
}

// initialize/fill and get a new 'cpu_percent_pid'
// structure with error control.
cpu_percent_pid *cpu_percent_pid_get(unsigned int pid)
{
	cpu_percent_pid *curr_cpu_percent_pid = cpu_percent_pid_new(pid);
	if (stat_cpu_update(curr_cpu_percent_pid->cpu)) {
		if (proc_pid_stat_update(curr_cpu_percent_pid->pid))
			return curr_cpu_percent_pid;
		else
			sprintf(ERROR_MESSAGE, "Error (pid): [%s]", "Unable to retrieve information from processors");
	} else
		sprintf(ERROR_MESSAGE, "Error (cpu): [%s]", "Unable to retrieve information from processors");

	ERROR_IS_SET = true;
	return NULL;
}

// initialize a new 'cpu_percent_pid' structure.
cpu_percent_pid *cpu_percent_pid_new(unsigned int pid)
{
	cpu_percent_pid *curr_cpu_percent_pid = malloc(sizeof(cpu_percent_pid));
	curr_cpu_percent_pid->page_size = sysconf(_SC_PAGESIZE);
	curr_cpu_percent_pid->cpu = stat_cpu_get();
	curr_cpu_percent_pid->cpu->count = 1; // get only 1st row
	curr_cpu_percent_pid->pid = proc_pid_stat_get(pid);
	return curr_cpu_percent_pid;
}

// Free a store_stat structure.
void cpu_percent_pid_free(cpu_percent_pid *cpu_percent_pid_to_free)
{
	stat_cpu_free(cpu_percent_pid_to_free->cpu);
	proc_pid_stat_free(cpu_percent_pid_to_free->pid);

	free(cpu_percent_pid_to_free);
	cpu_percent_pid_to_free = NULL;
}
