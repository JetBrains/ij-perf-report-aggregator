package server

import (
  "compress/gzip"
  "context"
  "crypto/tls"
  "github.com/NYTimes/gziphandler"
  "github.com/alecthomas/kingpin"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/didip/tollbooth"
  "github.com/didip/tollbooth/limiter"
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

  requestLimit := tollbooth.NewLimiter(2, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
  requestLimit.SetBurst(10)

  mux := http.NewServeMux()

  gzipWrapper, _ := gziphandler.NewGzipLevelHandler(gzip.DefaultCompression)

  mux.Handle("/info", tollbooth.LimitHandler(requestLimit, gzipWrapper(http.HandlerFunc(statsServer.handleInfoRequest))))
  mux.Handle("/metrics", tollbooth.LimitHandler(requestLimit, gzipWrapper(http.HandlerFunc(statsServer.handleMetricsRequest))))
  mux.Handle("/groupedMetrics", tollbooth.LimitHandler(requestLimit, gzipWrapper(http.HandlerFunc(statsServer.handleGroupedMetricsRequest))))

  server := listenAndServe("9044", mux, logger)

  logger.Info("started",
    zap.String("address", server.Addr),
  )

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

    // https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b
    ReadTimeout:  4 * time.Second,
    WriteTimeout: 60 * time.Second,

    TLSConfig: &tls.Config{
      MinVersion: tls.VersionTLS12,
    },
  }

  go func() {
    var err error
    if os.Getenv("USE_SSL") == "true" {
      err = server.ListenAndServeTLS("/etc/secrets/tls.cert", "/etc/secrets/tls.key")
    } else {
      err = server.ListenAndServe()
    }
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
