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
	// GithubAppIDEnv and GithubAppKeyEnv name the id and the PEM private key of the GitHub App that reads the alerts
	// repository. This is what the deployment uses.
	GithubAppIDEnv  = "ALERTS_GITHUB_APP_ID"
	GithubAppKeyEnv = "ALERTS_GITHUB_APP_KEY"
	// GithubTokenEnv names a GitHub token that can read the releases of the alerts repository, for a local run:
	// ALERTS_GITHUB_TOKEN=$(gh auth token). It wins over the app.
	GithubTokenEnv = "ALERTS_GITHUB_TOKEN"

	githubAPI  = "https://api.github.com"
	repository = "JetBrains/ide-infrastructure-alerts"
	assetName  = "ij-perf.json"
	githubJSON = "application/vnd.github+json"
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
	config, err := download(ctx, client, githubAPI, credentialsFromEnv())
	if err == nil {
		return config, nil
	}
	builtIn, builtInErr := Parse(snapshot)
	if builtInErr != nil {
		return nil, errors.Join(err, builtInErr)
	}
	return builtIn, err
}

// credentials is how the download proves its right to read the repository: a token, or the app that gets one.
type credentials struct {
	token  string
	appID  string
	appKey string
}

func credentialsFromEnv() credentials {
	return credentials{token: os.Getenv(GithubTokenEnv), appID: os.Getenv(GithubAppIDEnv), appKey: os.Getenv(GithubAppKeyEnv)}
}

func (c credentials) accessToken(ctx context.Context, client *http.Client, apiURL string) (string, error) {
	if c.token != "" {
		return c.token, nil
	}
	if c.appID == "" || c.appKey == "" {
		return "", errors.New(GithubAppIDEnv + " and " + GithubAppKeyEnv + " are not set")
	}
	app, err := newGithubApp(c.appID, c.appKey)
	if err != nil {
		return "", err
	}
	return app.token(ctx, client, apiURL, repository)
}

func download(ctx context.Context, client *http.Client, apiURL string, credentials credentials) (*Config, error) {
	token, err := credentials.accessToken(ctx, client, apiURL)
	if err != nil {
		return nil, err
	}
	release, err := get(ctx, client, apiURL+"/repos/"+repository+"/releases/latest", token, githubJSON)
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
