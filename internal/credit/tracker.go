package credit

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ReportForkStatus sends a non-blocking HTTP POST to analytics endpoint
func ReportForkStatus(botUsername string, creditIntact bool) {
	go func() {
		payload := map[string]interface{}{
			"bot_username":  botUsername,
			"credit_intact": creditIntact,
			"timestamp":     time.Now().Unix(),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return
		}

		// Mock analytics endpoint
		req, err := http.NewRequest("POST", "https://api.deepanshu.in/v1/analytics/fork", bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		_, _ = client.Do(req) // Ignoring response to avoid blocking
	}()
}

// CheckCreditPeriodically runs a goroutine that checks credit integrity periodically
func CheckCreditPeriodically(interval time.Duration, checkFunc func() bool) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if !checkFunc() {
				NotifyDeveloper("CREDIT TAMPERING DETECTED! Unauthorized modifications found in bot engine.")
			}
		}
	}()
}

// NotifyDeveloper sends alert
func NotifyDeveloper(message string) {
	// Just logs for now, could be extended to send Telegram message to developer
	log.Printf("[STD BOTS ALERTS] %s\n", message)
}
