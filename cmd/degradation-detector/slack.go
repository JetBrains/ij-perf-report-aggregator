package main

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"
)

func sendSlackMessage(ctx context.Context, degradation Degradation, analysisSettings AnalysisSettings) error {
  channel := "ij-perf-report-aggregator"
  slackMessage := createSlackMessage(channel, degradation, analysisSettings)
  slackMessageJson, err := json.Marshal(slackMessage)
  if err != nil {
    return fmt.Errorf("failed to marshal slack message: %w", err)
  }
  webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
  if len(webhookUrl) == 0 {
    return fmt.Errorf("SLACK_WEBHOOK_URL is not set")
  }
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

  log.Printf("Slack message was sent to " + channel)
  return nil
}

type SlackMessage struct {
  Text    string `json:"text"`
  Channel string `json:"channel"`
}

func createSlackMessage(channel string, degradation Degradation, settings AnalysisSettings) SlackMessage {
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
    Channel: channel,
  }
}
