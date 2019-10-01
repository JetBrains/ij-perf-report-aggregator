package server

import (
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "net/http"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
)

func (t *StatsServer) handleMetricsRequest(w http.ResponseWriter, request *http.Request) {
  query := request.URL.Query()
  product := query.Get("product")
  if len(product) == 0 {
    http.Error(w, `{"error": "product parameter is required"}`, 400)
    return
  }

  machine := query.Get("machine")
  if len(product) == 0 {
    http.Error(w, `{"error": "machine parameter is required"}`, 400)
    return
  }

	//noinspection SqlResolve
	statement, err := t.db.Prepare(`select generated_time, build_c1, build_c2, build_c3, metrics from report where product = ? and machine = ? order by build_c1, build_c2, build_c3, generated_time`, product, machine)
	if err != nil {
		t.logger.Error("cannot query", zap.Error(err))
		t.httpError(err, w)
		return
	}

	defer util.Close(statement, t.logger)

	w.Header().Set("Content-Type", "application/json")
	jsonWriter := jsoniter.NewStream(jsoniter.ConfigFastest, w, 64*1024)

	jsonWriter.WriteArrayStart()

  var stringBuilder strings.Builder

	isFirst := true
	lastBuildWithoutUniqueSuffix := ""
	for {
		hasRow, err := statement.Step()
		if err != nil {
			t.httpError(err, w)
			return
		}

		if !hasRow {
			break
		}

		var generatedTime int64
		var buildC1 int
		var buildC2 int
		var buildC3 int
		var metrics sqlite3.RawString
		err = statement.Scan(&generatedTime, &buildC1, &buildC2, &buildC3, &metrics)
		if err != nil {
			t.httpError(err, w)
			return
		}

		if isFirst {
			isFirst = false
		} else {
			jsonWriter.WriteMore()
		}

		jsonWriter.WriteObjectStart()
		// timestamp
		jsonWriter.WriteObjectField("t")
		// seconds to milliseconds
		jsonWriter.WriteInt64(generatedTime * 1000)
		jsonWriter.WriteMore()

		// build number with addition if not unique (build as x coordinate)
    stringBuilder.Reset()
    stringBuilder.WriteString(strconv.Itoa(buildC1))
    stringBuilder.WriteRune('.')
    stringBuilder.WriteString(strconv.Itoa(buildC2))
		if buildC3 != 0 {
      stringBuilder.WriteRune('.')
      stringBuilder.WriteString(strconv.Itoa(buildC3))
    }

    buildAsString := stringBuilder.String()
    // https://www.amcharts.com/docs/v4/tutorials/handling-repeating-categories-on-category-axis/
    if lastBuildWithoutUniqueSuffix == buildAsString {
      // not unique - add time
      stringBuilder.WriteRune(' ')
      stringBuilder.WriteRune('(')
      //stringBuilder.WriteString(time.Unix(generatedTime, 0).Format(time.UnixDate))
      stringBuilder.WriteString(strconv.FormatInt(generatedTime, 10))
      stringBuilder.WriteRune(')')
      buildAsString = stringBuilder.String()
    } else {
      lastBuildWithoutUniqueSuffix = buildAsString
    }

    jsonWriter.WriteObjectField("build")
    jsonWriter.WriteString(buildAsString)
    jsonWriter.WriteMore()

		// skip first '{'
		jsonWriter.WriteRaw(string(metrics[1:]))
	}

	jsonWriter.WriteArrayEnd()

	err = jsonWriter.Flush()
	if err != nil {
		t.logger.Error("cannot flush", zap.Error(err))
		return
	}
}

