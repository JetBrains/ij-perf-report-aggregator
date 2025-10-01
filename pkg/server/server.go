package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/auth"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/puddle/v2"
	"github.com/nats-io/nats.go"
	"github.com/rs/cors"
	"github.com/valyala/bytebufferpool"
	"go.deanishe.net/env"
)

const DefaultDbUrl = "127.0.0.1:9000"

type StatsServer struct {
	dbUrl        string
	nameToDbPool sync.Map

	poolMutex sync.Mutex
}

func Serve(dbUrl string, natsUrl string) error {
	if dbUrl == "" {
		dbUrl = DefaultDbUrl
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	dbpool.Config().MaxConns = 10
	if err != nil {
		return err
	}

	statsServer := &StatsServer{
		dbUrl: dbUrl,
	}

	defer func() {
		statsServer.nameToDbPool.Range(func(_, pool interface{}) bool {
			p, ok := pool.(*puddle.Pool[*ch.Client])
			if ok {
				p.Close()
			}
			return true
		})
	}()

	cacheManager, err := NewResponseCacheManager()
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	disposer := util.NewDisposer()
	defer disposer.Dispose()
	if natsUrl != "" {
		err = listenNats(cacheManager, natsUrl, disposer)
		if err != nil {
			return err
		}
	} else {
		slog.Info("no nats server configured")
	}

	router.Use(middleware.AllowContentType("application/octet-stream", "application/json"))
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"*"},
		MaxAge:         50,
	}).Handler)
	router.Use(middleware.Heartbeat("/health-check"))
	router.Use(middleware.Recoverer)

	router.Route("/api/meta", func(r chi.Router) {
		r.Route("/accidents", func(r chi.Router) {
			r.Post("/*", meta.CreatePostAccidentRequestHandler(dbpool))
			r.Delete("/*", meta.CreateDeleteAccidentRequestHandler(dbpool))
			r.Get("/*", meta.CreateGetAccidentByIdHandler(dbpool))
		})
		r.Post("/getAccidents*", meta.CreateGetManyAccidentsRequestHandler(dbpool))
		r.Get("/description*", meta.CreateGetDescriptionRequestHandler(dbpool))
		r.Post("/accidentsAroundDate*", meta.CreateGetAccidentsAroundDateRequestHandler(dbpool))
		r.Post("/missingData", meta.CreatePostMissingDataRequestHandler(dbpool))
		r.Route("/youtrack", func(r chi.Router) {
			r.Post("/createIssue", meta.CreatePostCreateIssueByAccident(dbpool))
			r.Post("/uploadAttachments", meta.CreatePostUploadAttachmentsToIssue())
		})
		r.Route("/teamcity", func(r chi.Router) {
			r.Post("/startBisect", meta.CreatePostStartBisect())
			r.Get("/changes", meta.HandleGetTeamCityChanges())
			r.Get("/buildType", meta.HandleGetTeamCityBuildType())
			r.Get("/buildCounter", meta.HandleGetTeamCityBuildCounter())
		})
	})

	router.Route("/api/auth", func(r chi.Router) {
		r.Route("/userinfo", func(r chi.Router) {
			r.Get("/*", auth.CreateGetUserInfoHandler())
		})
	})

	router.Post("/api/evaluateMetric*", statsServer.CreateProcessMetricDataHandler())

	router.Group(func(r chi.Router) {
		compressor := middleware.NewCompressor(5)
		compressor.SetEncoder("br", func(w io.Writer, level int) io.Writer {
			return brotli.NewWriterV2(w, level)
		})
		r.Use(compressor.Handler)

		r.Route("/api/", func(r chi.Router) {
			r.Route("/v1", func(r chi.Router) {
				r.Handle("/meta/measure", cacheManager.CreateHandler(statsServer.handleMetaMeasureRequest))
				r.Handle("/load/*", cacheManager.CreateHandler(statsServer.handleLoadRequest))
			})
			r.Handle("/q/*", cacheManager.CreateHandler(statsServer.handleLoadRequestV2))
			r.Handle("/highlightingPasses*", cacheManager.CreateHandler(statsServer.getDistinctHighlightingPasses))
			r.Handle("/compareBranches*", cacheManager.CreateHandler(statsServer.getBranchComparison))
			r.Handle("/compareModes*", cacheManager.CreateHandler(statsServer.getModeComparison))
			r.Handle("/zstd-dictionary/*", &CachingHandler{
				handler: func(_ *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
					return &bytebufferpool.ByteBuffer{B: util.ZstdDictionary}, false, nil
				},
				manager: cacheManager,
			})
		})
	})

	server := listenAndServe(env.Get("SERVER_PORT", "9044"), router)
	slog.Info("started", "server", server.Addr, "clickhouse", dbUrl, "nats", natsUrl)

	waitUntilTerminated(server, 1*time.Minute)
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
		return nil, fmt.Errorf("cannot create pool: %w", err)
	}
	if !isCorrectPool {
		return nil, errors.New("pool can't be casted to (*puddle.Pool[*ch.YoutrackClient])")
	}

	resource, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire from pool: %w", err)
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

func listenNats(cacheManager *ResponseCacheManager, natsUrl string, disposer *util.Disposer) error {
	// wait when nats service will be deployed
	nc, err := nats.Connect(natsUrl, nats.Timeout(30*time.Second))
	if err != nil {
		return fmt.Errorf("can't connect to nats: %w", err)
	}

	ncSubscription, err := nc.Subscribe("server.clearCache", func(m *nats.Msg) {
		cacheManager.Clear()
		slog.Info("cache cleared", "sender", m.Data)
	})
	if err != nil {
		return fmt.Errorf("can't subscribe to nats: %w", err)
	}

	disposer.Add(func() {
		err := ncSubscription.Unsubscribe()
		if err != nil {
			slog.Error("cannot unsubscribe", "error", err)
		}
	})
	return nil
}

func listenAndServe(port string, mux http.Handler) *http.Server {
	// buffer size is 4096 https://github.com/golang/go/issues/13870
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,

		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,

		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Debug("server closed")
		} else {
			slog.Error("cannot serve", "error", err, "port", port)
			os.Exit(1)
		}
	}()

	return server
}

func waitUntilTerminated(server *http.Server, shutdownTimeout time.Duration) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	shutdownHttpServer(server, shutdownTimeout)
}

func shutdownHttpServer(server *http.Server, shutdownTimeout time.Duration) {
	if server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	slog.Info("shutdown server", "timeout", shutdownTimeout)
	start := time.Now()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("cannot shutdown server", "error", err)
		return
	}

	slog.Info("server is shutdown", "duration", time.Since(start))
}
