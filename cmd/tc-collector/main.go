package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/alecthomas/kingpin"
  "github.com/araddon/dateparse"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "log"
  "net/http"
  "os"
  "time"
)

func main() {
	logger := util.CreateLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("perf-db-server", "perf-db-server").Version("0.0.1")
  configureCollectFromTeamCity(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func configureCollectFromTeamCity(app *kingpin.Application, logger *zap.Logger) {
  buildTypeIds := app.Flag("build-type-id", "The TeamCity build type id.").Short('c').Required().Strings()
  clickHouseUrl := app.Flag("clickhouse", "The ClickHouse server URL.").Default("localhost:9000").String()
  tcUrl := app.Flag("tc", "The TeamCity server URL.").Required().String()
  sinceDate := app.Flag("since", "The date to force collecting since").String()

  app.Action(func(context *kingpin.ParseContext) error {
    var since time.Time
    if len(*sinceDate) > 0 {
      var err error
      since, err = dateparse.ParseStrict(*sinceDate)
      if err != nil {
        return errors.WithStack(err)
      }
    }

    var httpClient = &http.Client{
      Timeout: 30 * time.Second,
    }

    err := collectFromTeamCity(*clickHouseUrl, *tcUrl, *buildTypeIds, since, httpClient, logger)
    if err != nil {
      return err
    }

    natsUrl := os.Getenv("NATS")
    if len(natsUrl) > 0 {
      err = doNotifyServer(natsUrl, logger)
      if err != nil {
        return err
      }
    }

    return nil
  })
}
