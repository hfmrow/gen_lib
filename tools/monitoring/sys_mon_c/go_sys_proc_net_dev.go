// go_sys_proc_net_dev.go

package sys_mon

// #cgo linux LDFLAGS: -lcurl
// #include "sys_proc_net_dev.h"
// #include "get_wan_ip.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type ProcNetDev struct {
	proc_net_dev *C.proc_net_dev
	// Suffix, default: 'iB' !!! max char length = 15
	Suffix string
	// Unit, default: '/s' !!! max char length = 15
	Unit       string
	Interfaces []iface
}

// GetWanAdress: Retrieve wan adress using http get method or using
// a 'stun' server whether 'useStunSrv' was toggled.
// http get: "ifconfig.co"
// stun srv: "stun1.l.google.com:19302"
func (s *ProcNetDev) GetWanAdress(adress string, useStunSrv ...bool) (string, error) {
	var outStr string
	cstr := (*C.char)(C.CString(adress))
	defer C.free(unsafe.Pointer(cstr))

	if len(useStunSrv) > 0 && useStunSrv[0] {
		c := C.return_wan_ip(cstr, 8888)
		if c == nil {
			return "", getErrorString()
		}
		defer C.free(unsafe.Pointer(c))
		outStr = C.GoString(c)
	} else {
		c := C.return_wan_ip_http_get(cstr)
		if c == nil {
			return "", getErrorString()
		}
		defer C.free(unsafe.Pointer(c))
		outStr = C.GoString(c)
	}
	return outStr, nil
}

// Close: Freeing 'C' structure.
func (s *ProcNetDev) Close() {
	if s.proc_net_dev != nil {
		C.proc_net_dev_free(s.proc_net_dev)
	}
}

// Retrieve available network interfaces
func (s *ProcNetDev) GetAvailableInterfaces() []string {
	avIf := make([]string, len(s.Interfaces))
	for idx, i := range s.Interfaces {
		avIf[idx] = i.Name
	}
	return avIf
}

// SetSuffix:
func (s *ProcNetDev) SetSuffix(suffix string) error {
	if len(suffix) > 15 {
		return fmt.Errorf("SetSuffix: max char len > 15 not allowed, char len %d", len(suffix))
	}
	cstr := (*C.char)(C.CString(suffix))
	defer C.free(unsafe.Pointer(cstr))
	s.proc_net_dev.suffix = *(*[16]C.char)(unsafe.Pointer(cstr))
	// C.iface_set_suffix(s.proc_net_dev, cstr)
	return nil
}

// SetUnit:
func (s *ProcNetDev) SetUnit(unit string) error {
	if len(unit) > 15 {
		return fmt.Errorf("SetUnit: max char len > 15 not allowed, char len %d", len(unit))
	}
	cstr := (*C.char)(C.CString(unit))
	defer C.free(unsafe.Pointer(cstr))
	s.proc_net_dev.unit = *(*[16]C.char)(unsafe.Pointer(cstr))
	// C.iface_set_unit(s.proc_net_dev, cstr)
	return nil
}

// Update:
func (s *ProcNetDev) Update() error {
	c := C.proc_net_dev_update(s.proc_net_dev)
	if !bool(c) {
		return getErrorString()
	}
	s.Interfaces = wrapProcNetDev(s.proc_net_dev).Interfaces
	return nil
}

// ProcNetDevNew: Create and initialize the "C" structure.
// If a "pid" is given, the statistics relate to the process.
// Otherwise, it's the overall flow
func ProcNetDevNew(pid ...uint32) (*ProcNetDev, error) {

	var cpid *C.uint = nil
	if len(pid) > 0 {
		cpid = new(C.uint)
		*cpid = C.uint(pid[0])
	}

	c := C.proc_net_dev_get(cpid)
	if c == nil {
		return nil, getErrorString()
	}
	return wrapProcNetDev(c), nil
}

func wrapProcNetDev(proc_net_dev *C.proc_net_dev) *ProcNetDev {
	if proc_net_dev == nil {
		return nil
	}
	count := int(proc_net_dev.count)
	ifaces := make([]iface, count)
	for i := 0; i < count; i++ {
		ifaces[i] = *wrapIface(C.iface_get_single(proc_net_dev.interfaces, C.int(i)))
	}
	return &ProcNetDev{
		proc_net_dev,
		C.GoString(&proc_net_dev.suffix[0]),
		C.GoString(&proc_net_dev.unit[0]),
		ifaces,
	}
}

type iface struct {
	iface     *C.iface
	DeltaSec  float64
	DeltaTx   uint32
	DeltaRx   uint32
	TxByteSec float64
	RxByteSec float64
	TxString  string
	RxString  string
	Tx        *netInterfaceTx
	Rx        *netInterfaceRx
	Name      string
}

func wrapIface(ciface *C.iface) *iface {
	if ciface == nil {
		return nil
	}

	return &iface{
		ciface,
		float64(ciface.delta_sec),
		uint32(ciface.delta_tx),
		uint32(ciface.delta_rx),
		float64(ciface.tx_byte_sec),
		float64(ciface.rx_byte_sec),
		C.GoString(&ciface.tx_string[0]),
		C.GoString(&ciface.rx_string[0]),
		wrapNetInterfaceTx(&ciface.tx),
		wrapNetInterfaceRx(&ciface.rx),
		C.GoString(&ciface.name[0]),
	}
}

type netInterfaceRx struct {
	net_interface_rx *C.net_interface_rx
	Bytes            uint64
	Packets          uint64
	Errs             uint64
	Drop             uint64
	Fifo             uint64
	Frame            uint64
	Compressed       uint64
	Multicast        uint64
}

func wrapNetInterfaceRx(net_interface_rx *C.net_interface_rx) *netInterfaceRx {
	if net_interface_rx == nil {
		return nil
	}

	return &netInterfaceRx{
		net_interface_rx,
		uint64(net_interface_rx.bytes),
		uint64(net_interface_rx.packets),
		uint64(net_interface_rx.errs),
		uint64(net_interface_rx.drop),
		uint64(net_interface_rx.fifo),
		uint64(net_interface_rx.frame),
		uint64(net_interface_rx.compressed),
		uint64(net_interface_rx.multicast),
	}
}

type netInterfaceTx struct {
	net_interface_tx *C.net_interface_tx
	Bytes            uint64
	Packets          uint64
	Errs             uint64
	Drop             uint64
	Fifo             uint64
	Colls            uint64
	Carrier          uint64
	Compressed       uint64
}

func wrapNetInterfaceTx(net_interface_tx *C.net_interface_tx) *netInterfaceTx {
	if net_interface_tx == nil {
		return nil
	}

	return &netInterfaceTx{
		net_interface_tx,
		uint64(net_interface_tx.bytes),
		uint64(net_interface_tx.packets),
		uint64(net_interface_tx.errs),
		uint64(net_interface_tx.drop),
		uint64(net_interface_tx.fifo),
		uint64(net_interface_tx.colls),
		uint64(net_interface_tx.carrier),
		uint64(net_interface_tx.compressed),
	}
}
