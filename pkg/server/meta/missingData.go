package meta

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MissingDataInsertParams struct {
	BuildType    string `json:"build_type"`
	Project      string `json:"project"`
	Metric       string `json:"metric"`
	MissingSince int64  `json:"missing_since"`
}

func CreatePostMissingDataRequestHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			slog.Error("cannot read body", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer body.Close()

		var params MissingDataInsertParams
		if err = json.Unmarshal(all, &params); err != nil {
			slog.Error("cannot unmarshal parameters", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var id int
		idRow := metaDb.QueryRow(request.Context(), "INSERT INTO missing_metrics (build_type, project, metric, missing_since) VALUES ($1, $2, $3, $4) RETURNING id", params.BuildType, params.Project, params.Metric, time.UnixMilli(params.MissingSince))
		if err = idRow.Scan(&id); err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				http.Error(writer, "Conflict: Accident already exists", http.StatusConflict)
			} else {
				slog.Error("cannot insert missing data", "error", err, "params", params)
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
