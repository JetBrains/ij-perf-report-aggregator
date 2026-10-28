package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

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
	Link         string `json:"link"`
}

type degradationsRequestBody struct {
	Owners            []string `json:"owners"`
	AdditionalMetrics []string `json:"additionalMetrics"`
}

type degradationsParams struct {
	owners            []string
	from              string
	to                string
	includeInferred   bool
	additionalMetrics []string
}

var errDegradationsMissingParams = errors.New(`"owner" (query parameter, or "owners" in the POST body), "from" and "to" are all required`)

func dropBlank(values []string) []string {
	return slices.DeleteFunc(values, func(s string) bool { return s == "" })
}

// parseDegradationsRequest reads the query parameters and, for POST, merges in the optional JSON body.
func parseDegradationsRequest(r *http.Request) (degradationsParams, error) {
	q := r.URL.Query()
	params := degradationsParams{
		owners:            q["owner"],
		from:              q.Get("from"),
		to:                q.Get("to"),
		includeInferred:   q.Get("inferredEvents") == "true",
		additionalMetrics: q["additionalMetric"],
	}

	if r.Method == http.MethodPost && r.Body != nil {
		var body degradationsRequestBody
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&body)
		if err != nil && !errors.Is(err, io.EOF) {
			return degradationsParams{}, fmt.Errorf("invalid request body: %w", err)
		}
		params.owners = slices.Concat(params.owners, body.Owners)
		params.additionalMetrics = slices.Concat(params.additionalMetrics, body.AdditionalMetrics)
	}

	params.owners = dropBlank(params.owners)
	params.additionalMetrics = dropBlank(params.additionalMetrics)

	if len(params.owners) == 0 || params.from == "" || params.to == "" {
		return degradationsParams{}, errDegradationsMissingParams
	}
	return params, nil
}

// CreateGetDegradationsHandler returns a handler (served for GET and POST) that fetches pre-computed
// degradations (accidents) from the meta DB for one or more owners and a date range.
//
// Query parameters:
//
//	owner            – (required unless given in the body, repeatable) matched against project_owner.owner
//	from             – (required) start of date range, inclusive (ISO date, e.g. "2024-01-01")
//	to               – (required) end of date range, inclusive (ISO date, e.g. "2024-03-31")
//	inferredEvents   – (optional) include R2D2-detected regressions/improvements; default false
//	additionalMetric – (optional, repeatable) extra metric names beyond the default mainMetrics list
//
// POST additionally accepts an optional JSON body {"owners": [...], "additionalMetrics": [...]} whose
// values are merged with the query parameters; use it when the owner list is too long for a URL.
//
// Example:
//
//	GET /api/meta/degradations?owner=Java&owner=IntelliJ+Kotlin+Plugin&from=2024-01-01&to=2024-03-31&inferredEvents=true&additionalMetric=myMetric
//	POST /api/meta/degradations?from=2024-01-01&to=2024-03-31 {"owners": ["Java", "IntelliJ Kotlin Plugin"]}
func CreateGetDegradationsHandler(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := parseDegradationsRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		owners, from, to := params.owners, params.from, params.to
		includeInferred, additionalMetrics := params.includeInferred, params.additionalMetrics

		allowedMetrics := slices.Concat(mainMetrics, additionalMetrics)

		const query = `
			SELECT ac.date,
			       ac.affected_test,
			       ac.build_number,
			       ac.reason,
			       ac.kind,
			       ac.user_name,
			       po.db_name,
			       po.table_name,
			       po.project
			FROM   accidents ac
			JOIN LATERAL (
			    SELECT db_name, table_name, project
			    FROM   project_owner
			    WHERE  owner = ANY($1)
			      AND  (ac.affected_test = project OR ac.affected_test LIKE project || '/%')
			    ORDER BY LENGTH(project) DESC
			    LIMIT  1
			) po ON true
			WHERE  ac.date      >= $2::date
			  AND  ac.date      <= $3::date
			  AND  ($4 OR ac.kind NOT IN ('InferredRegression', 'InferredImprovement'))
			ORDER BY ac.date DESC, ac.affected_test, ac.build_number`

		rows, err := metaDb.Query(r.Context(), query, owners, from, to, includeInferred)
		if err != nil {
			slog.Error("degradations: query failed", "error", err, "owners", owners, "from", from, "to", to)
			http.Error(w, "query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var entries []degradationReportEntry
		for rows.Next() {
			var (
				date                                                      pgtype.Date
				affectedTest, buildNumber, reason, kind, dbName, userName string
				tableName, poProject                                      string
			)
			if err := rows.Scan(&date, &affectedTest, &buildNumber, &reason, &kind, &userName, &dbName, &tableName, &poProject); err != nil {
				slog.Error("degradations: failed to scan row", "error", err)
				http.Error(w, "failed to read results", http.StatusInternalServerError)
				return
			}

			// Derive the metric: if affected_test has po.project as a prefix, the remainder is the metric.
			metric, hasMetric := strings.CutPrefix(affectedTest, poProject+"/")
			if !hasMetric {
				metric = ""
			}

			// When a metric is present, skip entries not in the allowed metrics list.
			if metric != "" && !slices.Contains(allowedMetrics, metric) {
				continue
			}

			fromDate := date.Time.AddDate(0, 0, -14)
			toDate := date.Time.AddDate(0, 0, 14)
			entries = append(entries, degradationReportEntry{
				Date:         date.Time.Format("2006-01-02"),
				AffectedTest: affectedTest,
				BuildNumber:  buildNumber,
				Reason:       reason,
				Kind:         kind,
				UserName:     userName,
				Link:         buildDegradationLink(poProject, metric, reason, dbName, tableName, fromDate, toDate),
			})
		}
		if err := rows.Err(); err != nil {
			slog.Error("degradations: failed to iterate rows", "error", err)
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

var ticketRegex = regexp.MustCompile(`(?:^|[\s/])([A-Z]+-\d+)`)

func buildDegradationLink(project, metric, reason, dbName, tableName string, fromDate, toDate time.Time) string {
	if metric == "" {
		if m := ticketRegex.FindStringSubmatch(reason); m != nil {
			return "https://youtrack.jetbrains.com/issue/" + m[1]
		}
	}

	q := url.Values{}
	q.Set("dbName", dbName)
	q.Set("table", tableName)
	q.Set("project", project)
	if metric != "" {
		q.Set("measure", metric)
	}
	q.Set("timeRange", "custom")
	q.Set("customRange", fromDate.Format("2006-1-2")+":"+toDate.Format("2006-1-2"))
	return ijPerfBaseURL + "/owners/test?" + q.Encode()
}
