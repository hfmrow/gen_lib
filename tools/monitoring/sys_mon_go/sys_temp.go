// sys_temp.go

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
	"path/filepath"
	"strconv"
	"strings"
)

type device struct {
	Name  string
	Items []*item
}

func deviceNew() *device {
	return new(device)
}

type item struct {
	Label string
	Value,
	Min,
	Max,
	Crit int64
	Files map[string]string
}

func itemNew() *item {
	itm := new(item)
	itm.Files = make(map[string]string, 0)
	return itm
}

func (sms *SystemMonitorStruct) GetTemp() (err error) {

	var (
		ok, okSub bool = true, true
		mainIdx   int
		content   []os.FileInfo
		val,
		itemName string
		files []string
		dev   *device
		itm   *item
	)

	for ok {
		root := fmt.Sprintf("%s/hwmon%d", sms.LinuxHwmonDir, mainIdx)
		if content, err = ioutil.ReadDir(root); err != nil {
			if os.IsNotExist(err) {
				err = nil
				break
			}
			return err
		}
		mainIdx++

		// TODO '.files' is not yet used, find a way to not to have to doing readdir on each call
		files = sms.sortFilesOnly(content, root, "name", "temp*")

		for idx := 0; idx < len(files); idx++ {

			lastNum := regNum.FindString(filepath.Base(files[idx]))
			itm = itemNew()
			for idx < len(files) && regNum.FindString(filepath.Base(files[idx])) == lastNum {

				okSub = false
				filename := files[idx]
				fileNamed := filename

				switch {

				case strings.HasSuffix(fileNamed, "name"):
					dev = deviceNew()
					if dev.Name, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					itemName = "name"
				case strings.HasSuffix(fileNamed, "_crit"):
					if val, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					if itm.Crit, err = strconv.ParseInt(val, 10, 64); err != nil {
						return
					}
					okSub = true
					itemName = "crit"
				case strings.HasSuffix(fileNamed, "_input"):
					if val, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					if itm.Value, err = strconv.ParseInt(val, 10, 64); err != nil {
						return
					}
					okSub = true
					itemName = "value"
				case strings.HasSuffix(fileNamed, "_label"):
					if itm.Label, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					okSub = true
					itemName = "label"
				case strings.HasSuffix(fileNamed, "_max"):
					if val, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					if itm.Max, err = strconv.ParseInt(val, 10, 64); err != nil {
						return
					}
					okSub = true
					itemName = "max"
				case strings.HasSuffix(fileNamed, "_min"):
					if val, err = sms.getValue(filename, "", "", -1); err != nil {
						return
					}
					if itm.Min, err = strconv.ParseInt(val, 10, 64); err != nil {
						return
					}
					okSub = true
					itemName = "min"
				}
				itm.Files[itemName] = filename
				idx++
			}
			if okSub {

				dev.Items = append(dev.Items, itm)
			}
			idx--
		}
		sms.Temp = append(sms.Temp, dev)
	}
	return
}
