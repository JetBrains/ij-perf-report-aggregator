package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
  "github.com/alecthomas/kingpin"
  "log/slog"
  "os"
)

func main() {
  var app = kingpin.New("report-aggregator", "report-aggregator").Version("0.0.1")

  ConfigureServeCommand(app)

  _, err := app.Parse(os.Args[1:])
  if err != nil {
    slog.Error("Cannon parse command line arguments", "err", err)
    os.Exit(1)
  }
}

func ConfigureServeCommand(app *kingpin.Application) {
  command := app.Command("serve", "Start aggregated stats server.")
  dbUrl := command.Flag("db", "The ClickHouse database URL.").Required().String()
  command.Action(func(context *kingpin.ParseContext) error {
    err := server.Serve(*dbUrl, "")
    if err != nil {
      return err
    }

    return nil
  })
}
