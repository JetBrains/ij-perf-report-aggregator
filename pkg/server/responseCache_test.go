package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressData(t *testing.T) {
	rcm, err := NewResponseCacheManager()
	require.NoError(t, err)

	testData := []byte("sample data to compress")

	compressedData, err := rcm.compressData(testData)
	require.NoError(t, err)
	assert.NotEmptyf(t, compressedData, "Expected compressed data, got empty")
}
