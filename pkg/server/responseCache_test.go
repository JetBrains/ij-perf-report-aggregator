package server

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestCompressData(t *testing.T) {
  rcm, err := NewResponseCacheManager()
  assert.NoError(t, err)

  testData := []byte("sample data to compress")

  compressedData, err := rcm.compressData(testData)
  assert.NoError(t, err)
  assert.NotEmptyf(t, compressedData, "Expected compressed data, got empty")

}
