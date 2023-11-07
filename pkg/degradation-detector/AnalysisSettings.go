package degradation_detector

type Settings struct {
  Test                   string
  Metric                 string
  Db                     string
  Table                  string
  Branch                 string
  Machine                string
  Channel                string
  ProductLink            string
  DoNotReportImprovement bool
}
