import { LineSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { deepEqual } from "fast-equals"
import { combineLatest, debounceTime, distinctUntilChanged, forkJoin, map, Observable, of, switchMap } from "rxjs"
import { Ref, shallowRef } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { ChartConfigurator, collator } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryDimension, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery, toMutableArray } from "../dataQuery"
import { LineChartOptions } from "../echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, numberAxisLabelFormatter, numberFormat } from "../formatter"
import { DimensionConfigurator } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  public readonly data = shallowRef<Array<string>>([])
  private readonly _value = shallowRef<Array<string> | string | null>(null)

  createObservable(): Observable<unknown> {
    return refToObservable(this.value, true)
  }

  public get value() {
    if (typeof this._value.value == "string") {
      this._value.value = [this._value.value]
      this.persistentStateManager.add("measure", this._value)
    }
    return this._value as Ref<string[] | null>
  }

  constructor(serverConfigurator: ServerConfigurator,
              private readonly persistentStateManager: PersistentStateManager,
              filters: Array<DataQueryConfigurator> = [],
              readonly skipZeroValues: boolean = true) {
    this.persistentStateManager.add("measure", this.value);

    (filters.length === 0 ? serverConfigurator.createObservable() : combineLatest(filters.map(it => it.createObservable()).concat(serverConfigurator.createObservable())))
      .pipe(
        debounceTime(100),
        distinctUntilChanged(deepEqual),
        switchMap(() => {
          const isIj = serverConfigurator.databaseName === "ij"

          const loadMeasureListUrl = getLoadMeasureListUrl(isIj ? "measure" : "measures", serverConfigurator, filters)
          if (loadMeasureListUrl == null) {
            return of(null)
          }

          if (isIj) {
            const server = serverConfigurator.value.value
            if (server == null || server.length === 0) {
              return of(null)
            }

            return forkJoin([
              fromFetchWithRetryAndErrorHandling<Array<string>>(`${server}/api/v1/meta/measure?db=${serverConfigurator.databaseName}`),
              fromFetchWithRetryAndErrorHandling<Array<string>>(loadMeasureListUrl),
            ])
              .pipe(
                map(data => {
                  return data.flat(1)
                }),
              )
          }
          else {
            return fromFetchWithRetryAndErrorHandling<Array<string>>(loadMeasureListUrl)
          }
        }),
      )
      .subscribe(data => {
        if (data != null) {
          this.data.value = data
          if (Array.isArray(this._value.value)) {
            this._value.value = this._value.value.filter(it => {
              return this.data.value.indexOf(it) > -1
            })
          }
          else if (typeof this._value.value == "string") {
            this._value.value = this.data.value.indexOf(this._value.value) > 0 ? this._value.value : null
          }
        }
      })
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const measureNames = toMutableArray(this.value.value)
    if (measureNames.length === 0) {
      return false
    }

    configureQuery(measureNames, query, configuration, this.skipZeroValues)
    configuration.measures = measureNames
    configuration.chartConfigurator = this
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): LineChartOptions {
    return configureChart(configuration, data)
  }
}

function getLoadMeasureListUrl(structureName: string,
                               serverConfigurator: ServerConfigurator,
                               filters: Array<DataQueryConfigurator>): string | null {
  const query = new DataQuery()
  const configuration = new DataQueryExecutorConfiguration()
  if (!serverConfigurator.configureQuery(query, configuration)) {
    return null
  }

  if (filters.length !== 0) {
    for (const parent of filters) {
      if (parent instanceof DimensionConfigurator) {
        const value = parent.value.value
        if (value == null || value.length === 0) {
          return null
        }

        query.addFilter({f: parent.name, v: value})
      }
      else {
        // well, so, it is a filter configurator (e.g. MachineConfigurator)
        if (!parent.configureQuery(query, configuration)) {
          return null
        }
      }
    }
  }

  // "group by" is equivalent of distinct (https://clickhouse.tech/docs/en/sql-reference/statements/select/distinct/#alternatives)
  query.addDimension({n: structureName, subName: "name"})
  query.order = [`${structureName}.name`]
  query.table = "report"
  query.flat = true
  return `${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`
}

export class PredefinedMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(private readonly measures: Array<string>, readonly skipZeroValues: Ref<boolean> = shallowRef(true)) {
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.skipZeroValues)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configureQuery(this.measures, query, configuration, this.skipZeroValues.value)
    configuration.chartConfigurator = this
    configuration.measures = this.measures
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): LineChartOptions {
    return configureChart(configuration, data)
  }
}

