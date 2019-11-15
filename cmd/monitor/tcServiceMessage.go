package main

import (
  "fmt"
  "os"
)

type TeamCityActivity struct {
  name string

  startMarker string
  endMarker   string
}

func (t *TeamCityActivity) Start(value string) {
  if len(t.name) == 0 || t.name != value {
    t.End()
    t.name = value
    _, _ = fmt.Fprintf(os.Stdout, "##teamcity[%s name='%s']\n", t.startMarker, t.name)
  }
}

func (t *TeamCityActivity) End() {
  if len(t.name) != 0 {
    _, _ = fmt.Fprintf(os.Stdout, "##teamcity[%s name='%s']\n\n", t.endMarker, t.name)
  }
}

type TeamCityTest struct {
  TeamCityActivity
}

func (t *TeamCityTest) Error(message string) {
  _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testStdErr name='%s' out='%s']\n", t.name, message)
}

func (t *TeamCityTest) Output(message string) {
  _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testStdOut name='%s' out='%s']\n", t.name, message)
}

func (t *TeamCityTest) Failed(message string) {
  _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testFailed name='%s' message='%s']\n", t.name, message)
}

func (t *TeamCityTest) CompareFailed(message string, actual float64, expected float64) {
  _, _ = fmt.Fprintf(os.Stdout, "##teamcity[testFailed type='comparisonFailure' name='%s' message='%s' expected='%.1f' actual='%.1f']\n", t.name, message, expected, actual)
}
