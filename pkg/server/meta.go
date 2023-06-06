package server

import (
  "encoding/json"
  "errors"
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
  Kind         string `json:"kind"`
  ExternalId   string `json:"externalId"`
}

type AccidentRequestParams struct {
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

type AccidentDeleteParams struct {
  Id int64 `json:"id"`
}

func createGetManyAccidentsRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error(err.Error())
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()

    var params AccidentRequestParams
    err = json.Unmarshal(all, &params)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }

    sql := "SELECT id, date, affected_test, reason, build_number, kind, externalId FROM accidents WHERE date >= CURRENT_DATE - INTERVAL '" + params.Interval + "'"
    if params.Tests != nil {
      sql += " and affected_test in (" + stringArrayToSQL(params.Tests) + ")"
    }
    rows, err := metaDb.Query(request.Context(), sql)
    if err != nil {
      logger.Error("Unable to execute the query", zap.String("query", sql))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer rows.Close()

    accidents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Accident, error) {
      var id int64
      var date pgtype.Date
      var affected_test, reason, build_number, kind, externalId string
      err := row.Scan(&id, &date, &affected_test, &reason, &build_number, &kind, &externalId)
      return Accident{
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
    }
    jsonBytes, err := json.Marshal(accidents)
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

func createPostAccidentRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }

    var params AccidentInsertParams
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

    _, err = metaDb.Exec(request.Context(), "INSERT INTO accidents (date, affected_test, reason, build_number, kind, externalId) VALUES ($1, $2, $3, $4, $5, $6)", params.Date, params.Test, params.Reason, params.BuildNumber, kind, params.ExternalId)
    if err != nil {
      logger.Error("Cannot execute query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()
    writer.WriteHeader(http.StatusOK)
  }
}

func createDeleteAccidentRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    body := request.Body
    all, err := io.ReadAll(body)
    if err != nil {
      logger.Error("Cannot read body", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }

    var params AccidentDeleteParams
    err = json.Unmarshal(all, &params)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }

    _, err = metaDb.Exec(request.Context(), "DELETE FROM accidents WHERE id=$1", params.Id)
    if err != nil {
      logger.Error("Cannot execute query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
    }
    defer body.Close()
    writer.WriteHeader(http.StatusOK)
  }
}

type DescriptionRequestParams struct {
  Project string `json:"project"`
  Branch  string `json:"branch"`
}

type Description struct {
  Project     string `json:"project"`
  Branch      string `json:"branch"`
  URL         string `json:"url"`
  MethodName  string `json:"methodName"`
  Description string `json:"description"`
}

func createGetDescriptionRequestHandler(logger *zap.Logger, metaDb *pgxpool.Pool) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    objectStart := strings.IndexRune(request.URL.Path, '(')
    var params DescriptionRequestParams
    err := rison.Unmarshal([]byte(request.URL.Path[objectStart:]), &params, rison.Rison)
    if err != nil {
      logger.Error("Cannot unmarshal parameters", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }

    rows, err := metaDb.Query(request.Context(), "SELECT project, branch, url, methodname, description FROM project_description WHERE project=$1 and branch=$2", params.Project, params.Branch)
    if err != nil {
      logger.Error("Unable to execute the query", zap.Error(err))
      writer.WriteHeader(http.StatusInternalServerError)
      return
    }
    defer rows.Close()
    description, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (Description, error) {
      var project, branch, url, method_name, description string
      err := row.Scan(&project, &branch, &url, &method_name, &description)
      return Description{
        Project:     project,
        Branch:      branch,
        URL:         url,
        MethodName:  method_name,
        Description: description,
      }, err
    })
    if err != nil {
      if errors.Is(err, pgx.ErrNoRows) {
        _, err = writer.Write([]byte("{}"))
        if err != nil {
          logger.Error(err.Error())
          writer.WriteHeader(http.StatusInternalServerError)
        }
      } else {
        logger.Error(err.Error())
        writer.WriteHeader(http.StatusInternalServerError)
      }
      return
    }

    jsonBytes, err := json.Marshal(description)
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

func stringArrayToSQL(input []string) string {
  var str strings.Builder
  str.WriteRune('\'')
  str.WriteString(strings.Join(input, "','"))
  str.WriteRune('\'')
  return str.String()
}
