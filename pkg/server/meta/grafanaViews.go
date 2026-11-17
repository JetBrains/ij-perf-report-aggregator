package meta

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GrafanaViewRequest is sent by every intellij-quality-infrastructure-dashboards Grafana dashboard on load.
type GrafanaViewRequest struct {
	Dashboard string `json:"dashboard"`
	User      string `json:"user"`
}

type GrafanaUniqueUsers struct {
	UniqueUsers int `json:"uniqueUsers"`
}

func CreatePostGrafanaViewHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req GrafanaViewRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		if err := decoder.Decode(&req); err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		dashboard := strings.TrimSpace(req.Dashboard)
		user := strings.TrimSpace(req.User)
		// An empty or uninterpolated `${__user.login}` means Grafana had no user to substitute.
		if dashboard == "" || user == "" || strings.HasPrefix(dashboard, "$") || strings.HasPrefix(user, "$") {
			http.Error(writer, "dashboard and user are required", http.StatusBadRequest)
			return
		}

		const upsertSQL = "INSERT INTO grafana_views (dashboard_uid, user_login) VALUES ($1, $2) " +
			"ON CONFLICT (dashboard_uid, user_login, day) DO UPDATE SET views = grafana_views.views + 1"
		if _, err := metaDb.Exec(request.Context(), upsertSQL, dashboard, user); err != nil {
			slog.Error("cannot record grafana view", "error", err, "dashboard", dashboard)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Answer with JSON so the Infinity variable query that sends this has something to parse.
		writer.Header().Set("Content-Type", "application/json")
		if _, err := writer.Write([]byte("[]")); err != nil {
			slog.Error("cannot write response", "error", err)
		}
	}
}

// CreateGetGrafanaUniqueUsersHandler serves GET /api/meta/grafanaViews/uniqueUsers?days=7: it counts distinct users over the last `days` days, today included.
func CreateGetGrafanaUniqueUsersHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		days := 7
		if value := request.URL.Query().Get("days"); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed <= 0 {
				http.Error(writer, "days must be a positive integer", http.StatusBadRequest)
				return
			}
			days = parsed
		}

		var result GrafanaUniqueUsers
		const sql = "SELECT count(DISTINCT user_login) FROM grafana_views WHERE day > CURRENT_DATE - $1::int"
		if err := metaDb.QueryRow(request.Context(), sql, days).Scan(&result.UniqueUsers); err != nil {
			slog.Error("cannot count grafana unique users", "error", err, "days", days)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			slog.Error("cannot write grafana unique users response", "error", err)
		}
	}
}
