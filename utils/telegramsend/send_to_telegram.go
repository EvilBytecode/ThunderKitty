package requests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url" // Import the url package
	"os"
	"path/filepath"
)

// SendTelegramMessage sends a text message to a Telegram chat.
func SendTelegramMessage(botToken, chatID, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Prepare the request payload
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	// Send the POST request
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Telegram API request failed with status %d: %s\n", resp.StatusCode, string(responseBody))
		return fmt.Errorf("failed to send message: %s", string(responseBody))
	}

	return nil
}

// SendTelegramDocument sends a document to a Telegram chat.
func SendTelegramDocument(botToken, chatID, filePath string) error {
	fmt.Println("Sending document to Telegram:", filePath)
	apiBaseURL := "https://api.telegram.org/bot"
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Always add chat_id
	writer.WriteField("chat_id", chatID)

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

	apiMethod := "sendDocument"
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
		return fmt.Errorf("failed to send document: %s", string(responseBody))
	}

	return nil
}
