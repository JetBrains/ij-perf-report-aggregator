package server

import (
  "context"
  "crypto/tls"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/nats-io/nats.go"
  "github.com/rs/cors"
  "go.uber.org/zap"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
)

type StatsServer struct {
  db *sqlx.DB

  machineInfo MachineInfo

  logger *zap.Logger
}

func Serve(dbUrl string, natsUrl string, logger *zap.Logger) error {
  if len(dbUrl) == 0 {
    dbUrl = "127.0.0.1:9000"
  }

  chDb, err := sqlx.Open("clickhouse", "tcp://"+dbUrl+"?read_timeout=45&write_timeout=45&compress=1")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(chDb, logger)

  statsServer := &StatsServer{
    db: chDb,

    logger: logger,

    machineInfo: getMachineInfo(),
  }

  cacheManager := NewResponseCacheManager(logger)

  mux := http.NewServeMux()

  disposer := util.NewDisposer()
 	defer disposer.Dispose()
  if len(natsUrl) > 0 {
    err = listenNats(cacheManager, natsUrl, disposer, logger)
    if err != nil {
      return err
    }
  }

  mux.Handle("/api/v1/info", cacheManager.CreateHandler(statsServer.handleInfoRequest))
  mux.Handle("/api/v1/metrics/", cacheManager.CreateHandler(statsServer.handleMetricsRequest))
  mux.Handle("/api/v1/groupedMetrics/", cacheManager.CreateHandler(statsServer.handleGroupedMetricsRequest))
  mux.Handle("/api/v1/report/", cacheManager.CreateHandler(statsServer.handleReportRequest))

  mux.HandleFunc("/health-check", func(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(200)
  })

  serverPort := os.Getenv("SERVER_PORT")
  if len(serverPort) == 0 {
    serverPort = "9044"
  }
  server := listenAndServe(serverPort, mux, logger)

  logger.Info("started", zap.String("address", server.Addr), zap.String("clickhouse", dbUrl), zap.String("nats", natsUrl))

  waitUntilTerminated(server, 1*time.Minute, logger)
  return nil
}

func listenNats(cacheManager *ResponseCacheManager, natsUrl string, disposer *util.Disposer, logger *zap.Logger) error {
  nc, err := nats.Connect(natsUrl)
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
