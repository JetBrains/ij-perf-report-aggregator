package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"
)

var networkRequestCount = runtime.NumCPU() + 1

func (t *Collector) loadReports(builds []*Build, reportExistenceChecker *ReportExistenceChecker, reportAnalyzer *analyzer.ReportAnalyzer) error {
	networkRequestCount = 20
	t.logger.Info("Network request count", "count", networkRequestCount)

	for index, build := range builds {
		if reportExistenceChecker.has(build.Id) {
			t.logger.Info("build already processed", "id", build.Id, "finishDate", build.FinishDate)
			builds[index] = nil
		}
	}

	if t.config.HasInstallerField {
		err := t.loadInstallerInfo(builds, networkRequestCount)
		if err != nil {
			return err
		}
	}
	if t.config.HasNoInstallerButHasChanges {
		err := t.loadChanges(builds, networkRequestCount)
		if err != nil {
			return err
		}
	}

	duration := time.Duration(len(builds)*300) * time.Second
	t.logger.Debug("load", "timeout", duration.Seconds())
	taskContextWithTimeout, cancel := context.WithTimeout(t.taskContext, duration)
	defer cancel()
	errGroup, loadContext := errgroup.WithContext(taskContextWithTimeout)

	errGroup.SetLimit(networkRequestCount)
	for _, build := range builds {
		if build == nil || build.Agent.Name == "Dead agent" {
			continue
		}
		errGroup.Go(func() error {
			t.logger.Info("processing build", "id", build.Id)
			if t.config.HasInstallerField && build.installerInfo == nil {
				// or already processed or cannot compute installer info
				return nil
			}

			artifacts, err := t.downloadReports(loadContext, *build)
			if err != nil {
				return err
			}

			if len(artifacts) == 0 {
				t.logger.Error("cannot find any performance report", "id", build.Id, "status", build.Status)
				return nil
			}

			tcBuildProperties, err := t.downloadBuildProperties(loadContext, *build)
			if err != nil {
				return err
			}
			if tcBuildProperties == nil {
				return nil
			}

			for _, artifact := range artifacts {
				if loadContext.Err() != nil {
					return nil
				}

				data := model.ExtraData{
					Machine:           build.Agent.Name,
					TcBuildId:         build.Id,
					TcBuildType:       build.Type,
					TcBuildProperties: tcBuildProperties,
					ReportFile:        artifact.path,
				}

				currentBuildTime, err := analyzer.ParseTime(build.FinishDate)
				if err == nil {
					data.CurrentBuildTime = currentBuildTime
				}

				if build.Private && build.TriggeredBy.User != nil {
					data.TriggeredBy = build.TriggeredBy.User.Email
				}

				if t.config.HasInstallerField {
					installerInfo := build.installerInfo
					data.BuildTime = installerInfo.buildTime
					data.Changes = installerInfo.changes
					data.TcInstallerBuildId = installerInfo.id
				}
				if t.config.HasNoInstallerButHasChanges {
					data.Changes = build.buildInfo.changes
				}

				if t.config.HasBuildNumber {
					data.TcBuildNumber = build.BuildNumber
				}

				err = reportAnalyzer.Analyze(artifact.data, data)
				if err != nil {
					if build.Status == "FAILURE" {
						t.logger.Warn("cannot parse performance report in the failed build", "buildId", build.Id, "error", err)
					} else {
						return err
					}
				}
			}
			return nil
		})
	}
	return errGroup.Wait()
}

func (t *Collector) loadChanges(builds []*Build, networkRequestCount int) error {
	var notLoadedBuildIds []*BuildInfo
	for _, build := range builds {
		if build == nil {
			continue
		}

		id := build.Id
		buildInfo := t.buildIdToInfo[id]
		if buildInfo == nil {
			buildInfo = &BuildInfo{
				id: id,
			}
			notLoadedBuildIds = append(notLoadedBuildIds, buildInfo)
			t.buildIdToInfo[id] = buildInfo
		}
		build.buildInfo = buildInfo
	}

	if len(notLoadedBuildIds) == 0 {
		return nil
	}
	notLoadedIds := make([]int, 0, len(notLoadedBuildIds))
	for _, installerInfo := range notLoadedBuildIds {
		notLoadedIds = append(notLoadedIds, installerInfo.id)
	}
	t.logger.Debug("load build info", "count", len(notLoadedBuildIds), "ids", notLoadedIds)

	errGroup, loadContext := errgroup.WithContext(t.taskContext)
	errGroup.SetLimit(networkRequestCount)
	for _, buildInfo := range notLoadedBuildIds {
		if buildInfo.id == -1 {
			continue
		}

		errGroup.Go(func() error {
			var err error
			buildInfo.changes, err = t.loadBuildChanges(loadContext, buildInfo.id)
			if err != nil {
				return fmt.Errorf("failed to load changes: %w", err)
			}
			return nil
		})
	}
	return errGroup.Wait()
}

