package main

import (
  "context"
  "encoding/json"
  "errors"
  "flag"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/araddon/dateparse"
  e "github.com/develar/errors"
  "go.uber.org/zap"
  "net/http"
  "os"
  "strings"
  "time"
)

// 1. You need to provide CONFIG env variable that may look like:
// {"teamcityUrl": "http://buildserver.labs.intellij.net", "buildConfigurations":[{"db": "perfint", "configurations": ["ijplatform_IjPlatform221_IntegrationPerformanceTestsLinux"]}]}
// 2. You also need to provide TC_TOKEN env variable which can be generated at: https://buildserver.labs.intellij.net/profile.html?item=accessTokens#
// 3. Clickhouse DB should be up and running (see readme.md "Adding a New Database" section)
func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  err := configureCollectFromTeamCity(logger)
  if err != nil {
    if errors.Is(err, context.Canceled) {
      os.Exit(78)
    }
    logger.Fatal("cannot collect", zap.Error(err))
  }
}

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func configureCollectFromTeamCity(logger *zap.Logger) error {
  clickHouseUrl := util.GetEnv("CLICKHOUSE", "localhost:9000")
  sinceDate := flag.String("since", "", "The date to force collecting since")
  flag.Parse()

  var since time.Time
  if len(*sinceDate) > 0 {
    var err error
    since, err = dateparse.ParseStrict(*sinceDate)
    if err != nil {
      return e.WithStack(err)
    }
  }

  var config CollectorConfiguration
  rawJson, err := util.GetEnvOrFile("CONFIG", "/etc/config/config.json")
  if err != nil {
    return err
  }

  rawJson = strings.TrimSpace(rawJson)
  if len(rawJson) == 0 {
    return e.New("File /etc/config/config.json is empty or env CONFIG is not set")
  }

  err = json.Unmarshal([]byte(rawJson), &config)
  if err != nil {
    return e.WithMessage(err, "cannot parse json: "+rawJson)
  }

  var httpClient = &http.Client{}
  httpClient.CheckRedirect = checkRedirectFunc

  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  for _, chunk := range config.BuildConfigurations {
    if taskContext.Err() != nil {
      break
    }

    var initialSince time.Time

    var buildConfigurationIds []string
    switch {
    case len(chunk.Configurations) == 0:
      osList := []string{"Mac", "Linux", "Windows"}
      for _, product := range chunk.Products {
        for _, osName := range osList {
          buildConfigurationIds = append(buildConfigurationIds, "ijplatform_master_"+productCodeToBuildName[strings.ToUpper(product)]+"StartupPerfTest"+osName)
        }
        buildConfigurationIds = append(buildConfigurationIds, "ijplatform_master_"+productCodeToBuildName[strings.ToUpper(product)]+"StartupPerfTestMacM1")
      }
    case chunk.Database == "jbr":
      jbrTypes := []string{"macOS12aarch64Metal", "macOS12aarch64OGL", "macOS12x64Metal", "macOS12x64OGL", "macOS13aarch64Metal", "macOS13aarch64OGL", "macOS13x64Metal",
        "macOS13x64OGL", "Ubuntu2004x64", "Ubuntu2004x64OGL", "Ubuntu2204x64", "Ubuntu2204x64OGL", "Windows10x64", "Windows11x64"}
      for _, configuration := range chunk.Configurations {
        for _, jbrType := range jbrTypes {
          buildConfigurationIds = append(buildConfigurationIds, configuration+"_"+jbrType)
        }
      }
    default:
      if len(chunk.Products) != 0 {
        return e.New("Must be specified either configurations or products, but not both")
      }
      for _, configuration := range chunk.Configurations {
        collector := &Collector{
          serverUrl:  config.TeamcityUrl + "/app/rest",
          httpClient: httpClient,
          logger:     logger,
        }
        configurations, err := collector.getSnapshots(taskContext, configuration)
        logger.Info("get snapshots", zap.Strings("configurations", configurations))
        if err != nil {
          logger.Warn("cannot get snapshots", zap.Error(err))
        }
        buildConfigurationIds = append(buildConfigurationIds, configurations...)
      }
    }

    if len(chunk.InitialSince) != 0 {
      initialSince, err = dateparse.ParseStrict(chunk.InitialSince)
      if err != nil {
        return e.WithStack(err)
      }
    }

    err = collectFromTeamCity(taskContext, clickHouseUrl, config.TeamcityUrl, chunk.Database, buildConfigurationIds, initialSince, since, httpClient, logger)
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
}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
  req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
  return nil
}

type CollectorConfiguration struct {
  TeamcityUrl         string           `json:"teamcityUrl"`
  BuildConfigurations []CollectorChunk `json:"buildConfigurations"`
}

type CollectorChunk struct {
  Database       string   `json:"db"`
  InitialSince   string   `json:"initialSince"`
  Products       []string `json:"products"`
  Configurations []string `json:"configurations"`
}
