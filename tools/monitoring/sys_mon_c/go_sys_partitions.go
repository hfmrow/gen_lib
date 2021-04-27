// go_sys_partitions.go

package sys_mon

// #include "sys_partitions.h"
import "C"

type Partitions struct {
	partitions *C.partitions
	Details    []partition
}

// Close: Freeing 'C' structure.
func (s *Partitions) Close() {
	if s.partitions != nil {
		C.partitions_free(s.partitions)
	}
}

// Update: 'Partitions' structure.
func (s *Partitions) Update() error {
	c := C.partitions_update(s.partitions)
	if !bool(c) {
		return getErrorString()
	}
	s.wrap()
	return nil
}

// internal wraping main structure
func (s *Partitions) wrap() {
	count := int(s.partitions.count)
	s.Details = make([]partition, count)
	for i := 0; i < count; i++ {
		s.Details[i] = *wrapPartition(C.partitions_get_single(s.partitions, C.int(i)))
	}
}

// PartitionsNew: Create and initialise 'C' structure.
func PartitionsNew() (*Partitions, error) {
	s := new(Partitions)
	s.partitions = C.partitions_get()
	if s.partitions == nil {
		return nil, getErrorString()
	}
	s.wrap()
	return s, nil
}

type partition struct {
	partitions *C.partitions
	Major      int
	Minor      int
	Size       uint32
	Name       string
	ClassBlock *classBlock
}

func wrapPartition(partitions *C.partitions) *partition {
	if partitions == nil {
		return nil
	}

	return &partition{
		partitions,
		int(partitions.major),
		int(partitions.minor),
		uint32(partitions.blocks),
		C.GoString(&partitions.name[0]),
		wrapClassBlock(&partitions.class_block),
	}
}

type classBlock struct {
	class_block       *C.class_block
	HwSectorSize      uint32
	LogicalBlockSize  uint32
	MaxHwSectorsKb    uint32
	PhysicalBlockSize uint32
	ReadAheadKb       uint32
	WriteCache        uint32
	Removable         uint32
	Blocks            uint32
	DevType           string
	PartName          string
	Uuid              string
}

func wrapClassBlock(class_block *C.class_block) *classBlock {
	if class_block == nil {
		return nil
	}

	return &classBlock{
		class_block,
		uint32(class_block.hw_sector_size),
		uint32(class_block.logical_block_size),
		uint32(class_block.max_hw_sectors_kb),
		uint32(class_block.physical_block_size),
		uint32(class_block.read_ahead_kb),
		uint32(class_block.write_cache),
		uint32(class_block.removable),
		uint32(class_block.size),
		C.GoString(&class_block.dev_type[0]),
		C.GoString(&class_block.part_name[0]),
		C.GoString(&class_block.uuid[0]),
	}
}
