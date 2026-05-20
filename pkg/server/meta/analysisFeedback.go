package meta

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalysisFeedback struct {
	Id         int     `json:"id"`
	AnalysisId int     `json:"analysisId"`
	Rate       int     `json:"rate"`
	Feedback   *string `json:"feedback,omitempty"`
	UserEmail  *string `json:"userEmail,omitempty"`
	CreatedAt  string  `json:"createdAt"`
}

type AnalysisFeedbackRequest struct {
	Rate     int     `json:"rate"`
	Feedback *string `json:"feedback,omitempty"`
}

func CreatePostAnalysisFeedback(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}

		var req AnalysisFeedbackRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		if err := decoder.Decode(&req); err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if req.Rate < 1 || req.Rate > 5 {
			http.Error(writer, "rate must be between 1 and 5", http.StatusBadRequest)
			return
		}

		var feedbackArg *string
		if req.Feedback != nil {
			trimmed := strings.TrimSpace(*req.Feedback)
			if trimmed != "" {
				feedbackArg = &trimmed
			}
		}

		userEmail := request.Header.Get("X-Auth-Request-Email")
		var userEmailArg *string
		if userEmail != "" {
			userEmailArg = &userEmail
		}

		_, err = metaDb.Exec(request.Context(),
			"INSERT INTO analysis_feedback (analysis_id, rate, feedback, user_email) VALUES ($1, $2, $3, $4)",
			id, req.Rate, feedbackArg, userEmailArg)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
				http.Error(writer, "analysis not found", http.StatusNotFound)
				return
			}
			slog.Error("cannot execute insert analysis_feedback query", "error", err, "analysisId", id)
			http.Error(writer, "Failed to insert analysis feedback: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func CreateGetAnalysisFeedback(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}

		const sql = "SELECT id, analysis_id, rate, feedback, user_email, created_at " +
			"FROM analysis_feedback WHERE analysis_id = $1 ORDER BY created_at DESC"
		rows, err := metaDb.Query(request.Context(), sql, id)
		if err != nil {
			slog.Error("unable to execute select analysis_feedback query", "error", err, "analysisId", id)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		feedbacks, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (AnalysisFeedback, error) {
			var fb AnalysisFeedback
			var rate int16
			var createdAt time.Time
			if err := row.Scan(&fb.Id, &fb.AnalysisId, &rate, &fb.Feedback, &fb.UserEmail, &createdAt); err != nil {
				return AnalysisFeedback{}, err
			}
			fb.Rate = int(rate)
			fb.CreatedAt = createdAt.Format(time.RFC3339)
			return fb, nil
		})
		if err != nil {
			slog.Error("unable to collect analysis_feedback rows", "error", err, "analysisId", id)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if feedbacks == nil {
			feedbacks = []AnalysisFeedback{}
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(feedbacks); err != nil {
			slog.Error("unable to write analysis_feedback response", "error", err, "analysisId", id)
		}
	}
}
