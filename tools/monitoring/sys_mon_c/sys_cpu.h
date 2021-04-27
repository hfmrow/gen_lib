// sys_cpu.h

#ifndef FUNCTIONS_SYS_CPU_INCLUDED
#define FUNCTIONS_SYS_CPU_INCLUDED

#include <stdio.h>
#include <stdbool.h>
#include <time.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

/***************************
 * '/proc/stat' inforpation
 ***************************/
// Information 'man procfs' search '/stat' then press
// 'n' until '/proc/stat'.
typedef struct stat_cpu {
	long cpu_total_conf;
	long cpu_total_onln;
	long count;
	unsigned long long user;
	unsigned long long nice;
	unsigned long long system;
	unsigned long long idle;
	unsigned long long iowait;
	unsigned long long irq;
	unsigned long long softirq;
	unsigned long long steal;
	unsigned long long guest;
	unsigned long long guest_nice;
	char cpu[8];
} stat_cpu;

// Create and initialize a new stat_cpu structure.
stat_cpu *stat_cpu_new();

// Free a store_stat_cpu structure.
void stat_cpu_free(stat_cpu *stat_cpu_to_free);

// update/retrieve cpus information from '/proc/stat' file.
bool stat_cpu_update(stat_cpu *curr_stat_cpu);

// retrieve cpus information from '/proc/stat' file.
stat_cpu *stat_cpu_get();

/****************************************************************************************
 * Doc at: https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txtaffected_cpus
 * directory: /sys/devices/system/cpu/cpufreq/policy*
 ****************************************************************************************/
typedef struct cpu_fs {
	// 'cpu_count' used in the first entry to specify the
	// number of entries, (number of processors)
	int cpu_count;
	long long bios_limit;
	long long base_frequency;
	long long cpuinfo_cur_freq;
	long long cpuinfo_min_freq;
	long long cpuinfo_max_freq;
	long long cpuinfo_transition_latency;
	long long scaling_cur_freq;
	long long scaling_min_freq;
	long long scaling_max_freq;
	long long scaling_setspeed;
	// numeric list (as string)
	char *scaling_available_frequencies;
	char *related_cpus;
	// strings list
	char *scaling_available_governors;
	char *energy_performance_available_preferences;
	// string
	char *energy_performance_preference;
	char *scaling_driver;
	char *scaling_governor;
} cpu_fs;

// Create and initialize a new cpu_fs structure.
cpu_fs *cpu_fs_new();

// Free a cpu_fs structure.
void cpu_fs_free(cpu_fs *curr_cpu_fs);

// update/retrieve cpus information from '/sys/devices/system/cpu/cpufreq/policy*' files.
bool cpu_fs_update(cpu_fs *curr_cpu_fs);

// retrieve cpus information from '/sys/devices/system/cpu/cpufreq/policy*' files.
// this version retrieve only 'scaling_cur_freq' section to optimize operation.
bool cpu_fs_curr_freq_update(cpu_fs *curr_cpu_fs);

// retrieve cpus information from '/sys/devices/system/cpu/cpufreq/policy*' files.
cpu_fs *cpu_fs_get();

// only used for binding structure (with golang for example)
cpu_fs *cpu_fs_get_single(cpu_fs *cpu_fs_to_single, int ask_for);

/************************
 * Cpu clock calculation
 ************************/
typedef struct time_spent {
	struct timespec begin_nano;
	clock_t begin_ticks;
	double spent;

// Wall time (also known as clock time or wall-clock time) is simply
// the total time elapsed during the measurement. It’s the time you
// can measure with a stopwatch, assuming that you are able to start
// and stop it exactly at the execution points you want.
	int NANO_CLOCK_WALL;

// CPU Time, on the other hand, refers to the time the CPU was busy
// processing the program’s instructions. The time spent waiting for
// other things to complete (like I/O operations) is not included in
// the CPU time.
	int NANO_CLOCK_CPUTIME;

// selected methods 'NANO_CLOCK_WALL' or 'NANO_CLOCK_CPUTIME'
// must be set at 'time_spent' structure creation.
	int nano_measure_method;

// The number of clock ticks per second
	long SC_CLK_TCK;
} time_spent;

/* nanoseconds version */
// get current nano count measurement depend on 'curr_time_spent->nano_measure_method'
// argument 'NANO_CLOCK_WALL' or 'NANO_CLOCK_CPUTIME'
time_spent *time_nano_get(time_spent *curr_time_spent);

// calculate the nanoseconds between 2 measurement periods
time_spent *time_nano_calculate(time_spent *curr_time_spent);

/* ticks version */
// get current ticks count
time_spent *time_ticks_get(time_spent *curr_time_spent);

// calculate tick between 2 tick periods
time_spent *time_ticks_calculate(time_spent *curr_time_spent);

// Initialize a new a 'time_spent' structure.
time_spent *time_spent_new(int nano_measure_method);

// Free a 'time_spent' structure.
void time_spent_free(time_spent *curr_time_spent);


#endif
