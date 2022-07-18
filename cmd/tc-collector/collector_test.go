package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/stretchr/testify/assert"
  "go.uber.org/zap"
  "testing"
)

func getLogger() *zap.Logger {
  config := zap.NewDevelopmentConfig()
  config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
  config.DisableCaller = true
  config.DisableStacktrace = true
  logger, _ := config.Build()
  return logger
}

func TestMapping1(t *testing.T) {
  newName := analyzer.MapPerfMeasureName("delayType#max_awt_delay", []string{"delayType", "delayType_1", "delayType#max_awt_delay", "delayType_1#max_awt_delay"}, getLogger())
  assert.Equal(t, "typing_1#max_awt_delay", newName)
}

func TestMapping2(t *testing.T) {
  newName := analyzer.MapPerfMeasureName("delayType#max_awt_delay", []string{"delayType", "delayType#max_awt_delay"}, getLogger())
  assert.Equal(t, "typing#max_awt_delay", newName)
}

func TestMapping3(t *testing.T) {
  newName := analyzer.MapPerfMeasureName("delayType_1#max_awt_delay", []string{"delayType", "delayType_1", "delayType#max_awt_delay", "delayType_1#max_awt_delay"}, getLogger())
  assert.Equal(t, "typing_2#max_awt_delay", newName)
}
