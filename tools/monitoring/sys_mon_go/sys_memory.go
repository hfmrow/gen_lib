// sys_memory.go

/*
	Copyright ©2020 H.F.M - system monitor library v1.0 https://github.com/hfmrow

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package sys_monitor

import (
	"strconv"
)

/*
 * Available variables:
		"MemTotal", "MemFree", "MemAvailable", "Buffers", "Cached", "SwapCached", "Active", "Inactive", "Active(anon)",
		"Inactive(anon)", "Active(file)", "Inactive(file)", "Unevictable", "Mlocked", "SwapTotal", "SwapFree", "Dirty",
		"Writeback", "AnonPages", "Mapped", "Shmem", "KReclaimable", "Slab", "SReclaimable", "SUnreclaim", "KernelStack",
		"PageTables", "NFS_Unstable", "Bounce", "WritebackTmp", "CommitLimit", "Committed_AS", "VmallocTotal", "VmallocUsed",
		"VmallocChunk", "Percpu", "HardwareCorrupted", "AnonHugePages", "ShmemHugePages", "ShmemPmdMapped", "FileHugePages",
		"FilePmdMapped", "HugePages_Total", "HugePages_Free", "HugePages_Rsvd", "HugePages_Surp", "Hugepagesize", "Hugetlb",
		"DirectMap4k", "DirectMap2M", "DirectMap1G"
*/

// memory: structure that hold memory info values
type memory struct {
	Values            map[string]*valueMem
	CommonMemoryInfos []string
}

func memoryNew() *memory {
	m := new(memory)
	m.Values = make(map[string]*valueMem, 0)
	m.CommonMemoryInfos = []string{"MemTotal", "MemAvailable", "SwapTotal", "SwapFree"}
	return m
}

type valueMem struct {
	Value int64
}

func valueMemNew(value ...int64) *valueMem {
	v := new(valueMem)
	if len(value) > 0 {
		v.Value = value[0]
	}
	return v
}

func (v *valueMem) String() string {
	return smsLocal.humanReadableSize(v.Value)
}

// getMemory: Retrieve memory information.
// Note: Original value are given as KB, this function retrieve them as Bytes
func (m *memory) getMemory() error {

	values, err := smsLocal.readValues(smsLocal.LinuxMeminfo, ":")
	if err != nil {
		return err
	}

	for _, items := range values {

		vm := valueMemNew()
		length := len(items)
		switch {
		case length == 1:
			vm.Value = 0

		case length == 2:
			if val, err := strconv.ParseInt(regNum.FindString(items[1]), 10, 64); err != nil {
				return err
			} else {
				vm.Value = val * 1024 // Convert to bytes
			}
		}
		m.Values[items[0]] = vm
	}

	return nil
}

// GetCommonMemoryInfos: Scan memory information and return values for
// variables names contained in 'm.CommonMemoryInfos' string slice.
func (m *memory) GetCommonMemoryInfos() (map[string]*valueMem, error) {

	err := m.getMemory()
	if err != nil {
		return nil, err
	}

	out := make(map[string]*valueMem, 0)
	for _, name := range m.CommonMemoryInfos {
		out[name] = valueMemNew(m.Values[name].Value)
	}
	return out, nil
}
