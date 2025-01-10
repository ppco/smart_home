package switchbot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type client struct {
	client http.Client
	token  string
	secret string
}

const (
	commandURL = "https://api.switch-bot.com/v1.1/devices/%s/commands"
)

func NewClient() *client {
	return &client{
		token:  os.Getenv("SWITCHBOT_TOKEN"),
		secret: os.Getenv("SWITCHBOT_SECRET"),
		client: http.Client{
			Transport: http.DefaultTransport,
			Timeout:   30 * time.Second,
		},
	}
}

// https://github.com/OpenWonderLabs/SwitchBotAPI?tab=readme-ov-file#send-device-control-commands
func (c *client) ExecCommand(deviceID string, param Commands) error {
	reqBody := map[string]any{
		"command":     param.Command,
		"parameter":   param.Parameter,
		"commandType": param.CommandType,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", fmt.Sprintf(commandURL, deviceID), bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	c.addHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body:%s", resp.StatusCode, string(resBody))
	}

	return nil
}

func (c *client) addHeaders(req *http.Request) {
	// https://github.com/OpenWonderLabs/SwitchBotAPI?tab=readme-ov-file#javascript-example-code
	t := time.Now().UnixMilli()
	token := c.token
	secret := c.secret
	nonce := "requestID"
	// HMAC-SHA256署名を作成
	data := fmt.Sprintf("%s%d%s", token, t, nonce)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("sign", sign)
	req.Header.Set("nonce", "requestID")
	req.Header.Set("t", strconv.FormatInt(t, 10))
}
