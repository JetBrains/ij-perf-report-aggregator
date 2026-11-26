// Package alerts reads the notification settings of IJ Perf that JetBrains/ide-infrastructure-alerts publishes as
// ij-perf.json. The README of that repository describes the file and how a service applies it.
package alerts

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	service       = "ij-perf"
	formatVersion = 1
)

type Config struct {
	Version       string                  `json:"version"`
	Commit        string                  `json:"commit"`
	Notifications map[string]Notification `json:"notifications"`
}

type Notification struct {
	// Fallback is the channel of an event that no subscriber matched, empty when such an event is dropped.
	Fallback string `json:"fallback"`
	// Spec is the default settings of the notification. The spec of a subscriber or of a block is already merged
	// with it, so a reader never merges.
	Spec        json.RawMessage       `json:"spec"`
	Subscribers map[string]Subscriber `json:"subscribers"`
}

type Subscriber struct {
	Enabled bool            `json:"enabled"`
	Target  Target          `json:"target"`
	Spec    json.RawMessage `json:"spec"`
	Owners  []Owner         `json:"owners"`
	Match   []Block         `json:"match"`
}

// Owner is a code-owner group whose events the subscriber receives.
type Owner struct {
	Name      string `json:"name"`
	Subgroups bool   `json:"subgroups"`
}

// Block selects events by their fields. It matches an event when every field of the block matches.
type Block struct {
	Fields map[string][]string `json:"fields"`
	Target Target              `json:"target"`
	Spec   json.RawMessage     `json:"spec"`
}

// Target is a Slack channel or a YouTrack project. IJ Perf sends only to Slack.
type Target struct {
	Channel  string          `json:"channel"`
	Mention  []string        `json:"mention"`
	YouTrack json.RawMessage `json:"youtrack"`
}

// SlackChannel returns the channel name without the leading #, empty when the target is not a Slack channel.
func (t Target) SlackChannel() string {
	return strings.TrimPrefix(t.Channel, "#")
}

// Parse reads ij-perf.json. Unknown fields are ignored: the format gets new fields without a version change.
func Parse(data []byte) (*Config, error) {
	var file struct {
		Config

		FormatVersion int    `json:"formatVersion"`
		Service       string `json:"service"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("failed to parse the alerts config: %w", err)
	}
	if file.FormatVersion != formatVersion {
		return nil, fmt.Errorf("unsupported format version %d of the alerts config, expected %d", file.FormatVersion, formatVersion)
	}
	if file.Service != service {
		return nil, fmt.Errorf("the alerts config is for the service %q, expected %q", file.Service, service)
	}
	if len(file.Notifications) == 0 {
		return nil, errors.New("the alerts config has no notifications")
	}
	return &file.Config, nil
}

// SubscriberNames returns the names of the subscribers in a stable order.
func (n Notification) SubscriberNames() []string {
	names := make([]string, 0, len(n.Subscribers))
	for name := range n.Subscribers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Matches reports whether one of the patterns matches the value. In a pattern `*` is any sequence of characters, and
// nothing else is special.
func Matches(patterns []string, value string) bool {
	for _, pattern := range patterns {
		if matchGlob(pattern, value) {
			return true
		}
	}
	return false
}

func matchGlob(pattern, value string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == value
	}
	if !strings.HasPrefix(value, parts[0]) {
		return false
	}
	value = value[len(parts[0]):]
	last := parts[len(parts)-1]
	for _, part := range parts[1 : len(parts)-1] {
		index := strings.Index(value, part)
		if index < 0 {
			return false
		}
		value = value[index+len(part):]
	}
	return strings.HasSuffix(value, last)
}

// IsPattern reports whether the value of a block field is a pattern and not a literal name.
func IsPattern(value string) bool {
	return strings.Contains(value, "*")
}
