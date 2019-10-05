package main

import (
  "fmt"
  "github.com/alecthomas/kingpin"
  "log"
  "os"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/filling"
  "report-aggregator/pkg/ideaLog"
  "report-aggregator/pkg/server"
  "report-aggregator/pkg/teamcity"
  "report-aggregator/pkg/util"
)

func main() {
	logger := util.CreateLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("report-aggregator", "report-aggregator").Version("0.0.1")

	ideaLog.ConfigureCollectFromDirCommand(app, logger)
  teamcity.ConfigureCollectFromTeamCity(app, logger)

	server.ConfigureServeCommand(app, logger)
	filling.ConfigureFillCommand(app, logger)
	analyzer.ConfigureUpdateMetricsCommand(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}
