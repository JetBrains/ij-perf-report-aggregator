package meta

import (
  "encoding/json"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgtype"
  "github.com/jackc/pgx/v5/pgxpool"
  "go.uber.org/zap"
  "io"
  "net/http"
  "strconv"
  "strings"
)

type accident struct {
  ID           int64  `json:"id"`
  Date         string `json:"date"`
  AffectedTest string `json:"affectedTest"`
  Reason       string `json:"reason"`
  BuildNumber  string `json:"buildNumber"`
  Kind         string `json:"kind"`
  ExternalId   string `json:"externalId"`
}

type accidentRequestParams struct {
  Tests    []string `json:"tests"`
  Interval string   `json:"interval"`
}

type AccidentInsertParams struct {
  Date        string `json:"date"`
  Test        string `json:"affected_test"`
  Reason      string `json:"reason"`
  BuildNumber string `json:"build_number"`
  Kind        string `json:"kind,omitempty"`
  ExternalId  string `json:"externalId,omitempty"`
}

type accidentDeleteParams struct {
  Id int64 `json:"id"`
}

func CreateGetManyAccidentsRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    defer body.Close()

    var params accidentRequestParams
    if err = json.Unmarshal(all, &params); err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }

    sql := "SELECT id, date, affected_test, reason, build_number, kind, externalId FROM accidents WHERE date >= CURRENT_DATE - INTERVAL '" + params.Interval + "'"
    if params.Tests != nil {
      sql += " and affected_test in (" + stringArrayToSQL(params.Tests) + ")"
    }
    rows, err := metaDb.Query(request.Context(), sql)
    if err != nil {
      logger.Error("Unable to execute the query", zap.String("query", sql))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    defer rows.Close()

    accidents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (accident, error) {
      var id int64
      var date pgtype.Date
      var affected_test, reason, build_number, kind, externalId string
      err := row.Scan(&id, &date, &affected_test, &reason, &build_number, &kind, &externalId)
      return accident{
        ID:           id,
        Date:         date.Time.String(),
        AffectedTest: affected_test,
        Reason:       reason,
        BuildNumber:  build_number,
        Kind:         kind,
        ExternalId:   externalId,
      }, err
    })
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    jsonBytes, err := json.Marshal(accidents)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    _, err = writer.Write(jsonBytes)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }
  }
}

func CreatePostAccidentRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }

    defer body.Close()

    var params AccidentInsertParams
    if err = json.Unmarshal(all, &params); err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    var kind string
    if params.Kind == "" {
      kind = "regression"
    } else {
      kind = params.Kind
    }

    var id int
    idRow := metaDb.QueryRow(request.Context(), "INSERT INTO accidents (date, affected_test, reason, build_number, kind, externalId) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", params.Date, params.Test, params.Reason, params.BuildNumber, kind, params.ExternalId)
    if err = idRow.Scan(&id); err != nil {
      if strings.Contains(err.Error(), "unique constraint") {
        http.Error(writer, "Conflict: Accident already exists", http.StatusConflict)
      } else {
        logger.Error("Cannot execute query", zap.Error(err))
        writer.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    writer.WriteHeader(http.StatusOK)
    _, err = writer.Write([]byte(strconv.Itoa(id)))
    if err != nil {
      logger.Error("Cannot write response", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
  }
}

func CreateDeleteAccidentRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()

    var params accidentDeleteParams
    if err = json.Unmarshal(all, &params); err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }

    _, err = metaDb.Exec(request.Context(), "DELETE FROM accidents WHERE id=$1", params.Id)
    if err != nil {
      logger.Error("Cannot execute query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
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
