package degradation_detector

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/cenkalti/backoff/v4"
  "log/slog"
  "math"
  "net/http"
  "os"
)

func SendSlackMessage(ctx context.Context, client *http.Client, degradation DegradationWithContext) error {
  slackMessage := degradation.Settings.CreateSlackMessage(degradation.Details)
  slackMessageJson, err := json.Marshal(slackMessage)
  if err != nil {
    return fmt.Errorf("failed to marshal slack message: %w", err)
  }
  webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
  if len(webhookUrl) == 0 {
    return fmt.Errorf("SLACK_WEBHOOK_URL is not set")
  }
  err = backoff.Retry(func() error {
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookUrl, bytes.NewBuffer(slackMessageJson))
    if err != nil {
      return fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("sending slack message failed: %w", err)
    }
    defer resp.Body.Close()
    slog.Info("slack message was sent", "degradation", degradation.Settings)
    return nil
  }, backoff.NewExponentialBackOff())
  return err
}

type SlackMessage struct {
  Text    string `json:"text"`
  Channel string `json:"channel"`
}

func getMachineGroup(pattern string) string {
  switch pattern {
  case "intellij-linux-performance-aws-%":
    return "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  case "intellij-linux-hw-hetzner-%":
    return "linux-blade-hetzner"
  case "intellij-linux-hw-munit-%":
    return "Linux Munich i7-3770, 32 Gb"
  }
  return ""
}

func getMessageBasedOnMedianChange(medianValues MedianValues) string {
  percentageChange := math.Abs((medianValues.newValue - medianValues.previousValue) / medianValues.previousValue * 100)
  medianMessage := fmt.Sprintf("Median changed by: %.2f%%. Median was %.2f and now it is %.2f.", percentageChange, medianValues.previousValue, medianValues.newValue)
  if medianValues.newValue > medianValues.previousValue {
    return "Degradation detected. " + medianMessage
  }
  return "Improvement detected. " + medianMessage
}
