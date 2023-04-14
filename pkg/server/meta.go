package server

import (
  "encoding/json"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgtype"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/sakura-internet/go-rison/v4"
  "go.uber.org/zap"
  "io"
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

type RequestParams struct {
  Tests    []string `json:"tests"`
  Interval string   `json:"interval"`
}

type InsertParams struct {
  Date        string `json:"date"`
  Test        string `json:"affected_test"`
  Reason      string `json:"reason"`
  BuildNumber string `json:"build_number"`
  Kind        string `json:"kind,omitempty"`
}

type DeleteParams struct {
  Id int64 `json:"id"`
}

func createGetMetaRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    objectStart := strings.IndexRune(request.URL.Path, '(')
    var params RequestParams
    err := rison.Unmarshal([]byte(request.URL.Path[objectStart:]), &params, rison.Rison)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    conn, err := metaDb.Acquire(request.Context())
    if err != nil {
      logger.Error("Cannot acquire connection for Postgres", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer conn.Release()
    sql := "SELECT id, date, affected_test, reason, build_number FROM accidents WHERE date >= CURRENT_DATE - INTERVAL '" + params.Interval + "'"
    if params.Tests != nil {
      sql += " and affected_test in (" + stringArrayToSQL(params.Tests) + ")"
    }
    rows, err := conn.Query(request.Context(), sql)
    if err != nil {
      logger.Error("Unable to execute the query", zap.String("query", sql))
      writer.WriteHeader(http.StatusInternalServerError)
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
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }

    jsonBytes, err := json.Marshal(accidents)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }
    _, err = writer.Write(jsonBytes)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }
  }
}

func createPostMetaRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    conn, err := metaDb.Acquire(request.Context())
    if err != nil {
      logger.Error("Cannot acquire connection for Postgres", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer conn.Release()

    var params InsertParams
    err = json.Unmarshal(all, &params)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    var kind string
    if params.Kind == "" {
      kind = "regression"
    } else {
      kind = params.Kind
    }

    _, err = conn.Exec(request.Context(), "INSERT INTO accidents (date, affected_test, reason, build_number, kind) VALUES ($1, $2, $3, $4, $5)", params.Date, params.Test, params.Reason, params.BuildNumber, kind)
    if err != nil {
      logger.Error("Cannot execute query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()
    writer.WriteHeader(http.StatusOK)
  }
}

func createDeleteMetaRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    conn, err := metaDb.Acquire(request.Context())
    if err != nil {
      logger.Error("Cannot acquire connection for Postgres", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer conn.Release()

    var params DeleteParams
    err = json.Unmarshal(all, &params)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }

    _, err = conn.Exec(request.Context(), "DELETE FROM accidents WHERE id=$1", params.Id)
    if err != nil {
      logger.Error("Cannot execute query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()
    writer.WriteHeader(http.StatusOK)
  }
}

func stringArrayToSQL(input []string) string {
  var str strings.Builder
  str.WriteRune('\'')
  str.WriteString(strings.Join(input, "','"))
  str.WriteRune('\'')
  return str.String()
}
