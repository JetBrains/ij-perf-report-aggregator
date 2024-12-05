package degradation_detector

import (
	"testing"
	"time"
)

func TestGetDateLink(t *testing.T) {
	mockNow := time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name        string
		degradation Degradation
		expected    string
	}{
		{
			name: "One week range",
			degradation: Degradation{
				timestamp: mockNow.Unix(),
			},
			expected: "timeRange=custom&customRange=2024-11-28:2024-12-5",
		},
		{
			name: "Different timestamp",
			degradation: Degradation{
				timestamp: mockNow.AddDate(0, 0, -2).Unix(), // 2 days before mockNow
			},
			expected: "timeRange=custom&customRange=2024-11-26:2024-12-5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getCustomDateLinkBetweenDates(tc.degradation, mockNow)
			if result != tc.expected {
				t.Errorf("getCustomDateLinkBetweenDates() = %v, want %v", result, tc.expected)
			}
		})
	}
}
