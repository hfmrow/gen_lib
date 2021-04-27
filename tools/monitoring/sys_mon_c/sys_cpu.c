// sys_cpu.c
#include "sys_cpu.h"
#include "file_func.h"

/*********************************************
 * '/sys/devices/system/cpu/cpufreq/policy*'
 *********************************************/
#define CPUFREQ_FILES_SIZE 17
static const char *CPUFREQ_FILES_LIST[CPUFREQ_FILES_SIZE] = {
	"bios_limit",					// unsigned long long
	"base_frequency",
	"cpuinfo_cur_freq",
	"cpuinfo_min_freq",
	"cpuinfo_max_freq",
	"cpuinfo_transition_latency",
	"scaling_cur_freq",
	"scaling_min_freq",
	"scaling_max_freq",
	"scaling_setspeed",
	"related_cpus",					// numeric list (as string)
	"scaling_available_frequencies",// numeric list (as string)
	"scaling_available_governors",	// strings list
	"energy_performance_available_preferences",
	"energy_performance_preference",
	"scaling_driver",				// strings
	"scaling_governor",				// strings
};

// retrieve cpus information from '/sys/devices/system/cpu/cpufreq/policy*' files.
// this version retrieve only 'scaling_cur_freq' section to optimize operation.
bool cpu_fs_curr_freq_update(cpu_fs *curr_cpu_fs)
{
	FILE *fp;
	char path[PATH_LENGTH];
	int count = 0;
	int line_length = 32;
	char line[line_length];
	bool ret = false;

	for (; count < curr_cpu_fs->cpu_count; count++) {
		sprintf(path, "/sys/devices/system/cpu/cpufreq/policy%d/%s", count, "scaling_cur_freq");
		fp = fopen(path, "r");
		if (fp != NULL) {
			if (fgets (line, line_length, fp) == NULL)
				ret = false;
			fclose(fp);
			sscanf(line, "%lld", &curr_cpu_fs[count].scaling_cur_freq);
		}
	}
	if (ret != true) {
		sprintf(ERROR_MESSAGE, "Error: [%s]", "Unable to retrieve information from processors");
		ERROR_IS_SET = true;
		return false;
	}
	return true;
}

// update/retrieve cpus information from:
// '/sys/devices/system/cpu/cpufreq/policy*' files.
bool cpu_fs_update(cpu_fs *curr_cpu_fs)
{
	FILE *fp;
	char path[PATH_LENGTH];
	int count = 0;
	int line_length = 256;
	char line[line_length];
	bool ret = false;

	for (; count < curr_cpu_fs->cpu_count; count++) {
		int i = 0;
		while (CPUFREQ_FILES_LIST[i] != NULL) {
			sprintf(path, "/sys/devices/system/cpu/cpufreq/policy%d/%s", count, CPUFREQ_FILES_LIST[i]);
			fp = fopen(path, "r");
			if (fp != NULL) {
				ret = true;
				if (fgets (line, line_length, fp) == NULL)
					ret = false;
				fclose(fp);
				remove_lf(line);
			} else {
				if (i < 10)
					strcpy(line, UNAVAILABLE_STR_INT);
				else
					strcpy(line, UNAVAILABLE_STR);
			}

			switch (i) {
			case 0:
				sscanf(line, "%lld", &curr_cpu_fs[count].bios_limit);
				break;
			case 1:
				sscanf(line, "%lld", &curr_cpu_fs[count].base_frequency);
				break;
			case 2:
				sscanf(line, "%lld", &curr_cpu_fs[count].cpuinfo_cur_freq);
				break;
			case 3:
				sscanf(line, "%lld", &curr_cpu_fs[count].cpuinfo_min_freq);
				break;
			case 4:
				sscanf(line, "%lld", &curr_cpu_fs[count].cpuinfo_max_freq);
				break;
			case 5:
				sscanf(line, "%lld", &curr_cpu_fs[count].cpuinfo_transition_latency);
				break;
			case 6:
				sscanf(line, "%lld", &curr_cpu_fs[count].scaling_cur_freq);
				break;
			case 7:
				sscanf(line, "%lld", &curr_cpu_fs[count].scaling_min_freq);
				break;
			case 8:
				sscanf(line, "%lld", &curr_cpu_fs[count].scaling_max_freq);
				break;
			case 9:
				sscanf(line, "%lld", &curr_cpu_fs[count].scaling_setspeed);
				break;
			case 10:
				curr_cpu_fs[count].related_cpus = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].related_cpus, line);
				break;
			case 11:
				curr_cpu_fs[count].scaling_available_frequencies = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].scaling_available_frequencies, line);
				break;
			case 12:
				curr_cpu_fs[count].scaling_available_governors = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].scaling_available_governors, line);
				break;
			case 13:
				curr_cpu_fs[count].energy_performance_available_preferences = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].energy_performance_available_preferences, line);
				break;
			case 14:
				curr_cpu_fs[count].energy_performance_preference = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].energy_performance_preference, line);
				break;
			case 15:
				curr_cpu_fs[count].scaling_driver = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].scaling_driver, line);
				break;
			case 16:
				curr_cpu_fs[count].scaling_governor = (char*)malloc(sizeof(char) * strlen(line));
				strcpy(curr_cpu_fs[count].scaling_governor, line);
			}
			i++;
		}
	}
	if (ret != true) {
		sprintf(ERROR_MESSAGE, "Error: [%s]", "Unable to retrieve information from processors");
		ERROR_IS_SET = true;
		return false;
	}
	return true;
}

