// cpu_usage.go

package sys_monitor

// #include <time.h>
/*
#include "stdio.h"
#include "stdlib.h"
#include "unistd.h"
#include "ctype.h"

const int PATH_LENGTH = 60;

char blank_read[60];

// Read error checking.
void is_correctly_read(FILE *fp, char path[PATH_LENGTH], int read, int expected)
{
	if (read != expected) {
		fclose(fp);
		printf("Error, file read: [%s], expected: %d, read: %d variable(s)\n", path, expected, read);
		exit(1);
	}
}

// Place the file reader to desired field
void seek_to_field(FILE *fp, char path[PATH_LENGTH], int field)
{
	for(int p = 1; p < field; p++) {
		is_correctly_read(fp, path, fscanf(fp, "%s", blank_read), 1);
		//printf("%s\n", blank_read);
	}
}

// cpu % calculation.
float cpu_percent_pid(int sleep_time, int i_pid)
{
	FILE *fp;
	char path[PATH_LENGTH];
	long int  q1, r1, q2, r2, q, r, x, y, z, w;
	long int x1, y1, x2, y2, z1, z2, w2, w1;

	snprintf(path, PATH_LENGTH, "%s%d%s", "/proc/", i_pid, "/stat");
	fp = fopen(path, "r");
	seek_to_field(fp, path, 14);
	is_correctly_read(fp, path, fscanf(fp, "%ld%ld", &q1, &r1), 2);
	fclose(fp);
	snprintf(path, PATH_LENGTH, "%s%s", "/proc", "/stat");
	fp = fopen(path, "r");
	seek_to_field(fp, path, 2);
	is_correctly_read(fp, path, fscanf(fp, "%ld%ld%ld%ld", &x1, &y1, &z1, &w1), 4);
	fclose(fp);

	sleep(sleep_time);

	snprintf(path, PATH_LENGTH, "%s%d%s", "/proc/", i_pid, "/stat");
	fp = fopen(path, "r");
	seek_to_field(fp, path, 14);
	is_correctly_read(fp, path, fscanf(fp, "%ld%ld", &q2, &r2), 2);
	fclose(fp);
	snprintf(path, PATH_LENGTH, "%s%s", "/proc", "/stat");
	fp = fopen(path, "r");
	seek_to_field(fp, path, 2);
	is_correctly_read(fp, path, fscanf(fp, "%ld%ld%ld%ld", &x2, &y2, &z2, &w2), 4);
	fclose(fp);

	q = q2-q1;
	r = r2-r1;
	x = x2-x1;
	y = y2-y1;
	z = z2-z1;
	w = w2-w1;

	return (q+r)*100.0/(x+y+z+w);
}
*/
import "C"
import (
	"runtime"
	"time"
)

var (
	startTime  = time.Now()
	startTicks = C.clock()
	cpuCount   = float64(runtime.NumCPU())
)

func CpuPercentPid(sleepTime, pid int) float64 {
	return float64(C.cpu_percent_pid(C.int(sleepTime), C.int(pid)))
}

func CpuPercentCurrProcess() float64 {
	clockSeconds := float64(C.clock()-startTicks) / float64(C.CLOCKS_PER_SEC)
	realSeconds := time.Since(startTime).Seconds()
	return (clockSeconds / realSeconds * 100) / cpuCount
}
