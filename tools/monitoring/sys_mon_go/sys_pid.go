// sys_processor.go

/*
	Copyright ©2020-21 H.F.M - system monitor library v1.0 https://github.com/hfmrow

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package sys_monitor

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	gltsss "github.com/hfmrow/gen_lib/tools/structures"
)

type Process struct {
	Pid int

	// Get percent for oneshot usage
	CpuPercentGet func() float64

	// Handle retrieving percent over time using goroutine
	// until CpuPercentSetAutoRefresh(false)
	CpuPercentSetAutoRefresh func(_ bool)
	cpuPercentAutoRefresh,
	cpuPercentStarted bool // Flag to indicate whether its already running
	CpuPercentAuto     float64
	cpuPercentAutoLock sync.RWMutex

	// To hold fields: 'number', 'name', 'description'
	statFileDescr []*itemStatFileDescr
	Desc          []*itemStatFileDescr
	// translation 'scanf' format from "C" to "Go"
	cToGoFormat *cToGoFormat

	// TODO debug purpose (to remove)
	h *gltsss.StructureHandler
}

func processNew(pid int) (*Process, error) {
	prc := new(Process)

	// init
	prc.Pid = pid
	// init stat file description
	err := prc.getStatFileDescription()
	if err != nil {
		return nil, err
	}

	prc.CpuPercentGet = prc.cpuPercent
	prc.CpuPercentSetAutoRefresh = prc.cpuPercentSetAutoRefresh
	// init format from "C" to "Go"
	prc.cToGoFormat = cToGoFormatNew()

	// TODO debug purpose (for saving structure)
	// prc.h = gltsss.StructureHandlerNew()
	// prc.wr()

	return prc, nil
}
func (prc *Process) cpuPercentSetAutoRefresh(enabled bool) {

	if !prc.cpuPercentAutoRefresh {
		prc.cpuPercentAutoRefresh = enabled
		go func() {
			for prc.cpuPercentAutoRefresh {
				prc.cpuPercentAutoLock.Lock()
				prc.CpuPercentAuto = CpuPercentPid(1, prc.Pid)
				prc.cpuPercentAutoLock.Unlock()
				time.Sleep(time.Second)
			}
			fmt.Printf("CpuPercentAuto: %v\n", prc.cpuPercentAutoRefresh)
		}()
	} else if !enabled {
		prc.cpuPercentAutoRefresh = false
	}
}

func (prc *Process) cpuPercent() float64 {

	return CpuPercentPid(1, prc.Pid)
}

func (prc *Process) GetStatus() (*StatsPid, error) {

	var (
		filename = regFilePidRepl.ReplaceAllString(smsLocal.LinuxPidStatus, fmt.Sprintf("%d", prc.Pid))
	)

	values, err := smsLocal.readValues(filename, ":")
	if err != nil {
		return nil, err
	}

	inStat := StatsPidNew()

	for idx, value := range values {
		spid := statPidNew()
		spid.Field = idx

		switch {
		case len(value) > 1:
			spid.rawVal = regWhtSpaces.ReplaceAllString(value[1], "\t")
			inStat.Data[value[0]] = spid
		case len(value) > 0:
			inStat.Data[value[0]] = spid
		}
	}
	return &inStat, nil
}

func (prc *Process) GetStat() (*StatsPid, error) {

	filename := regFilePidRepl.ReplaceAllString(smsLocal.LinuxPidStat, fmt.Sprintf("%d", prc.Pid))

	lines, err := smsLocal.readLines(filename)
	if err != nil {
		return nil, err
	}

	line := strings.Join(lines, "")
	commSl := regInBrckts.FindAllString(line, 1)                          // isolate 'name' (case of [[:space:]] inside parenthesis
	comm := strings.Join(commSl, "")[1 : len(commSl[0])-1]                // Removing parenthesis
	tmpLine := strings.Replace(line, comm, "", 1)                         // Remove it from line
	bItems := regWhtSpaces.Split(tmpLine, -1)                             // Split line on white chars
	items := append(bItems[:1], append([]string{comm}, bItems[2:]...)...) // re-insert 'name' at his right place while removing '()'

	inStat := StatsPidNew()

	// Seeking for the rest
	for field, iSFD := range prc.statFileDescr {
		spid := statPidNew()
		spid.rawVal = items[field]
		spid.Field = field
		spid.nType = iSFD.Format
		spid.Descr = iSFD.Desc

		inStat.Data[iSFD.Id] = spid
	}

	return &inStat, nil
}

type StatsPid struct {
	Data map[string]*statPid
	sync.RWMutex
}

// Create a map to hold retrieved values.
func StatsPidNew() StatsPid {
	return StatsPid{Data: make(map[string]*statPid)}
}

type statPid struct {
	Descr  string
	Field  int
	rawVal string
	nType  string
}

func (s *statPid) String() string {

	// TODO i think there is a delay between creation and request
	// this delay produce a 'runtime error: invalid memory address or nil pointer dereference'
	if s != nil {
		switch {
		case strings.Contains(s.rawVal, "kB"):
			val := regNum.FindString(s.rawVal)
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				log.Printf("sys_Pid/statPid/String/ParseInt [%s]: %v\n", s.rawVal, err)
			}
			return smsLocal.humanReadableSize(i * 1024)
		}
		return s.rawVal
	}
	return ""
}

// Get as int64 in bytes (not Kb). Return -1 for error.
func (s *statPid) Int64() int64 {

	// TODO i think there is a delay between creation and request
	// this delay produce a 'runtime error: invalid memory address or nil pointer dereference'
	if s != nil {
		switch {
		case strings.Contains(s.rawVal, "kB"):
			val := regNum.FindString(s.rawVal)
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				log.Printf("sys_Pid/statPid/Int64/ParseInt [%s]: %v\n", s.rawVal, err)
			}
			return i * 1024
		}
	}
	return -1
}

func statPidNew() *statPid {
	return new(statPid)
}

func (p *Process) GetProcess(pid int) {
	p.h.StructWrite(p, "saved.json")
}

func (p *Process) Write(filename string) {
	p.h.StructWrite(p, filename)
}

/*
 * stat file description builder
 */
