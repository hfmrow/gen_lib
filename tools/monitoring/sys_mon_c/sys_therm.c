// sys_therm.c

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "file_func.h"
#include "sys_therm.h"



// initialize sys_therm structure with sensors ID
bool sys_therm_init(sys_therm *curr_sys_therm)
{
	FILE *fp;
	char path[PATH_LENGTH];
	char path_sub[PATH_LENGTH];
	int count = 0;
	int count_sub = 0;
	bool ret = false;

	for (; count < curr_sys_therm->count; count++) {
		sprintf(path, "/sys/class/hwmon/hwmon%d", count);
		//sprintf(path, "/media/syndicate/storage/Documents/dev/c/WorkSpaces/src-to-go/sys_stat/hwmon/hwmon%d", count);
		fp = fopen(path, "r");
		if (fp != NULL) {
			fclose(fp);

			sprintf(path_sub, "/sys/class/hwmon/hwmon%d/name", count);
			//sprintf(path_sub, "/media/syndicate/storage/Documents/dev/c/WorkSpaces/src-to-go/sys_stat/hwmon/hwmon%d/name", count);
			fp = fopen(path_sub, "r");
			if (fp  != NULL) {
				ret = true;
				// store base path
				strcpy(curr_sys_therm[count].path, path);
				// store ID
				if (fscanf(fp, "%s", curr_sys_therm[count].name) != 1)
					strcpy(curr_sys_therm[count].name, UNAVAILABLE_STR);

				fclose(fp);
				count_sub++;
			}
		} else {
			break;
		}
	}
	if (ret != true) {
		sprintf(ERROR_MESSAGE, "Error: [%s]", "Unable to retrieve thermal information");
		ERROR_IS_SET = true;
		return false;
	}
	// reallocate memory to fit new size
	curr_sys_therm = realloc(curr_sys_therm, sizeof(sys_therm) * count_sub);
	curr_sys_therm->count = count_sub;
	return true;
}

// get all unique files containing 'temp'
bool get_temp_files(char *path, char **out_str, int *files_count_out)
{
	bool ret = false;
	char **files_list = calloc_2d_array(files_list, 128, 128);
	int files_count = 0;
	int int_files_count_out = 0;

	if ( !list_files(path, files_list, &files_count))
		goto end;

	// sort files list
	sort(files_list, files_count);
	// copy first matching name to destination
	for (int i = 1; i < files_count; i++)
		if (reg_match(files_list[i], "temp")) {
			strcpy(out_str[0], cut_at_first(basename_get(files_list[i]), "_"));
			int_files_count_out++;
			break;
		}
	// add matching name to destination
	for (int i = 1; i < files_count; i++) {
		if (reg_match(files_list[i], "temp")) {
			char *tmp_str = cut_at_first(basename_get(files_list[i]), "_");
			if (strcmp(out_str[int_files_count_out - 1], tmp_str)) {
				strcpy(out_str[int_files_count_out], tmp_str);
				int_files_count_out++;
			}
		}
	}
	*files_count_out = int_files_count_out;
	ret = true;
end:
	free(files_list);
	return ret;
}

