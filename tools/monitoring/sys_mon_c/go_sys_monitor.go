// go_sys_monitor.go

package sys_mon

// #include "file_func.h"
import "C"
import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	// gltsmgsncc "github.com/hfmrow/gen_lib/tools/monitoring/sys_mon/C"
)

// func CCharLen(str string, length int) *C.char {
// 	if len(str) > length {
// 		_, file, line, _ := runtime.Caller(1)
// 		log.Printf("./%s:%d Warning!, given string (%d) exceed defined length (%d)\n", filepath.Base(file), line, len(str), length)
// 		return nil
// 	}
// 	const cConst = length
// 	cstr := (*C.char)(C.CString(str))
// 	defer C.free(unsafe.Pointer(cstr))
// 	return *(*[cConst]C.char)(unsafe.Pointer(cstr))
// }

// Clean up the cache used by the package, on the next launch,
// the source code 'C' will be forced to be rebuilt. typically
// used in the development process for debugging purposes.
// func ForceRebuildPackage(rebuild ...bool) bool {
// 	if len(rebuild) > 0 && !rebuild[0] {
// 		return false
// 	}
// 	_, b, _, _ := runtime.Caller(0)
// 	outTerm, err := execCommand("clean cache", "go", "clean", "-cache", filepath.Dir(b))
// 	if err != nil {
// 		fmt.Printf("\nError: %s\n", outTerm)
// 	}
// 	fmt.Printf("\nCleaning cache of the package: %s\nFresh build will be available on next launch!\n\n",
// 		filepath.Base(filepath.Dir(b)))
// 	return true
// }

/*
 * Embedded 'C' ERROR function
 */
// Used to include 'C' directory and its content when using vendoring system.
// func fakeCall() {
// 	gltsmgsncc.AddSourceC()
// }

func getErrorString() error {
	cstr := new(C.char)
	C.internal_error_get(cstr)
	str := C.GoString(cstr)
	return fmt.Errorf("%s", str)
}

// execCommand: launch commandline application with arguments
// return terminal output and error.
func execCommand(infos string, cmds ...string) (outTerm []byte, err error) {

	execCmd := exec.Command(cmds[0], cmds[1:]...)
	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	outTerm, err = execCmd.CombinedOutput()
	if err != nil {
		err = errors.New(
			fmt.Sprintf("[%s][%s]\nCommand failed: %v\nTerminal:\n%v",
				infos,
				strings.Join(cmds, " "),
				err,
				string(outTerm)))
		return
	}
	return
}
