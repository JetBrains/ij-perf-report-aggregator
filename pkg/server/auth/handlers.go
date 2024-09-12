package auth

import (
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

    userInfo, err := fetchUserInfo(accessToken)
    if err != nil {
      http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
      return
    }

    json.NewEncoder(w).Encode(userInfo)
  }
}

func fetchUserInfo(accessToken string) (*UserInfo, error) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
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
