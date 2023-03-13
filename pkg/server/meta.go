package server

import (
  "encoding/json"
  "github.com/develar/errors"
  "github.com/sakura-internet/go-rison/v4"
  "github.com/valyala/bytebufferpool"
  "net/http"
  "strings"
)

func (t *StatsServer) handleMetaRequest(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  type RequestParams struct {
    Branch string `json:"branch"`
    Table  string `json:"table"`
  }
  objectStart := strings.IndexRune(request.URL.Path, '(')
  buffer := byteBufferPool.Get()
  var params RequestParams
  err := rison.Unmarshal([]byte(request.URL.Path[objectStart:]), &params, rison.Rison)
  if err != nil {
    return nil, false, err
  }
  conn := t.metaDb.Get(request.Context())
  if conn == nil {
    return nil, false, errors.New("Can't get connection to sqlite from pool")
  }
  defer t.metaDb.Put(conn)
  stmt := conn.Prep("SELECT id, date, affected_test, reason, build_number FROM accident WHERE branch=$branch and db_table=$table;")
  stmt.SetText("$branch", params.Branch)
  stmt.SetText("$table", params.Table)
  type Accident struct {
    ID           int64  `json:"id"`
    Date         string `json:"date"`
    AffectedTest string `json:"affectedTest"`
    Reason       string `json:"reason"`
    BuildNumber  string `json:"buildNumber"`
  }
  var accidents []Accident
  for {
    hasRow, err := stmt.Step()
    if err != nil {
      return nil, false, err
    }
    if !hasRow {
      break
    }
    accident := Accident{
      ID:           stmt.GetInt64("id"),
      Date:         stmt.GetText("date"),
      AffectedTest: stmt.GetText("affected_test"),
      Reason:       stmt.GetText("reason"),
      BuildNumber:  stmt.GetText("build_number"),
    }
    accidents = append(accidents, accident)
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
