package main

import (
  "context"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "log"
  "os"
  "strings"
  "time"
)

func main() {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  err := execute(taskContext)
  if err != nil {
    log.Fatal(fmt.Sprintf("%+v", err))
  }
}

func execute(taskContext context.Context) error {
  backupDir := "/Volumes/data/ch-backup"
  entries, err := os.ReadDir(backupDir)
  if err != nil {
    return fmt.Errorf("%w", err)
  }

  db, err := clickhouse.Open(&clickhouse.Options{
    Addr:            []string{"127.0.0.1:9000"},
    ConnMaxLifetime: 6 * time.Hour,
    DialTimeout:     time.Hour,
    ReadTimeout:     time.Hour,
    Settings: map[string]interface{}{
      // https://github.com/ClickHouse/ClickHouse/issues/2833
      // ZSTD 19+ is used, read/write timeout should be quite large (10 minutes)
      "send_timeout":     30_000,
      "receive_timeout":  3000,
      "max_memory_usage": 100000000000,
    },
  })
  if err != nil {
    return fmt.Errorf("%w", err)
  }

  for _, entry := range entries {
    name := entry.Name()
    if name == "backup" || name[0] == '.' {
      continue
    }

    dbAndTable := strings.SplitN(name, "_", 2)
    dbName := dbAndTable[0]
    tableName := dbAndTable[1]

    err = db.Exec(taskContext, "create database IF NOT EXISTS "+dbName)
    if err != nil {
      return fmt.Errorf("%w", err)
    }

    query := "RESTORE TABLE " + dbName + "." + tableName + " FROM Disk('backups', '" + name + "')"
    log.Println(query)
    err = db.Exec(taskContext, query)
    if err != nil {
      return fmt.Errorf("%w", err)
    }

    //tableSqlFile := filepath.Join(backupDir, name, "metadata", dbName, tableName+".sql")
    //sql, err := os.ReadFile(tableSqlFile)
    //if err != nil {
    //  return fmt.Errorf("%w", err)
    //}
    //
    //fixedSql := bytes.Replace(sql, []byte(", storage_policy = 's3'"), []byte(""), 1)
    //err = os.WriteFile(tableSqlFile, fixedSql, 0666)
    //if err != nil {
    //  return fmt.Errorf("%w", err)
    //}
  }
  return nil
}
