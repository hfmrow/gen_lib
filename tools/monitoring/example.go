// monitor.go

/*
	Source file auto-generated on Mon, 26 Apr 2021 08:45:09 using Gotk3 Objects Handler v1.7.8
	©2018-21 hfmrow https://hfmrow.github.io

	Copyright ©2020-21 hfmrow - system monitor v1.0 github.com/hfmrow/gen_lib/tools/monitoring

		- network (lan/wan)
		- cpu
		- memory (map)
		- process inspection (pid)
		- thermal information
		- time measuring
		- stats (pid)
		- partitions info
		- disk stats

	This library is designed to display system information from some system files (debian file system)
	This package version, use "C" personal libraries to retrieve information from these files.

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

/*
	This is a demonstration source code to show how to use 'sys_monitor' library.
	Demonstrate lot of functionalities awailable and give an overview of possibilities.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	gltsmgsn "github.com/hfmrow/gen_lib/tools/monitoring/sys_mon_c"

	gltsushe "github.com/hfmrow/gen_lib/tools/units/human_readable"
)

var (
	HumanReadableSize = gltsushe.HumanReadableSize
)

// go clean -cache github.com/hfmrow/gen_lib/tools/monitoring/sys_mon
func main() {
	var err error

	// // var pid = 1695 // current process pid
	var pid = os.Getpid() // current process pid

	// Used to stress cpu while 'cpu %' tests
	stress := func() {
		done := make(chan int)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
					}
				}
			}()
		}
		time.Sleep(time.Second * 10)
		close(done)
	}

	// /*************************
	//  * Thermal sensors tests ok
	//  ************************/
	st, err := gltsmgsn.SysThermNew()
	defer st.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		for _, interf := range st.Interfaces {
			// fmt.Printf("Interface: %s\n", interf.Name)
			for _, sensor := range interf.Sensors {
				fmt.Printf("Label: %s, temp: %s, max: %s, crit:%s, critAlarm:%s\n",
					sensor.Label,
					sensor.TempStr,
					sensor.MaxStr,
					sensor.CritStr,
					sensor.CritAlarmStr)
			}
		}
		time.Sleep(time.Second)
		st.Update()
	}

	// /***************
	//  * Smaps tests ok
	//  **************/
	smaps, err := gltsmgsn.SmapsNew(pid)
	defer smaps.Close()
	if err != nil {
		log.Fatal(err)
	}
	// smaps_rollup test
	smaps.ConvertToBytes = true
	fmt.Println("/***********BEFORE UPDATE***********/")
	fmt.Println(smaps.Rollup.Header.ToString())
	fmt.Println(HumanReadableSize(smaps.Rollup.Rss))
	err = smaps.UpdateRollup()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("/***********AFTER UPDATE***********/")
	fmt.Println(smaps.Rollup.Header.ToString())
	fmt.Println(HumanReadableSize(smaps.Rollup.Rss))
	// smaps test
	fmt.Println("/***********BEFORE UPDATE***********/")
	for _, smp := range smaps.Smaps {
		fmt.Println(smp.Header.Pathname, smp.Rss)
	}
	err = smaps.UpdateSmaps()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("/***********AFTER UPDATE***********/")
	for _, smp := range smaps.Smaps {
		fmt.Println(smp.Header.Pathname, smp.Rss)
	}

	// /***************
	//  * CpuFs tests ok
	//  **************/
	cpuFs, err := gltsmgsn.CpuFsNew()
	defer cpuFs.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cpu count: %d\n", cpuFs.CpuCount)
	for idx, vals := range cpuFs.CurrFreq {
		fmt.Printf("Cpu #%d: %dHz\n", idx, vals)
	}
	for idx, vals := range cpuFs.CpuList {
		fmt.Printf("Cpu #%d base freq: %dHz, min freq: %dHz, max freq: %dHz, curr governor: %s, curr perfPref: %s, curr driver: %s\n",
			idx,
			vals.BaseFrequency,
			vals.CpuinfoMinFreq,
			vals.CpuinfoMaxFreq,
			vals.ScalingGovernor,
			vals.EnergyPerformancePreference,
			vals.ScalingDriver)
	}

	// /***************
	//  * Cpu% tests ok
	//  **************/
	cpuPercent, err := gltsmgsn.CpuPercentPidNew(pid)
	defer cpuPercent.Close()
	if err != nil {
		log.Fatal(err)
	}
	go stress()
	for i := 0; i < 11; i++ {
		err = cpuPercent.Update()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		if i != 0 {
			fmt.Printf("Cpu:%.2f%%, Rss: %s\n", cpuPercent.CpuPercent, HumanReadableSize(cpuPercent.MemoryRss))
		}
	}

	// /*******************
	//  * pid infos tests ok
	//  ******************/
	StoreFilesPid, err := gltsmgsn.PidInfosNew()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range StoreFilesPid.Details {
		fmt.Printf("Pid: %d, Name: %s, file: %s, dir: %s, uidReal %d, gidReal %d\n",
			item.Pid,
			item.Name,
			item.Filename,
			item.Dirname,
			item.Uid.Real,
			item.Gid.Real)
	}
	fmt.Printf("getpid_from_name: %d\n", StoreFilesPid.GetPidFromFilename("ghb"))

	// /*******************
	//  * pid stat tests ok
	//  ******************/
	stat, err := gltsmgsn.ProcPidStatNew(pid)
	defer stat.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Pid: %d, Comm: %s, State: %s, Prio: %d, Vsize %s, Cnswap %d\n",
		stat.Pid,
		stat.Comm,
		stat.State,
		stat.Priority,
		HumanReadableSize(stat.Vsize),
		stat.Cnswap)

	// /********************
	//  * Status file tests ok
	//  ********************/
	sf, err := gltsmgsn.StatusFileNew(pid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name: %s, VmPeak: %s, VmRss: %s, VmSwap: %s, State: %s\n",
		sf.Name,
		HumanReadableSize(sf.Vm.VmSwap*1024),
		HumanReadableSize(sf.Vm.VmRss*1024),
		HumanReadableSize(sf.Vm.VmPeak*1024),
		sf.State,
	)
	fmt.Printf("Groups:")
	for _, val := range sf.Groups {
		fmt.Printf("%d ", val)
	}
	fmt.Println()

	// /*****************
	//  * Network tests
	//  *****************/
	net, err := gltsmgsn.ProcNetDevNew()
	defer net.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		wan, err := net.GetWanAdress("stun1.l.google.com:19302", true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wan adress(stan-srv):", wan)
		wan, err = net.GetWanAdress("checkip.amazonaws.com")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wan adress(http-get):", wan)
	}
	for _, card := range net.Interfaces {
		fmt.Printf("%s, ", card.Name)
	}
	fmt.Println(net.GetAvailableInterfaces())
	net.SetSuffix("B")
	net.SetUnit("/sec")
	for i := 0; i < 5; i++ {
		fmt.Printf("name: %s, tx: %s, rx: %s, deltasec: %.8f\n", net.Interfaces[1].Name, net.Interfaces[1].TxString, net.Interfaces[1].RxString, net.Interfaces[1].DeltaSec)
		time.Sleep(time.Millisecond * 1000)
		err = net.Update()
		if err != nil {
			log.Fatal(err)
		}
	}

	// /***************************
	//  * Meminfo tests
	//  ***************************/
	meminfo, err := gltsmgsn.MeminfoNew()
	defer meminfo.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("Ram total: %s, Free: %s, Used: %s, Swap total: %s, Free: %s, Used: %s, Cached: %s\n",
			HumanReadableSize(meminfo.MemTotal),
			HumanReadableSize(meminfo.MemAvailable),
			HumanReadableSize(meminfo.MemTotal-meminfo.MemAvailable),
			HumanReadableSize(meminfo.SwapTotal),
			HumanReadableSize(meminfo.SwapFree),
			HumanReadableSize(meminfo.SwapTotal-meminfo.SwapFree),
			HumanReadableSize(meminfo.SwapCached))
		meminfo.Update()
		time.Sleep(time.Millisecond * 500)
	}

	// /***************************
	//  * Partitions tests
	//  ***************************/
	part, err := gltsmgsn.PartitionsNew()
	defer part.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range part.Details {
		fmt.Printf("Name: %s, major: %d, minor: %d, Size: %s, Type: %s, uuid: %s, partName: %s\n",
			item.Name,
			item.Major,
			item.Minor,
			HumanReadableSize(float64(item.Size)*1024),
			item.ClassBlock.DevType,
			item.ClassBlock.Uuid,
			item.ClassBlock.PartName,
		)
	}

	// /***************************
	//  * Diskstats tests
	//  ***************************/
	diskstats, err := gltsmgsn.DiskstatsNew()
	defer diskstats.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range diskstats.Details {
		fmt.Printf("Name: %s, Type: %s, SectorsRead: %d, ReadsCompleted: %d, SectorsDiscarded: %d, SectorsWritten: %d, WritesCompleted: %d, TimeWritingMs: %d, IosInProgress: %d, TimeDoingIosMs: %d\n",
			item.Device,
			item.DevType,
			item.SectorsRead,
			item.ReadsCompleted,
			item.SectorsDiscarded,
			item.SectorsWritten,
			item.WritesCompleted,
			item.TimeWritingMs,
			item.IosInProgress,
			item.TimeDoingIosMs,
		)
	}

	// /***************************
	//  * Measurement 'TimeSpent' tests
	//  ***************************/
	ts, err := gltsmgsn.TimeSpentNew(gltsmgsn.NANO_CLOCK_CPUTIME) // or
	// ts, err := gltsmgsn.TimeSpentNew()
	defer ts.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("NANO_CLOCK_WALL: %d, NANO_CLOCK_CPUTIME: %d\n", gltsmgsn.NANO_CLOCK_WALL, gltsmgsn.NANO_CLOCK_CPUTIME)
	/* nano tests */
	ts.NanoGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%0.9fSec\n", ts.NanoCalculate())
	fmt.Printf("method: %d\n", ts.MesurementMethodGet())
	ts.MesurementMethodSet(ts.NANO_CLOCK_WALL)
	ts.NanoGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%0.9fSec\n", ts.NanoCalculate())
	fmt.Printf("method: %d\n", ts.MesurementMethodGet())
	ts.MesurementMethodSet(ts.NANO_CLOCK_CPUTIME)
	ts.NanoGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%0.9fSec\n", ts.NanoCalculate())
	fmt.Printf("method: %d\n", ts.MesurementMethodGet())
	ts.MesurementMethodSet(ts.NANO_CLOCK_WALL)
	ts.NanoGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%0.9fSec\n", ts.NanoCalculate())
	fmt.Printf("method: %d\n", ts.MesurementMethodGet())
	/* ticks tests */
	ts.TicksGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%f\n", ts.TicksCalculate())
	ts.TicksGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%f\n", ts.TicksCalculate())
	ts.TicksGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%f\n", ts.TicksCalculate())
	ts.TicksGet()
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("%f\n", ts.TicksCalculate())

}
