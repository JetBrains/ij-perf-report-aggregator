package filling

type MetricResult struct {
  Product string
  Machine uint8

  GeneratedTime int64
  BuildTime     int64

  TcBuildId          int
  TcInstallerBuildId int
  TcBuildProperties  []byte

  RawReport string

  BuildC1 int `db:"build_c1"`
  BuildC2 int `db:"build_c2"`
  BuildC3 int `db:"build_c3"`
}
