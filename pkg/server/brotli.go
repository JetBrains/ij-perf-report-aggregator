package server

import (
  "bytes"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/andybalholm/brotli"
  "github.com/valyala/bytebufferpool"
)

func decompressData(input []byte) ([]byte, error) {
  buf := bytebufferpool.Get()
  defer bytebufferpool.Put(buf)
  reader := brotli.NewReader(bytes.NewReader(input))
  _, err := buf.ReadFrom(reader)
  if err != nil {
    return nil, err
  }
  return CopyBuffer(buf), nil
}

func (rcm *ResponseCacheManager) compressData(value []byte) ([]byte, error) {
  buffer := bytebufferpool.Get()
  defer bytebufferpool.Put(buffer)
  writer := brotli.NewWriter(buffer)
  _, err := writer.Write(value)
  if err != nil {
    util.Close(writer, rcm.logger)
    return nil, err
  }

  util.Close(writer, rcm.logger)
  return CopyBuffer(buffer), nil
}
