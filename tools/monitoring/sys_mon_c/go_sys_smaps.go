// go_sys_smaps.go

package sys_mon

// #include "file_func.h"
// #include "sys_smaps.h"
import "C"

/*
 * Structures functions, this section handle files like 'smaps', 'smaps_rollup',
 * 'maps' is not used here because information are contained inside 'smaps'
 * via 'Header' variable of 'Rollup' or 'Smaps' structures
 */
// Information 'man procfs' search '/smaps' then press 'n' until '/proc/[pid]/smaps'
type Smaps struct {
	Rollup *smapsRollup

	Smaps []smap
	// internal usage only
	smaps *C.smaps
	pid   int

	// Values are given as kB inside parsed files,
	// enable this flag to convert to Bytes
	ConvertToBytes bool
}

// StatNew: Create a new structure that will contains required
// information about 'proc/[pid]/stat' files
// 'maxReadEntries' define the length of the buffer to read
// 'smaps' file, default is set to 2k.
func SmapsNew(pid int, maxReadEntries ...int) (*Smaps, error) {

	var maxEntries = 396
	if len(maxReadEntries) > 0 {
		maxEntries = maxReadEntries[0]
	}

	s := new(Smaps)
	s.pid = pid

	// init and fill smaps_rollup
	s.Rollup = smapsRollupNew()
	err := s.UpdateRollup()
	if err != nil {
		return nil, err
	}

	// init and fill smaps (average number of 2k max entries that will be read)
	// Usually desktops computers have 1k entries, but for servers it could be more.
	s.smaps = C.smaps_get(C.int(pid), C.int(maxEntries), C.bool(s.ConvertToBytes))
	if s.smaps == nil {
		return nil, getErrorString()
	}
	s.Smaps = wrapSmaps(s.smaps)

	return s, nil
}

func wrapSmaps(smaps *C.smaps) []smap {

	count := int(smaps.count)
	var smps = make([]smap, count)
	for i := 0; i < count; i++ {
		c := C.smaps_get_smap(smaps, C.int(i))
		smps[i] = *wrapSmap(c)
	}
	return smps
}

// Update: 'C' structure content with actual values.
func (s *Smaps) UpdateSmaps() error {

	c := C.smaps_update(C.int(s.pid), s.smaps, C.bool(s.ConvertToBytes))
	if !bool(c) {
		return getErrorString()
	}
	s.Smaps = wrapSmaps(s.smaps)
	return nil
}

// Update: 'C' structure content with actual values.
func (s *Smaps) UpdateRollup() error {

	c := C.get_smaps_rollup(C.int(s.pid), s.Rollup.smaps_rollup, C.bool(s.ConvertToBytes))
	if !bool(c) {
		return getErrorString()
	}
	s.Rollup = wrapSmapsRollup(s.Rollup.smaps_rollup)
	return nil
}

// Close: Freeing 'C' structure.
func (s *Smaps) Close() {
	if s.Rollup.smaps_rollup != nil {
		C.smaps_rollup_free(s.Rollup.smaps_rollup)
	}
	if s.smaps != nil {
		C.smaps_free(s.smaps)
	}
}

/*
 * smap
 */
type smap struct {
	smap           *C.smap
	Header         *MapHeader
	Size           uint64
	KernelPageSize uint64
	MmupageSize    uint64
	Rss            uint64
	Pss            uint64
	SharedClean    uint64
	SharedDirty    uint64
	PrivateClean   uint64
	PrivateDirty   uint64
	Referenced     uint64
	Anonymous      uint64
	LazyFree       uint64
	AnonHugePages  uint64
	ShmemPmdMapped uint64
	FilePmdMapped  uint64
	SharedHugetlb  uint64
	PrivateHugetlb uint64
	Swap           uint64
	SwapPss        uint64
	Locked         uint64
	ProtectionKey  int
	VmFlags        string
}

func wrapSmap(csmap *C.smap) *smap {
	if csmap == nil {
		return nil
	}

	return &smap{
		csmap,
		wrapMapHeader(&csmap.header),
		uint64(csmap.size),
		uint64(csmap.kernel_page_size),
		uint64(csmap.mmupage_size),
		uint64(csmap.rss),
		uint64(csmap.pss),
		uint64(csmap.shared_clean),
		uint64(csmap.shared_dirty),
		uint64(csmap.private_clean),
		uint64(csmap.private_dirty),
		uint64(csmap.referenced),
		uint64(csmap.anonymous),
		uint64(csmap.lazy_free),
		uint64(csmap.anon_huge_pages),
		uint64(csmap.shmem_pmd_mapped),
		uint64(csmap.file_pmd_mapped),
		uint64(csmap.shared_hugetlb),
		uint64(csmap.private_hugetlb),
		uint64(csmap.swap),
		uint64(csmap.swap_pss),
		uint64(csmap.locked),
		int(csmap.protection_key),
		C.GoString(&csmap.vm_flags[0]),
	}
}

/*
 * smaps_rollup
 */
type smapsRollup struct {
	smaps_rollup   *C.smaps_rollup
	Header         *MapHeader
	Rss            uint64
	Pss            uint64
	PssAnon        uint64
	PssFile        uint64
	PssShmem       uint64
	SharedClean    uint64
	SharedDirty    uint64
	PrivateClean   uint64
	PrivateDirty   uint64
	Referenced     uint64
	Anonymous      uint64
	LazyFree       uint64
	AnonHugePages  uint64
	ShmemPmdMapped uint64
	FilePmdMapped  uint64
	SharedHugetlb  uint64
	PrivateHugetlb uint64
	Swap           uint64
	SwapPss        uint64
	Locked         uint64
}

func wrapSmapsRollup(smaps_rollup *C.smaps_rollup) *smapsRollup {
	return &smapsRollup{
		smaps_rollup,
		wrapMapHeader(&smaps_rollup.header),
		uint64(smaps_rollup.rss),
		uint64(smaps_rollup.pss),
		uint64(smaps_rollup.pss_anon),
		uint64(smaps_rollup.pss_file),
		uint64(smaps_rollup.pss_shmem),
		uint64(smaps_rollup.shared_clean),
		uint64(smaps_rollup.shared_dirty),
		uint64(smaps_rollup.private_clean),
		uint64(smaps_rollup.private_dirty),
		uint64(smaps_rollup.referenced),
		uint64(smaps_rollup.anonymous),
		uint64(smaps_rollup.lazy_free),
		uint64(smaps_rollup.anon_huge_pages),
		uint64(smaps_rollup.shmem_pmd_mapped),
		uint64(smaps_rollup.file_pmd_mapped),
		uint64(smaps_rollup.shared_hugetlb),
		uint64(smaps_rollup.private_hugetlb),
		uint64(smaps_rollup.swap),
		uint64(smaps_rollup.swap_pss),
		uint64(smaps_rollup.locked),
	}
}

func smapsRollupNew() *smapsRollup {
	c := C.smaps_rollup_new()
	return wrapSmapsRollup(c)
}
