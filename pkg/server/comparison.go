package server

import (
	"fmt"
	"math"
	"slices"
	"sync"

	"github.com/AndreyAkinshin/pragmastat/go/v4"
	degradation_detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/machine"
)

// groupedComparisonSQL renders the query feeding buildGroupedComparisonResponse: the latest 50
// values of each measure per project/machine group for both values of dimCol (branch or mode),
// scanning the last month. The caller supplies only its own WHERE predicates.
func groupedComparisonSQL(dimCol, table, where string) string {
	return fmt.Sprintf(
		"SELECT %[1]s AS Branch, project AS Project, measure_name AS MeasureName, machine_group AS Machine, "+
			"arraySlice(groupArray(measure_value), 1, 50) AS MeasureValues "+
			"FROM (SELECT %[1]s, project, %[2]s AS machine_group, measures.name AS measure_name, measures.value AS measure_value "+
			"FROM %[3]s ARRAY JOIN measures "+
			"WHERE %[4]s AND generated_time > subtractMonths(now(), 1) ORDER BY generated_time DESC) "+
			"GROUP BY %[1]s, project, measure_name, machine_group",
		dimCol, machine.GroupSQLExpr("machine"), table, where)
}

type branchComparisonResponseItem struct {
	Project     string
	MeasureName string
	Machine     string
	Median1     float64
	Median2     float64
	Diff        float64
}

// buildGroupedComparisonResponse compares value1 against value2 per project/metric/machine
// group — runs from different hardware classes are never pooled into one median (AT-4930).
// A project/metric pair therefore yields one row per group.
func buildGroupedComparisonResponse(results []ownerQueryResult, value1, value2 string) []branchComparisonResponseItem {
	type key struct {
		Project     string
		MeasureName string
		Machine     string
	}

	branch1Map := make(map[key][]int)
	branch2Map := make(map[key][]int)

	for _, item := range filterQueryResults(results) {
		k := key{Project: item.Project, MeasureName: item.MeasureName, Machine: item.Machine}
		if item.Branch == value1 {
			branch1Map[k] = item.MeasureValues
		} else if item.Branch == value2 {
			branch2Map[k] = item.MeasureValues
		}
	}

	response := make([]branchComparisonResponseItem, 0)
	for k, values1 := range branch1Map {
		values2, ok := branch2Map[k]
		if !ok {
			continue
		}

		center1, err := pragmastat.Center(values1)
		if err != nil {
			continue
		}
		center2, err := pragmastat.Center(values2)
		if err != nil {
			continue
		}

		var diff float64
		ratio, err := pragmastat.Ratio(values2, values1)
		if err == nil {
			diff = math.Round((ratio-1)*1000) / 10
		}

		response = append(response, branchComparisonResponseItem{
			Project:     k.Project,
			MeasureName: k.MeasureName,
			Machine:     k.Machine,
			Median1:     center1,
			Median2:     center2,
			Diff:        diff,
		})
	}

	return response
}

type measureQueryResult struct {
	Branch        string
	Project       string
	MeasureName   string
	MeasureValues []int
}

// filterQueryResults trims each series to the values after its last significant change point.
func filterQueryResults(queryResults []ownerQueryResult) []ownerQueryResult {
	resultChan := make(chan ownerQueryResult, len(queryResults))
	var wg sync.WaitGroup

	for _, result := range queryResults {
		wg.Go(func() {
			values := make([]int, len(result.MeasureValues))
			copy(values, result.MeasureValues)
			slices.Reverse(values)
			indexes := statistic.GetChangePointIndexes(values, min(5, len(values)/2))
			validIndexes := filterValidChangePoints(values, indexes, 3.0, 3.0)
			if len(validIndexes) == 0 {
				result.MeasureValues = values
			} else {
				result.MeasureValues = values[validIndexes[len(validIndexes)-1]:]
			}
			resultChan <- result
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	items := make([]ownerQueryResult, 0, len(queryResults))
	for item := range resultChan {
		items = append(items, item)
	}
	return items
}

// filterValidChangePoints filters change points based on median difference and effect size thresholds
func filterValidChangePoints(values []int, changePoints []int, medianDifferenceThreshold float64, effectSizeThreshold float64) []int {
	if len(changePoints) == 0 {
		return changePoints
	}

	segments := degradation_detector.GetSegmentsBetweenChangePoints(changePoints, values)
	if len(segments) < 2 {
		return []int{}
	}

	validChangePoints := make([]int, 0)

	for i := 1; i < len(segments); i++ {
		prevSegment := segments[i-1]
		currentSegment := segments[i]

		ratio, err := pragmastat.Ratio(currentSegment, prevSegment)
		if err != nil {
			continue
		}

		currentCenter, err := pragmastat.Center(currentSegment)
		if err != nil {
			continue
		}
		previousCenter, err := pragmastat.Center(prevSegment)
		if err != nil {
			continue
		}

		percentageChange := math.Abs((ratio - 1) * 100)
		absoluteChange := math.Abs(currentCenter - previousCenter)

		if absoluteChange < 10 || percentageChange < medianDifferenceThreshold {
			continue
		}

		effectSize := statistic.EffectSize(currentSegment, prevSegment)
		if effectSize < effectSizeThreshold {
			continue
		}

		validChangePoints = append(validChangePoints, changePoints[i-1])
	}

	return validChangePoints
}
