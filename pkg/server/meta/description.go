package meta

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakura-internet/go-rison/v4"
)

type descriptionRequestParams struct {
	Project string `json:"project"`
	Branch  string `json:"branch"`
}

type description struct {
	Project     string `json:"project"`
	Branch      string `json:"branch"`
	URL         string `json:"url"`
	MethodName  string `json:"methodName"`
	Description string `json:"description"`
}

func CreateGetDescriptionRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	cacheSize := 1000 * 1000 * 50
	cache := fastcache.New(cacheSize)
	return func(writer http.ResponseWriter, request *http.Request) {
		objectStart := strings.IndexRune(request.URL.Path, '(')
		var params descriptionRequestParams
		err := rison.Unmarshal([]byte(request.URL.Path[objectStart:]), &params, rison.Rison)
		if err != nil {
			slog.Error("cannot unmarshal parameters", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()
		_, _ = io.Copy(io.Discard, request.Body)

		cacheKey := []byte(params.Project + params.Branch)
		cachedValue, isInCache := cache.HasGet(nil, cacheKey)
		if isInCache {
			_, err = writer.Write(cachedValue)
			if err != nil {
				slog.Error("cannot write cached value", "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		rows, err := metaDb.Query(request.Context(), "SELECT project, branch, url, methodname, description FROM project_description WHERE project=$1 and branch=$2", params.Project, params.Branch)
		if err != nil {
			slog.Error("unable to execute the query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		desc, err := pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (description, error) {
			var project, branch, url, method_name, desc string
			err := row.Scan(&project, &branch, &url, &method_name, &desc)
			return description{
				Project:     project,
				Branch:      branch,
				URL:         url,
				MethodName:  method_name,
				Description: desc,
			}, err
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				_, err = writer.Write([]byte("{}"))
				if err != nil {
					slog.Error("unable to write response", "error", err)
					writer.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				slog.Error("unable to collect rows", "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		jsonBytes, err := json.Marshal(desc)
		if err != nil {
			slog.Error("unable to marshal description", "description", desc, "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(jsonBytes)
		cache.Set(cacheKey, jsonBytes)
		if err != nil {
			slog.Error("unable to write response", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}
