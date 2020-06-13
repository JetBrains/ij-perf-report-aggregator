package main

import (
  "encoding/json"
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/alecthomas/kingpin"
  "github.com/araddon/dateparse"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "log"
  "net/http"
  "os"
  "strings"
  "time"
)

func main() {
	logger := util.CreateLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("tc-collector", "tc-collector").Version("0.0.1")
  configureCollectFromTeamCity(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func configureCollectFromTeamCity(app *kingpin.Application, logger *zap.Logger) {
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

    var chunks []CollectorChunk
    rawJson := strings.TrimSpace(os.Getenv("CONFIG"))
    if len(rawJson) == 0 {
      return errors.New("Env CONFIG is not set")
    }

    err := json.Unmarshal([]byte(rawJson), &chunks)
    if err != nil {
      return errors.WithMessage(err, "cannot parse json: "+rawJson)
    }

    var httpClient = &http.Client{
    }

    taskContext, cancel := util.CreateCommandContext()
    defer cancel()

    for _, chunk := range chunks {
      if taskContext.Err() != nil {
        break
      }

      var buildConfigurationIds []string
      if len(chunk.Configurations) == 0 {
        osList := []string{"Mac", "Linux", "Windows"}
        for _, product := range chunk.Products {
          for _, osName := range osList {
            buildConfigurationIds = append(buildConfigurationIds, "ijplatform_master_"+productCodeToBuildName[strings.ToUpper(product)]+"StartupPerfTest"+osName)
          }
        }
      } else {
        if len(chunk.Products) != 0 {
          return errors.New("Must be specified or configurations, or products, but not both")
        }

        buildConfigurationIds = chunk.Configurations
      }

      err := collectFromTeamCity(*clickHouseUrl, *tcUrl, chunk.Database, buildConfigurationIds, since, httpClient, logger, taskContext, cancel)
      if err != nil {
        return err
      }
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

type CollectorChunk struct {
  Database       string   `json:"db"`
  Products       []string `json:"products"`
  Configurations []string `json:"configurations"`
}