var (
	regIdLine   = regexp.MustCompile(`(?m)(\(\d+\))\s+\w+\s+.+`)
	regStatDesc = regexp.MustCompile(`(?s)\s+/proc/\[pid\]/stat`)
)

type itemStatFileDescr struct {
	Field int
	Id,
	Format,
	Desc string
}

func itemStatFileDescrNew() *itemStatFileDescr {
	return new(itemStatFileDescr)
}

func (prc *Process) getStatFileDescription() error {
	var (
		isfd    *itemStatFileDescr
		getDesc = func(idx int, lines *[]string) (int, string) {

			var (
				tmpDesc []string
				lIdx    int
			)
			for lIdx = idx; lIdx < len(*lines); lIdx++ {
				line := strings.TrimSpace((*lines)[lIdx])

				switch {
				case len(line) == 0:
					continue

				case regIdLine.MatchString(line):
					return lIdx - 1, strings.Join(tmpDesc, " ")

				default:
					tmpDesc = append(tmpDesc, line)
				}
			}
			if len(tmpDesc) > 0 {
				return lIdx - 1, strings.Join(tmpDesc, " ")
			}
			return lIdx - 1, "Unavailable description."
		}
	)

	lines := strings.Split(descrStatStr, "\n")
	for lineIdx := 0; lineIdx < len(lines); lineIdx++ {

		line := strings.TrimSpace(lines[lineIdx])

		switch {
		case len(line) == 0:
			continue

		case regIdLine.MatchString(line):
			isfd = itemStatFileDescrNew()
			splitted := regWhtSpaces.Split(line, -1)
			if val, err := strconv.Atoi(regNum.FindString(splitted[0])); err != nil {
				return err
			} else {
				isfd.Field = val - 1
			}
			isfd.Id = splitted[1]
			isfd.Format = splitted[2]
			lineIdx, isfd.Desc = getDesc(lineIdx+1, &lines)
		}
		prc.statFileDescr = append(prc.statFileDescr, isfd)
		prc.Desc = append(prc.Desc, isfd)
	}
	return nil
}

/*
 * "C" scanf format conversion to Go
 */
type cToGoFormat map[string]*gFmt

