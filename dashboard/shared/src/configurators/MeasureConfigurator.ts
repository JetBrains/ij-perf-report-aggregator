import { LineSeriesOption } from "echarts/charts"
import { Ref, shallowRef, watch } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { ChartConfigurator, ChartOptions } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryDimension, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery, toMutableArray } from "../dataQuery"
import { DebouncedTask, TaskHandle } from "../util/debounce"
import { loadJson } from "../util/httpUtil"
import { DimensionConfigurator } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"

// natural sort of alphanumerical strings
const collator = new Intl.Collator(undefined, {numeric: true, sensitivity: "base"})

export class MeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  public readonly data = shallowRef<Array<string>>([])
  public readonly value = shallowRef<Array<string>>([])

  private readonly debouncedLoadMetadata = new DebouncedTask(taskHandle => this.loadMetadata(taskHandle))

  constructor(private readonly serverConfigurator: ServerConfigurator,
              persistentStateManager: PersistentStateManager,
              private readonly parent: DimensionConfigurator | null = null,
              readonly skipZeroValues: boolean = true) {
    persistentStateManager.add("measure", this.value)

    if (this.parent != null) {
      watch(this.parent.value, this.debouncedLoadMetadata.executeFunctionReference)
    }
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

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    return configureChart(configuration, data)
  }

  scheduleLoadMetadata(immediately: boolean): void {
    this.debouncedLoadMetadata.execute(immediately)
  }

  loadMetadata(taskHandle: TaskHandle): Promise<unknown> {
    if (this.serverConfigurator.databaseName === "ij") {
      const server = this.serverConfigurator.value.value
      if (server == null || server.length === 0) {
        return Promise.resolve()
      }

      return loadJson<Array<string>>(`${server}/api/v1/meta/measure?db=${this.serverConfigurator.databaseName}`, null, taskHandle, data => {
        this.data.value = data
      })
    }

    const parent = this.parent
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    if (!this.serverConfigurator.configureQuery(query, configuration)) {
      return Promise.resolve()
    }
    if (parent != null && !parent.configureQuery(query, configuration)) {
      return Promise.resolve()
    }

    // "group by" is equivalent of distinct (https://clickhouse.tech/docs/en/sql-reference/statements/select/distinct/#alternatives)
    query.addDimension({name: "measures", subName: "name"})
    query.order = ["measures.name"]
    query.table = "report"
    query.flat = true

    return loadJson<Array<string>>(`${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`, null, taskHandle, data => {
      this.data.value = data
    })
  }
}

export class PredefinedMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  private readonly measures: Array<string>

  constructor(measures: Array<string>, public skipZeroValues: Ref<boolean> = shallowRef(true)) {
    // analyzer replaces space to underscore, but configurator supports specifying original metric names
    // render window.end => render_window.end
    this.measures = measures.map(it => it.replaceAll(" ", "_"))
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configureQuery(this.measures, query, configuration, this.skipZeroValues.value)
    configuration.chartConfigurator = this
    configuration.measures = this.measures
    return true
  }

  configureChart(data: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    return configureChart(configuration, data)
  }
}

export function measureNameToLabel(key: string): string {
  const metricPathEndDotIndex = key.indexOf(".")
  if (metricPathEndDotIndex == -1) {
    return key
  }
  else {
    return key
      .substring(0, metricPathEndDotIndex)
      .replaceAll("_",  " ")
  }
}

function configureQuery(measureNames: Array<string>, query: DataQuery, configuration: DataQueryExecutorConfiguration, skipZeroValues: boolean): void {
  // stable order of series (UI) and fields in query (caching)
  measureNames.sort((a, b) => collator.compare(a, b))

  query.addField({
    name: "t",
    sql: "toUnixTimestamp(generated_time) * 1000",
  })

  // we cannot request several measures in one SQL query - for each measure separate SQl query with filter by measure name
  if (query.db === "ij") {
    const field: DataQueryDimension = {name: measureNames[0]}
    const filter: DataQueryFilter | null = skipZeroValues ? {field: measureNames[0], operator: "!=", value: 0} : null
    if (measureNames.length > 1) {
      configureQueryProducer(configuration, field, filter, measureNames)
    }

    query.addField(field)
    if (filter != null) {
      query.addFilter(filter)
    }
  }
  else {
    const pureNames = measureNames.map(it => it.endsWith(".end") ? it.substring(0, it.length - ".end".length) : it)
    const filter: DataQueryFilter = {field: "measures.name", value: pureNames[0]}
    if (measureNames.length > 1) {
      configureQueryProducer(configuration, null, filter, pureNames)
    }

    if (measureNames.some(it => it.endsWith(".end"))) {
      query.addField({name: "measures", subName: "end", sql: "(measures.start + measures.value)"})
    }
    else {
      query.addField({name: "measures", subName: "value"})
    }

    query.addFilter(filter)

    if (skipZeroValues) {
      // for end we also filter by raw value and not by sum of start + duration (that stored under "value" name)
      query.addFilter({field: "measures.value", operator: "!=", value: 0})
    }
  }

  if (query.order != null) {
    throw new Error("order must be configured only by MetricLoader")
  }
  query.order = ["t"]
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, field: DataQueryDimension | null, filter: DataQueryFilter | null, values: Array<string>): void {
  let index = 1
  if (configuration.extraQueryProducer != null) {
    throw new Error("extraQueryMutator is already set")
  }

  configuration.extraQueryProducer = {
    mutate() {
      if (field != null) {
        field.name = values[index]
        if (filter != null) {
          filter.field = field.name
        }
      }
      else if (filter != null) {
        filter.value = values[index]
      }
      index++
      return index !== values.length
    },
    getDataSetLabel(index: number): string {
      return values[index].replaceAll("_", " ")
    },
    getMeasureName(index: number): string {
      return values[index]
    }
  }
}

function configureChart(configuration: DataQueryExecutorConfiguration, dataList: DataQueryResult): ChartOptions {
  const series = new Array<LineSeriesOption>()
  const extraQueryProducer = configuration.extraQueryProducer
  for (let dataIndex = 0; dataIndex < dataList.length; dataIndex++) {
    const measureName = extraQueryProducer == null ? configuration.measures[0] : extraQueryProducer.getMeasureName(dataIndex)

    const label = extraQueryProducer == null ? null : extraQueryProducer.getDataSetLabel(dataIndex)
    let seriesName: string
    if (label === null) {
      seriesName = measureNameToLabel(measureName)
    }
    else if (measureName.length === 1) {
      seriesName = label
    }
    else {
      seriesName = `${label}-${measureNameToLabel(measureName)}`
    }
    series.push({
      name: seriesName,
      type: "line",
      smooth: true,
      showSymbol: false,
      legendHoverLink: true,
      sampling: "lttb",
      dimensions: [{name: "time", type: "time"}, {name: seriesName, type: "int"}],
      data: dataList[dataIndex],
    })
  }

  return {
    series,
  }
}