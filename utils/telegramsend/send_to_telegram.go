package requests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// SendToTelegram sends a message or a file to a Telegram chat.
func SendToTelegram(botToken, chatID, message, filePath string) error {
	fmt.Println("Sending to telegram", filePath)
	apiBaseURL := "https://api.telegram.org/bot"
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Always add chat_id
	writer.WriteField("chat_id", chatID)

	// Add message if provided
	if message != "" {
		writer.WriteField("text", message)
	}

	// Add file if provided
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("document", filepath.Base(filePath))
		if err != nil {
			return err
		}

		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}

	writer.Close()

	// Find correct endpoint (file/no file)
	apiMethod := "sendMessage"
	if filePath != "" {
		apiMethod = "sendDocument"
	}

	apiURL := fmt.Sprintf("%s%s/%s", apiBaseURL, botToken, apiMethod)
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Telegram API request failed with status %d: %s\n", resp.StatusCode, string(responseBody))
		return nil
	}

	return nil
}

// Separate function for sending only message
// This will have to be refactored
func SendTelegramMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", botToken, chatID, message)
	_, err := http.Get(url)
	return err
}
