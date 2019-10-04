package server

import (
  "context"
  "crypto/tls"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/rs/cors"
  "go.uber.org/zap"
  "net/http"
  "os"
  "os/signal"
  "report-aggregator/pkg/util"
  "syscall"
  "time"
)

func ConfigureServeCommand(app *kingpin.Application, log *zap.Logger) {
  command := app.Command("serve", "Serve SQLite database.")
  dbPath := command.Flag("db", "The SQLite database file.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    err := serve(*dbPath, log)
    if err != nil {
      return err
    }

    return nil
  })
}

func serve(dbPath string, logger *zap.Logger) error {
  db, err := sqlite3.Open(dbPath, sqlite3.OPEN_READONLY)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, logger)

  statsServer := &StatsServer{
    logger: logger,
    db:     db,
  }

  cacheManager := NewResponseCacheManager()

  mux := http.NewServeMux()

  mux.Handle("/info", cacheManager.CreateHandler(statsServer.handleInfoRequest))
  mux.Handle("/metrics", cacheManager.CreateHandler(statsServer.handleMetricsRequest))
  mux.Handle("/groupedMetrics", cacheManager.CreateHandler(statsServer.handleGroupedMetricsRequest))

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
  http.HandleFunc("/health-check", func(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(200)
  })

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

type StatsServer struct {
  db     *sqlite3.Conn
  logger *zap.Logger
}

func (t *StatsServer) httpError(err error, w http.ResponseWriter) {
  t.logger.Error("internal error", zap.Error(err))
  http.Error(w, err.Error(), 503)
}
