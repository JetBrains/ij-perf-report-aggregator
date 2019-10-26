package main

import (
  "os"
  "strings"
  "testing"
)

// test that config grabs env
func TestBackup(t *testing.T) {
  var err error

  err = os.Setenv("CLICKHOUSE_DATA_PATH", "test")
  if err != nil {
    t.Error(err)
  }

  err = backup()
  if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
    t.Error(err)
  }
}
