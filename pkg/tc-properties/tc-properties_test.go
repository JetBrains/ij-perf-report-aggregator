package tc_properties

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	// /Volumes/data/Downloads/build.finish.properties
	t.Setenv("CLICKHOUSE_DATA_PATH", "test")

	data, err := os.ReadFile("../../testData/build.finish.properties")
	require.NoError(t, err)
	properties, err := ReadProperties(data)
	require.NoError(t, err)

	err = os.WriteFile("/tmp/foo.txt", properties, 0o777)
	require.NoError(t, err)

	// err = backup()
	// if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
	//   t.Error(err)
	// }
}

func TestReadProperties(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../testData/build.finish.properties")
	require.NoError(t, err)
	properties, _ := LoadBytes(data, nil)
	v, ok := properties.Get("vcsroot.authMethod")
	assert.True(t, ok)
	assert.Equal(t, "TEAMCITY_SSH_KEY", v)
}
