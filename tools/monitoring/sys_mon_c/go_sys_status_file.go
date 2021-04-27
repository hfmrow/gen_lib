// go_sys_status_file.go

package sys_mon

// #include "sys_pid_status.h"
import "C"

// StatusFileNew: create and initialize the "C" structure.
// No need to 'free' (close) anything, everything is already handled.
func StatusFileNew(pid int) (*StatusFile, error) {
	c := C.status_file_new(C.uint(pid))
	if c == nil {
		return nil, getErrorString()
	}
	defer C.status_file_free(c)
	sf := wrapStatusFile(c)
	return sf, nil
}

type StatusFile struct {
	status_file              *C.status_file
	Name                     string
	Umask                    string
	State                    string
	Tgid                     uint
	Ngid                     uint
	Pid                      uint
	Ppid                     uint
	TracerPid                uint
	Uid                      *resfId
	Gid                      *resfId
	FdSize                   uint64
	Groups                   []uint
	NsTgid                   []uint
	NsPid                    []uint
	NsPgid                   []uint
	NsSid                    []uint
	Vm                       *statusFileVmem
	Threads                  int
	SigQ                     string
	SigPnd                   string
	ShdPnd                   string
	SigBlk                   string
	SigIgn                   string
	SigCgt                   string
	CapInh                   string
	CapPrm                   string
	CapEff                   string
	CapBnd                   string
	CapAmb                   string
	NoNewPrivs               int
	Seccomp                  int
	StoreBypass              string
	CpusAllowed              string
	CpusAllowedList          string
	MemsAllowed              string
	MemsAllowedList          string
	VoluntaryCtxtSwitches    uint64
	NonvoluntaryCtxtSwitches uint64
}

func wrapStatusFile(status_file *C.status_file) *StatusFile {
	if status_file == nil {
		return nil
	}

	return &StatusFile{
		status_file,
		C.GoString(&status_file.name[0]),
		C.GoString(&status_file.umask[0]),
		C.GoString(&status_file.state[0]),
		uint(status_file.tgid),
		uint(status_file.ngid),
		uint(status_file.pid),
		uint(status_file.ppid),
		uint(status_file.tracer_pid),
		wrapResfId(&status_file.uid),
		wrapResfId(&status_file.gid),
		uint64(status_file.fd_size),
		wrapUintArray(status_file.groups),
		wrapUintArray(status_file.ns_tgid),
		wrapUintArray(status_file.ns_pid),
		wrapUintArray(status_file.ns_pgid),
		wrapUintArray(status_file.ns_sid),
		wrapStatusFileVmem(status_file.vm),
		int(status_file.threads),
		C.GoString(&status_file.sig_q[0]),
		C.GoString(&status_file.sig_pnd[0]),
		C.GoString(&status_file.shd_pnd[0]),
		C.GoString(&status_file.sig_blk[0]),
		C.GoString(&status_file.sig_ign[0]),
		C.GoString(&status_file.sig_cgt[0]),
		C.GoString(&status_file.cap_inh[0]),
		C.GoString(&status_file.cap_prm[0]),
		C.GoString(&status_file.cap_eff[0]),
		C.GoString(&status_file.cap_bnd[0]),
		C.GoString(&status_file.cap_amb[0]),
		int(status_file.no_new_privs),
		int(status_file.seccomp),
		C.GoString(&status_file.speculation_Store_Bypass[0]),
		C.GoString(&status_file.cpus_allowed[0]),
		C.GoString(&status_file.cpus_allowed_list[0]),
		C.GoString(&status_file.mems_allowed[0]),
		C.GoString(&status_file.mems_allowed_list[0]),
		uint64(status_file.voluntary_ctxt_switches),
		uint64(status_file.nonvoluntary_ctxt_switches),
	}
}

/*
 * uintArray
 */
func wrapUintArray(uint_array *C.uint_array) []uint {
	if uint_array == nil {
		return nil
	}
	count := int(uint_array.count)
	out := make([]uint, count)
	for i := 0; i < count; i++ {
		out[i] = uint(*C.uint_array_pick_value(uint_array, C.int(i)))
	}
	return out
}

/*
 * statusFileVmem
 */
type statusFileVmem struct {
	VmPeak       uint64 // Excepting the two last values
	VmSize       uint64 // the others are given as kB
	VmLck        uint64 // as mentioned in 'man proc'
	VmPin        uint64 // definition.
	VmHwm        uint64
	VmRss        uint64
	RssAnon      uint64
	RssFile      uint64
	RssShmem     uint64
	VmData       uint64
	VmStk        uint64
	VmExe        uint64
	VmLib        uint64
	VmPte        uint64
	VmSwap       uint64
	HugetlbPages uint64
	CoreDumping  int
	ThpEnabled   int
}

func wrapStatusFileVmem(status_file_vmem *C.status_file_vmem) *statusFileVmem {
	if status_file_vmem == nil {
		return nil
	}

	return &statusFileVmem{
		uint64(status_file_vmem.vm_peak),
		uint64(status_file_vmem.vm_size),
		uint64(status_file_vmem.vm_lck),
		uint64(status_file_vmem.vm_pin),
		uint64(status_file_vmem.vm_hwm),
		uint64(status_file_vmem.vm_rss),
		uint64(status_file_vmem.rss_anon),
		uint64(status_file_vmem.rss_file),
		uint64(status_file_vmem.rss_shmem),
		uint64(status_file_vmem.vm_data),
		uint64(status_file_vmem.vm_stk),
		uint64(status_file_vmem.vm_exe),
		uint64(status_file_vmem.vm_lib),
		uint64(status_file_vmem.vm_pte),
		uint64(status_file_vmem.vm_swap),
		uint64(status_file_vmem.hugetlb_pages),
		int(status_file_vmem.core_dumping),
		int(status_file_vmem.thp_enabled),
	}
}

/*
 * statusFileVmem
 */
type resfId struct {
	resf_id    *C.resf_id
	Real       uint
	Effective  uint
	SavedSet   uint
	FileSystem uint
}

func wrapResfId(resf_id *C.resf_id) *resfId {
	if resf_id == nil {
		return nil
	}

	return &resfId{
		resf_id,
		uint(resf_id.real),
		uint(resf_id.effective),
		uint(resf_id.saved_set),
		uint(resf_id.file_system),
	}
}
