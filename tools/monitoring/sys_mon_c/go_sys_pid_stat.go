// go_sys_pid_stat.go

package sys_mon

// #include "sys_pid_stat.h"
import "C"

// ProcPidStatNew: Create and initialise 'C' structure.
func ProcPidStatNew(pid int) (*ProcPidStat, error) {
	c := C.proc_pid_stat_get(C.uint(pid))
	if c == nil {
		return nil, getErrorString()
	}
	return wrapProcPidStat(c), nil
}

// Update: 'ProcPidStat' structure.
func (s *ProcPidStat) Update() error {
	c := C.proc_pid_stat_update(s.proc_pid_stat)
	if !bool(c) {
		return getErrorString()
	}
	return nil
}

// Close: Freeing 'C' structure.
func (s *ProcPidStat) Close() {
	if s.proc_pid_stat != nil {
		C.proc_pid_stat_free(s.proc_pid_stat)
	}
}

type ProcPidStat struct {
	proc_pid_stat       *C.proc_pid_stat
	Pid                 uint
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

func wrapProcPidStat(proc_pid_stat *C.proc_pid_stat) *ProcPidStat {
	if proc_pid_stat == nil {
		return nil
	}

	return &ProcPidStat{
		proc_pid_stat,
		uint(proc_pid_stat.pid),
		C.GoString(&proc_pid_stat.comm[0]),
		C.GoString(&proc_pid_stat.state[0]),
		int(proc_pid_stat.ppid),
		int(proc_pid_stat.pgrp),
		int(proc_pid_stat.session),
		int(proc_pid_stat.tty_nr),
		int(proc_pid_stat.tpgid),
		uint(proc_pid_stat.flags),
		uint32(proc_pid_stat.minflt),
		uint32(proc_pid_stat.cminflt),
		uint32(proc_pid_stat.majflt),
		uint32(proc_pid_stat.cmajflt),
		uint32(proc_pid_stat.utime),
		uint32(proc_pid_stat.stime),
		int32(proc_pid_stat.cutime),
		int32(proc_pid_stat.cstime),
		int32(proc_pid_stat.priority),
		int32(proc_pid_stat.nice),
		int32(proc_pid_stat.num_threads),
		int32(proc_pid_stat.itrealvalue),
		uint64(proc_pid_stat.starttime),
		uint32(proc_pid_stat.vsize),
		int32(proc_pid_stat.rss),
		uint32(proc_pid_stat.rsslim),
		uint32(proc_pid_stat.startcode),
		uint32(proc_pid_stat.endcode),
		uint32(proc_pid_stat.startstack),
		uint32(proc_pid_stat.kstkesp),
		uint32(proc_pid_stat.kstkeip),
		uint32(proc_pid_stat.signal),
		uint32(proc_pid_stat.blocked),
		uint32(proc_pid_stat.sigignore),
		uint32(proc_pid_stat.sigcatch),
		uint32(proc_pid_stat.wchan),
		uint32(proc_pid_stat.nswap),
		uint32(proc_pid_stat.cnswap),
		int(proc_pid_stat.exit_signal),
		int(proc_pid_stat.processor),
		uint(proc_pid_stat.rt_priority),
		uint(proc_pid_stat.policy),
		uint64(proc_pid_stat.delayacct_blkio_ticks),
		uint32(proc_pid_stat.guest_time),
		int32(proc_pid_stat.cguest_time),
		uint32(proc_pid_stat.start_data),
		uint32(proc_pid_stat.end_data),
		uint32(proc_pid_stat.start_brk),
		uint32(proc_pid_stat.arg_start),
		uint32(proc_pid_stat.arg_end),
		uint32(proc_pid_stat.env_start),
		uint32(proc_pid_stat.env_end),
		int(proc_pid_stat.exit_code),
	}
}
