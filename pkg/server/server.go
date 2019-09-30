package server

import (
	"context"
	"crypto/tls"
	"github.com/NYTimes/gziphandler"
	"github.com/alecthomas/kingpin"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/develar/errors"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/json-iterator/go"
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

	requestLimit := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	requestLimit.SetBurst(10)

	mux := http.NewServeMux()
	mux.Handle("/info", gziphandler.GzipHandler(tollbooth.LimitFuncHandler(requestLimit, statsServer.handleInfoRequest)))
	mux.Handle("/stats", gziphandler.GzipHandler(tollbooth.LimitFuncHandler(requestLimit, statsServer.handleStatsRequest)))

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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: cors.Default().Handler(mux),
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

func (t *StatsServer) handleInfoRequest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonWriter := jsoniter.NewStream(jsoniter.ConfigFastest, w, 64*1024)

	jsonWriter.WriteObjectStart()
	jsonWriter.WriteObjectField("productNames")

	productNames, err := t.getProductNames()
	if err != nil {
		t.httpError(err, w)
		return
	}

	jsonWriter.WriteArrayStart()
	isFirst := true
	for _, product := range productNames {
		if isFirst {
			isFirst = false
		} else {
			jsonWriter.WriteMore()
		}
		jsonWriter.WriteString(product)
	}
	jsonWriter.WriteArrayEnd()

	jsonWriter.WriteMore()
	jsonWriter.WriteObjectField("productToMachineNames")
	jsonWriter.WriteObjectStart()

	statement, err := t.db.Prepare("select distinct machine from report where product = ? order by machine")
	if err != nil {
		t.httpError(err, w)
		return
	}

	isFirst = true
	for _, product := range productNames {
		err = statement.Bind(product)
		if err != nil {
			t.httpError(err, w)
			return
		}

		if isFirst {
			isFirst = false
		} else {
			jsonWriter.WriteMore()
		}
		jsonWriter.WriteObjectField(product)

		jsonWriter.WriteArrayStart()
		err = writeStringList(jsonWriter, statement)
		if err != nil {
			t.httpError(err, w)
			return
		}
		jsonWriter.WriteArrayEnd()
	}

	jsonWriter.WriteObjectEnd()

	jsonWriter.WriteObjectEnd()

	err = jsonWriter.Flush()
	if err != nil {
		t.httpError(err, w)
		return
	}
}

func (t *StatsServer) handleStatsRequest(w http.ResponseWriter, _ *http.Request) {
	//noinspection SqlResolve
	statement, err := t.db.Prepare(`
select machine, generated_time, metrics, 
       json_extract(raw_report, '$.productCode') as productCode,
       json_extract(raw_report, '$.build') as build
from report order by machine, generated_time
	`)
	if err != nil {
		t.logger.Error("cannot query", zap.Error(err))
		t.httpError(err, w)
		return
	}

	defer util.Close(statement, t.logger)

	w.Header().Set("Content-Type", "application/json")
	jsonWriter := jsoniter.NewStream(jsoniter.ConfigFastest, w, 64*1024)

	jsonWriter.WriteArrayStart()

	var lastMachine sqlite3.RawString
	isFirst := true
	for {
		hasRow, err := statement.Step()
		if err != nil {
			t.httpError(err, w)
			return
		}

		if !hasRow {
			break
		}

		var machine sqlite3.RawString
		var generatedTime int64
		var metrics sqlite3.RawString
		var productCode string
		var build string
		err = statement.Scan(&machine, &generatedTime, &metrics, &productCode, &build)
		if err != nil {
			t.httpError(err, w)
			return
		}

		if lastMachine != machine {
			if len(lastMachine) != 0 {
				jsonWriter.WriteArrayEnd()
				jsonWriter.WriteObjectEnd()
			}

			jsonWriter.WriteObjectStart()

			jsonWriter.WriteObjectField("machine")
			jsonWriter.WriteString(string(machine))

			jsonWriter.WriteMore()
			jsonWriter.WriteObjectField("metrics")
			jsonWriter.WriteArrayStart()

			lastMachine = machine
			isFirst = true
		}

		if isFirst {
			isFirst = false
		} else {
			jsonWriter.WriteMore()
		}

		jsonWriter.WriteObjectStart()
		// timestamp
		jsonWriter.WriteObjectField("t")
		// seconds to milliseconds
		jsonWriter.WriteInt64(generatedTime * 1000)
		jsonWriter.WriteMore()

		if len(productCode) != 0 {
			jsonWriter.WriteObjectField("_p")
			jsonWriter.WriteString(productCode)
			jsonWriter.WriteMore()

			jsonWriter.WriteObjectField("_v")
			jsonWriter.WriteString(build)
			jsonWriter.WriteMore()
		}

		// skip first '{'
		jsonWriter.WriteRaw(string(metrics[1:]))
	}

	// close machine array
	jsonWriter.WriteArrayEnd()
	// close machine object entry
	jsonWriter.WriteObjectEnd()

	// close root array
	jsonWriter.WriteArrayEnd()

	err = jsonWriter.Flush()
	if err != nil {
		t.logger.Error("cannot flush", zap.Error(err))
		return
	}
}
