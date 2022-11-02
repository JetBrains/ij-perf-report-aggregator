package main

import (
  _ "embed"
  "github.com/nats-io/nats.go"
  "log"
  "os"
  "os/exec"
  "strings"
  "time"
)

//go:embed config.xml
var clickhouseConfig []byte

func main() {
  s3Url := "https://" + os.Getenv("S3_BUCKET") + ".s3.eu-west-1.amazonaws.com/data/"

  s := strings.NewReplacer(
    "$S3_URL", s3Url,
    "$S3_ACCESS_KEY", os.Getenv("S3_ACCESS_KEY"),
    "$S3_SECRET_KEY", os.Getenv("S3_SECRET_KEY"),
  ).Replace(string(clickhouseConfig))

  log.Print("S3 URL: " + s3Url)

  // /etc is not writeable
  configFile := "/var/lib/clickhouse/config.xml"
  err := os.WriteFile(configFile, []byte(s), 0666)
  if err != nil {
    log.Fatal(err)
  }

  cmd := exec.Command("/usr/bin/clickhouse", "server", "--config-file="+configFile)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err = cmd.Start()
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
