package main

import (
  "fmt"
  "github.com/alecthomas/kingpin"
  "go.uber.org/zap"
  "log"
  "os"
  "report-aggregator/pkg/server"
  "report-aggregator/pkg/util"
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
  command := app
  dbUrl := command.Flag("db", "The ClickHouse database URL.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    err := server.Serve(*dbUrl, true, log)
    if err != nil {
      return err
    }

    return nil
  })
}
