package server

import (
  "github.com/develar/errors"
  "go.uber.org/zap"
  "net/http"
  "net/url"
  "report-aggregator/pkg/util"
  "time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func performRequest(u *url.URL, logger *zap.Logger) (*http.Response, error) {
  r, err := httpClient.Get(u.String())
  if err != nil {
    return nil, err
  }

  if r.StatusCode >= 400 {
    util.Close(r.Body, logger)
    return nil, errors.Errorf("Request failed: %s", r.Status)
  }
  return r, nil
}

func (t *StatsServer) performRequest(unescapedQuery string) (*http.Response, error) {
  u, err := t.buildUrl(unescapedQuery)
  if err != nil {
    return nil, err
  }
  return performRequest(u, t.logger)
}

func (t *StatsServer) buildUrl(unescapedQuery string) (*url.URL, error) {
  relativeUrl, err := url.Parse("/api/v1/query")
  if err != nil {
    return nil, err
  }

  u := t.victoriaMetricsServerUrl.ResolveReference(relativeUrl)
  q := u.Query()
  q.Set("query", unescapedQuery)
  u.RawQuery = q.Encode()
  return u, nil
}
