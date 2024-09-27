package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	"io"
	"net/http"
	"net/url"
	"time"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	HD            string `json:"hd"`
}

func CreateGetUserInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("X-Auth-Request-Access-Token")
		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userInfo, err := FetchUserInfo(r.Context(), accessToken)
		if err != nil {
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(userInfo)
		if err != nil {
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}
	}
}

func FetchUserInfo(ctx context.Context, googleToken string) (*UserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+googleToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

type YTAuth struct {
	cache         *fastcache.Cache
	youtrackToken string
	youtrackUrl   string
	client        http.Client
}

type YTUser struct {
	ID string `json:"id"`
}

func NewYTAuth(youtrackUrl string, youtrackToken string) *YTAuth {
	return &YTAuth{
		cache:         fastcache.New(10 * 1000 * 1000),
		youtrackToken: youtrackToken,
		youtrackUrl:   youtrackUrl,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ytAuth *YTAuth) GetUser(ctx context.Context, email string) (*YTUser, error) {
	cachedId, exists := ytAuth.cache.HasGet(nil, []byte(email))
	if exists {
		var v YTUser
		err := json.Unmarshal(cachedId, &v)
		return &v, err
	}
	hubId, err := ytAuth.getHubUserId(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("cannot get Hub user id for %s: %w", email, err)
	}
	if hubId == "" {
		return nil, fmt.Errorf("user %s not found in JetBrains Hub", email)
	}
	ytId, err := ytAuth.getYTUserId(ctx, hubId)
	if err != nil {
		return nil, fmt.Errorf("cannot get YouTrack user id for %s: %w", email, err)
	}
	v, err := json.Marshal(ytId)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal YTUser: %w", err)
	}
	ytAuth.cache.Set([]byte(email), v)
	return ytId, nil
}

func (ytAuth *YTAuth) getYTUserId(ctx context.Context, hubToken string) (*YTUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ytAuth.youtrackUrl+"/api/users/"+hubToken+"?fields=id", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+ytAuth.youtrackToken)
	resp, err := ytAuth.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to YouTrack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service YouTrack API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo = &YTUser{}
	err = json.Unmarshal(body, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YTUser: %w", err)
	}

	return userInfo, nil
}

type usersPage struct {
	Users []hubUser `json:"users"`
}

type hubUser struct {
	ID string `json:"id"`
}

func (ytAuth *YTAuth) getHubUserId(ctx context.Context, email string) (string, error) {
	emailEscaped := url.QueryEscape(email)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://hub.jetbrains.com/api/rest/users?fileds=id&query=email:"+emailEscaped, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+ytAuth.youtrackToken)
	resp, err := ytAuth.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to JetBrains Hub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("service Hub API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo usersPage
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal user info: %w", err)
	}
	if len(userInfo.Users) == 0 {
		return "", nil
	}

	return userInfo.Users[0].ID, nil
}
