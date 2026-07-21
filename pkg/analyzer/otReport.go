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
		"OpenGrepGlobalInspection",
		"Building IR",
		"Collecting Opengrep AST/Problems",
		"Building IR library symbols",
		"Building IR project",
		"Running RML + IFDS",
	}

	aggregatedSpansToReport := []string{
		"Matching taint entries",
		"Running DFA engine",
		"Running reverse debug",
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

	aggregatedSpanDurations := aggregateOtJson(runResult.RawReport, aggregatedSpansToReport)
	for spanName, aggregated := range aggregatedSpanDurations {
		measureNames = append(measureNames, spanName)
		measureValues = append(measureValues, aggregated.sum)
		measureTypes = append(measureTypes, "d")

		measureNames = append(measureNames, spanName+".count")
		measureValues = append(measureValues, aggregated.count)
		measureTypes = append(measureTypes, "c")
	}

	runResult.ExtraFieldData = []any{measureNames, measureValues, measureTypes /*mode*/, ""}
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

type aggregatedSpan struct {
	sum   int32
	count int32
}

// aggregateOtJson sums the durations and counts occurrences of spans that are expected to repeat
// multiple times with the same operation name within a single report (e.g. one span per chunk in a
// loop). A name with zero occurrences is simply absent from the returned map.
func aggregateOtJson(ot []byte, operationNames []string) map[string]aggregatedSpan {
	var trace Trace
	err := json.Unmarshal(ot, &trace)
	if err != nil {
		return nil
	}
	aggregatedMap := make(map[string]aggregatedSpan)
	for _, data := range trace.Data {
		for _, span := range data.Spans {
			for _, operationName := range operationNames {
				if span.OperationName == operationName {
					aggregated := aggregatedMap[operationName]
					aggregated.sum += int32(span.Duration / 1000)
					aggregated.count++
					aggregatedMap[operationName] = aggregated
				}
			}
		}
	}
	return aggregatedMap
}
