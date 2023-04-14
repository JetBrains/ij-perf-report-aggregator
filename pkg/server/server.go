package server

import (
  "context"
  "crypto/tls"
  "github.com/ClickHouse/ch-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  dataquery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/andybalholm/brotli"
  "github.com/develar/errors"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/jackc/puddle/v2"
  "github.com/nats-io/nats.go"
  "github.com/rs/cors"
  "github.com/valyala/bytebufferpool"
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
  dbUrl        string
  nameToDbPool sync.Map

  machineInfo analyzer.MachineInfo

  poolMutex sync.Mutex

  logger *zap.Logger
}

func Serve(dbUrl string, natsUrl string, logger *zap.Logger) error {
  if len(dbUrl) == 0 {
    dbUrl = DefaultDbUrl
  }

  dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
  if err != nil {
    return err
  }

  statsServer := &StatsServer{
    dbUrl: dbUrl,

    logger: logger,

    machineInfo: analyzer.GetMachineInfo(),
  }

  defer func() {
    statsServer.nameToDbPool.Range(func(name, pool interface{}) bool {
      p, ok := pool.(*puddle.Pool[*ch.Client])
      if ok {
        p.Close()
      }
      return true
    })
  }()

  cacheManager, err := NewResponseCacheManager(logger)
  if err != nil {
    return err
  }

  r := chi.NewRouter()

  disposer := util.NewDisposer()
  defer disposer.Dispose()
  if len(natsUrl) > 0 {
    err = listenNats(cacheManager, natsUrl, disposer, logger)
    if err != nil {
      return err
    }
  }

  r.Use(middleware.AllowContentType("application/octet-stream", "application/json"))
  r.Use(cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST"},
    MaxAge:         50,
  }).Handler)
  r.Use(middleware.Heartbeat("/health-check"))
  r.Use(middleware.Recoverer)
  compressor := middleware.NewCompressor(5, "/*")
  compressor.SetEncoder("br", func(w io.Writer, level int) io.Writer {
    return brotli.NewWriterLevel(w, level)
  })
  r.Use(compressor.Handler)

  r.Post("/api/meta*", createPostMetaRequestHandler(logger, dbpool))
  r.Get("/api/meta/*", createGetMetaRequestHandler(logger, dbpool))
  r.Handle("/api/v1/meta/measure/*", cacheManager.CreateHandler(statsServer.handleMetaMeasureRequest))
  r.Handle("/api/v1/load/*", cacheManager.CreateHandler(statsServer.handleLoadRequest))
  r.Handle("/api/q/*", cacheManager.CreateHandler(statsServer.handleLoadRequestV2))
  r.Handle("/api/zstd-dictionary/*", &CachingHandler{
    handler: func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
      return &bytebufferpool.ByteBuffer{B: dataquery.ZstdDictionary}, false, nil
    },
    manager: cacheManager,
  })

  server := listenAndServe(util.GetEnv("SERVER_PORT", "9044"), r, logger)

  logger.Info("started", zap.String("address", server.Addr), zap.String("clickhouse", dbUrl), zap.String("nats", natsUrl))

  waitUntilTerminated(server, 1*time.Minute, logger)
  return nil
}

func (t *StatsServer) AcquireDatabase(ctx context.Context, name string) (*puddle.Resource[*ch.Client], error) {
  untypedPool, exists := t.nameToDbPool.Load(name)
  var pool *puddle.Pool[*ch.Client]
  var err error
  isCorrectPool := true
  if exists {
    pool, isCorrectPool = untypedPool.(*puddle.Pool[*ch.Client])
  } else {
    pool, err = createStoreForDatabaseUnderLock(name, t)
  }
  if err != nil {
    return nil, errors.WithStack(err)
  }
  if !isCorrectPool {
    return nil, errors.New("Pool can't be casted to (*puddle.Pool[*ch.Client])")
  }

  resource, err := pool.Acquire(ctx)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return resource, nil
}

func createStoreForDatabaseUnderLock(name string, t *StatsServer) (*puddle.Pool[*ch.Client], error) {
  t.poolMutex.Lock()
  defer t.poolMutex.Unlock()
  pool, err := puddle.NewPool(&puddle.Config[*ch.Client]{
    MaxSize: 16,
    Destructor: func(value *ch.Client) {
      _ = value.Close()
    },
    Constructor: func(ctx context.Context) (*ch.Client, error) {
      client, err := ch.Dial(ctx, ch.Options{
        Address:  t.dbUrl,
        Database: name,
        Settings: []ch.Setting{
          ch.SettingInt("readonly", 1),
          ch.SettingInt("max_query_size", 1000000),
          ch.SettingInt("max_memory_usage", 3221225472),
        },
      })
      return client, err
    },
  })
  if err == nil && pool != nil {
    t.nameToDbPool.Store(name, pool)
  }
  return pool, err
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

    ReadTimeout:  30 * time.Second,
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
