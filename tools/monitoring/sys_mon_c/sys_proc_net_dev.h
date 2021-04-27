// sys_proc_net_dev.h

#ifndef FUNCTIONS_NET_DEV_INCLUDED
#define FUNCTIONS_NET_DEV_INCLUDED

#include <stdbool.h>
#include <time.h>
#include <stdio.h>

typedef struct net_interface_tx {
	unsigned long long bytes;
	unsigned long long packets;
	unsigned long long errs;
	unsigned long long drop;
	unsigned long long fifo;
	unsigned long long colls;
	unsigned long long carrier;
	unsigned long long compressed;
} net_interface_tx;

typedef struct net_interface_rx {
	unsigned long long bytes;
	unsigned long long packets;
	unsigned long long errs;
	unsigned long long drop;
	unsigned long long fifo;
	unsigned long long frame;
	unsigned long long compressed;
	unsigned long long multicast;
} net_interface_rx;

typedef struct iface {
	// time_t tv_sec 	whole seconds (valid values are >= 0),
	// long tv_nsec 	nanoseconds (valid values are [0, 999999999])
	struct timespec last_update;
	double delta_sec;
	unsigned long delta_tx;
	unsigned long delta_rx;
	double tx_byte_sec;
	double rx_byte_sec;
	char tx_string[64];
	char rx_string[64];
	net_interface_tx tx;
	net_interface_rx rx;
	char name[64];
} iface;

typedef struct proc_net_dev {
	char path[64];
	int count;
	long pid; // -1 means non-specific (global flow)
	char suffix[16];
	char unit[16];
	iface *interfaces;
} proc_net_dev;

// only used for binding structure (with golang for example)
iface *iface_get_single(iface *curr_iface, int ask_for);

// // only used for binding structure (with golang for example)
//void iface_set_suffix(proc_net_dev *curr_proc_net_dev, char *suffix);

// // only used for binding structure (with golang for example)
//void iface_set_unit(proc_net_dev *curr_proc_net_dev, char *unit);

// read values from a 'proc/net/dev' or 'proc/PID/net/dev',
// and fill the given 'iface' structure. 'iface_name'
// returned with network card name.
bool net_dev_read_values(FILE *fp, char *path, iface *curr_iface);

// update a 'proc_net_dev' structure
bool proc_net_dev_update(proc_net_dev *curr_proc_net_dev);

// initialize/fill a 'proc_net_dev' structure
bool proc_net_dev_init(proc_net_dev *curr_proc_net_dev);

// initialize/fill and get a new 'proc_net_dev' structure
proc_net_dev *proc_net_dev_get(unsigned int *pid);

// creat new 'proc_net_dev'structure
proc_net_dev *proc_net_dev_new();

// free 'proc_net_dev' structure
void proc_net_dev_free(proc_net_dev *curr_proc_net_dev);

#endif
