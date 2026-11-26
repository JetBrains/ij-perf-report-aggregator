package alerts

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// ConfigFileEnv names a local ij-perf.json to use instead of the published one, for a local run.
	ConfigFileEnv = "ALERTS_CONFIG_FILE"
	// GithubTokenEnv names a GitHub token that can read the releases of the alerts repository.
	GithubTokenEnv = "ALERTS_GITHUB_TOKEN"

	latestReleaseURL = "https://api.github.com/repos/JetBrains/ide-infrastructure-alerts/releases/latest"
	assetName        = "ij-perf.json"
)

// The config that a run uses when it cannot download the latest one. Refresh it with `make update-alerts-snapshot`.
//
//go:embed snapshot/ij-perf.json
var snapshot []byte

// Load returns the config of the latest release of the alerts repository.
//
// A run must not depend on GitHub, so when the download fails Load returns the config that is built into the binary
// together with the error. The config is nil only when the built-in one is broken too.
func Load(ctx context.Context, client *http.Client) (*Config, error) {
	if path := os.Getenv(ConfigFileEnv); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return Parse(data)
	}
	config, err := download(ctx, client, latestReleaseURL, os.Getenv(GithubTokenEnv))
	if err == nil {
		return config, nil
	}
	builtIn, builtInErr := Parse(snapshot)
	if builtInErr != nil {
		return nil, errors.Join(err, builtInErr)
	}
	return builtIn, err
}

func download(ctx context.Context, client *http.Client, releaseURL string, token string) (*Config, error) {
	if token == "" {
		return nil, errors.New(GithubTokenEnv + " is not set")
	}
	release, err := get(ctx, client, releaseURL, token, "application/vnd.github+json")
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest release of the alerts config: %w", err)
	}
	var parsed struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(release, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse the latest release of the alerts config: %w", err)
	}
	for _, asset := range parsed.Assets {
		if asset.Name != assetName {
			continue
		}
		data, err := get(ctx, client, asset.URL, token, "application/octet-stream")
		if err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", assetName, err)
		}
		return Parse(data)
	}
	return nil, fmt.Errorf("the latest release of the alerts config has no %s", assetName)
}

func get(ctx context.Context, client *http.Client, url string, token string, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", accept)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}
