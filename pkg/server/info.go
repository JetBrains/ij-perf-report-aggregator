package server

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/json-iterator/go"
	"report-aggregator/pkg/util"
)

func (t *StatsServer) getProductNames() ([]string, error) {
	statement, err := t.db.Prepare("select distinct product from report order by product")
	if err != nil {
		return nil, err
	}

	defer util.Close(statement, t.logger)

	return t.readStringList(statement)
}

func (t *StatsServer) readStringList(statement *sqlite3.Stmt) ([]string, error) {
	var result []string
	for {
		hasRow, err := statement.Step()
		if err != nil {
			return nil, err
		}

		value, _, err := statement.ColumnText(0)
		if err != nil {
			return nil, err
		}

		if !hasRow {
			break
		}

		result = append(result, value)
	}
	return result, nil
}

func writeStringList(w *jsoniter.Stream, statement *sqlite3.Stmt) error {
	isFirst := true
	for {
		hasRow, err := statement.Step()
		if err != nil {
			return err
		}

		value, _, err := statement.ColumnRawString(0)
		if err != nil {
			return err
		}

		if !hasRow {
			break
		}

		if isFirst {
			isFirst = false
		} else {
			w.WriteMore()
		}
		w.WriteString(string(value))
	}
	return nil
}
