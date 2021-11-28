package sysinfo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type SysInfo struct {
	OS OS `json:"os"`
	Board Board `json:"board"`
	CPU CPU `json:"cpu"`
	BIOS BIOS `json:"bios"`
	Memory Memory `json:"memory"`
	Storage []StorageDevice `json:"storage"`
	Network []NetworkDevice `json:"network"`
}

type OS struct {
	Host string `json:"host"`
	Name string `json:"name"`
	Vendor string `json:"vendor"`
	Version string `json:"version"`
	Arch string `json:"arch"`
	Serial string `json:"serial"`
	Installed string `json:"installed"`
}

type CPU struct {
	Vendor string `json:"vendor"`
	Model string `json:"model"`
	Cpus uint32 `json:"cpus"`
	Cores uint32 `json:"cores"`
	Threads uint32 `json:"threads"`
	Serial string `json:"serial"`
}

type BIOS struct {
	Vendor string `json:"vendor"`
	Version string `json:"version"`
	Date string `json:"date"`
}

type Board struct {
	Name string `json:"name"`
	Vendor string `json:"vendor"`
	Version string `json:"version"`
	Serial string `json:"serial"`
}

type Memory struct {
	Vendor string `json:"vendor"`
	Size string `json:"size"`
}

type StorageDevice struct {
	Name string `json:"name"`
	Vendor string `json:"vendor"`
	Model string `json:"model"`
	Serial string `json:"serial"`
	Size string `json:"size"`
}

type NetworkDevice struct {
	Name string `json:"name"`
	MAC string `json:"mac"`
}

const (
	KB = 1000
	MB = KB * 1000
	GB = MB * 1000
	TB = GB * 1000
	PB = TB * 1000
)

func HumanFriendlyTraffic(bytes uint64) string {
	if bytes <= KB {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes <= MB {
		return fmt.Sprintf("%.2f KB", float32(bytes)/KB)
	}
	if bytes <= GB {
		return fmt.Sprintf("%.2f MB", float32(bytes)/MB)
	}
	if bytes <= TB {
		return fmt.Sprintf("%.2f GB", float32(bytes)/GB)
	}
	return fmt.Sprintf("%.2f TB", float32(bytes)/TB)
}

func (c *SysInfo)ToJSON() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func Get() (element *SysInfo, err error) {
	switch strings.ToLower(runtime.GOOS) {
	case "windows":
		return getSysInfoForWindows()
	case "linux":
		return getSysInfoForWindows()
	case "darwin":
		return getSysInfoForWindows()
	}
	err = fmt.Errorf("bad request")
	return
}