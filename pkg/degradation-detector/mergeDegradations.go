package degradation_detector

type degradationKey struct {
  slackChannel string
  metric       string
  build        string
}

type MultipleDegradationWithContext struct {
  Details  []Degradation
  Settings Settings
}

func MergeDegradations(degradations <-chan DegradationWithContext) chan MultipleDegradationWithContext {
  c := make(chan MultipleDegradationWithContext)
  go func() {
    m := make(map[degradationKey]MultipleDegradationWithContext, 100)
    for degradation := range degradations {
      key := degradationKey{
        slackChannel: degradation.Settings.SlackChannel(),
        metric:       degradation.Settings.GetMetric(),
        build:        degradation.Details.Build,
      }
      if existing, found := m[key]; found {
        d := make([]Degradation, 0, len(existing.Details)+1)
        copy(d, existing.Details)
        d = append(d, degradation.Details)
        m[key] = MultipleDegradationWithContext{
          Details:  d,
          Settings: existing.Settings.MergeAnother(degradation.Settings),
        }
      } else {
        d := make([]Degradation, 0)
        d = append(d, degradation.Details)
        m[key] = MultipleDegradationWithContext{
          Details:  d,
          Settings: degradation.Settings,
        }
      }
    }
    for _, value := range m {
      c <- value
    }
    close(c)
  }()
  return c
}
