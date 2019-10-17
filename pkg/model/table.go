package model

import (
  "database/sql"
  "github.com/develar/errors"
  "strings"
)

var EssentialDurationMetricNames = []string{"bootstrap", "appInitPreparation", "appInit", "pluginDescriptorLoading", "appComponentCreation", "projectComponentCreation"}
var DurationMetricNames = append(EssentialDurationMetricNames, "moduleLoading")
var InstantMetricNames = []string{"splash", "startUpCompleted"}

// https://clickhouse.yandex/docs/en/query_language/alter/#manipulations-with-key-expressions
// To keep the property that data part rows are ordered by the sorting key expression
// you cannot add expressions containing existing columns to the sorting key (only columns added by the ADD COLUMN command in the same ALTER query).

func ProcessMetricName(handler func(name string, isInstant bool)) {
  for _, name := range DurationMetricNames {
    handler(name, false)
  }
  for _, name := range InstantMetricNames {
    handler(name, true)
  }
}

func CreateTable(db *sql.DB) error {
  _, err := db.Exec("set allow_experimental_data_skipping_indices = 1")
  if err != nil {
    return errors.WithStack(err)
  }

  // https://www.altinity.com/blog/2019/7/new-encodings-to-improve-clickhouse
  // see zstd-compression-level.txt
  var sb strings.Builder
  sb.WriteString(`
  create table if not exists report (
    product FixedString(2) Codec(ZSTD(19)),
    machine String Codec(ZSTD(19)),

    build_time DateTime Codec(Delta, ZSTD(19)),
    generated_time DateTime Codec(Delta, ZSTD(19)),
    
    tc_build_id UInt32 Codec(DoubleDelta, ZSTD(19)),
    
    raw_report String Codec(ZSTD(19)),
    
    build_c1 UInt8 Codec(DoubleDelta, ZSTD(19)),
    build_c2 UInt16 Codec(DoubleDelta, ZSTD(19)),
    build_c3 UInt16 Codec(DoubleDelta, ZSTD(19))
  `)
  ProcessMetricName(func(name string, isInstant bool) {
    sb.WriteRune(',')
    sb.WriteRune(' ')
    sb.WriteString(name)
    sb.WriteRune('_')
    if isInstant {
      sb.WriteRune('i')
    } else {
      sb.WriteRune('d')
    }
    sb.WriteRune(' ')
    //if _, ok := MetricToUint16DataType[name]; ok  {
    if !isInstant {
      sb.WriteString("UInt16")
    } else {
      sb.WriteString("Int32")
    }
    sb.WriteRune(' ')
    sb.WriteString("Codec(Gorilla, ZSTD(19))")
  })

  // https://github.com/ClickHouse/ClickHouse/issues/3758#issuecomment-444490724
  sb.WriteString(") engine MergeTree partition by (product, toYYYYMM(generated_time)) order by (product, machine, build_c1, build_c2, build_c3, build_time, generated_time) SETTINGS old_parts_lifetime = 10")

  _, err = db.Exec(sb.String())
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}
