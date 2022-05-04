package main

import (
  "github.com/nats-io/nats.go"
  "log"
  "os"
  "os/exec"
  "time"
)

func main() {
  cmd := exec.Command("/usr/bin/clickhouse", "server", "--config-file=/etc/clickhouse-server/config.xml")
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    // wait for clickhouse server start
    time.Sleep(10 * time.Second)
    requestClearCache()
  }()

  err = cmd.Wait()
  if err != nil {
    log.Fatal(err)
  }
}

func requestClearCache() {
  url := os.Getenv("NATS")
  if len(url) == 0 {
    url = "nats://nats:4222"
  }

  nc, err := nats.Connect(url, nats.Name("NATS Sample Publisher"))
  if err != nil {
    log.Fatal(err)
  }
  defer nc.Close()

  err = nc.Publish("server.clearCache", []byte("clickhouse"))
  if err != nil {
    log.Fatal(err)
  }

  err = nc.Flush()
  if err != nil {
    log.Fatal(err)
  }
}
