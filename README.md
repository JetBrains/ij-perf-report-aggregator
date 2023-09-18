![](https://camo.githubusercontent.com/be6f8b50b2400e8b0dc74e58dd9a68803fe6698f5f30d843a7504888879f8392/68747470733a2f2f6a622e67672f6261646765732f696e63756261746f722d706c61737469632e737667)

## IJ Perf

Tool to collect performance reports in JSON format from TeamCity, insert into ClickHouse and visualize.

- `clickhouse-backup` - backup clickhouse data.
- `clickhouse-restore` - restore clickhouse data.
- `report-aggregator` - stats server, provide access to ClickHouse via HTTP API tailored for analytics (inspired by [Cube.js](https://cube.dev/docs/query-format) query format).
- `tc-collector` - collect performance reports from [TeamCity artifacts](https://www.jetbrains.com/help/teamcity/build-artifact.html) and insert to ClickHouse.
- `transform` - transform existing data into another form. Raw JSON report is preserved as is, but for performance reasons maybe needed to pre-analyze and extract data into separate columns during collecting. And as data requirements changes, re-analyze the whole data set maybe required.

### Dashboard Editing

Directory `dashboard` contains Vue.js application built using [Element Plus](http://element-plus.org/) Desktop UI Library.

- `pnpm i` to install dependencies. [pnpm](https://pnpm.js.org/en/installation/) is recommended, do not use Yarn or NPM.
- Run `serve` Run Configuration (don't forget to restart in case you have new data in Database)
- `pnpm run dev` to start a dev server with hot module replacement.
- Change _Server_ to `http://localhost:9044`

To change dashboard, edit your dashboard page `*Dashboard.vue`, for example `IntelliJDashboard.vue` or `GenericMetricDashboard.vue`.

`LineChartCard` or `BarChartCard` supports `measures` property.
Specify desired metric. Multiple metrics are supported, but keep in mind that each metric means chart series and overuse can make chart unreadable.

If metric is extracted from report to field, just use it's field name.
Otherwise use:

- `activityCategory.activityName` to get duration value
- `activityCategory.activityName.s` to get start value.
- `activityCategory.activityName.e` to get end value.

For example, if activity `launch terminal` reported under category `prepareAppInitActivities`, use `prepareAppInitActivities.launch terminal` as metric name. Or `prepareAppInitActivities.first render.s` to get start value of `first render`.

### Adding a New Database

`cd ~/Documents/report-aggregator`

`clickhouse client -h 127.0.0.1` to use clickhouse client to perform SQL queries.

1. Set env `DB_NAME` to desired database name:
   ```shell
   export DB_NAME=
   ```
2. Create directory for your database in `db-schema` and copy `report.sql` from another database.
3. Execute SQL:
   ```shell
   clickhouse client -h 127.0.0.1 --query="create database $DB_NAME"
   clickhouse client -h 127.0.0.1 -d $DB_NAME --multiquery < ./db-schema/common/installer.sql
   clickhouse client -h 127.0.0.1 -d $DB_NAME --multiquery < ./db-schema/common/collector_state.sql
   clickhouse client -h 127.0.0.1 -d $DB_NAME --multiquery < "./db-schema/$DB_NAME/report.sql"
   ```

Addition docs in [Knowledge Base](https://youtrack.jetbrains.com/articles/IJPL-A-196/IJ-Perf-Development)