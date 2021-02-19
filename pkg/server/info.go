package server

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/valyala/quicktemplate"
  "net/http"
)

func (t *StatsServer) handleMetaMeasureRequest( _ *http.Request) ([]byte, error) {
  buffer := byteBufferPool.Get()
  defer byteBufferPool.Put(buffer)

  measureNames := make([]string, len(analyzer.IjMetricDescriptors))
  for index, descriptor := range analyzer.IjMetricDescriptors {
    measureNames[index] = descriptor.Name
  }

  templateWriter := quicktemplate.AcquireWriter(buffer)
  defer quicktemplate.ReleaseWriter(templateWriter)
  jsonWriter := templateWriter.N()
  jsonWriter.S("[")
  for index, name := range measureNames {
    if index != 0 {
      jsonWriter.S(",")
    }
    jsonWriter.Q(name)
  }
  jsonWriter.S("]")

  return CopyBuffer(buffer), nil
}