// Create a map that hold scanf format conversion to Go
func cToGoFormatNew() *cToGoFormat {
	c := new(cToGoFormat)
	*c = make(map[string]*gFmt)
	gofmt := [][]string{
		{"%u", "%d", "u"},
		{"%d", "%d", "s"},
		{"%s", "%s", "na"},
		{"%c", "%s", "na"},
		{"%lu", "%d", "u"},
		{"%ld", "%d", "s"},
		{"%llu", "%d", "u"},
	}
	for _, v := range gofmt {
		gF := gFmtNew()
		gF.GoFmt = v[1]
		gF.Unsigned = v[2] == "u"
		(*c)[v[0]] = gF

	}
	return c
}

type gFmt struct {
	GoFmt    string
	Unsigned bool
}

func gFmtNew() *gFmt {
	return new(gFmt)
}

/*
 *	Fields description for (PID) 'stat' file.
 */

var descrStatStr string = `
(1) pid  %d
          The process ID.

(2) comm  %s
          The filename of the executable, in parentheses.
          This is visible whether or not the executable is
          swapped out.

(3) state  %c
          One of the following characters, indicating process
          state:

          R  Running

          S  Sleeping in an interruptible wait

          D  Waiting in uninterruptible disk sleep

          Z  Zombie

          T  Stopped (on a signal) or (before Linux 2.6.33)
             trace stopped

          t  Tracing stop (Linux 2.6.33 onward)

          W  Paging (only before Linux 2.6.0)

          X  Dead (from Linux 2.6.0 onward)

          x  Dead (Linux 2.6.33 to 3.13 only)

          K  Wakekill (Linux 2.6.33 to 3.13 only)

          W  Waking (Linux 2.6.33 to 3.13 only)

          P  Parked (Linux 3.9 to 3.13 only)

(4) ppid  %d
          The PID of the parent of this process.

(5) pgrp  %d
          The process group ID of the process.

(6) session  %d
          The session ID of the process.

(7) tty_nr  %d
          The controlling terminal of the process.  (The minor
          device number is contained in the combination of
          bits 31 to 20 and 7 to 0; the major device number is
          in bits 15 to 8.)

(8) tpgid  %d
          The ID of the foreground process group of the con‐
          trolling terminal of the process.

(9) flags  %u
          The kernel flags word of the process.  For bit means‐
          ings, see the PF_* defines in the Linux kernel
          source file include/linux/sched.h.  Details depend
          on the kernel version.

          The format for this field was %lu before Linux 2.6.

(10) minflt  %lu
          The number of minor faults the process has made
          which have not required loading a memory page from
          disk.

(11) cminflt  %lu
          The number of minor faults that the process's
          waited-for children have made.

(12) majflt  %lu
          The number of major faults the process has made
          which have required loading a memory page from disk.

(13) cmajflt  %lu
          The number of major faults that the process's
          waited-for children have made.

(14) utime  %lu
          Amount of time that this process has been scheduled
          in user mode, measured in clock ticks (divide by
          sysconf(_SC_CLK_TCK)).  This includes guest time,
          guest_time (time spent running a virtual CPU, see
          below), so that applications that are not aware of
          the guest time field do not lose that time from
          their calculations.

(15) stime  %lu
          Amount of time that this process has been scheduled
          in kernel mode, measured in clock ticks (divide by
          sysconf(_SC_CLK_TCK)).

(16) cutime  %ld
          Amount of time that this process's waited-for chil‐
          dren have been scheduled in user mode, measured in
          clock ticks (divide by sysconf(_SC_CLK_TCK)).  (See
          also times(2).)  This includes guest time,
          cguest_time (time spent running a virtual CPU, see
          below).

(17) cstime  %ld
          Amount of time that this process's waited-for chil‐
          dren have been scheduled in kernel mode, measured in
          clock ticks (divide by sysconf(_SC_CLK_TCK)).

(18) priority  %ld
          (Explanation for Linux 2.6) For processes running a
          real-time scheduling policy (policy below; see
          sched_setscheduler(2)), this is the negated schedul‐
          ing priority, minus one; that is, a number in the
          range -2 to -100, corresponding to real-time priori‐
          ties 1 to 99.  For processes running under a non-
          real-time scheduling policy, this is the raw nice
          value (setpriority(2)) as represented in the kernel.
          The kernel stores nice values as numbers in the
          range 0 (high) to 39 (low), corresponding to the
          user-visible nice range of -20 to 19.

          Before Linux 2.6, this was a scaled value based on
          the scheduler weighting given to this process.

(19) nice  %ld
          The nice value (see setpriority(2)), a value in the
          range 19 (low priority) to -20 (high priority).

(20) num_threads  %ld
          Number of threads in this process (since Linux 2.6).
          Before kernel 2.6, this field was hard coded to 0 as
          a placeholder for an earlier removed field.

(21) itrealvalue  %ld
          The time in jiffies before the next SIGALRM is sent
          to the process due to an interval timer.  Since ker‐
          nel 2.6.17, this field is no longer maintained, and
          is hard coded as 0.

(22) starttime  %llu
          The time the process started after system boot.  In
          kernels before Linux 2.6, this value was expressed
          in jiffies.  Since Linux 2.6, the value is expressed
          in clock ticks (divide by sysconf(_SC_CLK_TCK)).

          The format for this field was %lu before Linux 2.6.

(23) vsize  %lu
          Virtual memory size in bytes.

(24) rss  %ld
          Resident Set Size: number of pages the process has
          in real memory.  This is just the pages which count
          toward text, data, or stack space.  This does not
          include pages which have not been demand-loaded in,
          or which are swapped out.

(25) rsslim  %lu
          Current soft limit in bytes on the rss of the
          process; see the description of RLIMIT_RSS in
          getrlimit(2).

(26) startcode  %lu  [PT]
          The address above which program text can run.

(27) endcode  %lu  [PT]
          The address below which program text can run.

(28) startstack  %lu  [PT]
          The address of the start (i.e., bottom) of the
          stack.

(29) kstkesp  %lu  [PT]
          The current value of ESP (stack pointer), as found
          in the kernel stack page for the process.

(30) kstkeip  %lu  [PT]
          The current EIP (instruction pointer).

(31) signal  %lu
          The bitmap of pending signals, displayed as a deci‐
          mal number.  Obsolete, because it does not provide
          information on real-time signals; use
          /proc/[pid]/status instead.

(32) blocked  %lu
          The bitmap of blocked signals, displayed as a deci‐
          mal number.  Obsolete, because it does not provide
          information on real-time signals; use
          /proc/[pid]/status instead.

(33) sigignore  %lu
          The bitmap of ignored signals, displayed as a deci‐
          mal number.  Obsolete, because it does not provide
          information on real-time signals; use
          /proc/[pid]/status instead.

(34) sigcatch  %lu
          The bitmap of caught signals, displayed as a decimal
          number.  Obsolete, because it does not provide
          information on real-time signals; use
          /proc/[pid]/status instead.

(35) wchan  %lu  [PT]
          This is the "channel" in which the process is wait‐
          ing.  It is the address of a location in the kernel
          where the process is sleeping.  The corresponding
          symbolic name can be found in /proc/[pid]/wchan.

(36) nswap  %lu
          Number of pages swapped (not maintained).

(37) cnswap  %lu
          Cumulative nswap for child processes (not main‐
          tained).

(38) exit_signal  %d  (since Linux 2.1.22)
          Signal to be sent to parent when we die.

(39) processor  %d  (since Linux 2.2.8)
          CPU number last executed on.

(40) rt_priority  %u  (since Linux 2.5.19)
          Real-time scheduling priority, a number in the range
          1 to 99 for processes scheduled under a real-time
          policy, or 0, for non-real-time processes (see
          sched_setscheduler(2)).

(41) policy  %u  (since Linux 2.5.19)
          Scheduling policy (see sched_setscheduler(2)).
          Decode using the SCHED_* constants in linux/sched.h.

          The format for this field was %lu before Linux
          2.6.22.

(42) delayacct_blkio_ticks  %llu  (since Linux 2.6.18)
          Aggregated block I/O delays, measured in clock ticks
          (centiseconds).

(43) guest_time  %lu  (since Linux 2.6.24)
          Guest time of the process (time spent running a vir‐
          tual CPU for a guest operating system), measured in
          clock ticks (divide by sysconf(_SC_CLK_TCK)).

(44) cguest_time  %ld  (since Linux 2.6.24)
          Guest time of the process's children, measured in
          clock ticks (divide by sysconf(_SC_CLK_TCK)).

(45) start_data  %lu  (since Linux 3.3)  [PT]
          Address above which program initialized and unini‐
          tialized (BSS) data are placed.

(46) end_data  %lu  (since Linux 3.3)  [PT]
          Address below which program initialized and unini‐
          tialized (BSS) data are placed.

(47) start_brk  %lu  (since Linux 3.3)  [PT]
          Address above which program heap can be expanded
          with brk(2).

(48) arg_start  %lu  (since Linux 3.5)  [PT]
          Address above which program command-line arguments
          (argv) are placed.

(49) arg_end  %lu  (since Linux 3.5)  [PT]
          Address below program command-line arguments (argv)
          are placed.

(50) env_start  %lu  (since Linux 3.5)  [PT]
          Address above which program environment is placed.

(51) env_end  %lu  (since Linux 3.5)  [PT]
          Address below which program environment is placed.

(52) exit_code  %d  (since Linux 3.5)  [PT]
          The thread's exit status in the form reported by
          waitpid(2).`

