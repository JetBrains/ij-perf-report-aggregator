package main

import (
  "fmt"
  "github.com/alecthomas/kingpin"
  "go.uber.org/zap"
  "log"
  "os"
  "report-aggregator/pkg/filling"
  "report-aggregator/pkg/ideaLog"
  "report-aggregator/pkg/server"
  "report-aggregator/pkg/util"
)

func main() {
	logger := util.CreateLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("report-aggregator", "report-aggregator").Version("0.0.1")

	ideaLog.ConfigureCollectFromDirCommand(app, logger)

	ConfigureServeCommand(app, logger)
	filling.ConfigureFillCommand(app, logger)
	filling.ConfigureTransformCommand(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}


func ConfigureServeCommand(app *kingpin.Application, log *zap.Logger) {
  command := app.Command("serve", "Start aggregated stats server.")
  dbUrl := command.Flag("db", "The ClickHouse database URL.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    err := server.Serve(*dbUrl, log)
    if err != nil {
      return err
    }

    return nil
  })
}
