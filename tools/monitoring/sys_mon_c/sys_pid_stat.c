// sys_pid_stat.c

#include "file_func.h"
#include "sys_pid_stat.h"

/****************
 * proc_pid_stat
 ****************/
// update '/proc/[pid]/stat' data
bool proc_pid_stat_update(proc_pid_stat *curr_proc_pid_stat)
{
	FILE *fp;
	char path[PATH_LENGTH];

	sprintf(path, "/%s/%d/%s", "proc", curr_proc_pid_stat->pid, "stat");
	fp = fopen(path, "r");
	if (fp == NULL) {
		memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
		sprintf(ERROR_MESSAGE, "Unable to open: %s", path);
		ERROR_IS_SET = true;
		return false;
	}
	if ( !is_correctly_read(
	         fp,
	         path,
	         fscanf(fp, "%d%s%s%d%d%d%d%d%u%lu%lu%lu%lu%lu%lu%ld%ld%ld%ld%ld%ld%llu%lu%ld%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%lu%d%d%u%u%llu%lu%ld%lu%lu%lu%lu%lu%lu%lu%d",
	                &curr_proc_pid_stat->pid,
	                curr_proc_pid_stat->comm,
	                curr_proc_pid_stat->state,
	                &curr_proc_pid_stat->ppid,
	                &curr_proc_pid_stat->pgrp,
	                &curr_proc_pid_stat->session,
	                &curr_proc_pid_stat->tty_nr,
	                &curr_proc_pid_stat->tpgid,
	                &curr_proc_pid_stat->flags,
	                &curr_proc_pid_stat->minflt,
	                &curr_proc_pid_stat->cminflt,
	                &curr_proc_pid_stat->majflt,
	                &curr_proc_pid_stat->cmajflt,
	                &curr_proc_pid_stat->utime,
	                &curr_proc_pid_stat->stime,
	                &curr_proc_pid_stat->cutime,
	                &curr_proc_pid_stat->cstime,
	                &curr_proc_pid_stat->priority,
	                &curr_proc_pid_stat->nice,
	                &curr_proc_pid_stat->num_threads,
	                &curr_proc_pid_stat->itrealvalue,
	                &curr_proc_pid_stat->starttime,
	                &curr_proc_pid_stat->vsize,
	                &curr_proc_pid_stat->rss,
	                &curr_proc_pid_stat->rsslim,
	                &curr_proc_pid_stat->startcode,
	                &curr_proc_pid_stat->endcode,
	                &curr_proc_pid_stat->startstack,
	                &curr_proc_pid_stat->kstkesp,
	                &curr_proc_pid_stat->kstkeip,
	                &curr_proc_pid_stat->signal,
	                &curr_proc_pid_stat->blocked,
	                &curr_proc_pid_stat->sigignore,
	                &curr_proc_pid_stat->sigcatch,
	                &curr_proc_pid_stat->wchan,
	                &curr_proc_pid_stat->nswap,
	                &curr_proc_pid_stat->cnswap,
	                &curr_proc_pid_stat->exit_signal,
	                &curr_proc_pid_stat->processor,
	                &curr_proc_pid_stat->rt_priority,
	                &curr_proc_pid_stat->policy,
	                &curr_proc_pid_stat->delayacct_blkio_ticks,
	                &curr_proc_pid_stat->guest_time,
	                &curr_proc_pid_stat->cguest_time,
	                &curr_proc_pid_stat->start_data,
	                &curr_proc_pid_stat->end_data,
	                &curr_proc_pid_stat->start_brk,
	                &curr_proc_pid_stat->arg_start,
	                &curr_proc_pid_stat->arg_end,
	                &curr_proc_pid_stat->env_start,
	                &curr_proc_pid_stat->env_end,
	                &curr_proc_pid_stat->exit_code), 52))
		return false;
	fclose(fp);

	return true;
}

// retrieve '/proc/[pid]/stat' data
proc_pid_stat *proc_pid_stat_get(unsigned int pid)
{
	proc_pid_stat *curr_proc_pid_stat = proc_pid_stat_new(pid);
	curr_proc_pid_stat->pid = pid;
	proc_pid_stat_update(curr_proc_pid_stat);
	return curr_proc_pid_stat;
}

// Initialize a new a proc_pid_stat structure.
proc_pid_stat *proc_pid_stat_new()
{
	proc_pid_stat *curr_proc_pid_stat = malloc(sizeof(proc_pid_stat));
	return curr_proc_pid_stat;
}

// Free a proc_pid_stat structure.
void proc_pid_stat_free(proc_pid_stat *proc_pid_stat_to_free)
{
	free(proc_pid_stat_to_free);
	proc_pid_stat_to_free = NULL;
}