/*
 * Unused
 */
var descrStatusStr = `
Name	Command run by this process.  Strings longer than TASK_COMM_LEN (16) characters (including the terminating null byte) are silently truncated.
Umask	Process umask, expressed in octal with a leading zero; see umask(2).  (Since Linux 4.7.)
State	Current state of the process.  One of "R (running)", "S (sleeping)", "D (disk sleep)", "T (stopped)", "t (tracing stop)", "Z (zombie)", or "X (dead)".
Tgid	Thread group ID (i.e., Process ID).
Ngid	NUMA group ID (0 if none; since Linux 3.13).
Pid	Thread ID (see gettid(2)).
PPid	PID of parent process.
TracerPid	PID of process tracing this process (0 if not being traced).
Uid	Real, effective, saved set, and filesystem UIDs.
Gid	Real, effective, saved set, and filesystem GIDs.
FDSize	Number of file descriptor slots currently allocated.
Groups	Supplementary group list.
NStgid	Thread group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member. The leftmost entry shows the value with respect to the PID namespace of the process that mounted this procfs (or the root namespace if mounted by the kernel), followed by the value in successively nested inner namespaces. (Since Linux 4.1.)
NSpid	Thread ID in each of the PID namespaces of which [pid] is a member.  The fields are ordered as for NStgid.  (Since Linux 4.1.)
NSpgid	Process group ID in each of the PID namespaces of which [pid] is a member. The fields are ordered as for NStgid. (Since Linux 4.1.)
NSsid	descendant namespace session ID hierarchy Session ID in each of the PID namespaces of which [pid] is a member. The fields are ordered as for NStgid. (Since Linux 4.1.)
VmPeak	Peak virtual memory size.
VmSize	Virtual memory size.
VmLck	Locked memory size (see mlock(2)).
VmPin	Pinned memory size (since Linux 3.2).  These are pages that can't be moved because something needs to directly access physical memory.
VmHWM	Peak resident set size ("high water mark"). This value is inaccurate; see /proc/[pid]/statm.
VmRSS	Resident set size.  Note that the value here is the sum of RssAnon, RssFile, and RssShmem. This value is inaccurate; see /proc/[pid]/statm.
RssAnon	Size of resident anonymous memory.  (since Linux 4.5).  This value is inaccurate; see /proc/[pid]/statm.
RssFile	Size of resident file mappings.  (since Linux 4.5). This value is inaccurate; see /proc/[pid]/statm.
RssShmem	Size of resident shared memory (includes System V shared memory, mappings from tmpfs(5), and shared anonymous mappings). (since Linux 4.5).
VmData	Size of data segments. This value is inaccurate; see /proc/[pid]/statm.
VmStk	Size of stack segments. This value is inaccurate; see /proc/[pid]/statm.
VmExe	Size of text segments. This value is inaccurate; see /proc/[pid]/statm.
VmLib	Shared library code size.
VmPTE	Page table entries size (since Linux 2.6.10).
VmPMD	Size of second-level page tables (added in Linux 4.0; removed in Linux 4.15).
VmSwap	Swapped-out virtual memory size by anonymous private pages; shmem swap usage is not included (since Linux 2.6.34). This value is inaccurate; see /proc/[pid]/statm.
HugetlbPages	Size of hugetlb memory portions (since Linux 4.4).
CoreDumping	Contains the value 1 if the process is currently dumping core, and 0 if it is not (since Linux 4.15). This information can be used by a monitoring process to avoid killing a process that is currently dumping core, which could result in a corrupted core dump file.
Threads	Number of threads in process containing this thread.
SigQ	This field contains two slash-separated numbers that relate to queued signals for the real user ID of this process. The first of these is the number of currently queued signals for this real user ID, and the second is the resource limit on the number of queued signals for this process (see the description of RLIMIT_SIGPENDING in getrlimit(2)).
SigPnd	Mask (expressed in hexadecimal) of signals pending for thread as a whole (see pthreads(7)).
ShdPnd	Mask (expressed in hexadecimal) of signals pending for process as a whole (see signal(7)).
SigBlk	Masks (expressed in hexadecimal) indicating signals being blocked (see signal(7)).
SigIgn	Masks (expressed in hexadecimal) indicating signals being ignored (see signal(7)).
SigCgt	Masks (expressed in hexadecimal) indicating signals being caught (see signal(7)).
CapInh	Masks (expressed in hexadecimal) of capabilities enabled in inheritable set (see capabilities(7)).
CapPrm	Masks (expressed in hexadecimal) of capabilities enabled in permitted set (see capabilities(7)).
CapEff	Masks (expressed in hexadecimal) of capabilities enabled in effective set (see capabilities(7)).
CapBnd	Capability bounding set, expressed in hexadecimal (since Linux 2.6.26, see capabilities(7)).
CapAmb	Ambient capability set, expressed in hexadecimal (since Linux 4.3, see capabilities(7)).
NoNewPrivs	Value of the no_new_privs bit (since Linux 4.10, see prctl(2)).
Seccomp	Seccomp mode of the process (since Linux 3.8, see seccomp(2)). 0 means SECCOMP_MODE_DISABLED; 1 means SECCOMP_MODE_STRICT; 2 means SECCOMP_MODE_FILTER. This field is provided only if the kernel was built with the CONFIG_SECCOMP kernel configuration option enabled.
Speculation_Store_Bypass	Speculation flaw mitigation state (since Linux 4.17, see prctl(2)).
Cpus_allowed	Hexadecimal mask of CPUs on which this process may run (since Linux 2.6.24, see cpuset(7)).
Cpus_allowed_list	Same as previous, but in "list format" (since Linux 2.6.26, see cpuset(7)).
Mems_allowed	Mask of memory nodes allowed to this process (since Linux 2.6.24, see cpuset(7)).
Mems_allowed_list	Same as previous, but in "list format" (since Linux 2.6.26, see cpuset(7)).
voluntary_ctxt_switches	Number of voluntary context switches (since Linux 2.6.23).
nonvoluntary_ctxt_switches	Number of involuntary context switches (since Linux 2.6.23).`

