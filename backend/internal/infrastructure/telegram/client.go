package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/n1x9s/second-brain/backend/internal/domain"
)

type Client struct {
	botToken string
	client   *http.Client
}

func NewClient(botToken string) Client {
	return Client{botToken: botToken, client: &http.Client{Timeout: 8 * time.Second}}
}

func (c Client) Send(ctx context.Context, chatID string, message string) error {
	if c.botToken == "" || chatID == "" {
		return domain.ErrProviderEmpty
	}
	payload, _ := json.Marshal(map[string]string{
		"chat_id": chatID,
		"text":    message,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.botToken), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram sendMessage failed: %s", resp.Status)
	}
	return nil
}
