package degradation_detector

import (
  "context"
  "errors"
  "fmt"
  "github.com/slack-go/slack"
  "log/slog"
  "net/http"
  "net/url"
  "os"
  "strings"
  "sync"
  "time"
)

type SlackMessage struct {
  Text    string `json:"text"`
  Channel string `json:"channel"`
}

type slackSettings interface {
  CreateSlackMessage(degradations Degradation) SlackMessage
  SlackChannel() string
}

type SlackSettings struct {
  Channel     string
  ProductLink string
}

func (s SlackSettings) SlackChannel() string {
  return s.Channel
}

func eventLink(tests string, build string, timestamp int64) string {
  projects := strings.Split(tests, ",")
  escapedProjects := make([]string, 0, len(projects))
  for _, p := range projects {
    escapedProjects = append(escapedProjects, url.QueryEscape(p))
  }
  date := time.UnixMilli(timestamp).UTC().Format("02-01-2006")
  project := strings.Join(escapedProjects, ",")
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/degradations/report?tests=%s&build=%s&date=%s", project, build, date)
}

func (s PerformanceSettings) CreateSlackMessage(d Degradation) SlackMessage {
  reason := getMessageBasedOnMedianChange(d.medianValues)
  date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
  link := s.slackLink()
  tests := strings.ReplaceAll(s.Project, ",", "\n")
  text := fmt.Sprintf(
    "%sTest(s): %s\n"+
      "Metric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s\n"+
      "Report event: %s", icon(d.medianValues), tests, s.Metric, d.Build, s.Branch, date, reason, link, eventLink(s.Project, d.Build, d.timestamp))
  return SlackMessage{
    Text:    text,
    Channel: s.Channel,
  }
}

func (s StartupSettings) CreateSlackMessage(d Degradation) SlackMessage {
  reason := getMessageBasedOnMedianChange(d.medianValues)
  date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
  link := s.slackLink()
  tests := strings.ReplaceAll(s.Project, ",", "\n")

  text := fmt.Sprintf(
    "%sProject(s): %s\n"+
      "Metric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s\n"+
      "Report event: %s", icon(d.medianValues), tests, s.Metric, d.Build, s.Branch, date, reason, link, eventLink(s.Project, d.Build, d.timestamp))
  return SlackMessage{
    Text:    text,
    Channel: s.Channel,
  }
}

func (s FleetStartupSettings) CreateSlackMessage(d Degradation) SlackMessage {
  reason := getMessageBasedOnMedianChange(d.medianValues)
  date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
  link := s.slackLink()

  text := fmt.Sprintf(
    "%sMetric: %s\n"+
      "Build: %s\n"+
      "Branch: %s\n"+
      "Date: %s\n"+
      "Reason: %s\n"+
      "Link: %s\n", icon(d.medianValues), s.Metric, d.Build, s.Branch, date, reason, link)
  return SlackMessage{
    Text:    text,
    Channel: s.Channel,
  }
}

func (s FleetStartupSettings) slackLink() string {
  machineGroup := getMachineGroup(s.Machine)
  measurements := strings.Split(s.Metric, ",")
  escapedMeasurements := make([]string, 0, len(measurements))
  for _, m := range measurements {
    escapedMeasurements = append(escapedMeasurements, url.QueryEscape(m))
  }
  measure := strings.Join(escapedMeasurements, "&measure=")
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/fleet/startupExplore?machine=%s&branch=%s&measure=%s&timeRange=1M",
    url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), measure)
}

func (s PerformanceSettings) slackLink() string {
  testPage := "tests"
  if strings.HasSuffix(s.Db, "Dev") {
    testPage = "testsDev"
  }
  machineGroup := getMachineGroup(s.Machine)
  projects := strings.Split(s.Project, ",")
  measurements := strings.Split(s.Metric, ",")
  escapedProjects := make([]string, 0, len(projects))
  escapedMeasurements := make([]string, 0, len(measurements))
  for _, p := range projects {
    escapedProjects = append(escapedProjects, url.QueryEscape(p))
  }
  for _, m := range measurements {
    escapedMeasurements = append(escapedMeasurements, url.QueryEscape(m))
  }
  project := strings.Join(escapedProjects, "&project=")
  measure := strings.Join(escapedMeasurements, "&measure=")
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/%s/%s?machine=%s&branch=%s&project=%s&measure=%s&timeRange=1M",
    s.ProductLink, testPage, url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), project, measure)
}

func (s StartupSettings) slackLink() string {
  machineGroup := getMachineGroup(s.Machine)
  return fmt.Sprintf("https://ij-perf.labs.jb.gg/%s/startup?machine=%s&branch=%s&product=%s&project=%s&timeRange=1M",
    s.ProductLink, url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), url.QueryEscape(s.Product), url.QueryEscape(s.Project))
}

func SendDegradationsToSlack(insertionResults <-chan DegradationWithSettings, client *http.Client) {
  var wg sync.WaitGroup
  for result := range insertionResults {
    wg.Add(1)
    go func() {
      defer wg.Done()
      ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
      defer cancel()
      message := result.Settings.CreateSlackMessage(result.Details)
      err := sendSlackMessage(ctx, client, message)
      if err != nil {
        slog.Error("error while sending slack message", "error", err, "message", message)
        return
      }
      slog.Info("slack message was sent", "degradation", message)
    }()
  }
  wg.Wait()
}

func sendSlackMessage(ctx context.Context, client *http.Client, slackMessage SlackMessage) error {
  slackToken := os.Getenv("SLACK_TOKEN")
  if slackToken == "" {
    return errors.New("SLACK_TOKEN is not set")
  }
  api := slack.New(slackToken, slack.OptionHTTPClient(client))
  _, _, _, err := api.SendMessageContext(ctx, slackMessage.Channel, slack.MsgOptionText(slackMessage.Text, false))
  return err
}

func getMachineGroup(pattern string) string {
  switch pattern {
  case "intellij-linux-performance-aws-%":
    return "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  case "intellij-linux-hw-hetzner%":
    return "linux-blade-hetzner"
  case "intellij-linux-%-hetzner-%":
    return "linux-blade-hetzner"
  case "intellij-linux-hw-munit-%":
    return "Linux Munich i7-3770, 32 Gb"
  case "intellij-windows-hw-munit-%":
    return "Windows Munich i7-3770, 32 Gb"
  case "intellij-macos-hw-munit-%":
    return "macMini M1, 16 Gb"
  }
  return ""
}

func getMessageBasedOnMedianChange(medianValues MedianValues) string {
  percentageChange := medianValues.PercentageChange()
  medianMessage := fmt.Sprintf("Median changed by: %.2f%%. Median was %.2f and now it is %.2f.", percentageChange, medianValues.previousValue, medianValues.newValue)
  if medianValues.newValue > medianValues.previousValue {
    return "Degradation detected. " + medianMessage
  }
  return "Improvement detected. " + medianMessage
}

func icon(v MedianValues) string {
  icon := ""
  if v.newValue > v.previousValue {
    icon = ":chart_with_upwards_trend:"
  } else {
    icon = ":chart_with_downwards_trend:"
  }
  return icon
}
