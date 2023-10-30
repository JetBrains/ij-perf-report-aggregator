package main

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/cenkalti/backoff/v4"
  "log"
  "net/http"
  "os"
  "time"
)

func sendSlackMessage(ctx context.Context, degradation Degradation, analysisSettings AnalysisSettings) error {
  slackMessage := createSlackMessage(degradation, analysisSettings)
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
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("sending slack message failed: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("Slack message was sent to " + analysisSettings.channel)
    return nil
  }, backoff.NewExponentialBackOff())
  return err
}

type SlackMessage struct {
  Text    string `json:"text"`
  Channel string `json:"channel"`
}

func createSlackMessage(degradation Degradation, settings AnalysisSettings) SlackMessage {
  reason := getMessageBasedOnMedianChange(degradation.medianValues)
  date := time.UnixMilli(degradation.timestamp).UTC().Format("02-01-2006 15:04:05")
  text := fmt.Sprintf("DB: %s\n"+
    "Table: %s\n"+
    "Build: %s\n"+
    "Date: %s\n"+
    "Test/metric: %s/%s\n"+
    "Reason: %s", settings.db, settings.table, degradation.build, date, settings.test, settings.metric, reason)
  return SlackMessage{
    Text:    text,
    Channel: settings.channel,
  }
}
