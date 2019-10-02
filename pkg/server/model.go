package server

type MedianResult struct {
  metricName   string
  buildToValue []Value
}

type Value struct {
  buildC1 int
  value   int
}
