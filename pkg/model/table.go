package model

var EssentialDurationMetricNames = []string{"bootstrap", "appInitPreparation", "appInit", "pluginDescriptorLoading", "appComponentCreation", "projectComponentCreation"}
var DurationMetricNames = append(EssentialDurationMetricNames, "moduleLoading", "projectDumbAware", "editorRestoring")
var InstantMetricNames = []string{"splash", "startUpCompleted"}

// https://clickhouse.yandex/docs/en/query_language/alter/#manipulations-with-key-expressions
// To keep the property that data part rows are ordered by the sorting key expression
// you cannot add expressions containing existing columns to the sorting key (only columns added by the ADD COLUMN command in the same ALTER query).

type IdAndName struct {
  Id   int
  Name string
}

func ProcessMetricName(handler func(name string, isInstant bool)) {
  for _, name := range DurationMetricNames {
    handler(name, false)
  }
  for _, name := range InstantMetricNames {
    handler(name, true)
  }
}