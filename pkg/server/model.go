package server

type MedianResult struct {
  metricName    string
  groupedValues []Value
}

type Value struct {
  group string
  value int
}