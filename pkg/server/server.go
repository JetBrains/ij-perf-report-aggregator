package server

import (
  "context"
  "crypto/tls"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/nats-io/nats.go"
  "github.com/rs/cors"
  "go.uber.org/zap"
  "io"
  "net/http"
  "os"
  "os/signal"
  "sync"
  "syscall"
  "time"
)

const DefaultDbUrl = "127.0.0.1:9000"

type StatsServer struct {
  dbUrl    string
  nameToDb sync.Map

  machineInfo analyzer.MachineInfo

  logger *zap.Logger
}

func Serve(dbUrl string, natsUrl string, logger *zap.Logger) error {
  if len(dbUrl) == 0 {
    dbUrl = DefaultDbUrl
  }

  statsServer := &StatsServer{
    dbUrl: dbUrl,

    logger: logger,

    machineInfo: analyzer.GetMachineInfo(),
  }

  defer func() {
    statsServer.nameToDb.Range(func(name, db interface{}) bool {
      util.Close(db.(io.Closer), logger)
      return true
    })
  }()

  cacheManager, err := NewResponseCacheManager(logger)
  if err != nil {
    return err
  }

  mux := http.NewServeMux()

  disposer := util.NewDisposer()
  defer disposer.Dispose()
  if len(natsUrl) > 0 {
    err = listenNats(cacheManager, natsUrl, disposer, logger)
    if err != nil {
      return err
    }
  }

  mux.Handle("/api/v1/meta/measure", cacheManager.CreateHandler(statsServer.handleMetaMeasureRequest))
  mux.Handle("/api/v1/load/", cacheManager.CreateHandler(statsServer.handleLoadRequest))
  mux.Handle("/api/v1/compareMetrics", cacheManager.CreateHandler(statsServer.handleStatusRequest))
  mux.Handle("/api/v1/report/", cacheManager.CreateHandler(statsServer.handleReportRequest))

  mux.HandleFunc("/health-check", func(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(200)
  })

  server := listenAndServe(util.GetEnv("SERVER_PORT", "9044"), mux, logger)

  logger.Info("started", zap.String("address", server.Addr), zap.String("clickhouse", dbUrl), zap.String("nats", natsUrl))

  waitUntilTerminated(server, 1*time.Minute, logger)
  return nil
}

func (t *StatsServer) GetDatabase(name string) (*sqlx.DB, error) {
  db, _ := t.nameToDb.Load(name)
  if db != nil {
    return db.(*sqlx.DB), nil
  }

  // limit max memory to use - 3GB
  // increase max query size, because we IN statements contains a lot of values (e.g. machines)
  db, err := sqlx.Open("clickhouse", "tcp://"+t.dbUrl+"?read_timeout=45&write_timeout=45&compress=1&readonly=1&max_query_size=1000000&max_memory_usage=3221225472&database="+name)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  actual, loaded := t.nameToDb.LoadOrStore(name, db)
  if loaded {
    // not stored because another thread already stored another value, so, close candidate
    util.Close(db.(*sqlx.DB), t.logger)
  }
  return actual.(*sqlx.DB), nil
}

func listenNats(cacheManager *ResponseCacheManager, natsUrl string, disposer *util.Disposer, logger *zap.Logger) error {
  // wait when nats service will be deployed
  nc, err := nats.Connect(natsUrl, nats.Timeout(30*time.Second))
  if err != nil {
    return errors.WithStack(err)
  }

  ncSubscription, err := nc.Subscribe("server.clearCache", func(m *nats.Msg) {
    cacheManager.Clear()
    logger.Info("cache cleared", zap.ByteString("sender", m.Data))
  })

  if err != nil {
    return errors.WithStack(err)
  }

  disposer.Add(func() {
    err := ncSubscription.Unsubscribe()
    if err != nil {
      logger.Error("cannot unsubscribe", zap.Error(err))
    }
  })
  return nil
}

func listenAndServe(port string, mux http.Handler, logger *zap.Logger) *http.Server {
  // buffer size is 4096 https://github.com/golang/go/issues/13870
  server := &http.Server{
    Addr:    ":" + port,
    Handler: cors.Default().Handler(mux),

    ReadTimeout:  4 * time.Second,
    WriteTimeout: 60 * time.Second,

    TLSConfig: &tls.Config{
      MinVersion: tls.VersionTLS12,
    },
  }

  go func() {
    err := server.ListenAndServe()
    if err == http.ErrServerClosed {
      logger.Debug("server closed")
    } else {
      logger.Fatal("cannot serve", zap.Error(err), zap.String("port", port))
    }
  }()

  return server
}

func waitUntilTerminated(server *http.Server, shutdownTimeout time.Duration, logger *zap.Logger) {
  signals := make(chan os.Signal, 1)
  signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
  <-signals

  shutdownHttpServer(server, shutdownTimeout, logger)
}

func shutdownHttpServer(server *http.Server, shutdownTimeout time.Duration, logger *zap.Logger) {
  if server == nil {
    return
  }

  ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
  defer cancel()
  logger.Info("shutdown server", zap.Duration("timeout", shutdownTimeout))
  start := time.Now()
  err := server.Shutdown(ctx)
  if err != nil {
    logger.Error("cannot shutdown server", zap.Error(err))
    return
  }

  logger.Info("server is shutdown", zap.Duration("duration", time.Since(start)))
}