// import (
// 	"errors"
// 	"io/ioutil"
// 	"math"
// 	"os/exec"
// 	"path"
// 	"runtime"
// 	"strconv"
// 	"strings"
// 	"sync"
// )

// // SysInfo will record cpu and memory data
// type SysInfo struct {
// 	CPU    float64
// 	Memory float64
// 	Rss    float64
// }

// // Stat will store CPU time struct
// type Stat struct {
// 	utime  float64
// 	stime  float64
// 	cutime float64
// 	cstime float64
// 	start  float64
// 	rss    float64
// 	uptime float64
// }

// type fn func(int) (*SysInfo, error)

// var fnMap map[string]fn
// var platform string
// var history map[int]Stat
// var historyLock sync.Mutex
// var eol string

// func wrapper(statType string) func(pid int) (*SysInfo, error) {
// 	return func(pid int) (*SysInfo, error) {
// 		return stat(pid, statType)
// 	}
// }
// func init() {
// 	platform = runtime.GOOS
// 	if eol = "\n"; strings.Index(platform, "win") == 0 {
// 		platform = "win"
// 		eol = "\r\n"
// 	}
// 	history = make(map[int]Stat)
// 	fnMap = make(map[string]fn)
// 	fnMap["darwin"] = wrapper("ps")
// 	fnMap["sunos"] = wrapper("ps")
// 	fnMap["freebsd"] = wrapper("ps")
// 	fnMap["aix"] = wrapper("ps")
// 	fnMap["linux"] = wrapper("proc")
// 	fnMap["netbsd"] = wrapper("proc")
// 	fnMap["win"] = wrapper("win")
// }
// func formatStdOut(stdout []byte, userfulIndex int) []string {
// 	infoArr := strings.Split(string(stdout), eol)[userfulIndex]
// 	ret := strings.Fields(infoArr)
// 	return ret
// }

