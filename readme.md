Tool to aggregate IntelliJ Platform performance reports.


1. Collect from `idea.log` (reported only for snapshot builds or in internal mode or if VM property `idea.log.perf.stats` is set to `true`).

    `report-aggregator collect --dir ~/Library/Logs/IntelliJIdea2019.3 --db ~/ij-perf-report-db/db.sqlite --machine "imac 2016"`
    
    Here `--machine` it is arbitrary id of machine where reports were generated (comparing reports make sense only on the same hardware). 
    
    See [Locating IDE log files](https://intellij-support.jetbrains.com/hc/en-us/articles/207241085-Locating-IDE-log-files).
    
2. Start server.
    
    `report-aggregator serve --db ~/ij-perf-report-db/db.sqlite`
    
3. Open [visualizer](https://ij-perf.jetbrains.com/#/aggregatedStats) (page `/#/aggregatedStats`), check `Server url` and click on `Load` to load and visualize data.


`aggregatedStats` visualizer page and this tool in an alpha stage. Aggregator throws error if anything is unclear to ensure that stats is representative.

Chart legend allows you to disable unneeded metrics to make chart more clear.
 
```
usage: report-aggregator [<flags>] <command> [<args> ...]

report-aggregator

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  collect --dir=DIR --db=DB --machine=MACHINE
    Collect reports from idea.log files.

  serve --db=DB
    Serve SQLite database.
```