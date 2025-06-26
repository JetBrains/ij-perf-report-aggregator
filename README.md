[![JetBrains incubator project](https://jb.gg/badges/incubator-flat-square.svg)](https://confluence.jetbrains.com/display/ALL/JetBrains+on+GitHub)
![](https://camo.githubusercontent.com/be6f8b50b2400e8b0dc74e58dd9a68803fe6698f5f30d843a7504888879f8392/68747470733a2f2f6a622e67672f6261646765732f696e63756261746f722d706c61737469632e737667)
[![lint-and-test](https://github.com/JetBrains/ij-perf-report-aggregator/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/JetBrains/ij-perf-report-aggregator/actions/workflows/lint-and-test.yml)

## IJ Perf

Tool to collect performance reports in various formats from TeamCity, insert into ClickHouse, send notifications about degradations and visualize the results.

- `clickhouse-backup` - backup clickhouse data.
- `clickhouse-restore` - restore clickhouse data.
- `backend` - server, provide access to ClickHouse and Postgres via HTTP API
- `tc-collector` - collect performance reports from [TeamCity artifacts](https://www.jetbrains.com/help/teamcity/build-artifact.html) and insert to ClickHouse.
- `degradation-detector` service that checks for performance degradations and sends notifications.

### Dashboard Editing

Directory `dashboard` contains Vue.js application built using [PrimeVue](https://primevue.org/) UI Library.

- `pnpm i` to install dependencies
- `pnpm vite serve` to start a dev server with hot module replacement.
- `pnpm vite serve --host` to start a dev server with hot module replacement on 0.0.0.0.
- `pnpm build` to build for production and `pnpm vite preview` to preview the production build.
- You can change server to the local one in `ServerWithCompressConfigurator.ts`
  - To test your changes against production server, replace the url to `https://ij-perf-api.labs.jb.gg`. Make sure you have access to the internal network.

To change dashboard, edit your dashboard page in `dashboard/new-dashboard/src/components`
