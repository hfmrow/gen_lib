// sys_proc_net_dev.c

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "file_func.h"
#include "sys_proc_net_dev.h"

// read values from a 'proc/net/dev' or 'proc/PID/net/dev',
// and fill given structure
bool net_dev_read_values(FILE *fp, char *path, iface *curr_iface)
{
	bool ret = is_correctly_read(
	               fp,
	               path,
	               fscanf(fp, "%s%llu%llu%llu%llu%llu%llu%llu%llu %llu%llu%llu%llu%llu%llu%llu%llu",
	                      curr_iface->name,
	                      &curr_iface->rx.bytes,
	                      &curr_iface->rx.packets,
	                      &curr_iface->rx.errs,
	                      &curr_iface->rx.drop,
	                      &curr_iface->rx.fifo,
	                      &curr_iface->rx.frame,
	                      &curr_iface->rx.compressed,
	                      &curr_iface->rx.multicast,
	                      &curr_iface->tx.bytes,
	                      &curr_iface->tx.packets,
	                      &curr_iface->tx.errs,
	                      &curr_iface->tx.drop,
	                      &curr_iface->tx.fifo,
	                      &curr_iface->tx.colls,
	                      &curr_iface->tx.carrier,
	                      &curr_iface->tx.compressed), 17);
	if ( !ret)
		goto end;
	//char *cut_at = ":";
	cut_at_first(curr_iface->name, ":");
	curr_iface->last_update = time_spec_get();

end:
	return ret;
}

char *human_readable_bandwidth(double bytes, char *suffix, char *unit, char *string)
{
	if (bytes < 2<<9)
		sprintf(string, "%dB%s", (int)bytes, unit);
	else if (bytes < (unsigned long long)2<<19)
		sprintf(string, "%0.2fK%s%s", bytes / 1024, suffix, unit);
	else if (bytes < (unsigned long long)2<<29)
		sprintf(string, "%0.2fM%s%s", bytes / 1024e3, suffix, unit);
	else if (bytes < (unsigned long long)2<<39)
		sprintf(string, "%0.2fG%s%s", bytes / 1024e6, suffix, unit);
	else if (bytes < (unsigned long long)2<<49)
		sprintf(string, "%0.2fT%s%s", bytes / 1024e9, suffix, unit);

	return string;
}

// update a 'proc_net_dev' structure
bool proc_net_dev_update(proc_net_dev *curr_proc_net_dev)
{
	FILE *fp;
	double time_diff = 0;
	//double corrected_delta_sec;
	bool ret = false;

	if ((fp = fopen(curr_proc_net_dev->path, "r")) != NULL) {
		seek_to_next_line(fp);
		seek_to_next_line(fp);

		iface *tmp_interface = malloc(sizeof(iface) * curr_proc_net_dev->count);

		for (int i=0; i<curr_proc_net_dev->count; i++)
			if ( !net_dev_read_values(fp, curr_proc_net_dev->path, &tmp_interface[i]))
				goto end;
			else {
				// update time
				tmp_interface[i].last_update = time_spec_get();

				// compute transfert rate
				tmp_interface[i].delta_tx =
				    tmp_interface[i].tx.bytes - curr_proc_net_dev->interfaces[i].tx.bytes;
				tmp_interface[i].delta_rx =
				    tmp_interface[i].rx.bytes - curr_proc_net_dev->interfaces[i].rx.bytes;
				tmp_interface[i].delta_sec = *time_spec_sec_sub(
				                                 &tmp_interface[i].last_update,
				                                 &curr_proc_net_dev->interfaces[i].last_update,
				                                 &time_diff);

				tmp_interface[i].tx_byte_sec = tmp_interface[i].delta_tx / tmp_interface[i].delta_sec;
				tmp_interface[i].rx_byte_sec = tmp_interface[i].delta_rx / tmp_interface[i].delta_sec;

				human_readable_bandwidth(
				    tmp_interface[i].tx_byte_sec,
				    curr_proc_net_dev->suffix,
				    curr_proc_net_dev->unit,
				    tmp_interface[i].tx_string);
				human_readable_bandwidth(
				    tmp_interface[i].rx_byte_sec,
				    curr_proc_net_dev->suffix,
				    curr_proc_net_dev->unit,
				    tmp_interface[i].rx_string);
			}

		memcpy(curr_proc_net_dev->interfaces, tmp_interface, sizeof(iface) * curr_proc_net_dev->count);
		free(tmp_interface);

	} else {
		sprintf(ERROR_MESSAGE, "Error: unable to open [%s]", curr_proc_net_dev->path);
		ERROR_IS_SET = true;
		goto end;
	}

	ret = true;
end:
	fclose(fp);
	return ret;
}

