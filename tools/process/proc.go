// proc.go

/*
	Copyright ©2020 H.F.M - Process library v1.0 part of H.F.M gen_lib
	Include part of "https://github.com/shirou/gopsutil/tree/master/process" library made under BSD license
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package process

import (
	"fmt"
	"os"
	"path/filepath"

	proc "github.com/shirou/gopsutil/process"
)

// GetProcessesNames(): Retrieve parents processes. Unlike the
// other Process() function (included in this package), this
// one get only currently active 1st parent processes.
// Target User(s) may be defined.
func GetProcessesNames(userFilter ...string) (nameList []string) {
	var allProcesses []*proc.Process
	var singleProcess, parentProcess *proc.Process
	var err error
	var nameProcParent, userName string
	var wantedUser bool

	allProcesses, err = proc.Processes()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, singleProcess = range allProcesses {
		userName, _ = singleProcess.Username()
		wantedUser = false
		// User(s) filtering
		if len(userFilter) != 0 {
			for _, uName := range userFilter {
				if uName == userName {
					wantedUser = true
				}
			}
			if !wantedUser {
				continue
			}
		}
		nameProcParent, _ = singleProcess.Name()
		// recursive Check for 1st parent
		for err == nil {
			if parentProcess, err = singleProcess.Parent(); err == nil {
				singleProcess = parentProcess
			} else {
				nameProcParent, _ = singleProcess.Name()
			}
		}
		nameList = append(nameList, filepath.Base(nameProcParent))
	}
	return
}
