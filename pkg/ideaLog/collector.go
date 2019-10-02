package ideaLog

import (
  "bufio"
  "bytes"
  "context"
  "facette.io/natsort"
  "github.com/alecthomas/kingpin"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "log"
  "os"
  "path/filepath"
  "regexp"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/util"
  "strings"
  "time"
)

func ConfigureCollectFromDirCommand(app *kingpin.Application, log *zap.Logger) {
  command := app.Command("collect", "Collect reports from idea.log files.")
  dirs := command.Flag("dir", "The input directory.").Short('i').Required().Strings()
  dbPath := command.Flag("db", "The output SQLite database file.").Short('o').Required().String()
  machine := command.Flag("machine", "The name of machine to associate report with.").Short('m').Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    return collectFromDirs(*dirs, *dbPath, *machine, log)
  })
}

func collectFromDirs(dirs []string, dbPath string, machine string, logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  reportAnalyzer, err := analyzer.CreateReportAnalyzer(dbPath, taskContext, logger)
  if err != nil {
    return err
  }

  go func() {
    err = <-reportAnalyzer.ErrorChannel
    cancel()

    if err != nil {
      // zap doesn't print with newlines
      log.Printf("%+v", err)
      logger.Error("cannot analyze", zap.Error(err))
    }
  }()

  defer util.Close(reportAnalyzer, logger)

  for _, dir := range dirs {
    logCollector := &LogCollector{
      reportAnalyzer:        reportAnalyzer,
      log:                   logger,
      productAndBuildInfoRe: regexp.MustCompile(`#([A-Z]{2})-([\d.]+)`),
      machine:               machine,
    }

    err = logCollector.collectFromDir(dir, taskContext)
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

func (t *LogCollector) collectFromDir(dir string, taskContext context.Context) error {
  dirInfo, err := os.Stat(dir)
  if err != nil {
    return errors.WithStack(err)
  }
  if !dirInfo.IsDir() {
    return errors.New("file " + dir + " is not a dir")
  }

  files, err := filepath.Glob(dir + "/idea*.log*")
  if err != nil {
    return errors.WithStack(err)
  }

  // product code and build are not specified in old report versions, so, it is inferred from log files.
  // but ide started log statement can be in another file (because log chunked across files), so, sort it and process from biggest to lesser (idea.log - latest, idea.8 - oldest).
  natsort.Sort(files)
  for i := len(files) - 1; i >= 0; i-- {
    if taskContext.Err() != nil {
      return nil
    }

    err := t.collectFromLogFile(files[i], taskContext)
    if err != nil {
      return err
    }
  }
  return nil
}

type LogCollector struct {
  reportAnalyzer *analyzer.ReportAnalyzer
  extraData      *analyzer.ExtraData

  productAndBuildInfoRe *regexp.Regexp

  log *zap.Logger

  machine string
}

func (t *LogCollector) collectFromLogFile(filePath string, taskContext context.Context) error {
  file, err := os.Open(filePath)
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(file, t.log)

  scanner := bufio.NewScanner(file)
  var jsonData bytes.Buffer
  state := 0

  startSuffix := []byte("=== Start: StartUp Measurement ===")
  endSuffix := []byte("=== Stop: StartUp Measurement ===")
  versionSlice := []byte("#com.intellij.idea.Main - IDE:")

  for scanner.Scan() {
    line := scanner.Bytes()
    if state == 1 {
      // idea start-up performance writer bug - end suffix has extra trailing space, so, HasSuffix cannot be used
      if bytes.HasPrefix(line, endSuffix) {
        if taskContext.Err() != nil {
          return nil
        }

        if t.extraData == nil {
          return errors.New("extraData not computed")
        }

        if len(t.extraData.ProductCode) == 0 {
          return errors.New("ProductCode not computed")
        }
        if len(t.extraData.BuildNumber) == 0 {
          return errors.New("BuildNumber not computed")
        }

        err = t.reportAnalyzer.Analyze(jsonData.Bytes(), *t.extraData)
        if err != nil {
          return err
        }

        state = 0
        t.extraData.LastGeneratedTime = -1
        jsonData.Reset()
      } else {
        jsonData.Write(line)
      }
    } else if bytes.Contains(line, versionSlice) {
      result := t.productAndBuildInfoRe.FindStringSubmatch(string(line))
      if result == nil {
        return errors.New("cannot find product and build number info")
      }

      t.extraData = &analyzer.ExtraData{
        ProductCode: result[1],
        BuildNumber: result[2],
        Machine:     t.machine,
      }
    } else if bytes.HasSuffix(line, startSuffix) {
      lineString := scanner.Text()
      // UTC, but it is ok, modern reports contain correct generated time
      parsedTime, err := time.Parse("2006-01-02 15:04:05", lineString[0:strings.IndexRune(lineString, ',')])
      if err != nil {
        return errors.WithStack(err)
      }

      state = 1
      if t.extraData == nil {
        return errors.New("cannot find product and build number info")
      }
      t.extraData.LastGeneratedTime = parsedTime.Unix()
    }
  }

  return nil
}
