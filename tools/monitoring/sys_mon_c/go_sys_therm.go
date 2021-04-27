// go_sys_therm.go

package sys_mon

// #include "sys_therm.h"
import "C"

// Structure to hold Thermal information retrieved from:
// '/sys/class/hwmon/hwmon*' directories.
// Note: "n/a", "-0°C" or "-1" value means not available data.
type SysTherm struct {
	sys_therm  *C.sys_therm
	Interfaces []sysTherm
}

func SysThermNew() (*SysTherm, error) {
	st := new(SysTherm)
	if st.sys_therm = C.sys_therm_get(); st.sys_therm == nil {
		return nil, getErrorString()
	}
	if err := st.wrapSysThermInternal(); err != nil {
		return nil, err
	}
	return st, nil
}

func (st *SysTherm) wrapSysThermInternal() error {
	count := int(st.sys_therm.count)
	st.Interfaces = make([]sysTherm, count)
	for i := 0; i < count; i++ {
		g := wrapSysTherm(C.sys_therm_get_single(st.sys_therm, C.int(i)))
		if g == nil {
			return getErrorString()
		}
		st.Interfaces[i] = *g
	}
	return nil
}

// Update: 'Interfaces' data
func (st *SysTherm) Update() error {
	C.sys_therm_update(st.sys_therm)
	err := st.wrapSysThermInternal()
	if err != nil {
		return getErrorString()
	}
	return nil
}

// Close: and free memory used for structure storage
func (st *SysTherm) Close() {
	if st.sys_therm != nil {
		C.sys_therm_free(st.sys_therm)
	}
}

type sysTherm struct {
	sys_therm *C.sys_therm
	Name      string
	Path      string
	Sensors   []sysThermSensor
}

func (v *sysTherm) native() *C.sys_therm {
	if v == nil || v.sys_therm == nil {
		return nil
	}

	return v.sys_therm
}

func wrapSysTherm(sys_therm *C.sys_therm) *sysTherm {
	if sys_therm == nil {
		return nil
	}
	count := int(sys_therm.sensors_count)
	sensors := make([]sysThermSensor, count)
	for i := 0; i < count; i++ {
		sensors[i] = *wrapSysThermSensor(C.sys_therm_sensor_get_single(sys_therm.sensors, C.int(i)))
	}
	return &sysTherm{
		sys_therm,
		C.GoString(&sys_therm.name[0]),
		C.GoString(&sys_therm.path[0]),
		sensors,
	}
}

type sysThermSensor struct {
	sys_therm_sensor *C.sys_therm_sensor
	Temp             int
	Max              int
	Crit             int
	CritAlarm        int
	TempStr          string
	MaxStr           string
	CritStr          string
	CritAlarmStr     string
	Label            string
}

func (v *sysThermSensor) native() *C.sys_therm_sensor {
	if v == nil || v.sys_therm_sensor == nil {
		return nil
	}

	return v.sys_therm_sensor
}

func wrapSysThermSensor(sys_therm_sensor *C.sys_therm_sensor) *sysThermSensor {
	if sys_therm_sensor == nil {
		return nil
	}

	return &sysThermSensor{
		sys_therm_sensor,
		int(sys_therm_sensor.temp),
		int(sys_therm_sensor.max),
		int(sys_therm_sensor.crit),
		int(sys_therm_sensor.crit_alarm),
		C.GoString(&sys_therm_sensor.temp_str[0]),
		C.GoString(&sys_therm_sensor.max_str[0]),
		C.GoString(&sys_therm_sensor.crit_str[0]),
		C.GoString(&sys_therm_sensor.crit_alarm_str[0]),
		C.GoString(&sys_therm_sensor.label[0]),
	}
}
