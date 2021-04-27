# sys\_mon

```go
import "github.com/hfmrow/gen_lib/tools/monitoring/sys_mon"
```

## Index

- [func ForceRebuildPackage(rebuild ...bool) bool](<#func-forcerebuildpackage>)
- [type CpuFs](<#type-cpufs>)
  - [func CpuFsNew() (*CpuFs, error)](<#func-cpufsnew>)
  - [func (cf *CpuFs) Close()](<#func-cpufs-close>)
  - [func (cf *CpuFs) CurrFreqUpdate() error](<#func-cpufs-currfrequpdate>)
  - [func (cf *CpuFs) ListUpdate() error](<#func-cpufs-listupdate>)
- [type CpuPercentPid](<#type-cpupercentpid>)
  - [func CpuPercentPidNew(pid int) (*CpuPercentPid, error)](<#func-cpupercentpidnew>)
  - [func (s *CpuPercentPid) Close()](<#func-cpupercentpid-close>)
  - [func (v *CpuPercentPid) Update() error](<#func-cpupercentpid-update>)
- [type Diskstats](<#type-diskstats>)
  - [func DiskstatsNew() (*Diskstats, error)](<#func-diskstatsnew>)
  - [func (s *Diskstats) Close()](<#func-diskstats-close>)
  - [func (s *Diskstats) Update() error](<#func-diskstats-update>)
- [type MapHeader](<#type-mapheader>)
  - [func (v *MapHeader) ToString() string](<#func-mapheader-tostring>)
- [type Meminfo](<#type-meminfo>)
  - [func MeminfoNew() (*Meminfo, error)](<#func-meminfonew>)
  - [func (s *Meminfo) Close()](<#func-meminfo-close>)
  - [func (s *Meminfo) Update() (*Meminfo, error)](<#func-meminfo-update>)
- [type NanoMeasureMethod](<#type-nanomeasuremethod>)
- [type Partitions](<#type-partitions>)
  - [func PartitionsNew() (*Partitions, error)](<#func-partitionsnew>)
  - [func (s *Partitions) Close()](<#func-partitions-close>)
  - [func (s *Partitions) Update() error](<#func-partitions-update>)
- [type PidInfos](<#type-pidinfos>)
  - [func PidInfosNew() (*PidInfos, error)](<#func-pidinfosnew>)
  - [func (sf *PidInfos) GetPidFromFilename(filename string) int](<#func-pidinfos-getpidfromfilename>)
  - [func (sf *PidInfos) GetPidFromName(name string) int](<#func-pidinfos-getpidfromname>)
- [type ProcNetDev](<#type-procnetdev>)
  - [func ProcNetDevNew(pid ...uint32) (*ProcNetDev, error)](<#func-procnetdevnew>)
  - [func (s *ProcNetDev) Close()](<#func-procnetdev-close>)
  - [func (s *ProcNetDev) GetAvailableInterfaces() []string](<#func-procnetdev-getavailableinterfaces>)
  - [func (s *ProcNetDev) GetWanAdress(adress string, useStunSrv ...bool) (string, error)](<#func-procnetdev-getwanadress>)
  - [func (s *ProcNetDev) SetSuffix(suffix string) error](<#func-procnetdev-setsuffix>)
  - [func (s *ProcNetDev) SetUnit(unit string) error](<#func-procnetdev-setunit>)
  - [func (s *ProcNetDev) Update() error](<#func-procnetdev-update>)
- [type ProcPidStat](<#type-procpidstat>)
  - [func ProcPidStatNew(pid int) (*ProcPidStat, error)](<#func-procpidstatnew>)
  - [func (s *ProcPidStat) Close()](<#func-procpidstat-close>)
  - [func (s *ProcPidStat) Update() error](<#func-procpidstat-update>)
- [type Smaps](<#type-smaps>)
  - [func SmapsNew(pid int, maxReadEntries ...int) (*Smaps, error)](<#func-smapsnew>)
  - [func (s *Smaps) Close()](<#func-smaps-close>)
  - [func (s *Smaps) UpdateRollup() error](<#func-smaps-updaterollup>)
  - [func (s *Smaps) UpdateSmaps() error](<#func-smaps-updatesmaps>)
- [type StatusFile](<#type-statusfile>)
  - [func StatusFileNew(pid int) (*StatusFile, error)](<#func-statusfilenew>)
- [type SysTherm](<#type-systherm>)
  - [func SysThermNew() (*SysTherm, error)](<#func-systhermnew>)
  - [func (st *SysTherm) Close()](<#func-systherm-close>)
  - [func (st *SysTherm) Update() error](<#func-systherm-update>)
- [type TimeSpent](<#type-timespent>)
  - [func TimeSpentNew(method ...NanoMeasureMethod) (*TimeSpent, error)](<#func-timespentnew>)
  - [func (s *TimeSpent) Close()](<#func-timespent-close>)
  - [func (s *TimeSpent) MesurementMethodGet() int](<#func-timespent-mesurementmethodget>)
  - [func (s *TimeSpent) MesurementMethodSet(method NanoMeasureMethod)](<#func-timespent-mesurementmethodset>)
  - [func (s *TimeSpent) NanoCalculate() float64](<#func-timespent-nanocalculate>)
  - [func (s *TimeSpent) NanoGet()](<#func-timespent-nanoget>)
  - [func (s *TimeSpent) SpentGet() float64](<#func-timespent-spentget>)
  - [func (s *TimeSpent) TicksCalculate() float64](<#func-timespent-tickscalculate>)
  - [func (s *TimeSpent) TicksGet()](<#func-timespent-ticksget>)


## func ForceRebuildPackage

```go
func ForceRebuildPackage(rebuild ...bool) bool
```

Clean up the cache used by the package\, on the next launch\, the source code 'C' will be forced to be rebuilt\. typically used in the development process for debugging purposes\.

## type CpuFs

Structure to hold values retrieved from: '/sys/devices/system/cpu/cpufreq/policy\*' directories

```go
type CpuFs struct {
    CpuCount int
    CurrFreq []int64
    CpuList  []cpuFs
    // contains filtered or unexported fields
}
```

### func CpuFsNew

```go
func CpuFsNew() (*CpuFs, error)
```

### func \(\*CpuFs\) Close

```go
func (cf *CpuFs) Close()
```

### func \(\*CpuFs\) CurrFreqUpdate

```go
func (cf *CpuFs) CurrFreqUpdate() error
```

### func \(\*CpuFs\) ListUpdate

```go
func (cf *CpuFs) ListUpdate() error
```

## type CpuPercentPid

```go
type CpuPercentPid struct {
    CpuPercent float32
    MemoryRss  int64
    // contains filtered or unexported fields
}
```

### func CpuPercentPidNew

```go
func CpuPercentPidNew(pid int) (*CpuPercentPid, error)
```

CpuPercentPidNew: Create and initialise 'C' structure\.

### func \(\*CpuPercentPid\) Close

```go
func (s *CpuPercentPid) Close()
```

Close: Freeing 'C' structure\.

### func \(\*CpuPercentPid\) Update

```go
func (v *CpuPercentPid) Update() error
```

update current values

## type Diskstats

```go
type Diskstats struct {
    Details []gdiskstats
    // contains filtered or unexported fields
}
```

### func DiskstatsNew

```go
func DiskstatsNew() (*Diskstats, error)
```

DiskstatsNew: Create and initialise 'C' structure\.

### func \(\*Diskstats\) Close

```go
func (s *Diskstats) Close()
```

Close: Freeing 'C' structure\.

### func \(\*Diskstats\) Update

```go
func (s *Diskstats) Update() error
```

Update: 'Diskstats' structure\.

## type MapHeader

```go
type MapHeader struct {
    Start    uint32
    End      uint32
    Flags    string
    Offset   uint64
    DevMaj   uint32
    DevMin   uint32
    Inode    uint32
    Pathname string
    // contains filtered or unexported fields
}
```

### func \(\*MapHeader\) ToString

```go
func (v *MapHeader) ToString() string
```

## type Meminfo

```go
type Meminfo struct {
    MemTotal          uint32
    MemFree           uint32
    MemAvailable      uint32
    Buffers           uint32
    Cached            uint32
    SwapCached        uint32
    Active            uint32
    Inactive          uint32
    ActiveAnon        uint32
    InactiveAnon      uint32
    ActiveFile        uint32
    InactiveFile      uint32
    Unevictable       uint32
    Mlocked           uint32
    SwapTotal         uint32
    SwapFree          uint32
    Dirty             uint32
    Writeback         uint32
    AnonPages         uint32
    Mapped            uint32
    Shmem             uint32
    Kreclaimable      uint32
    Slab              uint32
    Sreclaimable      uint32
    Sunreclaim        uint32
    KernelStack       uint32
    PageTables        uint32
    NfsUnstable       uint32
    Bounce            uint32
    WritebackTmp      uint32
    CommitLimit       uint32
    CommittedAs       uint32
    VmallocTotal      uint32
    VmallocUsed       uint32
    VmallocChunk      uint32
    Percpu            uint32
    HardwareCorrupted uint32
    AnonHugePages     uint32
    ShmemHugePages    uint32
    ShmemPmdMapped    uint32
    FileHugePages     uint32
    FilePmdMapped     uint32
    HugePagesTotal    uint32
    HugePagesFree     uint32
    HugePagesRsvd     uint32
    HugePagesSurp     uint32
    Hugepagesize      uint32
    Hugetlb           uint32
    DirectMap4k       uint32
    DirectMap2M       uint32
    DirectMap1G       uint32
    // contains filtered or unexported fields
}
```

### func MeminfoNew

```go
func MeminfoNew() (*Meminfo, error)
```

MeminfoNew: Create and initialise 'C' structure\.

### func \(\*Meminfo\) Close

```go
func (s *Meminfo) Close()
```

Close: Freeing 'C' structure\.

### func \(\*Meminfo\) Update

```go
func (s *Meminfo) Update() (*Meminfo, error)
```

Update: 'Meminfo' structure\.

## type NanoMeasureMethod

```go
type NanoMeasureMethod C.int
```

```go
const (
    // (0) Wall time (also known as clock time or wall-clock time) is simply
    // the total time elapsed during the measurement. It’s the time you
    // can measure with a stopwatch, assuming that you are able to start
    // and stop it exactly at the execution points you want.
    NANO_CLOCK_WALL NanoMeasureMethod = C.CLOCK_REALTIME
    // (2) CPU Time, on the other hand, refers to the time the CPU was busy
    // processing the program’s instructions. The time spent waiting for
    // other things to complete (like I/O operations) is not included in
    // the CPU time.
    NANO_CLOCK_CPUTIME NanoMeasureMethod = C.CLOCK_PROCESS_CPUTIME_ID
)
```

## type Partitions

```go
type Partitions struct {
    Details []partition
    // contains filtered or unexported fields
}
```

### func PartitionsNew

```go
func PartitionsNew() (*Partitions, error)
```

PartitionsNew: Create and initialise 'C' structure\.

### func \(\*Partitions\) Close

```go
func (s *Partitions) Close()
```

Close: Freeing 'C' structure\.

### func \(\*Partitions\) Update

```go
func (s *Partitions) Update() error
```

Update: 'Partitions' structure\.

## type PidInfos

```go
type PidInfos struct {
    Details []storeFile
    // contains filtered or unexported fields
}
```

### func PidInfosNew

```go
func PidInfosNew() (*PidInfos, error)
```

PidInfosNew: Create and initialise 'C' structure\. No need to 'free' \(close\) anything\, everything is already handled\.

### func \(\*PidInfos\) GetPidFromFilename

```go
func (sf *PidInfos) GetPidFromFilename(filename string) int
```

Get pid using the filename base\. Returns "\-1" if not found\.

### func \(\*PidInfos\) GetPidFromName

```go
func (sf *PidInfos) GetPidFromName(name string) int
```

Get pid using the name\. Note: Instead of a presious\, this function is based on a 'comm' field which contains only 16 bytes\, which means that if the name is greater than 16 characters\, it will be truncated\. Returns "\-1" if not found\.

## type ProcNetDev

```go
type ProcNetDev struct {

    // Suffix, default: 'iB' !!! max char length = 15
    Suffix string
    // Unit, default: '/s' !!! max char length = 15
    Unit       string
    Interfaces []iface
    // contains filtered or unexported fields
}
```

### func ProcNetDevNew

```go
func ProcNetDevNew(pid ...uint32) (*ProcNetDev, error)
```

ProcNetDevNew: Create and initialize the "C" structure\. If a "pid" is given\, the statistics relate to the process\. Otherwise\, it's the overall flow

### func \(\*ProcNetDev\) Close

```go
func (s *ProcNetDev) Close()
```

Close: Freeing 'C' structure\.

### func \(\*ProcNetDev\) GetAvailableInterfaces

```go
func (s *ProcNetDev) GetAvailableInterfaces() []string
```

Retrieve available network interfaces

### func \(\*ProcNetDev\) GetWanAdress

```go
func (s *ProcNetDev) GetWanAdress(adress string, useStunSrv ...bool) (string, error)
```

GetWanAdress: Retrieve wan adress using http get method or using a 'stun' server whether 'useStunSrv' was toggled\. http get: "ifconfig\.co" stun srv: "stun1\.l\.google\.com:19302"

### func \(\*ProcNetDev\) SetSuffix

```go
func (s *ProcNetDev) SetSuffix(suffix string) error
```

SetSuffix:

### func \(\*ProcNetDev\) SetUnit

```go
func (s *ProcNetDev) SetUnit(unit string) error
```

SetUnit:

### func \(\*ProcNetDev\) Update

```go
func (s *ProcNetDev) Update() error
```

Update:

## type ProcPidStat

```go
type ProcPidStat struct {
    Pid                 uint
    Comm                string
    State               string
    Ppid                int
    Pgrp                int
    Session             int
    TtyNr               int
    Tpgid               int
    Flags               uint
    Minflt              uint32
    Cminflt             uint32
    Majflt              uint32
    Cmajflt             uint32
    Utime               uint32
    Stime               uint32
    Cutime              int32
    Cstime              int32
    Priority            int32
    Nice                int32
    NumThreads          int32
    Itrealvalue         int32
    Starttime           uint64
    Vsize               uint32
    Rss                 int32
    Rsslim              uint32
    Startcode           uint32
    Endcode             uint32
    Startstack          uint32
    Kstkesp             uint32
    Kstkeip             uint32
    Signal              uint32
    Blocked             uint32
    Sigignore           uint32
    Sigcatch            uint32
    Wchan               uint32
    Nswap               uint32
    Cnswap              uint32
    ExitSignal          int
    Processor           int
    RtPriority          uint
    Policy              uint
    DelayacctBlkioTicks uint64
    GuestTime           uint32
    CguestTime          int32
    StartData           uint32
    EndData             uint32
    StartBrk            uint32
    ArgStart            uint32
    ArgEnd              uint32
    EnvStart            uint32
    EnvEnd              uint32
    ExitCode            int
    // contains filtered or unexported fields
}
```

### func ProcPidStatNew

```go
func ProcPidStatNew(pid int) (*ProcPidStat, error)
```

ProcPidStatNew: Create and initialise 'C' structure\.

### func \(\*ProcPidStat\) Close

```go
func (s *ProcPidStat) Close()
```

Close: Freeing 'C' structure\.

### func \(\*ProcPidStat\) Update

```go
func (s *ProcPidStat) Update() error
```

Update: 'ProcPidStat' structure\.

## type Smaps

\* Structures functions\, this section handle files like 'smaps'\, 'smaps\_rollup'\, \* 'maps' is not used here because information are contained inside 'smaps' \* via 'Header' variable of 'Rollup' or 'Smaps' structures

Information 'man procfs' search '/smaps' then press 'n' until '/proc/\[pid\]/smaps'

```go
type Smaps struct {
    Rollup *smapsRollup

    Smaps []smap

    // Values are given as kB inside parsed files,
    // enable this flag to convert to Bytes
    ConvertToBytes bool
    // contains filtered or unexported fields
}
```

### func SmapsNew

```go
func SmapsNew(pid int, maxReadEntries ...int) (*Smaps, error)
```

StatNew: Create a new structure that will contains required information about 'proc/\[pid\]/stat' files 'maxReadEntries' define the length of the buffer to read 'smaps' file\, default is set to 2k\.

### func \(\*Smaps\) Close

```go
func (s *Smaps) Close()
```

Close: Freeing 'C' structure\.

### func \(\*Smaps\) UpdateRollup

```go
func (s *Smaps) UpdateRollup() error
```

Update: 'C' structure content with actual values\.

### func \(\*Smaps\) UpdateSmaps

```go
func (s *Smaps) UpdateSmaps() error
```

Update: 'C' structure content with actual values\.

## type StatusFile

```go
type StatusFile struct {
    Name                     string
    Umask                    string
    State                    string
    Tgid                     uint
    Ngid                     uint
    Pid                      uint
    Ppid                     uint
    TracerPid                uint
    Uid                      *resfId
    Gid                      *resfId
    FdSize                   uint64
    Groups                   []uint
    NsTgid                   []uint
    NsPid                    []uint
    NsPgid                   []uint
    NsSid                    []uint
    Vm                       *statusFileVmem
    Threads                  int
    SigQ                     string
    SigPnd                   string
    ShdPnd                   string
    SigBlk                   string
    SigIgn                   string
    SigCgt                   string
    CapInh                   string
    CapPrm                   string
    CapEff                   string
    CapBnd                   string
    CapAmb                   string
    NoNewPrivs               int
    Seccomp                  int
    StoreBypass              string
    CpusAllowed              string
    CpusAllowedList          string
    MemsAllowed              string
    MemsAllowedList          string
    VoluntaryCtxtSwitches    uint64
    NonvoluntaryCtxtSwitches uint64
    // contains filtered or unexported fields
}
```

### func StatusFileNew

```go
func StatusFileNew(pid int) (*StatusFile, error)
```

StatusFileNew: create and initialize the "C" structure\. No need to 'free' \(close\) anything\, everything is already handled\.

## type SysTherm

Structure to hold Thermal information retrieved from: '/sys/class/hwmon/hwmon\*' directories\. Note: "n/a"\, "\-0°C" or "\-1" value means not available data\.

```go
type SysTherm struct {
    Interfaces []sysTherm
    // contains filtered or unexported fields
}
```

### func SysThermNew

```go
func SysThermNew() (*SysTherm, error)
```

### func \(\*SysTherm\) Close

```go
func (st *SysTherm) Close()
```

Close: and free memory used for structure storage

### func \(\*SysTherm\) Update

```go
func (st *SysTherm) Update() error
```

Update: 'Interfaces' data

## type TimeSpent

```go
type TimeSpent struct {
    Spent              float64
    NANO_CLOCK_WALL    NanoMeasureMethod
    NANO_CLOCK_CPUTIME NanoMeasureMethod
    SC_CLK_TCK         int64
    // contains filtered or unexported fields
}
```

### func TimeSpentNew

```go
func TimeSpentNew(method ...NanoMeasureMethod) (*TimeSpent, error)
```

TimeSpentNew: Create and initialise 'C' structure\. if argument is set to \-1\, the default value is 'NANO\_CLOCK\_WALL'

### func \(\*TimeSpent\) Close

```go
func (s *TimeSpent) Close()
```

Close: Freeing 'C' structure\.

### func \(\*TimeSpent\) MesurementMethodGet

```go
func (s *TimeSpent) MesurementMethodGet() int
```

GetMesurementMethod:

### func \(\*TimeSpent\) MesurementMethodSet

```go
func (s *TimeSpent) MesurementMethodSet(method NanoMeasureMethod)
```

SetMesurementMethod:

### func \(\*TimeSpent\) NanoCalculate

```go
func (s *TimeSpent) NanoCalculate() float64
```

NanoCalculate: calculate the nanoseconds between 2 measurement periods

### func \(\*TimeSpent\) NanoGet

```go
func (s *TimeSpent) NanoGet()
```

NanoGet: get current nano count measurement depend on defined 'method' argument 'NANO\_CLOCK\_WALL' or 'NANO\_CLOCK\_CPUTIME' Value is internally stored\.

### func \(\*TimeSpent\) SpentGet

```go
func (s *TimeSpent) SpentGet() float64
```

GetSpent: time previously calculated\.

### func \(\*TimeSpent\) TicksCalculate

```go
func (s *TimeSpent) TicksCalculate() float64
```

TicksCalculate: calculate tick between 2 tick periods

### func \(\*TimeSpent\) TicksGet

```go
func (s *TimeSpent) TicksGet()
```

TicksGet: get current ticks count\. Value is internally stored\.