// func parseFloat(val string) float64 {
// 	floatVal, _ := strconv.ParseFloat(val, 64)
// 	return floatVal
// }

// func stat(pid int, statType string) (*SysInfo, error) {
// 	sysInfo := &SysInfo{}
// 	_history := history[pid]
// 	if statType == "ps" {
// 		args := "-o pcpu,rss -p"
// 		if platform == "aix" {
// 			args = "-o pcpu,rssize -p"
// 		}
// 		stdout, _ := exec.Command("ps", args, strconv.Itoa(pid)).Output()
// 		ret := formatStdOut(stdout, 1)
// 		if len(ret) == 0 {
// 			return sysInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
// 		}
// 		sysInfo.CPU = parseFloat(ret[0])
// 		sysInfo.Memory = parseFloat(ret[1]) * 1024
// 	} else if statType == "proc" {
// 		// default clkTck and pageSize
// 		var clkTck float64 = 100
// 		var pageSize float64 = 4096

// 		uptimeFileBytes, err := ioutil.ReadFile(path.Join("/proc", "uptime"))
// 		uptime := parseFloat(strings.Split(string(uptimeFileBytes), " ")[0])

// 		clkTckStdout, err := exec.Command("getconf", "CLK_TCK").Output()
// 		if err == nil {
// 			clkTck = parseFloat(formatStdOut(clkTckStdout, 0)[0])
// 		}