export function measureNameToLabel(key: string): string {
  const metricPathEndDotIndex = key.indexOf(".")
  if (metricPathEndDotIndex == -1) {
    // remove _d or _i suffix
    return key.replace(/_[a-z]$/g, "")
  }
  else {
    return key.substring(0, metricPathEndDotIndex)
  }
}

function configureMeasureInANewFormat(measureNames: Array<string>,
                                      configuration: DataQueryExecutorConfiguration,
                                      query: DataQuery,
                                      structureName: string,
                                      valueName: string,
                                      skipZeroValues: boolean) {
  const pureNames = measureNames.map(it => it.endsWith(".end") ? it.substring(0, it.length - ".end".length) : it)
  const filter: DataQueryFilter = {f: `${structureName}.name`, v: pureNames[0]}
  if (measureNames.length > 1) {
    configureQueryProducer(configuration, null, filter, pureNames)
  }

  if (measureNames.some(it => it.endsWith(".end"))) {
    query.insertField({n: structureName, subName: "end", sql: `(${structureName}.start + ${structureName}.${valueName})`}, 1)
  }
  else {
    query.insertField({n: structureName, subName: valueName}, 1)
  }

  query.addFilter(filter)

  if (skipZeroValues) {
    // for end we also filter by raw value and not by sum of start + duration (that stored under "value" name)
    query.addFilter({f: `${structureName}.${valueName}`, o: "!=", v: 0})
  }
}

function configureQuery(measureNames: Array<string>, query: DataQuery, configuration: DataQueryExecutorConfiguration, skipZeroValues: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  query.insertField({
    n: "t",
    sql: "toUnixTimestamp(generated_time) * 1000",
  }, 0)

  // we cannot request several measures in one SQL query - for each measure separate SQl query with filter by measure name
  if (query.db === "ij") {
    if (measureNames[0].includes(" ")) {
      configureMeasureInANewFormat(measureNames, configuration, query, "measure", "duration", skipZeroValues)
    }
    else {
      const field: DataQueryDimension = {n: measureNames[0]}
      const filter: DataQueryFilter | null = skipZeroValues ? {f: measureNames[0], o: "!=", v: 0} : null
      if (measureNames.length > 1) {
        configureQueryProducer(configuration, field, filter, measureNames)
      }

      query.insertField(field, 1)
      if (filter != null) {
        query.addFilter(filter)
      }
    }
  }
  else {
    configureMeasureInANewFormat(measureNames, configuration, query, "measures", "value", skipZeroValues)
  }

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = ["t"]
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, field: DataQueryDimension | null, filter: DataQueryFilter | null, values: Array<string>): void {
  configuration.extraQueryProducers.push({
    size(): number {
      return values.length
    },

    mutate(index: number): void {
      if (field != null) {
        field.n = values[index]
        if (filter != null) {
          filter.f = field.n
        }
      }
      else if (filter != null) {
        filter.v = values[index]
      }
    },
    getSeriesName(index: number): string {
      return measureNameToLabel(values[index])
    },
    getMeasureName(index: number): string {
      return values[index]
    }
  })
}

function configureChart(configuration: DataQueryExecutorConfiguration, dataList: DataQueryResult): LineChartOptions {
  const series = new Array<LineSeriesOption>()
  let useDurationFormatter = true

  const dataset: Array<DatasetOption> = []

  for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
    const measureName = configuration.measureNames[dataIndex]
    const seriesName = configuration.seriesNames[dataIndex]
    const seriesData = dataList[dataIndex]
    const symbolSize = Math.min(800 / seriesData[0].length, 9)
    series.push({
      // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
      id: measureName === seriesName ? seriesName : `${measureName}@${seriesName}`,
      name: seriesName,
      type: "line",
      smooth: false,
      showSymbol: seriesData[0].length < 100,
      symbolSize,
      symbol: "circle",
      legendHoverLink: true,
      sampling: "lttb",
      seriesLayoutBy: "row",
      datasetIndex: dataIndex,
      dimensions: [{name: "time", type: "time"}, {name: seriesName, type: "int"}],
    })

    if (useDurationFormatter && !isDurationFormatterApplicable(measureName)) {
      useDurationFormatter = false
    }

    dataset.push({
      source: seriesData,
      sourceHeader: false,
    })
  }

  return {
    dataset,
    yAxis: {
      axisLabel: {
        formatter: useDurationFormatter ? durationAxisPointerFormatter : numberAxisLabelFormatter,
      },
      axisPointer: {
        label: {
          formatter(data): string {
            const value = data["value"] as number
            return useDurationFormatter ? durationAxisPointerFormatter(value) : numberFormat.format(value)
          },
        },
      },
    },
    series,
  }
}