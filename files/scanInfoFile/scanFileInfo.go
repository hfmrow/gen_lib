// scanFileInfo.go

// Retrieve information about files

package scanFileInfo

import (
	"os"
	"path/filepath"
	"regexp"
	"time"

	humanize "github.com/dustin/go-humanize"
	times "gopkg.in/djherbis/times.v1"
)

// TODO rewrite it with independant result to minimize calculation time

// Used in ScanFiles function to store file information.
type fileInfos struct {
	IsExists         bool
	PathBase         string
	Base             string
	Path             string
	Ext              string
	Mtime            time.Time
	Atime            time.Time
	MtimeYMDhm       string
	AtimeYMDhm       string
	MtimeYMDhms      string
	AtimeYMDhms      string
	MtimeYMDhmsShort string
	AtimeYMDhmsShort string
	MtimeFriendlyHR  string
	AtimeFriendlyHR  string
	Type             string
	Size             int64
	SizeHR           string
}

// ScanFiles: Scan given files and retreive information about them stored in a []fileInfos structure.
func ScanFiles(inFiles []string) (outFiles []fileInfos) {
	for _, file := range inFiles {
		outFiles = append(outFiles, ScanFile(file))
	}
	return outFiles
}

// ScanFile: Scan a file and retreive information about it stored in a fileInfos structure.
func ScanFile(file string) (fi fileInfos) {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`)
	var tmpStr string
	fi.PathBase = file

	if infos, err := os.Stat(file); os.IsNotExist(err) {
		fi.IsExists = false
		return fi
	} else {
		fi.IsExists = true
		switch {
		case (infos.Mode()&os.ModeDir != 0):
			fi.Type = "Dir"
		case (infos.Mode()&os.ModeSymlink != 0):
			fi.Type = "Link"
		case (infos.Mode()&os.ModeAppend != 0):
			fi.Type = "Append" // a: append-only
		case (infos.Mode()&os.ModeExclusive != 0):
			fi.Type = "Exclusive" // l: exclusive use
		case (infos.Mode()&os.ModeTemporary != 0):
			fi.Type = "Temp" // T: temporary file; Plan 9 only
		case (infos.Mode()&os.ModeDevice != 0):
			fi.Type = "Device" // D: device file
		case (infos.Mode()&os.ModeNamedPipe != 0):
			fi.Type = "Pipe" // p: named pipe (FIFO)
		case (infos.Mode()&os.ModeSocket != 0):
			fi.Type = "Socket" // S: Unix domain socket
		case (infos.Mode()&os.ModeCharDevice != 0):
			fi.Type = "CharDev" // c: Unix character device, when ModeDevice is set
		case (infos.Mode()&os.ModeSticky != 0):
			fi.Type = "Sticky" // t: sticky
		case (infos.Mode()&os.ModeIrregular != 0):
			fi.Type = "Unknown" // ?: non-regular file; nothing else is known about this file
		default:
			fi.Type = "File"
		}

		if fi.Type == "Dir" {
			fi.Base = filepath.Base(file)
			fi.Path = file
			fi.Ext = "Dir"
		} else {
			fi.Base = filepath.Base(file)
			fi.Path = filepath.Dir(file)
			fi.Ext = filepath.Ext(file)
		}

		fi.Size = infos.Size()
		fi.SizeHR = humanize.Bytes(uint64(fi.Size))

		fi.Atime = times.Get(infos).AccessTime()
		fi.Mtime = times.Get(infos).ModTime()
		fi.MtimeFriendlyHR = humanize.Time(fi.Mtime)
		fi.AtimeFriendlyHR = humanize.Time(fi.Atime)
		fi.MtimeYMDhm = fi.Mtime.String()[:16]
		fi.AtimeYMDhm = fi.Atime.String()[:16]
		fi.MtimeYMDhms = fi.Mtime.String()[:19]
		fi.AtimeYMDhms = fi.Atime.String()[:19]
		tmpStr = nonAlNum.ReplaceAllString(fi.MtimeYMDhms, "")
		fi.MtimeYMDhmsShort = tmpStr[2:len(tmpStr)]
		tmpStr = nonAlNum.ReplaceAllString(fi.AtimeYMDhms, "")
		fi.AtimeYMDhmsShort = tmpStr[2:len(tmpStr)]
	}
	return fi
}
