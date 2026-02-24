package server

import (
	"math"
	"slices"
	"sync"

	"github.com/AndreyAkinshin/pragmastat/go/v4"
	degradation_detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/statistic"
)

type branchComparisonResponseItem struct {
	Project     string
	MeasureName string
	Median1     float64
	Median2     float64
	Diff        float64
}

type filteredValues struct {
	Branch      string
	Project     string
	MeasureName string
	Values      []int
}

func filterQueryResults(queryResults []struct {
	Branch        string
	Project       string
	MeasureName   string
	MeasureValues []int
},
) []filteredValues {
	resultChan := make(chan filteredValues, len(queryResults))
	var wg sync.WaitGroup

	for _, result := range queryResults {
		wg.Go(func() {
			values := make([]int, len(result.MeasureValues))
			copy(values, result.MeasureValues)
			slices.Reverse(values)
			indexes := statistic.GetChangePointIndexes(values, min(5, len(values)/2))
			validIndexes := filterValidChangePoints(values, indexes, 3.0, 3.0)
			var filtered []int
			if len(validIndexes) == 0 {
				filtered = values
			} else {
				lastIndex := validIndexes[len(validIndexes)-1]
				filtered = values[lastIndex:]
			}
			resultChan <- filteredValues{
				Branch:      result.Branch,
				Project:     result.Project,
				MeasureName: result.MeasureName,
				Values:      filtered,
			}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	items := make([]filteredValues, 0, len(queryResults))
	for item := range resultChan {
		items = append(items, item)
	}
	return items
}

func buildBranchComparisonResponse(items []filteredValues, branch1, branch2 string) []branchComparisonResponseItem {
	type key struct {
		Project     string
		MeasureName string
	}

	branch1Map := make(map[key][]int)
	branch2Map := make(map[key][]int)

	for _, item := range items {
		k := key{Project: item.Project, MeasureName: item.MeasureName}
		if item.Branch == branch1 {
			branch1Map[k] = item.Values
		} else if item.Branch == branch2 {
			branch2Map[k] = item.Values
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
			Median1:     center1,
			Median2:     center2,
			Diff:        diff,
		})
	}

	return response
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
