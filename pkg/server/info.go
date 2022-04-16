package server

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "net/http"
)

func (t *StatsServer) handleMetaMeasureRequest(_ *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
  buffer := byteBufferPool.Get()
  isOk := false
  defer func() {
    if !isOk {
      byteBufferPool.Put(buffer)
    }
  }()

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

  isOk = true
  return buffer, true, nil
}
