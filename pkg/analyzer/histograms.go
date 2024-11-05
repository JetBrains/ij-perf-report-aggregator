package analyzer

import (
	"fmt"

	"github.com/valyala/fastjson"
)

type PMetric struct {
	Key   string
	Value float64
}

func getPercentile(p float64, boundaries []float64, counts []int) PMetric {
	pName := fmt.Sprintf("p%d", int(p*100))

	sum := 0
	for _, c := range counts {
		sum += c
	}

	targetCount := int(float64(sum) * p)

	b := len(boundaries)

	count := counts[0]
	if count >= targetCount {
		return PMetric{pName, boundaries[0]}
	}
	if sum-counts[b] < targetCount {
		return PMetric{pName, boundaries[b-1]}
	}

	for i, c := range counts[1:] {
		if (count + c) >= targetCount {
			delta := targetCount - count
			value := boundaries[i] + (boundaries[i+1]-boundaries[i])*float64(delta)/float64(c)
			return PMetric{pName, value}
		}
		count += c
	}

	return PMetric{}
}

/*
[
{"startEpochNanos":1725274564471190409, "epochNanos":1725274567798283106, "labels":{}, "sum":6422.472478000001, "count":256, "boundaries":[0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0, 500.0, 750.0, 1000.0, 2500.0, 5000.0, 7500.0, 10000.0], "counts":[1, 0, 0, 148, 101, 5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0]},
{"startEpochNanos":1725274567798283106, "epochNanos":1725274576388025039, "labels":{}, "sum":15996.745029000023, "count":747, "boundaries":[0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0, 500.0, 750.0, 1000.0, 2500.0, 5000.0, 7500.0, 10000.0], "counts":[0, 0, 0, 555, 187, 4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0]}
]
*/
func histogramToMetrics(json *fastjson.Value) []PMetric {
	values := json.GetArray("data")

	boundariesJson := values[0].GetArray("boundaries")
	boundaries := make([]float64, len(boundariesJson))
	for i, j := range boundariesJson {
		boundaries[i] = j.GetFloat64()
	}
	counts := make([]int, len(boundaries)+1)

	for _, value := range values {
		for bucket, count := range value.GetArray("counts") {
			counts[bucket] += count.GetInt()
		}
	}

	return []PMetric{
		getPercentile(0.5, boundaries, counts),
		getPercentile(0.95, boundaries, counts),
		getPercentile(0.99, boundaries, counts),
	}
}