// 		pageSizeStdout, err := exec.Command("getconf", "PAGESIZE").Output()
// 		if err == nil {
// 			pageSize = parseFloat(formatStdOut(pageSizeStdout, 0)[0])
// 		}

// 		procStatFileBytes, err := ioutil.ReadFile(path.Join("/proc", strconv.Itoa(pid), "stat"))
// 		splitAfter := strings.SplitAfter(string(procStatFileBytes), ")")

// 		if len(splitAfter) == 0 || len(splitAfter) == 1 {
// 			return sysInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
// 		}
// 		infos := strings.Split(splitAfter[1], " ")
// 		stat := &Stat{
// 			utime:  parseFloat(infos[12]),
// 			stime:  parseFloat(infos[13]),
// 			cutime: parseFloat(infos[14]),
// 			cstime: parseFloat(infos[15]),
// 			start:  parseFloat(infos[20]) / clkTck,
// 			rss:    parseFloat(infos[22]),
// 			uptime: uptime,
// 		}

// 		_stime := 0.0
// 		_utime := 0.0
// 		if _history.stime != 0 {
// 			_stime = _history.stime
// 		}

// 		if _history.utime != 0 {
// 			_utime = _history.utime
// 		}
// 		total := stat.stime - _stime + stat.utime - _utime
// 		total = total / clkTck

// 		seconds := stat.start - uptime
// 		if _history.uptime != 0 {
// 			seconds = uptime - _history.uptime
// 		}

// 		seconds = math.Abs(seconds)
// 		if seconds == 0 {
// 			seconds = 1
// 		}

// 		historyLock.Lock()
// 		history[pid] = *stat
// 		historyLock.Unlock()
// 		sysInfo.CPU = (total / seconds) * 100
// 		sysInfo.Memory = stat.rss * pageSize
// 		sysInfo.Rss = stat.rss
// 	}
// 	return sysInfo, nil

// }

// // GetStat will return current system CPU and memory data
// func GetStat(pid int) (*SysInfo, error) {
// 	sysInfo, err := fnMap[platform](pid)
// 	return sysInfo, err
// }
