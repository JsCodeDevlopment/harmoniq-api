package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func SendWebhook(url string, secret string, event string, data interface{}) error {
	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	sig := GenerateHMACSignature(jsonPayload, secret)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gost-Signature", sig)
	req.Header.Set("X-Gost-Event", event)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook delivery failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}

	return nil
}

func GenerateHMACSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
