// sys_therm.h

#ifndef FUNCTIONS_THERMAL_INCLUDED
#define FUNCTIONS_THERMAL_INCLUDED

#include <stdbool.h>

//#define SENSOR_MAX_COUNT 16
#define LABEL_SIZE 128
#define LABEL_MAX_COUNT 16
#define PATH_LENGTH_THERMAL 128

typedef struct sys_therm_sensor {
	int temp;
	int max;
	int crit;
	int crit_alarm;
	char temp_str[16];
	char max_str[16];
	char crit_str[16];
	char crit_alarm_str[16];
	char label[LABEL_SIZE];
} sys_therm_sensor;

typedef struct sys_therm {
	int count;
	int sensors_count;
	char name[LABEL_SIZE];
	char path[PATH_LENGTH_THERMAL];
	sys_therm_sensor *sensors;
} sys_therm;

// get all unique files containing 'temp'
bool get_temp_files(char *path, char **out_str, int *files_count_out);

// update/retrieve thermal information from '/sys/class/hwmon/hwmon*' files.
void sys_therm_update(sys_therm *curr_sys_therm);

// initialize sys_therm structure with sensors information
bool sys_therm_init(sys_therm *curr_sys_therm);

// only used for binding structure (with golang for example)
sys_therm *sys_therm_get_single(sys_therm *sys_therm_to_single, int ask_for);

// only used for binding structure (with golang for example)
sys_therm_sensor *sys_therm_sensor_get_single(sys_therm_sensor *sys_therm_sensor_to_single, int ask_for);

// retrieve thermal information from '/sys/class/hwmon/hwmon*' files.
// Note: "n/a", "-0°C" or "-1" value means data not available.
sys_therm *sys_therm_get();

// creat new 'sys_therm'structure
sys_therm *sys_therm_new();

// free 'sys_therm' structure
void sys_therm_free(sys_therm *sys_therm_to_free);

#endif
