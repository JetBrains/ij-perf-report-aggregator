package alerts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const downloaded = `{"formatVersion": 1, "service": "ij-perf", "version": "v100", "notifications": {"n": {}}}`

// releaseServer is GitHub with one release: /latest lists its assets, and /asset is ij-perf.json.
func releaseServer(t *testing.T, withAsset bool) *httptest.Server {
	t.Helper()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/latest":
			assets := `{"name": "manifest.json", "url": "` + server.URL + `/manifest"}`
			if withAsset {
				assets += `, {"name": "ij-perf.json", "url": "` + server.URL + `/asset"}`
			}
			_, _ = w.Write([]byte(`{"assets": [` + assets + `]}`))
		case "/asset":
			assert.Equal(t, "application/octet-stream", r.Header.Get("Accept"))
			_, _ = w.Write([]byte(downloaded))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestDownload(t *testing.T) {
	t.Parallel()
	t.Run("the asset of the latest release", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, true)
		config, err := download(t.Context(), server.Client(), server.URL+"/latest", "secret")
		require.NoError(t, err)
		assert.Equal(t, "v100", config.Version)
	})
	t.Run("no token", func(t *testing.T) {
		t.Parallel()
		_, err := download(t.Context(), http.DefaultClient, "http://unused.invalid", "")
		assert.ErrorContains(t, err, GithubTokenEnv)
	})
	t.Run("wrong token", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, false)
		_, err := download(t.Context(), server.Client(), server.URL+"/latest", "wrong")
		assert.ErrorContains(t, err, "401")
	})
	t.Run("release without the asset", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, false)
		_, err := download(t.Context(), server.Client(), server.URL+"/latest", "secret")
		assert.ErrorContains(t, err, "has no ij-perf.json")
	})
}

func TestLoadUsesBuiltInConfigWithoutToken(t *testing.T) {
	t.Setenv(GithubTokenEnv, "")
	t.Setenv(ConfigFileEnv, "")
	config, err := Load(t.Context(), http.DefaultClient)
	require.Error(t, err)
	require.NotNil(t, config)
	builtIn, parseErr := Parse(snapshot)
	require.NoError(t, parseErr)
	assert.Equal(t, builtIn.Version, config.Version)
}

func TestLoadUsesLocalFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ij-perf.json")
	require.NoError(t, os.WriteFile(path, []byte(downloaded), 0o600))
	t.Setenv(ConfigFileEnv, path)
	config, err := Load(t.Context(), http.DefaultClient)
	require.NoError(t, err)
	assert.Equal(t, "v100", config.Version)
}
