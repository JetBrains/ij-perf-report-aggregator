package sql

import (
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "report-aggregator/pkg/util"
)

func GetInt(db *sqlite3.Conn, query string, logger *zap.Logger) (int, error) {
	statement, err := db.Prepare(query)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	defer util.Close(statement, logger)

	_, err = statement.Step()
	if err != nil {
		return -1, errors.WithStack(err)
	}

	value, _, err := statement.ColumnInt(0)
	return value, errors.WithStack(err)
}
