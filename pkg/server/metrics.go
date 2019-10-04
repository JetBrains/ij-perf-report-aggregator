package server

import (
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/json-iterator/go"
  "github.com/pkg/errors"
  "github.com/valyala/quicktemplate"
  "io"
  "net/http"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
)

func (t *StatsServer) handleMetricsRequest(request *http.Request) ([]byte, error) {
  query := request.URL.Query()
  product := query.Get("product")
  if len(product) == 0 {
    return nil, NewHttpError(400, "product parameter is required")
  }

  machine := query.Get("machine")
  if len(product) == 0 {
    return nil, NewHttpError(400, "machine parameter is required")
  }

  buffer := quicktemplate.AcquireByteBuffer()
  defer quicktemplate.ReleaseByteBuffer(buffer)
  err := t.computeMetricsResponse(product, machine, buffer)
  if err != nil {
    return nil, err
  }
  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result, nil
}

func (t *StatsServer) computeMetricsResponse(product string, machine string, writer io.Writer) error {
  statement, err := t.db.Prepare(`
select generated_time, build_c1, build_c2, build_c3, duration_metrics, instant_metrics 
from report 
where product = ? and machine = ? 
order by build_c1, build_c2, build_c3, generated_time`, product, machine)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(statement, t.logger)

  jsonWriter := jsoniter.NewStream(jsoniter.ConfigFastest, writer, 8*1024)

  jsonWriter.WriteArrayStart()

  var stringBuilder strings.Builder

  isFirst := true
  lastBuildWithoutUniqueSuffix := ""
  for {
    hasRow, err := statement.Step()
    if err != nil {
      return errors.WithStack(err)
    }

    if !hasRow {
      break
    }

    var generatedTime int64
    var buildC1 int
    var buildC2 int
    var buildC3 int
    var durationMetrics sqlite3.RawString
    var instantMetrics sqlite3.RawString
    err = statement.Scan(&generatedTime, &buildC1, &buildC2, &buildC3, &durationMetrics, &instantMetrics)
    if err != nil {
      return errors.WithStack(err)
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
    jsonWriter.WriteRaw(string(durationMetrics[1 : len(durationMetrics)-1]))
    jsonWriter.WriteMore()
    jsonWriter.WriteRaw(string(instantMetrics[1:]))
  }

  jsonWriter.WriteArrayEnd()

  return jsonWriter.Flush()
}
