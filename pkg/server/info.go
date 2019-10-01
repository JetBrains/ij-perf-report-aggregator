package server

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/json-iterator/go"
  "net/http"
  "report-aggregator/pkg/util"
)

func (t *StatsServer) handleInfoRequest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonWriter := jsoniter.NewStream(jsoniter.ConfigFastest, w, 64*1024)

	jsonWriter.WriteObjectStart()
	jsonWriter.WriteObjectField("productNames")

	productNames, err := t.getProductNames()
	if err != nil {
		t.httpError(err, w)
		return
	}

	jsonWriter.WriteArrayStart()
	isFirst := true
	for _, product := range productNames {
		if isFirst {
			isFirst = false
		} else {
			jsonWriter.WriteMore()
		}
		jsonWriter.WriteString(product)
	}
	jsonWriter.WriteArrayEnd()

	jsonWriter.WriteMore()
  err = t.writeProductToMachineNames(jsonWriter, productNames)
  if err != nil {
    t.httpError(err, w)
    return
  }

  jsonWriter.WriteMore()
  writeDurationMetrics(jsonWriter)

  jsonWriter.WriteMore()
  writeInstantMetrics(jsonWriter)

  jsonWriter.WriteObjectEnd()

  err = jsonWriter.Flush()
  if err != nil {
    t.httpError(err, w)
    return
  }
}

func (t *StatsServer) writeProductToMachineNames(jsonWriter *jsoniter.Stream, productNames []string) error {
  jsonWriter.WriteObjectField("productToMachineNames")
  jsonWriter.WriteObjectStart()

  statement, err := t.db.Prepare("select distinct machine from report where product = ? order by machine")
  if err != nil {
    return err
  }

  defer util.Close(statement, t.logger)

  isFirst := true
  for _, product := range productNames {
    err = statement.Bind(product)
    if err != nil {
      return err
    }

    if isFirst {
      isFirst = false
    } else {
      jsonWriter.WriteMore()
    }
    jsonWriter.WriteObjectField(product)

    jsonWriter.WriteArrayStart()
    err = writeStringList(jsonWriter, statement)
    if err != nil {
      return err
    }
    jsonWriter.WriteArrayEnd()
  }

  jsonWriter.WriteObjectEnd()
  return nil
}

func writeDurationMetrics(w *jsoniter.Stream) {
  w.WriteObjectField("durationMetricsNames")
  w.WriteArrayStart()

  w.WriteString("bootstrap")
  w.WriteMore()
  w.WriteString("appInitPreparation")
  w.WriteMore()
  w.WriteString("appInit")
  w.WriteMore()
  w.WriteString("pluginDescriptorLoading")
  w.WriteMore()
  w.WriteString("appComponentCreation")
  w.WriteMore()
  w.WriteString("projectComponentCreation")
  w.WriteMore()
  w.WriteString("moduleLoading")

  w.WriteArrayEnd()
}

func writeInstantMetrics(w *jsoniter.Stream) {
  w.WriteObjectField("instantMetricsNames")
  w.WriteArrayStart()

  w.WriteString("splash")

  w.WriteArrayEnd()
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
