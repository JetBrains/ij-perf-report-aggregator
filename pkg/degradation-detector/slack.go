package degradation_detector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

type SlackMessage struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type slackSettings interface {
	CreateSlackMessage(degradations Degradation) SlackMessage
	SlackChannel() string
	ChartLink(timeRangeProvider TimeRangeProvider) string
}

type SlackSettings struct {
	Channel string
	// ProductLink is the part of the link on your dashboards after https://ij-perf.labs.jb.gg/.
	// For example: intellij,ijent, kmt, clion
	ProductLink string
	// Mention is an optional Slack mention prepended to PerformanceSettings messages for this channel.
	// Use the raw Slack syntax, e.g. "<@U01ABC2DEF>" for a user or "<!subteam^S01ABC2DEF>" for a user group.
	Mention string
}

func (s SlackSettings) SlackChannel() string {
	return s.Channel
}

// mentionPrefix returns the mention followed by a newline when a mention is configured, otherwise an empty string.
func (s SlackSettings) mentionPrefix() string {
	if s.Mention == "" {
		return ""
	}
	return s.Mention + "\n"
}

// eventLink repeats the test parameter for each test, because a test name can contain a comma.
func eventLink(tests []string, build string, timestamp int64) string {
	params := make([]string, 0, len(tests))
	for _, test := range tests {
		params = append(params, "test="+url.QueryEscape(test))
	}
	date := time.UnixMilli(timestamp).UTC().Format("02-01-2006")
	return fmt.Sprintf("<https://ij-perf.labs.jb.gg/degradations/report?%s&build=%s&date=%s|Report event>", strings.Join(params, "&"), url.QueryEscape(build), date)
}

func (s PerformanceSettings) CreateSlackMessage(d Degradation) SlackMessage {
	reason := getReason(d)
	date := time.UnixMilli(d.timestamp).UTC().Format("02-01-2006 15:04:05")
	link := s.ChartLink(d)
	tests := strings.Join(s.Projects(), "\n")
	mode := "default"
	if s.Mode != "" {
		mode = s.Mode
	}
	text := s.mentionPrefix() + fmt.Sprintf(
		"%sTest(s): %s\n"+
			"Metric: `%s`\n"+
			"Mode: %s\n"+
			"Build: %s\n"+
			"Branch: %s\n"+
			"Date: %s\n"+
			"Reason: %s\n"+
			"%s\n"+
			"%s", icon(d.medianValues), tests, s.Metric, mode, d.Build, s.Branch, date, reason, link, eventLink(s.Projects(), d.Build, d.timestamp))
	return SlackMessage{
		Text:    text,
		Channel: s.Channel,
	}
}

func getCustomRange(start, end time.Time) string {
	currentDate := fmt.Sprintf("%d-%02d-%d", end.Year(), end.Month(), end.Day())
	startDate := fmt.Sprintf("%d-%02d-%d", start.Year(), start.Month(), start.Day())
	return fmt.Sprintf("timeRange=custom&customRange=%s:%s", startDate, currentDate)
}

type TimeRangeProvider interface {
	GetRangeStartTime() time.Time
}

func (d Degradation) GetRangeStartTime() time.Time {
	t := time.Unix(d.timestamp/1000, 0)
	return t.AddDate(0, -1, 0)
}

func (s PerformanceSettings) ChartLink(d TimeRangeProvider) string {
	build := ""
	if d, ok := d.(Degradation); ok {
		build = d.Build
	}

	testPage := "tests"
	if strings.HasSuffix(s.Db, "Dev") {
		testPage = "testsDev"
	}
	machineGroup := getMachineGroup(s.Machine)
	projects := s.Projects()
	measurements := strings.Split(s.Metric, ",")
	escapedProjects := make([]string, 0, len(projects))
	escapedMeasurements := make([]string, 0, len(measurements))
	for _, p := range projects {
		escapedProjects = append(escapedProjects, url.QueryEscape(p))
	}
	for _, m := range measurements {
		escapedMeasurements = append(escapedMeasurements, url.QueryEscape(m))
	}
	mode := "default"
	if s.Mode != "" {
		mode = s.Mode
	}
	project := strings.Join(escapedProjects, "&project=")
	measure := strings.Join(escapedMeasurements, "&measure=")
	return fmt.Sprintf("<https://ij-perf.labs.jb.gg/%s/%s?mode=%s&machine=%s&branch=%s&project=%s&measure=%s&%s&point=%s|See charts>",
		s.ProductLink, testPage, mode, url.QueryEscape(machineGroup), url.QueryEscape(s.Branch), project, measure, getCustomRange(d.GetRangeStartTime(), time.Now()), build)
}

func SendDegradationsToSlack(insertionResults <-chan DegradationWithSettings, client *http.Client) {
	var wg sync.WaitGroup
	for result := range insertionResults {
		wg.Go(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			message := result.Settings.CreateSlackMessage(result.Details)
			if message.Channel == "" {
				slog.Info("slack channel is not set, message not sent", "message", message)
				return
			}
			err := SendSlackMessage(ctx, client, message)
			if err != nil {
				slog.Error("error while sending slack message", "error", err, "message", message)
				return
			}
			slog.Info("slack message was sent", "degradation", message)
		})
	}
	wg.Wait()
}

