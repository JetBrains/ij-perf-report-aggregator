package server

import (
  "github.com/develar/errors"
  "net/http"
)

func (t *StatsServer) handleReportRequest(request *http.Request) ([]byte, error) {
  query, err := ReadQuery(request)
  if err != nil {
    return nil, err
  }

  query.Fields = []DataQueryDimension{{Name: "raw_report"}}
  var rawReport []byte
  row, err := SelectRow(query, "report", t.db, request.Context())
  if err != nil {
    return nil, err
  }

  err = row.Scan(&rawReport)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return rawReport, nil
}