package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/develar/errors"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"report-aggregator/pkg/analyzer"
	"report-aggregator/pkg/util"
	"strings"
	"time"
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

	configureCollectFromDirCommand(app, logger)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}

func configureCollectFromDirCommand(app *kingpin.Application, log *zap.Logger) {
	command := app.Command("collect", "Collect reports from idea.log files")
	dir := command.Flag("dir", "The input directory").Short('i').Required().String()
	dbPath := command.Flag("db", "The output SQLite database file.").Short('o').Required().String()
	command.Action(func(context *kingpin.ParseContext) error {
		err := collectFromDir(*dir, *dbPath, log)
		if err != nil {
			return err
		}

		return nil
	})
}

func collectFromDir(dir string, dbPath string, logger *zap.Logger) error {
	taskContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()

	files, err := filepath.Glob(dir + "/idea*.log*")
	if err != nil {
		return errors.WithStack(err)
	}

	reportAnalyzer, err := analyzer.CreateReportAnalyzer(dbPath, taskContext, logger)
	if err != nil {
		return err
	}

	go func() {
		err = <-reportAnalyzer.ErrorChannel
		cancel()

		if err != nil {
			logger.Error("cannot analyze", zap.Error(err))
		}
	}()

	defer util.Close(reportAnalyzer, logger)

	for _, file := range files {
		if taskContext.Err() != nil {
			return nil
		}

		err := collectFromLogFile(file, logger, reportAnalyzer, taskContext)
		if err != nil {
			return err
		}
	}

	select {
	case analyzeError := <-reportAnalyzer.ErrorChannel:
		cancel()
		return analyzeError

	case <-reportAnalyzer.Done():
		cancel()
		return nil

	case <-taskContext.Done():
		return nil
	}
}

func collectFromLogFile(filePath string, log *zap.Logger, reportAnalyzer *analyzer.ReportAnalyzer, taskContext context.Context) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(file, log)

	scanner := bufio.NewScanner(file)
	var jsonData bytes.Buffer
	state := 0

	startSuffix := []byte("=== Start: StartUp Measurement ===")
	endSuffix := []byte("=== Stop: StartUp Measurement ===")

	lastGeneratedTime := int64(-1)
	for scanner.Scan() {
		line := scanner.Bytes()
		if state == 1 {
			// idea start-up performance writer bug - end suffix has extra trailing space, so, HasSuffix cannot be used
			if bytes.HasPrefix(line, endSuffix) {
				if taskContext.Err() != nil {
					return nil
				}

				err = reportAnalyzer.Analyze(jsonData.Bytes(), lastGeneratedTime)
				if err != nil {
					return err
				}

				state = 0
				lastGeneratedTime = -1
				jsonData.Reset()
			} else {
				jsonData.Write(line)
			}
		} else if bytes.HasSuffix(line, startSuffix) {
			lineString := scanner.Text()
			// UTC, but it is ok, modern reports contain correct generated time
			parsedTime, err := time.Parse("2006-01-02 15:04:05", lineString[0:strings.IndexRune(lineString, ',')])
			if err != nil {
				return errors.WithStack(err)
			}

			state = 1
			lastGeneratedTime = parsedTime.Unix()
		}
	}

	return nil
}
