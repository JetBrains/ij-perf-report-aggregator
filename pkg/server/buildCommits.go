package server

import (
	"context"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/installer"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
)

// The installer tables of all databases hold the changes of every build by TeamCity build id and outlive TeamCity
// build retention. A build may be in several databases, with the same changes.
const installerTables = "merge(REGEXP('.'), '^installer$')"

// buildStoreConn returns one connection for all build lookups: a clickhouse-go connection is a pool safe for concurrent
// use, so a request doing several lookups doesn't pay a handshake for each.
func (t *StatsServer) buildStoreConn() (driver.Conn, error) {
	t.buildsConnOnce.Do(func() {
		t.buildsConn, t.buildsConnErr = t.openDatabaseConnection()
	})
	return t.buildsConn, t.buildsConnErr
}

// BuildCommits returns the commits of the build's range.
func (t *StatsServer) BuildCommits(ctx context.Context, buildId string) ([]string, error) {
	ids := parseBuildIds([]string{buildId})
	if len(ids) == 0 {
		return nil, nil
	}
	conn, err := t.buildStoreConn()
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, "SELECT changes FROM "+installerTables+" WHERE id = ? ORDER BY length(changes) DESC LIMIT 1", ids[0])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commits []string
	for rows.Next() {
		var changes []string
		if err := rows.Scan(&changes); err != nil {
			return nil, err
		}
		commits = append(commits, installer.DecodeChanges(changes)...)
	}
	return commits, rows.Err()
}

// MatchBuildCommits intersects the commits with the changes of the builds in ClickHouse, so only matched commits are
// sent back, not whole build ranges.
func (t *StatsServer) MatchBuildCommits(ctx context.Context, buildIds []string, commits []string) (map[string]meta.BuildCommitsMatch, error) {
	result := make(map[string]meta.BuildCommitsMatch)
	ids := parseBuildIds(buildIds)
	encodedToCommit := make(map[string]string, len(commits))
	encoded := make([]string, 0, len(commits))
	for _, commit := range commits {
		if e, ok := installer.EncodeCommit(commit); ok {
			encodedToCommit[e] = commit
			encoded = append(encoded, e)
		}
	}
	if len(ids) == 0 || len(encoded) == 0 {
		return result, nil
	}

	conn, err := t.buildStoreConn()
	if err != nil {
		return nil, err
	}

	const sql = "SELECT id, max(length(changes)), groupUniqArrayArray(arrayIntersect(changes, ?)) " +
		"FROM " + installerTables + " WHERE id IN ? AND hasAny(changes, ?) GROUP BY id"
	rows, err := conn.Query(ctx, sql, encoded, ids, encoded)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id      uint32
			size    uint64
			matched []string
		)
		if err := rows.Scan(&id, &size, &matched); err != nil {
			return nil, err
		}
		match := meta.BuildCommitsMatch{RangeSize: int(size)}
		for _, e := range matched {
			match.MatchedCommits = append(match.MatchedCommits, encodedToCommit[e])
		}
		result[strconv.FormatUint(uint64(id), 10)] = match
	}
	return result, rows.Err()
}

// BuildRuns looks the runs up in a report table; db and table come from the request, so they must be plain identifiers.
func (t *StatsServer) BuildRuns(ctx context.Context, db string, table string, buildIds []string, projects []string, from time.Time, to time.Time) (map[meta.BuildProject]meta.BuildRun, error) {
	if err := sql_util.ValidateIdentifier("db", db); err != nil {
		return nil, err
	}
	if err := sql_util.ValidateIdentifier("table", table); err != nil {
		return nil, err
	}
	result := make(map[meta.BuildProject]meta.BuildRun)
	ids := parseBuildIds(buildIds)
	if len(ids) == 0 || len(projects) == 0 {
		return result, nil
	}

	conn, err := t.buildStoreConn()
	if err != nil {
		return nil, err
	}

	sql := "SELECT tc_build_id, project, any(machine), any(branch), min(generated_time) FROM " + db + "." + table + " WHERE tc_build_id IN ? AND project IN ?"
	args := []any{ids, projects}
	if !from.IsZero() && !to.IsZero() {
		// tc_build_id is not in the sorting key, the time bounds let ClickHouse skip partitions
		sql += " AND generated_time BETWEEN ? AND ?"
		args = append(args, from, to)
	}
	rows, err := conn.Query(ctx, sql+" GROUP BY tc_build_id, project", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id      uint32
			project string
			run     meta.BuildRun
		)
		if err := rows.Scan(&id, &project, &run.Machine, &run.Branch, &run.GeneratedTime); err != nil {
			return nil, err
		}
		result[meta.BuildProject{BuildId: strconv.FormatUint(uint64(id), 10), Project: project}] = run
	}
	return result, rows.Err()
}

// parseBuildIds keeps the TeamCity build ids, dropping anything that is not one
func parseBuildIds(buildIds []string) []uint32 {
	ids := make([]uint32, 0, len(buildIds))
	for _, id := range buildIds {
		if parsed, err := strconv.ParseUint(id, 10, 32); err == nil {
			ids = append(ids, uint32(parsed))
		}
	}
	return ids
}
