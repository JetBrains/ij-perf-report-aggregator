package degradation_detector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// Hook posts the new degradations of the matching tests to a webhook, besides the Slack message.
//
// A hook matches the PerformanceSettings of its Db and Table whose project starts with ProjectPrefix. The URL and the
// token of the webhook come from the environment variables URLEnv and TokenEnv, so no secret is in the code. A hook
// whose variables are not set sends nothing. A degradation reaches a hook once: the detector sends only the
// degradations that it stored as new accidents.
type Hook struct {
	// Name identifies the hook in the logs.
	Name          string
	Db            string
	Table         string
	ProjectPrefix string
	URLEnv        string
	TokenEnv      string
	// AuthScheme is the scheme of the Authorization header: "Bearer" when empty, "ApiKey" for Air Automation.
	AuthScheme string
}

// HookEvent is the JSON body of a hook request.
type HookEvent struct {
	// EventID is the same for every delivery of one degradation, so the receiver can drop a repeated one.
	EventID       string  `json:"eventId"`
	Source        string  `json:"source"`
	Kind          string  `json:"kind"`
	Db            string  `json:"db"`
	Table         string  `json:"table"`
	Project       string  `json:"project"`
	Metric        string  `json:"metric"`
	Mode          string  `json:"mode"`
	Branch        string  `json:"branch"`
	Machine       string  `json:"machine"`
	Build         string  `json:"build"`
	Date          string  `json:"date"`
	PreviousValue float64 `json:"previousValue"`
	NewValue      float64 `json:"newValue"`
	ChangePercent float64 `json:"changePercent"`
	Reason        string  `json:"reason"`
	ChartURL      string  `json:"chartUrl"`
	ReportURL     string  `json:"reportUrl"`
}

const (
	maxHookAttempts = 4
	hookRetryPause  = 5 * time.Second
	hookTimeout     = 2 * time.Minute
	maxHookRequests = 5
)

func (h Hook) matches(s PerformanceSettings) bool {
	if s.Db != h.Db || s.Table != h.Table {
		return false
	}
	// MergeDegradations joins the projects of one metric and build with commas.
	for project := range strings.SplitSeq(s.Project, ",") {
		if strings.HasPrefix(project, h.ProjectPrefix) {
			return true
		}
	}
	return false
}

func (h Hook) authorization(token string) string {
	scheme := h.AuthScheme
	if scheme == "" {
		scheme = "Bearer"
	}
	return scheme + " " + token
}

func newHookEvent(s PerformanceSettings, d Degradation) HookEvent {
	kind := "improvement"
	if d.IsDegradation {
		kind = "degradation"
	}
	mode := "default"
	if s.Mode != "" {
		mode = s.Mode
	}
	return HookEvent{
		EventID:       s.DBTestName() + ":" + d.Build,
		Source:        "ij-perf",
		Kind:          kind,
		Db:            s.Db,
		Table:         s.Table,
		Project:       s.Project,
		Metric:        s.Metric,
		Mode:          mode,
		Branch:        s.Branch,
		Machine:       s.Machine,
		Build:         d.Build,
		Date:          time.UnixMilli(d.timestamp).UTC().Format(time.RFC3339),
		PreviousValue: d.medianValues.previousValue,
		NewValue:      d.medianValues.newValue,
		ChangePercent: d.medianValues.PercentageChange(),
		Reason:        getMessageBasedOnMedianChange(d.medianValues),
		ChartURL:      slackLinkURL(s.ChartLink(d)),
		ReportURL:     slackLinkURL(eventLink(s.Project, d.Build, d.timestamp)),
	}
}

// slackLinkURL returns the URL of a Slack link "<url|text>", or the text itself when it is no Slack link.
func slackLinkURL(link string) string {
	if !strings.HasPrefix(link, "<") || !strings.HasSuffix(link, ">") {
		return link
	}
	url, _, _ := strings.Cut(strings.TrimSuffix(strings.TrimPrefix(link, "<"), ">"), "|")
	return url
}

// SendDegradationsToHooks posts each degradation to every configured hook that matches its settings. A delivery that
// fails after its retries is logged and dropped: the accident is stored already, so a later run does not send it again.
func SendDegradationsToHooks(degradations <-chan DegradationWithSettings, client *http.Client, hooks []Hook) {
	type target struct {
		hook  Hook
		url   string
		token string
	}
	targets := make([]target, 0, len(hooks))
	for _, hook := range hooks {
		url, token := os.Getenv(hook.URLEnv), os.Getenv(hook.TokenEnv)
		if url == "" || token == "" {
			slog.Info("hook is not configured, it sends nothing", "hook", hook.Name, "urlEnv", hook.URLEnv, "tokenEnv", hook.TokenEnv)
			continue
		}
		targets = append(targets, target{hook, url, token})
	}

	var group errgroup.Group
	group.SetLimit(maxHookRequests)
	for result := range degradations {
		settings, ok := result.Settings.(PerformanceSettings)
		if !ok {
			continue
		}
		for _, t := range targets {
			if !t.hook.matches(settings) {
				continue
			}
			event := newHookEvent(settings, result.Details)
			group.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), hookTimeout)
				defer cancel()
				if err := sendHookEvent(ctx, client, t.url, t.hook.authorization(t.token), event, hookRetryPause); err != nil {
					slog.Error("hook did not accept the event", "hook", t.hook.Name, "event", event.EventID, "error", err)
					return nil
				}
				slog.Info("hook event was sent", "hook", t.hook.Name, "event", event.EventID)
				return nil
			})
		}
	}
	_ = group.Wait()
}

// sendHookEvent retries a network error and a 5xx or 429 answer after retryPause times the attempt, and fails at once
// on another answer that is not 2xx.
func sendHookEvent(ctx context.Context, client *http.Client, url string, authorization string, event HookEvent, retryPause time.Duration) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal the hook event: %w", err)
	}
	for attempt := 1; ; attempt++ {
		retry, err := postHookEvent(ctx, client, url, authorization, body)
		if err == nil {
			return nil
		}
		if !retry || attempt == maxHookAttempts {
			return err
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(time.Duration(attempt) * retryPause):
		}
	}
}

func postHookEvent(ctx context.Context, client *http.Client, url string, authorization string, body []byte) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("failed to create the request: %w", err)
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return true, fmt.Errorf("failed to send the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return false, nil
	}
	return resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests, fmt.Errorf("the hook answered %v", resp.Status)
}
