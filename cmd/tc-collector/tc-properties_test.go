package main

import (
  "io/ioutil"
  "os"
  "testing"
)

func TestFilter(t *testing.T) {
  // /Volumes/data/Downloads/build.finish.properties
  err := os.Setenv("CLICKHOUSE_DATA_PATH", "test")
  if err != nil {
    t.Error(err)
  }

  data, err := ioutil.ReadFile("/Volumes/data/Downloads/build.finish.properties copy")
  if err != nil {
    t.Error(err)
  }
  properties, err := readProperties(data)
  if err != nil {
    t.Error(err)
  }

  err = ioutil.WriteFile("/Volumes/data/Downloads/build.finish.properties", properties, 0777)
  if err != nil {
    t.Error(err)
  }

  //err = backup()
  //if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
  //  t.Error(err)
  //}
}