// initialize/fill a 'proc_net_dev' structure
bool proc_net_dev_init(proc_net_dev *curr_proc_net_dev)
{
	FILE *fp;
	bool ret;

	if (curr_proc_net_dev->pid > -1)
		sprintf(curr_proc_net_dev->path, "/proc/%ld/net/dev", curr_proc_net_dev->pid);
	else
		strcpy(curr_proc_net_dev->path, "/proc/net/dev");

	if ((fp = fopen(curr_proc_net_dev->path, "r")) != NULL) {
		seek_to_next_line(fp);
		seek_to_next_line(fp);

		iface *tmp_interface = malloc(sizeof(iface) * 150);
		int size = sizeof(tmp_interface->rx_string); // char[64]

		while(net_dev_read_values(fp, curr_proc_net_dev->path, &tmp_interface[curr_proc_net_dev->count])) {
			// Clean char[64] to be sure we does not have strange characters on 1st run
			memset(tmp_interface[curr_proc_net_dev->count].rx_string, 0, size);
			memset(tmp_interface[curr_proc_net_dev->count].tx_string, 0, size);
			curr_proc_net_dev->count++;
		}

		curr_proc_net_dev->interfaces = malloc(sizeof(iface) * curr_proc_net_dev->count);
		memcpy(curr_proc_net_dev->interfaces, tmp_interface, sizeof(iface) * curr_proc_net_dev->count);

		free(tmp_interface);
		tmp_interface = NULL;

		if (ERROR_IS_SET) {
			ret = false;
			goto end;

		} else {
			ret = true;
			goto end;
		}
	}

	memset(ERROR_MESSAGE, 0, sizeof(ERROR_MESSAGE));
	sprintf(ERROR_MESSAGE, "Unable to open: %s", curr_proc_net_dev->path);
	ERROR_IS_SET = true;
	ret = false;

end:
	return ret;
}

// initialize/fill and get a new 'proc_net_dev' structure
// // note: pid = NULL means non-specific (global flow)
proc_net_dev *proc_net_dev_get(unsigned int *pid)
{
	proc_net_dev *curr_proc_net_dev = proc_net_dev_new();
	curr_proc_net_dev->count = 0;
	if (pid != NULL)
		curr_proc_net_dev->pid = *pid;
	else
		curr_proc_net_dev->pid = -1;

	strcpy(curr_proc_net_dev->suffix, "iB");
	strcpy(curr_proc_net_dev->unit, "/s");

	if (proc_net_dev_init(curr_proc_net_dev))
		return curr_proc_net_dev;

	return NULL;
}

// only used for binding structure (with golang for example)
iface *iface_get_single(iface *curr_iface, int ask_for)
{
	return &curr_iface[ask_for];
}

// //only used for binding structure (with golang for example)
//void iface_set_suffix(proc_net_dev *curr_proc_net_dev, char *suffix)
//{
//	strcpy(&curr_proc_net_dev->suffix[0], suffix);
//}
//
// //only used for binding structure (with golang for example)
//void iface_set_unit(proc_net_dev *curr_proc_net_dev, char *unit)
//{
//	strcpy(&curr_proc_net_dev->unit[0], unit);
//}

// creat new 'proc_net_dev'structure
proc_net_dev *proc_net_dev_new()
{
	proc_net_dev *curr_proc_net_dev = malloc(sizeof(proc_net_dev));
	return curr_proc_net_dev;
}

// free 'proc_net_dev' structure
void proc_net_dev_free(proc_net_dev *curr_proc_net_dev)
{
	free(curr_proc_net_dev->interfaces);
	curr_proc_net_dev->interfaces = NULL;

	free(curr_proc_net_dev);
	curr_proc_net_dev = NULL;
}
