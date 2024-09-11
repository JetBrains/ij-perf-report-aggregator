package analyzer

import (
	"github.com/valyala/fastjson"
	"reflect"
	"testing"
)

func Test_getPercentile(t *testing.T) {
	type args struct {
		p          float64
		boundaries []float64
		counts     []int
	}
	tests := []struct {
		name string
		args args
		want PMetric
	}{
		{
			name: "50th percentile",
			args: args{
				p:          0.5,
				boundaries: []float64{0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0},
				counts:     []int{1, 0, 0, 148, 101, 5, 1, 0, 0},
			},
			want: PMetric{Key: "p50", Value: 22.87162162162162},
		},
		{
			name: "95th percentile",
			args: args{
				p:          0.95,
				boundaries: []float64{0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0},
				counts:     []int{0, 0, 0, 555, 187, 4, 1, 0, 0},
			},
			want: PMetric{Key: "p95", Value: 45.58823529411765},
		},
		{
			name: "99th percentile",
			args: args{
				p:          0.99,
				boundaries: []float64{0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0},
				counts:     []int{0, 0, 0, 555, 187, 4, 1, 0, 0},
			},
			want: PMetric{Key: "p99", Value: 49.59893048128342},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPercentile(tt.args.p, tt.args.boundaries, tt.args.counts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPercentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_histogramToMetrics(t *testing.T) {
	rawJSON := `{
		"data": [
			{
				"boundaries": [0.0, 5.0, 10.0, 25.0, 50.0, 75.0, 100.0, 250.0],
				"counts": [1, 0, 0, 148, 101, 5, 1, 0, 0]
			}
		]
	}`

	parsedJSON, _ := fastjson.Parse(rawJSON)

	tests := []struct {
		name string
		args *fastjson.Value
		want []PMetric
	}{
		{
			name: "single histogram",
			args: parsedJSON,
			want: []PMetric{
				{Key: "p50", Value: 22.87162162162162},
				{Key: "p95", Value: 48.26732673267327},
				{Key: "p99", Value: 65},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compare the two slices in an order-agnostic way.
			if got := histogramToMetrics(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("histogramToMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
