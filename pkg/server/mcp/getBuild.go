package mcp

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const teamCityBaseURL = "https://buildserver.labs.intellij.net"

type getBuildInput struct {
	BuildID  int64  `json:"tc_build_id"        jsonschema:"TeamCity build id (the same value returned by search_metric_values)"`
	Database string `json:"database,omitempty" jsonschema:"Optional database to restrict the scan to"`
	Table    string `json:"table,omitempty"    jsonschema:"Optional table to restrict the scan to"`
}

type buildProject struct {
	Database string `json:"database"`
	Table    string `json:"table"`
	Project  string `json:"project"`
}

type getBuildOutput struct {
	BuildID     int64          `json:"tc_build_id"`
	Branch      string         `json:"branch"`
	BuildTime   string         `json:"build_time"`
	Machine     string         `json:"machine"`
	TeamCityURL string         `json:"teamcity_url"`
	FirstCommit string         `json:"first_commit,omitempty" jsonschema:"Oldest git commit covered by this build's installer (short hex SHA). Empty if no source table is linked to an installer."`
	LastCommit  string         `json:"last_commit,omitempty"  jsonschema:"Newest git commit covered by this build's installer (short hex SHA). Empty if no source table is linked to an installer."`
	Projects    []buildProject `json:"projects"               jsonschema:"Distinct (database, table, project) tuples that this build produced data for."`
	Count       int            `json:"count"                  jsonschema:"Number of project entries"`
}

