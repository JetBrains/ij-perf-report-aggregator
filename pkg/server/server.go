package server

import (
  "context"
  "crypto/tls"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/rs/cors"
  "go.uber.org/zap"
  "net/http"
  "net/url"
  "os"
  "os/signal"
  "report-aggregator/pkg/util"
  "syscall"
  "time"
)

type StatsServer struct {
  db                       *sqlite3.Conn
  victoriaMetricsServerUrl *url.URL

  logger *zap.Logger
}

func Serve(dbPath string, victoriaMetricsServerUrl string, logger *zap.Logger) error {
  db, err := sqlite3.Open(dbPath, sqlite3.OPEN_READONLY)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  if len(victoriaMetricsServerUrl) == 0 {
    victoriaMetricsServerUrl = "http://localhost:8428"
  }

  parsedPromServerUrl, err := url.Parse(victoriaMetricsServerUrl + "/api/v1/query")
  if err != nil {
    return errors.WithStack(err)
  }

  statsServer := &StatsServer{
    db:                       db,
    victoriaMetricsServerUrl: parsedPromServerUrl,

    logger: logger,
  }

  cacheManager := NewResponseCacheManager(logger)

  mux := http.NewServeMux()

  mux.Handle("/info", cacheManager.CreateHandler(statsServer.handleInfoRequest))
  mux.Handle("/metrics/", cacheManager.CreateHandler(statsServer.handleMetricsRequest))
  mux.Handle("/groupedMetrics/", cacheManager.CreateHandler(statsServer.handleGroupedMetricsRequest))

  mux.HandleFunc("/health-check", func(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(200)
  })

  serverPort := os.Getenv("SERVER_PORT")
  if len(serverPort) == 0 {
    serverPort = "9044"
  }
  server := listenAndServe(serverPort, mux, logger)

  logger.Info("started", zap.String("address", server.Addr))

  waitUntilTerminated(server, 1*time.Minute, logger)

  return nil
}

func listenAndServe(port string, mux *http.ServeMux, logger *zap.Logger) *http.Server {
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
