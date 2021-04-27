//+build ignore
// go_sys_stat_pid.go

package sys_mon

// #include "sys_pid_stat.h"
// #include "file_func.h"
import "C"

type Stat struct {
	pid int
	Current,
	last *StoreStat
}

// StatNew: Create a new structure that will contains required
// information about 'proc/[pid]/stat' files
func StatNew(pid int) *Stat {

	s := new(Stat)
	s.Current = storeStatNew()
	s.last = storeStatNew()
	s.pid = pid
	return s
}

// Update: 'C' structure content with actual values.
func (s *Stat) Update() error {
	c := C.stat_process_get(C.int(s.pid), s.Current.native(), s.last.native())
	if !bool(c) {
		return getErrorString()
	}
	s.Current = wrapStoreStat(s.Current.store_stat)
	return nil
}

// Close: Freeing 'C' structure.
func (s *Stat) Close() {
	if s.Current.store_stat != nil {
		C.store_stat_free(s.Current.store_stat)
	}
	if s.last.store_stat != nil {
		C.store_stat_free(s.last.store_stat)
	}
}

/*
 * stat
 */
type StoreStat struct {
	store_stat *C.store_stat
	Cpu        *storeStatCpu
	Pid        *StoreStatPid
	CpuPercent float64
	// Resident set size: memory rss (page_size * rss)
	MemoryRss int64
}

func (v *StoreStat) native() *C.store_stat {
	if v == nil || v.store_stat == nil {
		return nil
	}

	return v.store_stat
}

func wrapStoreStat(store_stat *C.store_stat) *StoreStat {
	return &StoreStat{
		store_stat,
		wrapStoreStatCpu(store_stat.cpu),
		wrapStoreStatPid(store_stat.pid),
		float64(store_stat.cpu_percent),
		int64(store_stat.memory_rss),
	}
}

func storeStatNew() *StoreStat {
	c := C.store_stat_new()
	return wrapStoreStat(c)
}

/*
 * stat_pid
 */
// Information 'man procfs' search '/stat' then press 'n' until '/proc/[pid]/stat'
type StoreStatPid struct {
	store_stat_pid      *C.store_stat_pid
	Pid                 int
	Comm                string
	State               string
	Ppid                int
	Pgrp                int
	Session             int
	TtyNr               int
	Tpgid               int
	Flags               uint
	Minflt              uint32
	Cminflt             uint32
	Majflt              uint32
	Cmajflt             uint32
	Utime               uint32
	Stime               uint32
	Cutime              int32
	Cstime              int32
	Priority            int32
	Nice                int32
	NumThreads          int32
	Itrealvalue         int32
	Starttime           uint64
	Vsize               uint32
	Rss                 int32
	Rsslim              uint32
	Startcode           uint32
	Endcode             uint32
	Startstack          uint32
	Kstkesp             uint32
	Kstkeip             uint32
	Signal              uint32
	Blocked             uint32
	Sigignore           uint32
	Sigcatch            uint32
	Wchan               uint32
	Nswap               uint32
	Cnswap              uint32
	ExitSignal          int
	Processor           int
	RtPriority          uint
	Policy              uint
	DelayacctBlkioTicks uint64
	GuestTime           uint32
	CguestTime          int32
	StartData           uint32
	EndData             uint32
	StartBrk            uint32
	ArgStart            uint32
	ArgEnd              uint32
	EnvStart            uint32
	EnvEnd              uint32
	ExitCode            int
}

func (v *StoreStatPid) native() *C.store_stat_pid {
	if v == nil || v.store_stat_pid == nil {
		return nil
	}

	return v.store_stat_pid
}

func wrapStoreStatPid(store_stat_pid *C.store_stat_pid) *StoreStatPid {
	return &StoreStatPid{
		store_stat_pid,
		int(store_stat_pid.pid),
		C.GoString(&store_stat_pid.comm[0]),
		C.GoString(&store_stat_pid.state[0]),
		int(store_stat_pid.ppid),
		int(store_stat_pid.pgrp),
		int(store_stat_pid.session),
		int(store_stat_pid.tty_nr),
		int(store_stat_pid.tpgid),
		uint(store_stat_pid.flags),
		uint32(store_stat_pid.minflt),
		uint32(store_stat_pid.cminflt),
		uint32(store_stat_pid.majflt),
		uint32(store_stat_pid.cmajflt),
		uint32(store_stat_pid.utime),
		uint32(store_stat_pid.stime),
		int32(store_stat_pid.cutime),
		int32(store_stat_pid.cstime),
		int32(store_stat_pid.priority),
		int32(store_stat_pid.nice),
		int32(store_stat_pid.num_threads),
		int32(store_stat_pid.itrealvalue),
		uint64(store_stat_pid.starttime),
		uint32(store_stat_pid.vsize),
		int32(store_stat_pid.rss),
		uint32(store_stat_pid.rsslim),
		uint32(store_stat_pid.startcode),
		uint32(store_stat_pid.endcode),
		uint32(store_stat_pid.startstack),
		uint32(store_stat_pid.kstkesp),
		uint32(store_stat_pid.kstkeip),
		uint32(store_stat_pid.signal),
		uint32(store_stat_pid.blocked),
		uint32(store_stat_pid.sigignore),
		uint32(store_stat_pid.sigcatch),
		uint32(store_stat_pid.wchan),
		uint32(store_stat_pid.nswap),
		uint32(store_stat_pid.cnswap),
		int(store_stat_pid.exit_signal),
		int(store_stat_pid.processor),
		uint(store_stat_pid.rt_priority),
		uint(store_stat_pid.policy),
		uint64(store_stat_pid.delayacct_blkio_ticks),
		uint32(store_stat_pid.guest_time),
		int32(store_stat_pid.cguest_time),
		uint32(store_stat_pid.start_data),
		uint32(store_stat_pid.end_data),
		uint32(store_stat_pid.start_brk),
		uint32(store_stat_pid.arg_start),
		uint32(store_stat_pid.arg_end),
		uint32(store_stat_pid.env_start),
		uint32(store_stat_pid.env_end),
		int(store_stat_pid.exit_code),
	}
}
