// go_sys_pid.go

package sys_mon

// #include "stdlib.h"
// #include "sys_pid_status.h"
import "C"
import (
	"log"
	"unsafe"
)

// PidInfosNew: Create and initialise 'C' structure.
// No need to 'free' (close) anything, everything is already handled.
func PidInfosNew() (*PidInfos, error) {
	c := wrapStoreFiles(C.get_pid_infos())
	if c == nil {
		return nil, getErrorString()
	}
	defer C.store_files_free(c.store_files)

	return c, nil
}

// Get pid using the filename base. Returns "-1" if not found.
func (sf *PidInfos) GetPidFromFilename(filename string) int {
	cstr := (*C.char)(C.CString(filename))
	defer C.free(unsafe.Pointer(cstr))
	pid := new(C.uint)
	c := C.get_pid_by_filename(cstr, pid)
	if c == nil {
		log.Printf(getErrorString().Error())
		return -1
	}
	return int(*pid)
}

// Get pid using the name. Note: Instead of the previous, this function
// is based on a 'comm' field which contains only 16 bytes, which means
// that if the name is greater than 16 characters, it will be truncated.
// Returns "-1" if not found.
func (sf *PidInfos) GetPidFromName(name string) int {
	cstr := (*C.char)(C.CString(name))
	defer C.free(unsafe.Pointer(cstr))
	pid := new(C.uint)
	c := C.get_pid_by_name(cstr, pid)
	if c == nil {
		log.Printf(getErrorString().Error())
		return -1
	}
	return int(*pid)
}

// PidInfos: structure to store PID information for all processes.
type PidInfos struct {
	store_files *C.store_files
	count       uint32
	Details     []storeFile
}

func wrapStoreFiles(store_files *C.store_files) *PidInfos {
	if store_files == nil {
		return nil
	}
	count := uint32(store_files.count)
	var Details = make([]storeFile, count)
	for idx := 0; idx < int(count); idx++ {
		Details[idx] = *wrapStoreFile(C.store_files_get_single(store_files, C.int(idx)))
	}

	return &PidInfos{
		store_files,
		count,
		Details,
	}
}

type storeFile struct {
	store_file *C.store_file
	Pid        uint32
	Ppid       uint32
	Uid        resfId
	Gid        resfId
	Name       string
	State      string
	Dirname    string
	Filename   string
}

func wrapStoreFile(store_file *C.store_file) *storeFile {
	return &storeFile{
		store_file,
		uint32(store_file.pid),
		uint32(store_file.ppid),
		*wrapResfId(&store_file.uid),
		*wrapResfId(&store_file.gid),
		C.GoString(&store_file.name[0]),
		C.GoString(&store_file.state[0]),
		C.GoString(&store_file.dirname[0]),
		C.GoString(&store_file.filename[0]),
	}
}
