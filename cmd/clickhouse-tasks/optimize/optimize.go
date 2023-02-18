package main

import (
  "context"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "golang.org/x/sync/errgroup"
  "log"
)

func main() {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  err := execute(taskContext)
  if err != nil {
    log.Fatalf("%+v", err)
  }
}

type Item struct {
  Table    string `ch:"table"`
  Database string `ch:"database"`
}

func execute(taskContext context.Context) error {
  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{"127.0.0.1:9000"},
  })
  if err != nil {
    return fmt.Errorf("%w", err)
  }

  var result []Item
  err = db.Select(taskContext, &result, "select database, table from system.tables where database != 'system' and is_temporary = 0 and engine like '%MergeTree' order by database, name")
  if err != nil {
    return fmt.Errorf("%w", err)
  }

  log.Printf("optimize %d tables in parallel", len(result))
  group, ctx := errgroup.WithContext(taskContext)
  group.SetLimit(4)
  for index, item := range result {
    // if item.Table == "report" && item.Database == "ij" {
    //   continue
    // }
    // if item.Database != "perfint" || item.Table == "collector_state" || item.Table == "installer" {
    //   continue
    // }
    // if item.Database != "perfint" || !strings.HasSuffix(item.Table, "2") {
    //   continue
    // }

    group.Go(func(index int, item Item) func() error {
      return func() error {
        // log.Printf("optimize %s.%s", item.Database, item.Table)
        query := fmt.Sprintf("optimize table %s.%s", item.Database, item.Table)
        // query := fmt.Sprintf("alter table %s.%s reset setting storage_policy", item.Database, item.Table)
        // query := "create table IF NOT EXISTS " + item.Database + "." + item.Table + "2 as " + item.Database + "." + item.Table
        // query := "insert into " + item.Database + "." + item.Table + "2 select * from " + item.Database + "." + item.Table
        // query := "alter table " + item.Database + "." + item.Table + " drop column if exists raw_report"
        // query := "rename table " + item.Database + "." + item.Table + " to " + item.Database + "." + strings.TrimSuffix(item.Table, "2")
        // query := fmt.Sprintf("BACKUP TABLE %s.%s TO Disk('backups', '%s_%s')", item.Database, item.Table, item.Database, item.Table)
        log.Println(query)
        err = db.Exec(ctx, query)
        if err != nil {
          return fmt.Errorf("%w", err)
        }

        log.Printf("optimized %s.%s (%d of %d)", item.Database, item.Table, index, len(result))
        return nil
      }
    }(index, item))
  }
  return group.Wait()
}
