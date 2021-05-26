package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// AMAlert definition
type AMAlert struct {
	Status       string            `json:"status,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	startsAt     time.Time         `json:"startsAt,omitempty"`
	endsAt       time.Time         `json:"endsAt,omitempty"`
	generatorURL string            `json:"generatorURL,omitempty"`
	fingerprint  string            `json:"fingerprint,omitempty"`
}

// AMPayload definition
type AMPayload struct {
	Receiver          string            `json:"receiver,omitempty"`
	Status            string            `json:"status,omitempty"`
	Alerts            []AMAlert         `json:"alerts,omitempty"`
	groupLabels       map[string]string `json:"groupLabels,omitempty"`
	commonLabels      map[string]string `json:"commonLabels,omitempty"`
	commonAnnotations map[string]string `json:"commonAnnotations,omitempty"`
	externalURL       string            `json:"externalURL,omitempty"`
	groupKey          string            `json:"groupKey,omitempty"`
}

// DiscordPayload definition
type DiscordPayload struct {
	Content string `json:"content,omitempty"`
}

var (
	webhookURL = os.Getenv("DISCORD_WEBHOOK")
	e          = echo.New()
)

func main() {

	e.POST("/", func(ctx echo.Context) error {
		payload := new(AMPayload)

		if err := ctx.Bind(payload); err != nil {
			log.Fatalf("%v", payload)
			return ctx.NoContent(http.StatusBadRequest)
		}

		for _, alert := range payload.Alerts {
			discordAlertNotification(alert)
		}

		return ctx.NoContent(http.StatusCreated)
	})

	e.Logger.Fatal(e.Start(":9094"))
}

func discordAlertNotification(alert AMAlert) error {
	labels := []string{
		"alertname", "instance", "severity", "monitor", "job",
	}

	if !hasKeys(labels, &alert.Labels) {
		return fmt.Errorf("No labels to derive alert content")
	}

	payload := DiscordPayload{
		Content: fmt.Sprintf(
			"[%s] %s @ %s by %s", alert.Labels["severity"], alert.Labels["alertname"],
			alert.Labels["instance"], alert.Labels["job"],
		),
	}

	b, err := json.Marshal(&payload)

	if err != nil {
		return err
	}
	res, err := http.Post(webhookURL, "application/json", strings.NewReader(string(b)))

	if err != nil || res.StatusCode != 200 {
		return fmt.Errorf("Unable to send notification")
	}

	return nil
}

func hasKeys(keys []string, m *map[string]string) bool {
	for _, k := range keys {
		if _, ok := (*m)[k]; !ok {
			return ok
		}
	}

	return true
}
