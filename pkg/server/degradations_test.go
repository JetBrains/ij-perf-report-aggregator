package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDegradationsRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		method          string
		query           string
		body            string
		wantOwners      []string
		wantMetrics     []string
		wantInferred    bool
		wantErrContains string
	}{
		{
			name:         "GET query only",
			method:       http.MethodGet,
			query:        "owner=Java&owner=&owner=Kotlin&from=2024-01-01&to=2024-03-31&inferredEvents=true&additionalMetric=m1&additionalMetric=",
			wantOwners:   []string{"Java", "Kotlin"},
			wantMetrics:  []string{"m1"},
			wantInferred: true,
		},
		{
			name:       "GET ignores body",
			method:     http.MethodGet,
			query:      "owner=Java&from=2024-01-01&to=2024-03-31",
			body:       `{"owners": ["Kotlin"]}`,
			wantOwners: []string{"Java"},
		},
		{
			name:        "POST body only",
			method:      http.MethodPost,
			query:       "from=2024-01-01&to=2024-03-31",
			body:        `{"owners": ["Java", "", "Kotlin"], "additionalMetrics": ["m1", ""]}`,
			wantOwners:  []string{"Java", "Kotlin"},
			wantMetrics: []string{"m1"},
		},
		{
			name:       "POST query only with empty body",
			method:     http.MethodPost,
			query:      "owner=Java&from=2024-01-01&to=2024-03-31",
			wantOwners: []string{"Java"},
		},
		{
			name:         "POST union of query and body",
			method:       http.MethodPost,
			query:        "owner=Java&from=2024-01-01&to=2024-03-31&inferredEvents=true&additionalMetric=m1",
			body:         `{"owners": ["Kotlin"], "additionalMetrics": ["m2"]}`,
			wantOwners:   []string{"Java", "Kotlin"},
			wantMetrics:  []string{"m1", "m2"},
			wantInferred: true,
		},
		{
			name:            "POST malformed JSON",
			method:          http.MethodPost,
			query:           "owner=Java&from=2024-01-01&to=2024-03-31",
			body:            `{"owners": [`,
			wantErrContains: "invalid request body",
		},
		{
			name:            "POST empty body and no owners",
			method:          http.MethodPost,
			query:           "from=2024-01-01&to=2024-03-31",
			wantErrContains: "are all required",
		},
		{
			name:            "POST blank owners only",
			method:          http.MethodPost,
			query:           "owner=&from=2024-01-01&to=2024-03-31",
			body:            `{"owners": [""]}`,
			wantErrContains: "are all required",
		},
		{
			name:            "POST missing from",
			method:          http.MethodPost,
			query:           "to=2024-03-31",
			body:            `{"owners": ["Java"]}`,
			wantErrContains: "are all required",
		},
		{
			name:            "POST missing to",
			method:          http.MethodPost,
			query:           "from=2024-01-01",
			body:            `{"owners": ["Java"]}`,
			wantErrContains: "are all required",
		},
		{
			name:            "GET missing from and to",
			method:          http.MethodGet,
			query:           "owner=Java",
			wantErrContains: "are all required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(tt.method, "/api/meta/degradations?"+tt.query, strings.NewReader(tt.body))

			params, err := parseDegradationsRequest(r)

			if tt.wantErrContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantOwners, params.owners)
			assert.Equal(t, tt.wantMetrics, params.additionalMetrics)
			assert.Equal(t, tt.wantInferred, params.includeInferred)
			assert.Equal(t, "2024-01-01", params.from)
			assert.Equal(t, "2024-03-31", params.to)
		})
	}
}

func TestDegradationsHandlerRejectsBadRequestsWithoutDb(t *testing.T) {
	t.Parallel()

	handler := CreateGetDegradationsHandler(nil)
	for _, tc := range []struct{ method, query, body string }{
		{http.MethodGet, "owner=Java", ""},
		{http.MethodPost, "from=2024-01-01&to=2024-03-31", ""},
		{http.MethodPost, "from=2024-01-01&to=2024-03-31", "not json"},
	} {
		r := httptest.NewRequest(tc.method, "/api/meta/degradations?"+tc.query, strings.NewReader(tc.body))
		w := httptest.NewRecorder()

		handler(w, r)

		assert.Equalf(t, http.StatusBadRequest, w.Code, "%s ?%s body=%q", tc.method, tc.query, tc.body)
	}
}
