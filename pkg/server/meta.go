package server

import (
  "encoding/json"
  "github.com/develar/errors"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgtype"
  "github.com/sakura-internet/go-rison/v4"
  "github.com/valyala/bytebufferpool"
  "go.uber.org/zap"
  "net/http"
  "strings"
)

type Accident struct {
  ID           int64  `json:"id"`
  Date         string `json:"date"`
  AffectedTest string `json:"affectedTest"`
  Reason       string `json:"reason"`
  BuildNumber  string `json:"buildNumber"`
}

func (t *StatsServer) handleMetaRequest(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  type RequestParams struct {
    Branches []string `json:"branches"`
    Table    string   `json:"table"`
    Tests    []string `json:"tests"`
  }
  objectStart := strings.IndexRune(request.URL.Path, '(')
  buffer := byteBufferPool.Get()
  var params RequestParams
  err := rison.Unmarshal([]byte(request.URL.Path[objectStart:]), &params, rison.Rison)
  if err != nil {
    return nil, false, err
  }
  conn, err := t.metaDb.Acquire(request.Context())
  if err != nil {
    t.logger.Error("Cannot acquire connection for Postgres")
    return nil, false, err
  }
  if conn == nil {
    return nil, false, errors.New("Can't get connection to sqlite from pool")
  }
  defer conn.Release()
  sql := "SELECT id, date, affected_test, reason, build_number FROM accidents WHERE branch in (" + stringArrayToSQL(params.Branches) + ") and db_table=$1"
  if params.Tests != nil {
    sql += " and affected_test in (" + stringArrayToSQL(params.Tests) + ")"
  }
  rows, err := conn.Query(request.Context(), sql, params.Table)
  if err != nil {
    t.logger.Error("Unable to execute the query", zap.String("query", sql))
    return nil, false, err
  }
  defer rows.Close()
  var id int64
  var date pgtype.Date
  var affected_test, reason, build_number pgtype.Text
  var accidents []Accident
  _, err = pgx.ForEachRow(rows, []any{&id, &date, &affected_test, &reason, &build_number}, func() error {
    accident := Accident{
      ID:           id,
      Date:         date.Time.String(),
      AffectedTest: affected_test.String,
      Reason:       reason.String,
      BuildNumber:  build_number.String,
    }
    accidents = append(accidents, accident)
    return nil
  })
  if err != nil {
    t.logger.Error(err.Error())
    return nil, false, err
  }

  jsonBytes, err := json.Marshal(accidents)
  if err != nil {
    return nil, false, err
  }
  _, err = buffer.Write(jsonBytes)
  if err != nil {
    return nil, false, err
  }
  return buffer, true, nil
}

func stringArrayToSQL(input []string) string {
  var str strings.Builder
  str.WriteRune('\'')
  str.WriteString(strings.Join(input, "','"))
  str.WriteRune('\'')
  return str.String()
}
