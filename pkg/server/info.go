package server

import (
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/valyala/quicktemplate"
  "net/http"
  "report-aggregator/pkg/util"
)

var essentialMetricNames = []string{"bootstrap", "appInitPreparation", "appInit", "pluginDescriptorLoading", "appComponentCreation", "projectComponentCreation"}

func (t *StatsServer) handleInfoRequest(w http.ResponseWriter, _ *http.Request) {
  var errRef error
  t.infoResponseCacheOnce.Do(func() {
    productNames, err := t.getProductNames()
    if err != nil {
      errRef = err
      return
    }

    statement, err := t.db.Prepare("select rowid as id, name from machine where rowid in (select distinct machine from report where product = ?) order by name")
    if err != nil {
      errRef = err
      return
    }

    defer util.Close(statement, t.logger)

    buffer := quicktemplate.AcquireByteBuffer()
    defer quicktemplate.ReleaseByteBuffer(buffer)
    WriteInfo(buffer, productNames, essentialMetricNames, statement, &errRef)
    if errRef != nil {
      return
    }

    result := make([]byte, len(buffer.B))
    copy(result, buffer.B)
    t.infoResponseCache = result
  })

  if errRef != nil {
    t.httpError(errRef, w)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  _, err := w.Write(t.infoResponseCache)
  if err != nil {
    t.httpError(err, w)
    return
  }
}

func (t *StatsServer) getProductNames() ([]string, error) {
  statement, err := t.db.Prepare("select distinct product from report order by product")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  defer util.Close(statement, t.logger)

  return t.readStringList(statement)
}

func (t *StatsServer) readStringList(statement *sqlite3.Stmt) ([]string, error) {
  var result []string
  for {
    hasRow, err := statement.Step()
    if err != nil {
      return nil, errors.WithStack(err)
    }

    value, _, err := statement.ColumnText(0)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    if !hasRow {
      break
    }

    result = append(result, value)
  }
  return result, nil
}