func SendSlackMessage(ctx context.Context, client *http.Client, slackMessage SlackMessage) error {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		return errors.New("SLACK_TOKEN is not set")
	}
	return sendSlackMessage(ctx, slack.New(slackToken, slack.OptionHTTPClient(client)), slackMessage)
}

const maxSlackAttempts = 5

// sendSlackMessage retries when Slack rate-limits the request, waiting as long as Slack asks.
// A run sends all its messages at once, so bursts to one channel exceed Slack's per-channel limit.
func sendSlackMessage(ctx context.Context, api *slack.Client, slackMessage SlackMessage) error {
	for attempt := 1; ; attempt++ {
		_, _, _, err := api.SendMessageContext(ctx, slackMessage.Channel, slack.MsgOptionText(slackMessage.Text, false))
		var rateLimited *slack.RateLimitedError
		if !errors.As(err, &rateLimited) || attempt == maxSlackAttempts {
			return err
		}
		// waiting past the deadline would only turn rate_limited into a less telling context error
		if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < rateLimited.RetryAfter {
			return err
		}
		slog.Info("slack rate limited, retrying", "channel", slackMessage.Channel, "retryAfter", rateLimited.RetryAfter, "attempt", attempt)
		select {
		case <-ctx.Done():
			return err
		case <-time.After(rateLimited.RetryAfter):
		}
	}
}

var machineGroupMap = map[string]string{
	"intellij-linux-performance-aws-%":   "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
	"intellij-linux-performance-aws%":    "Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)",
	"intellij-windows-performance-%":     "Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)",
	"intellij-windows-performance-aws-%": "Windows EC2 C6id.4xlarge or i4i.4xlarge (16 vCPU Xeon, 32 or 128 GB)",
	"intellij-linux-hw-hetzner%":         "linux-blade-hetzner",
	"intellij-linux-%-hetzner-%":         "linux-blade-hetzner",
	"ij-linux-x64-perf-hw-blade-%":       "linux-unit-perf-blade",
	"intellij-linux-hw-munit-%":          "Linux Munich i7-3770, 32 Gb",
	"intellij-windows-hw-munit-%":        "Windows Munich i7-3770, 32 Gb",
	"intellij-linux-hw-de-unit-%":        "Linux Munich i7-13700, 64 Gb",
	"intellij-windows-hw-de-%":           "Windows Munich i7-13700, 64 Gb",
	"intellij-macos-perf-eqx-%":          "Apple M4 Pro (13 vCPU, 64 GB)",
	"intellij-macos-hw-munit-%":          "macMini M1, 16 Gb",
}

func getMachineGroup(pattern string) string {
	return machineGroupMap[pattern]
}

// getReason says why the detector reported d: the step of the median, or the median of the latest values over a threshold.
func getReason(d Degradation) string {
	if d.step != nil {
		return getStepReason(d)
	}
	if d.threshold == nil {
		return getMessageBasedOnMedianChange(d.medianValues)
	}
	latest := "The latest value"
	if d.threshold.window > 1 {
		latest = fmt.Sprintf("The median of the latest %d values", d.threshold.window)
	}
	return fmt.Sprintf("Threshold exceeded. %s is %.2f, and the limit is %.2f. One run earlier it was %.2f.",
		latest, d.medianValues.newValue, d.threshold.value, d.medianValues.previousValue)
}

func getStepReason(d Degradation) string {
	change := "Degradation"
	if !d.IsDegradation {
		change = "Improvement"
	}
	ratio := ""
	if d.medianValues.previousValue != 0 {
		ratio = fmt.Sprintf(" (x%.2f)", d.medianValues.newValue/d.medianValues.previousValue)
	}
	if d.step.kind == StepDrift {
		return fmt.Sprintf("%s: slow drift. The median of the latest %d runs is %.2f, the median of the %d runs before them is %.2f%s, and the noise band is ±%.2f.",
			change, stepDriftWindow, d.medianValues.newValue, d.step.baseline, d.medianValues.previousValue, ratio, d.step.margin)
	}
	return fmt.Sprintf("%s in one run. The run is %.2f, the median of the %d runs before it is %.2f%s, and the noise band is ±%.2f. Changed between builds %s and %s.",
		change, d.medianValues.newValue, d.step.baseline, d.medianValues.previousValue, ratio, d.step.margin, d.PreviousBuild, d.Build)
}

func getMessageBasedOnMedianChange(medianValues CenterValues) string {
	percentageChange := medianValues.PercentageChange()
	medianMessage := fmt.Sprintf("Median changed by: %.2f%%. Median was %.2f and now it is %.2f.", percentageChange, medianValues.previousValue, medianValues.newValue)
	if medianValues.newValue > medianValues.previousValue {
		return "Degradation detected. " + medianMessage
	}
	return "Improvement detected. " + medianMessage
}

func icon(v CenterValues) string {
	var icon string
	if v.newValue > v.previousValue {
		icon = ":chart_with_upwards_trend:"
	} else {
		icon = ":chart_with_downwards_trend:"
	}
	return icon
}
