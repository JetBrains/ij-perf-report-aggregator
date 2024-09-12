package auth

import (
  "context"
  "encoding/json"
  "io"
  "net/http"
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

    userInfo, err := fetchUserInfo(r.Context(), accessToken)
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

func fetchUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
  client := &http.Client{}
  req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", http.NoBody)
  if err != nil {
    return nil, err
  }

  req.Header.Add("Authorization", "Bearer "+accessToken)

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
