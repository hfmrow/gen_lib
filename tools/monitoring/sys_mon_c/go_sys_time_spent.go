// go_sys_time_spent.go

package sys_mon

// #include "sys_cpu.h"
import "C"

type NanoMeasureMethod C.int

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

// TimeSpentNew: Create and initialise 'C' structure.
// if argument is set to -1, the default value is 'NANO_CLOCK_WALL'
func TimeSpentNew(method ...NanoMeasureMethod) (*TimeSpent, error) {
	var mthd C.int = -1 // set to default value 'NANO_CLOCK_WALL'
	if len(method) > 0 {
		mthd = C.int(method[0])
	}
	c := C.time_spent_new(mthd)
	if c == nil {
		return nil, getErrorString()
	}
	return wrapTimeSpent(c), nil
}

// Close: Freeing 'C' structure.
func (s *TimeSpent) Close() {
	if s.time_spent != nil {
		C.time_spent_free(s.time_spent)
	}
}

// SetMesurementMethod:
func (s *TimeSpent) MesurementMethodSet(method NanoMeasureMethod) {
	s.time_spent.nano_measure_method = C.int(method)
}

// GetMesurementMethod:
func (s *TimeSpent) MesurementMethodGet() int {
	return int(s.time_spent.nano_measure_method)
}

// GetSpent: time previously calculated.
func (s *TimeSpent) SpentGet() float64 {
	return float64(s.time_spent.spent)
}

// NanoGet: get current nano count measurement depend on
// defined 'method' argument 'NANO_CLOCK_WALL' or 'NANO_CLOCK_CPUTIME'
// Value is internally stored.
func (s *TimeSpent) NanoGet() {
	C.time_nano_get(s.time_spent)
}

// NanoCalculate: calculate the nanoseconds between 2 measurement periods
func (s *TimeSpent) NanoCalculate() float64 {
	C.time_nano_calculate(s.time_spent)
	return float64(s.time_spent.spent)
}

// TicksGet: get current ticks count.
// Value is internally stored.
func (s *TimeSpent) TicksGet() {
	C.time_ticks_get(s.time_spent)
}

// TicksCalculate: calculate tick between 2 tick periods
func (s *TimeSpent) TicksCalculate() float64 {
	C.time_ticks_calculate(s.time_spent)
	return float64(s.time_spent.spent)
}

type TimeSpent struct {
	time_spent         *C.time_spent
	Spent              float64
	NANO_CLOCK_WALL    NanoMeasureMethod
	NANO_CLOCK_CPUTIME NanoMeasureMethod
	SC_CLK_TCK         int64
}

func wrapTimeSpent(time_spent *C.time_spent) *TimeSpent {
	if time_spent == nil {
		return nil
	}

	return &TimeSpent{
		time_spent,
		float64(time_spent.spent),
		NanoMeasureMethod(time_spent.NANO_CLOCK_WALL),
		NanoMeasureMethod(time_spent.NANO_CLOCK_CPUTIME),
		int64(time_spent.SC_CLK_TCK),
	}
}
