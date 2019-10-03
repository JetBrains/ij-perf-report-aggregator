package model

type Report struct {
  Version string `json:"version"`

  Generated   string `json:"generated"`
  Build       string `json:"build"`
  Os          string `json:"os"`
  ProductCode string `json:"productCode"`
  Runtime     string `json:"runtime"`

  // not used yet
  TraceEvents []TraceEvent `json:"traceEvents"`

  MainActivities           []Activity `json:"items"`
  PrepareAppInitActivities []Activity `json:"prepareAppInitActivities"`

  RawData       []byte    `json:"-"`
  GeneratedTime int64     `json:"-"`
  ExtraData     ExtraData `json:"-"`
}

type ExtraData struct {
  LastGeneratedTime int64

  ProductCode string
  BuildNumber string

  Machine   string
  TcBuildId int
}

type TraceEvent struct {
  Name  string `json:"name"`
  Phase string `json:"ph"`
  // in microseconds
  Timestamp int `json:"ts"`

  // in old reports (v10) can be int instead of string
  //Thread   string `json:"tid"`
  Category string `json:"cat"`
}

type Activity struct {
  Name   string `json:"name"`
  Thread string `json:"thread"`

  // in milliseconds
  Start    int `json:"start"`
  End      int `json:"end"`
  Duration int `json:"duration"`
}

// computed metrics
type DurationEventMetrics struct {
  Bootstrap               int `json:"bootstrap"`
  AppInitPreparation      int `json:"appInitPreparation"`
  AppInit                 int `json:"appInit"`
  PluginDescriptorLoading int `json:"pluginDescriptorLoading"`

  AppComponentCreation     int `json:"appComponentCreation"`
  ProjectComponentCreation int `json:"projectComponentCreation"`

  ModuleLoading int `json:"moduleLoading"`
}

type InstantEventMetrics struct {
  // value - not duration, but start, because it is instant event and not duration event
  Splash int `json:"splash"`
}
