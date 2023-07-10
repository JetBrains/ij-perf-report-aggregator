package tc_properties

import (
  "os"
  "testing"
)

func TestFilter(t *testing.T) {
  // /Volumes/data/Downloads/build.finish.properties
  t.Setenv("CLICKHOUSE_DATA_PATH", "test")

  data, err := os.ReadFile("../../testData/build.finish.properties")
  if err != nil {
    t.Error(err)
  }
  properties, err := ReadProperties(data)
  if err != nil {
    t.Error(err)
  }

  err = os.WriteFile("/tmp/foo.txt", properties, 0777)
  if err != nil {
    t.Error(err)
  }

  // err = backup()
  // if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
  //   t.Error(err)
  // }
}
