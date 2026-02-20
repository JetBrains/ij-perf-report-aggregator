package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accidentResponse struct {
	ID           int64  `json:"id"`
	Date         string `json:"date"`
	AffectedTest string `json:"affectedTest"`
	Reason       string `json:"reason"`
	BuildNumber  string `json:"buildNumber"`
	Kind         string `json:"kind"`
	ExternalId   string `json:"externalId,omitempty"`
	Stacktrace   string `json:"stacktrace"`
	UserName     string `json:"userName"`
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
	Stacktrace  string `json:"stacktrace"`
	UserName    string `json:"user_name,omitempty"`
}

type accidentIdParams struct {
	Id int64 `json:"id"`
}

func CreateGetAccidentsAroundDateRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	type date struct {
		Date string `json:"date"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			slog.Error("Cannot read body", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer body.Close()

		var params date
		if err = json.Unmarshal(all, &params); err != nil {
			slog.Error("Cannot unmarshal parameters", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		sql := "SELECT id, date, affected_test, reason, build_number, kind, stacktrace, user_name FROM accidents WHERE (LOWER(kind)='regression' or LOWER(kind)='improvement' or LOWER(kind)='investigation') AND date BETWEEN '" + params.Date + "'::date - INTERVAL '1 days' AND '" + params.Date + "'::date + INTERVAL '1 days'"
		rows, err := metaDb.Query(request.Context(), sql)
		if err != nil {
			slog.Error("unable to execute the query", "query", sql, "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		accidents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (accidentResponse, error) {
			var id int64
			var date pgtype.Date
			var affectedTest, reason, buildNumber, kind, stacktrace, userName string
			err := row.Scan(&id, &date, &affectedTest, &reason, &buildNumber, &kind, &stacktrace, &userName)
			return accidentResponse{
				ID:           id,
				Date:         date.Time.String(),
				AffectedTest: affectedTest,
				Reason:       reason,
				BuildNumber:  buildNumber,
				Kind:         kind,
				Stacktrace:   stacktrace,
				UserName:     userName,
			}, err
		})
		if err != nil {
			slog.Error("unable to collect rows", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonBytes, err := json.Marshal(accidents)
		if err != nil {
			slog.Error("unable to marshal accidents", "accidents", accidents, "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(jsonBytes)
		if err != nil {
			slog.Error("unable to write response", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CreateGetManyAccidentsRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			slog.Error("Cannot read body", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer body.Close()

		var params accidentRequestParams
		if err = json.Unmarshal(all, &params); err != nil {
			slog.Error("Cannot unmarshal parameters", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		sql := "SELECT id, date, affected_test, reason, build_number, kind, externalId, stacktrace, user_name FROM accidents WHERE date >= CURRENT_DATE - INTERVAL '" + params.Interval + "'"
		if params.Tests != nil {
			sql += " and affected_test in (" + stringArrayToSQL(params.Tests) + ") or affected_test = ''"
		}
		rows, err := metaDb.Query(request.Context(), sql)
		if err != nil {
			slog.Error("unable to execute the query", "query", sql, "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if _, err := writer.Write([]byte("[")); err != nil {
			slog.Error("Failed to write JSON array start", "error", err)
			return
		}

		firstItem := true
		for rows.Next() {
			accident, err := getAccidentFromRow(rows)
			if err != nil {
				slog.Error("unable to scan row", "error", err)
				// We've already started sending the response, so we can't change the status code
				// Best we can do is log the error and stop
				return
			}

			// Add comma separator between items (but not before the first item)
			if !firstItem {
				if _, err := writer.Write([]byte(",")); err != nil {
					slog.Error("Failed to write comma separator", "error", err)
					return
				}
			} else {
				firstItem = false
			}

			// Marshal and write this individual item
			itemBytes, err := json.Marshal(accident)
			if err != nil {
				slog.Error("unable to marshal accident", "error", err)
				return
			}

			if _, err := writer.Write(itemBytes); err != nil {
				slog.Error("Failed to write item", "error", err)
				return
			}

			// Flush the response writer if it supports flushing
			if flusher, ok := writer.(http.Flusher); ok {
				flusher.Flush()
			}
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			slog.Error("Error iterating over rows", "error", err)
			return
		}

		// Close the JSON array
		if _, err := writer.Write([]byte("]")); err != nil {
			slog.Error("Failed to write JSON array end", "error", err)
			return
		}
	}
}

func CreatePostAccidentRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			slog.Error("cannot read body", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer body.Close()

		var params AccidentInsertParams
		if err = json.Unmarshal(all, &params); err != nil {
			slog.Error("cannot unmarshal parameters", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var kind string
		if params.Kind == "" {
			kind = "regression"
		} else {
			kind = params.Kind
		}

		// For inferred kinds, check if any accident already exists for this date/test/build
		// regardless of kind. This prevents duplicates when a user reclassifies an
		// InferredRegression to Regression and the detector runs again.
		if kind == "InferredRegression" || kind == "InferredImprovement" {
			var exists bool
			err := metaDb.QueryRow(request.Context(),
				"SELECT EXISTS(SELECT 1 FROM accidents WHERE date = $1 AND affected_test = $2 AND build_number = $3 AND kind <> 'EXCEPTION')",
				params.Date, params.Test, params.BuildNumber).Scan(&exists)
			if err != nil {
				slog.Error("cannot check for existing accident", "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			if exists {
				http.Error(writer, "Conflict: Accident already exists", http.StatusConflict)
				return
			}
		}

		var id int
		idRow := metaDb.QueryRow(request.Context(), "INSERT INTO accidents (date, affected_test, reason, build_number, kind, externalId, stacktrace, user_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", params.Date, params.Test, params.Reason, params.BuildNumber, kind, params.ExternalId, params.Stacktrace, params.UserName)
		if err = idRow.Scan(&id); err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				http.Error(writer, "Conflict: Accident already exists", http.StatusConflict)
			} else {
				slog.Error("cannot execute insert accidents query", "error", err, "date", params.Date, "affected_test", params.Test, "reason", params.Reason, "build_number", params.BuildNumber, "kind", kind, "externalId", params.ExternalId)
				writer.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		_, err = writer.Write([]byte(strconv.Itoa(id)))
		if err != nil {
			slog.Error("cannot write response", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func CreateDeleteAccidentRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			slog.Error("cannot read body", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer body.Close()

		var params accidentIdParams
		if err = json.Unmarshal(all, &params); err != nil {
			slog.Error("cannot unmarshal parameters", "body", all, "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = metaDb.Exec(request.Context(), "DELETE FROM accidents WHERE id=$1", params.Id)
		if err != nil {
			slog.Error("cannot execute query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func CreateGetAccidentByIdHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")

		accidentById, err := getAccidentById(request.Context(), metaDb, id)
		if err != nil {
			slog.Error("cannot get accident by id", "error", err, "id", id)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(accidentById)
		if err != nil {
			slog.Error("cannot marshal accident", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(jsonBytes)
		if err != nil {
			slog.Error("cannot write response", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func stringArrayToSQL(input []string) string {
	var str strings.Builder
	str.WriteRune('\'')

	for i, s := range input {
		// Escape any single quotes in the string
		escapedStr := strings.ReplaceAll(s, "'", "''")
		str.WriteString(escapedStr)

		// Add a separator if it's not the last element
		if i < len(input)-1 {
			str.WriteString("','")
		}
	}

	str.WriteRune('\'')
	return str.String()
}

func getAccidentFromRow(row pgx.CollectableRow) (accidentResponse, error) {
	var id int64
	var date pgtype.Date
	var affected_test, reason, build_number, kind, externalId, stacktrace, user_name string
	err := row.Scan(&id, &date, &affected_test, &reason, &build_number, &kind, &externalId, &stacktrace, &user_name)
	return accidentResponse{
		ID:           id,
		Date:         date.Time.String(),
		AffectedTest: affected_test,
		Reason:       reason,
		BuildNumber:  build_number,
		Kind:         kind,
		ExternalId:   externalId,
		Stacktrace:   stacktrace,
		UserName:     user_name,
	}, err
}

func getAccidentById(ctx context.Context, metaDb *pgxpool.Pool, accidentId string) (*accidentResponse, error) {
	sql := "SELECT id, date, affected_test, reason, build_number, kind, externalId, stacktrace, user_name FROM accidents WHERE id=$1"
	rows, err := metaDb.Query(ctx, sql, accidentId)
	if err != nil {
		log.Println("unable to execute the query", "query", sql, "error", err)
		return nil, err
	}
	defer rows.Close()

	accidents, err := pgx.CollectRows(rows, getAccidentFromRow)
	if err != nil {
		log.Println("unable to collect rows", "error", err)
		return nil, err
	}

	if len(accidents) == 0 {
		return nil, fmt.Errorf("no accident found with id: %s", accidentId)
	}

	return &accidents[0], nil
}

func updateAccidentReason(ctx context.Context, metaDb *pgxpool.Pool, accident *accidentResponse) error {
	sql := `UPDATE accidents SET reason = $2 WHERE id = $1`
	_, err := metaDb.Exec(ctx, sql,
		accident.ID,
		accident.Reason,
	)
	if err != nil {
		log.Println("unable to update the query", "query", sql, "error", err)
		return err
	}
	return nil
}
