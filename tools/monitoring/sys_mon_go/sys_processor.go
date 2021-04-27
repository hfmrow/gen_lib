// sys_processor.go

/*
	Copyright ©2020 H.F.M - system monitor library v1.0 https://github.com/hfmrow

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package sys_monitor

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpus struct {
	Items map[string]*processor
	UpTime,
	IdleTime *valueTime
}

// MeanFreq: Returns the average frequency of previously stored data
// from CPU information.
func (c *cpus) AverageFreq(hideUnit ...bool) string {

	var (
		count,
		total float64
	)

	for _, proc := range c.Items {
		count++
		total += float64(proc.Freq.Value)
	}
	return smsLocal.humanReadableFreq(total/count, hideUnit...)
}

func cpusNew() *cpus {

	c := new(cpus)
	c.Items = make(map[string]*processor, 0)
	c.UpTime = valueTimeNew()
	c.IdleTime = valueTimeNew()
	return c
}

type processor struct {
	Id,
	BaseFreq,
	Freq,
	Min,
	Max *valProc
	Governor string
	Files    map[string]string
}

type valProc struct {
	Item  string // Used when value is string format
	Value int64
}

func valProcNew(value ...int64) *valProc {

	v := new(valProc)
	if len(value) > 0 {
		v.Value = value[0]
	}
	return v
}

type valueTime struct {
	Value time.Duration
}

func valueTimeNew(value ...time.Duration) *valueTime {
	v := new(valueTime)
	if len(value) > 0 {
		v.Value = value[0]
	}
	return v
}

func (v *valueTime) String(unitSep ...string) string {
	return smsLocal.humanReadableTime(v.Value, unitSep...)
}

func (vp *valProc) String(hideUnit ...bool) string {
	if len(vp.Item) > 0 {
		return vp.Item
	}
	return smsLocal.humanReadableFreq(vp.Value, hideUnit...)
}

func processorNew() *processor {
	p := new(processor)
	p.Files = make(map[string]string, 0)
	return p
}

func (c *cpus) GetUpTime() error {
	var f float64

	lines, err := smsLocal.readLines(smsLocal.LinuxUptime)
	if err != nil {
		return err
	}
	// Gel values names. Erguments:
	// 1st, seconds since the last start
	// 2nd, time spent in idle mode)
	vars := regWhtSpaces.Split(lines[0], -1)
	if f, err = strconv.ParseFloat(vars[0], 64); err != nil {
		return err
	}
	c.UpTime.Value = time.Duration(float64(time.Second) * f)

	if f, err = strconv.ParseFloat(vars[1], 64); err != nil {
		return err
	}
	c.IdleTime.Value = time.Duration(float64(time.Second) * f)

	return nil
}

// GetProcessor: Retrieve processors frequency information
func (c *cpus) initProcessors() error {

	var (
		err     error
		mainIdx int
		i       int64
		content []os.FileInfo
		val,
		itemName string
		files []string
		cpu   *processor
	)

	for {
		root := fmt.Sprintf("%s/policy%d", smsLocal.LinuxCpufreq, mainIdx)
		if content, err = ioutil.ReadDir(root); err != nil {
			if os.IsNotExist(err) {
				err = nil
				break
			}
			return err
		}

		mainIdx++
		files = smsLocal.sortFilesOnly(content, root, "base_frequency", "scaling_*", "related_cpus")
		for idx := 0; idx < len(files); idx++ {

			cpu = processorNew()
			for idx < len(files) {

				filename := files[idx]
				fileNamed := filename
				vp := valProcNew()
				switch {
				case strings.HasSuffix(fileNamed, "related_cpus"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					if i, err = strconv.ParseInt(val, 10, 64); err != nil {
						return err
					}
					vp.Value = i
					cpu.Id = vp
					itemName = "id"
				case strings.HasSuffix(fileNamed, "scaling_governor"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					cpu.Governor = val
					itemName = "governor"
				case strings.HasSuffix(fileNamed, "base_frequency"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					if i, err = strconv.ParseInt(val, 10, 64); err != nil {
						return err
					}
					vp.Value = i
					cpu.BaseFreq = vp
					itemName = "baseFreq"
				case strings.HasSuffix(fileNamed, "scaling_cur_freq"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					if i, err = strconv.ParseInt(val, 10, 64); err != nil {
						return err
					}
					vp.Value = i
					cpu.Freq = vp
					itemName = "freq"
				case strings.HasSuffix(fileNamed, "scaling_max_freq"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					if i, err = strconv.ParseInt(val, 10, 64); err != nil {
						return err
					}
					vp.Value = i
					cpu.Max = vp
					itemName = "max"
				case strings.HasSuffix(fileNamed, "scaling_min_freq"):
					if val, err = smsLocal.getValue(filename, "", "", -1); err != nil {
						return err
					}
					if i, err = strconv.ParseInt(val, 10, 64); err != nil {
						return err
					}
					vp.Value = i
					cpu.Min = vp
					itemName = "min"
				}
				cpu.Files[itemName] = filename
				idx++
			}
			c.Items[fmt.Sprintf("%d", cpu.Id)] = cpu
			idx--
		}
	}

	return c.GetUpTime()
}
