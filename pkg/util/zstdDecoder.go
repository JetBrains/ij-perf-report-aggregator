package util

import (
  _ "embed"
  "encoding/base64"
  "fmt"
  "github.com/klauspost/compress/zstd"
)

//go:embed zstd.dictionary
var ZstdDictionary []byte

func DecodeQuery(encoded string) ([]byte, error) {
  compressed, err := base64.RawURLEncoding.DecodeString(encoded)
  if err != nil {
    return nil, fmt.Errorf("cannot decode query: %w", err)
  }

  reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0), zstd.WithDecoderDicts(ZstdDictionary))
  if err != nil {
    return nil, fmt.Errorf("cannot create zstd reader: %w", err)
  }
  defer reader.Close()

  decompressed, err := reader.DecodeAll(compressed, nil)
  if err != nil {
    return nil, fmt.Errorf("cannot decompress query: %w", err)
  }
  return decompressed, nil
}

func EncodeQuery(data []byte) (string, error) {
  // Create a new ZSTD encoder with the dictionary
  writer, err := zstd.NewWriter(nil, zstd.WithEncoderDict(ZstdDictionary))
  if err != nil {
    return "", fmt.Errorf("cannot create zstd writer: %w", err)
  }

  // Compress the data
  compressed := writer.EncodeAll(data, nil)

  // Base64 URL encode the compressed data
  encoded := base64.RawURLEncoding.EncodeToString(compressed)
  return encoded, nil
}
