package server

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"net/http"
	"report-aggregator/pkg/util"
)

var essentialMetricNames = []string{"bootstrap", "appInitPreparation", "appInit", "pluginDescriptorLoading", "appComponentCreation", "projectComponentCreation"}

func (t *StatsServer) handleInfoRequest(w http.ResponseWriter, _ *http.Request) {
	productNames, err := t.getProductNames()
	if err != nil {
		t.httpError(err, w)
		return
	}

	statement, err := t.db.Prepare("select rowid as id, name from machine where rowid in (select distinct machine from report where product = ?) order by name")
	if err != nil {
		t.httpError(err, w)
		return
	}

	defer util.Close(statement, t.logger)

	w.Header().Set("Content-Type", "application/json")
	var errRef error
	WriteInfo(w, productNames, essentialMetricNames, statement, &errRef)

	if errRef != nil {
		t.httpError(err, w)
		return
	}
}

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