// retrieve cpus information from:
// '/sys/devices/system/cpu/cpufreq/policy*' files.
cpu_fs *cpu_fs_get()
{
	cpu_fs *curr_cpu_fs = cpu_fs_new();
	if (cpu_fs_update(curr_cpu_fs))
		return curr_cpu_fs;

	internal_error_set("Unable to retrieve information");
	return NULL;
}

// only used for binding structure (with golang for example)
cpu_fs *cpu_fs_get_single(cpu_fs *cpu_fs_to_single, int ask_for)
{
	return &cpu_fs_to_single[ask_for];
}

// Create and initialize a new cpu_fs structure.
cpu_fs *cpu_fs_new()
{
	cpu_fs *curr_cpu_fs = malloc(sizeof(cpu_fs) * sysconf(_SC_NPROCESSORS_CONF));
	curr_cpu_fs->cpu_count = sysconf(_SC_NPROCESSORS_CONF);
	return curr_cpu_fs;
}

// Free a cpu_fs structure.
void cpu_fs_free(cpu_fs *curr_cpu_fs)
{
	for (int i = 0; i < curr_cpu_fs->cpu_count; i++) {
		free(curr_cpu_fs[i].scaling_available_frequencies);
		curr_cpu_fs[i].scaling_available_frequencies = NULL;
		free(curr_cpu_fs[i].related_cpus);
		curr_cpu_fs[i].related_cpus = NULL;
		free(curr_cpu_fs[i].scaling_available_governors);
		curr_cpu_fs[i].scaling_available_governors = NULL;
		free(curr_cpu_fs[i].scaling_driver);
		curr_cpu_fs[i].scaling_driver = NULL;
		free(curr_cpu_fs[i].scaling_governor);
		curr_cpu_fs[i].scaling_governor = NULL;
	}
	free(curr_cpu_fs);
}

/****************
 * '/proc/stat'
 ****************/
// retrieving cpus information from '/proc/stat' file.
bool stat_cpu_update(stat_cpu *curr_stat_cpu)
{
	FILE *fp;
	fp = fopen("/proc/stat", "r");

	for (int i = 0; i < curr_stat_cpu->count; i++)
		if ( ! is_correctly_read(
		         fp,
		         "/proc/stat",
		         fscanf(fp, "%s%llu%llu%llu%llu%llu%llu%llu%llu%llu%llu",
		                curr_stat_cpu[i].cpu,
		                &curr_stat_cpu[i].user,
		                &curr_stat_cpu[i].nice,
		                &curr_stat_cpu[i].system,
		                &curr_stat_cpu[i].idle,
		                &curr_stat_cpu[i].iowait,
		                &curr_stat_cpu[i].irq,
		                &curr_stat_cpu[i].softirq,
		                &curr_stat_cpu[i].steal,
		                &curr_stat_cpu[i].guest,
		                &curr_stat_cpu[i].guest_nice), 11))
			return false;

	fclose(fp);

	return true;
}

