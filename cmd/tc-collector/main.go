package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/araddon/dateparse"
	"go.deanishe.net/env"
	_ "go.uber.org/automaxprocs"
	"log/slog"
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
	err := configureCollectFromTeamCity()
	if err != nil {
		slog.Error("cannot collect", "err", err)
		os.Exit(78)
	}
}

func hasOSSuffix(osList []string, configuration string) bool {
	for _, osName := range osList {
		if strings.HasSuffix(configuration, osName) {
			return true
		}
	}
	return false
}

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func configureCollectFromTeamCity() error {
	clickHouseUrl := env.Get("CLICKHOUSE", "127.0.0.1:9000")
	sinceDate := flag.String("since", "", "The date to force collecting since")
	flag.Parse()

	var since time.Time
	if *sinceDate != "" {
		var err error
		since, err = dateparse.ParseStrict(*sinceDate)
		if err != nil {
			return fmt.Errorf("cannot parse since date: %w", err)
		}
	}

	var config CollectorConfiguration
	rawJson, err := util.GetEnvOrFile("CONFIG", "/etc/config/config.json")
	if err != nil {
		return err
	}

	rawJson = strings.TrimSpace(rawJson)
	if rawJson == "" {
		return errors.New("file /etc/config/config.json is empty or env CONFIG is not set")
	}

	err = json.Unmarshal([]byte(rawJson), &config)
	if err != nil {
		return errors.New("cannot parse json: " + rawJson)
	}

	var httpClient = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}
	httpClient.CheckRedirect = checkRedirectFunc

	taskContext, cancel := util.CreateCommandContext()
	defer cancel()

	for _, chunk := range config.BuildConfigurations {
		if taskContext.Err() != nil {
			break
		}

		var buildConfigurationIds []string
		switch {
		case chunk.Database == "ij":
			osList := []string{"Linux", "Windows", "MacM2"}
			for _, configuration := range chunk.Configurations {
				if hasOSSuffix(osList, configuration) {
					buildConfigurationIds = append(buildConfigurationIds, configuration)
				} else {
					for _, osName := range osList {
						buildConfigurationIds = append(buildConfigurationIds, configuration+osName)
					}
				}
			}
		case chunk.Database == "jbr":
			jbrTypes := []string{"MacOS12aarch64Metal", "MacOS12aarch64OGL", "MacOS12x86_64Metal", "MacOS12x86_64OGL", "MacOS13aarch64Metal", "MacOS13aarch64OGL", "MacOS13x86_64Metal",
				"MacOS13x86_64OGL", "Ubuntu2004x86_64", "Ubuntu2004x86_64OGL", "Ubuntu2204x86_64", "Ubuntu2204x86_64OGL", "Windows10x86_64", "Windows11x86_64"}
			for _, configuration := range chunk.Configurations {
				for _, jbrType := range jbrTypes {
					buildConfigurationIds = append(buildConfigurationIds, configuration+"_"+jbrType)
				}
			}
		default:
			for _, configuration := range chunk.Configurations {
				collector := &Collector{
					serverUrl:  config.TeamcityUrl + "/app/rest",
					httpClient: httpClient,
				}
				configurations, err := collector.getSnapshots(taskContext, configuration)
				slog.Info("get snapshots", "configurations", configurations)
				if err != nil {
					slog.Warn("cannot get snapshots", "err", err)
				}
				buildConfigurationIds = append(buildConfigurationIds, configurations...)
			}
		}

		err = collectFromTeamCity(taskContext, clickHouseUrl, config.TeamcityUrl, chunk.Database, buildConfigurationIds, since, httpClient)
		if err != nil {
			return err
		}
	}

	natsUrl := os.Getenv("NATS")
	if natsUrl != "" {
		err = doNotifyServer(natsUrl)
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
	Configurations []string `json:"configurations"`
}
