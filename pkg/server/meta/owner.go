package meta

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateGetOwnerByProjectHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		project := request.URL.Query().Get("project")
		if project == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := metaDb.Query(request.Context(), "SELECT owner FROM project_owner WHERE project=$1", project)
		if err != nil {
			slog.Error("unable to execute the query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		owner, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (string, error) {
			var o string
			err := row.Scan(&o)
			return o, err
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				owner = ""
			} else {
				slog.Error("unable to collect row", "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		jsonBytes, err := json.Marshal(map[string]string{"owner": owner})
		if err != nil {
			slog.Error("unable to marshal owner", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(jsonBytes)
		if err != nil {
			slog.Error("unable to write response", "error", err)
		}
	}
}

func CreateGetProjectsByOwnerHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		owner := request.URL.Query().Get("owner")
		if owner == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := metaDb.Query(request.Context(), "SELECT project FROM project_owner WHERE owner=$1", owner)
		if err != nil {
			slog.Error("unable to execute the query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		projects, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
			var p string
			err := row.Scan(&p)
			return p, err
		})
		if err != nil {
			slog.Error("unable to collect rows", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if projects == nil {
			projects = []string{}
		}

		jsonBytes, err := json.Marshal(map[string][]string{"projects": projects})
		if err != nil {
			slog.Error("unable to marshal project owners", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(jsonBytes)
		if err != nil {
			slog.Error("unable to write response", "error", err)
		}
	}
}
