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

const (
	telegramAPIBase = "https://api.telegram.org/bot"
)

func SendToTelegram(botToken, chatID, title, description, filePath string) error {
	apiURL := fmt.Sprintf("%s%s/", telegramAPIBase, botToken)
	client := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if description != "" {
		writer.WriteField("chat_id", chatID)
		writer.WriteField("text", description)
	}
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

		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}
	writer.Close()
	req, err := http.NewRequest("POST", apiURL+"sendDocument", body)
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
		return fmt.Errorf("telegram API request failed with status %d", resp.StatusCode)
	}

	return nil
}
