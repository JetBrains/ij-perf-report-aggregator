package degradation_detector

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type MissingData struct {
	LastBuild     string
	LastTimestamp int64
	TCBuildType   string
	Settings      Settings
}

const (
	DaysToCheckMissing     = -3
	StaleDataThresholdDays = -30 // Don't report if data has been missing longer than this
)

func InferMissingData(data <-chan QueryResultWithSettings) <-chan MissingData {
	output := make(chan MissingData, 100)
	go func() {
		for datum := range data {
			for _, missingData := range inferMissingData(datum.builds, datum.timestamps, datum.buildTypes, datum.Settings) {
				output <- missingData
			}
		}
		close(output)
	}()
	return output
}

func inferMissingData(builds []string, timestamps []int64, buildTypes []string, analysisSettings Settings) []MissingData {
	result := make([]MissingData, 0)

	threeDaysAgo := time.Now().AddDate(0, 0, DaysToCheckMissing).UnixMilli()

	monthAgo := time.Now().AddDate(0, 0, StaleDataThresholdDays).UnixMilli()

	lastTimestamp := timestamps[len(timestamps)-1]
	// if there the data is missing for the last 3 days but existed a month ago - report it
	if lastTimestamp < threeDaysAgo && lastTimestamp > monthAgo {
		result = append(result, MissingData{
			LastBuild:     builds[len(builds)-1],
			LastTimestamp: lastTimestamp,
			Settings:      analysisSettings,
			TCBuildType:   buildTypes[len(buildTypes)-1],
		})
	}

	return result
}

type metricsMerger interface {
	MergeMetrics(settings Settings) Settings
}

func mergeMetricsHelper(settings Settings, newSettings Settings) Settings {
	switch s := settings.(type) {
	case FleetStartupSettings:
		s.Metric = fmt.Sprintf("%s,%s", s.Metric, newSettings.GetMetric())
		return s
	case PerformanceSettings:
		s.Metric = fmt.Sprintf("%s,%s", s.Metric, newSettings.GetMetric())
		return s
	case StartupSettings:
		s.Metric = fmt.Sprintf("%s,%s", s.Metric, newSettings.GetMetric())
		return s
	default:
		return settings
	}
}

func (s FleetStartupSettings) MergeMetrics(settings Settings) Settings {
	return mergeMetricsHelper(s, settings)
}

func (s PerformanceSettings) MergeMetrics(settings Settings) Settings {
	return mergeMetricsHelper(s, settings)
}

func (s StartupSettings) MergeMetrics(settings Settings) Settings {
	return mergeMetricsHelper(s, settings)
}

type MissingDataByBuildType map[string]map[string]MissingData

func MergeMissingData(missingData <-chan MissingData) MissingDataByBuildType {
	// map of tc_build_type to project to missingData
	missingDataByBuildType := make(map[string]map[string][]MissingData)

	for datum := range missingData {
		if missingDataByBuildType[datum.TCBuildType] == nil {
			missingDataByBuildType[datum.TCBuildType] = make(map[string][]MissingData)
		}
		project := datum.Settings.GetProject()
		missingDataByBuildType[datum.TCBuildType][project] = append(
			missingDataByBuildType[datum.TCBuildType][project],
			datum,
		)
	}

	result := make(map[string]map[string]MissingData)
	for buildType, buildMissingData := range missingDataByBuildType {
		result[buildType] = make(map[string]MissingData)

		for project, projectMissingData := range buildMissingData {
			if len(projectMissingData) == 0 {
				continue
			}

			// Merge metrics within this project
			baseData := projectMissingData[0]
			for i := 1; i < len(projectMissingData); i++ {
				baseData.Settings = baseData.Settings.MergeMetrics(projectMissingData[i].Settings)
			}

			result[buildType][project] = baseData
		}
	}

	return result
}

func SendMissingDataMessages(data MissingDataByBuildType, client *http.Client) {
	// Messages grouped by Slack channel
	channelMessages := make(map[string][]string)

	// First, group all messages by Slack channel
	for _, projects := range data {
		for project, missingData := range projects {
			channel := missingData.Settings.SlackChannel()
			if channel == "" {
				continue // Skip if no channel is specified
			}

			// Format message for this entry
			readableDate := time.UnixMilli(missingData.LastTimestamp).Format("02-01-2006 15:04:05")
			message := fmt.Sprintf("Project: %s (%s)\nLast Date: %s\n<https://buildserver.labs.intellij.net/buildConfiguration/%s|TC Configuration>\n",
				project,
				missingData.Settings.GetMetric(),
				readableDate,
				missingData.TCBuildType)

			// Add to channel's message list
			channelMessages[channel] = append(channelMessages[channel], message)
		}
	}

	// Combine messages for each channel
	result := make(map[string]string)
	for channel, messages := range channelMessages {
		result[channel] = strings.Join(messages, "\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	for channel, message := range result {
		err := SendSlackMessage(ctx, client, SlackMessage{
			Text:    message,
			Channel: channel,
		})
		if err != nil {
			slog.Error("failed to send slack message", "error", err)
		}
	}
}
