// go_sys_meminfo.go

package sys_mon

// #include "sys_meminfo.h"
import "C"

// Close: Freeing 'C' structure.
func (s *Meminfo) Close() {
	if s.meminfo != nil {
		C.meminfo_free(s.meminfo)
	}
}

// MeminfoNew: Create and initialise 'C' structure.
func MeminfoNew() (*Meminfo, error) {
	c := C.meminfo_get()
	if c == nil {
		return nil, getErrorString()
	}
	return wrapMeminfo(c), nil
}

// Update: 'Meminfo' structure.
func (s *Meminfo) Update() (*Meminfo, error) {
	c := C.meminfo_update(s.meminfo)
	if !bool(c) {
		return nil, getErrorString()
	}
	return wrapMeminfo(s.meminfo), nil
}

type Meminfo struct {
	meminfo           *C.meminfo
	MemTotal          uint32
	MemFree           uint32
	MemAvailable      uint32
	Buffers           uint32
	Cached            uint32
	SwapCached        uint32
	Active            uint32
	Inactive          uint32
	ActiveAnon        uint32
	InactiveAnon      uint32
	ActiveFile        uint32
	InactiveFile      uint32
	Unevictable       uint32
	Mlocked           uint32
	SwapTotal         uint32
	SwapFree          uint32
	Dirty             uint32
	Writeback         uint32
	AnonPages         uint32
	Mapped            uint32
	Shmem             uint32
	Kreclaimable      uint32
	Slab              uint32
	Sreclaimable      uint32
	Sunreclaim        uint32
	KernelStack       uint32
	PageTables        uint32
	NfsUnstable       uint32
	Bounce            uint32
	WritebackTmp      uint32
	CommitLimit       uint32
	CommittedAs       uint32
	VmallocTotal      uint32
	VmallocUsed       uint32
	VmallocChunk      uint32
	Percpu            uint32
	HardwareCorrupted uint32
	AnonHugePages     uint32
	ShmemHugePages    uint32
	ShmemPmdMapped    uint32
	FileHugePages     uint32
	FilePmdMapped     uint32
	HugePagesTotal    uint32
	HugePagesFree     uint32
	HugePagesRsvd     uint32
	HugePagesSurp     uint32
	Hugepagesize      uint32
	Hugetlb           uint32
	DirectMap4k       uint32
	DirectMap2M       uint32
	DirectMap1G       uint32
}

func (v *Meminfo) native() *C.meminfo {
	if v == nil || v.meminfo == nil {
		return nil
	}

	return v.meminfo
}

func wrapMeminfo(meminfo *C.meminfo) *Meminfo {
	if meminfo == nil {
		return nil
	}

	return &Meminfo{
		meminfo,
		uint32(meminfo.mem_total),
		uint32(meminfo.mem_free),
		uint32(meminfo.mem_available),
		uint32(meminfo.buffers),
		uint32(meminfo.cached),
		uint32(meminfo.swap_cached),
		uint32(meminfo.active),
		uint32(meminfo.inactive),
		uint32(meminfo.active_anon_),
		uint32(meminfo.inactive_anon_),
		uint32(meminfo.active_file_),
		uint32(meminfo.inactive_file_),
		uint32(meminfo.unevictable),
		uint32(meminfo.mlocked),
		uint32(meminfo.swap_total),
		uint32(meminfo.swap_free),
		uint32(meminfo.dirty),
		uint32(meminfo.writeback),
		uint32(meminfo.anon_pages),
		uint32(meminfo.mapped),
		uint32(meminfo.shmem),
		uint32(meminfo.kreclaimable),
		uint32(meminfo.slab),
		uint32(meminfo.sreclaimable),
		uint32(meminfo.sunreclaim),
		uint32(meminfo.kernel_stack),
		uint32(meminfo.page_tables),
		uint32(meminfo.nfs_unstable),
		uint32(meminfo.bounce),
		uint32(meminfo.writeback_tmp),
		uint32(meminfo.commit_limit),
		uint32(meminfo.committed_as),
		uint32(meminfo.vmalloc_total),
		uint32(meminfo.vmalloc_used),
		uint32(meminfo.vmalloc_chunk),
		uint32(meminfo.percpu),
		uint32(meminfo.hardware_corrupted),
		uint32(meminfo.anon_huge_pages),
		uint32(meminfo.shmem_huge_pages),
		uint32(meminfo.shmem_pmd_mapped),
		uint32(meminfo.file_huge_pages),
		uint32(meminfo.file_pmd_mapped),
		uint32(meminfo.huge_pages_total),
		uint32(meminfo.huge_pages_free),
		uint32(meminfo.huge_pages_rsvd),
		uint32(meminfo.huge_pages_surp),
		uint32(meminfo.hugepagesize),
		uint32(meminfo.hugetlb),
		uint32(meminfo.direct_map4k),
		uint32(meminfo.direct_map2_m),
		uint32(meminfo.direct_map1_g),
	}
}
