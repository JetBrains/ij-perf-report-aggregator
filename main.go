package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"go.uber.org/zap"
	"log"
	"os"
	"report-aggregator/pkg/ideaLog"
	"report-aggregator/pkg/server"
)

func createLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func main() {
	logger := createLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var app = kingpin.New("report-aggregator", "report-aggregator").Version("0.0.1")

	ideaLog.ConfigureCollectFromDirCommand(app, logger)
	server.ConfigureServeCommand(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}
