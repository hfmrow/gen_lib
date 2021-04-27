// sys_memory.go

/*
	Copyright ©2020 H.F.M - system monitor library v1.0 https://github.com/hfmrow

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package sys_monitor

import (
	"fmt"
	"strconv"
)

// swaps: structure that hold swaps info values
type swaps struct {
	Values map[string]*valueSwap
	names  []string
}

func swapsNew() *swaps {
	s := new(swaps)
	s.Values = make(map[string]*valueSwap, 0)
	s.names = []string{"MemTotal", "MemAvailable", "SwapTotal", "SwapFree"}
	return s
}

type valueSwap struct {
	Value int64
	Filename,
	Name,
	Type string
}

func valueSwapNew(value ...int64) *valueSwap {

	v := new(valueSwap)
	if len(value) > 0 {
		v.Value = value[0]
	}
	return v
}

func (v *valueSwap) String() string {

	switch {
	case v.Name == "Filename":
		return v.Filename
	case v.Name == "Type":
		return v.Type
	case v.Name == "Priority":
		return fmt.Sprintf("%d", v.Value)
	}
	return smsLocal.humanReadableSize(v.Value)
}

// getSwaps: Retrieve Swaps information.
// Note: Original value are given as KB, this function retrieve them as bytes
func (s *swaps) getSwaps() error {

	var multiplier int64

	lines, err := smsLocal.readLines(smsLocal.LinuxSwaps)
	if err != nil {
		return err
	}

	// Gel values names
	vars := regWhtSpaces.Split(lines[0], -1)

	// Retrieve values
	for idx := 1; idx < len(lines); idx++ {
		values := regWhtSpaces.Split(lines[idx], -1)
		for i, v := range values {

			vs := valueSwapNew()
			if val, err := strconv.ParseInt(v, 10, 64); err != nil {
				switch i {
				case 0:
					vs.Filename = v
				case 1:
					vs.Type = v
				}
			} else {
				switch vars[i] {
				case "Priority":
					multiplier = 1
				default:
					multiplier = 1024
				}
				vs.Value = val * multiplier // Convert to bytes
			}
			vs.Name = vars[i]
			s.Values[vars[i]] = vs
		}
	}

	return nil
}
