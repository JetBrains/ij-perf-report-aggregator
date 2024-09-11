package analyzer

import (
	"encoding/json"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/valyala/fastjson"
)

var otParsers fastjson.ParserPool

type Trace struct {
	Data []struct {
		Spans []struct {
			OperationName string `json:"operationName"`
			Duration      int    `json:"duration"`
		} `json:"spans"`
	} `json:"data"`
}

func analyzeQodanaReport(runResult *RunResult, data model.ExtraData) bool {
	spansToReport := []string{
		"project.opening",
		"gradle.sync.duration",
		"application.exit",
		"qodana run",
		"qodanaProjectOpening",
		"qodanaScriptRun",
		"qodanaProjectConfiguration",
		"qodanaProjectAnalysis",
		"refGraphBuilding",
		"globalInspectionsAnalysis",
		"localInspectionsAnalysis",
	}

	runResult.Report = &model.Report{}

	parser := otParsers.Get()
	defer otParsers.Put(parser)
	props, err := parser.ParseBytes(data.TcBuildProperties)
	if err != nil {
		return true
	}
	buildName := props.GetStringBytes("env.TEAMCITY_BUILDCONF_NAME")
	runResult.Report.Project = string(buildName)

	runResult.Report.Generated = data.CurrentBuildTime.String()
	measureNames := make([]string, 0)
	measureValues := make([]int32, 0)
	measureTypes := make([]string, 0)

	spanDurations := analyzeOtJson(runResult.RawReport, spansToReport)
	for spanName, values := range spanDurations {
		for _, value := range values {
			measureNames = append(measureNames, spanName)
			measureValues = append(measureValues, value)
			measureTypes = append(measureTypes, "d")
		}
	}
	runResult.ExtraFieldData = []interface{}{measureNames, measureValues, measureTypes}
	return false
}

func analyzeOtJson(ot []byte, operationNames []string) map[string][]int32 {
	var trace Trace
	err := json.Unmarshal(ot, &trace)
	if err != nil {
		return nil
	}
	durationMap := make(map[string][]int32)
	for _, data := range trace.Data {
		for _, span := range data.Spans {
			for _, operationName := range operationNames {
				if span.OperationName == operationName {
					durationMap[operationName] = append(durationMap[operationName], int32(span.Duration/1000))
				}
			}
		}
	}
	return durationMap
}
