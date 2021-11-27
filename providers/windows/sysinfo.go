package windows

import (
	"github.com/StackExchange/wmi"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

type BoardDevice struct {
	Manufacturer string `json:"vendor"`
	Product string `json:"device"`
}

type ProcessorDevice struct {
	Name string `json:"name"`
	NumberOfCores uint32 `json:"cores"`
	ThreadCount uint32 `json:"threads"`
	ProcessorId string `json:"id"`
	Manufacturer string `json:"vendor"`
}

type BiosEx struct {
	Manufacturer string `json:"vendor"`
	ReleaseDate time.Time `json:"date"`
	Version string `json:"version"`
}

type MemoryDevice struct {
	TotalVisibleMemorySize uint64 `json:"size"`
}

type DiskDrive struct {
	Model string `json:"model"`
	Name string `json:"name"`
	SerialNumber string `json:"serial"`
	Size uint64 `json:"size"`
}

type NetworkDevice struct {
	Name string `json:"name"`
	MAC string `json:"mac"`
}

type OperatingSystem struct {
	Caption string `json:"os"`
	CSName string `json:"device"`
	InstallDate  time.Time `json:"installed"`
	LastBootUpTime time.Time `json:"last"`
	Name string `json:"name"`
	OSArchitecture string `json:"arch"`
	SerialNumber string `json:"sn"`
	TotalVisibleMemorySize uint64 `json:"memory"`
	Version string `json:"version"`
	Manufacturer string `json:"vendor"`
}

func Board() (elements []BoardDevice, err error){
	err = wmi.Query("Select * from Win32_BaseBoard", &elements)
	return
}

func CPU() (elements []ProcessorDevice, err error){
	err = wmi.Query("Select * from Win32_Processor", &elements)
	return
}

func sysTotalMemory() uint64 {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return 0
	}
	// GetPhysicallyInstalledSystemMemory is simpler, but broken on
	// older versions of windows (and uses this under the hood anyway).
	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return 0
	}
	msx := &memoryStatusEx{
		Length: 64,
	}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return 0
	}
	return msx.TotalPhys
}

func Memory() *MemoryDevice{
	return &MemoryDevice{TotalVisibleMemorySize: sysTotalMemory()}
}

func Disk() (elements []DiskDrive, err error){
	err = wmi.Query("Select * from Win32_DiskDrive", &elements)
	return
}

func validInterface(ips []net.Addr) bool {
	for _, ip := range ips {
		if addr, ok := ip.(*net.IPNet); ok && !addr.IP.IsLoopback() {
			return true
		}
	}
	return false
}

func Network() (elements []NetworkDevice, err error){
	inets, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, v := range inets {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		if strings.Contains(v.Name, "Bluetooth") || strings.Contains(v.Name, "VMware") || strings.Contains(v.Name, "Virtual") {
			continue
		}
		if validInterface(ips) {
			elements = append(elements, NetworkDevice{Name: v.Name, MAC: v.HardwareAddr.String()})
		}
	}

	return
}

func System() (elements []OperatingSystem, err error){
	err = wmi.Query("Select * from Win32_OperatingSystem", &elements)
	return
}

func BIOS() (elements []BiosEx, err error){
	err = wmi.Query("Select * from Win32_BIOS", &elements)
	return
}