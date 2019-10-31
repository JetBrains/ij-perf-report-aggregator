package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/alecthomas/kingpin"
  "go.uber.org/zap"
  "log"
  "os"
)

func main() {
	logger := util.CreateLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("perf-db-server", "perf-db-server").Version("0.0.1")

	configureServeCommand(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}

func configureServeCommand(app *kingpin.Application, log *zap.Logger) {
  dbUrl := app.Flag("db", "The ClickHouse database URL.").Required().String()
  natsUrl := app.Flag("nats", "The NATS URL.").String()
  app.Action(func(context *kingpin.ParseContext) error {
    err := server.Serve(*dbUrl, *natsUrl, log)
    if err != nil {
      return err
    }

    return nil
  })
}
