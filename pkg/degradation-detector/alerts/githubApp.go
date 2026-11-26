package alerts

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// githubApp is the GitHub App that reads the private alerts repository. An app and not a personal token, so the
// credential belongs to no person. The app needs the permission to read the contents and an installation on the
// repository. Socrat reads its settings with the same app.
type githubApp struct {
	id         string
	privateKey *rsa.PrivateKey
}

// newGithubApp reads the private key of the app: the PEM file that GitHub gives (PKCS #1), the same key in PKCS #8, or
// in the format of OpenSSH. So the key needs no conversion. The key can be on one line with `\n` instead of line breaks.
func newGithubApp(id string, privateKey string) (*githubApp, error) {
	parsed, err := ssh.ParseRawPrivateKey([]byte(strings.ReplaceAll(strings.TrimSpace(privateKey), `\n`, "\n") + "\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to read the private key of the GitHub App, it must be a PEM or OpenSSH private key without a passphrase: %w", err)
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("the private key of the GitHub App must be an RSA key, and it is %T", parsed)
	}
	return &githubApp{id: strings.TrimSpace(id), privateKey: key}, nil
}

// token returns an access token of the installation of the app on the repository. The token lives an hour, which is
// enough for the one download of a run.
func (a *githubApp) token(ctx context.Context, client *http.Client, apiURL string, repository string) (string, error) {
	jwt, err := a.jwt(time.Now())
	if err != nil {
		return "", err
	}
	var installation struct {
		ID int64 `json:"id"`
	}
	if err := a.request(ctx, client, http.MethodGet, apiURL+"/repos/"+repository+"/installation", jwt, http.StatusOK, &installation); err != nil {
		return "", fmt.Errorf("failed to find the installation of the GitHub App %s on %s: %w", a.id, repository, err)
	}
	var minted struct {
		Token string `json:"token"`
	}
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", apiURL, installation.ID)
	if err := a.request(ctx, client, http.MethodPost, url, jwt, http.StatusCreated, &minted); err != nil {
		return "", fmt.Errorf("failed to get an access token of the GitHub App %s: %w", a.id, err)
	}
	if minted.Token == "" {
		return "", fmt.Errorf("GitHub returned no access token of the GitHub App %s", a.id)
	}
	return minted.Token, nil
}

func (a *githubApp) request(ctx context.Context, client *http.Client, method string, url string, jwt string, status int, result any) error {
	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", githubJSON)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != status {
		return fmt.Errorf("%s: %s", url, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

// jwt returns the token that proves to GitHub that the request comes from the app. GitHub accepts it for at most
// ten minutes, and the minute in the past covers a clock that runs ahead.
func (a *githubApp) jwt(now time.Time) (string, error) {
	encode := base64.RawURLEncoding.EncodeToString
	claims, err := json.Marshal(map[string]any{
		"iss": a.id,
		"iat": now.Add(-time.Minute).Unix(),
		"exp": now.Add(9 * time.Minute).Unix(),
	})
	if err != nil {
		return "", err
	}
	unsigned := encode([]byte(`{"alg":"RS256","typ":"JWT"}`)) + "." + encode(claims)
	hash := sha256.Sum256([]byte(unsigned))
	signature, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign the token of the GitHub App: %w", err)
	}
	return unsigned + "." + encode(signature), nil
}
