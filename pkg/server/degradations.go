package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type degradationReportEntry struct {
	Date         string `json:"date"`
	AffectedTest string `json:"affectedTest"`
	BuildNumber  string `json:"buildNumber"`
	Reason       string `json:"reason"`
	Kind         string `json:"kind"`
	UserName     string `json:"userName"`
}

// CreateGetDegradationsHandler returns a GET handler that fetches pre-computed
// degradations (accidents) from the meta DB for one or more owners and a date range.
//
// Query parameters (all required):
//
//	owner  – one or more values matched against project_owner.owner (repeat param for multiple)
//	from   – start of date range, inclusive (ISO date, e.g. "2024-01-01")
//	to     – end   of date range, inclusive (ISO date, e.g. "2024-03-31")
//
// Example:
//
//	GET /api/meta/degradations?owner=Java&owner=IntelliJ+Kotlin+Plugin&from=2024-01-01&to=2024-03-31
func CreateGetDegradationsHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owners := slices.DeleteFunc(r.URL.Query()["owner"], func(s string) bool { return s == "" })
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if len(owners) == 0 || from == "" || to == "" {
			http.Error(w, `query parameters "owner", "from", and "to" are all required`, http.StatusBadRequest)
			return
		}

		const query = `
			SELECT ac.date,
			       ac.affected_test,
			       ac.build_number,
			       ac.reason,
			       ac.kind,
			       ac.user_name
			FROM   accidents ac
			JOIN   project_owner po ON ac.affected_test = po.project
			WHERE  po.owner      = ANY($1)
			  AND  ac.date      >= $2::date
			  AND  ac.date      <= $3::date
			ORDER BY ac.date DESC, ac.affected_test, ac.build_number`

		rows, err := metaDb.Query(r.Context(), query, owners, from, to)
		if err != nil {
			slog.Error("degradations: query failed", "error", err, "owners", owners, "from", from, "to", to)
			http.Error(w, "query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (degradationReportEntry, error) {
			var (
				date                                    pgtype.Date
				affectedTest, buildNumber, reason, kind string
				userName                                string
			)
			err := row.Scan(&date, &affectedTest, &buildNumber, &reason, &kind, &userName)
			return degradationReportEntry{
				Date:         date.Time.Format("2006-01-02"),
				AffectedTest: affectedTest,
				BuildNumber:  buildNumber,
				Reason:       reason,
				Kind:         kind,
				UserName:     userName,
			}, err
		})
		if err != nil {
			slog.Error("degradations: failed to collect rows", "error", err)
			http.Error(w, "failed to read results", http.StatusInternalServerError)
			return
		}

		if entries == nil {
			entries = []degradationReportEntry{}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			slog.Error("degradations: failed to encode response", "error", err)
		}
	}
}
