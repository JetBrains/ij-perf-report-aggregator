package model

import "time"

type Report struct {
	Generated          string `json:"generated"`
	Project            string `json:"project"`
	ProjectURL         string `json:"projectURL"`
	ProjectDescription string `json:"projectDescription"`
	Owner              string `json:"owner"`

	MethodName string `json:"methodName"`

	Build     string `json:"build"`
	BuildDate string `json:"buildDate"`
}

type ExtraData struct {
	LastGeneratedTime time.Time
	BuildTime         time.Time

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