// update/retrieve thermal information from '/sys/class/hwmon/hwmon*' files.
void sys_therm_update(sys_therm *curr_sys_therm)
{
	char path[PATH_LENGTH];
	int count = 0;
	int count_sorted_temp_files;
	bool ret_sub;
	char labels[LABEL_MAX_COUNT][LABEL_SIZE] = {0};
	float values[LABEL_MAX_COUNT][3] = {0};
	char **sorted_temp_files = calloc_2d_array(sorted_temp_files, 128, 128);

	for (; count < curr_sys_therm->count; count++) {

		// some systems do not start at 'temp0' so we need to
		// have a list of files to avoid this complication.
		get_temp_files(curr_sys_therm[count].path, sorted_temp_files, &count_sorted_temp_files);

		for (int i = 0; i < count_sorted_temp_files; i++) {
			ret_sub = false;

			sprintf(path, "%s/%s_crit", curr_sys_therm[count].path, sorted_temp_files[i]);
			if (fgets_void(path, "%f", &values[i][0]))
				ret_sub = true;

			sprintf(path, "%s/%s_crit_alarm", curr_sys_therm[count].path, sorted_temp_files[i]);
			if (fgets_void(path, "%f", &values[i][1]))
				ret_sub = true;

			sprintf(path, "%s/%s_input", curr_sys_therm[count].path, sorted_temp_files[i]);
			if (fgets_void(path, "%f", &values[i][2]))
				ret_sub = true;

			sprintf(path, "%s/%s_max", curr_sys_therm[count].path, sorted_temp_files[i]);
			if (fgets_void(path, "%f", &values[i][3]))
				ret_sub = true;

			sprintf(path, "%s/%s_label", curr_sys_therm[count].path, sorted_temp_files[i]);
			if (fgets_void (path, "%s", &labels[i][0]))
				ret_sub = true;

			if (!ret_sub)
				break;
		}

		// allocate memory for count of sensors found
		curr_sys_therm[count].sensors = (sys_therm_sensor*)malloc(sizeof(sys_therm_sensor) * count_sorted_temp_files);
		curr_sys_therm[count].sensors_count = count_sorted_temp_files;

		for (int j = 0; j < curr_sys_therm[count].sensors_count; j++) {
			curr_sys_therm[count].sensors[j].crit = (int)values[j][0];
			curr_sys_therm[count].sensors[j].crit_alarm = (int)values[j][1];
			curr_sys_therm[count].sensors[j].temp = (int)values[j][2];
			curr_sys_therm[count].sensors[j].max = (int)values[j][3];

			sprintf(curr_sys_therm[count].sensors[j].crit_str, "%.0f°C", values[j][0]/1000);
			sprintf(curr_sys_therm[count].sensors[j].crit_alarm_str, "%.0f°C", values[j][1]/1000);
			sprintf(curr_sys_therm[count].sensors[j].temp_str, "%.0f°C", values[j][2]/1000);
			sprintf(curr_sys_therm[count].sensors[j].max_str, "%.0f°C", values[j][3]/1000);
			sprintf(curr_sys_therm[count].sensors[j].label, "%s", &labels[j][0]);
		}
	}

	free(sorted_temp_files);
}

// retrieve thermal information from '/sys/class/hwmon/hwmon*' files.
// Note: "n/a", "-0°C" or "-1" value means not available data.
sys_therm *sys_therm_get()
{
	sys_therm *curr_sys_therm = sys_therm_new(32);

	// 1st pass to get sensors ID
	if (sys_therm_init(curr_sys_therm)) {

		// realloc memory to fit real size
		curr_sys_therm = realloc(curr_sys_therm, sizeof(sys_therm) * curr_sys_therm->count);

		// 2nd pass to get values
		sys_therm_update(curr_sys_therm);
		return curr_sys_therm;
	}

	if ( !curr_sys_therm) {
		memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
		sprintf(ERROR_MESSAGE, "Unable to retrieve thermal information");
		ERROR_IS_SET = true;
	}

	return NULL;
}

// only used for binding structure (with golang for example)
sys_therm *sys_therm_get_single(sys_therm *sys_therm_to_single, int ask_for)
{
	return &sys_therm_to_single[ask_for];
}

// only used for binding structure (with golang for example)
sys_therm_sensor *sys_therm_sensor_get_single(sys_therm_sensor *sys_therm_sensor_to_single, int ask_for)
{
	return &sys_therm_sensor_to_single[ask_for];
}

// creat new 'sys_therm'structure
sys_therm *sys_therm_new(int count)
{
	sys_therm *curr_sys_therm = malloc(sizeof(sys_therm) * count);
	curr_sys_therm->count = count;
	return curr_sys_therm;
}

// free 'sys_therm' structure
void sys_therm_free(sys_therm *curr_sys_therm)
{
	for (int i = 0; i < curr_sys_therm->count; i++) {
		free(curr_sys_therm[i].sensors);
		curr_sys_therm[i].sensors = NULL;
	}
	free(curr_sys_therm);
	curr_sys_therm = NULL;
}
