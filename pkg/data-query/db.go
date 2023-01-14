package data_query

import (
  "context"
  "errors"
  "github.com/ClickHouse/ch-go"
  "github.com/ClickHouse/ch-go/chpool"
  "github.com/ClickHouse/ch-go/proto"
  "io"
  "net"
)

type DatabaseConnectionSupplier interface {
  AcquireDatabase(name string, ctx context.Context) (*chpool.Client, error)
}

func executeQuery(
  sqlQuery string,
  query DataQuery,
  dbSupplier DatabaseConnectionSupplier,
  ctx context.Context,
  resultHandler func(ctx context.Context, block proto.Block, result *proto.Results) error,
) error {
  for attempt := 0; attempt <= 8; attempt++ {
    dbResource, err := dbSupplier.AcquireDatabase(query.Database, ctx)
    if err != nil {
      return err
    }

    err, done := doExecution(sqlQuery, dbResource, ctx, resultHandler)
    if err != nil {
      return err
    }

    if done {
      return nil
    }
  }

  return errors.New("cannot acquire database")
}

func doExecution(
  sqlQuery string,
  client *chpool.Client,
  ctx context.Context,
  resultHandler func(ctx context.Context, block proto.Block, result *proto.Results) error,
) (error, bool) {
  isDestroyed := false
  defer func() {
    if !isDestroyed {
      client.Release()
    }
  }()

  var result proto.Results
  err := client.Do(ctx, ch.Query{
    Body:   sqlQuery,
    Result: result.Auto(),
    OnResult: func(ctx context.Context, block proto.Block) error {
      return resultHandler(ctx, block, &result)
    },
  })

  if err == nil {
    return nil, true
  }

  // if net error or io error - connection was closed due to inactivity, destroy it and acquire a new one
  if !isNetError(err) {
    return err, true
  }

  isDestroyed = true
  client.Release()
  return nil, false
}

func isNetError(err error) bool {
  for err != nil {
    if err == io.ErrUnexpectedEOF {
      return true
    }

    if _, ok := err.(net.Error); ok {
      return true
    }

    err = errors.Unwrap(err)
  }
  return false
}
