package main

import (
  "bytes"
  "context"
  "encoding/json"
  "fmt"
  "github.com/cenkalti/backoff/v4"
  "log"
  "net/http"
  "net/url"
  "os"
  "strings"
  "time"
)

func sendSlackMessage(ctx context.Context, degradation Degradation) error {
  analysisSettings := degradation.analysisSettings
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
  testPage := "tests"
  if strings.HasSuffix(degradation.analysisSettings.db, "Dev") {
    testPage = "testsDev"
  }
  machineGroup := getMachineGroup(settings.machine)
  link := fmt.Sprintf("https://ij-perf.labs.jb.gg/%s/%s?machine=%s&branch=%s&project=%s&measure=%s&timeRange=1M",
    settings.productLink, testPage, url.QueryEscape(machineGroup), url.QueryEscape(settings.branch), url.QueryEscape(settings.test), url.QueryEscape(settings.metric))

  icon := ""
  if degradation.medianValues.newValue > degradation.medianValues.previousValue {
    icon = ":chart_with_upwards_trend:"
  } else {
    icon = ":chart_with_downwards_trend:"
  }

  text := fmt.Sprintf(
    "%sTest: %s\n"+
      "Metric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s", icon, settings.test, settings.metric, degradation.build, degradation.analysisSettings.branch, date, reason, link)
  return SlackMessage{
    Text:    text,
    Channel: settings.channel,
  }
}

func getMachineGroup(pattern string) string {
  switch pattern {
  case "intellij-linux-performance-aws-%":
    return "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  case "intellij-linux-hw-hetzner-%":
    return "linux-blade-hetzner"
  }
  return ""
}
