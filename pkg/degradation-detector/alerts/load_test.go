package alerts

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

const downloaded = `{"formatVersion": 1, "service": "ij-perf", "version": "v100", "notifications": {"n": {}}}`

const releasePath = "/repos/JetBrains/ide-infrastructure-alerts/releases/latest"

// releaseServer is GitHub with one release of the alerts repository and one installation of the app 42 on it.
// It accepts the token "secret" and gives it to the app for a token signed with the key.
func releaseServer(t *testing.T, appKey *rsa.PublicKey, withAsset bool) *httptest.Server {
	t.Helper()
	var server *httptest.Server
	authorized := func(r *http.Request) bool { return r.Header.Get("Authorization") == "Bearer secret" }
	fromApp := func(r *http.Request) bool {
		return appKey != nil && verifyAppToken(appKey, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	}
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/JetBrains/ide-infrastructure-alerts/installation" && fromApp(r):
			_, _ = w.Write([]byte(`{"id": 7}`))
		case r.URL.Path == "/app/installations/7/access_tokens" && r.Method == http.MethodPost && fromApp(r):
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"token": "secret"}`))
		case !authorized(r):
			w.WriteHeader(http.StatusUnauthorized)
		case r.URL.Path == releasePath:
			assets := `{"name": "manifest.json", "url": "` + server.URL + `/manifest"}`
			if withAsset {
				assets += `, {"name": "ij-perf.json", "url": "` + server.URL + `/asset"}`
			}
			_, _ = w.Write([]byte(`{"assets": [` + assets + `]}`))
		case r.URL.Path == "/asset":
			assert.Equal(t, "application/octet-stream", r.Header.Get("Accept"))
			_, _ = w.Write([]byte(downloaded))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// verifyAppToken checks the token as GitHub does: signed by the app 42 and valid now for at most ten minutes.
func verifyAppToken(key *rsa.PublicKey, token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	hash := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], signature) != nil {
		return false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var claims struct {
		Iss string `json:"iss"`
		Iat int64  `json:"iat"`
		Exp int64  `json:"exp"`
	}
	if json.Unmarshal(payload, &claims) != nil {
		return false
	}
	now := time.Now().Unix()
	return claims.Iss == "42" && claims.Iat <= now && claims.Exp > now && claims.Exp-claims.Iat <= 600
}

func appKeyPem(t *testing.T, pkcs8 bool) (string, *rsa.PublicKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	if pkcs8 {
		der, err := x509.MarshalPKCS8PrivateKey(key)
		require.NoError(t, err)
		block = &pem.Block{Type: "PRIVATE KEY", Bytes: der}
	}
	return string(pem.EncodeToMemory(block)), &key.PublicKey
}

func TestDownload(t *testing.T) {
	t.Parallel()
	t.Run("with a token", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, nil, true)
		config, err := download(t.Context(), server.Client(), server.URL, credentials{token: "secret"})
		require.NoError(t, err)
		assert.Equal(t, "v100", config.Version)
	})
	t.Run("with the GitHub App", func(t *testing.T) {
		t.Parallel()
		key, public := appKeyPem(t, false)
		server := releaseServer(t, public, true)
		config, err := download(t.Context(), server.Client(), server.URL, credentials{appID: " 42\n", appKey: key})
		require.NoError(t, err)
		assert.Equal(t, "v100", config.Version)
	})
	t.Run("with the key of the GitHub App on one line and in PKCS#8", func(t *testing.T) {
		t.Parallel()
		key, public := appKeyPem(t, true)
		server := releaseServer(t, public, true)
		oneLine := strings.ReplaceAll(strings.TrimSpace(key), "\n", `\n`)
		config, err := download(t.Context(), server.Client(), server.URL, credentials{appID: "42", appKey: oneLine})
		require.NoError(t, err)
		assert.Equal(t, "v100", config.Version)
	})
	t.Run("with the key of another app", func(t *testing.T) {
		t.Parallel()
		key, _ := appKeyPem(t, false)
		_, public := appKeyPem(t, false)
		server := releaseServer(t, public, true)
		_, err := download(t.Context(), server.Client(), server.URL, credentials{appID: "42", appKey: key})
		assert.ErrorContains(t, err, "failed to find the installation of the GitHub App 42")
	})
	t.Run("with the key of the GitHub App in the OpenSSH format", func(t *testing.T) {
		t.Parallel()
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)
		block, err := ssh.MarshalPrivateKey(key, "")
		require.NoError(t, err)
		server := releaseServer(t, &key.PublicKey, true)
		config, err := download(t.Context(), server.Client(), server.URL, credentials{appID: "42", appKey: string(pem.EncodeToMemory(block))})
		require.NoError(t, err)
		assert.Equal(t, "v100", config.Version)
	})
	t.Run("with a public key", func(t *testing.T) {
		t.Parallel()
		_, err := download(t.Context(), http.DefaultClient, "http://unused.invalid", credentials{appID: "42", appKey: "ssh-rsa AAAA"})
		assert.ErrorContains(t, err, "private key")
	})
	t.Run("with a key that is not RSA", func(t *testing.T) {
		t.Parallel()
		_, key, err := ed25519.GenerateKey(rand.Reader)
		require.NoError(t, err)
		block, err := ssh.MarshalPrivateKey(key, "")
		require.NoError(t, err)
		_, err = download(t.Context(), http.DefaultClient, "http://unused.invalid", credentials{appID: "42", appKey: string(pem.EncodeToMemory(block))})
		assert.ErrorContains(t, err, "must be an RSA key")
	})
	t.Run("without credentials", func(t *testing.T) {
		t.Parallel()
		_, err := download(t.Context(), http.DefaultClient, "http://unused.invalid", credentials{appID: "42"})
		assert.ErrorContains(t, err, GithubAppKeyEnv)
	})
	t.Run("wrong token", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, nil, true)
		_, err := download(t.Context(), server.Client(), server.URL, credentials{token: "wrong"})
		assert.ErrorContains(t, err, "401")
	})
	t.Run("release without the asset", func(t *testing.T) {
		t.Parallel()
		server := releaseServer(t, nil, false)
		_, err := download(t.Context(), server.Client(), server.URL, credentials{token: "secret"})
		assert.ErrorContains(t, err, "has no ij-perf.json")
	})
}

func TestLoadUsesBuiltInConfigWithoutCredentials(t *testing.T) {
	t.Setenv(GithubTokenEnv, "")
	t.Setenv(GithubAppIDEnv, "")
	t.Setenv(GithubAppKeyEnv, "")
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
