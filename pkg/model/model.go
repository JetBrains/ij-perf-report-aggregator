package model

import "time"

type Report struct {
	Version string `json:"version"`

	Generated          string `json:"generated"`
	Project            string `json:"project"`
	ProjectURL         string `json:"projectURL"`
	ProjectDescription string `json:"projectDescription"`
	Owner              string `json:"owner"`

	MethodName string `json:"methodName"`

	Build     string `json:"build"`
	BuildDate string `json:"buildDate"`

	ProductCode string `json:"productCode"`

	// not used yet
	TraceEvents []TraceEvent `json:"traceEvents"`

	Activities               []Activity
	PrepareAppInitActivities []Activity

	TotalDuration int
}

type ExtraData struct {
	LastGeneratedTime time.Time
	BuildTime         time.Time

	ProductCode string
	BuildNumber string

	Machine string

	TcBuildId          int
	TcBuildType        string
	TcInstallerBuildId int
	TcBuildProperties  []byte
	TcBuildNumber      string
	Changes            []string
	TriggeredBy        string

	CurrentBuildTime time.Time

	// for logging only
	ReportFile string
}

type TraceEvent struct {
	Name  string `json:"name"`
	Phase string `json:"ph"`
	// in microseconds
	Timestamp int `json:"ts"`

	// in old reports (v10) can be int instead of string
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
