// network-monitor.go

/*
	Based on the work from: 'github.com/cs8425/NetTop' made under MIT license.

	Copyright ©2020 H.F.M - Network bandwidth monitor library v1.0 https://github.com/hfmrow

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package network_monitor

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var interv time.Duration

// NetMonitor: structure that hold methods to monitoring
// newtwork interfaces bandwidth.
type NetMonitor struct {
	done    chan bool
	Running bool

	// Hold error if exists
	Err error

	delta,
	last *NetStat

	t0 time.Time
	dt,
	interval time.Duration

	LinuxNetDevDir,
	currInterface string

	Interfaces []string
}

// NetMonitorNew: Create a new structure that hold methods
// to monitoring newtwork interfaces.
func NetMonitorNew(interf string, interval float64) (nt *NetMonitor, err error) {

	nt = new(NetMonitor)
	nt = &NetMonitor{
		delta: newNetStat(),
		last:  newNetStat(),
		t0:    time.Now(),
		// dt:             time.Duration(interval*1000) * time.Millisecond,
		LinuxNetDevDir: "/proc/net/dev",
		done:           make(chan bool),
	}

	// Getting all available interfaces and store them to structure
	if nt.Interfaces, err = nt.GetAvailableInterfaces(); err != nil {
		return nil, err
	}
	nt.SetInterval(interval)
	return nt, nt.SetInterface(interf)
}

type NetStat struct {
	Dev  []string
	Stat map[string]*devStat
}

func newNetStat() *NetStat {
	return &NetStat{
		Dev:  make([]string, 0),
		Stat: make(map[string]*devStat),
	}
}

type devStat struct {
	Name string
	Rx,
	Tx,
	RxSinceStart,
	TxSinceStart trans
}

type trans struct {
	val uint64
	dt  time.Duration
}

// String: output as human readable string.
// 'reduceUnit': display 'K/s' instead of 'KiB/s'
func (tr *trans) String(reduceUnit ...bool) string {

	unit := "iB"
	if len(reduceUnit) > 0 && reduceUnit[0] {
		unit = ""
	}
	return tr.vSize(unit)
}

func (tr *trans) Value() uint64 {
	return tr.val
}

func (tr *trans) vSize(ext string) string {

	var (
		tmp   = float64(tr.val) / tr.dt.Seconds()
		bytes = uint64(tmp)
		unit  string
	)

	switch {
	case bytes < uint64(2<<9):
		return fmt.Sprintf("%.0f%sB/s", tmp, unit)

	case bytes < uint64(2<<19):
		tmp = tmp / float64(2<<9)
		unit = "K"

	case bytes < uint64(2<<29):
		tmp = tmp / float64(2<<19)
		unit = "M"

	case bytes < uint64(2<<39):
		tmp = tmp / float64(2<<29)
		unit = "G"

	case bytes < uint64(2<<49):
		tmp = tmp / float64(2<<39)
		unit = "T"

	}
	return fmt.Sprintf("%2.2f%s%s/s", tmp, unit, ext)
}

// GetAvailableInterfaces: Getting all available interfaces
func (nt *NetMonitor) GetAvailableInterfaces() ([]string, error) {

	var backInterf = nt.currInterface

	nt.currInterface = "*"
	if info, err := nt.getInfo(); err == nil {
		nt.Interfaces = info.Dev
	} else {
		return []string{}, err
	}
	nt.currInterface = backInterf
	return nt.Interfaces, nil
}

func (nt *NetMonitor) SetInterval(interval float64) {
	if interval < 0.01 {
		interval = 0.01
	}
	nt.interval = time.Duration(interval*1000) * time.Millisecond
	interv = nt.interval
}

func (nt *NetMonitor) GetCurrInterface() string {
	return nt.currInterface
}

func (nt *NetMonitor) SetInterface(interf string) error {

	if len(interf) > 0 {
		// Check if the desired interface exists
		for _, iface := range nt.Interfaces {
			if iface == interf {
				nt.currInterface = interf
				return nil
			}
		}
	} else {
		if len(nt.Interfaces) > 0 {
			nt.currInterface = nt.Interfaces[0]
			return nil
		}
		return nil
	}

	return fmt.Errorf("Unavailable network interface: %s", interf)
}

func (nt *NetMonitor) Stop() {
	if nt.Running {
		nt.done <- true
	}
}

func (nt *NetMonitor) Start(callback func(stats *NetStat)) {
	go nt.start(callback)
}

func (nt *NetMonitor) start(callback func(stats *NetStat)) {

	nt.Running = true
	for {
		select {
		case <-nt.done:
			nt.Running = false
			return
		default:
			if stat, ok := nt.update(); ok {
				callback(stat)
			} else {
				if nt.Err != nil {
					nt.Stop()
					callback(nil)
				}
			}
		}
		time.Sleep(nt.interval)
	}
}

func (nt *NetMonitor) update() (*NetStat, bool) {

	nt.Err = nil
	outOk := false
	stat1, err := nt.getInfo()
	if err != nil {
		nt.Err = err
		return nil, false
	}
	nt.dt = time.Since(nt.t0)

	for _, value := range stat1.Dev {
		t0, ok := nt.last.Stat[value]

		if !ok {
			continue
		}

		dev, ok := nt.delta.Stat[value]
		if !ok {
			nt.delta.Stat[value] = new(devStat)
			dev = nt.delta.Stat[value]
			nt.delta.Dev = append(nt.delta.Dev, value)
		}

		t1 := stat1.Stat[value]
		dev.Rx.val = t1.Rx.val - t0.Rx.val
		dev.RxSinceStart.val += dev.Rx.val // TODO for further usage (cumulative value)
		dev.Rx.dt = nt.dt

		dev.Tx.val = t1.Tx.val - t0.Tx.val
		dev.TxSinceStart.val += dev.Tx.val // Same as above
		dev.Tx.dt = nt.dt

		outOk = true
	}
	nt.last = &stat1
	nt.t0 = time.Now()

	return nt.delta, outOk
}

func (nt *NetMonitor) getInfo() (ret NetStat, errOut error) {

	lines, _ := nt.readLines(nt.LinuxNetDevDir)

	ret.Dev = make([]string, 0)
	ret.Stat = make(map[string]*devStat)

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.Fields(strings.TrimSpace(fields[1]))

		if nt.currInterface != "*" && nt.currInterface != key {
			continue
		}

		c := new(devStat)
		c.Name = key
		r, err := strconv.ParseInt(value[0], 10, 64)
		if err != nil {
			errOut = fmt.Errorf("Rx: %v, %v", value[0], err)
			break
		}
		c.Rx.val = uint64(r)

		t, err := strconv.ParseInt(value[8], 10, 64)
		if err != nil {
			errOut = fmt.Errorf("Tx: %v, %v", value[8], err)
			break
		}
		c.Tx.val = uint64(t)

		ret.Dev = append(ret.Dev, key)
		ret.Stat[key] = c
	}

	return
}

func (nt *NetMonitor) readLines(filename string) ([]string, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, nil
	}

	return strings.Split(string(data), "\n"), nil
}
