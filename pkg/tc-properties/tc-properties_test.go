package tc_properties

import (
  "github.com/stretchr/testify/require"
  "os"
  "testing"
)

func TestFilter(t *testing.T) {
  // /Volumes/data/Downloads/build.finish.properties
  t.Setenv("CLICKHOUSE_DATA_PATH", "test")

  data, err := os.ReadFile("../../testData/build.finish.properties")
  require.NoError(t, err)
  properties, err := ReadProperties(data)
  require.NoError(t, err)

  err = os.WriteFile("/tmp/foo.txt", properties, 0777)
  require.NoError(t, err)

  // err = backup()
  // if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
  //   t.Error(err)
  // }
}