stat_cpu *stat_cpu_get()
{
	stat_cpu *curr_stat_cpu = stat_cpu_new();
	stat_cpu_update(curr_stat_cpu);
	return curr_stat_cpu;
}

// Create and initialize a new stat_cpu structure.
stat_cpu *stat_cpu_new()
{
	stat_cpu *stat_cpu_new = malloc(sizeof(stat_cpu) * (sysconf(_SC_NPROCESSORS_CONF) + 1));
	stat_cpu_new->cpu_total_conf = sysconf(_SC_NPROCESSORS_CONF);
	stat_cpu_new->cpu_total_onln = sysconf(_SC_NPROCESSORS_ONLN);

	// we add 1 to the total number of processors since
	// the first is the sum of the others
	stat_cpu_new->count = stat_cpu_new->cpu_total_conf + 1;

	return stat_cpu_new;
}

// Free a stat_cpu structure.
void stat_cpu_free(stat_cpu *stat_cpu_to_free)
{
	free(stat_cpu_to_free);
	stat_cpu_to_free = NULL;
}

/************************
 * Cpu clock calculation
 ************************/
/* nanoseconds version */
// get current nano count measurement depend on 'curr_time_spent->nano_measure_method'
// argument 'NANO_CLOCK_WALL' or 'NANO_CLOCK_CPUTIME'
time_spent *time_nano_get(time_spent *curr_time_spent)
{
	clock_gettime(curr_time_spent->nano_measure_method, &curr_time_spent->begin_nano);
	return curr_time_spent;
}

// calculate the nanoseconds between 2 measurement periods
time_spent *time_nano_calculate(time_spent *curr_time_spent)
{
	struct timespec end;
	clock_gettime(curr_time_spent->nano_measure_method, &end);
	long seconds = end.tv_sec - curr_time_spent->begin_nano.tv_sec;
	long nanoseconds = end.tv_nsec - curr_time_spent->begin_nano.tv_nsec;
	curr_time_spent->spent = seconds + nanoseconds*1e-9;
	curr_time_spent->begin_nano = end;
	return curr_time_spent;
}

/* ticks version */
// get current ticks count
time_spent *time_ticks_get(time_spent *curr_time_spent)
{
	curr_time_spent->begin_ticks = clock();
	return curr_time_spent;
}

// calculate tick between 2 tick periods
time_spent *time_ticks_calculate(time_spent *curr_time_spent)
{
	clock_t end = clock();
	curr_time_spent->spent= (double)(end - curr_time_spent->begin_ticks); // / CLOCKS_PER_SEC;
	curr_time_spent->begin_ticks = end;
	return curr_time_spent;
}

//get current ticks count
//time_spent *time_ticks_get(time_spent *curr_time_spent)
//{
//	curr_time_spent->begin_ticks = times(&curr_time_spent->buf);
//	return curr_time_spent;
//}

//calculate tick between 2 tick periods
//time_spent *time_ticks_calculate(time_spent *curr_time_spent)
//{
//
//	clock_t end = times(&curr_time_spent->buf);
//	curr_time_spent->spent= (double)(end - curr_time_spent->begin_ticks);  / CLOCKS_PER_SEC;
//	curr_time_spent->begin_ticks = end;
//	return curr_time_spent;
//}


// Initialize a new a 'time_spent' structure.
time_spent *time_spent_new(int nano_measure_method)
{
	time_spent *curr_time_spent = malloc(sizeof(time_spent));
	curr_time_spent->NANO_CLOCK_WALL = CLOCK_REALTIME;
	curr_time_spent->NANO_CLOCK_CPUTIME = CLOCK_PROCESS_CPUTIME_ID;
	curr_time_spent->SC_CLK_TCK = sysconf(_SC_CLK_TCK);

	if (nano_measure_method == -1)
		curr_time_spent->nano_measure_method = CLOCK_REALTIME;
	else
		curr_time_spent->nano_measure_method = nano_measure_method;

	return curr_time_spent;
}

// Free a 'time_spent' structure.
void time_spent_free(time_spent *curr_time_spent)
{
	free(curr_time_spent);
	curr_time_spent = NULL;
}