func (s *service) getBuild(ctx context.Context, _ *sdk.CallToolRequest, in getBuildInput) (*sdk.CallToolResult, getBuildOutput, error) {
	if in.BuildID <= 0 {
		return nil, getBuildOutput{}, errors.New("tc_build_id is required and must be positive")
	}
	tables, err := s.resolveTables(ctx, in.Database, in.Table)
	if err != nil {
		return nil, getBuildOutput{}, err
	}

	perTable := func(r tableRef) (string, []any) {
		// Fall back to generated_time when build_time is missing (column absent)
		// or unset (epoch-zero — some ingestion paths leave it 0).
		buildTimeExpr := "toString(generated_time) as bld_time"
		if r.HasBuildTime {
			buildTimeExpr = "toString(if(toUnixTimestamp(build_time) = 0, generated_time, build_time)) as bld_time"
		}
		// Only tables with tc_installer_build_id can join the per-database `installer`
		// table. Others get an empty Array(String) so UNION ALL columns line up.
		commitsExpr := "[]::Array(String) as installer_changes"
		fromExpr := fmt.Sprintf("from %s.%s", r.Database, r.Table)
		whereExpr := "where tc_build_id = ?"
		args := []any{r.Database, r.Table, in.BuildID}
		if r.HasInstallerID {
			// Restrict the right side to only the installer rows referenced by this build,
			// otherwise the join blows the 3 GiB query memory limit.
			commitsExpr = "arrayMap(c -> toString(c), i.changes) as installer_changes"
			fromExpr = fmt.Sprintf(
				"from %s.%s r left join (select id, changes from %s.installer where id in "+
					"(select tc_installer_build_id from %s.%s where tc_build_id = ?)) i "+
					"on r.tc_installer_build_id = i.id",
				r.Database, r.Table, r.Database, r.Database, r.Table)
			whereExpr = "where r.tc_build_id = ?"
			args = []any{r.Database, r.Table, in.BuildID, in.BuildID}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb,
			"select ? as db_name, ? as table_name, "+
				"toString(project) as project_name, toString(branch) as branch_name, "+
				"toString(machine) as machine_name, "+
				"%s, toString(generated_time) as gen_time, "+
				"%s "+
				"%s %s",
			buildTimeExpr, commitsExpr, fromExpr, whereExpr)
		return sb.String(), args
	}

	innerSQL, args := buildUnion(tables, perTable)
	// Order by gen_time desc (still selected internally) so the most recent row "wins"
	// when picking the root-level branch/build_time/machine/commit values.
	sql := "select db_name, table_name, project_name, branch_name, machine_name, " +
		"bld_time, installer_changes " +
		"from (" + innerSQL + ") as u order by gen_time desc"

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, getBuildOutput{}, fmt.Errorf("get_build: %w", err)
	}
	defer rows.Close()

	out := getBuildOutput{
		BuildID:     in.BuildID,
		TeamCityURL: fmt.Sprintf("%s/viewLog.html?buildId=%d", teamCityBaseURL, in.BuildID),
		Projects:    []buildProject{},
	}
	type projectKey struct{ db, table, project string }
	seenProject := make(map[projectKey]struct{})
	rootSet := false

	for rows.Next() {
		var dbName, tableName, project, branch, machine, bldTime string
		var changes []string
		if err := rows.Scan(
			&dbName, &tableName,
			&project, &branch, &machine,
			&bldTime,
			&changes,
		); err != nil {
			return nil, getBuildOutput{}, fmt.Errorf("scan: %w", err)
		}

		// Take root metadata from the first row (latest by generated_time).
		if !rootSet {
			out.Branch = branch
			out.BuildTime = bldTime
			out.Machine = machine
			out.FirstCommit, out.LastCommit = commitRange(changes)
			rootSet = true
		} else if out.FirstCommit == "" && len(changes) > 0 {
			// Some report tables for the same build may not link an installer; pick up
			// commits from any row that does carry them.
			out.FirstCommit, out.LastCommit = commitRange(changes)
		}

		key := projectKey{dbName, tableName, project}
		if _, dup := seenProject[key]; dup {
			continue
		}
		seenProject[key] = struct{}{}
		out.Projects = append(out.Projects, buildProject{
			Database: dbName,
			Table:    tableName,
			Project:  project,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, getBuildOutput{}, fmt.Errorf("rows: %w", err)
	}

	// Some report tables don't have a tc_installer_build_id column at all (e.g.
	// perfintDev.kotlin). For those builds the commits are still stored — under the
	// build's own id — in the per-database installer table. Fall back to a direct
	// lookup by tc_build_id across the databases we saw.
	if out.FirstCommit == "" {
		seenDB := map[string]struct{}{}
		for _, p := range out.Projects {
			seenDB[p.Database] = struct{}{}
		}
		for db := range seenDB {
			if changes := s.fetchInstallerChanges(ctx, db, in.BuildID); len(changes) > 0 {
				out.FirstCommit, out.LastCommit = commitRange(changes)
				break
			}
		}
	}

	out.Count = len(out.Projects)
	return nil, out, nil
}

// fetchInstallerChanges looks up the installer.changes array for the given build id
// in <db>.installer. Returns nil on any error (table missing, no row, etc.) — the
// caller treats absent commits as "this build doesn't expose them".
func (s *service) fetchInstallerChanges(ctx context.Context, db string, buildID int64) []string {
	if err := validateIdentifier("database", db); err != nil {
		slog.Warn("mcp: skipping installer fallback, invalid database identifier", "db", db, "err", err)
		return nil
	}
	sql := fmt.Sprintf(
		"select arrayMap(c -> toString(c), changes) from %s.installer where id = ? limit 1",
		db)
	rows, err := s.db.Query(ctx, sql, buildID)
	if err != nil {
		slog.Warn("mcp: installer fallback query failed", "db", db, "build_id", buildID, "err", err)
		return nil
	}
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	var changes []string
	if err := rows.Scan(&changes); err != nil {
		slog.Warn("mcp: installer fallback scan failed", "db", db, "build_id", buildID, "err", err)
		return nil
	}
	return changes
}

// shortCommitLen matches git's default short SHA length — long enough to be
// unambiguous in IntelliJ-sized repos while keeping LLM token cost low.
const shortCommitLen = 12

// commitRange returns the (oldest, newest) commit pair from an installer changes
// array. The collector stores commits newest-first as base64 raw-std SHA-1 bytes;
// we decode to hex and truncate. A non-base64 input is returned as-is so a future
// format switch (e.g. raw hex) doesn't silently corrupt output.
func commitRange(changes []string) (string, string) {
	if len(changes) == 0 {
		return "", ""
	}
	decode := func(s string) string {
		if s == "" {
			return ""
		}
		b, err := base64.RawStdEncoding.DecodeString(s)
		if err != nil {
			return s
		}
		return hex.EncodeToString(b)[:min(len(b)*2, shortCommitLen)]
	}
	return decode(changes[len(changes)-1]), decode(changes[0])
}
