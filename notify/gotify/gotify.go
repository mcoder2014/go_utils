package gotify

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerOption struct {
	// https://host:port
	ServerURL string
	// abcdefg
	Token string
}

// GotifyMessage represents the payload for sending a message
type GotifyMessage struct {
	Title    string                 `json:"title,omitempty"`
	Message  string                 `json:"message"`
	Priority int                    `json:"priority,omitempty"`
	Extras   map[string]interface{} `json:"extras,omitempty"`

	ServerOption *ServerOption `json:"server_option,omitempty"`
}

func SetDefaultNotifyServerConfig(config *ServerOption) {
	defaultServerOption = config
}

var defaultServerOption = &ServerOption{
	ServerURL: "https://demo.com",
	Token:     "abcedefg",
}

// SendNotification sends a message to the Gotify server
// msg: the message to send
func SendNotification(ctx context.Context, msg *GotifyMessage) error {
	// Serialize payload to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	if msg.ServerOption == nil {
		msg.ServerOption = defaultServerOption
	}

	// Create request
	fullURL := fmt.Sprintf("%s/message?token=%s", msg.ServerOption.ServerURL, msg.ServerOption.Token)
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Configure HTTP client
	// Note: In production, use a proper certificate.
	// Here we skip verification for self-signed certs as requested in the context.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: tr,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gotify returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	logrus.Infof("Message sent successfully! Status: %s", resp.Status)
	return nil
}
