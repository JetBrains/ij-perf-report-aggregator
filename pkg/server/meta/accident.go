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

type accident struct {
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

		accidents, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (accident, error) {
			var id int64
			var date pgtype.Date
			var affectedTest, reason, buildNumber, kind, stacktrace, userName string
			err := row.Scan(&id, &date, &affectedTest, &reason, &buildNumber, &kind, &stacktrace, &userName)
			return accident{
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

		accidents, err := pgx.CollectRows(rows, getAccidentFromRow)
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

func getAccidentFromRow(row pgx.CollectableRow) (accident, error) {
	var id int64
	var date pgtype.Date
	var affected_test, reason, build_number, kind, externalId, stacktrace, user_name string
	err := row.Scan(&id, &date, &affected_test, &reason, &build_number, &kind, &externalId, &stacktrace, &user_name)
	return accident{
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

func getAccidentById(ctx context.Context, metaDb *pgxpool.Pool, accidentId string) (*accident, error) {
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

func updateAccidentReason(ctx context.Context, metaDb *pgxpool.Pool, accident *accident) error {
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
