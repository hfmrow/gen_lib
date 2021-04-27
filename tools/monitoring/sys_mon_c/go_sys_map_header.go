// go_sys_map_header.go

package sys_mon

// #include "sys_smaps.h"
// #include "stdlib.h"
import "C"

type MapHeader struct {
	map_header *C.map_header
	Start      uint32
	End        uint32
	Flags      string
	Offset     uint64
	DevMaj     uint32
	DevMin     uint32
	Inode      uint32
	Pathname   string
}

func (v *MapHeader) native() *C.map_header {
	if v == nil || v.map_header == nil {
		return nil
	}
	return v.map_header
}

func wrapMapHeader(map_header *C.map_header) *MapHeader {
	return &MapHeader{
		map_header,
		uint32(map_header.start),
		uint32(map_header.end),
		C.GoString(&map_header.flags[0]),
		uint64(map_header.offset),
		uint32(map_header.dev_maj),
		uint32(map_header.dev_min),
		uint32(map_header.inode),
		C.GoString(&map_header.pathname[0]),
	}
}

func (v *MapHeader) free() {
	C.map_header_free(v.native())
}

func mapHeaderNew() *MapHeader {
	c := C.map_header_new()
	return wrapMapHeader(c)
}

func (v *MapHeader) ToString() string {
	cstr := new(C.char)
	C.map_header_to_string(cstr, v.native())
	str := C.GoString(cstr)
	return str
}
