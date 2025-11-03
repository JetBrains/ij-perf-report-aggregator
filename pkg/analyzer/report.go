package analyzer

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/valyala/fastjson"
)

var parserPool fastjson.ParserPool

func ReadReport(runResult *RunResult, config DatabaseConfiguration) error {
	parser := parserPool.Get()
	defer parserPool.Put(parser)

	report, err := parser.ParseBytes(runResult.RawReport)
	if err != nil {
		endIndex := min(len(runResult.RawReport), 10000)
		slog.Warn("invalid report. corrupted JSON, report will be skipped", "error", err, "file", runResult.ReportFileName, "rawReport", runResult.RawReport[:endIndex])
		runResult.Report = nil
		return nil
	}

	runResult.Report = &model.Report{
		Version:            string(report.GetStringBytes("version")),
		Generated:          string(report.GetStringBytes("generated")),
		Project:            string(report.GetStringBytes("project")),
		ProjectURL:         string(report.GetStringBytes("projectURL")),
		ProjectDescription: string(report.GetStringBytes("projectDescription")),

		ProductCode: string(report.GetStringBytes("productCode")),
		MethodName:  string(report.GetStringBytes("methodName")),
	}

	if config.HasInstallerField {
		runResult.Report.Build = string(report.GetStringBytes("build"))
		runResult.Report.BuildDate = string(report.GetStringBytes("buildDate"))
	}

	err = config.ReportReader(runResult, report)
	if err != nil {
		return nil
	}

	runResult.RawReport = report.MarshalTo(nil)
	return nil
}

func getBuildTimeFromReport(report *model.Report) (time.Time, error) {
	var buildTimeUnix time.Time
	buildTime, err := ParseTime(report.BuildDate)
	if err != nil {
		return time.Time{}, err
	}
	buildTimeUnix = buildTime
	return buildTimeUnix, nil
}

func ParseTime(s string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC1123, s)
	}

	if err != nil {
		parsedTime, err = time.Parse("Jan 2, 2006, 3:04:05 PM MST", s)
	}

	if err != nil {
		parsedTime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", s)
	}

	if err != nil {
		parsedTime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", s)
	}

	if err != nil {
		parsedTime, err = time.Parse("20060102T150405+0000", s)
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse time %s: %w", s, err)
	}
	return parsedTime, nil
}
