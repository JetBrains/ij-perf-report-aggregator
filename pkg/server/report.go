package server

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "net/http"
)

func (t *StatsServer) handleReportRequest(request *http.Request) ([]byte, error) {
  queries, _, err := data_query.ReadQuery(request)
  if err != nil {
    return nil, err
  }

  query := queries[0]
  query.Fields = []data_query.DataQueryDimension{{Name: "raw_report"}}
  query.Order = nil
  rawReport, err := data_query.SelectString(query, "report", t, request.Context())
  if err != nil {
    return nil, err
  }
  return rawReport, nil
}
