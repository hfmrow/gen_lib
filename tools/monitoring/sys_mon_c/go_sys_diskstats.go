// go_sys_diskstats.go

package sys_mon

// #include "sys_partitions.h"
import "C"

type Diskstats struct {
	diskstats *C.diskstats
	Details   []gdiskstats
}

// Close: Freeing 'C' structure.
func (s *Diskstats) Close() {
	if s.diskstats != nil {
		C.diskstats_free(s.diskstats)
	}
}

// DiskstatsNew: Create and initialise 'C' structure.
func DiskstatsNew() (*Diskstats, error) {
	s := new(Diskstats)
	s.diskstats = C.diskstats_get()
	if s.diskstats == nil {
		return nil, getErrorString()
	}
	s.wrap()
	return s, nil
}

func (s *Diskstats) wrap() {
	count := int(s.diskstats.count)
	s.Details = make([]gdiskstats, count)
	for i := 0; i < count; i++ {
		s.Details[i] = *wrapDiskstats(C.diskstats_get_single(s.diskstats, C.int(i)))
	}
}

// Update: 'Diskstats' structure.
func (s *Diskstats) Update() error {
	c := C.diskstats_update(s.diskstats)
	if !bool(c) {
		return getErrorString()
	}
	s.wrap()
	return nil
}

type gdiskstats struct {
	Device            string
	DevType           string
	Count             int
	ReadsCompleted    uint32
	ReadsMerged       uint32
	SectorsRead       uint32
	TimeReadingMs     uint32
	WritesCompleted   uint32
	WritesMerged      uint32
	SectorsWritten    uint32
	TimeWritingMs     uint32
	IosInProgress     uint32
	TimeDoingIosMs    uint32
	WeightedTimeIosMs uint32
	DiscardsCompleted uint32
	DiscardsMerged    uint32
	SectorsDiscarded  uint32
	SpentDiscardingMs uint32
	Undefined1        uint32
	Undefined2        uint32
	Undefined3        uint32
}

func wrapDiskstats(diskstats *C.diskstats) *gdiskstats {
	if diskstats == nil {
		return nil
	}

	return &gdiskstats{
		C.GoString(&diskstats.device[0]),
		C.GoString(&diskstats.dev_type[0]),
		int(diskstats.count),
		uint32(diskstats.reads_completed),
		uint32(diskstats.reads_merged),
		uint32(diskstats.sectors_read),
		uint32(diskstats.time_reading_ms),
		uint32(diskstats.writes_completed),
		uint32(diskstats.writes_merged),
		uint32(diskstats.sectors_written),
		uint32(diskstats.time_writing_ms),
		uint32(diskstats.ios_in_progress),
		uint32(diskstats.time_doing_ios_ms),
		uint32(diskstats.weighted_time_ios_ms),
		// the following values will be read if they are present
		// otherwise they will be ignored (depend on system)
		uint32(diskstats.discards_completed),
		uint32(diskstats.discards_merged),
		uint32(diskstats.sectors_discarded),
		uint32(diskstats.spent_discarding_ms),
		uint32(diskstats.undefined1),
		uint32(diskstats.undefined2),
		uint32(diskstats.undefined3),
	}
}
