package degradation_detector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMissingDataMessages(t *testing.T) {
	// Test case 1: Basic test with one project
	t.Parallel()
	t.Run("Single project test", func(t *testing.T) {
		t.Parallel()
		// Create a mock for MissingDataMerged
		mockData := createMockData(
			"channel1",
			"build_type1",
			"project1",
			createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
		)

		// Generate messages
		messages := GenerateMissingDataMessages(mockData)

		// Check channel
		message, exists := messages["channel1"]
		assert.True(t, exists, "Expected message for channel1")

		// Verify message content
		assert.Contains(t, message, "Data is missing for more than 3 days")
		assert.Contains(t, message, "*Projects:* project1")
		assert.Contains(t, message, "Metrics: metric1")
		assert.Contains(t, message, "Last Recorded: 15-05-2023")
		assert.Contains(t, message, "<https://buildserver.labs.intellij.net/buildConfiguration/build_type1|TC Configuration>")
		assert.Contains(t, message, "<https://ij-perf.labs.jb.gg/product/tests?mode=default&machine=&branch=&project=&measure=metric1&timeRange=custom&customRange=2025-04-12:2025-05-12|See charts>")
	})

	// Test case 2: Multiple projects with same metrics and timestamp
	t.Run("Multiple projects with same metrics", func(t *testing.T) {
		t.Parallel()
		mockData := MissingDataMerged{
			"channel1": {
				"build_type1": {
					"project1": createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
					"project2": createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
			},
		}

		messages := GenerateMissingDataMessages(mockData)
		message := messages["channel1"]

		assert.Contains(t, message, "*Projects:* project1, project2")
		assert.Contains(t, message, "Metrics: metric1")
	})

	// Test case 3: Multiple projects with different metrics
	t.Run("Multiple projects with different metrics", func(t *testing.T) {
		t.Parallel()
		mockData := MissingDataMerged{
			"channel1": {
				"build_type1": {
					"project1": createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
					"project2": createMockMissingData("metric2", time.Date(2023, 5, 16, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
			},
		}

		messages := GenerateMissingDataMessages(mockData)
		message := messages["channel1"]

		assert.Contains(t, message, "*Projects:* project1")
		assert.Contains(t, message, "*Projects:* project2")
		assert.Contains(t, message, "Metrics: metric1")
		assert.Contains(t, message, "Metrics: metric2")
	})

	// Test case 4: Multiple build types
	t.Run("Multiple build types", func(t *testing.T) {
		t.Parallel()
		mockData := MissingDataMerged{
			"channel1": {
				"build_type1": {
					"project1": createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
				"build_type2": {
					"project2": createMockMissingData("metric2", time.Date(2023, 5, 16, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
			},
		}

		messages := GenerateMissingDataMessages(mockData)
		message := messages["channel1"]

		assert.Contains(t, message, "build_type1")
		assert.Contains(t, message, "build_type2")
	})

	// Test case 5: Multiple channels
	t.Run("Multiple channels", func(t *testing.T) {
		t.Parallel()
		mockData := MissingDataMerged{
			"channel1": {
				"build_type1": {
					"project1": createMockMissingData("metric1", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
			},
			"channel2": {
				"build_type2": {
					"project2": createMockMissingData("metric2", time.Date(2023, 5, 16, 0, 0, 0, 0, time.UTC).UnixMilli()),
				},
			},
		}

		messages := GenerateMissingDataMessages(mockData)

		message1, exists1 := messages["channel1"]
		assert.True(t, exists1, "Expected message for channel1")
		assert.Contains(t, message1, "Data is missing for more than 3 days")

		message2, exists2 := messages["channel2"]
		assert.True(t, exists2, "Expected message for channel2")
		assert.Contains(t, message2, "Data is missing for more than 3 days")
	})
}

// Helper function to create mock data
func createMockData(channel, buildType, project string, missingData MissingData) MissingDataMerged {
	return MissingDataMerged{
		channel: {
			buildType: {
				project: missingData,
			},
		},
	}
}

// Helper function to create a mock MissingData object
func createMockMissingData(metric string, timestamp int64) MissingData {
	return MissingData{
		LastBuild:     "build1",
		LastTimestamp: timestamp,
		TCBuildType:   "tc_build_type1",
		Settings: PerformanceSettings{
			BaseSettings: BaseSettings{
				Metric: metric,
				SlackSettings: SlackSettings{
					Channel:     "test",
					ProductLink: "product",
				},
			},
		},
	}
}
