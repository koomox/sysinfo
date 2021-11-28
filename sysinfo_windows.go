package sysinfo

import (
	"github.com/koomox/sysinfo/providers/windows"
	"strings"
)

type sysInfoForWindows struct {
	OperatingSystem []windows.OperatingSystem
	BoardDevice []windows.BoardDevice
	ProcessorDevice []windows.ProcessorDevice
	BiosEx []windows.BiosEx
	MemoryDevice windows.MemoryDevice
	DiskDrive []windows.DiskDrive
	NetworkDevice []windows.NetworkDevice
}

func getSysInfoForWindows() (element *SysInfo, err error) {
	var (
		info = &sysInfoForWindows{}
		cpu *CPU
		system *OS
		board *Board
		memory *Memory
		bios *BIOS
		disk []StorageDevice
		network []NetworkDevice
	)
	if cpu, err = info.CPU(); err != nil {
		return
	}
	if system, err = info.System(); err != nil {
		return
	}
	if board, err = info.Board(); err != nil {
		return
	}
	memory = info.Memory()
	if bios, err = info.BIOS(); err != nil {
		return
	}
	if disk, err = info.Disk(); err != nil {
		return
	}
	if network, err = info.Network(); err != nil {
		return
	}
	element = &SysInfo{
		OS: *system,
		Board: *board,
		CPU: *cpu,
		BIOS: *bios,
		Memory: *memory,
		Storage: disk,
		Network: network,
	}
	return
}

func (c *sysInfoForWindows)CPU() (cpu *CPU, err error) {
	if c.ProcessorDevice, err = windows.CPU(); err != nil {
		return
	}
	for _, v := range c.ProcessorDevice {
		cpu = &CPU{
			Vendor: v.Manufacturer,
			Model: v.Name,
			Cpus: uint32(len(c.ProcessorDevice)),
			Cores: v.NumberOfCores,
			Threads: v.ThreadCount,
			Serial: v.ProcessorId,
		}
		break
	}

	return
}

func (c *sysInfoForWindows)Memory() *Memory {
	mem := windows.Memory()
	return &Memory{Size: HumanFriendlyTraffic(mem.TotalVisibleMemorySize)}
}

func (c *sysInfoForWindows)System() (s *OS, err error) {
	if c.OperatingSystem, err = windows.System(); err != nil {
		return
	}
	for _, v := range c.OperatingSystem {
		s = &OS{
			Host: v.CSName,
			Name: v.Caption,
			Vendor: v.Manufacturer,
			Version: v.Version,
			Arch: v.OSArchitecture,
			Serial: v.SerialNumber,
			Installed: v.InstallDate.Format("2006-01-02 15:04:05"),
		}
		break
	}

	return
}

func (c *sysInfoForWindows)Board() (b *Board, err error) {
	if c.BoardDevice, err = windows.Board(); err != nil {
		return
	}
	for _, v := range c.BoardDevice {
		b = &Board{
			Name: v.Product,
			Vendor: v.Manufacturer,
		}
		break
	}

	return
}

func (c *sysInfoForWindows)BIOS() (b *BIOS, err error) {
	if c.BiosEx, err = windows.BIOS(); err != nil {
		return
	}
	for _, v := range c.BiosEx {
		b = &BIOS{
			Vendor: v.Manufacturer,
			Version: v.Version,
			Date: v.ReleaseDate.Format("2006-01-02"),
		}
		break
	}

	return
}

func (c *sysInfoForWindows)Disk() (elements []StorageDevice, err error) {
	if c.DiskDrive, err = windows.Disk(); err != nil {
		return
	}
	for _, v := range c.DiskDrive {
		elements = append(elements, StorageDevice{
			Model: v.Model,
			Serial: strings.TrimSpace(v.SerialNumber),
			Size: HumanFriendlySize(v.Size),
		})
	}

	return
}

func (c *sysInfoForWindows)Network() (elements []NetworkDevice, err error) {
	if c.NetworkDevice, err = windows.Network(); err != nil {
		return
	}
	for _, v := range c.NetworkDevice {
		elements = append(elements, NetworkDevice{
			Name: v.Name,
			MAC: v.MAC,
		})
	}

	return
}