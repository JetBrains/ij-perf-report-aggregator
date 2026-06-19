package meta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serveJSON(body string) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(body))
	}
}

func newTestCodeOwnersClient(t *testing.T, handler http.HandlerFunc) *CodeOwnersClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &CodeOwnersClient{
		baseURL: srv.URL,
		token:   "test-token",
		http:    srv.Client(),
	}
}

func TestSlackChannelOf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		attrs []coAttribute
		want  string
	}{
		{
			name:  "alerts channel preferred over the general channel",
			attrs: []coAttribute{{Key: "Slack Channel", Value: "general"}, {Key: "Slack Alerts Channel", Value: "alerts"}},
			want:  "alerts",
		},
		{
			name:  "general channel used when no alerts channel is set",
			attrs: []coAttribute{{Key: "Slack Channel", Value: "general"}},
			want:  "general",
		},
		{
			name:  "leading hash and surrounding whitespace are trimmed",
			attrs: []coAttribute{{Key: "Slack Channel", Value: "  #general  "}},
			want:  "general",
		},
		{
			name:  "empty alerts channel falls through to the general channel",
			attrs: []coAttribute{{Key: "Slack Alerts Channel", Value: ""}, {Key: "Slack Channel", Value: "general"}},
			want:  "general",
		},
		{
			name:  "hash-only value is treated as empty",
			attrs: []coAttribute{{Key: "Slack Channel", Value: "#"}},
			want:  "",
		},
		{
			name:  "no slack attribute yields no channel",
			attrs: []coAttribute{{Key: "Some Other Attribute", Value: "x"}},
			want:  "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, slackChannelOf(tc.attrs))
		})
	}
}

func TestFetchOwnerChannels_RequestContract(t *testing.T) {
	t.Parallel()
	var gotFields, gotLimit, gotAuth string
	client := newTestCodeOwnersClient(t, func(writer http.ResponseWriter, request *http.Request) {
		gotFields = request.URL.Query().Get("fields")
		gotLimit = request.URL.Query().Get("limit")
		gotAuth = request.Header.Get("Authorization")
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"items":[]}`))
	})

	_, err := client.FetchOwnerChannels(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "aliases,attributes", gotFields)
	assert.Equal(t, "-1", gotLimit)
	assert.Equal(t, "Bearer test-token", gotAuth)
}

func TestFetchOwnerChannels_MapsNamesAndAliases(t *testing.T) {
	t.Parallel()
	// Channel-value normalization (alerts-vs-general preference, '#'/whitespace trimming) is
	// covered by TestSlackChannelOf; here the values are plain so this test focuses solely on
	// name/alias mapping and omission of channel-less groups.
	body := `{"items":[
		{"name":"java","aliases":[{"name":"java-legacy"}],"attributes":[{"key":"Slack Channel","value":"idea-java-alerts"}]},
		{"name":"ruby","aliases":[],"attributes":[{"key":"Slack Channel","value":"rubymine-alerts"}]},
		{"name":"no-channel","aliases":[{"name":"nc-alias"}],"attributes":[]}
	]}`
	client := newTestCodeOwnersClient(t, serveJSON(body))

	channels, err := client.FetchOwnerChannels(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "idea-java-alerts", channels["java"])
	assert.Equal(t, "idea-java-alerts", channels["java-legacy"], "alias resolves to the same channel as its canonical name")
	assert.Equal(t, "rubymine-alerts", channels["ruby"], "a second group also maps to its channel")

	_, ok := channels["no-channel"]
	assert.False(t, ok, "a group without a slack channel is omitted")
	_, ok = channels["nc-alias"]
	assert.False(t, ok, "the alias of a channel-less group is also omitted")
}

func TestFetchOwnerChannels_CanonicalNameWinsOverStaleAlias(t *testing.T) {
	t.Parallel()
	// Group A was renamed "payments" -> "payments-platform" and still carries "payments" as a
	// stale alias. Group B is a new active group that took the freed name "payments".
	// The active group's canonical name must win regardless of the service's response order.
	groupA := `{"name":"payments-platform","aliases":[{"name":"payments"}],"attributes":[{"key":"Slack Channel","value":"pay-platform-alerts"}]}`
	groupB := `{"name":"payments","aliases":[],"attributes":[{"key":"Slack Channel","value":"new-payments-alerts"}]}`

	orders := []struct {
		name string
		body string
	}{
		{name: "A before B", body: `{"items":[` + groupA + `,` + groupB + `]}`},
		{name: "B before A", body: `{"items":[` + groupB + `,` + groupA + `]}`},
	}
	for _, tc := range orders {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := newTestCodeOwnersClient(t, serveJSON(tc.body))
			channels, err := client.FetchOwnerChannels(context.Background())
			require.NoError(t, err)
			assert.Equal(t, "new-payments-alerts", channels["payments"], "active group's canonical name wins over a stale alias")
			assert.Equal(t, "pay-platform-alerts", channels["payments-platform"])
		})
	}
}

func TestFetchOwnerChannels_NoTokenReturnsError(t *testing.T) {
	t.Parallel()
	client := &CodeOwnersClient{baseURL: "http://example.invalid", token: "", http: &http.Client{}}
	_, err := client.FetchOwnerChannels(context.Background())
	require.Error(t, err)
}

func TestFetchOwnerChannels_NonOKStatusReturnsError(t *testing.T) {
	t.Parallel()
	client := newTestCodeOwnersClient(t, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
	_, err := client.FetchOwnerChannels(context.Background())
	require.Error(t, err)
}
