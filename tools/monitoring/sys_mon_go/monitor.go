// sys_monitor.go

/*
	Copyright ©2020-21 hfmrow - system monitor v1.0 github.com/hfmrow/gen_lib/tools/monitoring/sys_mon_go

	This library is designed to display system information from some system files (debian file system)
	This package version, use pure Go code (except for 'cpu_usage.go' wich use a bit of "C") to retrieve
	information from these files.

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package sys_monitor

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	gltsss "github.com/hfmrow/gen_lib/tools/structures"
)

var (

	// Used regular expressions
	regFilePidRepl = regexp.MustCompile(`(\{pid\})`)
	regWhtSpaces   = regexp.MustCompile(`\s+`)
	regNum         = regexp.MustCompile(`(\d+)`)
	regInBrckts    = regexp.MustCompile(`(\(.+\))`)
	regSpcBrckts   = regexp.MustCompile(`(?m)(\(.+?\))|\s`)

	// For methods using in unexported functions
	smsLocal *SystemMonitorStruct

	// TODO used for debuging
	sh = gltsss.StructureHandlerNew()
)

func (sms *SystemMonitorStruct) Write(filename string) error {
	return sh.StructWrite(sms, filename)
}

type SystemMonitorStruct struct {
	Temp    []*device
	CPUs    *cpus
	Memory  *memory
	Swaps   *swaps
	Process *Process

	LinuxMeminfo, // memory information.
	LinuxSwaps, // swap file information.
	LinuxHwmonDir, // temperature information.
	LinuxCpufreq, // frequency information.
	LinuxUptime, // uttime information.
	LinuxPidStat, // process/system information.
	LinuxPidStatus string // process/system information.
}

func SystemMonitorStructNew() *SystemMonitorStruct {

	smsLocal = new(SystemMonitorStruct)
	smsLocal.LinuxMeminfo = "/proc/meminfo"
	smsLocal.LinuxHwmonDir = "/sys/class/hwmon"
	smsLocal.LinuxCpufreq = "/sys/devices/system/cpu/cpufreq"
	smsLocal.LinuxUptime = "/proc/uptime"
	smsLocal.LinuxSwaps = "/proc/swaps"
	smsLocal.LinuxPidStat = "/proc/{pid}/stat"
	smsLocal.LinuxPidStatus = "/proc/{pid}/status"
	return smsLocal
}

//  GetMemory:
func (sms *SystemMonitorStruct) GetMemory() error {
	sms.Memory = memoryNew()
	return sms.Memory.getMemory()
}

//  GetSwaps:
func (sms *SystemMonitorStruct) GetSwaps() error {
	sms.Swaps = swapsNew()
	return sms.Swaps.getSwaps()
}

// GetProcessors:
func (sms *SystemMonitorStruct) GetProcessors() error {
	sms.CPUs = cpusNew()
	return sms.CPUs.initProcessors()
}

// GetProcess:
func (sms *SystemMonitorStruct) InitProcess(pid int) (err error) {
	sms.Process, err = processNew(pid)
	return
}

// sortFilesOnly: remove directories, symlink and all that does not match 'matchWith'
func (sms *SystemMonitorStruct) sortFilesOnly(files []os.FileInfo, root string, matchWith ...string) (filenames []string) {

	var match = func(name string) bool {
		for _, pattern := range matchWith {
			if ok, _ := filepath.Match(pattern, name); ok {
				return true
			}
		}
		return false
	}

	for _, file := range files {
		if !file.IsDir() && file.Mode()&os.ModeSymlink != os.ModeSymlink && match(file.Name()) {

			filenames = append(filenames, filepath.Join(root, file.Name()))
		}
	}
	return
}

// getValue: retrieve a single value from a file. It can be a single value (field = -1)
// or a value of the "field" corresponding to "valName" where the separator is "sep".
func (sms *SystemMonitorStruct) getValue(filename, valName, sep string, field int) (out string, err error) {

	lines, err := sms.readLines(filename)
	if err != nil {
		return "", err
	}
	if field == -1 {
		return strings.TrimSpace(lines[0]), nil
	}
	for _, line := range lines {
		if strings.Contains(line, valName) {
			items := strings.Split(line, sep)
			if len(items) >= field {
				return strings.TrimSpace(items[field]), nil
			} else {
				return "", fmt.Errorf("Fields count < %d", field)
			}
		}
	}
	return "", fmt.Errorf("'%s', Not found", valName)
}

// readLines: getting lines from file
func (sms *SystemMonitorStruct) readLines(filename string) ([]string, error) {

	data, err := ioutil.ReadFile(filename)
	return strings.Split(string(data), "\n"), err
}

// readValues: getting values formatted like "name : value"
func (sms *SystemMonitorStruct) readValues(filename, sep string) ([][]string, error) {

	var (
		ret       [][]string
		fileItems []string
	)

	lines, err := sms.readLines(filename)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) > 0 {
			if fileItems = strings.Split(line, sep); len(fileItems) > 0 {
				var tmpValue []string
				for _, v := range fileItems {
					tmpValue = append(tmpValue, strings.TrimSpace(v))
				}
				ret = append(ret, tmpValue)
			}
		}
	}
	if len(ret) == 0 {
		return ret, fmt.Errorf("File is empty: %s", filename)
	}
	return ret, nil
}

// humanReadableSize: Convert file size (octets/bytes) to human readable
// version. Note: 'options' in order means => 'useDecimal', 'hideUnit'.
// 'useDecimal' argument define kilo = 1000 instead of 1024.
func (sms *SystemMonitorStruct) humanReadableSize(size interface{}, options ...bool) string {

	var (
		hideUnit               bool
		val                    float64
		kilo                   float64 = 1024
		sP, sT, sG, sM, sK, sb string  = "PiB", "TiB", "GiB", "MiB", "KiB", "b"
	)
	if len(options) > 0 && options[0] {
		kilo = 1000
		sP, sT, sG, sM, sK, sb = "PB", "TB", "GB", "MB", "kB", "b"
	}
	if len(options) > 1 && options[1] {
		hideUnit = true
	}
	switch v := size.(type) {
	case uint64:
		val = float64(v)
	case uint32:
		val = float64(v)
	case uint:
		val = float64(v)
	case int64:
		val = float64(v)
	case int32:
		val = float64(v)
	case int:
		val = float64(v)
	case float64:
		val = float64(v)
	case float32:
		val = float64(v)
	default:
		log.Printf("Unable to define type of: %v\n", size)
	}
	unit := sb
	switch {
	case val < kilo:
		val = val
		return fmt.Sprintf("%.0f%s", val, unit)
	case val < math.Pow(kilo, 2):
		val = val / kilo
		unit = sK
	case val < math.Pow(kilo, 3):
		val = val / math.Pow(kilo, 2)
		unit = sM
	case val < math.Pow(kilo, 4):
		val = val / math.Pow(kilo, 3)
		unit = sG
	case val < math.Pow(kilo, 5):
		val = val / math.Pow(kilo, 4)
		unit = sT
	case val < math.Pow(kilo, 6):
		val = val / math.Pow(kilo, 5)
		unit = sP
	}
	if hideUnit {
		unit = ""
	}
	return fmt.Sprintf("%2.2f%s", val, unit)
}

// String: Convert frenquency value int64 to string human readable version.
func (sms *SystemMonitorStruct) humanReadableFreq(value interface{}, hideUnit ...bool) string {

	var (
		val  float64
		unit string
	)

	switch v := value.(type) {
	case int64:
		val = float64(v)
	case int32:
		val = float64(v)
	case int:
		val = float64(v)
	default:
		val = value.(float64)
	}

	switch {
	case val < 1000:
		if len(hideUnit) > 0 && hideUnit[0] {
			unit = ""
		} else {
			unit = "hz"
		}
		fmt.Sprintf("%2.0f%s", val, unit)
	case val < 1000e3:
		val = val / 1000
		unit = "Mhz"
	case val < float64(1000e6):
		val = val / 1000e3
		unit = "Ghz"
	case val < float64(1000e9):
		val = val / 1000e6
		unit = "Thz"
	}

	if len(hideUnit) > 0 && hideUnit[0] {
		unit = ""
	}
	return fmt.Sprintf("%2.2f%s", val, unit)
}

// humanReadableTime: Convert time.Duration to string human readable version.
func (sms *SystemMonitorStruct) humanReadableTime(duration time.Duration, unitSep ...string) string {

	var (
		out,
		sep string
		stop bool
		min  = time.Second * 60
		hour = min * 60
		day  = hour * 24
		d    int64
	)

	if len(unitSep) > 0 {
		sep = unitSep[0]
	}

	for !stop {
		switch {

		case duration >= day: // day
			d = int64(duration / day)
			out += fmt.Sprintf("%dd%s", d, sep)
			stop = duration == day
			duration = duration - (time.Duration(d) * (day))

		case duration >= hour: // hour
			d = int64(duration / hour)
			out += fmt.Sprintf("%dh%s", d, sep)
			stop = duration == hour
			duration = duration - (time.Duration(d) * (hour))

		case duration >= min: // minute
			d = int64(duration / min)
			out += fmt.Sprintf("%dm%s", d, sep)
			stop = duration == min
			duration = duration - (time.Duration(d) * (min))

		case int64(duration) > 0: // seconds
			out += fmt.Sprintf("%.3fs%s", duration.Seconds(), sep)
			stop = true

		default:
			stop = true
		}
	}
	return out
}