func (t *Collector) loadInstallerInfo(builds []*Build, networkRequestCount int) error {
	var notLoadedInstallerBuildIds []*InstallerInfo
	for _, build := range builds {
		if build == nil {
			continue
		}

		id, buildTime, err := computeBuildDate(build)
		if err != nil {
			return err
		}

		if id == -1 {
			if t.config.HasInstallerField {
				t.logger.Error("cannot find installer build", "buildId", build.Id)
			}
			continue
		}

		installerInfo := t.installerBuildIdToInfo[id]
		if installerInfo == nil {
			installerInfo = &InstallerInfo{
				id:        id,
				buildTime: buildTime,
			}
			notLoadedInstallerBuildIds = append(notLoadedInstallerBuildIds, installerInfo)
			t.installerBuildIdToInfo[id] = installerInfo
		}
		build.installerInfo = installerInfo
	}

	if len(notLoadedInstallerBuildIds) == 0 {
		return nil
	}

	notLoadedIds := make([]int, 0, len(notLoadedInstallerBuildIds))
	for _, installerInfo := range notLoadedInstallerBuildIds {
		notLoadedIds = append(notLoadedIds, installerInfo.id)
	}
	t.logger.Debug("load installer info", "count", len(notLoadedInstallerBuildIds), "ids", notLoadedIds)

	errGroup, loadContext := errgroup.WithContext(t.taskContext)
	errGroup.SetLimit(networkRequestCount)
	for _, installerInfo := range notLoadedInstallerBuildIds {
		if installerInfo.id == -1 {
			continue
		}

		errGroup.Go(func() error {
			var err error
			installerInfo.changes, err = t.loadBuildChanges(loadContext, installerInfo.id)
			if err != nil {
				return fmt.Errorf("failed to load changes: %w", err)
			}
			return nil
		})
	}
	return errGroup.Wait()
}

func computeBuildDate(build *Build) (int, time.Time, error) {
	for _, dependencyBuild := range build.ArtifactDependencies.Builds {
		if strings.Contains(dependencyBuild.BuildTypeId, "Installer") || strings.Contains(dependencyBuild.BuildTypeId, "Distribution") {
			result, err := time.Parse(tcTimeFormat, dependencyBuild.FinishDate)
			if err != nil {
				return -1, time.Time{}, err
			}

			return dependencyBuild.Id, result, nil
		}
	}
	return -1, time.Time{}, nil
}

func (t *Collector) loadBuildChanges(ctx context.Context, buildId int) ([]string, error) {
	artifactUrl, err := url.Parse(t.serverUrl + "/changes?locator=build:(id:" + strconv.Itoa(buildId) + ")&fields=change(version)&count=10000")
	if err != nil {
		return nil, err
	}

	response, err := t.get(ctx, artifactUrl.String())
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode > 300 {
		responseBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid response (%s): %s", response.Status, responseBody)
	}

	t.storeSessionIdCookie(response)

	var changeList ChangeList
	err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&changeList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	encoding := base64.RawStdEncoding

	var b bytes.Buffer
	result := make([]string, len(changeList.List))
	for index, change := range changeList.List {
		if strings.Contains(change.Version, " ") {
			// private build with custom change, format: 13 04 2022 12:14
			continue
		}
		data, err := hex.DecodeString(change.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to decode change version: %w", err)
		}

		b.Reset()

		buf := make([]byte, encoding.EncodedLen(len(data)))
		encoding.Encode(buf, data)
		result[index] = string(buf)
	}

	return result, nil
}
