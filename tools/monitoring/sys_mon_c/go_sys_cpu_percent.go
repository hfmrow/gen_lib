// go_cpu_percent_pid.go

package sys_mon

// #include "sys_pid_stat.h"
// #include "sys_cpu_percent.h"
import "C"

// CpuPercentPidNew: Create and initialise 'C' structure.
func CpuPercentPidNew(pid int) (*CpuPercentPid, error) {
	c := C.cpu_percent_pid_get(C.uint(pid))
	if c == nil {
		return nil, getErrorString()
	}
	return wrapCpuPercentPid(c), nil
}

// update current values
func (v *CpuPercentPid) Update() error {
	c := C.cpu_percent_pid_update(v.cpu_percent_pid)
	if bool(c) {
		v.CpuPercent = float32(v.cpu_percent_pid.cpu_percent)
		v.MemoryRss = int64(v.cpu_percent_pid.memory_rss)
		return nil
	}
	return getErrorString()
}

type CpuPercentPid struct {
	cpu_percent_pid *C.cpu_percent_pid
	CpuPercent      float32
	MemoryRss       int64
}

func wrapCpuPercentPid(cpu_percent_pid *C.cpu_percent_pid) *CpuPercentPid {
	if cpu_percent_pid == nil {
		return nil
	}

	return &CpuPercentPid{
		cpu_percent_pid,
		float32(cpu_percent_pid.cpu_percent),
		int64(cpu_percent_pid.memory_rss),
	}
}

// Close: Freeing 'C' structure.
func (s *CpuPercentPid) Close() {
	if s.cpu_percent_pid != nil {
		C.cpu_percent_pid_free(s.cpu_percent_pid)
	}
}
