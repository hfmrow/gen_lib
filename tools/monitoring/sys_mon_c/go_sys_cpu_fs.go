// go_sys_cpu_fs.go

package sys_mon

// #include "sys_cpu.h"
import "C"

// Structure to hold values retrieved from:
// '/sys/devices/system/cpu/cpufreq/policy*' directories
type CpuFs struct {
	cpuFsList     *C.cpu_fs
	cpuFsCurrFreq *C.cpu_fs

	CpuCount int
	CurrFreq []int64
	CpuList  []cpuFs
}

func wrapCpuFsList(cpu_fs *C.cpu_fs) []cpuFs {
	cfl := make([]cpuFs, int(cpu_fs.cpu_count))
	for i := 0; i < len(cfl); i++ {
		c := C.cpu_fs_get_single(cpu_fs, C.int(i))
		cfl[i] = *wrapCpuFs(c)
	}
	return cfl
}

func wrapCpuFsCurrFreq(cpu_fs *C.cpu_fs) []int64 {
	cfcf := make([]int64, int(cpu_fs.cpu_count))
	for i := 0; i < len(cfcf); i++ {
		c := C.cpu_fs_get_single(cpu_fs, C.int(i))
		cfcf[i] = int64(c.scaling_cur_freq)
	}
	return cfcf
}

func CpuFsNew() (*CpuFs, error) {
	cf := new(CpuFs)
	c := C.cpu_fs_get()
	if c == nil {
		return nil, getErrorString()
	}
	cf.cpuFsList = c
	c = C.cpu_fs_get()
	if c == nil {
		return nil, getErrorString()
	}
	cf.cpuFsCurrFreq = c
	cf.CpuCount = int(cf.cpuFsList.cpu_count)
	cf.CpuList = wrapCpuFsList(cf.cpuFsList)
	cf.CurrFreq = wrapCpuFsCurrFreq(cf.cpuFsCurrFreq)
	return cf, nil
}

func (cf *CpuFs) CurrFreqUpdate() error {
	c := C.cpu_fs_curr_freq_update(cf.cpuFsCurrFreq)
	if !bool(c) {
		return getErrorString()
	}
	cf.CurrFreq = wrapCpuFsCurrFreq(cf.cpuFsCurrFreq)
	return nil
}

func (cf *CpuFs) ListUpdate() error {
	c := C.cpu_fs_update(cf.cpuFsList)
	if !bool(c) {
		return getErrorString()
	}
	cf.CpuList = wrapCpuFsList(cf.cpuFsList)
	return nil
}

func (cf *CpuFs) Close() {
	if cf.cpuFsList != nil {
		C.cpu_fs_free(cf.cpuFsList)
	}
	if cf.cpuFsCurrFreq != nil {
		C.cpu_fs_free(cf.cpuFsCurrFreq)
	}
}

type cpuFs struct {
	cpu_fs                                *C.cpu_fs
	BiosLimit                             int64
	BaseFrequency                         int64
	CpuinfoCurFreq                        int64
	CpuinfoMinFreq                        int64
	CpuinfoMaxFreq                        int64
	CpuinfoTransitionLatency              int64
	ScalingCurFreq                        int64
	ScalingMinFreq                        int64
	ScalingMaxFreq                        int64
	ScalingSetspeed                       int64
	ScalingAvailableFrequencies           string
	RelatedCpus                           string
	ScalingAvailableGovernors             string
	EnergyPerformanceAvailablePreferences string
	EnergyPerformancePreference           string
	ScalingDriver                         string
	ScalingGovernor                       string
}

func wrapCpuFs(cpu_fs *C.cpu_fs) *cpuFs {
	if cpu_fs == nil {
		return nil
	}

	return &cpuFs{
		cpu_fs,
		int64(cpu_fs.bios_limit),
		int64(cpu_fs.base_frequency),
		int64(cpu_fs.cpuinfo_cur_freq),
		int64(cpu_fs.cpuinfo_min_freq),
		int64(cpu_fs.cpuinfo_max_freq),
		int64(cpu_fs.cpuinfo_transition_latency),
		int64(cpu_fs.scaling_cur_freq),
		int64(cpu_fs.scaling_min_freq),
		int64(cpu_fs.scaling_max_freq),
		int64(cpu_fs.scaling_setspeed),
		C.GoString(cpu_fs.scaling_available_frequencies),
		C.GoString(cpu_fs.related_cpus),
		C.GoString(cpu_fs.scaling_available_governors),
		C.GoString(cpu_fs.energy_performance_available_preferences),
		C.GoString(cpu_fs.energy_performance_preference),
		C.GoString(cpu_fs.scaling_driver),
		C.GoString(cpu_fs.scaling_governor),
	}
}
