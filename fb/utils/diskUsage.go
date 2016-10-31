package utils

import "syscall"

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

//import (
//"syscall"
//"unsafe"
//)
//
//func DiskUsage(path string) (Info, error) {
//	i := Info{}
//	h, err := syscall.LoadDLL("kernel32.dll")
//	if err != nil {
//		return i, err
//	}
//	c, err := h.FindProc("GetDiskFreeSpaceExW")
//	if err != nil {
//		return i, err
//	}
//	_, _, err = c.Call(
//		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
//		uintptr(unsafe.Pointer(&i.Free)),
//		uintptr(unsafe.Pointer(&i.Total)),
//		uintptr(unsafe.Pointer(&i.Available)))
//
//	return i, err
//}

