// Package comparison compares two series of runs, such as two branches or two modes, per project, metric and
// machine group. The dashboard compare pages and the compare_branches MCP tool share it.
package comparison

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

// SQL renders the query feeding Compare: the latest 50
// values of each measure per project/machine group for both values of dimCol (branch or mode),
// scanning the last month. The caller supplies only its own WHERE predicates.
// groupArraySorted keeps at most 50 values per group (newest first) instead of sorting the whole
// month: on perfUnitTests (~1k projects) the ORDER BY variant needed ~1.8 GiB, this one ~90 MiB.
func SQL(dimCol, table, where string) string {
	return fmt.Sprintf(
		"SELECT %[1]s AS Branch, project AS Project, measure_name AS MeasureName, machine_group AS Machine, "+
			"arrayMap(x -> x.2, groupArraySorted(50)((-toUnixTimestamp(generated_time), measure_value))) AS MeasureValues "+
			"FROM (SELECT %[1]s, project, %[2]s AS machine_group, measures.name AS measure_name, measures.value AS measure_value, generated_time "+
			"FROM %[3]s ARRAY JOIN measures "+
			"WHERE %[4]s AND generated_time > subtractMonths(now(), 1)) "+
			"GROUP BY %[1]s, project, measure_name, machine_group",
		dimCol, machine.GroupSQLExpr("machine"), table, where)
}

// shiftBoundsMisrate is the chance that the true shift lies outside DiffLower..DiffUpper. Kept low
// because a page compares hundreds of tests at once.
const shiftBoundsMisrate = 1e-3

// Row is marshalled as is by /api/compareBranches and /api/compareModes: its field names are their JSON.
type Row struct {
	Project     string
	MeasureName string
	Machine     string
	Median1     float64
	Median2     float64
	Diff        float64
	// bounds of the shift of value2 against value1, in percent of Median1
	DiffLower float64
	DiffUpper float64
	// the bounds exclude zero, i.e. Diff is not noise
	Significant bool
}

// Compare compares value1 against value2 per project/metric/machine
// group — runs from different hardware classes are never pooled into one median (AT-4930).
// A project/metric pair therefore yields one row per group.
func Compare(results []Series, value1, value2 string) []Row {
	type key struct {
		Project     string
		MeasureName string
		Machine     string
	}

	branch1Map := make(map[key][]int)
	branch2Map := make(map[key][]int)

	for _, item := range trimToLastChangePoint(results) {
		k := key{Project: item.Project, MeasureName: item.MeasureName, Machine: item.Machine}
		if item.Branch == value1 {
			branch1Map[k] = item.MeasureValues
		} else if item.Branch == value2 {
			branch2Map[k] = item.MeasureValues
		}
	}

	response := make([]Row, 0)
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

		// derived from the two centers shown next to it: pragmastat.Ratio (median of pairwise
		// ratios) can disagree in sign with them, e.g. 63 vs 62.5 reported as +1.6%
		var diff float64
		if center1 != 0 {
			diff = math.Round((center2/center1-1)*1000) / 10
		}

		item := Row{
			Project:     k.Project,
			MeasureName: k.MeasureName,
			Machine:     k.Machine,
			Median1:     center1,
			Median2:     center2,
			Diff:        diff,
		}
		// a single run on each side collapses the bounds to one point, which says nothing about noise
		if center1 != 0 && len(values1)*len(values2) > 1 {
			bounds, err := pragmastat.ShiftBounds(values2, values1, shiftBoundsMisrate)
			if err == nil {
				item.DiffLower = math.Round(bounds.Lower/center1*1000) / 10
				item.DiffUpper = math.Round(bounds.Upper/center1*1000) / 10
				item.Significant = bounds.Lower > 0 || bounds.Upper < 0
			}
		}
		response = append(response, item)
	}

	return response
}

// Series is one row of SQL: the newest values of one project, metric and machine group on one value of the
// compared column. Branch holds that value, which is a mode for a mode comparison.
type Series struct {
	Branch        string
	Project       string
	MeasureName   string
	Machine       string
	MeasureValues []int
}

// trimToLastChangePoint trims each series to the values after its last significant change point.
func trimToLastChangePoint(queryResults []Series) []Series {
	resultChan := make(chan Series, len(queryResults))
	var wg sync.WaitGroup

	for _, result := range queryResults {
		wg.Go(func() {
			values := make([]int, len(result.MeasureValues))
			copy(values, result.MeasureValues)
			slices.Reverse(values)
			indexes := statistic.GetChangePointIndexes(values, min(5, len(values)/2))
			validIndexes := degradation_detector.FilterValidChangePoints(values, indexes, 3.0, 3.0)
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

	items := make([]Series, 0, len(queryResults))
	for item := range resultChan {
		items = append(items, item)
	}
	return items
}
