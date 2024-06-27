package Mutex

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

func Create() bool {
	const AppID = "3575651c-bb47-448e-a514-22865732bbc"

	_, err := windows.CreateMutex(nil, false, syscall.StringToUTF16Ptr(fmt.Sprintf("Global\\%s", AppID)))
	return err != nil
}
